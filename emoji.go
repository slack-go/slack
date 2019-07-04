package slack

import (
	"context"
	"fmt"
	"io"
	"net/url"
)

type emojiAddOk struct {
	Ok bool `json:"ok"`
}

type emojiResponseFull struct {
	Emoji map[string]string `json:"emoji"`
	SlackResponse
}

// EmojiAddParameters contains all the parameters necessary (including the optional ones) for an AddEmoji() request.
type EmojiAddParameters struct {
	File     string
	Content  string
	Reader   io.Reader
	Filename string
}

// GetEmoji retrieves all the emojis
func (api *Client) GetEmoji() (map[string]string, error) {
	return api.GetEmojiContext(context.Background())
}

// GetEmojiContext retrieves all the emojis with a custom context
func (api *Client) GetEmojiContext(ctx context.Context) (map[string]string, error) {
	values := url.Values{
		"token": {api.token},
	}
	response := &emojiResponseFull{}

	err := api.postMethod(ctx, "emoji.list", values, response)
	if err != nil {
		return nil, err
	}

	if response.Err() != nil {
		return nil, response.Err()
	}

	return response.Emoji, nil
}

// AddEmoji uploads a emoji
func (api *Client) AddEmoji(name string, params EmojiAddParameters) (emoji *emojiAddOk, err error) {
	return api.AddEmojiContext(context.Background(), name, params)
}

// AddEmojiContext uploads a emoji and setting a custom context
func (api *Client) AddEmojiContext(ctx context.Context, name string, params EmojiAddParameters) (emoji *emojiAddOk, err error) {
	// Test if user token is valid. This helps because client.Do doesn't like this for some reason. XXX: More
	// investigation needed, but for now this will do.
	_, err = api.AuthTest()
	if err != nil {
		return nil, err
	}
	if name == "" {
		return nil, fmt.Errorf("emoji.add: Emoji name is mandatory")
	}
	response := &emojiAddOk{}
	values := url.Values{
		"token": {api.token},
	}
	values.Add("mode", "data")
	if name != "" {
		values.Add("name", name)
	}
	if params.Content != "" {
		values.Add("image", params.Content)
		err = api.postMethod(ctx, "emoji.add", values, response)
	} else {
		workspace, err := api.AuthTest()
		if err != nil {
			return nil, err
		}
		endpoint := fmt.Sprintf(api.webendpointformat, workspace.URL)

		if params.File != "" {
			err = postLocalWithMultipartResponse(ctx, api.httpclient, endpoint+"emoji.add", params.File, "image", values, response, api)
		} else if params.Reader != nil {
			err = postWithMultipartResponse(ctx, api.httpclient, endpoint+"emoji.add", params.Filename, "image", values, params.Reader, response, api)
		}
	}
	if err != nil {
		return nil, err
	}

	return response, err
}
