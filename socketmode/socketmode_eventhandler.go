package socketmode

import (
	"log"
)

type SocketmodeHandler struct {
	Client *Client

	EventMap map[EventType][]SocketmodeHandlerFunc
}

type SocketmodeHandlerFunc func(*Event, *Client)

func NewsSocketmodeHandler(client *Client) *SocketmodeHandler {
	eventMap := make(map[EventType][]SocketmodeHandlerFunc)

	return &SocketmodeHandler{
		Client:   client,
		EventMap: eventMap,
	}
}

func (r *SocketmodeHandler) Handle(et EventType, f SocketmodeHandlerFunc) {
	r.EventMap[et] = append(r.EventMap[et], f)
}

// RunSlackEventLoop receives the event via the socket
// It receives events from Slack and each is handled as needed
func (r *SocketmodeHandler) RunEventLoop() {

	go r.runEventLoop()

	r.Client.Run()
}

func (r *SocketmodeHandler) runEventLoop() {
	for evt := range r.Client.Events {
		if handlers, ok := r.EventMap[evt.Type]; ok {
			// If we registered an event
			for _, f := range handlers {
				go f(&evt, r.Client)
			}
		} else {
			// We need to explicitely subscribe to event in the Application Dashboard
			// So every event sould be handle otherwise this is an error
			log.Printf("Unexpected event type received: %v\n", evt.Type)
		}

	}
}
