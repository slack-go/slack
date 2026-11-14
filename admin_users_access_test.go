package slack

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAdminUsersAccessMethodPaths(t *testing.T) {
	ctx := context.Background()
	userIDs := []string{"W123"}
	testAdminUsersMethodPaths(t, []adminUsersMethodPathCase{
		{"Pause", "admin.users.access.pause", func(api *Client) error {
			_, err := api.AdminUsersAccessPause(userIDs)
			return err
		}},
		{"PauseContext", "admin.users.access.pause", func(api *Client) error {
			_, err := api.AdminUsersAccessPauseContext(ctx, userIDs)
			return err
		}},
		{"Resume", "admin.users.access.resume", func(api *Client) error {
			_, err := api.AdminUsersAccessResume(userIDs)
			return err
		}},
		{"ResumeContext", "admin.users.access.resume", func(api *Client) error {
			_, err := api.AdminUsersAccessResumeContext(ctx, userIDs)
			return err
		}},
	})
}

func TestAdminUsersAccessPause(t *testing.T) {
	api := newAdminUsersListTestClient(t, func(rw http.ResponseWriter, r *http.Request) {
		assert.NoError(t, r.ParseForm())
		assert.Equal(t, "testing-token", r.FormValue("token"))
		assert.Equal(t, "W123,W456", r.FormValue("user_ids"))
		assert.Equal(t, "Shift ended", r.FormValue("reason"))
		assert.Len(t, r.Form, 3)

		rw.Header().Set("Content-Type", "application/json")
		_, _ = rw.Write([]byte(`{"ok": true, "invalid_users": []}`))
	})

	response, err := api.AdminUsersAccessPause([]string{"W123", "W456"},
		nil, AdminUsersAccessPauseOptionReason("Shift ended"))
	require.NoError(t, err)
	assert.NotNil(t, response.InvalidUsers)
	assert.Empty(t, response.InvalidUsers)
}

func TestAdminUsersAccessResume(t *testing.T) {
	api := newAdminUsersListTestClient(t, func(rw http.ResponseWriter, r *http.Request) {
		assert.NoError(t, r.ParseForm())
		assert.Equal(t, "testing-token", r.FormValue("token"))
		assert.Equal(t, "W123,W456", r.FormValue("user_ids"))
		assert.Equal(t, "Shift started", r.FormValue("reason"))
		assert.Len(t, r.Form, 3)

		rw.Header().Set("Content-Type", "application/json")
		_, _ = rw.Write([]byte(`{"ok": true, "invalid_users": []}`))
	})

	response, err := api.AdminUsersAccessResume([]string{"W123", "W456"},
		nil, AdminUsersAccessResumeOptionReason("Shift started"))
	require.NoError(t, err)
	assert.NotNil(t, response.InvalidUsers)
	assert.Empty(t, response.InvalidUsers)
}

func TestAdminUsersAccessInvalidUsers(t *testing.T) {
	tests := []struct {
		name string
		call func(*Client, []string) (*AdminUsersAccessResponse, error)
	}{
		{"Pause", func(api *Client, userIDs []string) (*AdminUsersAccessResponse, error) {
			return api.AdminUsersAccessPause(userIDs)
		}},
		{"Resume", func(api *Client, userIDs []string) (*AdminUsersAccessResponse, error) {
			return api.AdminUsersAccessResume(userIDs)
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			api := newAdminUsersListTestClient(t, func(rw http.ResponseWriter, r *http.Request) {
				assert.NoError(t, r.ParseForm())
				assert.Equal(t, "U111AAA111,U222BBB222", r.FormValue("user_ids"))
				// An unset reason is left out.
				assert.Len(t, r.Form, 2)

				rw.Header().Set("Content-Type", "application/json")
				_, _ = rw.Write([]byte(`{
					"ok": false,
					"error": "failed_for_some_users",
					"invalid_users": [
						{"user_id": "U222BBB222", "reason": "user_not_frontline_member"}
					]
				}`))
			})

			response, err := tt.call(api, []string{"U111AAA111", "U222BBB222"})
			require.EqualError(t, err, "failed_for_some_users")
			require.NotNil(t, response)
			assert.Equal(t, []AdminUsersAccessInvalidUser{
				{UserID: "U222BBB222", Reason: "user_not_frontline_member"},
			}, response.InvalidUsers)
		})
	}
}
