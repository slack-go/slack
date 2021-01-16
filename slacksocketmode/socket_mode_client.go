package slacksocketmode

import (
	"encoding/json"
	"github.com/slack-go/slack"
	"net/url"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type SocketModeConnectedEvent struct {
	ConnectionCount int // 1 = first time, 2 = second time
	Info            *slack.SocketModeConnection
}

// Request maps to the content of each WebSocket message received via a Socket Mode WebSocket connection
//
// We call this a "request" rather than e.g. a WebSocket message or an Socket Mode "event" following python-slack-sdk:
//
//   https://github.com/slackapi/python-slack-sdk/blob/3f1c4c6e27bf7ee8af57699b2543e6eb7848bcf9/slack_sdk/socket_mode/request.py#L6
//
// We know that node-slack-sdk calls it an "event", that makes it hard for us to distinguish our client's own event
// that wraps both internal events and Socket Mode "events", vs node-slack-sdk's is for the latter only.
//
// https://github.com/slackapi/node-slack-sdk/blob/main/packages/socket-mode/src/SocketModeClient.ts#L537
type Request struct {
	Type string `json:"type"`

	// `hello` type only
	NumConnections int            `json:"num_connections"`
	ConnectionInfo ConnectionInfo `json:"connection_info"`

	// `disconnect` type only

	// Reason can be "warning" or else
	Reason string `json:"reason"`

	// `hello` and `disconnect` types only
	DebugInfo DebugInfo `json:"debug_info"`

	// `events_api` type only
	EnvelopeID             string          `json:"envelope_id"`
	Payload                json.RawMessage `json:"payload"`
	AcceptsResponsePayload bool            `json:"accepts_response_payload"`
	RetryAttempt           int             `json:"retry_attempt"`
	RetryReason            string          `json:"retry_reason"`
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
	pingDeadman  *time.Timer

	// Connection life-cycle
	conn             *websocket.Conn
	IncomingEvents   chan ClientEvent
	outgoingMessages chan slack.OutgoingMessage
	killChannel      chan bool
	disconnected     chan struct{}
	disconnectedm    *sync.Once

	// UserDetails upon connection
	info *slack.SocketModeConnection

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
