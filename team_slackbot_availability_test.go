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

func TestTeamSlackbotAvailability(t *testing.T) {
	for _, teamID := range []string{"", "T123"} {
		for _, tt := range []struct {
			name    string
			options []GetTeamInfoOption
			want    url.Values
		}{
			{"default", nil, url.Values{}},
			{"enabled", []GetTeamInfoOption{GetTeamInfoOptionIncludeSlackbotAvailability(true)}, url.Values{"include_slackbot_availability": {"true"}}},
			{"disabled", []GetTeamInfoOption{GetTeamInfoOptionIncludeSlackbotAvailability(false)}, url.Values{"include_slackbot_availability": {"false"}}},
			{"last option wins", []GetTeamInfoOption{nil, GetTeamInfoOptionIncludeSlackbotAvailability(true), GetTeamInfoOptionIncludeSlackbotAvailability(false)}, url.Values{"include_slackbot_availability": {"false"}}},
		} {
			t.Run(teamID+"/"+tt.name, func(t *testing.T) {
				server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					assert.Equal(t, "/team.info", r.URL.Path)
					if !assert.NoError(t, r.ParseForm()) {
						return
					}
					assert.Equal(t, "testing-token", r.PostForm.Get("token"))
					if teamID == "" {
						assert.NotContains(t, r.PostForm, "team")
					} else {
						assert.Equal(t, teamID, r.PostForm.Get("team"))
					}
					r.PostForm.Del("token")
					r.PostForm.Del("team")
					assert.Equal(t, tt.want, r.PostForm)
					_, _ = w.Write([]byte(`{"ok":true,"team":{"id":"T123","slackbot_availability":"restricted"}}`))
				}))
				defer server.Close()
				api := New("testing-token", OptionAPIURL(server.URL+"/"))
				var team *TeamInfo
				var err error
				if teamID == "" {
					team, err = api.GetTeamInfo(tt.options...)
				} else {
					team, err = api.GetOtherTeamInfo(teamID, tt.options...)
				}
				require.NoError(t, err)
				require.NotNil(t, team)
				assert.Equal(t, "restricted", team.SlackbotAvailability)
			})
		}
	}
}

func TestTeamSlackbotAvailabilityValues(t *testing.T) {
	for _, value := range []string{"available", "restricted", "unavailable", "unknown", "future_value"} {
		var team TeamInfo
		require.NoError(t, json.Unmarshal([]byte(`{"slackbot_availability":"`+value+`"}`), &team))
		assert.Equal(t, value, team.SlackbotAvailability)
	}
	encoded, err := json.Marshal(TeamInfo{})
	require.NoError(t, err)
	assert.NotContains(t, string(encoded), "slackbot_availability")
}

func TestTeamInfoErrors(t *testing.T) {
	for _, method := range []struct {
		name string
		call func(*Client, context.Context) error
	}{
		{"team info", func(api *Client, ctx context.Context) error { _, err := api.GetTeamInfoContext(ctx, nil); return err }},
		{"other team info", func(api *Client, ctx context.Context) error {
			_, err := api.GetOtherTeamInfoContext(ctx, "T123", nil)
			return err
		}},
	} {
		t.Run(method.name, func(t *testing.T) {
			for _, tt := range []struct {
				name, body string
				status     int
			}{
				{"Slack error", `{"ok":false,"error":"access_denied","response_metadata":{"messages":["Request rejected"]}}`, http.StatusOK},
				{"invalid JSON", `{"ok":`, http.StatusOK},
				{"HTTP error", `server error`, http.StatusInternalServerError},
			} {
				t.Run(tt.name, func(t *testing.T) {
					server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
						w.WriteHeader(tt.status)
						_, _ = w.Write([]byte(tt.body))
					}))
					defer server.Close()
					api := New("testing-token", OptionAPIURL(server.URL+"/"))
					err := method.call(api, context.Background())
					require.Error(t, err)
					if tt.name == "Slack error" {
						var slackErr SlackErrorResponse
						require.ErrorAs(t, err, &slackErr)
						assert.Equal(t, "access_denied", slackErr.Err)
						assert.Equal(t, []string{"Request rejected"}, slackErr.ResponseMetadata.Messages)
					}
				})
			}
			t.Run("canceled context", func(t *testing.T) {
				server := httptest.NewServer(http.HandlerFunc(func(http.ResponseWriter, *http.Request) { t.Error("canceled request reached server") }))
				defer server.Close()
				api := New("testing-token", OptionAPIURL(server.URL+"/"))
				ctx, cancel := context.WithCancel(context.Background())
				cancel()
				require.ErrorIs(t, method.call(api, ctx), context.Canceled)
			})
		})
	}
}
