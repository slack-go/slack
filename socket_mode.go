package slack

import (
	"context"
	"net/url"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

const (
	socketModeEventTypeHello      = "hello"
)

// StartSocketMode calls the "rtm.start" endpoint and returns the provided URL and the full Info block.
//
// To have a fully managed Websocket connection, use `NewRTM`, and call `ManageConnection()` on it.
func (api *Client) StartSocketMode() (info *SocketModeConnection, websocketURL string, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), websocketDefaultTimeout)
	defer cancel()

	return api.StartSocketModeContext(ctx)
}

type openResponseFull struct {
	SlackResponse
	SocketModeConnection
}

// StartSocketModeContext calls the "apps.connections.open" endpoint and returns the provided URL and the full Info block with a custom context.
//
// To have a fully managed Websocket connection, use `NewRTM`, and call `ManageConnection()` on it.
func (api *Client) StartSocketModeContext(ctx context.Context) (info *SocketModeConnection, websocketURL string, err error) {
	response := &openResponseFull{}
	err = api.postMethod(ctx, "apps.connections.open", url.Values{"token": {api.token}}, response)
	if err != nil {
		return nil, "", err
	}

	api.Debugln("Using URL:", response.SocketModeConnection.URL)
	return &response.SocketModeConnection, response.SocketModeConnection.URL, response.Err()
}

// SocketModeOption options for the managed SocketModeClient.
type SocketModeOption func(client *SocketModeClient)

// SocketModeOptionDialer takes a gorilla websocket Dialer and uses it as the
// Dialer when opening the websocket for the RTM connection.
func SocketModeOptionDialer(d *websocket.Dialer) SocketModeOption {
	return func(rtm *SocketModeClient) {
		rtm.dialer = d
	}
}

// SocketModeOptionPingInterval determines how often to deliver a ping message to slack.
func SocketModeOptionPingInterval(d time.Duration) SocketModeOption {
	return func(rtm *SocketModeClient) {
		rtm.pingInterval = d
		rtm.resetDeadman()
	}
}

// SocketModeOptionConnParams installs parameters to embed into the connection URL.
func SocketModeOptionConnParams(connParams url.Values) SocketModeOption {
	return func(rtm *SocketModeClient) {
		rtm.connParams = connParams
	}
}

// NewRTM returns a RTM, which provides a fully managed connection to
// Slack's websocket-based Real-Time Messaging protocol.
func (api *Client) NewSocketModeClient(options ...SocketModeOption) *SocketModeClient {
	result := &SocketModeClient{
		Client:           *api,
		IncomingEvents:   make(chan SocketModeEvent, 50),
		outgoingMessages: make(chan OutgoingMessage, 20),
		pingInterval:     defaultPingInterval,
		pingDeadman:      time.NewTimer(deadmanDuration(defaultPingInterval)),
		killChannel:      make(chan bool),
		disconnected:     make(chan struct{}),
		disconnectedm:    &sync.Once{},
		idGen:            NewSafeID(1),
		mu:               &sync.Mutex{},
	}

	for _, opt := range options {
		opt(result)
	}

	return result
}
