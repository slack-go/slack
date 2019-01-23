package slack

import (
	"context"
	"net/url"
)

// AppsUninstallResponse contains our Auth response from the auth.revoke endpoint
type AppsUninstallResponse struct {
	SlackResponse // Contains the "ok", and "Error", if any
}

// appsRequest sends the actual request, and unmarshals the response
func appsRequest(ctx context.Context, client httpClient, path string, values url.Values, d debug) (*AppsUninstallResponse, error) {
	response := &AppsUninstallResponse{}
	err := postSlackMethod(ctx, client, path, values, response, d)
	if err != nil {
		return nil, err
	}

	return response, response.Err()
}

// SendAppsUninstall will send a revocation for our token. We are re-supplying the token here, but, if
// "" is sent instead, it will use the token stored in New()
func (api *Client) SendAppsUninstall(clientID, clientSecret, token string) (*AppsUninstallResponse, error) {
	return api.SendAppsUninstallContext(context.Background(), clientID, clientSecret, token)
}

// SendAppsUninstallContext will retrieve the satus from api.test
func (api *Client) SendAppsUninstallContext(ctx context.Context, clientID, clientSecret, token string) (*AppsUninstallResponse, error) {
	// Allow a default token, but, if user specifies one, use that one.
	if token == "" {
		token = api.token
	}

	values := url.Values{
		"token":         {token},
		"client_id":     {clientID},
		"client_secret": {clientSecret},
	}

	return appsRequest(ctx, api.httpclient, "apps.uninstall", values, api)
}
