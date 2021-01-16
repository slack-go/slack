package slacksocketmode

import (
	"encoding/json"
	"github.com/slack-go/slack"
	"net/url"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type ConnectedEvent struct {
	ConnectionCount int // 1 = first time, 2 = second time
	Info            *slack.SocketModeConnection
}

type DebugInfo struct {
	// Host is the name of the host name on the Slack end, that can be something like `applink-7fc4fdbb64-4x5xq`
	Host string `json:"host"`

	// `hello` type only
	BuildNumber               int `json:"build_number"`
	ApproximateConnectionTime int `json:"approximate_connection_time"`
}

type ConnectionInfo struct {
	AppID string `json:"app_id"`
}

type SocketModeMessagePayload struct {
	Event json.RawMessage `json:"´event"`
}

// ClientEvent is the event sent to the consumer of Client
type ClientEvent struct {
	Type string
	Data interface{}

	// Request is the json-decoded raw WebSocket message that is received via the Slack Socket Mode
	// WebSocket connection.
	Request *Request
}

// Client allows allows programs to communicate with the
// [Events API](https://api.slack.com/events-api) over WebSocket.
//
// The implementation is highly inspired by https://www.npmjs.com/package/@slack/socket-mode,
// but the structure and the design has been adapted as much as possible to that of our RTM client for consistency
// within the library.
//
// You can instantiate the socket mode client with
// Client's New() or NewSocketModeClientWithOptions(*SocketModeClientOptions)
type Client struct {
	// Client is the main API, embedded
	slack.Client

	idGen        slack.IDGenerator
	pingInterval time.Duration

	// pingDeadman must be intiailized in New()
	pingDeadman *time.Timer

	// Connection life-cycle
	conn                *websocket.Conn
	IncomingEvents      chan ClientEvent
	socketModeResponses chan *Response
	killChannel         chan bool
	disconnected        chan struct{}
	disconnectedm       *sync.Once

	// UserDetails upon connection
	info *slack.SocketModeConnection

	// dialer is a gorilla/websocket Dialer. If nil, use the default
	// Dialer.
	dialer *websocket.Dialer

	// mu is mutex used to prevent RTM connection race conditions
	mu *sync.Mutex

	wsWriteMu *sync.Mutex

	// connParams is a map of flags for connection parameters.
	connParams url.Values
}

// signal that we are disconnected by closing the channel.
// protect it with a mutex to ensure it only happens once.
func (smc *Client) disconnect() {
	smc.disconnectedm.Do(func() {
		close(smc.disconnected)
	})
}

// Disconnect and wait, blocking until a successful disconnection.
func (smc *Client) Disconnect() error {
	// always push into the kill channel when invoked,
	// this lets the ManagedConnection() function properly clean up.
	// if the buffer is full then just continue on.
	select {
	case smc.killChannel <- true:
		return nil
	case <-smc.disconnected:
		return slack.ErrAlreadyDisconnected
	}
}

// GetInfo returns the info structure received when calling
// "startrtm", holding metadata needed to implement a full
// chat client. It will be non-nil after a call to StartRTM().
func (smc *Client) GetInfo() *slack.SocketModeConnection {
	return smc.info
}

func (smc *Client) resetDeadman() {
	smc.pingDeadman.Reset(deadmanDuration(smc.pingInterval))
}
