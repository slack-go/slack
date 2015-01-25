package slack

type OutgoingMessage struct {
	Id        int    `json:"id"`
	ChannelId string `json:"channel,omitempty"`
	Text      string `json:"text,omitempty"`
	Type      string `json:"type,omitempty"`
}

type Message struct {
	Msg
	SubMessage Msg `json:"message,omitempty"`
}

type Msg struct {
	Id        string `json:"id"`
	UserId    string `json:"user,omitempty"`
	Username  string `json:"username,omitempty"`
	ChannelId string `json:"channel,omitempty"`
	Timestamp string `json:"ts,omitempty"`
	Text      string `json:"text,omitempty"`
	Team      string `json:"team,omitempty"`
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
	UserId   string `json:"user"`
}

type Event struct {
	Type string `json:"type,omitempty"`
}

type Ping struct {
	Id   int    `json:"id"`
	Type string `json:"type"`
}

type AckMessage struct {
	ReplyTo   int    `json:"reply_to"`
	Timestamp string `json:"ts"`
	Text      string `json:"text"`
	SlackResponse
}

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
