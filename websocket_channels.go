package slack

type ChannelCreatedEvent struct {
	Type           string             `json:"type"`
	Channel        ChannelCreatedInfo `json:"channel"`
	EventTimestamp JSONTimeString     `json:"event_ts"`
}

type ChannelCreatedInfo struct {
	Id        string `json:"id"`
	IsChannel bool   `json:"is_channel"`
	Name      string `json:"name"`
	Created   int    `json:"created"`
	Creator   string `json:"creator"`
}

type ChannelJoinedEvent struct {
	Type    string  `json:"type"`
	Channel Channel `json:"channel"`
}

type ChannelInfoEvent struct {
	// channel_left
	// channel_deleted
	// channel_archive
	// channel_unarchive
	Type      string          `json:"type"`
	Channel   string          `json:"channel"`
	User      string          `json:"user,omitempty"`
	Timestamp *JSONTimeString `json:"ts,omitempty"`
}

type ChannelRenameEvent struct {
	Type    string            `json:"type"`
	Channel ChannelRenameInfo `json:"channel"`
}

type ChannelRenameInfo struct {
	ID      string         `json:"id"`
	Name    string         `json:"name"`
	Created JSONTimeString `json:"created"`
}

type ChannelHistoryChangedEvent struct {
	Type           string         `json:"type"`
	Latest         JSONTimeString `json:"latest"`
	Timestamp      JSONTimeString `json:"ts"`
	EventTimestamp JSONTimeString `json:"event_ts"`
}

type ChannelMarkedEvent ChannelInfoEvent
type ChannelLeftEvent ChannelInfoEvent
type ChannelDeletedEvent ChannelInfoEvent
type ChannelArchiveEvent ChannelInfoEvent
type ChannelUnarchiveEvent ChannelInfoEvent
