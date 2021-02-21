package socketmode

import (
	"log"

	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
)

type SocketmodeHandler struct {
	Socket *Client

	EventMap map[EventType][]SocketmodeHandlerFunc
}

type SocketmodeHandlerFunc func(Event)

func NewsSocketmodeHandler(socket *Client) *SocketmodeHandler {
	eventMap := make(map[EventType][]SocketmodeHandlerFunc)

	return &SocketmodeHandler{
		Socket:   socket,
		EventMap: eventMap,
	}
}

func (r *SocketmodeHandler) Handle(et EventType, f func(Event)) {
	r.EventMap[et] = append(r.EventMap[et], f)
}

// RunSlackEventLoop receives the event via the socket
// It receives events from Slack and each is handled as needed
func (r *SocketmodeHandler) RunEventLoop() {

	go r.runEventLoop()

	r.Socket.Run()
}

func (r *SocketmodeHandler) runEventLoop() {
	for evt := range r.Socket.Events {
		if handlers, ok := r.EventMap[evt.Type]; ok {
			// If we registered an event
			for _, f := range handlers {
				go f(evt)
			}
		} else {
			// We need to explicitely subscribe to event in the Application Dashboard
			// So every event sould be handle otherwise this is an error
			log.Printf("Unexpected event type received: %v\n", evt.Type)
		}

	}
}

func (r *SocketmodeHandler) EventTypeEventsAPIHandler(evt *Event) {
	eventsAPIEvent, ok := evt.Data.(slackevents.EventsAPIEvent)
	if !ok {
		log.Println("Unable to handle recieved EventsAPIEvent")
		return
	}

	// We need to confirm that we recived that event
	r.Socket.Ack(*evt.Request)

	switch eventsAPIEvent.Type {
	case slackevents.CallbackEvent:
		innerEvent := eventsAPIEvent.InnerEvent
		switch ev := innerEvent.Data.(type) {
		case *slackevents.AppMentionEvent:

		case *slackevents.MessageEvent:

		case *slackevents.MemberJoinedChannelEvent:

		default:
			log.Printf("unsupported CallbackEvent event received: %T", ev)
		}
	default:
		log.Printf("unsupported eventsAPIEvent event received: %s", eventsAPIEvent.Type)
	}
}

func (r *SocketmodeHandler) EventTypeInteractiveHandler(evt *Event) {
	callback, ok := evt.Data.(slack.InteractionCallback)
	if !ok {
		return
	}

	var payload interface{}

	switch callback.Type {
	case slack.InteractionTypeBlockActions:
		// See https://api.slack.com/apis/connections/socket-implement#button
	case slack.InteractionTypeShortcut:
		// See https://api.slack.com/interactivity/shortcuts
	case slack.InteractionTypeViewSubmission:
		// See https://api.slack.com/apis/connections/socket-implement#modal
	case slack.InteractionTypeDialogSubmission:

	default:
		log.Printf("unsupported callbackEvent event received: %s", callback.Type)
	}

	r.Socket.Ack(*evt.Request, payload)
}

func (r *SocketmodeHandler) EventTypeSlashCommandHandler(evt *Event) {
	cmd, ok := evt.Data.(slack.SlashCommand)
	if !ok {
		return
	}

	log.Printf("%s", cmd)
}
