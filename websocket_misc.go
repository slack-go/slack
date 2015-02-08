package slack

import (
	"encoding/json"
	"time"
)

// TODO: Probably need an error event

type HelloEvent struct{}

type PresenceChangeEvent struct {
	Type     string `json:"type"`
	Presence string `json:"presence"`
	UserId   string `json:"user"`
}

type UserTypingEvent struct {
	Type      string `json:"type"`
	UserId    string `json:"user"`
	ChannelId string `json:"channel"`
}

type LatencyReport struct {
	Value time.Duration
}

type PrefChangeEvent struct {
	Type  string          `json:"type"`
	Name  string          `json:"name"`
	Value json.RawMessage `json:"value"`
}

type ManualPresenceChangeEvent struct {
	Type     string `json:"type"`
	Presence string `json:"presence"`
}
type UserChangeEvent struct {
	Type string `json:"type"`
	User User   `json:"user"`
}
type EmojiChangedEvent struct {
	Type           string         `json:"type"`
	EventTimestamp JSONTimeString `json:"event_ts"`
}
type CommandsChangedEvent struct {
	Type           string         `json:"type"`
	EventTimestamp JSONTimeString `json:"event_ts"`
}
type EmailDomainChangedEvent struct {
	Type           string         `json:"type"`
	EventTimestamp JSONTimeString `json:"event_ts"`
	EmailDomain    string         `json:"email_domain"`
}
type BotAddedEvent struct {
	Type string `json:"type"`
	Bot  Bot    `json:"bot"`
}
type BotChangedEvent struct {
	Type string `json:"type"`
	Bot  Bot    `json:"bot"`
}
type AccountsChangedEvent struct {
	Type string `json:"type"`
}
