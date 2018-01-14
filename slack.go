package slack

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/url"
	"os"
)

// Added as a var so that we can change this for testing purposes
var SLACK_API string = "https://slack.com/api/"
var SLACK_WEB_API_FORMAT string = "https://%s.slack.com/api/users.admin.%s?t=%s"

type SlackResponse struct {
	Ok    bool   `json:"ok"`
	Error string `json:"error"`
}

type AuthTestResponse struct {
	URL    string `json:"url"`
	Team   string `json:"team"`
	User   string `json:"user"`
	TeamID string `json:"team_id"`
	UserID string `json:"user_id"`
}

type authTestResponseFull struct {
	SlackResponse
	AuthTestResponse
}

type Client struct {
	config struct {
		token string
	}
	info  Info
	debug bool
}

// New creates new Client.
func New(token string) *Client {
	s := &Client{}
	s.config.token = token
	return s
}

// AuthTest tests if the user is able to do authenticated requests or not
func (api *Client) AuthTest() (response *AuthTestResponse, error error) {
	return api.AuthTestContext(context.Background())
}

// AuthTestContext tests if the user is able to do authenticated requests or not with a custom context
func (api *Client) AuthTestContext(ctx context.Context) (response *AuthTestResponse, error error) {
	api.Debugf("Challenging auth...")
	responseFull := &authTestResponseFull{}
	err := post(ctx, "auth.test", url.Values{"token": {api.config.token}}, responseFull, api.debug)
	if err != nil {
		api.Debugf("failed to test for auth: %s", err)
		return nil, err
	}
	if !responseFull.Ok {
		api.Debugf("auth response was not Ok: %s", responseFull.Error)
		return nil, errors.New(responseFull.Error)
	}

	api.Debugf("Auth challenge was successful with response %+v", responseFull.AuthTestResponse)
	return &responseFull.AuthTestResponse, nil
}

// SetDebug switches the api into debug mode
// When in debug mode, it logs various info about what its doing
// If you ever use this in production, don't call SetDebug(true)
func (api *Client) SetDebug(debug bool) {
	api.debug = debug
	if debug && logger == nil {
		SetLogger(log.New(os.Stdout, "nlopes/slack", log.LstdFlags|log.Lshortfile))
	}
}

// Debugf print a formatted debug line.
func (api *Client) Debugf(format string, v ...interface{}) {
	if api.debug {
		logger.Output(2, fmt.Sprintf(format, v...))
	}
}

// Debugln print a debug line.
func (api *Client) Debugln(v ...interface{}) {
	if api.debug {
		logger.Output(2, fmt.Sprintln(v...))
	}
}
