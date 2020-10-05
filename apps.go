package slack

import (
	"context"
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

// ListEventAuthorizations lists authed users and teams for the given event_context. You must provide an app-level token to the client using OptionAppLevelToken. More info: https://api.slack.com/methods/apps.event.authorizations.list
func (api *Client) ListEventAuthorizations(eventContext string) ([]EventAuthorization, error) {
	resp := &listEventAuthorizationsResponse{}

	err := api.postMethodWithBearerToken(context.Background(), "apps.event.authorizations.list", url.Values{
		"event_context": {event_context},
	}, &resp, api.appLevelToken)

	if err != nil {
		return nil, err
	}
	if !resp.Ok {
		return nil, resp.Err()
	}

	return resp.Authorizations, nil
}
