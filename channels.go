package slack

import (
	"context"
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

// Channel contains information about the channel
type Channel struct {
	groupConversation
	IsChannel bool `json:"is_channel"`
	IsGeneral bool `json:"is_general"`
	IsMember  bool `json:"is_member"`
}

func channelRequest(ctx context.Context, path string, values url.Values, debug bool) (*channelResponseFull, error) {
	response := &channelResponseFull{}
	err := post(ctx, path, values, response, debug)
	if err != nil {
		return nil, err
	}
	if !response.Ok {
		return nil, errors.New(response.Error)
	}
	return response, nil
}

// ArchiveChannel archives the given channel
// see https://api.slack.com/methods/channels.archive
func (api *Client) ArchiveChannel(channelId string) error {
	return api.ArchiveChannelContext(context.Background(), channelId)
}

// ArchiveChannelContext archives the given channel with a custom context
// see https://api.slack.com/methods/channels.archive
func (api *Client) ArchiveChannelContext(ctx context.Context, channelId string) error {
	values := url.Values{
		"token":   {api.config.token},
		"channel": {channelId},
	}
	_, err := channelRequest(ctx, "channels.archive", values, api.debug)
	if err != nil {
		return err
	}
	return nil
}

// UnarchiveChannel unarchives the given channel
// see https://api.slack.com/methods/channels.unarchive
func (api *Client) UnarchiveChannel(channelId string) error {
	return api.UnarchiveChannelContext(context.Background(), channelId)
}

// UnarchiveChannelContext unarchives the given channel with a custom context
// see https://api.slack.com/methods/channels.unarchive
func (api *Client) UnarchiveChannelContext(ctx context.Context, channelId string) error {
	values := url.Values{
		"token":   {api.config.token},
		"channel": {channelId},
	}
	_, err := channelRequest(ctx, "channels.unarchive", values, api.debug)
	if err != nil {
		return err
	}
	return nil
}

// CreateChannel creates a channel with the given name and returns a *Channel
// see https://api.slack.com/methods/channels.create
func (api *Client) CreateChannel(channelName string) (*Channel, error) {
	return api.CreateChannelContext(context.Background(), channelName)
}

// CreateChannelContext creates a channel with the given name and returns a *Channel with a custom context
// see https://api.slack.com/methods/channels.create
func (api *Client) CreateChannelContext(ctx context.Context, channelName string) (*Channel, error) {
	values := url.Values{
		"token": {api.config.token},
		"name":  {channelName},
	}
	response, err := channelRequest(ctx, "channels.create", values, api.debug)
	if err != nil {
		return nil, err
	}
	return &response.Channel, nil
}

// GetChannelHistory retrieves the channel history
// see https://api.slack.com/methods/channels.history
func (api *Client) GetChannelHistory(channelId string, params HistoryParameters) (*History, error) {
	return api.GetChannelHistoryContext(context.Background(), channelId, params)
}

// GetChannelHistoryContext retrieves the channel history with a custom context
// see https://api.slack.com/methods/channels.history
func (api *Client) GetChannelHistoryContext(ctx context.Context, channelId string, params HistoryParameters) (*History, error) {
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
	if params.Unreads != DEFAULT_HISTORY_UNREADS {
		if params.Unreads {
			values.Add("unreads", "1")
		} else {
			values.Add("unreads", "0")
		}
	}
	response, err := channelRequest(ctx, "channels.history", values, api.debug)
	if err != nil {
		return nil, err
	}
	return &response.History, nil
}

// GetChannelInfo retrieves the given channel
// see https://api.slack.com/methods/channels.info
func (api *Client) GetChannelInfo(channelId string) (*Channel, error) {
	return api.GetChannelInfoContext(context.Background(), channelId)
}

// GetChannelInfoContext retrieves the given channel with a custom context
// see https://api.slack.com/methods/channels.info
func (api *Client) GetChannelInfoContext(ctx context.Context, channelId string) (*Channel, error) {
	values := url.Values{
		"token":   {api.config.token},
		"channel": {channelId},
	}
	response, err := channelRequest(ctx, "channels.info", values, api.debug)
	if err != nil {
		return nil, err
	}
	return &response.Channel, nil
}

// InviteUserToChannel invites a user to a given channel and returns a *Channel
// see https://api.slack.com/methods/channels.invite
func (api *Client) InviteUserToChannel(channelId, user string) (*Channel, error) {
	return api.InviteUserToChannelContext(context.Background(), channelId, user)
}

// InviteUserToChannelCustom invites a user to a given channel and returns a *Channel with a custom context
// see https://api.slack.com/methods/channels.invite
func (api *Client) InviteUserToChannelContext(ctx context.Context, channelId, user string) (*Channel, error) {
	values := url.Values{
		"token":   {api.config.token},
		"channel": {channelId},
		"user":    {user},
	}
	response, err := channelRequest(ctx, "channels.invite", values, api.debug)
	if err != nil {
		return nil, err
	}
	return &response.Channel, nil
}

// JoinChannel joins the currently authenticated user to a channel
// see https://api.slack.com/methods/channels.join
func (api *Client) JoinChannel(channelName string) (*Channel, error) {
	return api.JoinChannelContext(context.Background(), channelName)
}

// JoinChannelContext joins the currently authenticated user to a channel with a custom context
// see https://api.slack.com/methods/channels.join
func (api *Client) JoinChannelContext(ctx context.Context, channelName string) (*Channel, error) {
	values := url.Values{
		"token": {api.config.token},
		"name":  {channelName},
	}
	response, err := channelRequest(ctx, "channels.join", values, api.debug)
	if err != nil {
		return nil, err
	}
	return &response.Channel, nil
}

// LeaveChannel makes the authenticated user leave the given channel
// see https://api.slack.com/methods/channels.leave
func (api *Client) LeaveChannel(channelId string) (bool, error) {
	return api.LeaveChannelContext(context.Background(), channelId)
}

// LeaveChannelContext makes the authenticated user leave the given channel with a custom context
// see https://api.slack.com/methods/channels.leave
func (api *Client) LeaveChannelContext(ctx context.Context, channelId string) (bool, error) {
	values := url.Values{
		"token":   {api.config.token},
		"channel": {channelId},
	}
	response, err := channelRequest(ctx, "channels.leave", values, api.debug)
	if err != nil {
		return false, err
	}
	if response.NotInChannel {
		return response.NotInChannel, nil
	}
	return false, nil
}

// KickUserFromChannel kicks a user from a given channel
// see https://api.slack.com/methods/channels.kick
func (api *Client) KickUserFromChannel(channelId, user string) error {
	return api.KickUserFromChannelContext(context.Background(), channelId, user)
}

// KickUserFromChannelContext kicks a user from a given channel with a custom context
// see https://api.slack.com/methods/channels.kick
func (api *Client) KickUserFromChannelContext(ctx context.Context, channelId, user string) error {
	values := url.Values{
		"token":   {api.config.token},
		"channel": {channelId},
		"user":    {user},
	}
	_, err := channelRequest(ctx, "channels.kick", values, api.debug)
	if err != nil {
		return err
	}
	return nil
}

// GetChannels retrieves all the channels
// see https://api.slack.com/methods/channels.list
func (api *Client) GetChannels(excludeArchived bool) ([]Channel, error) {
	return api.GetChannelsContext(context.Background(), excludeArchived)
}

// GetChannelsContext retrieves all the channels with a custom context
// see https://api.slack.com/methods/channels.list
func (api *Client) GetChannelsContext(ctx context.Context, excludeArchived bool) ([]Channel, error) {
	values := url.Values{
		"token": {api.config.token},
	}
	if excludeArchived {
		values.Add("exclude_archived", "1")
	}
	response, err := channelRequest(ctx, "channels.list", values, api.debug)
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
// see https://api.slack.com/methods/channels.mark
func (api *Client) SetChannelReadMark(channelId, ts string) error {
	return api.SetChannelReadMarkContext(context.Background(), channelId, ts)
}

// SetChannelReadMarkContext sets the read mark of a given channel to a specific point with a custom context
// For more details see SetChannelReadMark documentation
// see https://api.slack.com/methods/channels.mark
func (api *Client) SetChannelReadMarkContext(ctx context.Context, channelId, ts string) error {
	values := url.Values{
		"token":   {api.config.token},
		"channel": {channelId},
		"ts":      {ts},
	}
	_, err := channelRequest(ctx, "channels.mark", values, api.debug)
	if err != nil {
		return err
	}
	return nil
}

// RenameChannel renames a given channel
// see https://api.slack.com/methods/channels.rename
func (api *Client) RenameChannel(channelId, name string) (*Channel, error) {
	return api.RenameChannelContext(context.Background(), channelId, name)
}

// RenameChannelContext renames a given channel with a custom context
// see https://api.slack.com/methods/channels.rename
func (api *Client) RenameChannelContext(ctx context.Context, channelId, name string) (*Channel, error) {
	values := url.Values{
		"token":   {api.config.token},
		"channel": {channelId},
		"name":    {name},
	}
	// XXX: the created entry in this call returns a string instead of a number
	// so I may have to do some workaround to solve it.
	response, err := channelRequest(ctx, "channels.rename", values, api.debug)
	if err != nil {
		return nil, err
	}
	return &response.Channel, nil
}

// SetChannelPurpose sets the channel purpose and returns the purpose that was successfully set
// see https://api.slack.com/methods/channels.setPurpose
func (api *Client) SetChannelPurpose(channelId, purpose string) (string, error) {
	return api.SetChannelPurposeContext(context.Background(), channelId, purpose)
}

// SetChannelPurposeContext sets the channel purpose and returns the purpose that was successfully set with a custom context
// see https://api.slack.com/methods/channels.setPurpose
func (api *Client) SetChannelPurposeContext(ctx context.Context, channelId, purpose string) (string, error) {
	values := url.Values{
		"token":   {api.config.token},
		"channel": {channelId},
		"purpose": {purpose},
	}
	response, err := channelRequest(ctx, "channels.setPurpose", values, api.debug)
	if err != nil {
		return "", err
	}
	return response.Purpose, nil
}

// SetChannelTopic sets the channel topic and returns the topic that was successfully set
// see https://api.slack.com/methods/channels.setTopic
func (api *Client) SetChannelTopic(channelId, topic string) (string, error) {
	return api.SetChannelTopicContext(context.Background(), channelId, topic)
}

// SetChannelTopicContext sets the channel topic and returns the topic that was successfully set with a custom context
// see https://api.slack.com/methods/channels.setTopic
func (api *Client) SetChannelTopicContext(ctx context.Context, channelId, topic string) (string, error) {
	values := url.Values{
		"token":   {api.config.token},
		"channel": {channelId},
		"topic":   {topic},
	}
	response, err := channelRequest(ctx, "channels.setTopic", values, api.debug)
	if err != nil {
		return "", err
	}
	return response.Topic, nil
}

// GetChannelReplies gets an entire thread (a message plus all the messages in reply to it).
// see https://api.slack.com/methods/channels.replies
func (api *Client) GetChannelReplies(channelId, thread_ts string) ([]Message, error) {
	return api.GetChannelRepliesContext(context.Background(), channelId, thread_ts)
}

// GetChannelRepliesContext gets an entire thread (a message plus all the messages in reply to it) with a custom context
// see https://api.slack.com/methods/channels.replies
func (api *Client) GetChannelRepliesContext(ctx context.Context, channelId, thread_ts string) ([]Message, error) {
	values := url.Values{
		"token":     {api.config.token},
		"channel":   {channelId},
		"thread_ts": {thread_ts},
	}
	response, err := channelRequest(ctx, "channels.replies", values, api.debug)
	if err != nil {
		return nil, err
	}
	return response.History.Messages, nil
}
