package slack

import (
	"errors"
	"net/url"
)

const (
	DEFAULT_REACTION_NAME         = ""
	DEFAULT_REACTION_FILE         = ""
	DEFAULT_REACTION_FILE_COMMENT = ""
	DEFAULT_REACTION_CHANNEL      = ""
	DEFAULT_REACTION_TIMESTAMP    = ""
)

type reactionResponseFull struct {
	SlackResponse
}

type ReactionParameters struct {
	Name        string `json:"name"`
	File        string `json:"file"`
	FileComment string `json:"file_comment"`
	Channel     string `json:"channel"`
	Timestamp   string `json:"timestamp"`
}

// NewReactioneParameters provides an instance of
// ReactionParameters with all of the sane default value set.
func NewReactionParameters() ReactionParameters {
	return ReactionParameters{
		Name:        DEFAULT_REACTION_NAME,
		File:        DEFAULT_REACTION_FILE,
		FileComment: DEFAULT_REACTION_FILE_COMMENT,
		Channel:     DEFAULT_REACTION_CHANNEL,
		Timestamp:   DEFAULT_REACTION_TIMESTAMP,
	}
}

func reactionRequest(path string, values url.Values, debug bool) (*reactionResponseFull, error) {
	response := &reactionResponseFull{}
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
// One of file, file_comment, or the combination of channel and timestamp
// must be specified.
func (api *Slack) AddReaction(params ReactionParameters) error {
	values := url.Values{
		"token": {api.config.token},
	}
	if params.Name != DEFAULT_REACTION_NAME {
		values.Set("name", string(params.Name))
	}
	if params.File != DEFAULT_REACTION_FILE {
		values.Set("file", string(params.File))
	}
	if params.FileComment != DEFAULT_REACTION_FILE_COMMENT {
		values.Set("file_comment", string(params.FileComment))
	}
	if params.Channel != DEFAULT_REACTION_CHANNEL {
		values.Set("channel", string(params.Channel))
	}
	if params.Timestamp != DEFAULT_REACTION_TIMESTAMP {
		values.Set("timestamp", string(params.Timestamp))
	}
	_, err := reactionRequest("reactions.add", values, api.debug)
	if err != nil {
		return err
	}
	return nil
}
