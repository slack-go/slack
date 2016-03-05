package slack

type MPIMCloseEvent struct {
	Type    string `json:"type"`
	User    string `json:"user"`
	Channel string `json:"channel"`
	mpimChannelEvent
}

type MPIMHistoryChangedEvent struct {
	ChannelHistoryChangedEvent
	IsMPIM int `json:"is_mpim"`
}

type mpimChannelEvent struct {
	IsMPIM bool `json:"is_mpim,omitempty"`
	IsOpen bool `json:"is_open,omitempty"`
}

type MPIMJoinedEvent struct {
	Type    string           `json:"type"`
	Channel mpimChannelEvent `json:"channel"`
}

type MPIMMarkedEvent struct {
	Type                string `json:"type"`
	Channel             string `json:"channel"`
	Timestamp           string `json:"ts,omitempty"`
	UnreadCount         int    `json:"unread_count,omitempty"`
	UnreadCountDisplay  int    `json:"unread_count_display,omitempty"`
	NumMentions         int    `json:"num_mentions,omitempty"`
	NumMentionsDisplay  int    `json:"num_mentions_display,omitempty"`
	MentionCount        int    `json:"mention_count,omitempty"`
	MentionCountDisplay int    `json:"mention_count_display,omitempty"`
	mpimChannelEvent
}

type MPIMOpenEvent MPIMCloseEvent
