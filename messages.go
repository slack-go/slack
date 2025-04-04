package slack

// OutgoingMessage is used for the realtime API, and seems incomplete.
type OutgoingMessage struct {
	ID int `json:"id"`
	// channel ID
	Channel         string   `json:"channel,omitempty"`
	Text            string   `json:"text,omitempty"`
	Type            string   `json:"type,omitempty"`
	ThreadTimestamp string   `json:"thread_ts,omitempty"`
	ThreadBroadcast bool     `json:"reply_broadcast,omitempty"`
	IDs             []string `json:"ids,omitempty"`
}

// Message is an auxiliary type to allow us to have a message containing sub messages
type Message struct {
	Msg
	SubMessage      *Msg `json:"message,omitempty"`
	PreviousMessage *Msg `json:"previous_message,omitempty"`
	// Root is the message that was broadcast to the channel when the SubType is
	// thread_broadcast. If this is not a thread_broadcast message event, this
	// value is nil.
	Root *Msg `json:"root,omitempty"`
}

// Msg SubTypes (https://api.slack.com/events/message)
const (
	MsgSubTypeBotMessage                = "bot_message"                 // [Events API, RTM] A message was posted by an integration
	MsgSubTypeMeMessage                 = "me_message"                  // [Events API, RTM] A /me message was sent
	MsgSubTypeMessageChanged            = "message_changed"             // [Events API, RTM] A message was changed
	MsgSubTypeMessageDeleted            = "message_deleted"             // [Events API, RTM] A message was deleted
	MsgSubTypeMessageReplied            = "message_replied"             // [Events API, RTM] A message thread received a reply
	MsgSubTypeReplyBroadcast            = "reply_broadcast"             // @Deprecated (No longer served) A message thread's reply was broadcast to a channel
	MsgSubTypeThreadBroadcast           = "thread_broadcast"            // [Events API, RTM] A message thread's reply was broadcast to a channel
	MsgSubTypeChannelJoin               = "channel_join"                // [Events API, RTM] A member joined a channel
	MsgSubTypeChannelLeave              = "channel_leave"               // [Events API, RTM] A member left a channel
	MsgSubTypeChannelTopic              = "channel_topic"               // [Events API, RTM] A channel topic was updated
	MsgSubTypeChannelPurpose            = "channel_purpose"             // [Events API, RTM] A channel purpose was updated
	MsgSubTypeChannelName               = "channel_name"                // [Events API, RTM] A channel was renamed
	MsgSubTypeChannelArchive            = "channel_archive"             // [Events API, RTM] A channel was archived
	MsgSubTypeChannelUnarchive          = "channel_unarchive"           // [Events API, RTM] A channel was unarchived
	MsgSubTypeGroupJoin                 = "group_join"                  // [RTM] A member joined a group
	MsgSubTypeGroupLeave                = "group_leave"                 // [RTM] A member left a group
	MsgSubTypeGroupTopic                = "group_topic"                 // [RTM] A group topic was updated
	MsgSubTypeGroupPurpose              = "group_purpose"               // [RTM] A group purpose was updated
	MsgSubTypeGroupName                 = "group_name"                  // [RTM] A group was renamed
	MsgSubTypeGroupArchive              = "group_archive"               // [RTM] A group was archived
	MsgSubTypeGroupUnarchive            = "group_unarchive"             // [RTM] A group was unarchived
	MsgSubTypeFileShare                 = "file_share"                  // [Events API, RTM] A file was shared into a channel
	MsgSubTypeFileComment               = "file_comment"                // [RTM] A comment was added to a file
	MsgSubTypeGileMention               = "file_mention"                // [RTM] A file was mentioned in a channel
	MsgSubTypePinnedItem                = "pinned_item"                 // [RTM] An item was pinned in a channel
	MsgSubTypeUnpinnedItem              = "unpinned_item"               // [RTM] An item was unpinned from a channel
	MsgSubTypeEkmAccessDenied           = "ekm_access_denied"           // [Events API, RTM] Message content redacted due to Enterprise Key Management (EKM)
	MsgSubTypeChannelPostingPermissions = "channel_posting_permissions" // [Events API, RTM] The posting permissions for a channel changed
	MsgSubTypeAssistantAppThread        = "assistant_app_thread"        // [Events API, RTM] The message is an app assistant thread
)

// Msg contains information about a slack message
type Msg struct {
	// Basic Message
	ClientMsgID     string       `json:"client_msg_id,omitempty"`
	Type            string       `json:"type,omitempty"`
	Channel         string       `json:"channel,omitempty"`
	User            string       `json:"user,omitempty"`
	Text            string       `json:"text,omitempty"`
	Timestamp       string       `json:"ts,omitempty"`
	ThreadTimestamp string       `json:"thread_ts,omitempty"`
	IsStarred       bool         `json:"is_starred,omitempty"`
	PinnedTo        []string     `json:"pinned_to,omitempty"`
	Attachments     []Attachment `json:"attachments,omitempty"`
	Edited          *Edited      `json:"edited,omitempty"`
	LastRead        string       `json:"last_read,omitempty"`
	Subscribed      bool         `json:"subscribed,omitempty"`
	UnreadCount     int          `json:"unread_count,omitempty"`

	// Message Subtypes
	SubType string `json:"subtype,omitempty"`

	// Hidden Subtypes
	Hidden           bool   `json:"hidden,omitempty"`     // message_changed, message_deleted, unpinned_item
	DeletedTimestamp string `json:"deleted_ts,omitempty"` // message_deleted
	EventTimestamp   string `json:"event_ts,omitempty"`

	// bot_message (https://api.slack.com/events/message/bot_message)
	BotID      string      `json:"bot_id,omitempty"`
	Username   string      `json:"username,omitempty"`
	Icons      *Icon       `json:"icons,omitempty"`
	BotProfile *BotProfile `json:"bot_profile,omitempty"`

	// channel_join, group_join
	Inviter string `json:"inviter,omitempty"`

	// channel_topic, group_topic
	Topic string `json:"topic,omitempty"`

	// channel_purpose, group_purpose
	Purpose string `json:"purpose,omitempty"`

	// channel_name, group_name
	Name    string `json:"name,omitempty"`
	OldName string `json:"old_name,omitempty"`

	// channel_archive, group_archive
	Members []string `json:"members,omitempty"`

	// channels.replies, groups.replies, im.replies, mpim.replies
	ReplyCount   int      `json:"reply_count,omitempty"`
	ReplyUsers   []string `json:"reply_users,omitempty"`
	Replies      []Reply  `json:"replies,omitempty"`
	ParentUserId string   `json:"parent_user_id,omitempty"`
	LatestReply  string   `json:"latest_reply,omitempty"`

	// file_share, file_comment, file_mention
	Files []File `json:"files,omitempty"`

	// file_share
	Upload bool `json:"upload,omitempty"`

	// file_comment
	Comment *Comment `json:"comment,omitempty"`

	// pinned_item
	ItemType string `json:"item_type,omitempty"`

	// https://api.slack.com/rtm
	ReplyTo int    `json:"reply_to,omitempty"`
	Team    string `json:"team,omitempty"`

	// reactions
	Reactions []ItemReaction `json:"reactions,omitempty"`

	// slash commands and interactive messages
	ResponseType    string `json:"response_type,omitempty"`
	ReplaceOriginal bool   `json:"replace_original"`
	DeleteOriginal  bool   `json:"delete_original"`

	// metadata
	Metadata SlackMetadata `json:"metadata,omitempty"`

	// Block type Message
	Blocks Blocks `json:"blocks,omitempty"`
	// permalink
	Permalink string `json:"permalink,omitempty"`
}

const (
	// ResponseTypeInChannel in channel response for slash commands.
	ResponseTypeInChannel = "in_channel"
	// ResponseTypeEphemeral ephemeral response for slash commands.
	ResponseTypeEphemeral = "ephemeral"
)

// ScheduledMessage contains information about a slack scheduled message
type ScheduledMessage struct {
	ID          string `json:"id"`
	Channel     string `json:"channel_id"`
	PostAt      int    `json:"post_at"`
	DateCreated int    `json:"date_created"`
	Text        string `json:"text"`
}

// Icon is used for bot messages
type Icon struct {
	IconURL   string `json:"icon_url,omitempty"`
	IconEmoji string `json:"icon_emoji,omitempty"`
}

// BotProfile contains information about a bot
type BotProfile struct {
	AppID   string `json:"app_id,omitempty"`
	Deleted bool   `json:"deleted,omitempty"`
	Icons   *Icons `json:"icons,omitempty"`
	ID      string `json:"id,omitempty"`
	Name    string `json:"name,omitempty"`
	TeamID  string `json:"team_id,omitempty"`
	Updated int64  `json:"updated,omitempty"`
}

// Edited indicates that a message has been edited.
type Edited struct {
	User      string `json:"user,omitempty"`
	Timestamp string `json:"ts,omitempty"`
}

// Reply contains information about a reply for a thread
type Reply struct {
	User      string `json:"user,omitempty"`
	Timestamp string `json:"ts,omitempty"`
}

// Event contains the event type
type Event struct {
	Type string `json:"type,omitempty"`
}

// Ping contains information about a Ping Event
type Ping struct {
	ID        int    `json:"id"`
	Type      string `json:"type"`
	Timestamp int64  `json:"timestamp"`
}

// Pong contains information about a Pong Event
type Pong struct {
	Type      string `json:"type"`
	ReplyTo   int    `json:"reply_to"`
	Timestamp int64  `json:"timestamp"`
}

// NewOutgoingMessage prepares an OutgoingMessage that the user can
// use to send a message. Use this function to properly set the
// messageID.
func (rtm *RTM) NewOutgoingMessage(text string, channelID string, options ...RTMsgOption) *OutgoingMessage {
	id := rtm.idGen.Next()
	msg := OutgoingMessage{
		ID:      id,
		Type:    "message",
		Channel: channelID,
		Text:    text,
	}
	for _, option := range options {
		option(&msg)
	}
	return &msg
}

// NewSubscribeUserPresence prepares an OutgoingMessage that the user can
// use to subscribe presence events for the specified users.
func (rtm *RTM) NewSubscribeUserPresence(ids []string) *OutgoingMessage {
	return &OutgoingMessage{
		Type: "presence_sub",
		IDs:  ids,
	}
}

// NewTypingMessage prepares an OutgoingMessage that the user can
// use to send as a typing indicator. Use this function to properly set the
// messageID.
func (rtm *RTM) NewTypingMessage(channelID string) *OutgoingMessage {
	id := rtm.idGen.Next()
	return &OutgoingMessage{
		ID:      id,
		Type:    "typing",
		Channel: channelID,
	}
}

// RTMsgOption allows configuration of various options available for sending an RTM message
type RTMsgOption func(*OutgoingMessage)

// RTMsgOptionTS sets thead timestamp of an outgoing message in order to respond to a thread
func RTMsgOptionTS(threadTimestamp string) RTMsgOption {
	return func(msg *OutgoingMessage) {
		msg.ThreadTimestamp = threadTimestamp
	}
}

// RTMsgOptionBroadcast sets broadcast reply to channel to "true"
func RTMsgOptionBroadcast() RTMsgOption {
	return func(msg *OutgoingMessage) {
		msg.ThreadBroadcast = true
	}
}
