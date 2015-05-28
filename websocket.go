package slack

import (
	"log"
	"sync"
	"time"

	"golang.org/x/net/websocket"
)

// RTM represents a managed websocket connection. It also supports
// all the methods of the `Slack` type.
type RTM struct {
	mutex     sync.Mutex
	messageId int
	pings     map[int]time.Time

	// Connection life-cycle
	conn             *websocket.Conn
	IncomingEvents   chan SlackEvent
	outgoingMessages chan OutgoingMessage

	// Slack is the main API, embedded
	Client
	websocketURL string

	// UserDetails upon connection
	info *Info
}

// NewRTM returns a RTM, which provides a fully managed connection to
// Slack's websocket-based Real-Time Messaging protocol.
func newRTM(api *Client) *RTM {
	return &RTM{
		Client:         *api,
		pings:          make(map[int]time.Time),
		IncomingEvents: make(chan SlackEvent, 50),
	}
}

// Disconnect and wait, blocking until a successful disconnection.
func (rtm *RTM) Disconnect() error {
	log.Println("RTM::Disconnect not implemented!")
	return nil
}

// Reconnect, only makes sense if you've successfully disconnectd with Disconnect().
func (rtm *RTM) Reconnect() error {
	log.Println("RTM::Reconnect not implemented!")
	return nil
}

// GetInfo returns the info structure received when calling
// "startrtm", holding all channels, groups and other metadata needed
// to implement a full chat client. It will be non-nil after a call to
// StartRTM().
func (rtm *RTM) GetInfo() *Info {
	return rtm.info
}

// SendMessage submits a simple message through the websocket.  For
// more complicated messages, use `rtm.PostMessage` with a complete
// struct describing your attachments and all.
func (rtm *RTM) SendMessage(msg *OutgoingMessage) {
	if msg == nil {
		rtm.Debugln("Error: Attempted to SendMessage(nil)")
		return
	}

	rtm.outgoingMessages <- *msg
}
