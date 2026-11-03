package slackevents

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUserConnectionEvents(t *testing.T) {
	for _, tt := range []struct {
		name, payload string
		want          any
	}{
		{"connect", `{"type":"user_connection","subtype":"connect","user":"U123","trigger_id":"trigger","event_ts":"1.0"}`, &UserConnectionEvent{Type: "user_connection", Subtype: "connect", User: "U123", TriggerID: "trigger", EventTimestamp: "1.0"}},
		{"disconnect", `{"type":"user_connection","subtype":"disconnect","user":"U123","enterprise_id":"E123","event_ts":"1.0"}`, &UserConnectionEvent{Type: "user_connection", Subtype: "disconnect", User: "U123", EnterpriseID: "E123", EventTimestamp: "1.0"}},
		{"disconnect without enterprise", `{"type":"user_connection","subtype":"disconnect","user":"U123","event_ts":"1.0"}`, &UserConnectionEvent{Type: "user_connection", Subtype: "disconnect", User: "U123", EventTimestamp: "1.0"}},
		{"status requested", `{"type":"user_connection_status_requested","user":"U123","event_ts":"1.0"}`, &UserConnectionStatusRequestedEvent{Type: "user_connection_status_requested", User: "U123", EventTimestamp: "1.0"}},
	} {
		t.Run(tt.name, func(t *testing.T) {
			event, err := ParseEvent(json.RawMessage(`{"type":"event_callback","event":`+tt.payload+`}`), OptionNoVerifyToken())
			require.NoError(t, err)
			assert.Equal(t, tt.want, event.InnerEvent.Data)
			encoded, err := json.Marshal(event.InnerEvent.Data)
			require.NoError(t, err)
			assert.JSONEq(t, tt.payload, string(encoded))
		})
	}
}
