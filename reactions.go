package slack

import (
	"errors"
	"fmt"
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

// ReactedToItem is an item that was reacted to, and the details of the
// reaction.
type ReactedToItem struct {
	Item
	ItemReaction
}

// AddReactionParameters is the inputs to create a new reaction.
type AddReactionParameters struct {
	ItemRef
	Name string
}

// NewAddReactionParameters initialies the inputs to react to an item.
func NewAddReactionParameters(name string, ref ItemRef) AddReactionParameters {
	return AddReactionParameters{Name: name, ItemRef: ref}
}

// GetReactionParameters is the inputs to get reactions on an item.
type GetReactionParameters struct {
	ItemRef
	Full bool
}

// NewGetReactionParameters initializes the inputs to get reactions on an item.
func NewGetReactionParameters(ref ItemRef) GetReactionParameters {
	return GetReactionParameters{ItemRef: ref}
}

type getReactionsResponseFull struct {
	Message struct {
		Type    string
		Message struct {
			Reactions []ItemReaction
		}
		File struct {
			Reactions []ItemReaction
		}
		FileComment struct {
			Comment struct {
				Reactions []ItemReaction
			}
		} `json:"file_comment"`
	}
	SlackResponse
}

func (res getReactionsResponseFull) FindReactions() []ItemReaction {
	switch res.Message.Type {
	case "message":
		return res.Message.Message.Reactions
	case "file":
		return res.Message.File.Reactions
	case "file_comment":
		return res.Message.FileComment.Comment.Reactions
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
	User  string
	Count int
	Page  int
	Full  bool
}

// NewListReactionsParameters initializes the inputs to find all reactions by a user.
func NewListReactionsParameters(userID string) ListReactionsParameters {
	return ListReactionsParameters{
		User:  userID,
		Count: DEFAULT_REACTIONS_COUNT,
		Page:  DEFAULT_REACTIONS_PAGE,
		Full:  DEFAULT_REACTIONS_FULL,
	}
}

type listReactionsResponseFull struct {
	Items []struct {
		Message struct {
			Reactions []ItemReaction
		}
	}
	SlackResponse
}

// AddReaction adds a reaction emoji to a message, file or file comment.
func (api *Slack) AddReaction(params AddReactionParameters) error {
	values := url.Values{
		"token": {api.config.token},
	}
	if params.Name != "" {
		values.Set("name", params.Name)
	}
	if params.ChannelId != "" {
		values.Set("channel", string(params.ChannelId))
	}
	if params.Timestamp != "" {
		values.Set("timestamp", string(params.Timestamp))
	}
	if params.FileId != "" {
		values.Set("file", string(params.FileId))
	}
	if params.FileCommentId != "" {
		values.Set("file_comment", string(params.FileCommentId))
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

// GetReactions returns details about the reactions on an item.
func (api *Slack) GetReactions(params GetReactionParameters) ([]ItemReaction, error) {
	values := url.Values{
		"token": {api.config.token},
	}
	if params.ChannelId != "" {
		values.Set("channel", string(params.ChannelId))
	}
	if params.Timestamp != "" {
		values.Set("timestamp", string(params.Timestamp))
	}
	if params.FileId != "" {
		values.Set("file", string(params.FileId))
	}
	if params.FileCommentId != "" {
		values.Set("file_comment", string(params.FileCommentId))
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
	return response.FindReactions(), nil
}

// ListReactions returns information about the items a user reacted to.
func (api *Slack) ListReactions(params ListReactionsParameters) ([]ReactedToItem, Paging, error) {
	values := url.Values{
		"token": {api.config.token},
	}
	if params.User != DEFAULT_REACTIONS_USERID {
		values.Add("user", params.User)
	}
	if params.Count != DEFAULT_REACTIONS_COUNT {
		values.Add("count", string(params.Count))
	}
	if params.Page != DEFAULT_REACTIONS_PAGE {
		values.Add("count", string(params.Page))
	}
	if params.Full != DEFAULT_REACTIONS_FULL {
		values.Add("count", fmt.Sprintf("%t", params.Full))
	}
	response := &listReactionsResponseFull{}
	err := parseResponse("reactions.list", values, response, api.debug)
	if err != nil {
		return nil, Paging{}, err
	}
	if !response.Ok {
		return nil, Paging{}, errors.New(response.Error)
	}
	return nil, Paging{}, nil
}
