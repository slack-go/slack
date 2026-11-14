package slack

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAdminUsersSessionMethodPaths(t *testing.T) {
	ctx := context.Background()
	userIDs := []string{"W123"}
	testAdminUsersMethodPaths(t, []adminUsersMethodPathCase{
		{"ClearSettings", "admin.users.session.clearSettings", func(api *Client) error {
			return api.AdminUsersSessionClearSettings(userIDs)
		}},
		{"ClearSettingsContext", "admin.users.session.clearSettings", func(api *Client) error {
			return api.AdminUsersSessionClearSettingsContext(ctx, userIDs)
		}},
		{"GetSettings", "admin.users.session.getSettings", func(api *Client) error {
			_, err := api.AdminUsersSessionGetSettings(userIDs)
			return err
		}},
		{"GetSettingsContext", "admin.users.session.getSettings", func(api *Client) error {
			_, err := api.AdminUsersSessionGetSettingsContext(ctx, userIDs)
			return err
		}},
		{"List", "admin.users.session.list", func(api *Client) error {
			_, err := api.AdminUsersSessionList()
			return err
		}},
		{"ListContext", "admin.users.session.list", func(api *Client) error {
			_, err := api.AdminUsersSessionListContext(ctx)
			return err
		}},
		{"Reset", "admin.users.session.reset", func(api *Client) error {
			return api.AdminUsersSessionReset("W123")
		}},
		{"ResetContext", "admin.users.session.reset", func(api *Client) error {
			return api.AdminUsersSessionResetContext(ctx, "W123")
		}},
		{"ResetBulk", "admin.users.session.resetBulk", func(api *Client) error {
			return api.AdminUsersSessionResetBulk(userIDs)
		}},
		{"ResetBulkContext", "admin.users.session.resetBulk", func(api *Client) error {
			return api.AdminUsersSessionResetBulkContext(ctx, userIDs)
		}},
		{"SetSettings", "admin.users.session.setSettings", func(api *Client) error {
			return api.AdminUsersSessionSetSettings(userIDs)
		}},
		{"SetSettingsContext", "admin.users.session.setSettings", func(api *Client) error {
			return api.AdminUsersSessionSetSettingsContext(ctx, userIDs)
		}},
	})
}

func TestAdminUsersSessionClearSettings(t *testing.T) {
	api := newAdminUsersListTestClient(t, func(rw http.ResponseWriter, r *http.Request) {
		assert.NoError(t, r.ParseForm())
		assert.Equal(t, "testing-token", r.FormValue("token"))
		assert.Equal(t, "W123,W456", r.FormValue("user_ids"))
		assert.Len(t, r.Form, 2)

		rw.Header().Set("Content-Type", "application/json")
		_, _ = rw.Write([]byte(`{"ok": true}`))
	})

	require.NoError(t, api.AdminUsersSessionClearSettings([]string{"W123", "W456"}))
}

func TestAdminUsersSessionGetSettings(t *testing.T) {
	api := newAdminUsersListTestClient(t, func(rw http.ResponseWriter, r *http.Request) {
		assert.NoError(t, r.ParseForm())
		assert.Equal(t, "W123,W456", r.FormValue("user_ids"))
		assert.Len(t, r.Form, 2)

		rw.Header().Set("Content-Type", "application/json")
		_, _ = rw.Write([]byte(`{
			"ok": true,
			"session_settings": [{
				"user_id": "W123",
				"desktop_app_browser_quit": true,
				"duration": 315569520
			}],
			"no_settings_applied": ["W456"]
		}`))
	})

	response, err := api.AdminUsersSessionGetSettings([]string{"W123", "W456"})
	require.NoError(t, err)
	assert.Equal(t, []AdminUserSessionSettings{{
		UserID:                "W123",
		DesktopAppBrowserQuit: true,
		Duration:              315569520,
	}}, response.SessionSettings)
	assert.Equal(t, []string{"W456"}, response.NoSettingsApplied)
}

func TestAdminUsersSessionList(t *testing.T) {
	api := newAdminUsersListTestClient(t, func(rw http.ResponseWriter, r *http.Request) {
		assert.NoError(t, r.ParseForm())
		assert.Equal(t, "testing-token", r.FormValue("token"))
		assert.Equal(t, "T123", r.FormValue("team_id"))
		assert.Equal(t, "W123", r.FormValue("user_id"))
		assert.Equal(t, "cursor-1", r.FormValue("cursor"))
		assert.Equal(t, "100", r.FormValue("limit"))
		assert.Len(t, r.Form, 5)

		rw.Header().Set("Content-Type", "application/json")
		_, _ = rw.Write([]byte(`{
			"ok": true,
			"active_sessions": [{
				"user_id": "W123",
				"team_id": "E123",
				"session_id": 1112275520242,
				"recent": {
					"device_hardware": "Intel",
					"os": "OS X",
					"os_version": "10.15.7",
					"slack_client_version": "91.0.4472.77",
					"ip": "192.0.2.1"
				},
				"created": {
					"device_hardware": "Intel",
					"os": "OS X",
					"os_version": "10.15.6",
					"slack_client_version": "90.0.4430.93",
					"ip": "192.0.2.2"
				}
			}, {
				"user_id": "W123",
				"team_id": "E123",
				"session_id": 1112275520243,
				"created": {"os": "iOS"}
			}],
			"response_metadata": {"next_cursor": "cursor-2"}
		}`))
	})

	response, err := api.AdminUsersSessionListContext(context.Background(),
		nil,
		AdminUsersSessionListOptionUser("T123", "W123"),
		AdminUsersSessionListOptionCursor("cursor-1"),
		AdminUsersSessionListOptionLimit(100),
	)
	require.NoError(t, err)
	require.Len(t, response.ActiveSessions, 2)
	assert.Equal(t, AdminUserSession{
		UserID:    "W123",
		TeamID:    "E123",
		SessionID: 1112275520242,
		Created: AdminUserSessionDevice{
			DeviceHardware:     "Intel",
			OS:                 "OS X",
			OSVersion:          "10.15.6",
			SlackClientVersion: "90.0.4430.93",
			IP:                 "192.0.2.2",
		},
		Recent: &AdminUserSessionDevice{
			DeviceHardware:     "Intel",
			OS:                 "OS X",
			OSVersion:          "10.15.7",
			SlackClientVersion: "91.0.4472.77",
			IP:                 "192.0.2.1",
		},
	}, response.ActiveSessions[0])
	// Slack omits recent when the session hasn't changed since it was created.
	assert.Nil(t, response.ActiveSessions[1].Recent)
	assert.Equal(t, "iOS", response.ActiveSessions[1].Created.OS)
	assert.Equal(t, "cursor-2", response.ResponseMetadata.Cursor)
}

func TestAdminUsersSessionListOmitsOptionalArguments(t *testing.T) {
	api := newAdminUsersListTestClient(t, func(rw http.ResponseWriter, r *http.Request) {
		assert.NoError(t, r.ParseForm())
		assert.Equal(t, []string{"testing-token"}, r.Form["token"])
		assert.Len(t, r.Form, 1)

		rw.Header().Set("Content-Type", "application/json")
		_, _ = rw.Write([]byte(`{"ok": true, "active_sessions": []}`))
	})

	_, err := api.AdminUsersSessionList()
	require.NoError(t, err)
}

func TestAdminUsersSessionListError(t *testing.T) {
	api := newAdminUsersListTestClient(t, func(rw http.ResponseWriter, r *http.Request) {
		rw.Header().Set("Content-Type", "application/json")
		_, _ = rw.Write([]byte(`{"ok": false, "error": "no_active_sessions"}`))
	})

	response, err := api.AdminUsersSessionList()
	require.EqualError(t, err, "no_active_sessions")
	assert.NotNil(t, response)
}

func TestAdminUsersSessionReset(t *testing.T) {
	api := newAdminUsersListTestClient(t, func(rw http.ResponseWriter, r *http.Request) {
		assert.NoError(t, r.ParseForm())
		assert.Equal(t, "testing-token", r.FormValue("token"))
		assert.Equal(t, "W123", r.FormValue("user_id"))
		assert.Len(t, r.Form, 2)

		rw.Header().Set("Content-Type", "application/json")
		_, _ = rw.Write([]byte(`{"ok": true}`))
	})

	require.NoError(t, api.AdminUsersSessionReset("W123"))
}

func TestAdminUsersSessionResetBulk(t *testing.T) {
	api := newAdminUsersListTestClient(t, func(rw http.ResponseWriter, r *http.Request) {
		assert.NoError(t, r.ParseForm())
		assert.Equal(t, "testing-token", r.FormValue("token"))
		assert.Equal(t, "W123,W456", r.FormValue("user_ids"))
		assert.Len(t, r.Form, 2)

		rw.Header().Set("Content-Type", "application/json")
		_, _ = rw.Write([]byte(`{"ok": false, "error": "failed_for_some_users"}`))
	})

	err := api.AdminUsersSessionResetBulk([]string{"W123", "W456"})
	require.EqualError(t, err, "failed_for_some_users")
}

func TestAdminUsersSessionResetOptions(t *testing.T) {
	userIDs := []string{"W123"}
	tests := []struct {
		name  string
		field string
		call  func(*Client) error
	}{
		{"ResetMobileOnly", "mobile_only", func(api *Client) error {
			return api.AdminUsersSessionReset("W123", nil,
				AdminUsersSessionResetOptionMobileOnly(true),
				AdminUsersSessionResetOptionWebOnly(false))
		}},
		{"ResetWebOnly", "web_only", func(api *Client) error {
			return api.AdminUsersSessionReset("W123", nil,
				AdminUsersSessionResetOptionMobileOnly(false),
				AdminUsersSessionResetOptionWebOnly(true))
		}},
		{"ResetBulkMobileOnly", "mobile_only", func(api *Client) error {
			return api.AdminUsersSessionResetBulk(userIDs, nil,
				AdminUsersSessionResetBulkOptionMobileOnly(true),
				AdminUsersSessionResetBulkOptionWebOnly(false))
		}},
		{"ResetBulkWebOnly", "web_only", func(api *Client) error {
			return api.AdminUsersSessionResetBulk(userIDs, nil,
				AdminUsersSessionResetBulkOptionMobileOnly(false),
				AdminUsersSessionResetBulkOptionWebOnly(true))
		}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			api := newAdminUsersListTestClient(t, func(rw http.ResponseWriter, r *http.Request) {
				assert.NoError(t, r.ParseForm())
				assert.Equal(t, "true", r.FormValue(tt.field))
				// Both filters default to false, so the other one is left out.
				assert.Len(t, r.Form, 3)

				rw.Header().Set("Content-Type", "application/json")
				_, _ = rw.Write([]byte(`{"ok": true}`))
			})

			require.NoError(t, tt.call(api))
		})
	}
}

func TestAdminUsersSessionSetSettings(t *testing.T) {
	api := newAdminUsersListTestClient(t, func(rw http.ResponseWriter, r *http.Request) {
		assert.NoError(t, r.ParseForm())
		assert.Equal(t, "W123,W456", r.FormValue("user_ids"))
		assert.Equal(t, "28800", r.FormValue("duration"))
		// An explicit false changes the setting, so it must be sent.
		assert.Equal(t, "false", r.FormValue("desktop_app_browser_quit"))
		assert.Len(t, r.Form, 4)

		rw.Header().Set("Content-Type", "application/json")
		_, _ = rw.Write([]byte(`{"ok": true}`))
	})

	err := api.AdminUsersSessionSetSettings([]string{"W123", "W456"},
		nil,
		AdminUsersSessionSetSettingsOptionDuration(28800),
		AdminUsersSessionSetSettingsOptionDesktopAppBrowserQuit(false),
	)
	require.NoError(t, err)
}

func TestAdminUsersSessionSetSettingsSendsExplicitDuration(t *testing.T) {
	api := newAdminUsersListTestClient(t, func(rw http.ResponseWriter, r *http.Request) {
		assert.NoError(t, r.ParseForm())
		// Slack rejects a zero duration, so it is sent rather than silently dropped.
		assert.Equal(t, "0", r.FormValue("duration"))
		assert.Len(t, r.Form, 3)

		rw.Header().Set("Content-Type", "application/json")
		_, _ = rw.Write([]byte(`{"ok": false, "error": "invalid_arguments"}`))
	})

	err := api.AdminUsersSessionSetSettings([]string{"W123"},
		AdminUsersSessionSetSettingsOptionDuration(0),
	)
	require.EqualError(t, err, "invalid_arguments")
}

func TestAdminUsersSessionSetSettingsOmitsOptionalArguments(t *testing.T) {
	api := newAdminUsersListTestClient(t, func(rw http.ResponseWriter, r *http.Request) {
		assert.NoError(t, r.ParseForm())
		assert.Len(t, r.Form, 2)

		rw.Header().Set("Content-Type", "application/json")
		_, _ = rw.Write([]byte(`{"ok": true}`))
	})

	require.NoError(t, api.AdminUsersSessionSetSettings([]string{"W123"}))
}
