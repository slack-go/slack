package socketmode

import (
	"fmt"
	"log"

	"github.com/slack-go/slack"
)

type SocketmodeHandler struct {
	Client *Client

	EventMap            map[EventType][]SocketmodeHandlerFunc
	InteractionEventMap map[slack.InteractionType][]SocketmodeHandlerFunc
}

type SocketmodeHandlerFunc func(*Event, *Client)

func NewsSocketmodeHandler(client *Client) *SocketmodeHandler {
	eventMap := make(map[EventType][]SocketmodeHandlerFunc)
	interactioneventMap := make(map[slack.InteractionType][]SocketmodeHandlerFunc)

	return &SocketmodeHandler{
		Client:              client,
		EventMap:            eventMap,
		InteractionEventMap: interactioneventMap,
	}
}

func (r *SocketmodeHandler) Handle(et EventType, f SocketmodeHandlerFunc) {
	r.EventMap[et] = append(r.EventMap[et], f)
}

func (r *SocketmodeHandler) HandleInteraction(et slack.InteractionType, f SocketmodeHandlerFunc) {
	r.InteractionEventMap[et] = append(r.InteractionEventMap[et], f)
}

// RunSlackEventLoop receives the event via the socket
// It receives events from Slack and each is handled as needed
func (r *SocketmodeHandler) RunEventLoop() {

	go r.runEventLoop()

	r.Client.Run()
}

func (r *SocketmodeHandler) runEventLoop() {
	for evt := range r.Client.Events {

		// Some eventType can be further decomposed
		switch evt.Type {
		case EventTypeInteractive:
			go r.Interaction(&evt)
			// case EventTypeEventsAPI:
		}

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

func (r *SocketmodeHandler) Interaction(evt *Event) {
	interaction, ok := evt.Data.(slack.InteractionCallback)
	if !ok {
		fmt.Printf("Ignored %+v\n", evt)
		return
	}

	if handlers, ok := r.InteractionEventMap[interaction.Type]; ok {
		// If we registered an event
		for _, f := range handlers {
			go f(evt, r.Client)
		}
	} else {
		// We need to explicitely subscribe to event in the Application Dashboard
		// So every event sould be handle otherwise this is an error
		log.Printf("Unexpected event type received: %v\n", evt.Type)
	}

}
