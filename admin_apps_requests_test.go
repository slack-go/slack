package slack

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAdminAppsRequestsList(t *testing.T) {
	api := newAdminAppsTestClient(t, func(rw http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/admin.apps.requests.list", r.URL.Path)
		assert.NoError(t, r.ParseForm())
		assert.Equal(t, "testing-token", r.FormValue("token"))
		assert.Equal(t, "250", r.FormValue("limit"))
		assert.Equal(t, "cursor-1", r.FormValue("cursor"))
		assert.Equal(t, "E123", r.FormValue("enterprise_id"))
		assert.Equal(t, "true", r.FormValue("certified"))
		assert.NotContains(t, r.Form, "team_id")
		assert.Len(t, r.Form, 5)

		rw.Header().Set("Content-Type", "application/json")
		_, _ = rw.Write([]byte(`{
			"ok": true,
			"app_requests": [{
				"id": "Ar123",
				"app": {"id": "A789", "name": "Requested App"},
				"previous_resolution": {
					"status": "restricted",
					"scopes": [{"name": "files:read", "is_sensitive": false, "token_type": "user"}]
				},
				"user": {"id": "W123", "name": "Jane Doe", "email": "jane@example.com"},
				"team": {"id": "T123", "name": "Example", "domain": "example"},
				"scopes": [{"name": "chat:write", "is_sensitive": false, "token_type": "bot", "is_optional": true, "is_approved": false}],
				"message": "Please approve",
				"is_user_app_collaborator": true,
				"date_created": 1700000002
			}],
			"response_metadata": {"next_cursor": "cursor-2"}
		}`))
	})

	response, err := api.AdminAppsRequestsList(
		context.Background(),
		AdminAppsListOptionLimit(250),
		AdminAppsListOptionCursor("cursor-1"),
		AdminAppsListOptionEnterpriseID("E123"),
		AdminAppsListOptionCertified(true),
	)
	require.NoError(t, err)
	require.Len(t, response.AppRequests, 1)

	request := response.AppRequests[0]
	assert.Equal(t, "Ar123", request.ID)
	assert.Equal(t, "A789", request.App.ID)
	require.NotNil(t, request.PreviousResolution)
	assert.Equal(t, "restricted", request.PreviousResolution.Status)
	require.Len(t, request.PreviousResolution.Scopes, 1)
	assert.Equal(t, "files:read", request.PreviousResolution.Scopes[0].Name)
	assert.Equal(t, "W123", request.User.ID)
	assert.Equal(t, "T123", request.Team.ID)
	require.Len(t, request.Scopes, 1)
	require.NotNil(t, request.Scopes[0].IsOptional)
	assert.True(t, *request.Scopes[0].IsOptional)
	require.NotNil(t, request.Scopes[0].IsApproved)
	assert.False(t, *request.Scopes[0].IsApproved)
	assert.Equal(t, "Please approve", request.Message)
	assert.True(t, request.IsUserAppCollaborator)
	assert.Equal(t, int64(1700000002), request.DateCreated)
	assert.Equal(t, "cursor-2", response.ResponseMetadata.Cursor)
}

func TestAdminAppsRequestsListError(t *testing.T) {
	api := newAdminAppsTestClient(t, func(rw http.ResponseWriter, r *http.Request) {
		_, _ = rw.Write([]byte(`{"ok": false, "error": "invalid_cursor"}`))
	})

	response, err := api.AdminAppsRequestsList(context.Background())
	assert.EqualError(t, err, "invalid_cursor")
	assert.NotNil(t, response)
}
