package slack

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRequestUserConnectionAuth(t *testing.T) {
	provider := UserConnectionPrimaryProvider{
		Value: "example", Name: NewTextBlockObject(PlainTextType, "Example", false, false),
		Description: NewTextBlockObject(PlainTextType, "Connect your account", false, false),
		IconURL:     "https://example.com/icon.png", AuthURL: "https://example.com/auth",
		DisconnectURL:     "https://example.com/disconnect",
		ConnectedIdentity: NewTextBlockObject(PlainTextType, "alice", false, false),
		ConnectedContext:  NewTextBlockObject(PlainTextType, "Example org", false, false),
		Accessory:         NewOverflowBlockElement("disconnect", NewOptionBlockObject("disconnect", NewTextBlockObject(PlainTextType, "Disconnect", false, false), nil)),
	}
	for _, tt := range []struct {
		name      string
		options   []RequestUserConnectionAuthOption
		want      url.Values
		connected bool
	}{
		{"open disconnected modal", []RequestUserConnectionAuthOption{nil, RequestUserConnectionAuthOptionTriggerID("trigger")}, url.Values{"trigger_id": {"trigger"}}, false},
		{"update connected modal", []RequestUserConnectionAuthOption{RequestUserConnectionAuthOptionUserID("U123")}, url.Values{"user_id": {"U123"}}, true},
		{"both recipients", []RequestUserConnectionAuthOption{RequestUserConnectionAuthOptionTriggerID("trigger"), RequestUserConnectionAuthOptionUserID("U123")}, url.Values{"trigger_id": {"trigger"}, "user_id": {"U123"}}, false},
		{"last option wins", []RequestUserConnectionAuthOption{RequestUserConnectionAuthOptionUserID("old"), RequestUserConnectionAuthOptionUserID("U123")}, url.Values{"user_id": {"U123"}}, false},
	} {
		t.Run(tt.name, func(t *testing.T) {
			provider.IsConnected = tt.connected
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, http.MethodPost, r.Method)
				assert.Equal(t, "/apps.user.connection.requestAuth", r.URL.Path)
				if !assert.NoError(t, r.ParseForm()) {
					return
				}
				assert.Equal(t, "testing-token", r.PostForm.Get("token"))
				var fields map[string]json.RawMessage
				if !assert.NoError(t, json.Unmarshal([]byte(r.PostForm.Get("primary_provider")), &fields)) {
					return
				}
				assert.Equal(t, fmt.Sprint(tt.connected), string(fields["is_connected"]))
				var got UserConnectionPrimaryProvider
				if !assert.NoError(t, json.Unmarshal([]byte(r.PostForm.Get("primary_provider")), &got)) {
					return
				}
				assert.Equal(t, provider, got)
				r.PostForm.Del("token")
				r.PostForm.Del("primary_provider")
				assert.Equal(t, tt.want, r.PostForm)
				_, _ = w.Write([]byte(`{"ok":true}`))
			}))
			defer server.Close()
			api := New("testing-token", OptionAPIURL(server.URL+"/"))
			require.NoError(t, api.RequestUserConnectionAuth(provider, tt.options...))
		})
	}
}

func TestBulkUpdateUserConnections(t *testing.T) {
	for _, tt := range []struct {
		name                    string
		connected, disconnected []string
		want                    url.Values
	}{
		{"connected only", []string{"U1", "U2"}, nil, url.Values{"connected_user_ids": {"U1,U2"}}},
		{"disconnected only", nil, []string{"U3"}, url.Values{"disconnected_user_ids": {"U3"}}},
		{"both", []string{"U1"}, []string{"U3"}, url.Values{"connected_user_ids": {"U1"}, "disconnected_user_ids": {"U3"}}},
	} {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, http.MethodPost, r.Method)
				assert.Equal(t, "/apps.user.connection.bulkUpdate", r.URL.Path)
				if !assert.NoError(t, r.ParseForm()) {
					return
				}
				assert.Equal(t, "testing-token", r.PostForm.Get("token"))
				r.PostForm.Del("token")
				assert.Equal(t, tt.want, r.PostForm)
				_, _ = w.Write([]byte(`{"ok":true,"changed_user_ids":["U1"],"failures":[{"user_id":"U3","error":"user_not_found"}],"response_metadata":{"messages":["Some users failed"]}}`))
			}))
			defer server.Close()
			api := New("testing-token", OptionAPIURL(server.URL+"/"))
			response, err := api.BulkUpdateUserConnections(tt.connected, tt.disconnected)
			require.NoError(t, err)
			require.NotNil(t, response)
			assert.Equal(t, []string{"U1"}, response.ChangedUserIDs)
			assert.Equal(t, []UserConnectionFailure{{UserID: "U3", Error: "user_not_found"}}, response.Failures)
			assert.Equal(t, []string{"Some users failed"}, response.ResponseMetadata.Messages)
		})
	}
}

func TestUserConnectionValidation(t *testing.T) {
	api := New("testing-token", OptionAPIURL(":invalid"))
	for _, options := range [][]RequestUserConnectionAuthOption{nil, {nil}, {RequestUserConnectionAuthOptionUserID("")}} {
		err := api.RequestUserConnectionAuth(UserConnectionPrimaryProvider{}, options...)
		require.EqualError(t, err, "one of trigger_id or user_id is required")
	}
	for _, lists := range []struct{ connected, disconnected []string }{
		{}, {make([]string, 1001), nil}, {make([]string, 500), make([]string, 501)},
	} {
		response, err := api.BulkUpdateUserConnections(lists.connected, lists.disconnected)
		require.EqualError(t, err, "between 1 and 1000 user IDs are required")
		assert.Nil(t, response)
	}
}

func TestBulkUpdateUserConnectionsLimit(t *testing.T) {
	users := make([]string, 1000)
	for i := range users {
		users[i] = fmt.Sprintf("U%d", i)
	}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, strings.Join(users, ","), r.FormValue("connected_user_ids"))
		_, _ = w.Write([]byte(`{"ok":true,"changed_user_ids":[],"failures":[]}`))
	}))
	defer server.Close()
	api := New("testing-token", OptionAPIURL(server.URL+"/"))
	response, err := api.BulkUpdateUserConnections(users, nil)
	require.NoError(t, err)
	require.NotNil(t, response)
	assert.Empty(t, response.Failures)
}

func TestUserConnectionRequestErrors(t *testing.T) {
	for _, method := range []struct {
		name string
		call func(*Client, context.Context) error
	}{
		{"request auth", func(api *Client, ctx context.Context) error {
			return api.RequestUserConnectionAuthContext(ctx, UserConnectionPrimaryProvider{}, RequestUserConnectionAuthOptionUserID("U123"))
		}},
		{"bulk update", func(api *Client, ctx context.Context) error {
			_, err := api.BulkUpdateUserConnectionsContext(ctx, []string{"U123"}, nil)
			return err
		}},
	} {
		t.Run(method.name, func(t *testing.T) {
			for _, tt := range []struct {
				name, body string
				status     int
			}{
				{"Slack error", `{"ok":false,"error":"invalid_arguments","response_metadata":{"messages":["Request rejected"]}}`, http.StatusOK},
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
						assert.Equal(t, "invalid_arguments", slackErr.Err)
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
