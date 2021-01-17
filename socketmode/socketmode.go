package socketmode

import (
	"context"
	"net/url"
	"sync"
	"time"

	"github.com/slack-go/slack"

	"github.com/gorilla/websocket"
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
	EventTypeConnected        = EventType("connected")
	EventTypeErrorWriteFailed = EventType("write_error")

	//
	// The following event types are guaranteed to not change unless Slack changes
	//

	EventTypeEventsAPI    = EventType("events_api")
	EventTypeInteractive  = EventType("interactive")
	EventTypeSlashCommand = EventType("slash_command")

	websocketDefaultTimeout = 10 * time.Second
	defaultMaxPingInterval  = 30 * time.Second
)

// Open calls the "apps.connection.open" endpoint and returns the provided URL and the full Info block.
//
// To have a fully managed Websocket connection, use `New`, and call `Run()` on it.
func (smc *Client) Open() (info *slack.SocketModeConnection, websocketURL string, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), websocketDefaultTimeout)
	defer cancel()

	return smc.apiClient.StartSocketModeContext(ctx)
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
	return func(rtm *Client) {
		rtm.pingInterval = d
	}
}

// OptionConnParams installs parameters to embed into the connection URL.
func OptionConnParams(connParams url.Values) Option {
	return func(smc *Client) {
		smc.connParams = connParams
	}
}

// New returns a Socket Mode client which provides a fully managed connection to
// Slack's Websocket-based Socket Mode.
func New(api *slack.Client, options ...Option) *Client {
	result := &Client{
		apiClient:           *api,
		Events:              make(chan ClientEvent, 50),
		socketModeResponses: make(chan *Response, 20),
		pingInterval:        defaultMaxPingInterval,
		killChannel:         make(chan bool),
		disconnected:        make(chan struct{}),
		disconnectedm:       &sync.Once{},
		idGen:               slack.NewSafeID(1),
		mu:                  &sync.Mutex{},
		wsWriteMu:           &sync.Mutex{},
	}

	for _, opt := range options {
		opt(result)
	}

	result.pingDeadman = time.NewTimer(deadmanDuration(result.pingInterval))

	return result
}
