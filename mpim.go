package slack

import (
	"errors"
	"net/url"
	"strconv"
	"strings"
)

type mpimResponseFull struct {
	NoOp          bool   `json:"no_op,omitempty"`
	AlreadyClosed bool   `json:"already_closed,omitempty"`
	MPIM          MPIM   `json:"group,omitempty"`
	Groups        []MPIM `json:"groups,omitempty"`
	History
	SlackResponse
}

// MPIM contains information related to the Multi Party IM
type MPIM struct {
	Group
	IsMPIM bool `json:"is_mpim"`
}

func mpimRequest(path string, values url.Values, debug bool) (*mpimResponseFull, error) {
	response := &mpimResponseFull{}
	err := post(path, values, response, debug)
	if err != nil {
		return nil, err
	}
	if !response.Ok {
		return nil, errors.New(response.Error)
	}
	return response, nil
}

// CloseMPIMChannel closes the direct message channel
func (api *Client) CloseMPIMChannel(channel string) (bool, bool, error) {
	values := url.Values{
		"token":   {api.config.token},
		"channel": {channel},
	}
	response, err := mpimRequest("mpim.close", values, api.debug)
	if err != nil {
		return false, false, err
	}
	return response.NoOp, response.AlreadyClosed, nil
}

// OpenMPIMChannel opens a multi party IM
// Returns some status and the channel ID
func (api *Client) OpenMPIMChannel(users []string) (*MPIM, error) {
	values := url.Values{
		"token": {api.config.token},
		"users": {strings.Join(users, ",")},
	}
	response, err := mpimRequest("mpim.open", values, api.debug)
	if err != nil {
		return nil, err
	}
	return &response.MPIM, nil
}

// MarkMPIMChannel sets the read mark of a multiparty channel to a specific point
func (api *Client) MarkMPIMChannel(channel, ts string) (err error) {
	values := url.Values{
		"token":   {api.config.token},
		"channel": {channel},
		"ts":      {ts},
	}
	_, err = mpimRequest("mpim.mark", values, api.debug)
	if err != nil {
		return err
	}
	return
}

// GetMPIMHistory retrieves the multiparty channel history
func (api *Client) GetMPIMHistory(channel string, params HistoryParameters) (*History, error) {
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
	response, err := mpimRequest("mpim.history", values, api.debug)
	if err != nil {
		return nil, err
	}
	return &response.History, nil
}

// GetMPIMChannels returns the list of multiparty channels
func (api *Client) GetMPIMChannels() ([]MPIM, error) {
	values := url.Values{
		"token": {api.config.token},
	}
	response, err := mpimRequest("mpim.list", values, api.debug)
	if err != nil {
		return nil, err
	}
	return response.Groups, nil
}
