package slack

import (
	"context"
	"errors"
	"net/url"
)

// Conversation is the foundation for IM and BaseGroupConversation
type conversation struct {
	ID                 string   `json:"id"`
	Created            JSONTime `json:"created"`
	IsOpen             bool     `json:"is_open"`
	LastRead           string   `json:"last_read,omitempty"`
	Latest             *Message `json:"latest,omitempty"`
	UnreadCount        int      `json:"unread_count,omitempty"`
	UnreadCountDisplay int      `json:"unread_count_display,omitempty"`
}

// GroupConversation is the foundation for Group and Channel
type groupConversation struct {
	conversation
	Name       string   `json:"name"`
	Creator    string   `json:"creator"`
	IsArchived bool     `json:"is_archived"`
	Members    []string `json:"members"`
	Topic      Topic    `json:"topic"`
	Purpose    Purpose  `json:"purpose"`
}

// Topic contains information about the topic
type Topic struct {
	Value   string   `json:"value"`
	Creator string   `json:"creator"`
	LastSet JSONTime `json:"last_set"`
}

// Purpose contains information about the purpose
type Purpose struct {
	Value   string   `json:"value"`
	Creator string   `json:"creator"`
	LastSet JSONTime `json:"last_set"`
}

type GetUsersInConversationParameters struct {
	ChannelID string
	Cursor    string
	Limit     int
}

type responseMetaData struct {
	NextCursor string `json:"next_cursor"`
}

// GetUsersInConversation returns the list of users in a conversation
func (api *Client) GetUsersInConversation(params *GetUsersInConversationParameters) ([]string, string, error) {
	return api.GetUsersInConversationContext(context.Background(), params)
}

// GetUsersInConversation returns the list of users in a conversation with a custom context
func (api *Client) GetUsersInConversationContext(ctx context.Context, params *GetUsersInConversationParameters) ([]string, string, error) {
	values := url.Values{
		"token":   {api.config.token},
		"channel": {params.ChannelID},
	}
	if params.Cursor != "" {
		values.Add("cursor", params.Cursor)
	}
	if params.Limit != 0 {
		values.Add("limit", string(params.Limit))
	}
	response := struct {
		Members          []string         `json:"members"`
		ResponseMetaData responseMetaData `json:"response_metadata"`
		SlackResponse
	}{}
	err := post(ctx, "conversations.members", values, &response, api.debug)
	if err != nil {
		return nil, "", err
	}
	if !response.Ok {
		return nil, "", errors.New(response.Error)
	}
	return response.Members, response.ResponseMetaData.NextCursor, nil
}

// ArchiveConversation archives a conversation
func (api *Client) ArchiveConversation(channelID string) error {
	return api.ArchiveConversationContext(context.Background(), channelID)
}

// ArchiveConversationContext archives a conversation with a custom context
func (api *Client) ArchiveConversationContext(ctx context.Context, channelID string) error {
	values := url.Values{
		"token":   {api.config.token},
		"channel": {channelID},
	}
	response := SlackResponse{}
	err := post(ctx, "conversations.archive", values, &response, api.debug)
	if err != nil {
		return err
	}
	if !response.Ok {
		return errors.New(response.Error)
	}
	return nil
}

// UnArchiveConversation reverses conversation archival
func (api *Client) UnArchiveConversation(channelID string) error {
	return api.UnArchiveConversationContext(context.Background(), channelID)
}

// UnArchiveConversationContext reverses conversation archival with a custom context
func (api *Client) UnArchiveConversationContext(ctx context.Context, channelID string) error {
	values := url.Values{
		"token":   {api.config.token},
		"channel": {channelID},
	}
	response := SlackResponse{}
	err := post(ctx, "conversations.unarchive", values, &response, api.debug)
	if err != nil {
		return err
	}
	if !response.Ok {
		return errors.New(response.Error)
	}
	return nil
}

// SetTopicOfConversation sets the topic for a conversation
func (api *Client) SetTopicOfConversation(channelID, topic string) (*Channel, error) {
	return api.SetTopicOfConversationContext(context.Background(), channelID, topic)
}

// SetTopicOfConversationContext sets the topic for a conversation with a custom context
func (api *Client) SetTopicOfConversationContext(ctx context.Context, channelID, topic string) (*Channel, error) {
	values := url.Values{
		"token":   {api.config.token},
		"channel": {channelID},
		"topic":   {topic},
	}
	response := struct {
		SlackResponse
		Channel *Channel `json:"channel"`
	}{}
	err := post(ctx, "conversations.setTopic", values, &response, api.debug)
	if err != nil {
		return nil, err
	}
	if !response.Ok {
		return nil, errors.New(response.Error)
	}
	return response.Channel, nil
}

// SetPurposeOfConversation sets the purpose for a conversation
func (api *Client) SetPurposeOfConversation(channelID, purpose string) (*Channel, error) {
	return api.SetPurposeOfConversationContext(context.Background(), channelID, purpose)
}

// SetPurposeOfConversationContext sets the purpose for a conversation with a custom context
func (api *Client) SetPurposeOfConversationContext(ctx context.Context, channelID, purpose string) (*Channel, error) {
	values := url.Values{
		"token":   {api.config.token},
		"channel": {channelID},
		"purpose": {purpose},
	}
	response := struct {
		SlackResponse
		Channel *Channel `json:"channel"`
	}{}
	err := post(ctx, "conversations.setPurpose", values, &response, api.debug)
	if err != nil {
		return nil, err
	}
	if !response.Ok {
		return nil, errors.New(response.Error)
	}
	return response.Channel, nil
}

// RenameConversation renames a conversation
func (api *Client) RenameConversation(channelID, channelName string) (*Channel, error) {
	return api.RenameConversationContext(context.Background(), channelID, channelName)
}

// RenameConversationContext renames a conversation with a custom context
func (api *Client) RenameConversationContext(ctx context.Context, channelID, channelName string) (*Channel, error) {
	values := url.Values{
		"token":   {api.config.token},
		"channel": {channelID},
		"name":    {channelName},
	}
	response := struct {
		SlackResponse
		Channel *Channel `json:"channel"`
	}{}
	err := post(ctx, "conversations.rename", values, &response, api.debug)
	if err != nil {
		return nil, err
	}
	if !response.Ok {
		return nil, errors.New(response.Error)
	}
	return response.Channel, nil
}
