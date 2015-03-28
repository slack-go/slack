package slack

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"net/url"
	"strconv"
	"sync"
	"time"

	"golang.org/x/net/websocket"
)

type MessageEvent Message

type SlackWS struct {
	conn      *websocket.Conn
	messageId int
	mutex     sync.Mutex
	pings     map[int]time.Time
	Slack
}

// AckMessage is used for messages received in reply to other messages
type AckMessage struct {
	ReplyTo   int    `json:"reply_to"`
	Timestamp string `json:"ts"`
	Text      string `json:"text"`
	SlackWSResponse
}

type SlackWSResponse struct {
	Ok    bool          `json:"ok"`
	Error *SlackWSError `json:"error"`
}

type SlackWSError struct {
	Code int
	Msg  string
}

type SlackEvent struct {
	Type uint64
	Data interface{}
}

type JSONTimeString string

// String converts the unix timestamp into a string
func (t JSONTimeString) String() string {
	if t == "" {
		return ""
	}
	floatN, err := strconv.ParseFloat(string(t), 64)
	if err != nil {
		log.Panicln(err)
		return ""
	}
	timeStr := int64(floatN)
	tm := time.Unix(int64(timeStr), 0)
	return fmt.Sprintf("\"%s\"", tm.Format("Mon Jan _2"))
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
	wsApi.pings = make(map[int]time.Time)
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
	// TODO: What happens if we already have this id?
	api.pings[api.messageId] = time.Now()
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
	for {
		event := json.RawMessage{}
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
			api.handleEvent(ch, event)
		}
		time.Sleep(time.Millisecond * 500)
	}
}

func (api *SlackWS) handleEvent(ch chan SlackEvent, event json.RawMessage) {
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
			// TODO: Send the ack back (is this useful?)
			//ch <- SlackEvent{Type: EventAck, Data: ack}
			log.Printf("Received an ok for: %d", ack.ReplyTo)
			return
		}

		// Send the error to the user
		ch <- SlackEvent{Data: ack.Error}
	case "hello":
		ch <- SlackEvent{Data: HelloEvent{}}
	case "pong":
		pong := Pong{}
		if err = json.Unmarshal(event, &pong); err != nil {
			log.Fatal(err)
		}
		api.mutex.Lock()
		latency := time.Since(api.pings[pong.ReplyTo])
		api.mutex.Unlock()
		ch <- SlackEvent{Data: LatencyReport{Value: latency}}
	default:
		callEvent(em.Type, ch, event)
	}
}

func callEvent(eventType string, ch chan SlackEvent, event json.RawMessage) {
	eventMapping := map[string]interface{}{
		"message":         &MessageEvent{},
		"presence_change": &PresenceChangeEvent{},
		"user_typing":     &UserTypingEvent{},

		"channel_marked":          &ChannelMarkedEvent{},
		"channel_created":         &ChannelCreatedEvent{},
		"channel_joined":          &ChannelJoinedEvent{},
		"channel_left":            &ChannelLeftEvent{},
		"channel_deleted":         &ChannelDeletedEvent{},
		"channel_rename":          &ChannelRenameEvent{},
		"channel_archive":         &ChannelArchiveEvent{},
		"channel_unarchive":       &ChannelUnarchiveEvent{},
		"channel_history_changed": &ChannelHistoryChangedEvent{},

		"im_created":         &IMCreatedEvent{},
		"im_open":            &IMOpenEvent{},
		"im_close":           &IMCloseEvent{},
		"im_marked":          &IMMarkedEvent{},
		"im_history_changed": &IMHistoryChangedEvent{},

		"group_marked":          &GroupMarkedEvent{},
		"group_open":            &GroupOpenEvent{},
		"group_joined":          &GroupJoinedEvent{},
		"group_left":            &GroupLeftEvent{},
		"group_close":           &GroupCloseEvent{},
		"group_rename":          &GroupRenameEvent{},
		"group_archive":         &GroupArchiveEvent{},
		"group_unarchive":       &GroupUnarchiveEvent{},
		"group_history_changed": &GroupHistoryChangedEvent{},

		"file_created":         &FileCreatedEvent{},
		"file_shared":          &FileSharedEvent{},
		"file_unshared":        &FileUnsharedEvent{},
		"file_public":          &FilePublicEvent{},
		"file_private":         &FilePrivateEvent{},
		"file_change":          &FileChangeEvent{},
		"file_deleted":         &FileDeletedEvent{},
		"file_comment_added":   &FileCommentAddedEvent{},
		"file_comment_edited":  &FileCommentEditedEvent{},
		"file_comment_deleted": &FileCommentDeletedEvent{},

		"star_added":   &StarAddedEvent{},
		"star_removed": &StarRemovedEvent{},

		"pref_change": &PrefChangeEvent{},

		"team_join":              &TeamJoinEvent{},
		"team_rename":            &TeamRenameEvent{},
		"team_pref_change":       &TeamPrefChangeEvent{},
		"team_domain_change":     &TeamDomainChangeEvent{},
		"team_migration_started": &TeamMigrationStartedEvent{},

		"manual_presence_change": &ManualPresenceChangeEvent{},

		"user_change": &UserChangeEvent{},

		"emoji_changed": &EmojiChangedEvent{},

		"commands_changed": &CommandsChangedEvent{},

		"email_domain_changed": &EmailDomainChangedEvent{},

		"bot_added":   &BotAddedEvent{},
		"bot_changed": &BotChangedEvent{},

		"accounts_changed": &AccountsChangedEvent{},
	}

	msg := eventMapping[eventType]
	if msg == nil {
		log.Printf("XXX: Not implemented yet: %s -> %v", eventType, event)
	}
	if err := json.Unmarshal(event, &msg); err != nil {
		log.Fatal(err)
	}
	ch <- SlackEvent{Data: msg}
}
