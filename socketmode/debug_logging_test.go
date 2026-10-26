package socketmode

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
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

func TestSocketModeDebugLogsRedactIncomingCredentials(t *testing.T) {
	for _, tt := range []struct {
		name    string
		message string
		secrets []string
		visible []string
	}{
		{
			name:    "slash command",
			message: `{"type":"slash_commands","envelope_id":"debug-envelope","payload":{"is_enterprise_install":false,"token":"verification-secret","response_url":"https://hooks.slack.com/commands/response-secret","trigger_id":"trigger-debug-id","command":"/debug","text":"diagnostic text"}}`,
			secrets: []string{"verification-secret", "response-secret"},
			visible: []string{"trigger-debug-id"},
		},
		{
			name:    "nested response URLs",
			message: `{"type":"interactive","envelope_id":"debug-envelope","payload":{"type":"view_submission","response_urls":[{"response_url":"https://hooks.slack.com/commands/nested-secret","channel_id":"C123"}],"view":{"private_metadata":"diagnostic text"}}}`,
			secrets: []string{"nested-secret"},
		},
		{
			name:    "function token",
			message: `{"type":"events_api","envelope_id":"debug-envelope","payload":{"type":"event_callback","event":{"type":"function_executed","bot_access_token":"xoxb-function-secret","inputs":{"text":"diagnostic text"}}}}`,
			secrets: []string{"xoxb-function-secret"},
		},
		{
			name:    "escaped and case insensitive keys",
			message: `{"type":"slash_commands","envelope_id":"debug-envelope","payload":{"is_enterprise_install":false,"\u0074oken":"escaped-secret","RESPONSE_URL":"https://hooks.slack.com/commands/uppercase-secret","text":"diagnostic text"}}`,
			secrets: []string{"escaped-secret", "uppercase-secret"},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			var logs bytes.Buffer
			client := New(slack.New("test-token"), OptionDebug(true), OptionLog(log.New(&logs, "", 0)))
			conn := dialTestWebSocket(t, sendFrames([]byte(tt.message)))
			require.NoError(t, conn.SetReadDeadline(time.Now().Add(5*time.Second)))
			messages := make(chan json.RawMessage, 1)
			require.NoError(t, client.receiveMessagesInto(context.Background(), conn, messages))
			raw := <-messages
			require.Equal(t, tt.message, string(raw), "redaction must not change the received bytes")
			messages <- raw
			close(messages)
			require.NoError(t, client.runRequestHandler(context.Background(), messages))

			require.Len(t, client.Events, 1)
			event := <-client.Events
			require.NotNil(t, event.Request, "unexpected event: %+v", event)
			var original Request
			require.NoError(t, json.Unmarshal(raw, &original))
			require.Equal(t, original.Payload, event.Request.Payload, "handlers need the original credentials")
			for _, secret := range tt.secrets {
				assert.NotContains(t, logs.String(), secret)
			}
			for _, value := range tt.visible {
				assert.Contains(t, logs.String(), value)
			}
			assert.Contains(t, logs.String(), "debug-envelope")
			assert.Contains(t, logs.String(), "diagnostic text")
		})
	}
}

func TestSocketModeDebugLogsOmitMalformedPayload(t *testing.T) {
	for _, message := range []string{
		`{"token":"malformed-secret",`,
		`{"type":"hello"} {"token":"malformed-secret"}`,
	} {
		var logs bytes.Buffer
		client := New(slack.New("test-token"), OptionDebug(true), OptionLog(log.New(&logs, "", 0)))
		messages := make(chan json.RawMessage, 1)
		messages <- json.RawMessage(message)
		close(messages)
		require.NoError(t, client.runRequestHandler(context.Background(), messages))
		assert.NotContains(t, logs.String(), "malformed-secret")
		require.Len(t, client.Events, 1)
		event := <-client.Events
		require.Equal(t, EventTypeErrorBadMessage, event.Type)
		assert.Equal(t, message, string(event.Data.(*ErrorBadMessage).Message))
	}
}

func TestSocketModeDebugLogsRedactOutgoingCredentials(t *testing.T) {
	var logs bytes.Buffer
	client := New(slack.New("test-token"), OptionDebug(true), OptionLog(log.New(&logs, "", 0)))
	response := Response{EnvelopeID: "debug-envelope", Payload: map[string]any{
		"text": "diagnostic text", "token": "outgoing-secret", "action_token": "action-secret",
		"number": json.Number("1234567890123456789"),
	}}
	want, err := json.Marshal(response)
	require.NoError(t, err)
	require.NoError(t, client.SendCtx(context.Background(), response))

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	received := make(chan []byte, 1)
	conn := dialTestWebSocket(t, func(conn *websocket.Conn) {
		defer cancel()
		_ = conn.SetReadDeadline(time.Now().Add(5 * time.Second))
		_, message, err := conn.ReadMessage()
		if err != nil {
			t.Error(err)
			return
		}
		received <- message
	})
	require.ErrorIs(t, client.runResponseSender(ctx, conn), context.Canceled)
	select {
	case message := <-received:
		assert.Equal(t, want, message, "redaction must not change the response sent to Slack")
	default:
		t.Fatal("response was not received")
	}
	assert.NotContains(t, logs.String(), "outgoing-secret")
	assert.NotContains(t, logs.String(), "action-secret")
	assert.Contains(t, logs.String(), "debug-envelope")
	assert.Contains(t, logs.String(), "diagnostic text")
	assert.Contains(t, logs.String(), "1234567890123456789")
}

func TestSocketModeDebugLogsOmitConnectionTicket(t *testing.T) {
	var rootLogs, socketLogs bytes.Buffer
	tickets := make(chan string, 1)
	upgrader := websocket.Upgrader{CheckOrigin: func(*http.Request) bool { return true }}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/apps.connections.open" {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]any{
				"ok": true, "url": "ws://" + r.Host + "/ws?ticket=connection-secret",
			})
			return
		}
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			t.Error(err)
			return
		}
		defer conn.Close()
		tickets <- r.URL.Query().Get("ticket")
	}))
	defer server.Close()
	api := slack.New("test-token", slack.OptionAPIURL(server.URL+"/"), slack.OptionAppLevelToken("xapp-test"),
		slack.OptionDebug(true), slack.OptionLog(log.New(&rootLogs, "", 0)))
	client := New(api, OptionDebug(true), OptionLog(log.New(&socketLogs, "", 0)))
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	info, conn, err := client.openAndDial(ctx, nil)
	require.NoError(t, err)
	defer conn.Close()
	assert.Contains(t, info.URL, "ticket=connection-secret")
	select {
	case ticket := <-tickets:
		assert.Equal(t, "connection-secret", ticket)
	case <-ctx.Done():
		t.Fatal("server did not receive the ticket")
	}
	for _, output := range []string{rootLogs.String(), socketLogs.String()} {
		assert.NotContains(t, output, "connection-secret")
		assert.Contains(t, output, strings.TrimPrefix(server.URL, "http://"))
	}
}
