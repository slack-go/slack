package slack

type OutgoingMessage struct {
	Id        int    `json:"id"`
	ChannelID string `json:"channel,omitempty"`
	Text      string `json:"text,omitempty"`
	Type      string `json:"type,omitempty"`
}

type Message struct {
	Msg
	SubMessage Msg `json:"message,omitempty"`
}

type Msg struct {
	Id        string `json:"id"`
	UserID    string `json:"user,omitempty"`
	ChannelID string `json:"channel,omitempty"`
	Timestamp string `json:"ts,omitempty"`
	Text      string `json:"text,omitempty"`
	// Type may come if it's part of a message list
	// e.g.: channel.history
	Type      string `json:"type,omitempty"`
	IsStarred bool   `json:"is_starred,omitempty"`
	// Submessage
	SubType          string `json:"subtype,omitempty"`
	Hidden           bool   `json:"bool,omitempty"`
	DeletedTimestamp string `json:"deleted_ts,omitempty"`
}

type Presence struct {
	Presence string `json:"presence"`
	UserID   string `json:"user"`
}

type Event struct {
	Type string `json:"type,omitempty"`
}

type Ping struct {
	Id   int    `json:"id"`
	Type string `json:"type"`
}

type AckMessage struct {
	Ok        bool   `json:"ok"`
	ReplyTo   int    `json:"reply_to"`
	Timestamp string `json:"ts"`
	Text      string `json:"text"`
}

var message_id int

func NewOutgoingMessage(text string, channel string) *OutgoingMessage {
	message_id++
	return &OutgoingMessage{
		Id:        message_id,
		Type:      "message",
		ChannelID: channel,
		Text:      text,
	}
}

// XXX: maybe support variable arguments so that people
// can send stuff through their ping
func NewPing() *Ping {
	message_id++
	return &Ping{Id: message_id, Type: "ping"}
}

func (info Info) GetUserById(id string) *User {
	for _, user := range info.Users {
		if user.Id == id {
			return &user
		}
	}
	return nil
}

func (info Info) GetChannelById(id string) *Channel {
	for _, channel := range info.Channels {
		if channel.Id == id {
			return &channel
		}
	}
	return nil
}
