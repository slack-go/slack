package slacksocketmode

import (
	"context"
	"github.com/slack-go/slack"
	"net/url"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

const (
	RequestTypeHello      = "hello"
	RequestTypeEventsAPI  = "events_api"
	RequestTypeDisconnect = "disconnect"

	EventTypeConnected        = "connected"
	EventTypeErrorWriteFailed = "write_error"

	websocketDefaultTimeout = 10 * time.Second
	defaultMaxPingInterval  = 30 * time.Second
)

// Open calls the "apps.connection.open" endpoint and returns the provided URL and the full Info block.
//
// To have a fully managed Websocket connection, use `New`, and call `Run()` on it.
func (smc *Client) Open() (info *slack.SocketModeConnection, websocketURL string, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), websocketDefaultTimeout)
	defer cancel()

	return smc.StartSocketModeContext(ctx)
}

// Option options for the managed Client.
type Option func(client *Client)

// OptionDialer takes a gorilla websocket Dialer and uses it as the
// Dialer when opening the websocket for the RTM connection.
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

// NewRTM returns a RTM, which provides a fully managed connection to
// Slack's websocket-based Real-Time Messaging protocol.
func New(api *slack.Client, options ...Option) *Client {
	result := &Client{
		Client:              *api,
		IncomingEvents:      make(chan ClientEvent, 50),
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
