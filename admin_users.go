package slack

import (
	"context"
	"net/url"
	"strconv"
)

type adminUsersListParams struct {
	teamID                           string
	cursor                           string
	limit                            int
	isActive                         *bool
	includeDeactivatedUserWorkspaces bool
	onlyGuests                       bool
	includeAdmins                    bool
	includeOwners                    bool
}

// AdminUsersListOption is an option for AdminUsersList.
type AdminUsersListOption func(*adminUsersListParams)

// AdminUsersListOptionTeamID filters results to the specified workspace.
func AdminUsersListOptionTeamID(teamID string) AdminUsersListOption {
	return func(params *adminUsersListParams) {
		params.teamID = teamID
	}
}

// AdminUsersListOptionCursor sets the cursor for pagination.
func AdminUsersListOptionCursor(cursor string) AdminUsersListOption {
	return func(params *adminUsersListParams) {
		params.cursor = cursor
	}
}

// AdminUsersListOptionLimit sets the maximum number of users to return.
func AdminUsersListOptionLimit(limit int) AdminUsersListOption {
	return func(params *adminUsersListParams) {
		params.limit = limit
	}
}

// AdminUsersListOptionIsActive returns only active users (true) or only deactivated
// users (false). Slack returns only active users when it is not set.
func AdminUsersListOptionIsActive(isActive bool) AdminUsersListOption {
	return func(params *adminUsersListParams) {
		params.isActive = &isActive
	}
}

// AdminUsersListOptionIncludeDeactivatedUserWorkspaces includes workspaces where a
// user is deactivated in each user's workspaces.
func AdminUsersListOptionIncludeDeactivatedUserWorkspaces(include bool) AdminUsersListOption {
	return func(params *adminUsersListParams) {
		params.includeDeactivatedUserWorkspaces = include
	}
}

// AdminUsersListOptionOnlyGuests returns only guests and their expiration dates.
func AdminUsersListOptionOnlyGuests(onlyGuests bool) AdminUsersListOption {
	return func(params *adminUsersListParams) {
		params.onlyGuests = onlyGuests
	}
}

// AdminUsersListOptionIncludeAdmins returns only admins, or admins and owners when
// combined with AdminUsersListOptionIncludeOwners.
func AdminUsersListOptionIncludeAdmins(includeAdmins bool) AdminUsersListOption {
	return func(params *adminUsersListParams) {
		params.includeAdmins = includeAdmins
	}
}

// AdminUsersListOptionIncludeOwners returns only owners.
func AdminUsersListOptionIncludeOwners(includeOwners bool) AdminUsersListOption {
	return func(params *adminUsersListParams) {
		params.includeOwners = includeOwners
	}
}

// AdminUser represents a user returned by admin.users.list.
type AdminUser struct {
	ID                   string   `json:"id"`
	Email                string   `json:"email"`
	IsAdmin              bool     `json:"is_admin"`
	IsOwner              bool     `json:"is_owner"`
	IsPrimaryOwner       bool     `json:"is_primary_owner"`
	IsRestricted         bool     `json:"is_restricted"`
	IsUltraRestricted    bool     `json:"is_ultra_restricted"`
	IsBot                bool     `json:"is_bot"`
	Username             string   `json:"username"`
	FullName             string   `json:"full_name"`
	IsActive             bool     `json:"is_active"`
	DateCreated          int64    `json:"date_created"`
	DeactivatedTimestamp int64    `json:"deactivated_ts"`
	ExpirationTimestamp  int64    `json:"expiration_ts"`
	LastActiveTimestamp  int64    `json:"last_active_ts"`
	Workspaces           []string `json:"workspaces"`
	Roles                []string `json:"roles"`
	Has2FA               bool     `json:"has_2fa"`
	HasSSO               bool     `json:"has_sso"`
}

// AdminUsersListResponse represents the response from admin.users.list.
type AdminUsersListResponse struct {
	SlackResponse
	Users []AdminUser `json:"users"`
}

// AdminUsersList lists users on a workspace or Enterprise organization.
//
// Slack API docs: https://docs.slack.dev/reference/methods/admin.users.list
func (api *Client) AdminUsersList(ctx context.Context, options ...AdminUsersListOption) (*AdminUsersListResponse, error) {
	params := adminUsersListParams{}
	for _, opt := range options {
		opt(&params)
	}

	values := url.Values{
		"token": {api.token},
	}

	if params.teamID != "" {
		values.Add("team_id", params.teamID)
	}
	if params.cursor != "" {
		values.Add("cursor", params.cursor)
	}
	if params.limit > 0 {
		values.Add("limit", strconv.Itoa(params.limit))
	}
	if params.isActive != nil {
		values.Add("is_active", strconv.FormatBool(*params.isActive))
	}
	if params.includeDeactivatedUserWorkspaces {
		values.Add("include_deactivated_user_workspaces", "true")
	}
	if params.onlyGuests {
		values.Add("only_guests", "true")
	}
	if params.includeAdmins {
		values.Add("include_admins", "true")
	}
	if params.includeOwners {
		values.Add("include_owners", "true")
	}

	response := &AdminUsersListResponse{}
	err := api.postMethod(ctx, "admin.users.list", values, response)
	if err != nil {
		return nil, err
	}

	return response, response.Err()
}

// AdminUsersSetAdmin makes a regular user or owner an admin. teamID may be a
// workspace or org ID. With an org ID, Slack also makes the user a Workspace
// Owner on every workspace in the org they belong to.
//
// Slack API docs: https://docs.slack.dev/reference/methods/admin.users.setAdmin
func (api *Client) AdminUsersSetAdmin(teamID, userID string) error {
	return api.AdminUsersSetAdminContext(context.Background(), teamID, userID)
}

// AdminUsersSetAdminContext makes a user an admin with a custom context.
func (api *Client) AdminUsersSetAdminContext(ctx context.Context, teamID, userID string) error {
	return api.adminUsersSetAccountType(ctx, "admin.users.setAdmin", teamID, userID)
}

// AdminUserExpiration is a guest returned by admin.users.getExpiration.
type AdminUserExpiration struct {
	ID                  string `json:"id"`
	Email               string `json:"email"`
	IsRestricted        bool   `json:"is_restricted"`
	IsUltraRestricted   bool   `json:"is_ultra_restricted"`
	ExpirationTimestamp int64  `json:"expiration_ts"`
}

// AdminUsersGetExpirationResponse represents the response from admin.users.getExpiration.
type AdminUsersGetExpirationResponse struct {
	SlackResponse
	User AdminUserExpiration `json:"user"`
}

type adminUsersGetExpirationParams struct {
	targetTeam string
}

// AdminUsersGetExpirationOption is an option for AdminUsersGetExpiration.
type AdminUsersGetExpirationOption func(*adminUsersGetExpirationParams)

// AdminUsersGetExpirationOptionTargetTeam returns the guest's expiration on a
// workspace in the org. It requires an org-level token.
func AdminUsersGetExpirationOptionTargetTeam(teamID string) AdminUsersGetExpirationOption {
	return func(params *adminUsersGetExpirationParams) {
		params.targetTeam = teamID
	}
}

// AdminUsersGetExpiration gets when a guest account is disabled, as Unix seconds.
//
// Slack API docs: https://docs.slack.dev/reference/methods/admin.users.getExpiration
func (api *Client) AdminUsersGetExpiration(userID string, options ...AdminUsersGetExpirationOption) (*AdminUsersGetExpirationResponse, error) {
	return api.AdminUsersGetExpirationContext(context.Background(), userID, options...)
}

// AdminUsersGetExpirationContext gets a guest's expiration with a custom context.
func (api *Client) AdminUsersGetExpirationContext(ctx context.Context, userID string, options ...AdminUsersGetExpirationOption) (*AdminUsersGetExpirationResponse, error) {
	params := adminUsersGetExpirationParams{}
	for _, opt := range options {
		if opt != nil {
			opt(&params)
		}
	}

	values := url.Values{
		"token":   {api.token},
		"user_id": {userID},
	}
	if params.targetTeam != "" {
		values.Add("target_team", params.targetTeam)
	}

	response := &AdminUsersGetExpirationResponse{}
	if err := api.postMethod(ctx, "admin.users.getExpiration", values, response); err != nil {
		return nil, err
	}
	return response, response.Err()
}

type adminUsersSetExpirationParams struct {
	teamID string
}

// AdminUsersSetExpirationOption is an option for AdminUsersSetExpiration.
type AdminUsersSetExpirationOption func(*adminUsersSetExpirationParams)

// AdminUsersSetExpirationOptionTeamID sets the guest's workspace. It is required
// with an org-level token; without it, Slack returns
// team_id_required_for_enterprise.
func AdminUsersSetExpirationOptionTeamID(teamID string) AdminUsersSetExpirationOption {
	return func(params *adminUsersSetExpirationParams) {
		params.teamID = teamID
	}
}

// AdminUsersSetExpiration sets when a guest account is disabled, as Unix seconds.
//
// Slack API docs: https://docs.slack.dev/reference/methods/admin.users.setExpiration
func (api *Client) AdminUsersSetExpiration(userID string, expirationTS int64, options ...AdminUsersSetExpirationOption) error {
	return api.AdminUsersSetExpirationContext(context.Background(), userID, expirationTS, options...)
}

// AdminUsersSetExpirationContext sets a guest's expiration with a custom context.
func (api *Client) AdminUsersSetExpirationContext(ctx context.Context, userID string, expirationTS int64, options ...AdminUsersSetExpirationOption) error {
	params := adminUsersSetExpirationParams{}
	for _, opt := range options {
		if opt != nil {
			opt(&params)
		}
	}

	values := url.Values{
		"token":         {api.token},
		"user_id":       {userID},
		"expiration_ts": {strconv.FormatInt(expirationTS, 10)},
	}
	if params.teamID != "" {
		values.Add("team_id", params.teamID)
	}

	response := &SlackResponse{}
	if err := api.postMethod(ctx, "admin.users.setExpiration", values, response); err != nil {
		return err
	}
	return response.Err()
}

// AdminUsersSetOwner makes a regular user or admin an owner. teamID may be a
// workspace or org ID. With an org ID, Slack also makes the user a Workspace
// Owner on every workspace in the org they belong to.
//
// Slack API docs: https://docs.slack.dev/reference/methods/admin.users.setOwner
func (api *Client) AdminUsersSetOwner(teamID, userID string) error {
	return api.AdminUsersSetOwnerContext(context.Background(), teamID, userID)
}

// AdminUsersSetOwnerContext makes a user an owner with a custom context.
func (api *Client) AdminUsersSetOwnerContext(ctx context.Context, teamID, userID string) error {
	return api.adminUsersSetAccountType(ctx, "admin.users.setOwner", teamID, userID)
}

// AdminUsersSetRegular makes a guest, admin, or owner a regular user. teamID may
// be a workspace or org ID. With an org ID, Slack also removes every Workspace
// Owner role the user holds in the org; restore any with AdminUsersSetOwner and a
// workspace ID. On failed_to_set_regular_on_workspaces, the org-level change has
// been applied and the call can be retried.
//
// Unlike SetRegular, this calls the documented Admin API method.
//
// Slack API docs: https://docs.slack.dev/reference/methods/admin.users.setRegular
func (api *Client) AdminUsersSetRegular(teamID, userID string) error {
	return api.AdminUsersSetRegularContext(context.Background(), teamID, userID)
}

// AdminUsersSetRegularContext makes a user a regular user with a custom context.
func (api *Client) AdminUsersSetRegularContext(ctx context.Context, teamID, userID string) error {
	return api.adminUsersSetAccountType(ctx, "admin.users.setRegular", teamID, userID)
}

func (api *Client) adminUsersSetAccountType(ctx context.Context, method, teamID, userID string) error {
	values := url.Values{
		"token":   {api.token},
		"team_id": {teamID},
		"user_id": {userID},
	}

	response := &SlackResponse{}
	if err := api.postMethod(ctx, method, values, response); err != nil {
		return err
	}
	return response.Err()
}
