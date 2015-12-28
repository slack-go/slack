package slack

import (
	"errors"
	"net/http"
	"net/url"
)

type oAuthResponseFull struct {
	AccessToken string `json:"access_token"`
	Scope       string `json:"scope"`
	SlackResponse
}

// GetOAuthToken retrieves an AccessToken
func GetOAuthToken(client *http.Client, clientID, clientSecret, code, redirectURI string, debug bool) (accessToken string, scope string, err error) {
	values := url.Values{
		"client_id":     {clientID},
		"client_secret": {clientSecret},
		"code":          {code},
		"redirect_uri":  {redirectURI},
	}
	response := &oAuthResponseFull{}
	err = post(client, "oauth.access", values, response, debug)
	if err != nil {
		return "", "", err
	}
	if !response.Ok {
		return "", "", errors.New(response.Error)
	}
	return response.AccessToken, response.Scope, nil
}
