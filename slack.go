package slack

import (
	"errors"
	"log"
	"net/url"
)

// Well-known slack error strings.
const (
	TokenRevokedError    = "token_revoked"     // oauth token revoked by user
	ChannelNotFoundError = "channel_not_found" // Value passed for channel was invalid.
	NotInChannelError    = "not_in_channel"    // Cannot post user messages to a channel they are not in.
	IsArchivedError      = "is_archived"       // Channel has been archived.
	MsgTooLongError      = "msg_too_long"      // Message text is too long
	NoTextError          = "no_text"           // No message text provided
	RateLimitedError     = "rate_limited"      // Application has posted too many messages, read the Rate Limit documentation for more information
	NotAuthedError       = "not_authed"        // No authentication token provided.
	InvalidAuthError     = "invalid_auth"      // Invalid authentication token.
	AccountInactiveError = "account_inactive"  // Authentication token is for a deleted user or team.
)

/*
  Added as a var so that we can change this for testing purposes
*/
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

func New(token string) *Client {
	s := &Client{}
	s.config.token = token
	return s
}

// AuthTest tests if the user is able to do authenticated requests or not
func (api *Client) AuthTest() (response *AuthTestResponse, error error) {
	responseFull := &authTestResponseFull{}
	err := post("auth.test", url.Values{"token": {api.config.token}}, responseFull, api.debug)
	if err != nil {
		return nil, err
	}
	if !responseFull.Ok {
		return nil, errors.New(responseFull.Error)
	}
	return &responseFull.AuthTestResponse, nil
}

// SetDebug switches the api into debug mode
// When in debug mode, it logs various info about what its doing
// If you ever use this in production, don't call SetDebug(true)
func (api *Client) SetDebug(debug bool) {
	api.debug = debug
}

func (api *Client) Debugf(format string, v ...interface{}) {
	if api.debug {
		log.Printf(format, v...)
	}
}

func (api *Client) Debugln(v ...interface{}) {
	if api.debug {
		log.Println(v...)
	}
}
