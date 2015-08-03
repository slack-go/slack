package slack

import (
	"encoding/json"
	"time"
)

// TODO: Probably need an error event

// HelloEvent represents the hello event
type HelloEvent struct{}

// PresenceChangeEvent represents the presence change event
type PresenceChangeEvent struct {
	Type     string `json:"type"`
	Presence string `json:"presence"`
	User     string `json:"user"`
}

// UserTypingEvent represents the user typing event
type UserTypingEvent struct {
	Type    string `json:"type"`
	User    string `json:"user"`
	Channel string `json:"channel"`
}

// LatencyReport represents the latency report
type LatencyReport struct {
	Value time.Duration
}

// PrefChangeEvent represents the preference change event
type PrefChangeEvent struct {
	Type  string          `json:"type"`
	Name  string          `json:"name"`
	Value json.RawMessage `json:"value"`
}

// ManualPresenceChangeEvent represents the manual presence change event
type ManualPresenceChangeEvent struct {
	Type     string `json:"type"`
	Presence string `json:"presence"`
}

// UserChangeEvent represents the user change event
type UserChangeEvent struct {
	Type string `json:"type"`
	User User   `json:"user"`
}

// EmojiChangedEvent represents the emoji changed event
type EmojiChangedEvent struct {
	Type           string         `json:"type"`
	EventTimestamp JSONTimeString `json:"event_ts"`
}

// CommandsChangedEvent represents the commands changed event
type CommandsChangedEvent struct {
	Type           string         `json:"type"`
	EventTimestamp JSONTimeString `json:"event_ts"`
}

// EmailDomainChangedEvent represents the email domain changed event
type EmailDomainChangedEvent struct {
	Type           string         `json:"type"`
	EventTimestamp JSONTimeString `json:"event_ts"`
	EmailDomain    string         `json:"email_domain"`
}

// BotAddedEvent represents the bot added event
type BotAddedEvent struct {
	Type string `json:"type"`
	Bot  Bot    `json:"bot"`
}

// BotChangedEvent represents the bot changed event
type BotChangedEvent struct {
	Type string `json:"type"`
	Bot  Bot    `json:"bot"`
}

// AccountsChangedEvent represents the accounts changed event
type AccountsChangedEvent struct {
	Type string `json:"type"`
}
