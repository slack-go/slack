package slack

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestScopeDiagnosticsReturnedByAPIMethods(t *testing.T) {
	for _, tt := range []struct {
		name       string
		path       string
		errorsJSON string
		call       func(*Client) (*SlackResponse, error)
	}{
		{
			name: "agent response", path: "/agents.conversations.setProperties", errorsJSON: `[]`,
			call: func(api *Client) (*SlackResponse, error) {
				response, err := api.SetAgentConversationProperties("C123ABC456",
					SetAgentConversationPropertiesOptionCodeChannel(AgentCodeChannelProperties{
						SummaryMessage: &AgentSummaryMessage{MessageTS: "1717171717.123456"},
					}))
				if response == nil {
					return nil, err
				}
				return &response.SlackResponse, err
			},
		},
		{
			name: "agent method returning only an error", path: "/agents.conversations.archive", errorsJSON: `[]`,
			call: func(api *Client) (*SlackResponse, error) {
				return nil, api.ArchiveAgentConversation("C123ABC456")
			},
		},
		{
			name: "kick response with an empty errors object", path: "/conversations.kick", errorsJSON: `{}`,
			call: func(api *Client) (*SlackResponse, error) {
				return nil, api.KickUserFromConversation("C123ABC456", "U123ABC456")
			},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, tt.path, r.URL.Path)
				w.Header().Set("Content-Type", "application/json")
				_, _ = fmt.Fprintf(w, `{"ok":false,"error":"missing_scope","needed":"code_channels:manage","provided":"channels:read,chat:write","errors":%s,"response_metadata":{"messages":["Token lacks the required scope"]}}`, tt.errorsJSON)
			}))
			defer server.Close()
			api := New("testing-token", OptionAPIURL(server.URL+"/"))
			response, err := tt.call(api)
			require.EqualError(t, err, "missing_scope", "preserve the existing error text")
			var slackErr SlackErrorResponse
			require.ErrorAs(t, fmt.Errorf("API request failed: %w", err), &slackErr)
			assert.Equal(t, "code_channels:manage", slackErr.Needed)
			assert.Equal(t, "channels:read,chat:write", slackErr.Provided)
			assert.Equal(t, []string{"Token lacks the required scope"}, slackErr.ResponseMetadata.Messages)
			if response != nil {
				assert.Equal(t, slackErr.Needed, response.Needed)
				assert.Equal(t, slackErr.Provided, response.Provided)
			}
		})
	}
}

func TestScopeDiagnosticsOptional(t *testing.T) {
	for _, tt := range []struct {
		body string
		want string
	}{
		{body: `{"ok":true}`},
		{body: `{"ok":false,"error":"invalid_auth"}`, want: "invalid_auth"},
	} {
		var response SlackResponse
		require.NoError(t, json.Unmarshal([]byte(tt.body), &response))
		assert.Empty(t, response.Needed)
		assert.Empty(t, response.Provided)
		encoded, err := json.Marshal(response)
		require.NoError(t, err)
		assert.NotContains(t, string(encoded), `"needed"`)
		assert.NotContains(t, string(encoded), `"provided"`)
		if tt.want == "" {
			require.NoError(t, response.Err())
		} else {
			require.EqualError(t, response.Err(), tt.want)
			var slackErr SlackErrorResponse
			require.ErrorAs(t, response.Err(), &slackErr)
			assert.Empty(t, slackErr.Needed)
			assert.Empty(t, slackErr.Provided)
		}
	}
}
