package slack

import (
	"context"
	"errors"
	"net/url"
	"strconv"
	"strings"
)

type mpimChannel struct {
	ID string `json:"id"`
}

type mpimResponseFull struct {
	NoOp          bool        `json:"no_op"`
	AlreadyClosed bool        `json:"already_closed"`
	AlreadyOpen   bool        `json:"already_open"`
	Channel       mpimChannel `json:"channel"`
	MpIMs         []MpIM      `json:"groups"`
	History
	SlackResponse
}

// MpIM contains information about a multiparty IM.
type MpIM Group

func mpimRequest(ctx context.Context, path string, values url.Values, debug bool) (*mpimResponseFull, error) {
	response := &mpimResponseFull{}
	err := post(ctx, path, values, response, debug)
	if err != nil {
		return nil, err
	}
	if !response.Ok {
		return nil, errors.New(response.Error)
	}
	return response, nil
}

// CloseIMChannel closes a multiparty direct message channel.
func (api *Client) CloseMpIMChannel(channel string) (bool, bool, error) {
	return api.CloseMpIMChannelContext(context.Background(), channel)
}

// CloseIMChannelContext closes a multiparty direct message channel.
func (api *Client) CloseMpIMChannelContext(ctx context.Context, channel string) (bool, bool, error) {
	values := url.Values{
		"token":   {api.config.token},
		"channel": {channel},
	}
	response, err := mpimRequest(ctx, "mpim.close", values, api.debug)
	if err != nil {
		return false, false, err
	}
	return response.NoOp, response.AlreadyClosed, nil
}

// OpenIMChannel opens a multiparty direct message.
func (api *Client) OpenMpIMChannel(users []string) (bool, bool, string, error) {
	return api.OpenMpIMChannelContext(context.Background(), users)
}

// OpenIMChannelContext opens a multiparty direct message.
func (api *Client) OpenMpIMChannelContext(ctx context.Context, users []string) (bool, bool, string, error) {
	usersJoin := strings.Join(users,",") // Comma separated lists of users.
	values := url.Values{
		"token": {api.config.token},
		"users":  {usersJoin},
	}
	response, err := mpimRequest(ctx, "mpim.open", values, api.debug)
	if err != nil {
		return false, false, "", err
	}
	return response.NoOp, response.AlreadyOpen, response.Channel.ID, nil
}

// MarkMpIMChannel sets the read cursor in a multiparty direct message channel.
func (api *Client) MarkMpIMChannel(channel, ts string) (err error) {
	return api.MarkIMChannelContext(context.Background(), channel, ts)
}

// MarkMpIMChannelContext sets the read cursor in a multiparty direct message channel.
func (api *Client) MarkMpIMChannelContext(ctx context.Context, channel, ts string) (err error) {
	values := url.Values{
		"token":   {api.config.token},
		"channel": {channel},
		"ts":      {ts},
	}
	_, err = mpimRequest(ctx, "mpim.mark", values, api.debug)
	if err != nil {
		return err
	}
	return
}

// GetMpIMHistory retrieves the history of messages and events from a multiparty direct message.
func (api *Client) GetMpIMHistory(channel string, params HistoryParameters) (*History, error) {
	return api.GetMpIMHistoryContext(context.Background(), channel, params)
}

// GetIMHistoryContext retrieves the history of messages and events from a multiparty direct message.
func (api *Client) GetMpIMHistoryContext(ctx context.Context, channel string, params HistoryParameters) (*History, error) {
	values := url.Values{
		"token":   {api.config.token},
		"channel": {channel},
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
	if params.Unreads != DEFAULT_HISTORY_UNREADS {
		if params.Unreads {
			values.Add("unreads", "1")
		} else {
			values.Add("unreads", "0")
		}
	}
	response, err := mpimRequest(ctx, "mpim.history", values, api.debug)
	if err != nil {
		return nil, err
	}
	return &response.History, nil
}

// GetIMChannels lists multiparty direct message channels for the calling user.
func (api *Client) GetMpIMChannels() ([]MpIM, error) {
	return api.GetMpIMChannelsContext(context.Background())
}

// GetMpIMChannelsContext lists multiparty direct message channels for the calling user.
func (api *Client) GetMpIMChannelsContext(ctx context.Context) ([]MpIM, error) {
	values := url.Values{
		"token": {api.config.token},
	}
	response, err := mpimRequest(ctx, "mpim.list", values, api.debug)
	if err != nil {
		return nil, err
	}
	return response.MpIMs, nil
}
