package slack

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newAdminUsersListTestClient(t *testing.T, handler http.HandlerFunc) *Client {
	t.Helper()
	ts := httptest.NewServer(handler)
	t.Cleanup(ts.Close)
	return New("testing-token", OptionAPIURL(ts.URL+"/"))
}

// The handlers run on the server goroutine, so they use assert, not require.

func TestAdminUsersListFilters(t *testing.T) {
	api := newAdminUsersListTestClient(t, func(rw http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/admin.users.list", r.URL.Path)
		assert.NoError(t, r.ParseForm())
		assert.Equal(t, "testing-token", r.FormValue("token"))
		assert.Equal(t, "T123", r.FormValue("team_id"))
		assert.Equal(t, "cursor-1", r.FormValue("cursor"))
		assert.Equal(t, "250", r.FormValue("limit"))
		// is_active defaults to true, so false must be sent.
		assert.Equal(t, "false", r.FormValue("is_active"))
		assert.Equal(t, "true", r.FormValue("include_deactivated_user_workspaces"))
		assert.Equal(t, "true", r.FormValue("include_admins"))
		// The other filters default to false, so false is left out.
		assert.NotContains(t, r.Form, "only_guests")
		assert.NotContains(t, r.Form, "include_owners")

		rw.Header().Set("Content-Type", "application/json")
		_, _ = rw.Write([]byte(`{"ok": true, "users": []}`))
	})

	_, err := api.AdminUsersList(context.Background(),
		AdminUsersListOptionTeamID("T123"),
		AdminUsersListOptionCursor("cursor-1"),
		AdminUsersListOptionLimit(250),
		AdminUsersListOptionIsActive(false),
		AdminUsersListOptionIncludeDeactivatedUserWorkspaces(true),
		AdminUsersListOptionOnlyGuests(false),
		AdminUsersListOptionIncludeAdmins(true),
		AdminUsersListOptionIncludeOwners(false),
	)
	require.NoError(t, err)
}

func TestAdminUsersListOmitsOptionalFilters(t *testing.T) {
	api := newAdminUsersListTestClient(t, func(rw http.ResponseWriter, r *http.Request) {
		assert.NoError(t, r.ParseForm())
		assert.Equal(t, []string{"testing-token"}, r.Form["token"])
		assert.Len(t, r.Form, 1)

		rw.Header().Set("Content-Type", "application/json")
		_, _ = rw.Write([]byte(`{"ok": true, "users": []}`))
	})

	_, err := api.AdminUsersList(context.Background())
	require.NoError(t, err)
}

func TestAdminUsersListMultipleWorkspaces(t *testing.T) {
	api := newAdminUsersListTestClient(t, func(rw http.ResponseWriter, r *http.Request) {
		rw.Header().Set("Content-Type", "application/json")
		_, _ = rw.Write([]byte(`{
			"ok": true,
			"users": [{
				"id": "W123",
				"email": "alex@example.com",
				"is_admin": true,
				"is_owner": false,
				"is_primary_owner": false,
				"is_restricted": false,
				"is_ultra_restricted": false,
				"is_bot": false,
				"username": "alex",
				"full_name": "Alex Example",
				"is_active": true,
				"date_created": 1566922090,
				"deactivated_ts": 0,
				"expiration_ts": 0,
				"last_active_ts": 1700000000,
				"workspaces": ["T123", "T456"],
				"roles": ["Rl0A"],
				"has_2fa": true,
				"has_sso": true
			}]
		}`))
	})

	response, err := api.AdminUsersList(context.Background())
	require.NoError(t, err)
	require.Len(t, response.Users, 1)

	user := response.Users[0]
	assert.Equal(t, "W123", user.ID)
	assert.Equal(t, "alex@example.com", user.Email)
	assert.True(t, user.IsAdmin)
	assert.False(t, user.IsOwner)
	assert.False(t, user.IsPrimaryOwner)
	assert.False(t, user.IsRestricted)
	assert.False(t, user.IsUltraRestricted)
	assert.False(t, user.IsBot)
	assert.Equal(t, "alex", user.Username)
	assert.Equal(t, "Alex Example", user.FullName)
	assert.True(t, user.IsActive)
	assert.Equal(t, int64(1566922090), user.DateCreated)
	assert.Zero(t, user.DeactivatedTimestamp)
	assert.Zero(t, user.ExpirationTimestamp)
	assert.Equal(t, int64(1700000000), user.LastActiveTimestamp)
	assert.Equal(t, []string{"T123", "T456"}, user.Workspaces)
	assert.Equal(t, []string{"Rl0A"}, user.Roles)
	assert.True(t, user.Has2FA)
	assert.True(t, user.HasSSO)
}

func TestAdminUsersListCursorPagination(t *testing.T) {
	api := newAdminUsersListTestClient(t, func(rw http.ResponseWriter, r *http.Request) {
		rw.Header().Set("Content-Type", "application/json")
		_, _ = rw.Write([]byte(`{
			"ok": true,
			"users": [],
			"response_metadata": {"next_cursor": "cursor-2"}
		}`))
	})

	response, err := api.AdminUsersList(context.Background())
	require.NoError(t, err)
	assert.Equal(t, "cursor-2", response.ResponseMetadata.Cursor)
}

func TestAdminUsersListError(t *testing.T) {
	api := newAdminUsersListTestClient(t, func(rw http.ResponseWriter, r *http.Request) {
		rw.Header().Set("Content-Type", "application/json")
		_, _ = rw.Write([]byte(`{"ok": false, "error": "invalid_arguments"}`))
	})

	response, err := api.AdminUsersList(context.Background())
	require.EqualError(t, err, "invalid_arguments")
	assert.NotNil(t, response)
}

func TestAdminUsersListIgnoresUnknownFields(t *testing.T) {
	api := newAdminUsersListTestClient(t, func(rw http.ResponseWriter, r *http.Request) {
		rw.Header().Set("Content-Type", "application/json")
		_, _ = rw.Write([]byte(`{
			"ok": true,
			"future_response_field": "value",
			"users": [{"id": "W123", "future_user_field": {"nested": true}}]
		}`))
	})

	response, err := api.AdminUsersList(context.Background())
	require.NoError(t, err)
	require.Len(t, response.Users, 1)
	assert.Equal(t, "W123", response.Users[0].ID)
}

func TestAdminUsersSetAccountType(t *testing.T) {
	tests := []struct {
		name   string
		method string
		call   func(*Client) error
	}{
		{"SetAdmin", "admin.users.setAdmin", func(api *Client) error {
			return api.AdminUsersSetAdmin("E123", "W123")
		}},
		{"SetAdminContext", "admin.users.setAdmin", func(api *Client) error {
			return api.AdminUsersSetAdminContext(context.Background(), "E123", "W123")
		}},
		{"SetOwner", "admin.users.setOwner", func(api *Client) error {
			return api.AdminUsersSetOwner("E123", "W123")
		}},
		{"SetOwnerContext", "admin.users.setOwner", func(api *Client) error {
			return api.AdminUsersSetOwnerContext(context.Background(), "E123", "W123")
		}},
		{"SetRegular", "admin.users.setRegular", func(api *Client) error {
			return api.AdminUsersSetRegular("E123", "W123")
		}},
		{"SetRegularContext", "admin.users.setRegular", func(api *Client) error {
			return api.AdminUsersSetRegularContext(context.Background(), "E123", "W123")
		}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			api := newAdminUsersListTestClient(t, func(rw http.ResponseWriter, r *http.Request) {
				assert.Equal(t, http.MethodPost, r.Method)
				assert.Equal(t, "/"+tt.method, r.URL.Path)
				assert.NoError(t, r.ParseForm())
				assert.Equal(t, "testing-token", r.FormValue("token"))
				assert.Equal(t, "E123", r.FormValue("team_id"))
				assert.Equal(t, "W123", r.FormValue("user_id"))
				assert.Len(t, r.Form, 3)

				rw.Header().Set("Content-Type", "application/json")
				_, _ = rw.Write([]byte(`{"ok": true}`))
			})

			require.NoError(t, tt.call(api))
		})
	}
}

func TestAdminUsersSetRegularWorkspaceFailure(t *testing.T) {
	api := newAdminUsersListTestClient(t, func(rw http.ResponseWriter, r *http.Request) {
		rw.Header().Set("Content-Type", "application/json")
		_, _ = rw.Write([]byte(`{"ok": false, "error": "failed_to_set_regular_on_workspaces"}`))
	})

	err := api.AdminUsersSetRegular("E123", "W123")
	var slackErr SlackErrorResponse
	require.ErrorAs(t, err, &slackErr)
	assert.Equal(t, "failed_to_set_regular_on_workspaces", slackErr.Err)
}

func TestAdminUsersGetExpiration(t *testing.T) {
	api := newAdminUsersListTestClient(t, func(rw http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/admin.users.getExpiration", r.URL.Path)
		assert.NoError(t, r.ParseForm())
		assert.Equal(t, "testing-token", r.FormValue("token"))
		assert.Equal(t, "W123", r.FormValue("user_id"))
		assert.NotContains(t, r.Form, "target_team")
		assert.Len(t, r.Form, 2)

		rw.Header().Set("Content-Type", "application/json")
		_, _ = rw.Write([]byte(`{
			"ok": true,
			"user": {
				"id": "W123",
				"email": "guest@example.com",
				"is_restricted": false,
				"is_ultra_restricted": true,
				"expiration_ts": 1767225600
			}
		}`))
	})

	response, err := api.AdminUsersGetExpiration("W123")
	require.NoError(t, err)
	assert.Equal(t, AdminUserExpiration{
		ID:                  "W123",
		Email:               "guest@example.com",
		IsRestricted:        false,
		IsUltraRestricted:   true,
		ExpirationTimestamp: 1767225600,
	}, response.User)
}

func TestAdminUsersGetExpirationTargetTeam(t *testing.T) {
	api := newAdminUsersListTestClient(t, func(rw http.ResponseWriter, r *http.Request) {
		assert.NoError(t, r.ParseForm())
		assert.Equal(t, "T123", r.FormValue("target_team"))
		assert.NotContains(t, r.Form, "team_id")
		assert.Len(t, r.Form, 3)

		rw.Header().Set("Content-Type", "application/json")
		_, _ = rw.Write([]byte(`{"ok": true, "user": {"id": "W123"}}`))
	})

	_, err := api.AdminUsersGetExpirationContext(context.Background(), "W123",
		nil,
		AdminUsersGetExpirationOptionTargetTeam("T123"),
	)
	require.NoError(t, err)
}

func TestAdminUsersGetExpirationError(t *testing.T) {
	api := newAdminUsersListTestClient(t, func(rw http.ResponseWriter, r *http.Request) {
		rw.Header().Set("Content-Type", "application/json")
		_, _ = rw.Write([]byte(`{"ok": false, "error": "user_is_not_guest"}`))
	})

	response, err := api.AdminUsersGetExpiration("W123")
	require.EqualError(t, err, "user_is_not_guest")
	assert.NotNil(t, response)
}

func TestAdminUsersSetExpiration(t *testing.T) {
	api := newAdminUsersListTestClient(t, func(rw http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/admin.users.setExpiration", r.URL.Path)
		assert.NoError(t, r.ParseForm())
		assert.Equal(t, "testing-token", r.FormValue("token"))
		assert.Equal(t, "W123", r.FormValue("user_id"))
		assert.Equal(t, "1767225600", r.FormValue("expiration_ts"))
		assert.NotContains(t, r.Form, "team_id")
		assert.Len(t, r.Form, 3)

		rw.Header().Set("Content-Type", "application/json")
		_, _ = rw.Write([]byte(`{"ok": true}`))
	})

	require.NoError(t, api.AdminUsersSetExpiration("W123", 1767225600))
}

func TestAdminUsersSetExpirationTeamID(t *testing.T) {
	api := newAdminUsersListTestClient(t, func(rw http.ResponseWriter, r *http.Request) {
		assert.NoError(t, r.ParseForm())
		assert.Equal(t, "T123", r.FormValue("team_id"))
		assert.Len(t, r.Form, 4)

		rw.Header().Set("Content-Type", "application/json")
		_, _ = rw.Write([]byte(`{"ok": true}`))
	})

	err := api.AdminUsersSetExpirationContext(context.Background(), "W123", 1767225600,
		nil,
		AdminUsersSetExpirationOptionTeamID("T123"),
	)
	require.NoError(t, err)
}

func TestAdminUsersSetExpirationError(t *testing.T) {
	api := newAdminUsersListTestClient(t, func(rw http.ResponseWriter, r *http.Request) {
		rw.Header().Set("Content-Type", "application/json")
		_, _ = rw.Write([]byte(`{"ok": false, "error": "user_is_not_a_guest"}`))
	})

	err := api.AdminUsersSetExpiration("W123", 1767225600)
	require.EqualError(t, err, "user_is_not_a_guest")
}
