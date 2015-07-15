package slack

import (
	"errors"
	"net/url"
)

// Reaction describes the reaction and the item reacted to. One of file,
// file_comment, or the combination of channel and timestamp must be specified.
type Reaction struct {
	Name        string `json:"name"`
	File        string `json:"file"`
	FileComment string `json:"file_comment"`
	Channel     string `json:"channel"`
	Timestamp   string `json:"timestamp"`
}

// NewMessageReaction initializes a reaction to a message.
func NewMessageReaction(name, channel, timestamp string) Reaction {
	return Reaction{Channel: channel, Timestamp: timestamp}
}

// NewFileReaction initializes a reaction to a file.
func NewFileReaction(name, file string) Reaction {
	return Reaction{Name: name, File: file}
}

// NewFileCommentReaction initializes a reaction to a file comment.
func NewFileCommentReaction(name, fileComment string) Reaction {
	return Reaction{Name: name, FileComment: fileComment}
}

func addReactionRequest(path string, values url.Values, debug bool) (*SlackResponse, error) {
	response := &SlackResponse{}
	err := parseResponse(path, values, response, debug)
	if err != nil {
		return nil, err
	}
	if !response.Ok {
		return nil, errors.New(response.Error)
	}
	return response, nil
}

// AddReaction adds a reaction emoji to a message, file or file comment.
func (api *Slack) AddReaction(reaction Reaction) error {
	values := url.Values{
		"token": {api.config.token},
	}
	if reaction.Name != "" {
		values.Set("name", reaction.Name)
	}
	if reaction.File != "" {
		values.Set("file", string(reaction.File))
	}
	if reaction.FileComment != "" {
		values.Set("file_comment", string(reaction.FileComment))
	}
	if reaction.Channel != "" {
		values.Set("channel", string(reaction.Channel))
	}
	if reaction.Timestamp != "" {
		values.Set("timestamp", string(reaction.Timestamp))
	}
	_, err := addReactionRequest("reactions.add", values, api.debug)
	if err != nil {
		return err
	}
	return nil
}
