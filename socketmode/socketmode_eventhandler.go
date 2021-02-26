package socketmode

import (
	"fmt"
	"log"

	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
)

type SocketmodeHandler struct {
	Client *Client

	EventMap            map[EventType][]SocketmodeHandlerFunc
	InteractionEventMap map[slack.InteractionType][]SocketmodeHandlerFunc
	EventApiMap         map[slackevents.EventAPIType][]SocketmodeHandlerFunc

	Default SocketmodeHandlerFunc
}

// Handler have access to the event and socketmode client
type SocketmodeHandlerFunc func(*Event, *Client)

// Middleware accept SocketmodeHandlerFunc, and return SocketmodeHandlerFunc
type SocketmodeMiddlewareFunc func(SocketmodeHandlerFunc) SocketmodeHandlerFunc

func NewsSocketmodeHandler(client *Client) *SocketmodeHandler {
	eventMap := make(map[EventType][]SocketmodeHandlerFunc)
	interactioneventMap := make(map[slack.InteractionType][]SocketmodeHandlerFunc)
	eventApiMap := make(map[slackevents.EventAPIType][]SocketmodeHandlerFunc)

	return &SocketmodeHandler{
		Client:              client,
		EventMap:            eventMap,
		EventApiMap:         eventApiMap,
		InteractionEventMap: interactioneventMap,
		Default: func(e *Event, c *Client) {
			log.Printf("Unexpected event type received: %v\n", e.Type)
		},
	}
}

// Register the middleare funtion to use to handle an Event (from socketmode)
func (r *SocketmodeHandler) Handle(et EventType, f SocketmodeHandlerFunc) {
	r.EventMap[et] = append(r.EventMap[et], f)
}

// Register the middleare funtion to use to handle an Interaction
func (r *SocketmodeHandler) HandleInteraction(et slack.InteractionType, f SocketmodeHandlerFunc) {
	r.InteractionEventMap[et] = append(r.InteractionEventMap[et], f)
}

// Register the middleare funtion to use to handle an Event (from slackevents)
func (r *SocketmodeHandler) HandleEventsAPI(et slackevents.EventAPIType, f SocketmodeHandlerFunc) {
	r.EventApiMap[et] = append(r.EventApiMap[et], f)
}

// Register the middleare funtion to use as a last resort
func (r *SocketmodeHandler) HandleDefault(f SocketmodeHandlerFunc) {
	r.Default = f
}

// RunSlackEventLoop receives the event via the socket
func (r *SocketmodeHandler) RunEventLoop() {

	go r.runEventLoop()

	r.Client.Run()
}

func (r *SocketmodeHandler) runEventLoop() {
	for evt := range r.Client.Events {
		r.dispatcher(evt)
	}
}

func (r *SocketmodeHandler) dispatcher(evt Event) {
	var ishandled bool

	// Some eventType can be further decomposed
	switch evt.Type {
	case EventTypeInteractive:
		ishandled = r.interactionDispatcher(&evt)
	case EventTypeEventsAPI:
		ishandled = r.eventAPIDispatcher(&evt)
	default:
		ishandled = r.socketmodeDispatcher(&evt)
	}

	if !ishandled {
		go r.Default(&evt, r.Client)
	}
}

// Dispatch socketmode events to the registered middleware
func (r *SocketmodeHandler) socketmodeDispatcher(evt *Event) bool {
	if handlers, ok := r.EventMap[evt.Type]; ok {
		// If we registered an event
		for _, f := range handlers {
			go f(evt, r.Client)
		}

		return true
	}

	return false
}

// Dispatch interactions to the registered middleware
func (r *SocketmodeHandler) interactionDispatcher(evt *Event) bool {
	interaction, ok := evt.Data.(slack.InteractionCallback)
	if !ok {
		fmt.Printf("Ignored %+v\n", evt)
		return false
	}

	if handlers, ok := r.InteractionEventMap[interaction.Type]; ok {
		// If we registered an event
		for _, f := range handlers {
			go f(evt, r.Client)
		}

		return true
	} else if handlers, ok := r.EventMap[evt.Type]; ok {
		// fallback it this event is not handle by a more specific handler

		for _, f := range handlers {
			go f(evt, r.Client)
		}

		return true
	}

	return false
}

// Dispatch eventAPI events to the registered middleware
func (r *SocketmodeHandler) eventAPIDispatcher(evt *Event) bool {
	eventsAPIEvent, ok := evt.Data.(slackevents.EventsAPIEvent)
	if !ok {
		fmt.Printf("Ignored %+v\n", evt)
		return false
	}

	innerEventType := slackevents.EventAPIType(eventsAPIEvent.InnerEvent.Type)

	if handlers, ok := r.EventApiMap[innerEventType]; ok {
		// If we registered an event
		for _, f := range handlers {
			go f(evt, r.Client)
		}

		return true
	} else if handlers, ok := r.EventMap[evt.Type]; ok {
		// fallback it this event is not handle by a more specific handler
		for _, f := range handlers {
			go f(evt, r.Client)
		}

		return true
	}

	return false
}
