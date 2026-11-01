package slack

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/require"
)

func TestRTMGetInfoZeroValue(t *testing.T) {
	var rtm RTM
	if info := rtm.GetInfo(); info != nil {
		t.Fatalf("GetInfo() = %v, want nil", info)
	}
}

func TestRTMGetInfoConcurrent(t *testing.T) {
	rtm := New("test-token").NewRTM()
	defer rtm.pingDeadman.Stop()
	var wg sync.WaitGroup
	wg.Go(func() {
		for range 1000 {
			// The connection manager uses this lock when replacing the info.
			rtm.mu.Lock()
			rtm.info = &Info{Team: &Team{ID: "T1"}}
			rtm.mu.Unlock()
		}
	})
	for range 1000 {
		_ = rtm.GetInfo()
	}
	wg.Wait()
}

func TestRTMIncomingEventsTruncatedStream(t *testing.T) {
	upgrader := websocket.Upgrader{}
	truncate := make(chan struct{})
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			t.Error(err)
			return
		}
		defer conn.Close()
		<-truncate
		// End the TCP write side partway through a TLS record header, while
		// leaving the read side open to accept the client's ping messages.
		tcp := conn.UnderlyingConn().(*tls.Conn).NetConn().(*net.TCPConn)
		if _, err := tcp.Write([]byte{0x17, 0x03}); err != nil {
			t.Error(err)
			return
		}
		if err := tcp.CloseWrite(); err != nil {
			t.Error(err)
			return
		}
		// Keep accepting pings: a successful write does not establish that the
		// client's read side is healthy.
		for {
			if _, _, err := conn.ReadMessage(); err != nil {
				return
			}
		}
	}))
	t.Cleanup(server.Close)

	dialer := websocket.Dialer{
		TLSClientConfig: server.Client().Transport.(*http.Transport).TLSClientConfig.Clone(),
	}
	conn, _, err := dialer.Dial("wss"+strings.TrimPrefix(server.URL, "https"), nil)
	close(truncate)
	require.NoError(t, err)
	t.Cleanup(func() { _ = conn.Close() })

	rtm := New("test-token").NewRTM()
	t.Cleanup(func() { rtm.pingDeadman.Stop() })
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	pingDone := make(chan struct{})
	go func() {
		defer close(pingDone)
		for {
			select {
			case <-rtm.forcePing:
				if err := conn.WriteJSON(Ping{Type: "ping"}); err != nil {
					return
				}
			case <-ctx.Done():
				return
			}
		}
	}()
	t.Cleanup(func() {
		cancel()
		_ = conn.Close()
		<-pingDone
	})

	require.NotPanics(t, func() {
		err = rtm.handleIncomingEvents(ctx, conn, make(chan json.RawMessage, 1))
	})
	require.ErrorIs(t, err, io.ErrUnexpectedEOF)
	require.Len(t, rtm.IncomingEvents, 1)
	event := <-rtm.IncomingEvents
	require.Equal(t, "incoming_error", event.Type)
	incomingErr, ok := event.Data.(*IncomingEventError)
	require.True(t, ok)
	require.ErrorIs(t, incomingErr.ErrorObj, io.ErrUnexpectedEOF)
}

func TestRTMIncomingEventsEmptyFrames(t *testing.T) {
	const hello = `{"type":"hello"}`
	const message = `{"type":"message","text":"still connected"}`
	repeatedFrames := append(make([][]byte, 1024), []byte(hello))

	for _, tt := range []struct {
		name     string
		frames   [][]byte
		messages []string
	}{
		{name: "isolated empty frame", frames: [][]byte{nil, []byte(hello)}, messages: []string{hello}},
		{name: "repeated empty frames", frames: repeatedFrames, messages: []string{hello}},
		{name: "empty and incomplete frames between messages", frames: [][]byte{nil, []byte(hello), []byte(" "), []byte("{"), []byte(message)}, messages: []string{hello, message}},
	} {
		t.Run(tt.name, func(t *testing.T) {
			upgrader := websocket.Upgrader{}
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				conn, err := upgrader.Upgrade(w, r, nil)
				if err != nil {
					t.Error(err)
					return
				}
				defer conn.Close()
				for _, frame := range tt.frames {
					if err := conn.WriteMessage(websocket.TextMessage, frame); err != nil {
						return
					}
				}
				_ = conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, "done"))
			}))
			t.Cleanup(server.Close)
			conn, _, err := websocket.DefaultDialer.Dial("ws"+strings.TrimPrefix(server.URL, "http"), nil)
			require.NoError(t, err)
			t.Cleanup(func() { _ = conn.Close() })
			require.NoError(t, conn.SetReadDeadline(time.Now().Add(5*time.Second)))

			rtm := New("test-token").NewRTM()
			t.Cleanup(func() { rtm.pingDeadman.Stop() })
			// Buffer ping requests so the reader can run without a manager.
			rtm.forcePing = make(chan bool, len(tt.frames))
			events := make(chan json.RawMessage, len(tt.frames))
			err = rtm.handleIncomingEvents(context.Background(), conn, events)
			var closeErr *websocket.CloseError
			require.ErrorAs(t, err, &closeErr)
			require.Equal(t, websocket.CloseNormalClosure, closeErr.Code)
			require.Len(t, events, len(tt.messages))
			for _, message := range tt.messages {
				require.JSONEq(t, message, string(<-events))
			}
		})
	}
}
