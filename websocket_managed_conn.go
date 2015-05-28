package slack

import (
	"encoding/json"
	"fmt"
	"reflect"
	"time"

	"golang.org/x/net/websocket"
)

// ManageConnection is a long-running goroutine that handles
// reconnections and piping messages back and to `rtm.IncomingEvents`
// and `rtm.OutgoingMessages`.
//
// Usage would look like:
//
//     bot := slack.New("my-token")
//     rtm := bot.NewRTM()  // check err
//     setupYourHandlers(rtm.IncomingEvents, rtm.OutgoingMessages)
//     rtm.ManageConnection()
//
func (rtm *RTM) ManageConnection() {
	boff := &backoff{
		Min:    100 * time.Millisecond,
		Max:    5 * time.Minute,
		Factor: 2,
		Jitter: true,
	}
	connectionCount := 0

	for {
		var conn *websocket.Conn // use as first
		var err error
		var info *Info

		connectionCount += 1

		attempts := 1
		boff.Reset()
		for {
			rtm.IncomingEvents <- SlackEvent{"connecting", &ConnectingEvent{
				Attempt:         attempts,
				ConnectionCount: connectionCount,
			}}

			info, conn, err = rtm.startRTMAndDial()
			if err == nil {
				break // connected
			}

			dur := boff.Duration()
			rtm.Debugf("reconnection %d failed: %s", attempts, err)
			rtm.Debugln(" -> reconnecting in", dur)
			attempts += 1
			time.Sleep(dur)
		}

		rtm.IncomingEvents <- SlackEvent{"connected", &ConnectedEvent{
			ConnectionCount: connectionCount,
			Info:            info,
		}}

		killCh := make(chan bool, 3)
		connErrors := make(chan error, 10) // in case we get many such errors

		go rtm.keepalive(30*time.Second, conn, killCh, connErrors)
		go rtm.handleIncomingEvents(conn, killCh, connErrors)
		go rtm.handleOutgoingMessages(conn, killCh, connErrors)

		// Here, block and await for disconnection, if it ever happens.
		err = <-connErrors

		rtm.Debugln("RTM connection error:", err)
		rtm.IncomingEvents <- SlackEvent{"disconnected", &DisconnectedEvent{}}
		killCh <- true // 3 child go-routines
		killCh <- true
		killCh <- true
	}
}

func (rtm *RTM) startRTMAndDial() (*Info, *websocket.Conn, error) {
	info, url, err := rtm.StartRTM()
	if err != nil {
		return nil, nil, err
	}

	conn, err := websocket.Dial(url, "", "http://api.slack.com")
	if err != nil {
		return nil, nil, err
	}

	return info, conn, err
}

func (rtm *RTM) keepalive(interval time.Duration, conn *websocket.Conn, killCh chan bool, errors chan error) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-killCh:
			return
		case <-ticker.C:
			rtm.ping(conn, errors)
		}
	}
}

func (rtm *RTM) ping(conn *websocket.Conn, errors chan error) {
	rtm.mutex.Lock()
	defer rtm.mutex.Unlock()

	rtm.messageId++
	rtm.pings[rtm.messageId] = time.Now()

	msg := &Ping{Id: rtm.messageId, Type: "ping"}
	if err := websocket.JSON.Send(conn, msg); err != nil {
		errors <- fmt.Errorf("error sending 'ping': %s", err)
	}
}

func (rtm *RTM) handleOutgoingMessages(conn *websocket.Conn, killCh chan bool, errors chan error) {
	// we pass "conn" in case we do a reconnection, in that case we'll
	// have a new `conn` even though we're dealing with the same
	// incoming and outgoing channels for messages/events.
	for {
		select {
		case <-killCh:
			return

		case msg := <-rtm.outgoingMessages:
			rtm.Debugln("Sending message:", msg)

			rtm.mutex.Lock()
			err := websocket.JSON.Send(conn, msg)
			rtm.mutex.Unlock()
			if err != nil {
				errors <- fmt.Errorf("error sending 'message': %s", err)
				return
			}
		}
	}
}

func (rtm *RTM) handleIncomingEvents(conn *websocket.Conn, killCh chan bool, errors chan error) {
	for {

		select {
		case <-killCh:
			return
		default:
		}

		event := json.RawMessage{}
		err := websocket.JSON.Receive(conn, &event)
		if err != nil {
			errors <- err
			return
		}
		if len(event) == 0 {
			//log.Println("Event Empty. WTF?")
			continue
		}

		rtm.Debugln("Incoming event:", string(event[:]))

		rtm.handleEvent(event)

		// FIXME: please I hope we don't need to sleep!!!
		//time.Sleep(time.Millisecond * 500)
	}
}

func (rtm *RTM) handleEvent(event json.RawMessage) {
	em := Event{}
	err := json.Unmarshal(event, &em)
	if err != nil {
		rtm.Debugln("RTM Error unmarshalling event:", err)
		rtm.Debugln(" -> Erroneous event:", string(event))
		return
	}

	switch em.Type {
	case "":
		// try ok
		ack := AckMessage{}
		if err = json.Unmarshal(event, &ack); err != nil {
			rtm.Debugln("RTM Error unmarshalling 'ack' event:", err)
			rtm.Debugln(" -> Erroneous 'ack' event:", string(event))
			return
		}

		if ack.Ok {
			rtm.IncomingEvents <- SlackEvent{"ack", ack}
		} else {
			rtm.IncomingEvents <- SlackEvent{"error", ack.Error}
		}

	case "hello":
		rtm.IncomingEvents <- SlackEvent{"hello", &HelloEvent{}}

	case "pong":
		pong := Pong{}
		if err = json.Unmarshal(event, &pong); err != nil {
			rtm.Debugln("RTM Error unmarshalling 'pong' event:", err)
			rtm.Debugln(" -> Erroneous 'ping' event:", string(event))
			return
		}

		rtm.mutex.Lock()
		latency := time.Since(rtm.pings[pong.ReplyTo])
		rtm.mutex.Unlock()

		rtm.IncomingEvents <- SlackEvent{"latency-report", &LatencyReport{Value: latency}}

	default:
		for k, v := range eventMapping {
			if em.Type == k {
				t := reflect.TypeOf(v)
				recvEvent := reflect.New(t).Interface()

				err := json.Unmarshal(event, recvEvent)
				if err != nil {
					rtm.Debugf("RTM Error unmarshalling %q event: %s", em.Type, err)
					rtm.Debugf(" -> Erroneous %q event: %s", em.Type, string(event))
					return
				}

				rtm.IncomingEvents <- SlackEvent{em.Type, recvEvent}
				return
			}
		}

		rtm.Debugf("RTM Error, received unmapped event %q: %s\n", em.Type, string(event))
	}
}

var eventMapping = map[string]interface{}{
	"message":         MessageEvent{},
	"presence_change": PresenceChangeEvent{},
	"user_typing":     UserTypingEvent{},

	"channel_marked":          ChannelMarkedEvent{},
	"channel_created":         ChannelCreatedEvent{},
	"channel_joined":          ChannelJoinedEvent{},
	"channel_left":            ChannelLeftEvent{},
	"channel_deleted":         ChannelDeletedEvent{},
	"channel_rename":          ChannelRenameEvent{},
	"channel_archive":         ChannelArchiveEvent{},
	"channel_unarchive":       ChannelUnarchiveEvent{},
	"channel_history_changed": ChannelHistoryChangedEvent{},

	"im_created":         IMCreatedEvent{},
	"im_open":            IMOpenEvent{},
	"im_close":           IMCloseEvent{},
	"im_marked":          IMMarkedEvent{},
	"im_history_changed": IMHistoryChangedEvent{},

	"group_marked":          GroupMarkedEvent{},
	"group_open":            GroupOpenEvent{},
	"group_joined":          GroupJoinedEvent{},
	"group_left":            GroupLeftEvent{},
	"group_close":           GroupCloseEvent{},
	"group_rename":          GroupRenameEvent{},
	"group_archive":         GroupArchiveEvent{},
	"group_unarchive":       GroupUnarchiveEvent{},
	"group_history_changed": GroupHistoryChangedEvent{},

	"file_created":         FileCreatedEvent{},
	"file_shared":          FileSharedEvent{},
	"file_unshared":        FileUnsharedEvent{},
	"file_public":          FilePublicEvent{},
	"file_private":         FilePrivateEvent{},
	"file_change":          FileChangeEvent{},
	"file_deleted":         FileDeletedEvent{},
	"file_comment_added":   FileCommentAddedEvent{},
	"file_comment_edited":  FileCommentEditedEvent{},
	"file_comment_deleted": FileCommentDeletedEvent{},

	"star_added":   StarAddedEvent{},
	"star_removed": StarRemovedEvent{},

	"pref_change": PrefChangeEvent{},

	"team_join":              TeamJoinEvent{},
	"team_rename":            TeamRenameEvent{},
	"team_pref_change":       TeamPrefChangeEvent{},
	"team_domain_change":     TeamDomainChangeEvent{},
	"team_migration_started": TeamMigrationStartedEvent{},

	"manual_presence_change": ManualPresenceChangeEvent{},

	"user_change": UserChangeEvent{},

	"emoji_changed": EmojiChangedEvent{},

	"commands_changed": CommandsChangedEvent{},

	"email_domain_changed": EmailDomainChangedEvent{},

	"bot_added":   BotAddedEvent{},
	"bot_changed": BotChangedEvent{},

	"accounts_changed": AccountsChangedEvent{},
}
