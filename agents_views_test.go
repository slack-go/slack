package slack

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSetAgentConversationView(t *testing.T) {
	const html = "<!doctype html><html><body>Coverage: 97%</body></html>"
	const diff = "diff --git a/a.go b/a.go\n--- a/a.go\n+++ b/a.go\n@@ -1 +1 @@\n-old\n+new\n"
	for _, tt := range []struct {
		name         string
		options      []SetAgentConversationViewOption
		want         url.Values
		wantJSON     map[string]string
		responseType string
	}{
		{
			name: "html uses server default type",
			options: []SetAgentConversationViewOption{
				SetAgentConversationViewOptionViewKey("reports/coverage.html"),
				SetAgentConversationViewOptionContent(html),
				SetAgentConversationViewOptionName("Coverage"),
				SetAgentConversationViewOptionCSP(AgentViewCSP{
					ResourceDomains: []string{"https://cdn.jsdelivr.net"},
					ConnectDomains:  []string{"https://api.example.com"},
				}),
			},
			want: url.Values{
				"view_key": {"reports/coverage.html"}, "content": {html}, "name": {"Coverage"},
			},
			wantJSON: map[string]string{
				"csp": `{"resource_domains":["https://cdn.jsdelivr.net"],"connect_domains":["https://api.example.com"]}`,
			},
			responseType: AgentViewTypeHTML,
		},
		{
			name: "diff singleton without view key",
			options: []SetAgentConversationViewOption{
				SetAgentConversationViewOptionType(AgentViewTypeDiff),
				SetAgentConversationViewOptionContent(diff),
				SetAgentConversationViewOptionBaseBranch("main"),
				SetAgentConversationViewOptionHeadBranch("feature"),
			},
			want: url.Values{
				"type": {"diff"}, "content": {diff}, "base_branch": {"main"}, "head_branch": {"feature"},
			},
			responseType: AgentViewTypeDiff,
		},
		{
			name: "block kit uses a JSON array",
			options: []SetAgentConversationViewOption{
				SetAgentConversationViewOptionType(AgentViewTypeBlockKit),
				SetAgentConversationViewOptionViewKey("plan"),
				SetAgentConversationViewOptionBlocks(
					NewSectionBlock(NewTextBlockObject(MarkdownType, "*Plan*", false, false), nil, nil),
					NewDividerBlock(),
				),
			},
			want: url.Values{"type": {"block_kit"}, "view_key": {"plan"}},
			wantJSON: map[string]string{
				"blocks": `[{"type":"section","text":{"type":"mrkdwn","text":"*Plan*"}},{"type":"divider"}]`,
			},
			responseType: AgentViewTypeBlockKit,
		},
		{
			name: "canvas with access and content hash",
			options: []SetAgentConversationViewOption{
				SetAgentConversationViewOptionType(AgentViewTypeCanvas),
				SetAgentConversationViewOptionViewKey("design"),
				SetAgentConversationViewOptionCanvasID("F1234567890"),
				SetAgentConversationViewOptionAccessLevel("comment"),
				SetAgentConversationViewOptionAgentContentHash("sha256:abc123"),
			},
			want: url.Values{
				"type": {"canvas"}, "view_key": {"design"}, "canvas_id": {"F1234567890"},
				"access_level": {"comment"}, "agent_content_hash": {"sha256:abc123"},
			},
			responseType: AgentViewTypeCanvas,
		},
		{
			name: "pull request",
			options: []SetAgentConversationViewOption{
				SetAgentConversationViewOptionType(AgentViewTypePullRequest),
				SetAgentConversationViewOptionPRURL("https://github.com/example/repo/pull/42"),
			},
			want: url.Values{
				"type": {"pull_request"}, "pr_url": {"https://github.com/example/repo/pull/42"},
			},
			responseType: AgentViewTypePullRequest,
		},
		{
			name: "explicit empty values and last option wins",
			options: []SetAgentConversationViewOption{
				nil,
				SetAgentConversationViewOptionType(AgentViewTypeDiff),
				SetAgentConversationViewOptionContent("old"),
				SetAgentConversationViewOptionContent(""),
				SetAgentConversationViewOptionBlocks(NewDividerBlock()),
				SetAgentConversationViewOptionBlocks(),
				SetAgentConversationViewOptionCSP(AgentViewCSP{ResourceDomains: []string{"https://example.com"}}),
				SetAgentConversationViewOptionCSP(AgentViewCSP{}),
			},
			want:         url.Values{"type": {"diff"}, "content": {""}},
			wantJSON:     map[string]string{"blocks": `[]`, "csp": `{}`},
			responseType: AgentViewTypeDiff,
		},
		{
			name: "name and legacy label passed to Slack",
			options: []SetAgentConversationViewOption{
				SetAgentConversationViewOptionType(AgentViewTypeHTML),
				SetAgentConversationViewOptionViewKey("coverage.html"),
				SetAgentConversationViewOptionContent(html),
				SetAgentConversationViewOptionLabel("Old label"),
				SetAgentConversationViewOptionName("Old name"),
				SetAgentConversationViewOptionName("Coverage"),
			},
			want: url.Values{
				"type": {"html"}, "view_key": {"coverage.html"}, "content": {html},
				"label": {"Old label"}, "name": {"Coverage"},
			},
			responseType: AgentViewTypeHTML,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, http.MethodPost, r.Method)
				assert.Equal(t, "/agents.conversations.setView", r.URL.Path)
				require.NoError(t, r.ParseForm())
				assert.Equal(t, "testing-token", r.PostForm.Get("token"))
				assert.Equal(t, "C9876543210", r.PostForm.Get("channel_id"))
				r.PostForm.Del("token")
				r.PostForm.Del("channel_id")
				for key, value := range tt.wantJSON {
					assert.JSONEq(t, value, r.PostForm.Get(key))
					r.PostForm.Del(key)
				}
				assert.Equal(t, tt.want, r.PostForm)
				w.Header().Set("Content-Type", "application/json")
				response := map[string]any{
					"ok": true, "channel_id": "C9876543210", "view_id": "Ct1234567890",
					"file_id": "F1234567890", "content_version": 3, "type": tt.responseType,
				}
				if tt.responseType == AgentViewTypeCanvas {
					response["canvas_id"] = "F1234567890"
				}
				assert.NoError(t, json.NewEncoder(w).Encode(response))
			}))
			defer server.Close()
			api := New("testing-token", OptionAPIURL(server.URL+"/"))
			response, err := api.SetAgentConversationView("C9876543210", tt.options...)
			require.NoError(t, err)
			require.NotNil(t, response)
			assert.True(t, response.Ok)
			assert.Equal(t, "C9876543210", response.ChannelID)
			assert.Equal(t, "Ct1234567890", response.ViewID)
			assert.Equal(t, "F1234567890", response.FileID)
			assert.Equal(t, 3, response.ContentVersion)
			assert.Equal(t, tt.responseType, response.Type)
			if tt.responseType == AgentViewTypeCanvas {
				assert.Equal(t, "F1234567890", response.CanvasID)
			} else {
				assert.Empty(t, response.CanvasID)
			}
		})
	}
}

func TestSetAgentConversationViewErrors(t *testing.T) {
	for _, tt := range []struct {
		name   string
		status int
		body   string
	}{
		{"Slack error", http.StatusOK, `{"ok":false,"error":"too_many_views","response_metadata":{"messages":["View limit reached"]}}`},
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
			response, err := api.SetAgentConversationView("C9876543210")
			require.Error(t, err)
			if tt.name == "Slack error" {
				var slackErr SlackErrorResponse
				require.ErrorAs(t, err, &slackErr)
				assert.Equal(t, "too_many_views", slackErr.Err)
				assert.Equal(t, []string{"View limit reached"}, slackErr.ResponseMetadata.Messages)
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
		response, err := api.SetAgentConversationViewContext(ctx, "C9876543210")
		require.ErrorIs(t, err, context.Canceled)
		assert.Nil(t, response)
	})

	t.Run("block encoding error", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			t.Error("invalid block reached server")
		}))
		defer server.Close()
		api := New("testing-token", OptionAPIURL(server.URL+"/"))
		response, err := api.SetAgentConversationView("C9876543210",
			SetAgentConversationViewOptionType(AgentViewTypeBlockKit),
			SetAgentConversationViewOptionBlocks(&UnknownBlock{raw: json.RawMessage(`{`)}))
		require.Error(t, err)
		var marshalErr *json.MarshalerError
		require.ErrorAs(t, err, &marshalErr)
		assert.Nil(t, response)
	})
}

func TestListAgentConversationViews(t *testing.T) {
	for _, tt := range []struct {
		name string
		body string
		want []AgentConversationView
	}{
		{
			name: "views with and without a stable key",
			body: `{"ok":true,"views":[
				{"view_id":"T1234567890","view_key":"reports/coverage.html","file_id":"F1234567890","label":"Coverage","content_version":3,"date_added":1717171717},
				{"view_id":"T0987654321","file_id":"F0987654321","label":"Diff","content_version":5,"date_added":1717171800},
				{"view_id":"T1111111111","view_key":"plan","file_id":"F1111111111","label":"Plan","content_version":1,"date_added":1717171900}
			]}`,
			want: []AgentConversationView{
				{ViewID: "T1234567890", ViewKey: "reports/coverage.html", FileID: "F1234567890", Label: "Coverage", ContentVersion: 3, DateAdded: 1717171717},
				{ViewID: "T0987654321", FileID: "F0987654321", Label: "Diff", ContentVersion: 5, DateAdded: 1717171800},
				{ViewID: "T1111111111", ViewKey: "plan", FileID: "F1111111111", Label: "Plan", ContentVersion: 1, DateAdded: 1717171900},
			},
		},
		{
			name: "empty channel",
			body: `{"ok":true,"views":[]}`,
			want: []AgentConversationView{},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, http.MethodPost, r.Method)
				assert.Equal(t, "/agents.conversations.listViews", r.URL.Path)
				require.NoError(t, r.ParseForm())
				assert.Equal(t, url.Values{
					"token": {"testing-token"}, "channel_id": {"C9876543210"},
				}, r.PostForm)
				w.Header().Set("Content-Type", "application/json")
				_, _ = w.Write([]byte(tt.body))
			}))
			defer server.Close()
			api := New("testing-token", OptionAPIURL(server.URL+"/"))
			response, err := api.ListAgentConversationViews("C9876543210")
			require.NoError(t, err)
			require.NotNil(t, response)
			assert.True(t, response.Ok)
			assert.Equal(t, tt.want, response.Views)
		})
	}
}

func TestRemoveAgentConversationView(t *testing.T) {
	for _, tt := range []struct {
		name    string
		options []RemoveAgentConversationViewOption
		want    url.Values
	}{
		{
			name:    "by key",
			options: []RemoveAgentConversationViewOption{RemoveAgentConversationViewOptionViewKey("reports/coverage.html")},
			want:    url.Values{"view_key": {"reports/coverage.html"}},
		},
		{
			name:    "by ID",
			options: []RemoveAgentConversationViewOption{RemoveAgentConversationViewOptionViewID("T1234567890")},
			want:    url.Values{"view_id": {"T1234567890"}},
		},
		{
			name: "nil and empty options",
			options: []RemoveAgentConversationViewOption{
				nil,
				RemoveAgentConversationViewOptionViewKey(""),
				RemoveAgentConversationViewOptionViewID("T1234567890"),
			},
			want: url.Values{"view_id": {"T1234567890"}},
		},
		{
			name: "last option for the same selector wins",
			options: []RemoveAgentConversationViewOption{
				RemoveAgentConversationViewOptionViewID("T1234567890"),
				RemoveAgentConversationViewOptionViewID(""),
				RemoveAgentConversationViewOptionViewKey("old-key"),
				RemoveAgentConversationViewOptionViewKey("plan"),
			},
			want: url.Values{"view_key": {"plan"}},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, http.MethodPost, r.Method)
				assert.Equal(t, "/agents.conversations.removeView", r.URL.Path)
				require.NoError(t, r.ParseForm())
				assert.Equal(t, "testing-token", r.PostForm.Get("token"))
				assert.Equal(t, "C9876543210", r.PostForm.Get("channel_id"))
				r.PostForm.Del("token")
				r.PostForm.Del("channel_id")
				assert.Equal(t, tt.want, r.PostForm)
				w.Header().Set("Content-Type", "application/json")
				_, _ = w.Write([]byte(`{"ok":true,"channel_id":"C9876543210","view_id":"T1234567890"}`))
			}))
			defer server.Close()
			api := New("testing-token", OptionAPIURL(server.URL+"/"))
			response, err := api.RemoveAgentConversationView("C9876543210", tt.options...)
			require.NoError(t, err)
			require.NotNil(t, response)
			assert.True(t, response.Ok)
			assert.Equal(t, "C9876543210", response.ChannelID)
			assert.Equal(t, "T1234567890", response.ViewID)
		})
	}
}

func TestRemoveAgentConversationViewInvalidSelectors(t *testing.T) {
	for _, tt := range []struct {
		name    string
		options []RemoveAgentConversationViewOption
		want    url.Values
	}{
		{name: "neither", want: url.Values{}},
		{
			name: "both",
			options: []RemoveAgentConversationViewOption{
				RemoveAgentConversationViewOptionViewKey("plan"),
				RemoveAgentConversationViewOptionViewID("T1234567890"),
			},
			want: url.Values{"view_key": {"plan"}, "view_id": {"T1234567890"}},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				require.NoError(t, r.ParseForm())
				r.PostForm.Del("token")
				r.PostForm.Del("channel_id")
				assert.Equal(t, tt.want, r.PostForm, "let Slack validate the selector combination")
				w.Header().Set("Content-Type", "application/json")
				_, _ = w.Write([]byte(`{"ok":false,"error":"invalid_arguments"}`))
			}))
			defer server.Close()
			api := New("testing-token", OptionAPIURL(server.URL+"/"))
			response, err := api.RemoveAgentConversationView("C9876543210", tt.options...)
			require.Error(t, err)
			var slackErr SlackErrorResponse
			require.ErrorAs(t, err, &slackErr)
			assert.Equal(t, "invalid_arguments", slackErr.Err)
			require.NotNil(t, response)
			assert.False(t, response.Ok)
		})
	}
}

func TestAgentConversationViewManagementErrors(t *testing.T) {
	for _, method := range []struct {
		name       string
		slackError string
		call       func(*Client, context.Context) (*SlackResponse, error)
	}{
		{
			name: "list", slackError: "not_agent_channel",
			call: func(api *Client, ctx context.Context) (*SlackResponse, error) {
				response, err := api.ListAgentConversationViewsContext(ctx, "C9876543210")
				if response == nil {
					return nil, err
				}
				return &response.SlackResponse, err
			},
		},
		{
			name: "remove", slackError: "view_not_found",
			call: func(api *Client, ctx context.Context) (*SlackResponse, error) {
				response, err := api.RemoveAgentConversationViewContext(ctx, "C9876543210", RemoveAgentConversationViewOptionViewKey("plan"))
				if response == nil {
					return nil, err
				}
				return &response.SlackResponse, err
			},
		},
	} {
		t.Run(method.name, func(t *testing.T) {
			for _, tt := range []struct {
				name   string
				status int
				body   string
			}{
				{"Slack error", http.StatusOK, `{"ok":false,"error":"` + method.slackError + `","response_metadata":{"messages":["Request rejected"]}}`},
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
					response, err := method.call(api, context.Background())
					require.Error(t, err)
					if tt.name == "Slack error" {
						var slackErr SlackErrorResponse
						require.ErrorAs(t, err, &slackErr)
						assert.Equal(t, method.slackError, slackErr.Err)
						assert.Equal(t, []string{"Request rejected"}, slackErr.ResponseMetadata.Messages)
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
				response, err := method.call(api, ctx)
				require.ErrorIs(t, err, context.Canceled)
				assert.Nil(t, response)
			})
		})
	}
}
