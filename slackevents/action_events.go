package slackevents

import (
	"encoding/json"

	"github.com/slack-go/slack"
)

// Deprecated: MessageActionResponse is associated with [MessageAction] which cannot
// handle block_actions. Use [slack.InteractionCallback] instead.
type MessageActionResponse struct {
	ResponseType    string `json:"response_type"`
	ReplaceOriginal bool   `json:"replace_original"`
	Text            string `json:"text"`
}

// Deprecated: MessageActionEntity is associated with [MessageAction] which cannot
// handle block_actions. Use [slack.InteractionCallback] instead.
type MessageActionEntity struct {
	ID     string `json:"id"`
	Domain string `json:"domain"`
	Name   string `json:"name"`
}

// Deprecated: MessageAction cannot represent block_actions payloads. Use
// [slack.InteractionCallback] instead, which handles all interaction types.
// See [slack.InteractionCallbackParse] for parsing from an HTTP request.
type MessageAction struct {
	Type             string                   `json:"type"`
	Actions          []slack.AttachmentAction `json:"actions"`
	CallbackID       string                   `json:"callback_id"`
	Team             MessageActionEntity      `json:"team"`
	Channel          MessageActionEntity      `json:"channel"`
	User             MessageActionEntity      `json:"user"`
	ActionTimestamp  json.Number              `json:"action_ts"`
	MessageTimestamp json.Number              `json:"message_ts"`
	AttachmentID     json.Number              `json:"attachment_id"`
	Token            string                   `json:"token"`
	Message          slack.Message            `json:"message"`
	OriginalMessage  slack.Message            `json:"original_message"`
	ResponseURL      string                   `json:"response_url"`
	TriggerID        string                   `json:"trigger_id"`
}
