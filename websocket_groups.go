package slack

// GroupCreatedEvent represents the Group created event
type GroupCreatedEvent struct {
	Type    string             `json:"type"`
	User    string             `json:"user"`
	Channel ChannelCreatedInfo `json:"channel"`
}

// XXX: Should we really do this? event.Group is probably nicer than event.Channel
// even though the api returns "channel"

// GroupMarkedEvent represents the Group marked event
type GroupMarkedEvent ChannelInfoEvent

// GroupOpenEvent represents the Group open event
type GroupOpenEvent ChannelInfoEvent

// GroupCloseEvent represents the Group close event
type GroupCloseEvent ChannelInfoEvent

// GroupArchiveEvent represents the Group archive event
type GroupArchiveEvent ChannelInfoEvent

// GroupUnarchiveEvent represents the Group unarchive event
type GroupUnarchiveEvent ChannelInfoEvent

// GroupLeftEvent represents the Group left event
type GroupLeftEvent ChannelInfoEvent

// GroupJoinedEvent represents the Group joined event
type GroupJoinedEvent ChannelJoinedEvent

// GroupRenameEvent represents the Group rename event
type GroupRenameEvent ChannelRenameEvent

// GroupHistoryChangedEvent represents the Group history changed event
type GroupHistoryChangedEvent ChannelHistoryChangedEvent
