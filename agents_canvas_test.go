package slack

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetAgentConversationCanvas(t *testing.T) {
	for _, tt := range []struct {
		name    string
		options []GetAgentConversationCanvasOption
		want    url.Values
	}{
		{name: "server defaults", want: url.Values{}},
		{
			name: "all options",
			options: []GetAgentConversationCanvasOption{
				GetAgentConversationCanvasOptionContentFormat("markdown"),
				GetAgentConversationCanvasOptionIncludeResolved(true),
			},
			want: url.Values{"content_format": {"markdown"}, "include_resolved": {"true"}},
		},
		{
			name:    "explicit false",
			options: []GetAgentConversationCanvasOption{GetAgentConversationCanvasOptionIncludeResolved(false)},
			want:    url.Values{"include_resolved": {"false"}},
		},
		{
			name: "nil and empty options",
			options: []GetAgentConversationCanvasOption{
				nil,
				GetAgentConversationCanvasOptionContentFormat(""),
			},
			want: url.Values{},
		},
		{
			name: "last option wins",
			options: []GetAgentConversationCanvasOption{
				GetAgentConversationCanvasOptionContentFormat("other"),
				GetAgentConversationCanvasOptionContentFormat("markdown"),
				GetAgentConversationCanvasOptionIncludeResolved(true),
				GetAgentConversationCanvasOptionIncludeResolved(false),
			},
			want: url.Values{"content_format": {"markdown"}, "include_resolved": {"false"}},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, http.MethodPost, r.Method)
				assert.Equal(t, "/agents.conversations.getCanvas", r.URL.Path)
				require.NoError(t, r.ParseForm())
				assert.Equal(t, "testing-token", r.PostForm.Get("token"))
				assert.Equal(t, "C9876543210", r.PostForm.Get("channel"))
				assert.Equal(t, "F1234567890", r.PostForm.Get("canvas_id"))
				r.PostForm.Del("token")
				r.PostForm.Del("channel")
				r.PostForm.Del("canvas_id")
				assert.Equal(t, tt.want, r.PostForm)
				w.Header().Set("Content-Type", "application/json")
				_, _ = w.Write([]byte(`{
					"ok":true,"canvas_id":"F1234567890","title":"Migration plan",
					"content":"# Migration plan\n\n1. Inventory cron jobs…",
					"comments":[{
						"id":"An12345","ts":"1717180000.000100","user_id":"U2345678901",
						"text":"Can we do the billing jobs last?","quoted_text":"1. Inventory cron jobs",
						"is_resolved":false,"has_more_replies":true,
						"replies":[{"ts":"1717180100.000200","user_id":"U3456789012","text":"+1, those are the riskiest"}]
					},{
						"id":"An67890","ts":"1717180200.000300","user_id":"U2345678901",
						"text":"Looks good","quoted_text":"Migration plan","is_resolved":true,
						"replies":[],"has_more_replies":false
					}],"has_more_comments":true
				}`))
			}))
			defer server.Close()
			api := New("testing-token", OptionAPIURL(server.URL+"/"))
			response, err := api.GetAgentConversationCanvas("C9876543210", "F1234567890", tt.options...)
			require.NoError(t, err)
			require.NotNil(t, response)
			assert.True(t, response.Ok)
			assert.Equal(t, "F1234567890", response.CanvasID)
			assert.Equal(t, "Migration plan", response.Title)
			assert.Equal(t, "# Migration plan\n\n1. Inventory cron jobs…", response.Content)
			assert.True(t, response.HasMoreComments)
			assert.Equal(t, []AgentCanvasComment{
				{
					ID: "An12345", Timestamp: "1717180000.000100", UserID: "U2345678901",
					Text: "Can we do the billing jobs last?", QuotedText: "1. Inventory cron jobs",
					HasMoreReplies: true,
					Replies: []AgentCanvasCommentReply{{
						Timestamp: "1717180100.000200", UserID: "U3456789012", Text: "+1, those are the riskiest",
					}},
				},
				{
					ID: "An67890", Timestamp: "1717180200.000300", UserID: "U2345678901",
					Text: "Looks good", QuotedText: "Migration plan", IsResolved: true,
					Replies: []AgentCanvasCommentReply{},
				},
			}, response.Comments)
		})
	}
}

func TestGetAgentConversationCanvasWithoutComments(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"ok":true,"canvas_id":"F1234567890","title":"Empty plan","content":"","comments":[],"has_more_comments":false}`))
	}))
	defer server.Close()
	api := New("testing-token", OptionAPIURL(server.URL+"/"))
	response, err := api.GetAgentConversationCanvas("C9876543210", "F1234567890")
	require.NoError(t, err)
	require.NotNil(t, response)
	assert.Empty(t, response.Content)
	assert.Equal(t, []AgentCanvasComment{}, response.Comments)
	assert.False(t, response.HasMoreComments)
}

func TestSetAgentConversationCanvasContent(t *testing.T) {
	const markdown = "# Migration plan\n\n1. Inventory cron jobs\n2. Port billing jobs last\n"
	for _, tt := range []struct {
		name         string
		content      string
		changedCount int
	}{
		{name: "revised markdown", content: markdown, changedCount: 2},
		{name: "no changes", content: markdown, changedCount: 0},
		{name: "empty content is sent", content: "", changedCount: 1},
	} {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, http.MethodPost, r.Method)
				assert.Equal(t, "/agents.conversations.setCanvasContent", r.URL.Path)
				require.NoError(t, r.ParseForm())
				assert.Equal(t, url.Values{
					"token": {"testing-token"}, "channel": {"C9876543210"},
					"canvas_id": {"F1234567890"}, "content": {tt.content},
				}, r.PostForm)
				w.Header().Set("Content-Type", "application/json")
				_, _ = fmt.Fprintf(w, `{"ok":true,"canvas_id":"F1234567890","sections_changed_count":%d}`, tt.changedCount)
			}))
			defer server.Close()
			api := New("testing-token", OptionAPIURL(server.URL+"/"))
			response, err := api.SetAgentConversationCanvasContent("C9876543210", "F1234567890", tt.content)
			require.NoError(t, err)
			require.NotNil(t, response)
			assert.True(t, response.Ok)
			assert.Equal(t, "F1234567890", response.CanvasID)
			assert.Equal(t, tt.changedCount, response.SectionsChangedCount)
		})
	}
}

func TestAgentConversationCanvasErrors(t *testing.T) {
	for _, method := range []struct {
		name       string
		slackError string
		call       func(*Client, context.Context) (*SlackResponse, error)
	}{
		{
			name: "get", slackError: "comments_fetch_failed",
			call: func(api *Client, ctx context.Context) (*SlackResponse, error) {
				response, err := api.GetAgentConversationCanvasContext(ctx, "C9876543210", "F1234567890")
				if response == nil {
					return nil, err
				}
				return &response.SlackResponse, err
			},
		},
		{
			name: "set content", slackError: "canvas_content_update_failed",
			call: func(api *Client, ctx context.Context) (*SlackResponse, error) {
				response, err := api.SetAgentConversationCanvasContentContext(ctx, "C9876543210", "F1234567890", "# Plan")
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
