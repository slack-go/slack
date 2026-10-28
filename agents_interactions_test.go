package slack

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Code-channel Block Kit views use the standard interaction parser. These
// fixtures exercise the documented view/container fields, rather than assuming
// a separate code-channel interaction type.
func TestAgentViewInteractionHTTP(t *testing.T) {
	for _, tt := range []struct {
		name    string
		payload string
		kind    InteractionType
	}{
		{
			name: "button and input state",
			kind: InteractionTypeBlockActions,
			payload: `{
				"type":"block_actions","api_app_id":"A123ABC456",
				"team":{"id":"T123ABC456"},"user":{"id":"U2345678901"},
				"channel":{"id":"C9876543210"},"trigger_id":"13345224609.738474920",
				"container":{"type":"view","view_id":"Ct1234567890","channel_id":"C9876543210"},
				"actions":[{"type":"button","action_id":"approve","block_id":"decision","value":"approve"}],
				"state":{"values":{"notes":{"comment":{"type":"plain_text_input","value":"Ship it"}}}}
			}`,
		},
		{
			name: "external select suggestion",
			kind: InteractionTypeBlockSuggestion,
			payload: `{
				"type":"block_suggestion","api_app_id":"A123ABC456",
				"team":{"id":"T123ABC456"},"user":{"id":"U2345678901"},
				"channel":{"id":"C9876543210"},
				"container":{"type":"view","view_id":"Ct1234567890","channel_id":"C9876543210"},
				"block_id":"repository","action_id":"repo_lookup","value":"billing"
			}`,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			form := url.Values{"payload": {tt.payload}}
			request := httptest.NewRequest(http.MethodPost, "/interactions", strings.NewReader(form.Encode()))
			request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			callback, err := InteractionCallbackParse(request)
			require.NoError(t, err)
			assert.Equal(t, tt.kind, callback.Type)
			assert.Equal(t, "A123ABC456", callback.APIAppID)
			assert.Equal(t, "C9876543210", callback.Channel.ID)
			assert.Equal(t, "Ct1234567890", callback.Container.ViewID)
			assert.Equal(t, "C9876543210", callback.Container.ChannelID)
			if tt.kind == InteractionTypeBlockActions {
				require.Len(t, callback.ActionCallback.BlockActions, 1)
				assert.Equal(t, "approve", callback.ActionCallback.BlockActions[0].ActionID)
				assert.Equal(t, "decision", callback.ActionCallback.BlockActions[0].BlockID)
				require.NotNil(t, callback.BlockActionState)
				assert.Equal(t, "Ship it", callback.BlockActionState.Values["notes"]["comment"].Value)
			} else {
				assert.Equal(t, "repo_lookup", callback.ActionID)
				assert.Equal(t, "repository", callback.BlockID)
				assert.Equal(t, "billing", callback.Value)
			}

			encoded, err := json.Marshal(callback)
			require.NoError(t, err)
			var decoded InteractionCallback
			require.NoError(t, json.Unmarshal(encoded, &decoded))
			assert.Equal(t, callback.Container.Type, decoded.Container.Type)
			assert.Equal(t, callback.Container.ViewID, decoded.Container.ViewID)
			assert.Equal(t, callback.Container.ChannelID, decoded.Container.ChannelID)
			if tt.kind == InteractionTypeBlockActions {
				require.Len(t, decoded.ActionCallback.BlockActions, 1)
				assert.Equal(t, "approve", decoded.ActionCallback.BlockActions[0].ActionID)
				assert.Equal(t, "decision", decoded.ActionCallback.BlockActions[0].BlockID)
				assert.Equal(t, "approve", decoded.ActionCallback.BlockActions[0].Value)
				require.NotNil(t, decoded.BlockActionState)
				assert.Equal(t, "Ship it", decoded.BlockActionState.Values["notes"]["comment"].Value)
			}
			assert.Equal(t, callback.Value, decoded.Value)
		})
	}
}

func TestAgentSlashCommandHTTP(t *testing.T) {
	// Agent commands use the same form payload as manifest slash commands:
	// https://docs.slack.dev/ai/slack-code/agent-slash-commands/
	form := url.Values{
		"token": {"XXYYZZ"}, "team_id": {"T123ABC456"}, "team_domain": {"example"},
		"channel_id": {"C9876543210"}, "channel_name": {"fix-flaky-login"},
		"user_id": {"U2345678901"}, "user_name": {"alice"},
		"command": {"/create-pr"}, "text": {"Fix login & preserve café names"},
		"api_app_id":   {"A123ABC456"},
		"response_url": {"https://hooks.slack.com/commands/T123ABC456/1234567890/example"},
		"trigger_id":   {"13345224609.738474920"},
	}
	request := httptest.NewRequest(http.MethodPost, "/agent-commands", strings.NewReader(form.Encode()))
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	command, err := SlashCommandParse(request)
	require.NoError(t, err)
	assert.Equal(t, "/create-pr", command.Command)
	assert.Equal(t, "Fix login & preserve café names", command.Text)
	assert.Equal(t, "C9876543210", command.ChannelID)
	assert.Equal(t, "U2345678901", command.UserID)
	assert.Equal(t, "A123ABC456", command.APIAppID)
	assert.Equal(t, form.Get("response_url"), command.ResponseURL)
	assert.Equal(t, form.Get("trigger_id"), command.TriggerID)
}
