package slack

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"net/url"
	"sync"
	"time"

	"golang.org/x/net/websocket"
)

const (
	EV_MESSAGE = iota
	EV_USER_TYPING
)

type SlackWS struct {
	conn      *websocket.Conn
	messageId int
	mutex     sync.Mutex
	Slack
}

type SlackWSResponse struct {
	Ok    bool          `json:"ok"`
	Error *SlackWSError `json:"error"`
}

type SlackWSError struct {
	Code int
	Msg  string
}

func (s SlackWSError) Error() string {
	return s.Msg
}

var portMapping = map[string]string{"ws": "80", "wss": "443"}

func fixUrlPort(orig string) (string, error) {
	urlObj, err := url.ParseRequestURI(orig)
	if err != nil {
		return "", err
	}
	_, _, err = net.SplitHostPort(urlObj.Host)
	if err != nil {
		return urlObj.Scheme + "://" + urlObj.Host + ":" + portMapping[urlObj.Scheme] + urlObj.Path, nil
	}
	return orig, nil
}

func (api *Slack) StartRTM(protocol, origin string) (*SlackWS, error) {
	response := &infoResponseFull{}
	err := parseResponse("rtm.start", url.Values{"token": {api.config.token}}, response, api.debug)
	if err != nil {
		return nil, err
	}
	if !response.Ok {
		return nil, response.Error
	}
	api.info = response.Info
	// websocket.Dial does not accept url without the port (yet)
	// Fixed by: https://github.com/golang/net/commit/5058c78c3627b31e484a81463acd51c7cecc06f3
	// but slack returns the address with no port, so we have to fix it
	api.info.Url, err = fixUrlPort(api.info.Url)
	if err != nil {
		return nil, err
	}
	api.config.protocol, api.config.origin = protocol, origin
	wsApi := &SlackWS{Slack: *api}
	wsApi.conn, err = websocket.Dial(api.info.Url, api.config.protocol, api.config.origin)
	if err != nil {
		return nil, err
	}
	return wsApi, nil
}

func (api *SlackWS) Ping() error {
	api.mutex.Lock()
	defer api.mutex.Unlock()
	api.messageId++
	msg := &Ping{Id: api.messageId, Type: "ping"}
	if err := websocket.JSON.Send(api.conn, msg); err != nil {
		return err
	}
	return nil
}

func (api *SlackWS) Keepalive(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if err := api.Ping(); err != nil {
				log.Fatal(err)
			}
		}
	}
}

func (api *SlackWS) SendMessage(msg *OutgoingMessage) error {
	if msg == nil {
		return fmt.Errorf("Can't send a nil message")
	}

	if err := websocket.JSON.Send(api.conn, *msg); err != nil {
		return err
	}
	return nil
}

func (api *SlackWS) HandleIncomingEvents(ch chan SlackEvent) {
	event := json.RawMessage{}
	for {
		if err := websocket.JSON.Receive(api.conn, &event); err == io.EOF {
			//log.Println("Derpi derp, should we destroy conn and start over?")
			//if err = api.StartRTM(); err != nil {
			//	log.Fatal(err)
			//}
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
			if api.debug {
				log.Println(string(event[:]))
			}
			handleEvent(ch, event)
		}
		time.Sleep(time.Millisecond * 500)
	}
}

func handleEvent(ch chan SlackEvent, event json.RawMessage) {
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
			return
		}

		// TODO: errors end up in this bucket. They shouldn't.
		log.Printf("Got error(?): %s", event)
	case "hello":
		return
	case "pong":
		// XXX: Eventually check to which ping this matched with
		//      Allows us to have stats about latency and what not
		return
	case "presence_change":
		//log.Printf("`%s is %s`\n", info.GetUserById(event.PUserId).Name, event.Presence)
	case "message":
		handleMessage(ch, event)
	case "channel_marked":
		log.Printf("XXX: To implement %s", em)
	case "user_typing":
		handleUserTyping(ch, event)
	default:
		log.Println("XXX: " + string(event))
	}
}

func handleUserTyping(ch chan SlackEvent, event json.RawMessage) {
	msg := UserTyping{}
	if err := json.Unmarshal(event, &msg); err != nil {
		log.Fatal(err)
	}
	ch <- SlackEvent{Type: EV_USER_TYPING, Data: msg}
}

func handleMessage(ch chan SlackEvent, event json.RawMessage) {
	msg := Message{}
	err := json.Unmarshal(event, &msg)
	if err != nil {
		log.Fatal(err)
	}
	ch <- SlackEvent{Type: EV_MESSAGE, Data: msg}
}
