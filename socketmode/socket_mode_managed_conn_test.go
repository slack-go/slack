//go:build go1.13
// +build go1.13

package socketmode

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/websocket"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slacktest"

	"github.com/stretchr/testify/assert"
)

func Test_passContext(t *testing.T) {
	s := slacktest.NewTestServer()
	go s.Start()

	api := slack.New("ABCDEFG", slack.OptionAPIURL(s.GetAPIURL()))
	cli := New(api)

	ctx, cancel := context.WithTimeout(context.TODO(), time.Nanosecond)
	defer cancel()

	t.Run("RunWithContext", func(t *testing.T) {
		// should fail imidiatly.
		assert.EqualError(t, cli.RunContext(ctx), context.DeadlineExceeded.Error())
	})

	t.Run("openAndDial", func(t *testing.T) {
		_, _, err := cli.openAndDial(ctx, func(_ string) error { return nil })

		// should fail imidiatly.
		assert.EqualError(t, errors.Unwrap(err), context.DeadlineExceeded.Error())
	})

	t.Run("OpenWithContext", func(t *testing.T) {
		_, _, err := cli.OpenContext(ctx)

		// should fail imidiatly.
		assert.EqualError(t, errors.Unwrap(err), context.DeadlineExceeded.Error())
	})
}

func TestSendCtx_PayloadQueued(t *testing.T) {
	api := slack.New("ABCDEFG")
	cli := New(api)

	err := cli.SendCtx(context.Background(), Response{
		EnvelopeID: "test-envelope",
		Payload:    "small payload",
	})

	assert.NoError(t, err)

	select {
	case res := <-cli.socketModeResponses:
		assert.Equal(t, "test-envelope", res.EnvelopeID)
	default:
		t.Fatal("expected response to be queued on socketModeResponses channel")
	}
}

func TestSendCtx_PayloadAtLimit(t *testing.T) {
	api := slack.New("ABCDEFG")
	cli := New(api)

	// Build a response that serializes to exactly maxResponseSize bytes.
	// The envelope and JSON overhead consume some bytes, so we pad the payload.
	envelope := "test-envelope"
	overhead, _ := json.Marshal(Response{EnvelopeID: envelope, Payload: ""})
	padding := strings.Repeat("x", maxResponseSize-len(overhead))

	// Verify we hit exactly the limit.
	exact, _ := json.Marshal(Response{EnvelopeID: envelope, Payload: padding})
	assert.Equal(t, maxResponseSize, len(exact), "test setup: payload should serialize to exactly %d bytes", maxResponseSize)

	err := cli.SendCtx(context.Background(), Response{
		EnvelopeID: envelope,
		Payload:    padding,
	})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "silently dropped")
}

func TestSendCtx_PayloadJustUnderLimit(t *testing.T) {
	api := slack.New("ABCDEFG")
	cli := New(api)

	envelope := "test-envelope"
	overhead, _ := json.Marshal(Response{EnvelopeID: envelope, Payload: ""})
	padding := strings.Repeat("x", maxResponseSize-len(overhead)-1)

	// Verify we're one byte under.
	under, _ := json.Marshal(Response{EnvelopeID: envelope, Payload: padding})
	assert.Equal(t, maxResponseSize-1, len(under), "test setup: payload should serialize to exactly %d bytes", maxResponseSize-1)

	err := cli.SendCtx(context.Background(), Response{
		EnvelopeID: envelope,
		Payload:    padding,
	})
	assert.NoError(t, err)

	select {
	case res := <-cli.socketModeResponses:
		assert.Equal(t, envelope, res.EnvelopeID)
	default:
		t.Fatal("expected response to be queued on socketModeResponses channel")
	}
}

func TestSendCtx_MarshalError(t *testing.T) {
	api := slack.New("ABCDEFG")
	cli := New(api)

	err := cli.SendCtx(context.Background(), Response{
		EnvelopeID: "test-envelope",
		Payload:    func() {},
	})

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "marshalling socket mode response")
}

func TestAck_ReturnsNoError(t *testing.T) {
	api := slack.New("ABCDEFG")
	cli := New(api)

	err := cli.Ack(Request{EnvelopeID: "test-envelope"}, "payload")
	assert.NoError(t, err)

	select {
	case res := <-cli.socketModeResponses:
		assert.Equal(t, "test-envelope", res.EnvelopeID)
	default:
		t.Fatal("expected response to be queued on socketModeResponses channel")
	}
}

func TestAck_ReturnsErrorWhenPayloadTooLarge(t *testing.T) {
	api := slack.New("ABCDEFG")
	cli := New(api)

	largePayload := strings.Repeat("x", maxResponseSize)
	err := cli.Ack(Request{EnvelopeID: "test-envelope"}, largePayload)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "silently dropped")
}

// dialTestWebSocket starts a slacktest.Server with a custom /ws handler,
// then returns a client-side *websocket.Conn connected to it.
func dialTestWebSocket(t *testing.T, serverFunc func(conn *websocket.Conn)) *websocket.Conn {
	t.Helper()
	srv := slacktest.NewTestServer(func(c slacktest.Customize) {
		c.Handle("/ws", slacktest.Websocket(serverFunc))
	})
	srv.Start()
	t.Cleanup(srv.Stop)

	conn, _, err := websocket.DefaultDialer.Dial(srv.GetWSURL(), nil)
	if err != nil {
		t.Fatalf("dial: %v", err)
	}
	t.Cleanup(func() { conn.Close() })
	return conn
}

func TestReceiveMessagesInto_WebSocketCloseError(t *testing.T) {
	conn := dialTestWebSocket(t, func(srvConn *websocket.Conn) {
		// Cleanly close the server side to trigger a CloseError on the client.
		srvConn.WriteMessage(
			websocket.CloseMessage,
			websocket.FormatCloseMessage(websocket.CloseGoingAway, "bye"),
		)
	})

	api := slack.New("ABCDEFG")
	cli := New(api)
	sink := make(chan json.RawMessage, 1)

	err := cli.receiveMessagesInto(context.Background(), conn, sink)

	// WebSocket close errors SHOULD force a reconnect (err != nil).
	assert.Error(t, err)
}

// TestRunMessageReceiver_SurvivesMalformedJSON demonstrates the full message
// receiver loop: valid messages are forwarded, malformed JSON produces an error
// event without dropping the connection, and a WebSocket close terminates the
// loop so the caller can reconnect.
func TestRunMessageReceiver_SurvivesMalformedJSON(t *testing.T) {
	ready := make(chan struct{})
	conn := dialTestWebSocket(t, func(srvConn *websocket.Conn) {
		// Wait for the receiver loop to start before sending.
		<-ready

		// 1. Valid JSON message
		srvConn.WriteMessage(websocket.TextMessage, []byte(`{"type":"hello"}`))
		// 2. Malformed JSON — should NOT kill the connection
		srvConn.WriteMessage(websocket.TextMessage, []byte(`{not json`))
		// 3. Another valid message — proves the connection survived
		srvConn.WriteMessage(websocket.TextMessage, []byte(`{"type":"disconnect"}`))
		// 4. Close the WebSocket — should terminate the loop
		srvConn.WriteMessage(
			websocket.CloseMessage,
			websocket.FormatCloseMessage(websocket.CloseNormalClosure, "done"),
		)
	})

	api := slack.New("ABCDEFG")
	cli := New(api)
	sink := make(chan json.RawMessage, 10)

	// Run the receiver loop in a goroutine.
	loopErr := make(chan error, 1)
	go func() {
		close(ready)
		loopErr <- cli.runMessageReceiver(context.Background(), conn, sink)
	}()

	// Wait for the loop to finish.
	var err error
	select {
	case err = <-loopErr:
	case <-time.After(5 * time.Second):
		t.Fatal("timed out waiting for receiver loop to finish")
	}

	// Drain all messages that made it through to the sink.
	var messages []string
	for len(sink) > 0 {
		messages = append(messages, string(<-sink))
	}

	// The loop should exit with an error (WebSocket close triggers reconnect).
	assert.Error(t, err)

	// Both valid messages should have been forwarded.
	assert.Equal(t, []string{`{"type":"hello"}`, `{"type":"disconnect"}`}, messages)

	// The malformed JSON should have produced an IncomingError event.
	select {
	case evt := <-cli.Events:
		assert.Equal(t, EventTypeIncomingError, evt.Type)
		incomingErr, ok := evt.Data.(*slack.IncomingEventError)
		if assert.True(t, ok) {
			var syntaxErr *json.SyntaxError
			assert.ErrorAs(t, incomingErr.ErrorObj, &syntaxErr)
		}
	default:
		t.Fatal("expected an IncomingError event for the malformed JSON")
	}
}
