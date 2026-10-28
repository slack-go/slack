package main

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/slack-go/slack"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAdminWorkflow(t *testing.T) {
	for _, failDeploy := range []bool{false, true} {
		t.Run(map[bool]string{false: "complete lifecycle", true: "cleanup published fixture after failure"}[failDeploy], func(t *testing.T) {
			published, deleted := false, false
			permission := "no_one"
			groups := []string{}
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
					result["skill"] = map[string]string{"id": "F_TEST"}
				case "/admin.slackbot.skills.publish":
					published = true
				case "/admin.slackbot.skills.unpublish":
					published = false
				case "/admin.slackbot.skills.setAutoEnabled":
					if failDeploy && r.PostForm.Get("permission_type") == "named_entities" {
						result = map[string]any{"ok": false, "error": "skill_not_visible_to_audience", "usergroup_ids": []string{"S_TEST"}}
						break
					}
					permission = r.PostForm.Get("permission_type")
					groups = []string{}
					if permission == "named_entities" {
						groups = []string{"S_A", "S_Z"}
					}
					result["skill_id"] = "F_TEST"
					result["permission_type"] = permission
					result["usergroup_ids"] = groups
				case "/admin.slackbot.skills.list":
					cursor := "next"
					skills := []map[string]any{}
					if r.PostForm.Get("cursor") == "next" {
						cursor = ""
						inventoryGroups := []string{}
						if permission == "named_entities" {
							inventoryGroups = []string{"S_Z", "S_A"}
						}
						skills = append(skills, map[string]any{"id": "F_TEST", "is_official": published, "auto_enabled": map[string]any{"permission_type": permission, "usergroup_ids": inventoryGroups}})
					}
					result["skills"] = skills
					result["response_metadata"] = map[string]string{"next_cursor": cursor}
				case "/slackbot.skills.delete":
					deleted = true
				default:
					t.Errorf("unexpected method %s", r.URL.Path)
				}
				w.Header().Set("Content-Type", "application/json")
				assert.NoError(t, json.NewEncoder(w).Encode(result))
			}))
			defer server.Close()
			err := runWorkflow(context.Background(), slack.New("test", slack.OptionAPIURL(server.URL+"/")), []string{"S_Z", "S_A", "S_Z"})
			if failDeploy {
				require.ErrorContains(t, err, "skill_not_visible_to_audience")
			} else {
				require.NoError(t, err)
				assert.Len(t, seen, 6)
			}
			assert.True(t, deleted)
			assert.False(t, published)
			assert.Equal(t, "no_one", permission)
		})
	}
}
