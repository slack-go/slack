package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/slack-go/slack"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMCPWorkflow(t *testing.T) {
	for _, tt := range []struct {
		name, wantError                                      string
		ignore, reject, wrongServer, failCreate, failCleanup bool
	}{
		{name: "app and server declaration"},
		{name: "declaration rejected", reject: true, wantError: "mcp_server_not_found"},
		{name: "ok without storage is not verification", ignore: true, wantError: "expected 1 MCP declarations, info returned 0"},
		{name: "partial create failure cleans both skills", failCreate: true, wantError: "create with declarations"},
		{name: "wrong server is not verification", wrongServer: true, wantError: "info did not return declaration"},
		{name: "cleanup error prevents success", failCleanup: true, wantError: "cleanup temporary skill"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			skills := map[string][]slack.SlackbotSkillMCPServer{}
			created, deletes, nonEmptyUpdates := 0, 0, 0
			acceptedDeclaration := ""
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if !assert.NoError(t, r.ParseForm()) {
					return
				}
				assert.Equal(t, "application/x-www-form-urlencoded", r.Header.Get("Content-Type"))
				result := map[string]any{"ok": true}
				id := r.PostForm.Get("skill_id")
				declared := []slack.SlackbotSkillMCPServer{}
				if raw, present := r.PostForm["mcp_servers"]; present {
					require.Len(t, raw, 1)
					var entries []json.RawMessage
					if !assert.NoError(t, json.Unmarshal([]byte(raw[0]), &entries)) {
						return
					}
					if len(entries) > 0 {
						if r.URL.Path == "/slackbot.skills.update" {
							nonEmptyUpdates++
						}
						var object map[string]any
						if !assert.NoError(t, json.Unmarshal(entries[0], &object)) {
							return
						}
						assert.Equal(t, "app", object["service_type"])
						assert.Equal(t, "A_TEST", object["app_id"])
						assert.Equal(t, "Amcp_TEST", object["server_id"])
						if tt.reject {
							// Simulate a failed update that partially changes state.
							// Cleanup must still delete the temporary skill.
							skills[id] = []slack.SlackbotSkillMCPServer{{AppID: "A_TEST"}}
							result = map[string]any{"ok": false, "error": "mcp_server_not_found", "invalid_mcp_servers": entries}
							w.Header().Set("Content-Type", "application/json")
							assert.NoError(t, json.NewEncoder(w).Encode(result))
							return
						}
						assert.Empty(t, skills[id], "write must start from an empty state")
						if !tt.ignore {
							declared = []slack.SlackbotSkillMCPServer{{ServiceType: "app", AppID: "A_TEST", ServerID: "Amcp_TEST", DisplayName: "Test MCP", IsAvailable: true}}
							if acceptedDeclaration == "" {
								acceptedDeclaration = raw[0]
							} else {
								assert.JSONEq(t, acceptedDeclaration, raw[0], "create must use the exact verified declaration")
							}
						}
					}
				}
				switch r.URL.Path {
				case "/slackbot.skills.create":
					created++
					id = fmt.Sprintf("F%d", created)
					skills[id] = declared
					if tt.failCreate && created > 1 {
						result["ok"] = false
						result["error"] = "internal_error"
					}
					result["skill"] = map[string]any{"id": id, "mcp_servers": declared}
				case "/slackbot.skills.update":
					skills[id] = declared
					result["skill"] = map[string]any{"id": id, "mcp_servers": declared}
				case "/slackbot.skills.info":
					if tt.wrongServer && len(skills[id]) > 0 {
						skills[id][0].ServerID = "Amcp_WRONG"
					}
					result["skill"] = map[string]any{"id": id, "mcp_servers": skills[id]}
				case "/slackbot.skills.delete":
					deletes++
					delete(skills, id)
					if tt.failCleanup {
						result["ok"] = false
						result["error"] = "restricted_action"
					}
				default:
					t.Errorf("unexpected method %s", r.URL.Path)
				}
				w.Header().Set("Content-Type", "application/json")
				assert.NoError(t, json.NewEncoder(w).Encode(result))
			}))
			defer server.Close()
			declarations := []json.RawMessage{json.RawMessage(`{"service_type":"app","app_id":"A_TEST","server_id":"Amcp_TEST"}`)}
			err := runMCPWorkflow(context.Background(), slack.New("testing", slack.OptionAPIURL(server.URL+"/")), "A_TEST", declarations)
			if tt.wantError == "" {
				require.NoError(t, err)
			} else {
				require.ErrorContains(t, err, tt.wantError)
			}
			assert.Equal(t, created, deletes, "cleanup must attempt deletion of every created skill")
			assert.Empty(t, skills)
			assert.Positive(t, nonEmptyUpdates)
		})
	}
}

func TestMCPWorkflowRejectsEmptyInput(t *testing.T) {
	for _, supplied := range [][]json.RawMessage{nil, {}} {
		err := runMCPWorkflow(context.Background(), slack.New("unused"), "A_TEST", supplied)
		require.ErrorContains(t, err, "non-empty declarations")
	}
}
