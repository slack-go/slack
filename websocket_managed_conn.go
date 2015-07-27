package slack

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"reflect"
	"time"

	"golang.org/x/net/websocket"
)

// ManageConnection can be called on a Slack RTM instance returned by the
// NewRTM method. It will connect to the slack RTM API and handle all incoming
// and outgoing events. If a connection fails then it will attempt to reconnect
// and will notify any listeners through an error event on the IncomingEvents
// channel.
//
// This detects failed and closed connections through the RTM's keepRunning
// channel. Once this channel is closed or has something sent to it, this will
// open a lock on the RTM's mutex and check if the disconnect was intentional
// or not. If it was not then it attempts to reconnect.
//
// The defined error events are located in websocket_internals.go.
func (rtm *RTM) ManageConnection() {
	var connectionCount int
	for {
		// open a lock - we want to close this before returning from the
		// function so we won't defer the mutex's close therefore we MUST
		// release the lock before returning on an error!
		rtm.mutex.Lock()
		connectionCount++
		// start trying to connect
		// the returned err is already passed onto the IncomingEvents channel
		info, conn, err := rtm.connect(connectionCount)
		log.Println(err)
		// if err != nil then the connection is sucessful
		// otherwise we need to send a Disconnected event
		if err != nil {
			rtm.IncomingEvents <- SlackEvent{"disconnected", &DisconnectedEvent{
				Intentional: false,
			}}
			rtm.mutex.Unlock()
			return
		}
		rtm.info = info
		rtm.IncomingEvents <- SlackEvent{"connected", &ConnectedEvent{
			ConnectionCount: connectionCount,
			Info:            info,
		}}
		// set the connection object and unlock the mutex
		rtm.conn = conn
		rtm.isConnected = true
		rtm.keepRunning = make(chan bool)
		rtm.mutex.Unlock()

		// we're now connected (or have failed fatally) so we can set up
		// listeners and monitor for stopping
		go rtm.sendKeepAlive(30 * time.Second)
		go rtm.handleIncomingEvents()
		go rtm.handleOutgoingMessages()

		// should return only once we are disconnected
		<-rtm.keepRunning

		// after being disconnected we need to check if it was intentional
		// if not then we should try to reconnect
		rtm.mutex.Lock()
		intentional := rtm.wasIntentional
		rtm.mutex.Unlock()
		if intentional {
			return
		}
		// else continue and run the loop again to connect
	}
}

// connect attempts to connect to the slack websocket API. It handles any
// errors that occur while connecting and will return once a connection
// has been successfully opened.
func (rtm *RTM) connect(connectionCount int) (*Info, *websocket.Conn, error) {
	// used to provide exponential backoff wait time with jitter before trying
	// to connect to slack again
	boff := &backoff{
		Min:    100 * time.Millisecond,
		Max:    5 * time.Minute,
		Factor: 2,
		Jitter: true,
	}

	for {
		// send connecting event
		rtm.IncomingEvents <- SlackEvent{"connecting", &ConnectingEvent{
			Attempt:         boff.attempts + 1,
			ConnectionCount: connectionCount,
		}}
		// attempt to start the connection
		info, conn, err := rtm.startRTMAndDial()
		if err == nil {
			return info, conn, nil
		}
		// check for fatal errors - currently only invalid_auth
		if sErr, ok := err.(*SlackWebError); ok && sErr.Error() == "invalid_auth" {
			rtm.IncomingEvents <- SlackEvent{"invalid_auth", &InvalidAuthEvent{}}
			return nil, nil, sErr
		}
		// any other errors are treated as recoverable and we try again after
		// sending the event along the IncomingEvents channel
		rtm.IncomingEvents <- SlackEvent{"connection_error", &ConnectionErrorEvent{
			Attempt:  boff.attempts,
			ErrorObj: err,
		}}
		// get time we should wait before attempting to connect again
		dur := boff.Duration()
		rtm.Debugf("reconnection %d failed: %s", boff.attempts+1, err)
		rtm.Debugln(" -> reconnecting in", dur)
		time.Sleep(dur)
	}
}

// startRTMAndDial attemps to connect to the slack websocket. It returns the
// full information returned by the "rtm.start" method on the slack API.
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

// killConnection stops the websocket connection and signals to all goroutines
// that they should cease listening to the connection for events.
//
// This requires that a lock on the RTM's mutex is held before being called.
func (rtm *RTM) killConnection(intentional bool) error {
	rtm.Debugln("killing connection")
	if rtm.isConnected {
		close(rtm.keepRunning)
	}
	rtm.isConnected = false
	rtm.wasIntentional = intentional
	err := rtm.conn.Close()
	rtm.IncomingEvents <- SlackEvent{"disconnected", &DisconnectedEvent{intentional}}
	return err
}

// handleOutgoingMessages listens on the outgoingMessages channel for any
// queued messages that have not been sent.
//
// This will stop executing once the RTM's keepRunning channel has been closed
// or has anything sent to it.
func (rtm *RTM) handleOutgoingMessages() {
	for {
		select {
		// catch "stop" signal on channel close
		case <-rtm.keepRunning:
			return
		// listen for messages that need to be sent
		case msg := <-rtm.outgoingMessages:
			rtm.sendOutgoingMessage(msg)
		}
	}
}

// sendOutgoingMessage sends the given OutgoingMessage to the slack websocket
// after acquiring a lock on the RTM's mutex.
//
// It does not currently detect if a outgoing message fails due to a disconnect
// and instead lets a future failed 'PING' detect the failed connection.
func (rtm *RTM) sendOutgoingMessage(msg OutgoingMessage) {
	rtm.mutex.Lock()
	defer rtm.mutex.Unlock()
	rtm.Debugln("Sending message:", msg)
	if !rtm.isConnected {
		// check for race condition of connection closed after lock
		// obtained
		rtm.IncomingEvents <- SlackEvent{"outgoing_error", &OutgoingErrorEvent{
			Message:  msg,
			ErrorObj: errors.New("Cannot send message - API is not connected"),
		}}
		return
	}
	if len(msg.Text) > maxMessageTextLength {
		rtm.IncomingEvents <- SlackEvent{"outgoing_error", &MessageTooLongEvent{
			Message:   msg,
			MaxLength: maxMessageTextLength,
		}}
		return
	}
	err := websocket.JSON.Send(rtm.conn, msg)
	if err != nil {
		rtm.IncomingEvents <- SlackEvent{"outgoing_error", &OutgoingErrorEvent{
			Message:  msg,
			ErrorObj: err,
		}}
	}
}

// sendKeepAlive is a blocking call that sends a 'PING' message once for every
// duration elapsed.
//
// This will stop executing once the RTM's keepRunning channel has been closed
// or has anything sent to it.
func (rtm *RTM) sendKeepAlive(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		// catch "stop" signal on channel close
		case <-rtm.keepRunning:
			return
		// send pings on ticker interval
		case <-ticker.C:
			go rtm.ping()
		}
	}
}

// ping sends a 'PING' message to the RTM's websocket. If the 'PING' message
// fails to send then this calls killConnection to signal an unintentional
// websocket disconnect.
//
// This does not handle incoming 'PONG' responses but does store the time of
// each successful 'PING' send so latency can be detected upon a 'PONG'
// response.
func (rtm *RTM) ping() {
	rtm.mutex.Lock()
	defer rtm.mutex.Unlock()
	rtm.Debugln("Sending PING")
	if !rtm.isConnected {
		// it's possible that the API has disconnected while we were waiting
		// for a lock on the mutex
		rtm.Debugln("Cannot send ping - API is not connected")
		// no need to send an error event since it really isn't an error
		return
	}
	rtm.messageID++
	rtm.pings[rtm.messageID] = time.Now()

	msg := &Ping{ID: rtm.messageID, Type: "ping"}
	err := websocket.JSON.Send(rtm.conn, msg)
	if err != nil {
		rtm.Debugf("RTM Error sending 'PING': %s", err.Error())
		rtm.killConnection(false)
	}
}

// handleIncomingEvents monitors the RTM's opened websocket for any incoming
// events.
//
// This will stop executing once the RTM's keepRunning channel has been closed
// or has anything sent to it.
func (rtm *RTM) handleIncomingEvents() {
	for {
		// non-blocking listen to see if channel is closed
		select {
		// catch "stop" signal on channel close
		case <-rtm.keepRunning:
			return
		default:
			rtm.receiveIncomingEvent()
		}
	}
}

// receiveIncomingEvent attempts to receive an event from the RTM's websocket.
// This will block until a frame is available from the websocket.
func (rtm *RTM) receiveIncomingEvent() {
	event := json.RawMessage{}
	err := websocket.JSON.Receive(rtm.conn, &event)
	if err == io.EOF {
		// EOF's don't seem to signify a failed connection so instead we ignore
		// them here and detect a failed connection upon attempting to send a
		// 'PING' message

		// trigger a 'PING' to detect pontential websocket disconnect
		go rtm.ping()
		return
	} else if err != nil {
		// TODO detect if this is a fatal error
		rtm.IncomingEvents <- SlackEvent{"incoming_error", &IncomingEventError{
			ErrorObj: err,
		}}
		return
	} else if len(event) == 0 {
		rtm.Debugln("Received empty event")
		return
	}
	rtm.Debugln("Incoming Event:", string(event[:]))
	rtm.handleRawEvent(event)
}

// handleEOF should be called after receiving an EOF on the RTM's websocket.
// It calls the internal killConnection method if the RTM was still considered
// to be connected. If it is not considered connected then it is because
// the killConnection method has already been called elsewhere.
func (rtm *RTM) handleEOF() {
	rtm.Debugln("Received EOF on websocket")
	// we need a lock in order to access isConnected and to call killConnection
	rtm.mutex.Lock()
	defer rtm.mutex.Unlock()
	// if isConnected is true then we didn't expect the EOF event
	// so for it to be intentional we need to have it be false
	if rtm.isConnected {
		// try to kill the connection - this should fail silently if the
		// API has already disconnected
		_ = rtm.killConnection(false)
	}
}

// handleRawEvent takes a raw JSON message received from the slack websocket
// and handles the encoded event.
func (rtm *RTM) handleRawEvent(rawEvent json.RawMessage) {
	event := &Event{}
	err := json.Unmarshal(rawEvent, event)
	if err != nil {
		rtm.IncomingEvents <- SlackEvent{"unmarshalling_error", &UnmarshallingErrorEvent{err}}
		return
	}
	switch event.Type {
	case "":
		rtm.handleAck(rawEvent)
	case "hello":
		rtm.IncomingEvents <- SlackEvent{"hello", &HelloEvent{}}
	case "pong":
		rtm.handlePong(rawEvent)
	default:
		rtm.handleEvent(event.Type, rawEvent)
	}
}

// handleAck handles an incoming 'ACK' message.
func (rtm *RTM) handleAck(event json.RawMessage) {
	ack := &AckMessage{}
	if err := json.Unmarshal(event, ack); err != nil {
		rtm.Debugln("RTM Error unmarshalling 'ack' event:", err)
		rtm.Debugln(" -> Erroneous 'ack' event:", string(event))
		return
	}
	if ack.Ok {
		rtm.IncomingEvents <- SlackEvent{"ack", ack}
	} else {
		rtm.IncomingEvents <- SlackEvent{"ack_error", &AckErrorEvent{ack.Error}}
	}
}

// handlePong handles an incoming 'PONG' message which should be in response to
// a previously sent 'PING' message. This is then used to compute the
// connection's latency.
func (rtm *RTM) handlePong(event json.RawMessage) {
	pong := &Pong{}
	if err := json.Unmarshal(event, pong); err != nil {
		rtm.Debugln("RTM Error unmarshalling 'pong' event:", err)
		rtm.Debugln(" -> Erroneous 'ping' event:", string(event))
		return
	}
	rtm.mutex.Lock()
	defer rtm.mutex.Unlock()
	if pingTime, exists := rtm.pings[pong.ReplyTo]; exists {
		latency := time.Since(pingTime)
		rtm.IncomingEvents <- SlackEvent{"latency_report", &LatencyReport{Value: latency}}
		delete(rtm.pings, pong.ReplyTo)
	} else {
		rtm.Debugln("RTM Error - unmatched 'pong' event:", string(event))
	}
}

// handleEvent is the "default" response to an event that does not have a
// special case. It matches the command's name to a mapping of defined events
// and then sends the corresponding event struct to the IncomingEvents channel.
// If the event type is not found or the event cannot be unmarshalled into the
// correct struct then this sends an UnmarshallingErrorEvent to the
// IncomingEvents channel.
func (rtm *RTM) handleEvent(typeStr string, event json.RawMessage) {
	v, exists := eventMapping[typeStr]
	if !exists {
		rtm.Debugf("RTM Error, received unmapped event %q: %s\n", typeStr, string(event))
		err := fmt.Errorf("RTM Error: Received unmapped event %q: %s\n", typeStr, string(event))
		rtm.IncomingEvents <- SlackEvent{"unmarshalling_error", &UnmarshallingErrorEvent{err}}
		return
	}
	t := reflect.TypeOf(v)
	recvEvent := reflect.New(t).Interface()
	err := json.Unmarshal(event, recvEvent)
	if err != nil {
		rtm.Debugf("RTM Error, received unmapped event %q: %s\n", typeStr, string(event))
		err := fmt.Errorf("RTM Error: Could not unmarshall event %q: %s\n", typeStr, string(event))
		rtm.IncomingEvents <- SlackEvent{"unmarshalling_error", &UnmarshallingErrorEvent{err}}
		return
	}
	rtm.IncomingEvents <- SlackEvent{typeStr, recvEvent}
}

// eventMapping holds a mapping of event names to their corresponding struct
// implementations. The structs should be instances of the unmarshalling
// target for the matching event type.
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

	"reaction_added":   ReactionAddedEvent{},
	"reaction_removed": ReactionRemovedEvent{},

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
