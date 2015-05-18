package slack

import (
	"errors"
	"net/url"
	"strconv"
)

type channelResponseFull struct {
	Channel      Channel   `json:"channel"`
	Channels     []Channel `json:"channels"`
	Purpose      string    `json:"purpose"`
	Topic        string    `json:"topic"`
	NotInChannel bool      `json:"not_in_channel"`
	History
	SlackResponse
}

// ChannelTopic contains information about the channel topic
type ChannelTopic struct {
	Value   string   `json:"value"`
	Creator string   `json:"creator"`
	LastSet JSONTime `json:"last_set"`
}

// ChannelPurpose contains information about the channel purpose
type ChannelPurpose struct {
	Value   string   `json:"value"`
	Creator string   `json:"creator"`
	LastSet JSONTime `json:"last_set"`
}

type BaseChannel struct {
	Id                 string         `json:"id"`
	Created            JSONTime       `json:"created"`
	IsOpen             bool           `json:"is_open"`
	LastRead           string         `json:"last_read,omitempty"`
	Latest             Message        `json:"latest,omitempty"`
	UnreadCount        int            `json:"unread_count,omitempty"`
	UnreadCountDisplay int            `json:"unread_count_display,omitempty"`
}

// Channel contains information about the channel
type Channel struct {
	BaseChannel
	Name               string         `json:"name"`
	IsChannel          bool           `json:"is_channel"`
	Creator            string         `json:"creator"`
	IsArchived         bool           `json:"is_archived"`
	IsGeneral          bool           `json:"is_general"`
	IsGroup            bool           `json:"is_group"`
	IsStarred          bool           `json:"is_starred"`
	Members            []string       `json:"members"`
	Topic              ChannelTopic   `json:"topic"`
	Purpose            ChannelPurpose `json:"purpose"`
	IsMember           bool           `json:"is_member"`
	NumMembers         int            `json:"num_members,omitempty"`
}

func channelRequest(path string, values url.Values, debug bool) (*channelResponseFull, error) {
	response := &channelResponseFull{}
	err := parseResponse(path, values, response, debug)
	if err != nil {
		return nil, err
	}
	if !response.Ok {
		return nil, errors.New(response.Error)
	}
	return response, nil
}

// ArchiveChannel archives the given channel
func (api *Slack) ArchiveChannel(channelId string) error {
	values := url.Values{
		"token":   {api.config.token},
		"channel": {channelId},
	}
	_, err := channelRequest("channels.archive", values, api.debug)
	if err != nil {
		return err
	}
	return nil
}

// UnarchiveChannel unarchives the given channel
func (api *Slack) UnarchiveChannel(channelId string) error {
	values := url.Values{
		"token":   {api.config.token},
		"channel": {channelId},
	}
	_, err := channelRequest("channels.unarchive", values, api.debug)
	if err != nil {
		return err
	}
	return nil
}

// CreateChannel creates a channel with the given name and returns a *Channel
func (api *Slack) CreateChannel(channel string) (*Channel, error) {
	values := url.Values{
		"token": {api.config.token},
		"name":  {channel},
	}
	response, err := channelRequest("channels.create", values, api.debug)
	if err != nil {
		return nil, err
	}
	return &response.Channel, nil
}

// GetChannelHistory retrieves the channel history
func (api *Slack) GetChannelHistory(channelId string, params HistoryParameters) (*History, error) {
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
	response, err := channelRequest("channels.history", values, api.debug)
	if err != nil {
		return nil, err
	}
	return &response.History, nil
}

// GetChannelInfo retrieves the given channel
func (api *Slack) GetChannelInfo(channelId string) (*Channel, error) {
	values := url.Values{
		"token":   {api.config.token},
		"channel": {channelId},
	}
	response, err := channelRequest("channels.info", values, api.debug)
	if err != nil {
		return nil, err
	}
	return &response.Channel, nil
}

// InviteUserToChannel invites a user to a given channel and returns a *Channel
func (api *Slack) InviteUserToChannel(channelId, userId string) (*Channel, error) {
	values := url.Values{
		"token":   {api.config.token},
		"channel": {channelId},
		"user":    {userId},
	}
	response, err := channelRequest("channels.invite", values, api.debug)
	if err != nil {
		return nil, err
	}
	return &response.Channel, nil
}

// JoinChannel joins the currently authenticated user to a channel
func (api *Slack) JoinChannel(channel string) (*Channel, error) {
	values := url.Values{
		"token": {api.config.token},
		"name":  {channel},
	}
	response, err := channelRequest("channels.join", values, api.debug)
	if err != nil {
		return nil, err
	}
	return &response.Channel, nil
}

// LeaveChannel makes the authenticated user leave the given channel
func (api *Slack) LeaveChannel(channelId string) (bool, error) {
	values := url.Values{
		"token":   {api.config.token},
		"channel": {channelId},
	}
	response, err := channelRequest("channels.leave", values, api.debug)
	if err != nil {
		return false, err
	}
	if response.NotInChannel {
		return response.NotInChannel, nil
	}
	return false, nil
}

// KickUserFromChannel kicks a user from a given channel
func (api *Slack) KickUserFromChannel(channelId, userId string) error {
	values := url.Values{
		"token":   {api.config.token},
		"channel": {channelId},
		"user":    {userId},
	}
	_, err := channelRequest("channels.kick", values, api.debug)
	if err != nil {
		return err
	}
	return nil
}

// GetChannels retrieves all the channels
func (api *Slack) GetChannels(excludeArchived bool) ([]Channel, error) {
	values := url.Values{
		"token": {api.config.token},
	}
	if excludeArchived {
		values.Add("exclude_archived", "1")
	}
	response, err := channelRequest("channels.list", values, api.debug)
	if err != nil {
		return nil, err
	}
	return response.Channels, nil
}

// SetChannelReadMark sets the read mark of a given channel to a specific point
// Clients should try to avoid making this call too often. When needing to mark a read position, a client should set a
// timer before making the call. In this way, any further updates needed during the timeout will not generate extra calls
// (just one per channel). This is useful for when reading scroll-back history, or following a busy live channel. A
// timeout of 5 seconds is a good starting point. Be sure to flush these calls on shutdown/logout.
func (api *Slack) SetChannelReadMark(channelId, ts string) error {
	values := url.Values{
		"token":   {api.config.token},
		"channel": {channelId},
		"ts":      {ts},
	}
	_, err := channelRequest("channels.mark", values, api.debug)
	if err != nil {
		return err
	}
	return nil
}

// RenameChannel renames a given channel
func (api *Slack) RenameChannel(channelId, name string) (*Channel, error) {
	values := url.Values{
		"token":   {api.config.token},
		"channel": {channelId},
		"name":    {name},
	}
	// XXX: the created entry in this call returns a string instead of a number
	// so I may have to do some workaround to solve it.
	response, err := channelRequest("channels.rename", values, api.debug)
	if err != nil {
		return nil, err
	}
	return &response.Channel, nil

}

// SetChannelPurpose sets the channel purpose and returns the purpose that was
// successfully set
func (api *Slack) SetChannelPurpose(channelId, purpose string) (string, error) {
	values := url.Values{
		"token":   {api.config.token},
		"channel": {channelId},
		"purpose": {purpose},
	}
	response, err := channelRequest("channels.setPurpose", values, api.debug)
	if err != nil {
		return "", err
	}
	return response.Purpose, nil
}

// SetChannelTopic sets the channel topic and returns the topic that was successfully set
func (api *Slack) SetChannelTopic(channelId, topic string) (string, error) {
	values := url.Values{
		"token":   {api.config.token},
		"channel": {channelId},
		"topic":   {topic},
	}
	response, err := channelRequest("channels.setTopic", values, api.debug)
	if err != nil {
		return "", err
	}
	return response.Topic, nil
}
