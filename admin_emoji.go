package slack

import (
	"context"
	"net/url"
)

type adminEmojiResponseFull struct {
	Emoji map[string]*Emoji `json:"emoji"`
	SlackResponse
}

type Emoji struct {
	URL         string `json:"url"`
	DateCreated int    `json:"date_created"`
	UploadedBy  string `json:"uploaded_by"`
}

// GetAdminEmoji retrieves all the emojis
func (api *Client) GetAdminEmoji() (map[string]*Emoji, error) {
	return api.GetAdminEmojiContext(context.Background())
}

// GetAdminEmojiContext retrieves all the emojis with a custom context
func (api *Client) GetAdminEmojiContext(ctx context.Context) (map[string]*Emoji, error) {
	values := url.Values{
		"token": {api.token},
	}
	response := &adminEmojiResponseFull{}

	err := api.postMethod(ctx, "admin.emoji.list", values, response)
	if err != nil {
		return nil, err
	}

	if response.Err() != nil {
		return nil, response.Err()
	}

	return response.Emoji, nil
}
