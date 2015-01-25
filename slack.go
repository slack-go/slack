package slack

import (
	"errors"
	"net/url"
)

/*
  Added as a var so that we can change this for testing purposes
*/
var SLACK_API string = "https://slack.com/api/"

type UserTyping struct {
	Type      string `json:"type"`
	UserId    string `json:"user"`
	ChannelId string `json:"channel"`
}

type SlackEvent struct {
	Type int
	Data interface{}
}

type SlackResponse struct {
	Ok    bool   `json:"ok"`
	Error string `json:"error"`
}

type AuthTestResponse struct {
	Url    string `json:"url"`
	Team   string `json:"team"`
	User   string `json:"user"`
	TeamId string `json:"team_id"`
	UserId string `json:"user_id"`
}

type authTestResponseFull struct {
	SlackResponse
	AuthTestResponse
}

type Slack struct {
	config Config
	info   Info
	debug  bool
}

func New(token string) *Slack {
	return &Slack{
		config: Config{token: token},
	}
}

func (api *Slack) GetInfo() Info {
	return api.info
}

func (api *Slack) AuthTest() (response *AuthTestResponse, error error) {
	response_full := &authTestResponseFull{}
	err := parseResponse("auth.test", url.Values{"token": {api.config.token}}, response_full, api.debug)
	if err != nil {
		return nil, err
	}
	if !response_full.Ok {
		return nil, errors.New(response_full.Error)
	}
	return &response_full.AuthTestResponse, nil
}

func (api *Slack) SetDebug(debug bool) {
	api.debug = debug
}
