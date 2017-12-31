package slack

import (
	"context"
	"errors"
	"log"
	"net/http"
	"net/url"
	"os"
	"fmt"
)

var (
	logger *log.Logger // A logger that can be set by consumers
)

/*
  Added as a var so that we can change this for testing purposes
*/
var SLACK_API string = "https://slack.com/api/"
var SLACK_WEB_API_FORMAT string = "https://%s.slack.com/api/users.admin.%s?t=%s"

// HTTPClient sets a custom http.Client
// deprecated: in favor of SetHTTPClient()
var HTTPClient = &http.Client{}

var customHTTPClient HTTPRequester = HTTPClient

// HTTPRequester defines the minimal interface needed for an http.Client to be implemented.
//
// Use it in conjunction with the SetHTTPClient function to allow for other capabilities
// like a tracing http.Client
type HTTPRequester interface {
	Do(*http.Request) (*http.Response, error)
}

// SetHTTPClient allows you to specify a custom http.Client
// Use this instead of the package level HTTPClient variable if you want to use a custom client like the
// Stackdriver Trace HTTPClient https://godoc.org/cloud.google.com/go/trace#HTTPClient
func SetHTTPClient(client HTTPRequester) {
	customHTTPClient = client
}

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
	token      string
	info       Info
	debug      bool
	httpclient HTTPRequester
}

// SetLogger let's library users supply a logger, so that api debugging
// can be logged along with the application's debugging info.
func SetLogger(l *log.Logger) {
	logger = l
}

// Option defines an option for a Client
type Option func(*Client)

// OptionHTTPClient - provide a custom http client to the slack client.
func OptionHTTPClient(c HTTPRequester) func(*Client) {
	return func(s *Client) {
		s.httpclient = c
	}
}

// New builds a slack client from the provided token and options.
func New(token string, options ...Option) *Client {
	s := &Client{
		token:      token,
		httpclient: customHTTPClient,
	}

	for _, opt := range options {
		opt(s)
	}

	return s
}

// AuthTest tests if the user is able to do authenticated requests or not
func (api *Client) AuthTest() (response *AuthTestResponse, error error) {
	return api.AuthTestContext(context.Background())
}

// AuthTestContext tests if the user is able to do authenticated requests or not with a custom context
func (api *Client) AuthTestContext(ctx context.Context) (response *AuthTestResponse, error error) {
	responseFull := &authTestResponseFull{}
	err := post(ctx, api.httpclient, "auth.test", url.Values{"token": {api.token}}, responseFull, api.debug)
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
	if debug && logger == nil {
		logger = log.New(os.Stdout, "nlopes/slack", log.LstdFlags|log.Lshortfile)
	}
}

func (api *Client) Debugf(format string, v ...interface{}) {
	if api.debug {
		logger.Output(2,fmt.Sprintf(format, v...))
	}
}

func (api *Client) Debugln(v ...interface{}) {
	if api.debug {
		logger.Output(2, fmt.Sprintln(v...))
	}
}
