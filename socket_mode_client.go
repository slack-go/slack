package slack

import (
	"encoding/json"
	"net/url"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type SocketModeConnectedEvent struct {
	ConnectionCount int // 1 = first time, 2 = second time
	Info            *SocketModeConnection
}

// SocketModeConnection contains various details about the SocketMode connection.
// It is returned by an "apps.connections.open" API call.
type SocketModeConnection struct {
	URL  string                 `json:"url,omitempty"`
	Data map[string]interface{} `json:"-"`
}

// SocketModeMessage maps to each message received via the WebSocket connection of SocketMode
type SocketModeMessage struct {
	Type       string                   `json:"type"`
	Reason     string                   `json:"reason"`
	Payload    SocketModeMessagePayload `json:"payload"`
	EnvelopeID string                   `json:"envelope_id"`
}

type SocketModeMessagePayload struct {
	Event json.RawMessage `json:"´event"`
}

// SocketModeEvent is the event sent to the consumer of SocketModeClient
type SocketModeEvent struct {
	Type string
	Data interface{}

	// Message is the json-decoded raw WebSocket message that is received over the Slack SocketMode
	// WebSocket connection.
	Message *SocketModeMessage
}

// SocketModeClient allows allows programs to communicate with the
// [Events API](https://api.slack.com/events-api) over WebSocket.
//
// The implementation is highly inspired by https://www.npmjs.com/package/@slack/socket-mode,
// but the structure and the design has been adapted as much as possible to that of our RTM client for consistency
// within the library.
//
// You can instantiate the socket mode client with
// Client's NewSocketModeClient() or NewSocketModeClientWithOptions(*SocketModeClientOptions)
type SocketModeClient struct {
	// Client is the main API, embedded
	Client

	idGen        IDGenerator
	pingInterval time.Duration
	pingDeadman  *time.Timer

	// Connection life-cycle
	conn             *websocket.Conn
	IncomingEvents   chan SocketModeEvent
	outgoingMessages chan OutgoingMessage
	killChannel      chan bool
	disconnected     chan struct{}
	disconnectedm    *sync.Once

	// UserDetails upon connection
	info *SocketModeConnection

	// dialer is a gorilla/websocket Dialer. If nil, use the default
	// Dialer.
	dialer *websocket.Dialer

	// mu is mutex used to prevent RTM connection race conditions
	mu *sync.Mutex

	// connParams is a map of flags for connection parameters.
	connParams url.Values
}

// signal that we are disconnected by closing the channel.
// protect it with a mutex to ensure it only happens once.
func (c *SocketModeClient) disconnect() {
	c.disconnectedm.Do(func() {
		close(c.disconnected)
	})
}

// Disconnect and wait, blocking until a successful disconnection.
func (c *SocketModeClient) Disconnect() error {
	// always push into the kill channel when invoked,
	// this lets the ManagedConnection() function properly clean up.
	// if the buffer is full then just continue on.
	select {
	case c.killChannel <- true:
		return nil
	case <-c.disconnected:
		return ErrAlreadyDisconnected
	}
}

// GetInfo returns the info structure received when calling
// "startrtm", holding metadata needed to implement a full
// chat client. It will be non-nil after a call to StartRTM().
func (c *SocketModeClient) GetInfo() *SocketModeConnection {
	return c.info
}

func (c *SocketModeClient) resetDeadman() {
	c.pingDeadman.Reset(deadmanDuration(c.pingInterval))
}
