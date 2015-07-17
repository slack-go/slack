package slack

import (
	"errors"
	"fmt"
	"net/url"
)

// ItemReaction is the reactions that have happened on an item.
type ItemReaction struct {
	Name  string   `json:"name"`
	Count int      `json:"count"`
	Users []string `json:"users"`
}

// ReactedItem is an item that was reacted to, and the details of the
// reactions.
type ReactedItem struct {
	Type      string
	Message   *Message
	File      *File
	Comment   *Comment
	Reactions []ItemReaction
}

// GetReactionsParameters is the inputs to get reactions to an item.
type GetReactionsParameters struct {
	Full bool
}

// NewGetReactionsParameters initializes the inputs to get reactions to an item.
func NewGetReactionsParameters() GetReactionsParameters {
	return GetReactionsParameters{
		Full: false,
	}
}

type getReactionsResponseFull struct {
	Type string
	M    struct {
		Reactions []ItemReaction
	} `json:"message"`
	F struct {
		Reactions []ItemReaction
	} `json:"file"`
	FC struct {
		Reactions []ItemReaction
	} `json:"comment"`
	SlackResponse
}

func (res getReactionsResponseFull) extractReactions() []ItemReaction {
	switch res.Type {
	case "message":
		return res.M.Reactions
	case "file":
		return res.F.Reactions
	case "file_comment":
		return res.FC.Reactions
	}
	return []ItemReaction{}
}

const (
	DEFAULT_REACTIONS_USERID = ""
	DEFAULT_REACTIONS_COUNT  = 100
	DEFAULT_REACTIONS_PAGE   = 1
	DEFAULT_REACTIONS_FULL   = false
)

// ListReactionsParameters is the inputs to find all reactions by a user.
type ListReactionsParameters struct {
	UserId string
	Count  int
	Page   int
	Full   bool
}

// NewListReactionsParameters initializes the inputs to find all reactions
// performed by a user.
func NewListReactionsParameters() ListReactionsParameters {
	return ListReactionsParameters{
		UserId: DEFAULT_REACTIONS_USERID,
		Count:  DEFAULT_REACTIONS_COUNT,
		Page:   DEFAULT_REACTIONS_PAGE,
		Full:   DEFAULT_REACTIONS_FULL,
	}
}

type listReactionsResponseFull struct {
	Items []struct {
		Type string
		M    struct {
			*Message
			Reactions []ItemReaction
		} `json:"message"`
		F struct {
			*File
			Reactions []ItemReaction
		} `json:"file"`
		FC struct {
			*Comment
			Reactions []ItemReaction
		} `json:"comment"`
	}
	Paging `json:"paging"`
	SlackResponse
}

func (res listReactionsResponseFull) extractReactedItems() []ReactedItem {
	items := make([]ReactedItem, len(res.Items))
	for i, input := range res.Items {
		item := ReactedItem{
			Type: input.Type,
		}
		switch input.Type {
		case "message":
			item.Message = input.M.Message
			item.Reactions = input.M.Reactions
		case "file":
			item.File = input.F.File
			item.Reactions = input.F.Reactions
		case "file_comment":
			item.Comment = input.FC.Comment
			item.Reactions = input.FC.Reactions
		}
		items[i] = item
	}
	return items
}

// AddReaction adds a reaction emoji to a message, file or file comment.
func (api *Slack) AddReaction(name string, item ItemRef) error {
	values := url.Values{
		"token": {api.config.token},
	}
	if name != "" {
		values.Set("name", name)
	}
	if item.ChannelId != "" {
		values.Set("channel", string(item.ChannelId))
	}
	if item.Timestamp != "" {
		values.Set("timestamp", string(item.Timestamp))
	}
	if item.FileId != "" {
		values.Set("file", string(item.FileId))
	}
	if item.CommentId != "" {
		values.Set("file_comment", string(item.CommentId))
	}
	response := &SlackResponse{}
	if err := parseResponse("reactions.add", values, response, api.debug); err != nil {
		return err
	}
	if !response.Ok {
		return errors.New(response.Error)
	}
	return nil
}

// RemoveReaction removes a reaction emoji from a message, file or file comment.
func (api *Slack) RemoveReaction(name string, item ItemRef) error {
	values := url.Values{
		"token": {api.config.token},
	}
	if name != "" {
		values.Set("name", name)
	}
	if item.ChannelId != "" {
		values.Set("channel", string(item.ChannelId))
	}
	if item.Timestamp != "" {
		values.Set("timestamp", string(item.Timestamp))
	}
	if item.FileId != "" {
		values.Set("file", string(item.FileId))
	}
	if item.CommentId != "" {
		values.Set("file_comment", string(item.CommentId))
	}
	response := &SlackResponse{}
	if err := parseResponse("reactions.remove", values, response, api.debug); err != nil {
		return err
	}
	if !response.Ok {
		return errors.New(response.Error)
	}
	return nil
}

// GetReactions returns details about the reactions on an item.
func (api *Slack) GetReactions(item ItemRef, params GetReactionsParameters) ([]ItemReaction, error) {
	values := url.Values{
		"token": {api.config.token},
	}
	if item.ChannelId != "" {
		values.Set("channel", string(item.ChannelId))
	}
	if item.Timestamp != "" {
		values.Set("timestamp", string(item.Timestamp))
	}
	if item.FileId != "" {
		values.Set("file", string(item.FileId))
	}
	if item.CommentId != "" {
		values.Set("file_comment", string(item.CommentId))
	}
	if params.Full != DEFAULT_REACTIONS_FULL {
		values.Set("full", fmt.Sprintf("%t", params.Full))
	}
	response := &getReactionsResponseFull{}
	if err := parseResponse("reactions.get", values, response, api.debug); err != nil {
		return nil, err
	}
	if !response.Ok {
		return nil, errors.New(response.Error)
	}
	return response.extractReactions(), nil
}

// ListReactions returns information about the items a user reacted to.
func (api *Slack) ListReactions(params ListReactionsParameters) ([]ReactedItem, *Paging, error) {
	values := url.Values{
		"token": {api.config.token},
	}
	if params.UserId != DEFAULT_REACTIONS_USERID {
		values.Add("user", params.UserId)
	}
	if params.Count != DEFAULT_REACTIONS_COUNT {
		values.Add("count", fmt.Sprintf("%d", params.Count))
	}
	if params.Page != DEFAULT_REACTIONS_PAGE {
		values.Add("page", fmt.Sprintf("%d", params.Page))
	}
	if params.Full != DEFAULT_REACTIONS_FULL {
		values.Add("full", fmt.Sprintf("%t", params.Full))
	}
	response := &listReactionsResponseFull{}
	err := parseResponse("reactions.list", values, response, api.debug)
	if err != nil {
		return nil, nil, err
	}
	if !response.Ok {
		return nil, nil, errors.New(response.Error)
	}
	return response.extractReactedItems(), &response.Paging, nil
}
