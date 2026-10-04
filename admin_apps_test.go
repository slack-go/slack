package slack

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newAdminAppsTestClient(t *testing.T, handler http.HandlerFunc) *Client {
	t.Helper()
	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)
	return New("testing-token", OptionAPIURL(server.URL+"/"))
}

// The handlers run on the server goroutine, so they use assert, not require.

func TestAdminAppsApprovedList(t *testing.T) {
	api := newAdminAppsTestClient(t, func(rw http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/admin.apps.approved.list", r.URL.Path)
		assert.NoError(t, r.ParseForm())
		assert.Equal(t, "testing-token", r.FormValue("token"))
		assert.Equal(t, "1000", r.FormValue("limit"))
		assert.Equal(t, "cursor-1", r.FormValue("cursor"))
		assert.Equal(t, "E123", r.FormValue("enterprise_id"))
		assert.Equal(t, "false", r.FormValue("certified"))
		assert.NotContains(t, r.Form, "team_id")

		rw.Header().Set("Content-Type", "application/json")
		_, _ = rw.Write([]byte(`{
			"ok": true,
			"approved_apps": [{
				"app": {
					"id": "A123",
					"name": "Test App",
					"description": "test",
					"help_url": "https://example.com/help",
					"privacy_policy_url": "https://example.com/privacy",
					"app_homepage_url": "https://example.com",
					"app_directory_url": "https://example.slack.com/apps/A123",
					"is_app_directory_approved": true,
					"is_internal": false,
					"developer_type": "third_party",
					"socket_mode_enabled": true,
					"icons": {"image_32": "https://example.com/32.png", "image_original": "https://example.com/original.png"},
					"additional_info": "additional"
				},
				"scopes": [{"name": "chat:write", "description": "Send messages", "is_sensitive": true, "token_type": "bot"}],
				"date_updated": 1700000000,
				"last_resolved_by": {"actor_id": "W123", "actor_type": "user"}
			}],
			"response_metadata": {"next_cursor": "cursor-2"}
		}`))
	})

	response, err := api.AdminAppsApprovedList(
		context.Background(),
		AdminAppsListOptionLimit(1000),
		AdminAppsListOptionCursor("cursor-1"),
		AdminAppsListOptionEnterpriseID("E123"),
		AdminAppsListOptionCertified(false),
	)
	require.NoError(t, err)
	require.Len(t, response.ApprovedApps, 1)

	item := response.ApprovedApps[0]
	assert.Equal(t, "A123", item.App.ID)
	assert.Equal(t, "Test App", item.App.Name)
	assert.Equal(t, "https://example.com/32.png", item.App.Icons.Image32)
	assert.Equal(t, "https://example.com/original.png", item.App.Icons.ImageOriginal)
	assert.Equal(t, "additional", item.App.AdditionalInfo)
	require.Len(t, item.Scopes, 1)
	assert.Equal(t, "chat:write", item.Scopes[0].Name)
	assert.True(t, item.Scopes[0].IsSensitive)
	assert.Equal(t, int64(1700000000), item.DateUpdated)
	assert.Equal(t, "W123", item.LastResolvedBy.ActorID)
	assert.Equal(t, "cursor-2", response.ResponseMetadata.Cursor)
}

func TestAdminAppsRestrictedList(t *testing.T) {
	api := newAdminAppsTestClient(t, func(rw http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/admin.apps.restricted.list", r.URL.Path)
		assert.NoError(t, r.ParseForm())
		assert.Equal(t, "T123", r.FormValue("team_id"))
		assert.Equal(t, "true", r.FormValue("certified"))

		rw.Header().Set("Content-Type", "application/json")
		_, _ = rw.Write([]byte(`{
			"ok": true,
			"restricted_apps": [{
				"app": {"id": "A456", "name": "Restricted App"},
				"date_updated": 1700000001
			}]
		}`))
	})

	response, err := api.AdminAppsRestrictedList(
		context.Background(),
		AdminAppsListOptionTeamID("T123"),
		AdminAppsListOptionCertified(true),
	)
	require.NoError(t, err)
	require.Len(t, response.RestrictedApps, 1)
	assert.Equal(t, "A456", response.RestrictedApps[0].App.ID)
	assert.Equal(t, "Restricted App", response.RestrictedApps[0].App.Name)
	assert.Equal(t, int64(1700000001), response.RestrictedApps[0].DateUpdated)
}

func TestAdminAppsListOmitsOptionalArguments(t *testing.T) {
	api := newAdminAppsTestClient(t, func(rw http.ResponseWriter, r *http.Request) {
		assert.NoError(t, r.ParseForm())
		assert.Equal(t, []string{"testing-token"}, r.Form["token"])
		assert.Len(t, r.Form, 1)
		_, _ = rw.Write([]byte(`{"ok": true, "approved_apps": []}`))
	})

	_, err := api.AdminAppsApprovedList(context.Background())
	require.NoError(t, err)
}

func TestAdminAppsApprovedListError(t *testing.T) {
	api := newAdminAppsTestClient(t, func(rw http.ResponseWriter, r *http.Request) {
		_, _ = rw.Write([]byte(`{"ok": false, "error": "invalid_cursor"}`))
	})

	response, err := api.AdminAppsApprovedList(context.Background())
	assert.EqualError(t, err, "invalid_cursor")
	assert.NotNil(t, response)
}
