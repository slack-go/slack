package slack

import (
	"errors"
	"net/url"
)

// Reaction is the act of reacting to an item.
type Reaction struct {
	Name string `json:"name"`
	ItemRef
}

// ItemReaction is the reactions that have happened on an item.
type ItemReaction struct {
	Name  string   `json:"name"`
	Count int      `json:"count"`
	Users []string `json:"users"`
}

type getReactionsResponseFull struct {
	Message struct {
		Message struct {
			Reactions []ItemReaction
		}
	}
	SlackResponse
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

func getReactionRequest(path string, values url.Values, debug bool) (*getReactionsResponseFull, error) {
	response := &getReactionsResponseFull{}
	err := parseResponse(path, values, response, debug)
	if err != nil {
		return nil, err
	}
	if !response.Ok {
		return nil, errors.New(response.Error)
	}
	return response, nil
}

func setupReactionItemRef(values url.Values, item ItemRef) {
	if item.FileId != "" {
		values.Set("file", string(item.FileId))
	}
	if item.FileCommentId != "" {
		values.Set("file_comment", string(item.FileCommentId))
	}
	if item.ChannelId != "" {
		values.Set("channel", string(item.ChannelId))
	}
	if item.Timestamp != "" {
		values.Set("timestamp", string(item.Timestamp))
	}
}

// AddReaction adds a reaction emoji to a message, file or file comment.
func (api *Slack) AddReaction(name string, item ItemRef) error {
	values := url.Values{
		"token": {api.config.token},
	}
	if name != "" {
		values.Set("name", name)
	}
	setupReactionItemRef(values, item)
	_, err := addReactionRequest("reactions.add", values, api.debug)
	if err != nil {
		return err
	}
	return nil
}

// GetReactions returns details about the reactions on an item.
func (api *Slack) GetReactions(item ItemRef) ([]ItemReaction, error) {
	values := url.Values{
		"token": {api.config.token},
	}
	setupReactionItemRef(values, item)
	response, err := getReactionRequest("reactions.get", values, api.debug)
	if err != nil {
		return nil, err
	}
	return response.Message.Message.Reactions, nil
}
