package slack

import (
	"context"
	"net/url"
)

// OAuthResponseIncomingWebhook ...
type OAuthResponseIncomingWebhook struct {
	URL              string `json:"url"`
	Channel          string `json:"channel"`
	ChannelID        string `json:"channel_id,omitempty"`
	ConfigurationURL string `json:"configuration_url"`
}

// OAuthResponseBot ...
type OAuthResponseBot struct {
	BotUserID      string `json:"bot_user_id"`
	BotAccessToken string `json:"bot_access_token"`
}

// OAuthResponse ...
type OAuthResponse struct {
	AccessToken     string                       `json:"access_token"`
	Scope           string                       `json:"scope"`
	TeamName        string                       `json:"team_name"`
	TeamID          string                       `json:"team_id"`
	IncomingWebhook OAuthResponseIncomingWebhook `json:"incoming_webhook"`
	Bot             OAuthResponseBot             `json:"bot"`
	UserID          string                       `json:"user_id,omitempty"`
	SlackResponse
}

// OAuthV2Response ...
type OAuthV2Response struct {
	AccessToken     string                       `json:"access_token"`
	TokenType       string                       `json:"token_type"`
	Scope           string                       `json:"scope"`
	BotUserID       string                       `json:"bot_user_id"`
	AppID           string                       `json:"app_id"`
	Team            OAuthV2ResponseTeam          `json:"team"`
	IncomingWebhook OAuthResponseIncomingWebhook `json:"incoming_webhook"`
	Enterprise      OAuthV2ResponseEnterprise    `json:"enterprise"`
	AuthedUser      OAuthV2ResponseAuthedUser    `json:"authed_user"`
	SlackResponse
}

// OAuthV2ResponseTeam ...
type OAuthV2ResponseTeam struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// OAuthV2ResponseEnterprise ...
type OAuthV2ResponseEnterprise struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// OAuthV2ResponseAuthedUser ...
type OAuthV2ResponseAuthedUser struct {
	ID          string `json:"id"`
	Scope       string `json:"scope"`
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
}

// GetOAuthToken retrieves an AccessToken
func (api *Client) GetOAuthToken(clientID, clientSecret, code, redirectURI string) (accessToken string, scope string, err error) {
	return api.GetOAuthTokenContext(context.Background(), clientID, clientSecret, code, redirectURI)
}

// GetOAuthTokenContext retrieves an AccessToken with a custom context
func (api *Client) GetOAuthTokenContext(ctx context.Context, clientID, clientSecret, code, redirectURI string) (accessToken string, scope string, err error) {
	response, err := api.GetOAuthResponseContext(ctx, clientID, clientSecret, code, redirectURI)
	if err != nil {
		return "", "", err
	}
	return response.AccessToken, response.Scope, nil
}

// GetBotOAuthToken retrieves top-level and bot AccessToken - https://api.slack.com/legacy/oauth#bot_user_access_tokens
func (api *Client) GetBotOAuthToken(clientID, clientSecret, code, redirectURI string) (accessToken string, scope string, bot OAuthResponseBot, err error) {
	return api.GetBotOAuthTokenContext(context.Background(), clientID, clientSecret, code, redirectURI)
}

// GetBotOAuthTokenContext retrieves top-level and bot AccessToken with a custom context
func (api *Client) GetBotOAuthTokenContext(ctx context.Context, clientID, clientSecret, code, redirectURI string) (accessToken string, scope string, bot OAuthResponseBot, err error) {
	response, err := api.GetOAuthResponseContext(ctx, clientID, clientSecret, code, redirectURI)
	if err != nil {
		return "", "", OAuthResponseBot{}, err
	}
	return response.AccessToken, response.Scope, response.Bot, nil
}

// GetOAuthResponse retrieves OAuth response
func (api *Client) GetOAuthResponse(clientID, clientSecret, code, redirectURI string) (resp *OAuthResponse, err error) {
	return api.GetOAuthResponseContext(context.Background(), clientID, clientSecret, code, redirectURI)
}

// GetOAuthResponseContext retrieves OAuth response with custom context
func (api *Client) GetOAuthResponseContext(ctx context.Context, clientID, clientSecret, code, redirectURI string) (resp *OAuthResponse, err error) {
	values := url.Values{
		"client_id":     {clientID},
		"client_secret": {clientSecret},
		"code":          {code},
		"redirect_uri":  {redirectURI},
	}
	response := &OAuthResponse{}
	if err = api.postMethod(ctx, "oauth.access", values, response); err != nil {
		return nil, err
	}
	return response, response.Err()
}

// GetOAuthV2Response gets a V2 OAuth access token response - https://api.slack.com/methods/oauth.v2.access
func (api *Client) GetOAuthV2Response(clientID, clientSecret, code, redirectURI string) (resp *OAuthV2Response, err error) {
	return api.GetOAuthV2ResponseContext(context.Background(), clientID, clientSecret, code, redirectURI)
}

// GetOAuthV2ResponseContext with a context, gets a V2 OAuth access token response
func (api *Client) GetOAuthV2ResponseContext(ctx context.Context, clientID, clientSecret, code, redirectURI string) (resp *OAuthV2Response, err error) {
	values := url.Values{
		"client_id":     {clientID},
		"client_secret": {clientSecret},
		"code":          {code},
		"redirect_uri":  {redirectURI},
	}
	response := &OAuthV2Response{}
	if err = api.postMethod(ctx, "oauth.v2.access", values, response); err != nil {
		return nil, err
	}
	return response, response.Err()
}
