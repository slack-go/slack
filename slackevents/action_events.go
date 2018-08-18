package slackevents

import (
	"encoding/json"

	"github.com/nlopes/slack"
)

type MessageActionResponse struct {
	ResponseType    string `json:"response_type"`
	ReplaceOriginal bool   `json:"replace_original"`
	Text            string `json:"text"`
}

type MessageActionEntity struct {
	Id     string `json:"id"`
	Domain string `json:"domain"`
}

type MessageAction struct {
	Type             string                   `json:"type"`
	Actions          []slack.AttachmentAction `json:"actions"`
	CallbackId       string                   `json:"callback_id"`
	Team             MessageActionEntity      `json:"team"`
	Channel          MessageActionEntity      `json:"channel"`
	User             MessageActionEntity      `json:"user"`
	ActionTimestamp  json.Number              `json:"action_ts"`
	MessageTimestamp json.Number              `json:"message_ts"`
	AttachmentId     json.Number              `json:"attachment_id"`
	Token            string                   `json:"token"`
	Message          slack.Message            `json:"message"`
	OriginalMessage  slack.Message            `json:"original_message"`
	ResponseUrl      string                   `json:"response_url"`
	TriggerId        string                   `json:"trigger_id"`
}
