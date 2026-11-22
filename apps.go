package slack

import (
	"context"
	"encoding/json"
	"net/url"
)

type listEventAuthorizationsResponse struct {
	SlackResponse
	Authorizations []EventAuthorization `json:"authorizations"`
}

type EventAuthorization struct {
	EnterpriseID        string `json:"enterprise_id"`
	TeamID              string `json:"team_id"`
	UserID              string `json:"user_id"`
	IsBot               bool   `json:"is_bot"`
	IsEnterpriseInstall bool   `json:"is_enterprise_install"`
}

// ListEventAuthorizations lists all authorized users and teams for the given event context.
// You must provide an app-level token to the client using OptionAppLevelToken.
// For more details, see ListEventAuthorizationsContext documentation.
func (api *Client) ListEventAuthorizations(eventContext string) ([]EventAuthorization, error) {
	return api.ListEventAuthorizationsContext(context.Background(), eventContext)
}

// ListEventAuthorizationsContext lists all authorized users and teams for the given event context.
// It follows Slack's pagination cursors using ctx for every request.
//
// Slack API docs: https://docs.slack.dev/reference/methods/apps.event.authorizations.list/
func (api *Client) ListEventAuthorizationsContext(ctx context.Context, eventContext string) ([]EventAuthorization, error) {
	params := map[string]string{"event_context": eventContext}
	var authorizations []EventAuthorization
	for {
		resp := &listEventAuthorizationsResponse{}
		request, _ := json.Marshal(params)
		if err := api.postJSONMethod(ctx, "apps.event.authorizations.list", api.appLevelToken, request, resp); err != nil {
			return nil, err
		}
		if !resp.Ok {
			return nil, resp.Err()
		}

		if authorizations == nil {
			authorizations = resp.Authorizations
		} else {
			authorizations = append(authorizations, resp.Authorizations...)
		}
		cursor := resp.ResponseMetadata.Cursor
		if cursor == "" {
			return authorizations, nil
		}
		params["cursor"] = cursor
	}
}

// UninstallApp uninstalls your app from a workspace.
// For more details, see UninstallAppContext documentation.
func (api *Client) UninstallApp(clientID, clientSecret string) error {
	return api.UninstallAppContext(context.Background(), clientID, clientSecret)
}

// UninstallAppContext uninstalls your app from a workspace with a custom context.
// Slack API docs: https://api.slack.com/methods/apps.uninstall
func (api *Client) UninstallAppContext(ctx context.Context, clientID, clientSecret string) error {
	values := url.Values{
		"token":         {api.token},
		"client_id":     {clientID},
		"client_secret": {clientSecret},
	}

	response := SlackResponse{}

	err := api.postMethod(ctx, "apps.uninstall", values, &response)
	if err != nil {
		return err
	}

	return response.Err()
}
