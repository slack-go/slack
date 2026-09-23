package socketmode

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/websocket"
	"github.com/slack-go/slack"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// dialFrameServer opens a real WebSocket connection to a server that keeps sending the
// frames produced by send, so tests can drive the receive loop with real gorilla reads.
func dialFrameServer(t *testing.T, send func(conn *websocket.Conn) error) *websocket.Conn {
	t.Helper()

	upgrader := websocket.Upgrader{}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}
		defer conn.Close()
		for {
			if err := send(conn); err != nil {
				return
			}
		}
	}))
	t.Cleanup(srv.Close)

	conn, _, err := websocket.DefaultDialer.Dial("ws"+strings.TrimPrefix(srv.URL, "http"), nil)
	require.NoError(t, err)
	t.Cleanup(func() { _ = conn.Close() })

	return conn
}

// runReceiver runs the receive loop and returns its error, failing the test if it either
// panics or never returns.
func runReceiver(t *testing.T, conn *websocket.Conn) error {
	t.Helper()

	smc := New(slack.New("xoxb-test-token"))
	// The receive loop reports unusable frames on Events; drain it so nothing blocks.
	go func() {
		for range smc.Events { //nolint:revive // draining
		}
	}()

	done := make(chan error, 1)
	go func() {
		defer func() {
			if r := recover(); r != nil {
				done <- errors.New("receive loop panicked")
			}
		}()
		done <- smc.runMessageReceiver(context.Background(), conn, make(chan json.RawMessage, 1))
	}()

	select {
	case err := <-done:
		return err
	case <-time.After(10 * time.Second):
		t.Fatal("receive loop never returned: it is spinning on a connection that will never deliver a message")
		return nil
	}
}

// TestRunMessageReceiverGivesUpOnUnusableFrames is the regression test for the busy loop.
//
// A frame that carries no JSON value leaves ReadJSON returning io.ErrUnexpectedEOF, which
// the receive loop deliberately tolerates. Tolerating it *unconditionally* means a peer
// that only ever sends such frames keeps the client pinned to a connection that will never
// deliver a message — the loop spins forever and never reconnects.
func TestRunMessageReceiverGivesUpOnUnusableFrames(t *testing.T) {
	conn := dialFrameServer(t, func(conn *websocket.Conn) error {
		return conn.WriteMessage(websocket.TextMessage, nil) // empty frame: no JSON value
	})

	err := runReceiver(t, conn)

	require.Error(t, err, "receiver must give up instead of spinning")
	assert.NotContains(t, err.Error(), "panicked")
	assert.ErrorIs(t, err, io.ErrUnexpectedEOF, "the cause must stay visible to the caller")
	assert.Contains(t, err.Error(), "giving up after 10 consecutive")
}

// TestRunMessageReceiverGivesUpOnMalformedJSON is the same guarantee for the other
// tolerated error: a peer that only ever sends malformed JSON.
func TestRunMessageReceiverGivesUpOnMalformedJSON(t *testing.T) {
	conn := dialFrameServer(t, func(conn *websocket.Conn) error {
		return conn.WriteMessage(websocket.TextMessage, []byte("}not json{"))
	})

	err := runReceiver(t, conn)

	require.Error(t, err, "receiver must give up instead of spinning")
	assert.Contains(t, err.Error(), "giving up after 10 consecutive")

	var syntaxErr *json.SyntaxError
	assert.ErrorAs(t, err, &syntaxErr, "the cause must stay visible to the caller")
}

// TestRunMessageReceiverKeepsGoingAfterRecoverableFrame guards the other half of the
// contract: a *stray* unusable frame must not cost us the connection, so the counter has
// to reset on the next frame that does carry a value.
func TestRunMessageReceiverKeepsGoingAfterRecoverableFrame(t *testing.T) {
	// Alternate unusable and usable frames, forever. With a limit but no reset, this
	// connection would be dropped after ten frames; it must survive instead.
	var n int
	conn := dialFrameServer(t, func(conn *websocket.Conn) error {
		n++
		if n%2 == 1 {
			return conn.WriteMessage(websocket.TextMessage, nil)
		}
		return conn.WriteMessage(websocket.TextMessage, []byte(`{"type":"hello"}`))
	})

	smc := New(slack.New("xoxb-test-token"))
	go func() {
		for range smc.Events { //nolint:revive // draining
		}
	}()

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	sink := make(chan json.RawMessage, 1)
	go func() {
		for range sink { //nolint:revive // draining
		}
	}()

	done := make(chan error, 1)
	go func() { done <- smc.runMessageReceiver(ctx, conn, sink) }()

	select {
	case err := <-done:
		t.Fatalf("receiver dropped a healthy connection over stray unusable frames: %v", err)
	case <-ctx.Done():
		// Still reading after many alternating frames, which is what we want.
	}
	assert.Greater(t, n, maxConsecutiveIgnoredReads, "test did not send enough frames to be meaningful")
}

// TestIgnoredReadErrorUnwraps guards the contract runMessageReceiver relies on: the
// sentinel must stay transparent to errors.Is/As so callers still see the real cause.
func TestIgnoredReadErrorUnwraps(t *testing.T) {
	err := error(ignoredReadError{io.ErrUnexpectedEOF})

	assert.ErrorIs(t, err, io.ErrUnexpectedEOF)

	ignorable, ok := errors.AsType[ignoredReadError](err)
	require.True(t, ok)
	assert.Equal(t, io.ErrUnexpectedEOF, ignorable.err)
}
