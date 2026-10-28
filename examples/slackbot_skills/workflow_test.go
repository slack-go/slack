package main

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/slack-go/slack"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWorkflow(t *testing.T) {
	for _, failUpdate := range []bool{false, true} {
		t.Run(map[bool]string{false: "complete lifecycle", true: "cleanup after failure"}[failUpdate], func(t *testing.T) {
			exists, enabled, deleted := false, false, false
			instructions, description := "", ""
			seen := map[string]bool{}
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, "application/x-www-form-urlencoded", r.Header.Get("Content-Type"))
				if !assert.NoError(t, r.ParseForm()) {
					return
				}
				seen[r.URL.Path] = true
				result := map[string]any{"ok": true}
				switch r.URL.Path {
				case "/slackbot.skills.create":
					exists = true
					enabled = true
					instructions = r.PostForm.Get("instructions")
				case "/slackbot.skills.update":
					if failUpdate {
						result = map[string]any{"ok": false, "error": "restricted_action"}
						break
					}
					instructions = r.PostForm.Get("instructions")
					description = r.PostForm.Get("description")
				case "/slackbot.skills.info":
					if !exists {
						result = map[string]any{"ok": false, "error": "skill_not_found"}
					}
				case "/slackbot.skills.list":
					// Every traversal starts with an empty page carrying a cursor.
					cursor := "next"
					skills := []map[string]any{}
					if r.PostForm.Get("cursor") == "next" {
						cursor = ""
						if r.PostForm.Get("filter") == "owned" || enabled {
							skills = append(skills, map[string]any{"id": "F_TEST"})
						}
					}
					result["skills"] = skills
					result["response_metadata"] = map[string]string{"next_cursor": cursor}
				case "/slackbot.skills.library.add":
					enabled = true
				case "/slackbot.skills.library.remove":
					enabled = false
				case "/slackbot.skills.delete":
					exists = false
					deleted = true
				default:
					t.Errorf("unexpected method %s", r.URL.Path)
				}
				if exists {
					result["skill"] = map[string]any{"id": "F_TEST", "instructions": instructions, "description": description}
				}
				w.Header().Set("Content-Type", "application/json")
				assert.NoError(t, json.NewEncoder(w).Encode(result))
			}))
			defer server.Close()
			err := runWorkflow(context.Background(), slack.New("test", slack.OptionAPIURL(server.URL+"/")))
			if failUpdate {
				require.ErrorContains(t, err, "restricted_action")
			} else {
				require.NoError(t, err)
				assert.Len(t, seen, 7)
			}
			assert.True(t, deleted, "fixture must be cleaned up on success and failure")
		})
	}
}

func TestWorkflowRevokedToken(t *testing.T) {
	calls := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls++
		assert.Equal(t, "/slackbot.skills.create", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"ok":false,"error":"token_revoked"}`))
	}))
	defer server.Close()
	err := runWorkflow(context.Background(), slack.New("revoked", slack.OptionAPIURL(server.URL+"/")))
	var slackErr slack.SlackErrorResponse
	require.True(t, errors.As(err, &slackErr))
	assert.Equal(t, "token_revoked", slackErr.Err)
	assert.Equal(t, 1, calls, "no fixture ID means no subsequent API or cleanup calls")
}
