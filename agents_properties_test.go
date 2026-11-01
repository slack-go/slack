package slack

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSetAgentConversationProperties(t *testing.T) {
	resourceURL := "https://github.com/example/repo/pull/42"
	resourceType, title, provider := "pull_request", "Fix login", "github"
	empty := ""
	for _, tt := range []struct {
		name    string
		options []SetAgentConversationPropertiesOption
		want    map[string]string
	}{
		{
			name: "all properties",
			options: []SetAgentConversationPropertiesOption{
				SetAgentConversationPropertiesOptionCodeChannel(AgentCodeChannelProperties{
					ContextBarItems: []AgentContextBarItem{
						{Key: "repo", Label: "example/repo", Icon: "folder", URL: "https://github.com/example/repo"},
						{Key: "run-tests", Label: "Run tests", ItemType: "action"},
						{Key: "ci", Label: "Tests passing"},
					},
					SummaryMessage: &AgentSummaryMessage{MessageTS: "1717182000.456789", ThreadTS: "1717171717.123456"},
				}),
				SetAgentConversationPropertiesOptionAgentResource(AgentResourceProperties{URL: &resourceURL, ResourceType: &resourceType, Title: &title, Provider: &provider}),
			},
			want: map[string]string{
				"code_channel":   `{"context_bar_items":[{"key":"repo","label":"example/repo","icon":"folder","url":"https://github.com/example/repo"},{"key":"run-tests","label":"Run tests","item_type":"action"},{"key":"ci","label":"Tests passing"}],"summary_message":{"message_ts":"1717182000.456789","thread_ts":"1717171717.123456"}}`,
				"agent_resource": `{"url":"https://github.com/example/repo/pull/42","resource_type":"pull_request","title":"Fix login","provider":"github"}`,
			},
		},
		{
			name: "summary leaves context bar unchanged",
			options: []SetAgentConversationPropertiesOption{
				SetAgentConversationPropertiesOptionCodeChannel(AgentCodeChannelProperties{
					SummaryMessage: &AgentSummaryMessage{MessageTS: "1717182000.456789"},
				}),
			},
			want: map[string]string{"code_channel": `{"summary_message":{"message_ts":"1717182000.456789"}}`},
		},
		{
			name: "empty slice clears context bar",
			options: []SetAgentConversationPropertiesOption{
				SetAgentConversationPropertiesOptionCodeChannel(AgentCodeChannelProperties{
					ContextBarItems: []AgentContextBarItem{},
				}),
			},
			want: map[string]string{"code_channel": `{"context_bar_items":[]}`},
		},
		{
			name: "resource only omits other fields",
			options: []SetAgentConversationPropertiesOption{
				SetAgentConversationPropertiesOptionAgentResource(AgentResourceProperties{URL: &resourceURL}),
			},
			want: map[string]string{"agent_resource": `{"url":"https://github.com/example/repo/pull/42"}`},
		},
		{
			name: "explicit empty resource fields",
			options: []SetAgentConversationPropertiesOption{
				SetAgentConversationPropertiesOptionAgentResource(AgentResourceProperties{URL: &empty, ResourceType: &empty, Title: &empty, Provider: &empty}),
			},
			want: map[string]string{"agent_resource": `{"url":"","resource_type":"","title":"","provider":""}`},
		},
		{
			name: "nil option is ignored",
			options: []SetAgentConversationPropertiesOption{
				nil,
				SetAgentConversationPropertiesOptionAgentResource(AgentResourceProperties{URL: &resourceURL}),
			},
			want: map[string]string{"agent_resource": `{"url":"https://github.com/example/repo/pull/42"}`},
		},
		{
			name: "last option for the same property wins",
			options: []SetAgentConversationPropertiesOption{
				SetAgentConversationPropertiesOptionAgentResource(AgentResourceProperties{URL: &resourceURL}),
				SetAgentConversationPropertiesOptionAgentResource(AgentResourceProperties{Title: &title}),
			},
			want: map[string]string{"agent_resource": `{"title":"Fix login"}`},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, http.MethodPost, r.Method)
				assert.Equal(t, "/agents.conversations.setProperties", r.URL.Path)
				if !assert.NoError(t, r.ParseForm()) {
					return
				}
				assert.Equal(t, "testing-token", r.PostForm.Get("token"))
				assert.Equal(t, "C9876543210", r.PostForm.Get("channel_id"))
				assert.Len(t, r.PostForm, 2+len(tt.want))
				for key, wantJSON := range tt.want {
					assert.JSONEq(t, wantJSON, r.PostForm.Get(key))
				}
				w.Header().Set("Content-Type", "application/json")
				_, _ = w.Write([]byte(`{"ok":true,"channel_id":"C9876543210"}`))
			}))
			defer server.Close()
			api := New("testing-token", OptionAPIURL(server.URL+"/"))
			response, err := api.SetAgentConversationProperties("C9876543210", tt.options...)
			require.NoError(t, err)
			require.NotNil(t, response)
			assert.True(t, response.Ok)
			assert.Equal(t, "C9876543210", response.ChannelID)
		})
	}
}

func TestSetAgentConversationPropertiesWarnings(t *testing.T) {
	items := make([]AgentContextBarItem, 6)
	for i := range items {
		items[i] = AgentContextBarItem{Key: fmt.Sprintf("item-%d", i), Label: "Context"}
	}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !assert.NoError(t, r.ParseForm()) {
			return
		}
		var properties AgentCodeChannelProperties
		if !assert.NoError(t, json.Unmarshal([]byte(r.PostForm.Get("code_channel")), &properties)) {
			return
		}
		assert.Len(t, properties.ContextBarItems, 6, "let Slack enforce limits and return warnings")
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"ok":true,"channel_id":"C9876543210","warning":"too_many_context_bar_items"}`))
	}))
	defer server.Close()
	warningReceived := false
	api := New("testing-token", OptionAPIURL(server.URL+"/"), OptionOnWarning(func(path string, request any, warning *Warning) {
		warningReceived = true
		assert.Equal(t, "agents.conversations.setProperties", path)
		assert.Equal(t, []string{"too_many_context_bar_items"}, warning.Codes)
	}))
	response, err := api.SetAgentConversationProperties("C9876543210",
		SetAgentConversationPropertiesOptionCodeChannel(AgentCodeChannelProperties{ContextBarItems: items}))
	require.NoError(t, err)
	require.NotNil(t, response)
	assert.Equal(t, "too_many_context_bar_items", response.Warning)
	assert.True(t, warningReceived)
}

func TestSetAgentConversationPropertiesErrors(t *testing.T) {
	for _, tt := range []struct {
		name   string
		status int
		body   string
	}{
		{"Slack error", http.StatusOK, `{"ok":false,"error":"no_properties_provided","response_metadata":{"messages":["Provide at least one property"]}}`},
		{"HTTP error", http.StatusInternalServerError, `server error`},
		{"invalid JSON", http.StatusOK, `{"ok":`},
	} {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tt.status)
				_, _ = w.Write([]byte(tt.body))
			}))
			defer server.Close()
			api := New("testing-token", OptionAPIURL(server.URL+"/"))
			response, err := api.SetAgentConversationPropertiesContext(context.Background(), "C9876543210")
			require.Error(t, err)
			if tt.name == "Slack error" {
				var slackErr SlackErrorResponse
				require.ErrorAs(t, err, &slackErr)
				assert.Equal(t, "no_properties_provided", slackErr.Err)
				assert.Equal(t, []string{"Provide at least one property"}, slackErr.ResponseMetadata.Messages)
				require.NotNil(t, response)
				assert.False(t, response.Ok)
			} else {
				assert.Nil(t, response)
			}
		})
	}

	t.Run("canceled context", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			t.Error("canceled request reached server")
		}))
		defer server.Close()
		api := New("testing-token", OptionAPIURL(server.URL+"/"))
		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		_, err := api.SetAgentConversationPropertiesContext(ctx, "C9876543210")
		require.ErrorIs(t, err, context.Canceled)
	})
}
