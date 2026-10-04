package slack

import "context"

// AdminAppRequestScope describes one OAuth scope in an app request.
type AdminAppRequestScope struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	IsSensitive bool   `json:"is_sensitive"`
	TokenType   string `json:"token_type"`
	IsOptional  *bool  `json:"is_optional"`
	IsApproved  *bool  `json:"is_approved"`
}

// AdminAppPreviousResolution describes the prior decision on an app request.
type AdminAppPreviousResolution struct {
	Status string                 `json:"status"`
	Scopes []AdminAppRequestScope `json:"scopes"`
}

// AdminAppRequestUser identifies the user who requested an app.
type AdminAppRequestUser struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// AdminAppRequestTeam identifies the workspace where an app was requested.
type AdminAppRequestTeam struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Domain string `json:"domain"`
}

// AdminAppRequest is one pending app installation request.
type AdminAppRequest struct {
	ID                    string                      `json:"id"`
	App                   AdminApp                    `json:"app"`
	PreviousResolution    *AdminAppPreviousResolution `json:"previous_resolution"`
	User                  AdminAppRequestUser         `json:"user"`
	Team                  AdminAppRequestTeam         `json:"team"`
	Scopes                []AdminAppRequestScope      `json:"scopes"`
	Message               string                      `json:"message"`
	IsUserAppCollaborator bool                        `json:"is_user_app_collaborator"`
	DateCreated           int64                       `json:"date_created"`
}

// AdminAppsRequestsListResponse represents the response from admin.apps.requests.list.
type AdminAppsRequestsListResponse struct {
	SlackResponse
	AppRequests []AdminAppRequest `json:"app_requests"`
}

// AdminAppsRequestsList lists pending app installation requests for an Enterprise organization or workspace.
//
// Slack API docs: https://docs.slack.dev/reference/methods/admin.apps.requests.list
func (api *Client) AdminAppsRequestsList(ctx context.Context, options ...AdminAppsListOption) (*AdminAppsRequestsListResponse, error) {
	params := adminAppsListParams{}
	for _, option := range options {
		option(&params)
	}

	response := &AdminAppsRequestsListResponse{}
	if err := api.postMethod(ctx, "admin.apps.requests.list", adminAppsListValues(api.token, params), response); err != nil {
		return nil, err
	}
	return response, response.Err()
}
