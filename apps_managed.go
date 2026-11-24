package slack

import (
	"context"
	"net/url"
)

// ManagedAppPermission sets who can interact with a managed app.
type ManagedAppPermission string

const (
	ManagedAppPermissionEveryone ManagedAppPermission = "everyone"
	ManagedAppPermissionAppOwner ManagedAppPermission = "app_owner"
)

// SetManagedAppPermissionsResponse is the response from apps.managed.permissions.set.
type SetManagedAppPermissionsResponse struct {
	SlackResponse
	// PermissionType is the app's permission, such as app_owner. It comes from
	// Slack's recorded API samples, not the reference.
	PermissionType string `json:"permission_type"`
	// ChannelRestrictionMode comes from Slack's recorded API samples, not the
	// reference.
	ChannelRestrictionMode string `json:"channel_restriction_mode"`
	// UserIDs comes from Slack's recorded API samples, not the reference.
	UserIDs []string `json:"user_ids"`
}

// SetManagedAppPermissions sets who can interact with a managed app. The client
// token must be issued to the manager app. Slack only applies permissions before
// the app is installed; afterwards it changes nothing and returns the current
// permissions. Enterprise Grid only.
//
// Slack API docs: https://docs.slack.dev/reference/methods/apps.managed.permissions.set
func (api *Client) SetManagedAppPermissions(appID string, permissions ManagedAppPermission) (*SetManagedAppPermissionsResponse, error) {
	return api.SetManagedAppPermissionsContext(context.Background(), appID, permissions)
}

// SetManagedAppPermissionsContext sets who can interact with a managed app with a
// custom context.
func (api *Client) SetManagedAppPermissionsContext(ctx context.Context, appID string, permissions ManagedAppPermission) (*SetManagedAppPermissionsResponse, error) {
	values := url.Values{
		"token":       {api.token},
		"app_id":      {appID},
		"permissions": {string(permissions)},
	}

	response := &SetManagedAppPermissionsResponse{}
	if err := api.postMethod(ctx, "apps.managed.permissions.set", values, response); err != nil {
		return nil, err
	}
	return response, response.Err()
}
