package socketmode

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/slack-go/slack"
	"github.com/stretchr/testify/require"
)

func TestRunRequestHandlerContinuesAfterMalformedCallback(t *testing.T) {
	for _, test := range []struct {
		name    string
		message string
	}{
		{
			name:    "missing inner event",
			message: `{"type":"events_api","payload":{"type":"event_callback"}}`,
		},
		{
			name:    "null inner event",
			message: `{"type":"events_api","payload":{"type":"event_callback","event":null}}`,
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			client := New(slack.New("test-token"))
			messages := make(chan json.RawMessage, 2)
			messages <- json.RawMessage(test.message)
			messages <- json.RawMessage(`{"type":"events_api","envelope_id":"valid-event","payload":{"type":"event_callback","event":{"type":"app_mention","user":"U123","text":"hello"}}}`)
			close(messages)

			require.NotPanics(t, func() {
				require.NoError(t, client.runRequestHandler(context.Background(), messages))
			})
			require.Len(t, client.Events, 2)

			event := <-client.Events
			require.Equal(t, EventTypeErrorBadMessage, event.Type)
			badMessage, ok := event.Data.(*ErrorBadMessage)
			require.True(t, ok, "expected *ErrorBadMessage, got %T", event.Data)
			require.ErrorContains(t, badMessage.Cause, "missing or null event")
			require.Equal(t, test.message, string(badMessage.Message))

			event = <-client.Events
			require.Equal(t, EventTypeEventsAPI, event.Type)
			require.NotNil(t, event.Request)
			require.Equal(t, "valid-event", event.Request.EnvelopeID)
		})
	}
}
