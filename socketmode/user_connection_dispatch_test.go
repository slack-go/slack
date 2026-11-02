package socketmode

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
)

func TestUserConnectionSocketModeDispatch(t *testing.T) {
	for _, tt := range []struct {
		name    string
		payload string
		want    slackevents.EventsAPIType
	}{
		{"user connection", `{"type":"event_callback","event":{"type":"user_connection","subtype":"connect","user":"U123","trigger_id":"trigger","event_ts":"1.0"}}`, slackevents.UserConnection},
		{"connection status", `{"type":"event_callback","event":{"type":"user_connection_status_requested","user":"U123","event_ts":"1.0"}}`, slackevents.UserConnectionStatusRequested},
	} {
		t.Run(tt.name, func(t *testing.T) {
			parsed, err := slackevents.ParseEvent(json.RawMessage(tt.payload), slackevents.OptionNoVerifyToken())
			require.NoError(t, err)
			handler := NewSocketmodeHandler(New(slack.New("test-token")))
			called := make(chan string, 2)
			handler.HandleEvents(tt.want, func(event *Event, _ *Client) { called <- "typed" })
			handler.HandleDefault(func(event *Event, _ *Client) { called <- "default" })
			handler.DispatchEvent(Event{Type: EventTypeEventsAPI, Data: parsed})
			require.Equal(t, "typed", <-called)
		})
	}
}
