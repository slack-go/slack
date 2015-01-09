package slack

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"net/url"
	"time"

	"golang.org/x/net/websocket"
)

type UserTyping struct {
	Type      string `json:"type"`
	UserID    string `json:"user"`
	ChannelID string `json:"channel"`
}

type SlackEvent struct {
	Type int
	Data interface{}
}

type SlackAPI struct {
	config Config
	conn   *websocket.Conn

	info Info
}

func New(token string) *SlackAPI {
	return &SlackAPI{
		config: Config{token: token},
	}
}

func (api *SlackAPI) GetInfo() Info {
	return api.info
}

func (api *SlackAPI) StartRTM(protocol string, origin string) error {
	resp, err := http.PostForm(SLACK_API+"rtm.start", url.Values{"token": {api.config.token}})
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	decoder := json.NewDecoder(resp.Body)
	if err = decoder.Decode(&api.info); err != nil {
		return err
	}
	if !api.info.Ok {
		return errors.New(api.info.Error)
	}
	api.config.protocol, api.config.origin = protocol, origin
	api.conn, err = websocket.Dial(api.info.Url, api.config.protocol, api.config.origin)
	if err != nil {
		return err
	}
	return nil
}

func (api *SlackAPI) Ping() error {
	if err := websocket.JSON.Send(api.conn, NewPing()); err != nil {
		return err
	}
	return nil
}

func (api *SlackAPI) SendMessage(msg OutgoingMessage) error {
	if err := websocket.JSON.Send(api.conn, msg); err != nil {
		return err
	}
	return nil
}

func (api *SlackAPI) HandleIncomingEvents(ch *chan SlackEvent) {
	event := json.RawMessage{}
	for {
		if err := websocket.JSON.Receive(api.conn, &event); err == io.EOF {
			// should we reconnect here?
			if !api.conn.IsClientConn() {
				api.conn, err = websocket.Dial(api.info.Url, api.config.protocol, api.config.origin)
				if err != nil {
					log.Panic(err)
				}
			}
			// XXX: check for timeout and implement exponential backoff
		} else if err != nil {
			log.Panic(err)
		}
		if len(event) == 0 {
			log.Println("Event Empty. WTF?")
		} else {
			handle_event(ch, event)
		}
		time.Sleep(time.Millisecond * 500)
	}
}

func handle_event(ch *chan SlackEvent, event json.RawMessage) {
	em := Event{}
	err := json.Unmarshal(event, &em)
	if err != nil {
		log.Fatal(err)
	}
	switch em.Type {
	case "":
		// try ok
		ack := AckMessage{}
		if err = json.Unmarshal(event, &ack); err != nil {
			log.Fatal(err)
		}
		if ack.Ok {
			log.Printf("Received an ok for: %d", ack.ReplyTo)
		} else {
			log.Println(event)
			log.Println("XXX: ?")
		}
	case "hello":
		return
	case "pong":
		// XXX: Eventually check to which ping this matched with
		//      Allows us to have stats about latency and what not
		return
	case "presence_change":
		//log.Printf("`%s is %s`\n", info.GetUserById(event.PUserID).Name, event.Presence)
	case "message":
		handle_message(ch, event)
	case "channel_marked":
		log.Printf("XXX: To implement %s", em)
	case "user_typing":
		handle_user_typing(ch, event)
	default:
		log.Println("XXX: " + string(event))
	}
}

func handle_user_typing(ch *chan SlackEvent, event json.RawMessage) {
	msg := UserTyping{}
	if err := json.Unmarshal(event, &msg); err != nil {
		log.Fatal(err)
	}
	*ch <- SlackEvent{Type: EV_USER_TYPING, Data: msg}
}

func handle_message(ch *chan SlackEvent, event json.RawMessage) {
	msg := Message{}
	err := json.Unmarshal(event, &msg)
	if err != nil {
		log.Fatal(err)
	}
	*ch <- SlackEvent{Type: EV_MESSAGE, Data: msg}
}
