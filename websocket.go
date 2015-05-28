package slack

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/url"
	"reflect"
	"sync"
	"time"

	"golang.org/x/net/websocket"
)

// SlackWS represents a managed websocket connection. It also supports all the methods of the `Slack` type.
type SlackWS struct {
	mutex     sync.Mutex
	messageId int
	pings     map[int]time.Time

	// Connection life-cycle
	conn             *websocket.Conn
	IncomingEvents   chan SlackEvent
	connectionErrors chan error
	killRoutines     chan bool

	// Slack is the main API, embedded
	Slack

	// UserDetails upon connection
	info *Info
}

// StartRTM starts a Websocket used to do all common chat client operations.
func (api *Slack) StartRTM() (*SlackWS, error) {
	response := &infoResponseFull{}
	err := post("rtm.start", url.Values{"token": {api.config.token}}, response, api.debug)
	if err != nil {
		return nil, err
	}
	if !response.Ok {
		return nil, response.Error
	}

	// websocket.Dial does not accept url without the port (yet)
	// Fixed by: https://github.com/golang/net/commit/5058c78c3627b31e484a81463acd51c7cecc06f3
	// but slack returns the address with no port, so we have to fix it
	websocketUrl, err := websocketizeUrlPort(response.Info.URL)
	if err != nil {
		return nil, err
	}

	ws := &SlackWS{
		Slack: *api,
		info:  &response.Info,
	}
	ws.pings = make(map[int]time.Time)

	ws.conn, err = websocket.Dial(websocketUrl, "", "")
	if err != nil {
		return nil, err
	}

	ws.IncomingEvents = make(chan SlackEvent, 50)
	ws.killRoutines = make(chan bool, 10)
	ws.connectionErrors = make(chan error, 10)
	go ws.manageConnection(websocketUrl)
	return ws, nil
}

func (ws *SlackWS) manageConnection(url string) {
	// launch keepalive with fixed timings
	// listen on
	// receive any connectionErrors, killall goroutines
	// reconnect and restart them all
}

// Disconnect and wait, blocking until a successful disconnection.
func (ws *SlackWS) Disconnect() error {
	return nil
}

// Reconnect, only makes sense if you've successfully disconnectd with Disconnect().
func (ws *SlackWS) Reconnect() error {
	return nil
}

// GetInfo returns the info structure received when calling
// "startrtm", holding all channels, groups and other metadata needed
// to implement a full chat client. It will be non-nil after a call to
// StartRTM().
func (ws *SlackWS) GetInfo() *Info {
	return ws.info
}

func (ws *SlackWS) Ping() error {
	ws.mutex.Lock()
	defer ws.mutex.Unlock()
	ws.messageId++
	msg := &Ping{Id: ws.messageId, Type: "ping"}
	if err := websocket.JSON.Send(ws.conn, msg); err != nil {
		return err
	}
	// TODO: What happens if we already have this id?
	ws.pings[ws.messageId] = time.Now()
	return nil
}

func (ws *SlackWS) Keepalive(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if err := ws.Ping(); err != nil {
				log.Fatal(err)
			}
		}
	}
}

func (ws *SlackWS) SendMessage(msg *OutgoingMessage) error {
	ws.mutex.Lock()
	defer ws.mutex.Unlock()

	if msg == nil {
		return fmt.Errorf("Can't send a nil message")
	}

	if err := websocket.JSON.Send(ws.conn, *msg); err != nil {
		return err
	}
	return nil
}

func (ws *SlackWS) HandleIncomingEvents(ch chan SlackEvent) {
	for {
		event := json.RawMessage{}
		if err := websocket.JSON.Receive(ws.conn, &event); err == io.EOF {
			//log.Println("Derpi derp, should we destroy conn and start over?")
			//if err = ws.StartRTM(); err != nil {
			//	log.Fatal(err)
			//}
			// should we reconnect here?

			if !ws.conn.IsClientConn() {
				// FIXME: take the URL from the connection manager.
				url := "boo"
				ws.conn, err = websocket.Dial(url, "", "")
				if err != nil {
					log.Panic(err)
				}
			}
			// XXX: check for timeout and implement exponential backoff
		} else if err != nil {
			log.Panic(err)
		}
		if len(event) == 0 {
			if ws.debug {
				log.Println("Event Empty. WTF?")
			}
		} else {
			if ws.debug {
				log.Println(string(event[:]))
			}
			ws.handleEvent(ch, event)
		}
		time.Sleep(time.Millisecond * 500)
	}
}

func (ws *SlackWS) handleEvent(ch chan SlackEvent, event json.RawMessage) {
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
			// FIXME: never do that mama!
			log.Fatal(err)
		}

		if ack.Ok {
			ch <- SlackEvent{"ack", ack}
		} else {
			ch <- SlackEvent{"error", ack.Error}
		}
	case "hello":
		ch <- SlackEvent{"hello", HelloEvent{}}
	case "pong":
		pong := Pong{}
		if err = json.Unmarshal(event, &pong); err != nil {
			log.Fatal(err)
		}
		ws.mutex.Lock()
		latency := time.Since(ws.pings[pong.ReplyTo])
		ws.mutex.Unlock()
		ch <- SlackEvent{"latency-report", LatencyReport{Value: latency}}
	default:
		callEvent(em.Type, ch, event)
	}
}

func callEvent(eventType string, ch chan SlackEvent, event json.RawMessage) {
	for k, v := range eventMapping {
		if eventType == k {
			t := reflect.TypeOf(v)
			recvEvent := reflect.New(t).Interface()
			err := json.Unmarshal(event, recvEvent)
			if err != nil {
				log.Println("Unable to unmarshal event:", eventType, event)
			}
			ch <- SlackEvent{k, recvEvent}
			return
		}
	}
	log.Printf("XXX: Not implemented yet: %s -> %v", eventType, event)
}

var eventMapping = map[string]interface{}{
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
