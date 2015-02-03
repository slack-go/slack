package slack

type GroupCreatedEvent struct {
	Type    string             `json:"type"`
	UserId  string             `json:"user"`
	Channel ChannelCreatedInfo `json:"channel"`
}

// XXX: Should we really do this? event.Group is probably nicer than event.Channel
// even though the api returns "channel"
type GroupMarkedEvent ChannelInfoEvent
type GroupOpenEvent ChannelInfoEvent
type GroupCloseEvent ChannelInfoEvent
type GroupArchiveEvent ChannelInfoEvent
type GroupUnarchiveEvent ChannelInfoEvent
type GroupLeftEvent ChannelInfoEvent
type GroupJoinedEvent ChannelJoinedEvent
type GroupRenameEvent ChannelRenameEvent
type GroupHistoryChangedEvent ChannelHistoryChangedEvent
