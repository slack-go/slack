package slack

import (
	"errors"
	"net/url"
	"strconv"
)

type imChannel struct {
	Id string `json:"id"`
}

type imResponseFull struct {
	NoOp          bool      `json:"no_op"`
	AlreadyClosed bool      `json:"already_closed"`
	AlreadyOpen   bool      `json:"already_open"`
	Channel       imChannel `json:"channel"`
	IMs           []IM      `json:"ims"`
	History
	SlackResponse
}

// IM contains information related to the Direct Message channel
type IM struct {
	BaseChannel
	IsIM               bool     `json:"is_im"`
	UserId             string   `json:"user"`
	IsUserDeleted      bool     `json:"is_user_deleted"`
}

func imRequest(path string, values url.Values, debug bool) (*imResponseFull, error) {
	response := &imResponseFull{}
	err := parseResponse(path, values, response, debug)
	if err != nil {
		return nil, err
	}
	if !response.Ok {
		return nil, errors.New(response.Error)
	}
	return response, nil
}

// CloseIMChannel closes the direct message channel
func (api *Slack) CloseIMChannel(channelId string) (bool, bool, error) {
	values := url.Values{
		"token":   {api.config.token},
		"channel": {channelId},
	}
	response, err := imRequest("im.close", values, api.debug)
	if err != nil {
		return false, false, err
	}
	return response.NoOp, response.AlreadyClosed, nil
}

// OpenIMChannel opens a direct message channel to the user provided as argument
// Returns some status and the channelId
func (api *Slack) OpenIMChannel(userId string) (bool, bool, string, error) {
	values := url.Values{
		"token": {api.config.token},
		"user":  {userId},
	}
	response, err := imRequest("im.open", values, api.debug)
	if err != nil {
		return false, false, "", err
	}
	return response.NoOp, response.AlreadyOpen, response.Channel.Id, nil
}

// MarkIMChannel sets the read mark of a direct message channel to a specific point
func (api *Slack) MarkIMChannel(channelId, ts string) (err error) {
	values := url.Values{
		"token":   {api.config.token},
		"channel": {channelId},
		"ts":      {ts},
	}
	_, err = imRequest("im.mark", values, api.debug)
	if err != nil {
		return err
	}
	return
}

// GetIMHistory retrieves the direct message channel history
func (api *Slack) GetIMHistory(channelId string, params HistoryParameters) (*History, error) {
	values := url.Values{
		"token":   {api.config.token},
		"channel": {channelId},
	}
	if params.Latest != DEFAULT_HISTORY_LATEST {
		values.Add("latest", params.Latest)
	}
	if params.Oldest != DEFAULT_HISTORY_OLDEST {
		values.Add("oldest", params.Oldest)
	}
	if params.Count != DEFAULT_HISTORY_COUNT {
		values.Add("count", strconv.Itoa(params.Count))
	}
	if params.Inclusive != DEFAULT_HISTORY_INCLUSIVE {
		if params.Inclusive {
			values.Add("inclusive", "1")
		} else {
			values.Add("inclusive", "0")
		}
	}
	response, err := imRequest("im.history", values, api.debug)
	if err != nil {
		return nil, err
	}
	return &response.History, nil
}

// GetIMChannels returns the list of direct message channels
func (api *Slack) GetIMChannels() ([]IM, error) {
	values := url.Values{
		"token": {api.config.token},
	}
	response, err := imRequest("im.list", values, api.debug)
	if err != nil {
		return nil, err
	}
	return response.IMs, nil
}
