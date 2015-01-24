package slack

import (
	"errors"
	"net/url"
)

type emojiResponseFull struct {
	Emoji map[string]string `json:"emoji"`
	SlackResponse
}

func (api *Slack) GetEmoji() (map[string]string, error) {
	values := url.Values{
		"token": {api.config.token},
	}
	response := &emojiResponseFull{}
	err := ParseResponse("emoji.list", values, response, api.debug)
	if err != nil {
		return nil, err
	}
	if !response.Ok {
		return nil, errors.New(response.Error)
	}
	return response.Emoji, nil
}
