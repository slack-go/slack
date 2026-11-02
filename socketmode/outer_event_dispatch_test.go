package socketmode

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
)

func TestDispatchEventsAPIOuterType(t *testing.T) {
	for _, tt := range []struct {
		name    string
		payload string
		want    slackevents.EventsAPIType
	}{
		{"outer event", `{"type":"app_rate_limited","team_id":"T123","minute_rate_limited":123,"api_app_id":"A123"}`, slackevents.EventsAPIType(slackevents.AppRateLimited)},
		{"inner event takes precedence", `{"type":"event_callback","event":{"type":"app_mention","user":"U123","text":"hello"}}`, slackevents.AppMention},
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
