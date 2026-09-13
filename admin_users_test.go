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

func TestAdminUsersListFilters(t *testing.T) {
	falseValue := false
	trueValue := true
	api := newAdminUsersListTestClient(t, func(rw http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/admin.users.list", r.URL.Path)
		require.NoError(t, r.ParseForm())
		assert.Equal(t, "testing-token", r.FormValue("token"))
		assert.Equal(t, "T123", r.FormValue("team_id"))
		assert.Equal(t, "cursor-1", r.FormValue("cursor"))
		assert.Equal(t, "false", r.FormValue("is_active"))
		assert.Equal(t, "true", r.FormValue("include_deactivated_user_workspaces"))
		assert.Equal(t, "false", r.FormValue("only_guests"))
		assert.Equal(t, "true", r.FormValue("include_admins"))
		assert.Equal(t, "false", r.FormValue("include_owners"))
		assert.Equal(t, "250", r.FormValue("limit"))

		rw.Header().Set("Content-Type", "application/json")
		_, _ = rw.Write([]byte(`{"ok": true, "users": []}`))
	})

	_, err := api.AdminUsersList(context.Background(), AdminUsersListParams{
		TeamID:                           "T123",
		Cursor:                           "cursor-1",
		IsActive:                         &falseValue,
		IncludeDeactivatedUserWorkspaces: &trueValue,
		OnlyGuests:                       &falseValue,
		IncludeAdmins:                    &trueValue,
		IncludeOwners:                    &falseValue,
		Limit:                            250,
	})
	require.NoError(t, err)
}

func TestAdminUsersListOmitsOptionalFilters(t *testing.T) {
	api := newAdminUsersListTestClient(t, func(rw http.ResponseWriter, r *http.Request) {
		require.NoError(t, r.ParseForm())
		assert.Equal(t, []string{"testing-token"}, r.Form["token"])
		assert.Len(t, r.Form, 1)

		rw.Header().Set("Content-Type", "application/json")
		_, _ = rw.Write([]byte(`{"ok": true, "users": []}`))
	})

	_, err := api.AdminUsersList(context.Background(), AdminUsersListParams{})
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
				"workspaces": ["T123", "T456"],
				"has_2fa": true,
				"has_sso": true
			}]
		}`))
	})

	response, err := api.AdminUsersList(context.Background(), AdminUsersListParams{})
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
	assert.Zero(t, user.DeactivatedTS)
	assert.Zero(t, user.ExpirationTS)
	assert.Equal(t, []string{"T123", "T456"}, user.Workspaces)
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

	response, err := api.AdminUsersList(context.Background(), AdminUsersListParams{})
	require.NoError(t, err)
	assert.Equal(t, "cursor-2", response.ResponseMetadata.Cursor)
}

func TestAdminUsersListError(t *testing.T) {
	api := newAdminUsersListTestClient(t, func(rw http.ResponseWriter, r *http.Request) {
		rw.Header().Set("Content-Type", "application/json")
		_, _ = rw.Write([]byte(`{"ok": false, "error": "invalid_arguments"}`))
	})

	response, err := api.AdminUsersList(context.Background(), AdminUsersListParams{})
	assert.EqualError(t, err, "invalid_arguments")
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

	response, err := api.AdminUsersList(context.Background(), AdminUsersListParams{})
	require.NoError(t, err)
	require.Len(t, response.Users, 1)
	assert.Equal(t, "W123", response.Users[0].ID)
}
