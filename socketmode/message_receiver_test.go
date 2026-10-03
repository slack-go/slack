package socketmode

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"testing"

	"github.com/gorilla/websocket"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slacktest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// sendFrames returns a server that writes frames in order, then closes the WebSocket.
func sendFrames(frames ...[]byte) func(conn *websocket.Conn) {
	return func(conn *websocket.Conn) {
		for _, frame := range frames {
			if err := conn.WriteMessage(websocket.TextMessage, frame); err != nil {
				return
			}
		}
		_ = conn.WriteMessage(
			websocket.CloseMessage,
			websocket.FormatCloseMessage(websocket.CloseNormalClosure, "done"),
		)
	}
}

func repeatFrame(frame []byte, n int) [][]byte {
	frames := make([][]byte, n)
	for i := range frames {
		frames[i] = frame
	}
	return frames
}

// receive runs the receive loop in the calling goroutine and turns a panic into an
// error, so a regression fails the test instead of crashing the test binary.
func receive(conn *websocket.Conn, sink chan json.RawMessage) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("receive loop panicked: %v", r)
		}
	}()
	return New(slack.New("xoxb-test-token")).runMessageReceiver(context.Background(), conn, sink)
}

// truncatedConn fails every read once truncated is set, as crypto/tls does when the
// stream ends in the middle of a record. gorilla maps only io.EOF to a CloseError, so
// this io.ErrUnexpectedEOF becomes a permanent read error.
type truncatedConn struct {
	net.Conn
	truncated bool
}

func (c *truncatedConn) Read(p []byte) (int, error) {
	if c.truncated {
		return 0, io.ErrUnexpectedEOF
	}
	return c.Conn.Read(p)
}

// TestRunMessageReceiverGivesUpOnTruncatedStream is the regression test for the crash:
// the loop took the permanent error for a stray frame and read again until gorilla
// panicked with "repeated read on failed websocket connection".
func TestRunMessageReceiverGivesUpOnTruncatedStream(t *testing.T) {
	srv := slacktest.NewTestServer(func(c slacktest.Customize) {
		c.Handle("/ws", slacktest.Websocket(func(*websocket.Conn) {}))
	})
	srv.Start()
	t.Cleanup(srv.Stop)

	var tc *truncatedConn
	dialer := websocket.Dialer{
		NetDialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
			c, err := (&net.Dialer{}).DialContext(ctx, network, addr)
			if err != nil {
				return nil, err
			}
			tc = &truncatedConn{Conn: c}
			return tc, nil
		},
	}
	conn, _, err := dialer.Dial(srv.GetWSURL(), nil)
	require.NoError(t, err)
	t.Cleanup(func() { _ = conn.Close() })
	tc.truncated = true

	err = receive(conn, make(chan json.RawMessage, 1))

	require.Error(t, err)
	assert.NotContains(t, err.Error(), "panicked")
	assert.ErrorIs(t, err, io.ErrUnexpectedEOF)
}

func TestRunMessageReceiverGivesUpOnUnusableFrames(t *testing.T) {
	// Frames with no JSON value make ReadJSON return io.ErrUnexpectedEOF.
	conn := dialTestWebSocket(t, sendFrames(repeatFrame(nil, maxConsecutiveIgnoredReads)...))

	err := receive(conn, make(chan json.RawMessage, 1))

	require.Error(t, err)
	assert.ErrorIs(t, err, io.ErrUnexpectedEOF)
	assert.Contains(t, err.Error(), fmt.Sprintf("giving up after %d consecutive", maxConsecutiveIgnoredReads))
}

func TestRunMessageReceiverGivesUpOnMalformedJSON(t *testing.T) {
	conn := dialTestWebSocket(t, sendFrames(repeatFrame([]byte("}not json{"), maxConsecutiveIgnoredReads)...))

	err := receive(conn, make(chan json.RawMessage, 1))

	require.Error(t, err)
	assert.Contains(t, err.Error(), fmt.Sprintf("giving up after %d consecutive", maxConsecutiveIgnoredReads))
	var syntaxErr *json.SyntaxError
	assert.ErrorAs(t, err, &syntaxErr)
}

// TestRunMessageReceiverKeepsGoingAfterRecoverableFrame checks that the count resets on
// the next good frame. Without the reset, the unusable frames below reach the limit.
func TestRunMessageReceiverKeepsGoingAfterRecoverableFrame(t *testing.T) {
	unusable := repeatFrame(nil, maxConsecutiveIgnoredReads-1)
	var frames [][]byte
	frames = append(frames, unusable...)
	frames = append(frames, []byte(`{"type":"hello"}`))
	frames = append(frames, unusable...)
	frames = append(frames, []byte(`{"type":"disconnect"}`))
	conn := dialTestWebSocket(t, sendFrames(frames...))

	sink := make(chan json.RawMessage, 2)
	err := receive(conn, sink)

	closeErr, ok := errors.AsType[*websocket.CloseError](err)
	require.True(t, ok, "the loop must end on the close frame, got: %v", err)
	assert.Equal(t, websocket.CloseNormalClosure, closeErr.Code)
	require.Len(t, sink, 2)
	assert.JSONEq(t, `{"type":"hello"}`, string(<-sink))
	assert.JSONEq(t, `{"type":"disconnect"}`, string(<-sink))
}

func TestIgnoredReadErrorUnwraps(t *testing.T) {
	err := error(ignoredReadError{io.ErrUnexpectedEOF})

	assert.ErrorIs(t, err, io.ErrUnexpectedEOF)

	ignorable, ok := errors.AsType[ignoredReadError](err)
	require.True(t, ok)
	assert.Equal(t, io.ErrUnexpectedEOF, ignorable.err)
}
