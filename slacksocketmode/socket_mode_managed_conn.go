package slacksocketmode

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/internal/backoff"
	"github.com/slack-go/slack/internal/misc"
	"github.com/slack-go/slack/slackevents"
	"io"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/slack-go/slack/internal/errorsx"
	"github.com/slack-go/slack/internal/timex"
)

// Run can be called on a Slack RTM instance returned by the
// NewRTM method. It will connect to the slack RTM API and handle all incoming
// and outgoing events. If a connection fails then it will attempt to reconnect
// and will notify any listeners through an error event on the IncomingEvents
// channel.
//
// If the connection ends and the disconnect was unintentional then this will
// attempt to reconnect.
//
// This should only be called once per slack API! Otherwise expect undefined
// behavior.
//
// The defined error events are located in websocket_internals.go.
func (smc *Client) Run() {
	var (
		err  error
		info *slack.SocketModeConnection
		conn *websocket.Conn
	)

	for connectionCount := 0; ; connectionCount++ {
		// start trying to connect
		// the returned err is already passed onto the IncomingEvents channel
		if info, conn, err = smc.connect(connectionCount); err != nil {
			// when the connection is unsuccessful its fatal, and we need to bail out.
			smc.Debugf("Failed to connect with Socket Mode on try %d: %s", connectionCount, err)
			smc.disconnect()
			return
		}

		// lock to prevent data races with Disconnect particularly around isConnected
		// and conn.
		smc.mu.Lock()
		smc.conn = conn
		smc.info = info
		smc.mu.Unlock()

		smc.IncomingEvents <- smc.internalEvent("connected", &SocketModeConnectedEvent{
			ConnectionCount: connectionCount,
			Info:            info,
		})

		smc.Debugf("RTM connection succeeded on try %d", connectionCount)

		rawEvents := make(chan json.RawMessage)
		// we're now connected so we can set up listeners
		go smc.handleIncomingEvents(rawEvents)
		// this should be a blocking call until the connection has ended
		smc.handleEvents(rawEvents)

		select {
		case <-smc.disconnected:
			// after handle events returns we need to check if we're disconnected
			// when this happens we need to cleanup the newly created connection.
			if err = conn.Close(); err != nil {
				smc.Debugln("failed to close conn on disconnected RTM", err)
			}
			return
		default:
			// otherwise continue and run the loop again to reconnect
		}
	}
}

// connect attempts to connect to the slack websocket API. It handles any
// errors that occur while connecting and will return once a connection
// has been successfully opened.
func (smc *Client) connect(connectionCount int) (*slack.SocketModeConnection, *websocket.Conn, error) {
	const (
		errInvalidAuth      = "invalid_auth"
		errInactiveAccount  = "account_inactive"
		errMissingAuthToken = "not_authed"
	)

	// used to provide exponential backoff wait time with jitter before trying
	// to connect to slack again
	boff := &backoff.Backoff{
		Max: 5 * time.Minute,
	}

	for {
		var (
			backoff time.Duration
		)

		// send connecting event
		smc.IncomingEvents <- smc.internalEvent("connecting", &slack.ConnectingEvent{
			Attempt:         boff.Attempts() + 1,
			ConnectionCount: connectionCount,
		})

		// attempt to start the connection
		info, conn, err := smc.openAndDial()
		if err == nil {
			return info, conn, nil
		}

		// check for fatal errors
		switch err.Error() {
		case errInvalidAuth, errInactiveAccount, errMissingAuthToken:
			smc.Debugf("invalid auth when connecting with SocketMode: %s", err)
			return nil, nil, err
		default:
		}

		switch actual := err.(type) {
		case misc.StatusCodeError:
			if actual.Code == http.StatusNotFound {
				smc.Debugf("invalid auth when connecting with RTM: %s", err)
				smc.IncomingEvents <- smc.internalEvent("invalid_auth", &slack.InvalidAuthEvent{})
				return nil, nil, err
			}
		case *slack.RateLimitedError:
			backoff = actual.RetryAfter
		default:
		}

		backoff = timex.Max(backoff, boff.Duration())
		// any other errors are treated as recoverable and we try again after
		// sending the event along the IncomingEvents channel
		smc.IncomingEvents <- smc.internalEvent("connection_error", &slack.ConnectionErrorEvent{
			Attempt:  boff.Attempts(),
			Backoff:  backoff,
			ErrorObj: err,
		})

		// get time we should wait before attempting to connect again
		smc.Debugf("reconnection %d failed: %s reconnecting in %v\n", boff.Attempts(), err, backoff)

		// wait for one of the following to occur,
		// backoff duration has elapsed, killChannel is signalled, or
		// the smc finishes disconnecting.
		select {
		case <-time.After(backoff): // retry after the backoff.
		case intentional := <-smc.killChannel:
			if intentional {
				smc.killConnection(intentional, slack.ErrRTMDisconnected)
				return nil, nil, slack.ErrRTMDisconnected
			}
		case <-smc.disconnected:
			return nil, nil, slack.ErrRTMDisconnected
		}
	}
}

// openAndDial attempts to open a Socket Mode connection and dial to the connection endpoint using WebSocket.
// It returns the  full information returned by the "apps.connections.open" method on the
// Slack API.
func (smc *Client) openAndDial() (info *slack.SocketModeConnection, _ *websocket.Conn, err error) {
	var (
		url string
	)

	smc.Debugf("Starting SocketMode")
	info, url, err = smc.Open()

	if err != nil {
		smc.Debugf("Failed to start or connect with SocketMode: %s", err)
		return nil, nil, err
	}

	smc.Debugf("Dialing to websocket on url %s", url)
	// Only use HTTPS for connections to prevent MITM attacks on the connection.
	upgradeHeader := http.Header{}
	upgradeHeader.Add("Origin", "https://api.slack.com")
	dialer := websocket.DefaultDialer
	if smc.dialer != nil {
		dialer = smc.dialer
	}
	conn, _, err := dialer.Dial(url, upgradeHeader)
	if err != nil {
		smc.Debugf("Failed to dial to the websocket: %s", err)
		return nil, nil, err
	}

	conn.SetPingHandler(func(appData string) error {
		smc.handlePing(json.RawMessage([]byte(appData)))

		return nil
	})

	conn.SetCloseHandler(func(code int, text string) error {
		smc.handleClose(code, text)

		return nil
	})

	return info, conn, err
}

// killConnection stops the websocket connection and signals to all goroutines
// that they should cease listening to the connection for events.
//
// This should not be called directly! Instead a boolean value (true for
// intentional, false otherwise) should be sent to the killChannel on the RTM.
func (smc *Client) killConnection(intentional bool, cause error) (err error) {
	smc.Debugln("killing connection", cause)

	if smc.conn != nil {
		err = smc.conn.Close()
	}

	smc.IncomingEvents <- smc.internalEvent("disconnected", &slack.DisconnectedEvent{Intentional: intentional, Cause: cause})

	if intentional {
		smc.disconnect()
	}

	return err
}

// handleEvents is a blocking function that handles all events. This sends
// pings when asked to (on rtm.forcePing) and upon every given elapsed
// interval. This also sends outgoing messages that are received from the RTM's
// outgoingMessages channel. This also handles incoming raw events from the RTM
// rawEvents channel.
func (smc *Client) handleEvents(events chan json.RawMessage) {
	ticker := time.NewTicker(smc.pingInterval)
	defer ticker.Stop()
	for {
		select {
		// catch "stop" signal on channel close
		case intentional := <-smc.killChannel:
			_ = smc.killConnection(intentional, errorsx.String("signaled"))
			return
		// detect when the connection is dead.
		case <-smc.pingDeadman.C:
			_ = smc.killConnection(false, slack.ErrRTMDeadman)
			return
		// listen for messages that need to be sent
		case msg := <-smc.outgoingMessages:
			smc.sendOutgoingMessage(msg)
			// listen for incoming messages that need to be parsed
		case rawEvent := <-events:
			_ = smc.handleRawEvent(rawEvent)
		}
	}
}

// handleIncomingEvents monitors the RTM's opened websocket for any incoming
// events. It pushes the raw events into the channel.
//
// This will stop executing once the RTM's when a fatal error is detected, or
// a disconnect occurs.
func (smc *Client) handleIncomingEvents(events chan json.RawMessage) {
	for {
		if err := smc.receiveIncomingEvent(events); err != nil {
			select {
			case smc.killChannel <- false:
			case <-smc.disconnected:
			}
			return
		}
	}
}

func (smc *Client) sendWithDeadline(msg interface{}) error {
	// set a write deadline on the connection
	if err := smc.conn.SetWriteDeadline(time.Now().Add(10 * time.Second)); err != nil {
		return err
	}
	if err := smc.conn.WriteJSON(msg); err != nil {
		return err
	}
	// remove write deadline
	return smc.conn.SetWriteDeadline(time.Time{})
}

func (smc *Client) internalEvent(tpe string, data interface{}) SocketModeEvent {
	return SocketModeEvent{Type: tpe, Data: data}
}

func (smc *Client) externalEvent(tpe string, data interface{}) SocketModeEvent {
	return SocketModeEvent{Type: tpe, Data: data}
}

// sendOutgoingMessage sends the given OutgoingMessage to the slack websocket.
//
// It does not currently detect if a outgoing message fails due to a disconnect
// and instead lets a future failed 'PING' detect the failed connection.
func (smc *Client) sendOutgoingMessage(msg slack.OutgoingMessage) {
	smc.Debugln("Sending message:", msg)
	if len([]rune(msg.Text)) > slack.MaxMessageTextLength {
		smc.IncomingEvents <- smc.internalEvent("outgoing_error", &slack.MessageTooLongEvent{
			Message:   msg,
			MaxLength: slack.MaxMessageTextLength,
		})
		return
	}

	if err := smc.sendWithDeadline(msg); err != nil {
		smc.IncomingEvents <- smc.internalEvent("outgoing_error", &slack.OutgoingErrorEvent{
			Message:  msg,
			ErrorObj: err,
		})
	}
}

// ack tells Slack that the we have received the SocketModeRequest denoted by the envelope ID,
// by sending back the envelope ID over the WebSocket connection.
func (smc *Client) ack(envelopeID string) error {
	smc.Debugln("Sending ACK ", envelopeID)

	// See https://github.com/slackapi/node-slack-sdk/blob/c3f4d7109062a0356fb765d53794b7b5f6b3b5ae/packages/socket-mode/src/SocketModeClient.ts#L417
	msg := map[string]interface{}{"envelope_id": envelopeID}

	if err := smc.sendWithDeadline(msg); err != nil {
		smc.Debugf("RTM Error sending 'ACK %s': %s", envelopeID, err.Error())
		return err
	}
	return nil
}

// receiveIncomingEvent attempts to receive an event from the RTM's websocket.
// This will block until a frame is available from the websocket.
// If the read from the websocket results in a fatal error, this function will return non-nil.
func (smc *Client) receiveIncomingEvent(events chan json.RawMessage) error {
	event := json.RawMessage{}
	err := smc.conn.ReadJSON(&event)

	// check if the connection was closed.
	if websocket.IsUnexpectedCloseError(err) {
		return err
	}

	switch {
	case err == io.ErrUnexpectedEOF:
		// EOF's don't seem to signify a failed connection so instead we ignore
		// them here and detect a failed connection upon attempting to send a
		// 'PING' message

		// Unlike RTM, we don't ping from the our end as there seem to have no client ping.
		// We just continue to the next loop so that we `smc.disconnected` should be received if
		// this EOF error was actually due to disconnection.

		return nil
	case err != nil:
		// All other errors from ReadJSON come from NextReader, and should
		// kill the read loop and force a reconnect.
		smc.IncomingEvents <- smc.internalEvent("incoming_error", &slack.IncomingEventError{
			ErrorObj: err,
		})

		return err
	case len(event) == 0:
		smc.Debugln("Received empty event")
	default:
		buf := &bytes.Buffer{}
		d := json.NewEncoder(buf)
		d.SetIndent("", "  ")
		if err := d.Encode(event); err != nil {
			smc.Debugln("Failed encoding decoded json:", err)
		}
		reencoded := buf.String()

		smc.Debugln("Incoming Event:", reencoded)
		select {
		case events <- event:
		case <-smc.disconnected:
			smc.Debugln("disonnected while attempting to send raw event")
		}
	}

	return nil
}

// handleRawEvent takes a raw JSON message received from the slack websocket
// and handles the encoded event.
// returns the event type of the message.
func (smc *Client) handleRawEvent(rawEvent json.RawMessage) string {
	event := &Message{}
	err := json.Unmarshal(rawEvent, event)
	if err != nil {
		smc.IncomingEvents <- smc.internalEvent("unmarshalling_error", &slack.UnmarshallingErrorEvent{err})
		return ""
	}

	str := string(rawEvent)
	println(str)

	// See https://github.com/slackapi/node-slack-sdk/blob/main/packages/socket-mode/src/SocketModeClient.ts#L533
	// for all the available event types.
	switch event.Type {
	case socketModeEventTypeHello:
		smc.IncomingEvents <- smc.externalEvent("hello", &slack.HelloEvent{})
	case "events_api":
		payloadEvent := event.Payload

		eventsAPIEvent, err := slackevents.ParseEvent(payloadEvent, slackevents.OptionNoVerifyToken())
		if err != nil {
			return ""
		}

		smc.IncomingEvents <- smc.externalEvent(eventsAPIEvent.Type, eventsAPIEvent)

		// We automatically ack the message.
		// TODO Should there be any way to manually ack the msg, like the official nodejs client?
		smc.ack(event.EnvelopeID)
	case "disconnect":

	default:
		panic(fmt.Errorf("unexpected type %q: %v", event.Type, event))
	}

	return event.Type
}

// handleAck handles an incoming 'ACK' message.
func (smc *Client) handleAck(event json.RawMessage) {
	ack := &slack.AckMessage{}
	if err := json.Unmarshal(event, ack); err != nil {
		smc.Debugln("RTM Error unmarshalling 'ack' event:", err)
		smc.Debugln(" -> Erroneous 'ack' event:", string(event))
		return
	}

	if ack.Ok {
		smc.IncomingEvents <- smc.externalEvent("ack", ack)
	} else if ack.RTMResponse.Error != nil {
		// As there is no documentation for RTM error-codes, this
		// identification of a rate-limit warning is very brittle.
		if ack.RTMResponse.Error.Code == -1 && ack.RTMResponse.Error.Msg == "slow down, too many messages..." {
			smc.IncomingEvents <- smc.internalEvent("ack_error", &slack.RateLimitEvent{})
		} else {
			smc.IncomingEvents <- smc.internalEvent("ack_error", &slack.AckErrorEvent{ack.Error})
		}
	} else {
		smc.IncomingEvents <- smc.internalEvent("ack_error", &slack.AckErrorEvent{fmt.Errorf("ack decode failure")})
	}
}

// handlePing handles an incoming 'PONG' message which should be in response to
// a previously sent 'PING' message. This is then used to compute the
// connection's latency.
func (smc *Client) handlePing(event json.RawMessage) {
	smc.resetDeadman()

	s := string(event)

	println("WebSocket ping message received:", s)

	//
	//latency := time.Since(time.Unix(p.Timestamp, 0))
	//smc.IncomingEvents <- smc.internalEvent("latency_report", &LatencyReport{Value: latency})
}

func (smc *Client) handleClose(code int, text string) {
	smc.killConnection(code == 200, errors.New(text))
}
