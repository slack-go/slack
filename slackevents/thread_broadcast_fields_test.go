package slackevents

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestThreadBroadcastDisplayAsBot(t *testing.T) {
	for _, value := range []string{"true", "false"} {
		t.Run(value, func(t *testing.T) {
			raw := `{"type":"event_callback","event":{"type":"message","subtype":"thread_broadcast","files":[{"id":"F123"}],"upload":true,"parent_user_id":"U123","display_as_bot":` + value + `}}`
			parsed, err := ParseEvent(json.RawMessage(raw), OptionNoVerifyToken())
			require.NoError(t, err)
			event, ok := parsed.InnerEvent.Data.(*MessageEvent)
			require.True(t, ok)
			require.NotNil(t, event.Message)
			require.Len(t, event.Message.Files, 1)
			require.Equal(t, "F123", event.Message.Files[0].ID)
			require.True(t, event.Message.Upload)
			require.Equal(t, "U123", event.Message.ParentUserId)
			encoded, err := json.Marshal(event.Message)
			require.NoError(t, err)
			var fields map[string]json.RawMessage
			require.NoError(t, json.Unmarshal(encoded, &fields))
			require.Equal(t, value, string(fields["display_as_bot"]))
		})
	}
	parsed, err := ParseEvent(json.RawMessage(`{"type":"event_callback","event":{"type":"message","text":"hello"}}`), OptionNoVerifyToken())
	require.NoError(t, err)
	event := parsed.InnerEvent.Data.(*MessageEvent)
	encoded, err := json.Marshal(event.Message)
	require.NoError(t, err)
	require.NotContains(t, string(encoded), `"display_as_bot"`)
}

func TestMessageDisplayAsBotWireShapes(t *testing.T) {
	for _, value := range []string{`"false"`, `0`, `[]`, `{}`} {
		_, err := ParseEvent(json.RawMessage(`{"type":"event_callback","event":{"type":"message","display_as_bot":`+value+`}}`), OptionNoVerifyToken())
		require.Error(t, err)
	}
	event, err := ParseEvent(json.RawMessage(`{"type":"event_callback","event":{"type":"message","display_as_bot":null}}`), OptionNoVerifyToken())
	require.NoError(t, err)
	require.Nil(t, event.InnerEvent.Data.(*MessageEvent).Message.DisplayAsBot)
}
