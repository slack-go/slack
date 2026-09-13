package slack

import (
	"context"
	"net/url"
	"strconv"
)

// AdminUsersListParams contains arguments for AdminUsersList.
type AdminUsersListParams struct {
	TeamID                           string
	Cursor                           string
	IsActive                         *bool
	IncludeDeactivatedUserWorkspaces *bool
	OnlyGuests                       *bool
	IncludeAdmins                    *bool
	IncludeOwners                    *bool
	Limit                            int
}

// AdminUser represents a user returned by admin.users.list.
type AdminUser struct {
	ID                string   `json:"id"`
	Email             string   `json:"email"`
	IsAdmin           bool     `json:"is_admin"`
	IsOwner           bool     `json:"is_owner"`
	IsPrimaryOwner    bool     `json:"is_primary_owner"`
	IsRestricted      bool     `json:"is_restricted"`
	IsUltraRestricted bool     `json:"is_ultra_restricted"`
	IsBot             bool     `json:"is_bot"`
	Username          string   `json:"username"`
	FullName          string   `json:"full_name"`
	IsActive          bool     `json:"is_active"`
	DateCreated       int64    `json:"date_created"`
	DeactivatedTS     int64    `json:"deactivated_ts"`
	ExpirationTS      int64    `json:"expiration_ts"`
	Workspaces        []string `json:"workspaces"`
	Has2FA            bool     `json:"has_2fa"`
	HasSSO            bool     `json:"has_sso"`
}

// AdminUsersListResponse represents the response from admin.users.list.
type AdminUsersListResponse struct {
	SlackResponse
	Users []AdminUser `json:"users"`
}

// AdminUsersList lists users on a workspace or Enterprise organization.
//
// Slack API docs: https://docs.slack.dev/reference/methods/admin.users.list
func (api *Client) AdminUsersList(ctx context.Context, params AdminUsersListParams) (*AdminUsersListResponse, error) {
	values := url.Values{
		"token": {api.token},
	}

	if params.TeamID != "" {
		values.Add("team_id", params.TeamID)
	}
	if params.Cursor != "" {
		values.Add("cursor", params.Cursor)
	}
	if params.IsActive != nil {
		values.Add("is_active", strconv.FormatBool(*params.IsActive))
	}
	if params.IncludeDeactivatedUserWorkspaces != nil {
		values.Add("include_deactivated_user_workspaces", strconv.FormatBool(*params.IncludeDeactivatedUserWorkspaces))
	}
	if params.OnlyGuests != nil {
		values.Add("only_guests", strconv.FormatBool(*params.OnlyGuests))
	}
	if params.IncludeAdmins != nil {
		values.Add("include_admins", strconv.FormatBool(*params.IncludeAdmins))
	}
	if params.IncludeOwners != nil {
		values.Add("include_owners", strconv.FormatBool(*params.IncludeOwners))
	}
	if params.Limit > 0 {
		values.Add("limit", strconv.Itoa(params.Limit))
	}

	response := &AdminUsersListResponse{}
	err := api.postMethod(ctx, "admin.users.list", values, response)
	if err != nil {
		return nil, err
	}

	return response, response.Err()
}
