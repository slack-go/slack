package socketmode

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/gorilla/websocket"

	"github.com/slack-go/slack"
)

// EventType is the type of events that are emitted by scoketmode.Client.
// You receive and handle those events from a socketmode.Client.Events channel.
// Those event types does not necessarily match 1:1 to those of Slack Events API events.
type EventType string

const (
	// The following request types are the types of requests sent from Slack via Socket Mode WebSocket connection
	// and handled internally by the socketmode.Client.
	// The consumer of socketmode.Client will never see it.

	RequestTypeHello         = "hello"
	RequestTypeEventsAPI     = "events_api"
	RequestTypeDisconnect    = "disconnect"
	RequestTypeSlashCommands = "slash_commands"
	RequestTypeInteractive   = "interactive"

	// The following event types are for events emitted by socketmode.Client itself and
	// does not originate from Slack.
	EventTypeConnecting       = EventType("connecting")
	EventTypeInvalidAuth      = EventType("invalid_auth")
	EventTypeConnectionError  = EventType("connection_error")
	EventTypeConnected        = EventType("connected")
	EventTypeIncomingError    = EventType("incoming_error")
	EventTypeErrorWriteFailed = EventType("write_error")
	EventTypeErrorBadMessage  = EventType("error_bad_message")

	//
	// The following event types are guaranteed to not change unless Slack changes
	//

	EventTypeHello        = EventType("hello")
	EventTypeDisconnect   = EventType("disconnect")
	EventTypeEventsAPI    = EventType("events_api")
	EventTypeInteractive  = EventType("interactive")
	EventTypeSlashCommand = EventType("slash_commands")

	websocketDefaultTimeout = 10 * time.Second
	defaultMaxPingInterval  = 30 * time.Second
	defaultHandshakeTimeout = 45 * time.Second

	// maxResponseSize is Slack's server-side limit for Socket Mode WebSocket response
	// messages (20KB). Responses at or above this size on the wire are silently dropped
	// by Slack — the WebSocket write succeeds but Slack ignores the payload. Empirically
	// verified using both this library and Slack's official Node SDK
	// (@slack/socket-mode).
	//
	// Note: gorilla/websocket's WriteJSON adds a trailing newline (via
	// json.Encoder.Encode), so the on-wire size is json.Marshal(response) + 1. The
	// pre-flight check in SendCtx uses >= maxResponseSize on the json.Marshal output,
	// which is equivalent to > maxResponseSize on the wire minus the newline —
	// conservative by exactly 1 byte.
	maxResponseSize = 20 * 1024 // 20480 bytes

	// defaultWriteBufferSize is the WebSocket write buffer size used when dialing Slack's
	// Socket Mode endpoint. gorilla/websocket's default is 4096 bytes; messages exceeding
	// the buffer are split into WebSocket continuation frames (I think). Slack's server
	// does not reassemble continuation frames, causing silent message drops (this is my
	// hypothesis from testing). We set this above Slack's 20KB response limit to ensure
	// all valid messages are sent as single frames.
	defaultWriteBufferSize = 32 * 1024
)

// Open calls the "apps.connections.open" endpoint and returns the provided URL and the
// full Info block.
//
// To have a fully managed Websocket connection, use `New`, and call `Run()` on it.
func (smc *Client) Open() (info *slack.SocketModeConnection, websocketURL string, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), websocketDefaultTimeout)
	defer cancel()

	return smc.StartSocketModeContext(ctx)
}

// OpenContext calls the "apps.connections.open" endpoint and returns the provided URL and
// the full Info block.
//
// To have a fully managed Websocket connection, use `New`, and call `Run()` on it.
func (smc *Client) OpenContext(ctx context.Context) (info *slack.SocketModeConnection, websocketURL string, err error) {
	return smc.StartSocketModeContext(ctx)
}

// Option options for the managed Client.
type Option func(client *Client)

// OptionDialer takes a gorilla websocket Dialer and uses it as the
// Dialer when opening the websocket for the Socket Mode connection.
func OptionDialer(d *websocket.Dialer) Option {
	return func(smc *Client) {
		smc.dialer = d
	}
}

// OptionPingInterval determines how often we expect Slack to deliver WebSocket ping to us.
// If no ping is delivered to us within this interval after the last ping, we assumes the WebSocket connection
// is dead and needs to be reconnected.
func OptionPingInterval(d time.Duration) Option {
	return func(smc *Client) {
		smc.maxPingInterval = d
	}
}

// OptionDebug enable debugging for the client
func OptionDebug(b bool) func(*Client) {
	return func(c *Client) {
		c.debug = b
	}
}

// OptionLog set logging for client.
func OptionLog(l logger) func(*Client) {
	return func(c *Client) {
		c.log = internalLog{logger: l}
	}
}

// New returns a Socket Mode client which provides a fully managed connection to
// Slack's Websocket-based Socket Mode.
func New(api *slack.Client, options ...Option) *Client {
	result := &Client{
		Client:              *api,
		Events:              make(chan Event, 50),
		socketModeResponses: make(chan *Response, 20),
		maxPingInterval:     defaultMaxPingInterval,
		log:                 log.New(os.Stderr, "slack-go/slack/socketmode", log.LstdFlags|log.Lshortfile),
	}

	for _, opt := range options {
		opt(result)
	}

	return result
}

// sendEvent safely sends an event into the Clients Events channel
// and blocks until buffer space is had, or the context is canceled.
// This prevents deadlocking in the event that Events buffer is full,
// other goroutines are waiting, and/or timing allows receivers to exit
// before all senders are finished.
func (smc *Client) sendEvent(ctx context.Context, event Event) {
	select {
	case smc.Events <- event:
	case <-ctx.Done():
	}
}
