package slack

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newSetManagedAppPermissionsTestClient(t *testing.T, permissions, body string) *Client {
	t.Helper()
	// The handler runs on the server goroutine, so it uses assert, not require.
	ts := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/apps.managed.permissions.set", r.URL.Path)
		assert.Equal(t, "application/x-www-form-urlencoded", r.Header.Get("Content-Type"))
		assert.NoError(t, r.ParseForm())
		assert.Equal(t, url.Values{
			"token":       {"client-token"},
			"app_id":      {"A123"},
			"permissions": {permissions},
		}, r.PostForm)

		rw.Header().Set("Content-Type", "application/json")
		_, _ = rw.Write([]byte(body))
	}))
	t.Cleanup(ts.Close)
	return New("client-token", OptionAPIURL(ts.URL+"/"))
}

func TestSetManagedAppPermissions(t *testing.T) {
	tests := []struct {
		name        string
		permissions ManagedAppPermission
		wire        string
		call        func(*Client, ManagedAppPermission) (*SetManagedAppPermissionsResponse, error)
	}{
		{
			name:        "SetManagedAppPermissions",
			permissions: ManagedAppPermissionAppOwner,
			wire:        "app_owner",
			call: func(api *Client, permissions ManagedAppPermission) (*SetManagedAppPermissionsResponse, error) {
				return api.SetManagedAppPermissions("A123", permissions)
			},
		},
		{
			name:        "SetManagedAppPermissionsContext",
			permissions: ManagedAppPermissionEveryone,
			wire:        "everyone",
			call: func(api *Client, permissions ManagedAppPermission) (*SetManagedAppPermissionsResponse, error) {
				return api.SetManagedAppPermissionsContext(context.Background(), "A123", permissions)
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			api := newSetManagedAppPermissionsTestClient(t, test.wire, `{
				"ok": true,
				"permission_type": "app_owner",
				"channel_restriction_mode": "example_restriction_mode",
				"user_ids": ["U012AB3CD"]
			}`)

			response, err := test.call(api, test.permissions)
			require.NoError(t, err)
			assert.Equal(t, "app_owner", response.PermissionType)
			assert.Equal(t, "example_restriction_mode", response.ChannelRestrictionMode)
			assert.Equal(t, []string{"U012AB3CD"}, response.UserIDs)
		})
	}
}

func TestSetManagedAppPermissionsError(t *testing.T) {
	api := newSetManagedAppPermissionsTestClient(t, "app_owner",
		`{"ok": false, "error": "app_not_managed"}`)

	response, err := api.SetManagedAppPermissions("A123", ManagedAppPermissionAppOwner)
	require.EqualError(t, err, "app_not_managed")
	var slackErr SlackErrorResponse
	require.ErrorAs(t, err, &slackErr)
	assert.Equal(t, "app_not_managed", slackErr.Err)
	require.NotNil(t, response)
	assert.False(t, response.Ok)
}
