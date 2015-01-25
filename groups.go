package slack

import (
	"errors"
	"net/url"
	"strconv"
)

type Group struct {
	Id          string         `json:"id"`
	Name        string         `json:"name"`
	IsGroup     bool           `json:"is_group"`
	Created     JSONTime       `json:"created"`
	Creator     string         `json:"creator"`
	IsArchived  bool           `json:"is_archived"`
	IsGeneral   bool           `json:"is_general"`
	IsOpen      bool           `json:"is_open,omitempty"`
	Members     []string       `json:"members"`
	Topic       ChannelTopic   `json:"topic"`
	Purpose     ChannelPurpose `json:"purpose"`
	LastRead    string         `json:"last_read,omitempty"`
	Latest      Message        `json:"latest,omitempty"`
	UnreadCount int            `json:"unread_count,omitempty"`
	NumMembers  int            `json:"num_members,omitempty"`

	// XXX: does this exist for a group too?
	IsMember bool `json:"is_member"`
}

type groupResponseFull struct {
	Group          Group   `json:"group"`
	Groups         []Group `json:"groups"`
	Purpose        string  `json:"purpose"`
	Topic          string  `json:"topic"`
	NotInGroup     bool    `json:"not_in_group"`
	NoOp           bool    `json:"no_op"`
	AlreadyClosed  bool    `json:"already_closed"`
	AlreadyOpen    bool    `json:"already_open"`
	AlreadyInGroup bool    `json:"already_in_group"`
	Channel        Channel `json:"channel"`
	History
	SlackResponse
}

func groupRequest(path string, values url.Values, debug bool) (*groupResponseFull, error) {
	response := &groupResponseFull{}
	err := parseResponse(path, values, response, debug)
	if err != nil {
		return nil, err
	}
	if !response.Ok {
		return nil, errors.New(response.Error)
	}
	return response, nil
}

func (api *Slack) ArchiveGroup(groupId string) error {
	values := url.Values{
		"token":   {api.config.token},
		"channel": {groupId},
	}
	_, err := groupRequest("groups.archive", values, api.debug)
	if err != nil {
		return err
	}
	return nil
}

func (api *Slack) UnarchiveGroup(groupId string) error {
	values := url.Values{
		"token":   {api.config.token},
		"channel": {groupId},
	}
	_, err := groupRequest("groups.unarchive", values, api.debug)
	if err != nil {
		return err
	}
	return nil
}

func (api *Slack) CreateGroup(group string) (*Group, error) {
	values := url.Values{
		"token": {api.config.token},
		"name":  {group},
	}
	response, err := groupRequest("groups.create", values, api.debug)
	if err != nil {
		return nil, err
	}
	return &response.Group, nil
}

func (api *Slack) CreateChildGroup(groupId string) (*Group, error) {
	values := url.Values{
		"token":   {api.config.token},
		"channel": {groupId},
	}
	response, err := groupRequest("groups.createChild", values, api.debug)
	if err != nil {
		return nil, err
	}
	return &response.Group, nil
}

func (api *Slack) CloseGroup(groupId string) (bool, bool, error) {
	values := url.Values{
		"token":   {api.config.token},
		"channel": {groupId},
	}
	response, err := imRequest("groups.close", values, api.debug)
	if err != nil {
		return false, false, err
	}
	return response.NoOp, response.AlreadyClosed, nil
}

func (api *Slack) GetGroupHistory(groupId string, params HistoryParameters) (*History, error) {
	values := url.Values{
		"token":   {api.config.token},
		"channel": {groupId},
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
	response, err := groupRequest("groups.history", values, api.debug)
	if err != nil {
		return nil, err
	}
	groupHistory := response.History
	return &groupHistory, nil
}

func (api *Slack) InviteUserToGroup(groupId, userId string) (*Group, bool, error) {
	values := url.Values{
		"token":   {api.config.token},
		"channel": {groupId},
		"user":    {userId},
	}
	response, err := groupRequest("groups.invite", values, api.debug)
	if err != nil {
		return nil, false, err
	}
	return &response.Group, response.AlreadyInGroup, nil
}

func (api *Slack) LeaveGroup(groupId string) error {
	values := url.Values{
		"token":   {api.config.token},
		"channel": {groupId},
	}
	_, err := groupRequest("groups.leave", values, api.debug)
	if err != nil {
		return err
	}
	return nil
}

func (api *Slack) KickUserFromGroup(groupId, userId string) error {
	values := url.Values{
		"token":   {api.config.token},
		"channel": {groupId},
		"user":    {userId},
	}
	_, err := groupRequest("groups.kick", values, api.debug)
	if err != nil {
		return err
	}
	return nil
}

func (api *Slack) GetGroups(excludeArchived bool) ([]Group, error) {
	values := url.Values{
		"token": {api.config.token},
	}
	if excludeArchived {
		values.Add("exclude_archived", "1")
	}
	response, err := groupRequest("groups.list", values, api.debug)
	if err != nil {
		return nil, err
	}
	return response.Groups, nil
}

/* Clients should try to avoid making this call too often. When needing to mark a read position, a client should set a
timer before making the call. In this way, any further updates needed during the timeout will not generate extra calls
(just one per channel). This is useful for when reading scroll-back history, or following a busy live channel. A timeout
of 5 seconds is a good starting point. Be sure to flush these calls on shutdown/logout.  */
func (api *Slack) SetGroupReadMark(groupId, ts string) error {
	values := url.Values{
		"token":   {api.config.token},
		"channel": {groupId},
		"ts":      {ts},
	}
	_, err := groupRequest("groups.mark", values, api.debug)
	if err != nil {
		return err
	}
	return nil
}

func (api *Slack) OpenGroup(groupId string) (bool, bool, error) {
	values := url.Values{
		"token": {api.config.token},
		"user":  {groupId},
	}
	response, err := groupRequest("groups.open", values, api.debug)
	if err != nil {
		return false, false, err
	}
	return response.NoOp, response.AlreadyOpen, nil
}

// XXX: They return a channel, not a group. What is this crap? :(
// Inconsistent api it seems.
func (api *Slack) RenameGroup(groupId, name string) (*Channel, error) {
	values := url.Values{
		"token":   {api.config.token},
		"channel": {groupId},
		"name":    {name},
	}
	// XXX: the created entry in this call returns a string instead of a number
	// so I may have to do some workaround to solve it.
	response, err := groupRequest("groups.rename", values, api.debug)
	if err != nil {
		return nil, err
	}
	return &response.Channel, nil

}

func (api *Slack) SetGroupPurpose(groupId, purpose string) (string, error) {
	values := url.Values{
		"token":   {api.config.token},
		"channel": {groupId},
		"purpose": {purpose},
	}
	response, err := groupRequest("groups.setPurpose", values, api.debug)
	if err != nil {
		return "", err
	}
	return response.Purpose, nil
}

func (api *Slack) SetGroupTopic(groupId, topic string) (string, error) {
	values := url.Values{
		"token":   {api.config.token},
		"channel": {groupId},
		"topic":   {topic},
	}
	response, err := groupRequest("groups.setTopic", values, api.debug)
	if err != nil {
		return "", err
	}
	return response.Topic, nil
}
