package main

import (
	"context"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"

	"github.com/slack-go/slack"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSetViewPullRequestAvailability(t *testing.T) {
	for _, tt := range []struct {
		name     string
		viewType string
		body     string
		apiError string
		wantHint bool
	}{
		{
			name: "rejected pull request type gets a workaround", viewType: slack.AgentViewTypePullRequest,
			body:     `{"ok":false,"error":"invalid_view_type","response_metadata":{"messages":["Type rejected"]}}`,
			apiError: "invalid_view_type", wantHint: true,
		},
		{
			name: "other pull request errors are preserved", viewType: slack.AgentViewTypePullRequest,
			body: `{"ok":false,"error":"invalid_pr_url"}`, apiError: "invalid_pr_url",
		},
		{
			name: "other rejected types get no pull request hint", viewType: slack.AgentViewTypeHTML,
			body: `{"ok":false,"error":"invalid_view_type"}`, apiError: "invalid_view_type",
		},
		{
			name: "accepted pull request response is preserved", viewType: slack.AgentViewTypePullRequest,
			body: `{"ok":true,"channel_id":"C_CODE","view_id":"Ct_VIEW","type":"pull_request"}`,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			var calls atomic.Int32
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				calls.Add(1)
				assert.Equal(t, "/agents.conversations.setView", r.URL.Path)
				require.NoError(t, r.ParseForm())
				assert.Equal(t, tt.viewType, r.PostForm.Get("type"))
				if tt.viewType == slack.AgentViewTypePullRequest {
					assert.Equal(t, "https://github.com/slack-go/slack/pull/1623", r.PostForm.Get("pr_url"))
				}
				w.Header().Set("Content-Type", "application/json")
				_, _ = w.Write([]byte(tt.body))
			}))
			defer server.Close()
			api := slack.New("testing-token", slack.OptionAPIURL(server.URL+"/"))
			response, err := setView(context.Background(), api, "C_CODE", tt.viewType,
				slack.SetAgentConversationViewOptionType(tt.viewType),
				slack.SetAgentConversationViewOptionPRURL("https://github.com/slack-go/slack/pull/1623"))
			assert.Equal(t, int32(1), calls.Load(), "a rejected request must not automatically publish another view")
			require.NotNil(t, response)
			if tt.apiError == "" {
				require.NoError(t, err)
				assert.Equal(t, "Ct_VIEW", response.ViewID)
				return
			}
			require.Error(t, err)
			var slackErr slack.SlackErrorResponse
			require.ErrorAs(t, err, &slackErr)
			assert.Equal(t, tt.apiError, slackErr.Err)
			if tt.wantHint {
				assert.Contains(t, err.Error(), "-type diff -content-file")
				assert.Equal(t, []string{"Type rejected"}, slackErr.ResponseMetadata.Messages)
			} else {
				assert.NotContains(t, err.Error(), "-type diff")
			}
		})
	}
}
