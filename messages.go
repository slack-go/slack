package slack

type OutgoingMessage struct {
	Id        int    `json:"id"`
	ChannelId string `json:"channel,omitempty"`
	Text      string `json:"text,omitempty"`
	Type      string `json:"type,omitempty"`
}

// Message is an auxiliary type to allow us to have a message containing sub messages
type Message struct {
	Msg
	SubMessage Msg `json:"message,omitempty"`
}

// Msg contains information about a slack message
type Msg struct {
	Id        string `json:"id"`
	BotId     string `json:"bot_id,omitempty"`
	UserId    string `json:"user,omitempty"`
	Username  string `json:"username,omitempty"`
	ChannelId string `json:"channel,omitempty"`
	Timestamp string `json:"ts,omitempty"`
	Text      string `json:"text,omitempty"`
	Team      string `json:"team,omitempty"`
	File      File   `json:"file,omitempty"`
	// Type may come if it's part of a message list
	// e.g.: channel.history
	Type      string `json:"type,omitempty"`
	IsStarred bool   `json:"is_starred,omitempty"`
	// Submessage
	SubType          string       `json:"subtype,omitempty"`
	Hidden           bool         `json:"bool,omitempty"`
	DeletedTimestamp string       `json:"deleted_ts,omitempty"`
	Attachments      []Attachment `json:"attachments,omitempty"`
	ReplyTo          int          `json:"reply_to,omitempty"`
	Upload           bool         `json:"upload,omitempty"`
}

// Presence XXX: not used yet
type Presence struct {
	Presence string `json:"presence"`
	UserId   string `json:"user"`
}

// Event contains the event type
type Event struct {
	Type string `json:"type,omitempty"`
}

// Ping contains information about a Ping Event
type Ping struct {
	Id   int    `json:"id"`
	Type string `json:"type"`
}

// Pong contains information about a Pong Event
type Pong struct {
	Type    string `json:"type"`
	ReplyTo int    `json:"reply_to"`
}

// NewOutGoingMessage prepares an OutgoingMessage that the user can use to send a message
func (api *SlackWS) NewOutgoingMessage(text string, channel string) *OutgoingMessage {
	api.mutex.Lock()
	defer api.mutex.Unlock()
	api.messageId++
	return &OutgoingMessage{
		Id:        api.messageId,
		Type:      "message",
		ChannelId: channel,
		Text:      text,
	}
}
