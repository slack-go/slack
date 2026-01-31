package slack

import (
	"context"
	"net/url"
	"strconv"
	"strings"
)

// AdminRolesAddAssignmentsParams contains arguments for AdminRolesAddAssignments method call.
type AdminRolesAddAssignmentsParams struct {
	RoleID        string
	EntityIDs     []string
	UserIDs       []string
	DateEffective int64
}

// AdminRolesRejectedUser represents a user that could not be assigned a role.
type AdminRolesRejectedUser struct {
	ID    string `json:"id"`
	Error string `json:"error"`
}

// AdminRolesRejectedEntity represents an entity that could not be assigned a role.
type AdminRolesRejectedEntity struct {
	ID    string `json:"id"`
	Error string `json:"error"`
}

// AdminRolesAddAssignmentsResponse represents the response from admin.roles.addAssignments.
type AdminRolesAddAssignmentsResponse struct {
	SlackResponse
	RejectedUsers    []AdminRolesRejectedUser   `json:"rejected_users"`
	RejectedEntities []AdminRolesRejectedEntity `json:"rejected_entities"`
}

// AdminRolesAddAssignments adds members to a specified role.
// For more information see the admin.roles.addAssignments docs:
// https://api.slack.com/methods/admin.roles.addAssignments
func (api *Client) AdminRolesAddAssignments(ctx context.Context, params AdminRolesAddAssignmentsParams) (*AdminRolesAddAssignmentsResponse, error) {
	values := url.Values{
		"token":   {api.token},
		"role_id": {params.RoleID},
	}

	if len(params.EntityIDs) > 0 {
		values.Add("entity_ids", strings.Join(params.EntityIDs, ","))
	}

	if len(params.UserIDs) > 0 {
		values.Add("user_ids", strings.Join(params.UserIDs, ","))
	}

	if params.DateEffective > 0 {
		values.Add("date_effective", strconv.FormatInt(params.DateEffective, 10))
	}

	response := &AdminRolesAddAssignmentsResponse{}
	err := api.postMethod(ctx, "admin.roles.addAssignments", values, response)
	if err != nil {
		return nil, err
	}

	return response, response.Err()
}

// AdminRolesListAssignmentsParams contains arguments for AdminRolesListAssignments method call.
type AdminRolesListAssignmentsParams struct {
	RoleIDs []string
	Limit   int
	Cursor  string
}

// RoleAssignment represents a single role assignment.
type RoleAssignment struct {
	RoleID      string `json:"role_id"`
	EntityID    string `json:"entity_id,omitempty"`
	UserID      string `json:"user_id,omitempty"`
	DateCreated int64  `json:"date_created,omitempty"`
}

// AdminRolesListAssignmentsResponse represents the response from admin.roles.listAssignments.
type AdminRolesListAssignmentsResponse struct {
	SlackResponse
	RoleAssignments  []RoleAssignment `json:"role_assignments"`
	ResponseMetadata ResponseMetadata `json:"response_metadata"`
}

// AdminRolesListAssignments lists assignments for roles.
// For more information see the admin.roles.listAssignments docs:
// https://api.slack.com/methods/admin.roles.listAssignments
func (api *Client) AdminRolesListAssignments(ctx context.Context, params AdminRolesListAssignmentsParams) (*AdminRolesListAssignmentsResponse, error) {
	values := url.Values{
		"token": {api.token},
	}

	if len(params.RoleIDs) > 0 {
		values.Add("role_ids", strings.Join(params.RoleIDs, ","))
	}

	if params.Limit > 0 {
		values.Add("limit", strconv.Itoa(params.Limit))
	}

	if params.Cursor != "" {
		values.Add("cursor", params.Cursor)
	}

	response := &AdminRolesListAssignmentsResponse{}
	err := api.postMethod(ctx, "admin.roles.listAssignments", values, response)
	if err != nil {
		return nil, err
	}

	return response, response.Err()
}

// AdminRolesRemoveAssignmentsParams contains arguments for AdminRolesRemoveAssignments method call.
type AdminRolesRemoveAssignmentsParams struct {
	RoleID    string
	EntityIDs []string
	UserIDs   []string
}

// AdminRolesRemoveAssignmentsResponse represents the response from admin.roles.removeAssignments.
type AdminRolesRemoveAssignmentsResponse struct {
	SlackResponse
	RejectedUsers    []AdminRolesRejectedUser   `json:"rejected_users"`
	RejectedEntities []AdminRolesRejectedEntity `json:"rejected_entities"`
}

// AdminRolesRemoveAssignments removes members from a specified role.
// For more information see the admin.roles.removeAssignments docs:
// https://api.slack.com/methods/admin.roles.removeAssignments
func (api *Client) AdminRolesRemoveAssignments(ctx context.Context, params AdminRolesRemoveAssignmentsParams) (*AdminRolesRemoveAssignmentsResponse, error) {
	values := url.Values{
		"token":   {api.token},
		"role_id": {params.RoleID},
	}

	if len(params.EntityIDs) > 0 {
		values.Add("entity_ids", strings.Join(params.EntityIDs, ","))
	}

	if len(params.UserIDs) > 0 {
		values.Add("user_ids", strings.Join(params.UserIDs, ","))
	}

	response := &AdminRolesRemoveAssignmentsResponse{}
	err := api.postMethod(ctx, "admin.roles.removeAssignments", values, response)
	if err != nil {
		return nil, err
	}

	return response, response.Err()
}
