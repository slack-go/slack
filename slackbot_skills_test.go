package slack

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// skillTestClient checks the wire request independently of the option helpers.
func skillTestClient(t *testing.T, method, wantBody, response string) *Client {
	t.Helper()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/"+method, r.URL.Path)
		assert.Equal(t, "application/x-www-form-urlencoded", r.Header.Get("Content-Type"))
		if !assert.NoError(t, r.ParseForm()) {
			return
		}
		assert.Equal(t, "testing-token", r.PostForm.Get("token"))
		// Expected arguments are written as JSON in fixtures for readability.
		// Check each actual form field, including presence and single encoding.
		var want map[string]json.RawMessage
		if !assert.NoError(t, json.Unmarshal([]byte(wantBody), &want)) {
			return
		}
		assert.Len(t, r.PostForm, len(want)+1)
		for key, raw := range want {
			if !assert.Len(t, r.PostForm[key], 1, "form field %s", key) {
				continue
			}
			var text string
			if json.Unmarshal(raw, &text) == nil {
				assert.Equal(t, text, r.PostForm.Get(key), "form field %s", key)
			} else {
				assert.JSONEq(t, string(raw), r.PostForm.Get(key), "form field %s", key)
			}
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(response))
	}))
	t.Cleanup(server.Close)
	return New("testing-token", OptionAPIURL(server.URL+"/"))
}

func TestCreateSlackbotSkill(t *testing.T) {
	for _, tt := range []struct {
		name    string
		options []CreateSlackbotSkillOption
		want    string
	}{
		{"markdown and all metadata", []CreateSlackbotSkillOption{
			nil, CreateSlackbotSkillOptionInstructions("# Review"),
			CreateSlackbotSkillOptionDescription("Release risks"),
			CreateSlackbotSkillOptionCategory("engineering"),
			CreateSlackbotSkillOptionExamplePrompt("Run {skill}"),
			CreateSlackbotSkillOptionMCPServers(json.RawMessage(`{"app_id":"A1"}`)),
		}, `{"display_name":"Review","instructions":"# Review","description":"Release risks","category":"engineering","example_prompt":"Run {skill}","mcp_servers":[{"app_id":"A1"}]}`},
		{"canvas and explicit empty array", []CreateSlackbotSkillOption{
			CreateSlackbotSkillOptionCanvasID("F1"), CreateSlackbotSkillOptionMCPServers(),
		}, `{"display_name":"Review","canvas_id":"F1","mcp_servers":[]}`},
		{"last option wins and empty values survive", []CreateSlackbotSkillOption{
			CreateSlackbotSkillOptionInstructions("old"), CreateSlackbotSkillOptionInstructions(""),
		}, `{"display_name":"Review","instructions":""}`},
	} {
		t.Run(tt.name, func(t *testing.T) {
			api := skillTestClient(t, "slackbot.skills.create", tt.want, `{"ok":true,"skill":{"id":"F1","type":"canvas","name":null,"instructions":null,"date_updated":null,"mcp_servers":[{"service_type":"app","app_id":"A1","server_id":"M1","display_name":"Calendar","is_available":true}]}}`)
			response, err := api.CreateSlackbotSkill("Review", tt.options...)
			require.NoError(t, err)
			assert.Equal(t, "F1", response.Skill.ID)
			assert.Nil(t, response.Skill.Name)
			assert.Nil(t, response.Skill.Instructions)
			assert.Nil(t, response.Skill.DateUpdated)
			require.Len(t, response.Skill.MCPServers, 1)
			assert.True(t, response.Skill.MCPServers[0].IsAvailable)
		})
	}
}

func TestCreateSlackbotSkillInvalidMCPServers(t *testing.T) {
	api := skillTestClient(t, "slackbot.skills.create", `{"display_name":"Review"}`, `{"ok":false,"error":"mcp_server_not_found","invalid_mcp_servers":[{"app_id":"A1"}],"needed":"skills:write"}`)
	response, err := api.CreateSlackbotSkill("Review")
	require.EqualError(t, err, "mcp_server_not_found")
	require.NotNil(t, response)
	assert.JSONEq(t, `[{"app_id":"A1"}]`, string(response.InvalidMCPServers))
	var slackErr SlackErrorResponse
	require.ErrorAs(t, err, &slackErr)
	assert.Equal(t, "skills:write", slackErr.Needed)
	body, err := json.Marshal(response)
	require.NoError(t, err)
	var decoded map[string]json.RawMessage
	require.NoError(t, json.Unmarshal(body, &decoded))
	assert.NotContains(t, decoded, "skill", "failed creation must not print a fabricated empty skill")
	assert.Contains(t, decoded, "invalid_mcp_servers", "structured error details remain available")
}

func TestCreateSlackbotSkillInvalidJSONOption(t *testing.T) {
	api := New("testing-token", OptionAPIURL("http://127.0.0.1:1/"))
	response, err := api.CreateSlackbotSkill("Review", CreateSlackbotSkillOptionMCPServers(json.RawMessage(`{`)))
	require.Error(t, err)
	assert.Nil(t, response)
}

func TestSlackbotSkillsFormUnicodeAndEscaping(t *testing.T) {
	const name = "Revisión 東京 🚀 + & = % ? #"
	const instructions = "# Revisión\n\nKeep + & = % ? # intact.\n[link](https://example.com/?a=1&b=a+b%20c#section)\n```\nx = \"a&b+c\"\n```"
	const app = `{"app_id":"A1","label":"東京 + & = % 🚀","nested":["a+b","c&d"]}`
	want, err := json.Marshal(map[string]any{
		"display_name": name, "instructions": instructions,
		"description": "", "mcp_servers": []json.RawMessage{json.RawMessage(app)},
	})
	require.NoError(t, err)
	api := skillTestClient(t, "slackbot.skills.create", string(want), `{"ok":true,"skill":{"id":"F1"}}`)
	_, err = api.CreateSlackbotSkill(name,
		CreateSlackbotSkillOptionInstructions(instructions),
		CreateSlackbotSkillOptionDescription(""),
		CreateSlackbotSkillOptionMCPServers(json.RawMessage(app)))
	require.NoError(t, err)
}

func TestSlackbotSkillsFormWarningRedactsToken(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !assert.NoError(t, r.ParseForm()) {
			return
		}
		assert.Equal(t, "testing-token", r.PostForm.Get("token"))
		assert.Equal(t, "", r.PostForm.Get("description"))
		assert.Contains(t, r.PostForm, "description")
		assert.NotContains(t, r.PostForm, "instructions")
		assert.Equal(t, "[]", r.PostForm.Get("mcp_servers"))
		_, _ = w.Write([]byte(`{"ok":true,"warning":"test_warning","skill":{"id":"F1"}}`))
	}))
	defer server.Close()
	called := false
	api := New("testing-token", OptionAPIURL(server.URL+"/"), OptionOnWarning(func(path string, request any, warning *Warning) {
		called = true
		assert.Equal(t, "slackbot.skills.update", path)
		values, ok := request.(url.Values)
		require.True(t, ok)
		assert.NotContains(t, values, "token")
		assert.Equal(t, "[]", values.Get("mcp_servers"))
		assert.Equal(t, []string{"test_warning"}, warning.Codes)
	}))
	_, err := api.UpdateSlackbotSkill("F1", UpdateSlackbotSkillOptionDescription(""), UpdateSlackbotSkillOptionMCPServers())
	require.NoError(t, err)
	assert.True(t, called)
}

func TestCreateSlackbotSkillCanceledContext(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	response, err := New("testing-token").CreateSlackbotSkillContext(ctx, "Review")
	require.ErrorIs(t, err, context.Canceled)
	assert.Nil(t, response)
}

func TestGetSlackbotSkillInfo(t *testing.T) {
	api := skillTestClient(t, "slackbot.skills.info", `{"skill_id":"builtin:review"}`, `{"ok":true,"skill":{"id":"builtin:review","display_name":"Review","instructions":"# Review","access":"read","date_created":1790000000,"date_updated":1790100000,"usage_count":42,"library_count":3,"permalink":"https://example.slack.com/skill"}}`)
	response, err := api.GetSlackbotSkillInfo("builtin:review")
	require.NoError(t, err)
	require.NotNil(t, response.Skill.Instructions)
	assert.Equal(t, "# Review", *response.Skill.Instructions)
	require.NotNil(t, response.Skill.DateUpdated)
	assert.EqualValues(t, 1790100000, *response.Skill.DateUpdated)
	assert.Equal(t, "read", response.Skill.Access)
	assert.Equal(t, 42, response.Skill.UsageCount)
}

func TestUpdateSlackbotSkill(t *testing.T) {
	for _, tt := range []struct {
		name    string
		options []UpdateSlackbotSkillOption
		want    string
	}{
		{"all fields", []UpdateSlackbotSkillOption{
			nil, UpdateSlackbotSkillOptionDisplayName("Updated"),
			UpdateSlackbotSkillOptionInstructions("# Updated"),
			UpdateSlackbotSkillOptionDescription("Description"),
			UpdateSlackbotSkillOptionCategory("productivity"),
			UpdateSlackbotSkillOptionExamplePrompt("Use {skill}"),
			UpdateSlackbotSkillOptionMCPServers(json.RawMessage(`{"app_id":"A1"}`)),
		}, `{"skill_id":"F1","display_name":"Updated","instructions":"# Updated","description":"Description","category":"productivity","example_prompt":"Use {skill}","mcp_servers":[{"app_id":"A1"}]}`},
		{"clears", []UpdateSlackbotSkillOption{
			UpdateSlackbotSkillOptionDescription("old"), UpdateSlackbotSkillOptionDescription(""),
			UpdateSlackbotSkillOptionCategory(""), UpdateSlackbotSkillOptionExamplePrompt(""),
			UpdateSlackbotSkillOptionMCPServers(),
		}, `{"skill_id":"F1","description":"","category":"","example_prompt":"","mcp_servers":[]}`},
		{"only supplied field changes", []UpdateSlackbotSkillOption{
			UpdateSlackbotSkillOptionInstructions(""),
		}, `{"skill_id":"F1","instructions":""}`},
	} {
		t.Run(tt.name, func(t *testing.T) {
			api := skillTestClient(t, "slackbot.skills.update", tt.want, `{"ok":true,"skill":{"id":"F1","instructions":null}}`)
			response, err := api.UpdateSlackbotSkill("F1", tt.options...)
			require.NoError(t, err)
			assert.Equal(t, "F1", response.Skill.ID)
		})
	}
}

func TestListSlackbotSkills(t *testing.T) {
	api := skillTestClient(t, "slackbot.skills.list", `{"filter":"owned","category":"engineering","cursor":"next","limit":0}`, `{"ok":true,"skills":[],"response_metadata":{"next_cursor":"continue"}}`)
	response, err := api.ListSlackbotSkills(nil,
		ListSlackbotSkillsOptionFilter("catalog"), ListSlackbotSkillsOptionFilter("owned"),
		ListSlackbotSkillsOptionCategory("engineering"), ListSlackbotSkillsOptionCursor("next"),
		ListSlackbotSkillsOptionLimit(0))
	require.NoError(t, err)
	assert.Empty(t, response.Skills)
	assert.Equal(t, "continue", response.ResponseMetadata.Cursor)

	api = skillTestClient(t, "slackbot.skills.list", `{}`, `{"ok":true,"skills":[{"id":"F1","instructions":null}],"response_metadata":{"next_cursor":""}}`)
	response, err = api.ListSlackbotSkills()
	require.NoError(t, err)
	require.Len(t, response.Skills, 1)
	assert.Nil(t, response.Skills[0].Instructions)
	assert.Empty(t, response.ResponseMetadata.Cursor)
}

func TestAddSlackbotSkillToLibrary(t *testing.T) {
	api := skillTestClient(t, "slackbot.skills.library.add", `{"skill_id":"F1"}`, `{"ok":true}`)
	require.NoError(t, api.AddSlackbotSkillToLibrary("F1"))
	api = skillTestClient(t, "slackbot.skills.library.add", `{"skill_id":"F1"}`, `{"ok":false,"error":"restricted_action"}`)
	require.EqualError(t, api.AddSlackbotSkillToLibraryContext(context.Background(), "F1"), "restricted_action")
}

func TestRemoveSlackbotSkillFromLibrary(t *testing.T) {
	api := skillTestClient(t, "slackbot.skills.library.remove", `{"skill_id":"F1"}`, `{"ok":true}`)
	require.NoError(t, api.RemoveSlackbotSkillFromLibrary("F1"))
	api = skillTestClient(t, "slackbot.skills.library.remove", `{"skill_id":"F1"}`, `{"ok":false,"error":"skill_not_found"}`)
	require.EqualError(t, api.RemoveSlackbotSkillFromLibraryContext(context.Background(), "F1"), "skill_not_found")
}

func TestDeleteSlackbotSkill(t *testing.T) {
	api := skillTestClient(t, "slackbot.skills.delete", `{"skill_id":"F1"}`, `{"ok":true}`)
	require.NoError(t, api.DeleteSlackbotSkill("F1"))
	api = skillTestClient(t, "slackbot.skills.delete", `{"skill_id":"builtin:test"}`, `{"ok":false,"error":"method_not_supported_for_skill_type"}`)
	require.EqualError(t, api.DeleteSlackbotSkillContext(context.Background(), "builtin:test"), "method_not_supported_for_skill_type")
}

// All endpoints must preserve Slack errors, propagate transport/decode errors,
// respect cancellation, and use the client's existing rate-limit handling.
func TestSlackbotSkillsFailures(t *testing.T) {
	calls := []struct {
		name     string
		response bool
		call     func(context.Context, *Client) (any, error)
	}{
		{"create", true, func(ctx context.Context, api *Client) (any, error) {
			return api.CreateSlackbotSkillContext(ctx, "Review", CreateSlackbotSkillOptionInstructions("# Review"))
		}},
		{"info", true, func(ctx context.Context, api *Client) (any, error) { return api.GetSlackbotSkillInfoContext(ctx, "F1") }},
		{"update", true, func(ctx context.Context, api *Client) (any, error) {
			return api.UpdateSlackbotSkillContext(ctx, "F1", UpdateSlackbotSkillOptionDescription(""))
		}},
		{"list", true, func(ctx context.Context, api *Client) (any, error) { return api.ListSlackbotSkillsContext(ctx) }},
		{"library.add", false, func(ctx context.Context, api *Client) (any, error) {
			return nil, api.AddSlackbotSkillToLibraryContext(ctx, "F1")
		}},
		{"library.remove", false, func(ctx context.Context, api *Client) (any, error) {
			return nil, api.RemoveSlackbotSkillFromLibraryContext(ctx, "F1")
		}},
		{"delete", false, func(ctx context.Context, api *Client) (any, error) {
			return nil, api.DeleteSlackbotSkillContext(ctx, "F1")
		}},
		{"admin.list", true, func(ctx context.Context, api *Client) (any, error) { return api.AdminSlackbotSkillsListContext(ctx) }},
		{"admin.publish", true, func(ctx context.Context, api *Client) (any, error) {
			return api.AdminSlackbotSkillsPublishContext(ctx, []string{"F1"})
		}},
		{"admin.setAutoEnabled", true, func(ctx context.Context, api *Client) (any, error) {
			return api.AdminSlackbotSkillsSetAutoEnabledContext(ctx, "F1", "no_one")
		}},
		{"admin.unpublish", true, func(ctx context.Context, api *Client) (any, error) {
			return api.AdminSlackbotSkillsUnpublishContext(ctx, []string{"F1"})
		}},
	}
	for _, call := range calls {
		t.Run(call.name, func(t *testing.T) {
			for _, tt := range []struct {
				name   string
				status int
				body   string
			}{
				{"Slack error", http.StatusOK, `{"ok":false,"error":"missing_scope","needed":"skills:write","provided":"skills:read","response_metadata":{"messages":["scope required"]}}`},
				{"HTTP error", http.StatusInternalServerError, `server error`},
				{"invalid JSON", http.StatusOK, `{"ok":`},
				{"rate limit", http.StatusTooManyRequests, `{"ok":false,"error":"ratelimited"}`},
			} {
				t.Run(tt.name, func(t *testing.T) {
					server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
						w.Header().Set("Content-Type", "application/json")
						w.Header().Set("Retry-After", "2")
						w.WriteHeader(tt.status)
						_, _ = w.Write([]byte(tt.body))
					}))
					defer server.Close()
					response, err := call.call(context.Background(), New("test", OptionAPIURL(server.URL+"/")))
					require.Error(t, err)
					if tt.name == "Slack error" {
						var slackErr SlackErrorResponse
						require.ErrorAs(t, err, &slackErr)
						assert.Equal(t, "skills:write", slackErr.Needed)
						assert.Equal(t, "skills:read", slackErr.Provided)
						assert.Equal(t, []string{"scope required"}, slackErr.ResponseMetadata.Messages)
						if call.response {
							assert.NotNil(t, response)
						}
					} else {
						assert.Nil(t, response)
					}
					if tt.name == "rate limit" {
						var limited *RateLimitedError
						require.ErrorAs(t, err, &limited)
						assert.Equal(t, 2*time.Second, limited.RetryAfter)
					}
				})
			}
			t.Run("canceled context", func(t *testing.T) {
				ctx, cancel := context.WithCancel(context.Background())
				cancel()
				response, err := call.call(ctx, New("test", OptionAPIURL("http://127.0.0.1:1/")))
				require.ErrorIs(t, err, context.Canceled)
				assert.Nil(t, response)
			})
		})
	}
}
