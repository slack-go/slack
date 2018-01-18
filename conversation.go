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
