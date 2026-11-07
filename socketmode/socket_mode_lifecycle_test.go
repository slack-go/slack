package socketmode

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"testing/synctest"

	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/require"

	"github.com/slack-go/slack"
)

type socketModeHTTPClientFunc func(*http.Request) (*http.Response, error)

func (f socketModeHTTPClientFunc) Do(r *http.Request) (*http.Response, error) { return f(r) }

type delayedReadConn struct {
	net.Conn
	failed  chan struct{}
	release <-chan struct{}
	once    sync.Once
}

func (c *delayedReadConn) Read(p []byte) (int, error) {
	n, err := c.Conn.Read(p)
	if err != nil {
		c.once.Do(func() { close(c.failed) })
		<-c.release
	}
	return n, err
}

func TestRunContextWaitsForMessageReceiver(t *testing.T) {
	stopServer := make(chan struct{})
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		upgrader := websocket.Upgrader{CheckOrigin: func(*http.Request) bool { return true }}
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			t.Errorf("upgrade local WebSocket: %v", err)
			return
		}
		defer conn.Close()
		<-stopServer
	}))
	defer server.Close()
	defer close(stopServer)
	wsURL := "ws" + strings.TrimPrefix(server.URL, "http")

	synctest.Test(t, func(t *testing.T) {
		failed, release := make(chan struct{}), make(chan struct{})
		var releaseOnce sync.Once
		defer releaseOnce.Do(func() { close(release) })
		transport := socketModeHTTPClientFunc(func(*http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusOK,
				Header:     make(http.Header),
				Body:       io.NopCloser(strings.NewReader(fmt.Sprintf(`{"ok":true,"url":%q}`, wsURL))),
			}, nil
		})
		dialer := &websocket.Dialer{NetDialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
			conn, err := (&net.Dialer{}).DialContext(ctx, network, addr)
			if err != nil {
				return nil, err
			}
			return &delayedReadConn{Conn: conn, failed: failed, release: release}, nil
		}}
		client := New(slack.New("token", slack.OptionHTTPClient(transport)), OptionDialer(dialer))
		ctx, cancel := context.WithCancel(t.Context())
		defer cancel()
		var runError error
		done := make(chan struct{})
		go func() {
			runError = client.RunContext(ctx)
			close(done)
		}()
		for event := range client.Events {
			require.NotEqual(t, EventTypeConnectionError, event.Type, "local connection failed: %v", event.Data)
			if event.Type == EventTypeConnected {
				break
			}
		}
		cancel()
		<-failed
		synctest.Wait()
		select {
		case <-done:
			t.Error("RunContext returned before its event-producing receiver exited")
		default:
		}
		releaseOnce.Do(func() { close(release) })
		<-done
		require.ErrorIs(t, runError, context.Canceled)
	})
}
