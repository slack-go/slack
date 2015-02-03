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

// TODO
type ManualPresenceChangeEvent struct{}
type UserChangeEvent struct{}
type EmojiChangedEvent struct{}
type CommandsChangedEvent struct{}
type EmailDomainChangedEvent struct{}
type BotAddedEvent struct{}
type BotChangedEvent struct{}
type AccountsChangedEvent struct{}
