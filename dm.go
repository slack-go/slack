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

type IM struct {
	Id            string   `json:"id"`
	IsIM          bool     `json:"is_im"`
	UserId        string   `json:"user"`
	Created       JSONTime `json:"created"`
	IsUserDeleted bool     `json:"is_user_deleted"`
}

func imRequest(path string, values url.Values, debug bool) (*imResponseFull, error) {
	response := &imResponseFull{}
	err := ParseResponse(path, values, response, debug)
	if err != nil {
		return nil, err
	}
	if !response.Ok {
		return nil, errors.New(response.Error)
	}
	return response, nil
}

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
	response, err := imRequest("im.history", values, api.debug)
	if err != nil {
		return nil, err
	}
	return &response.History, nil
}

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
