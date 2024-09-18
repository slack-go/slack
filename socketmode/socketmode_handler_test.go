package socketmode

import (
	"log"
	"os"
	"reflect"
	"runtime"
	"testing"

	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
)

func init_SocketmodeHandler() *SocketmodeHandler {
	eventMap := make(map[EventType][]SocketmodeHandlerFunc)
	interactioneventMap := make(map[slack.InteractionType][]SocketmodeHandlerFunc)
	eventApiMap := make(map[slackevents.EventsAPIType][]SocketmodeHandlerFunc)
	interactionBlockActionEventMap := make(map[string]SocketmodeHandlerFunc)
	slashCommandMap := make(map[string]SocketmodeHandlerFunc)

	return &SocketmodeHandler{
		Client: &Client{
			log: log.New(os.Stderr, "slack-go/slack/socketmode", log.LstdFlags|log.Lshortfile),
		},
		EventMap:                       eventMap,
		EventApiMap:                    eventApiMap,
		InteractionEventMap:            interactioneventMap,
		InteractionBlockActionEventMap: interactionBlockActionEventMap,
		SlashCommandMap:                slashCommandMap,
	}
}

// The goal of this function is to catch the name of the function that is behing called
// This let us validate that the dispatcher did its job correctly
func testing_wrapper(ch chan<- string, f SocketmodeHandlerFunc) SocketmodeHandlerFunc {
	return SocketmodeHandlerFunc(func(e *Event, c *Client) {
		f(e, c)

		var name_f string

		// test with the name of the function we called
		v := reflect.ValueOf(f)
		if v.Kind() == reflect.Func {
			if rf := runtime.FuncForPC(v.Pointer()); rf != nil {
				name_f = rf.Name()
			}
		} else {
			name_f = v.String()
		}

		ch <- name_f
	})
}

func middleware_interaction(evt *Event, client *Client) {
	//do nothing
}

func middleware_interaction_block_action(evt *Event, client *Client) {
	//do nothing
}

func middleware_eventapi(evt *Event, client *Client) {
	//do nothing
}

func middleware(evt *Event, client *Client) {
	//do nothing
}

func defaultmiddleware(evt *Event, client *Client) {
	//do nothing
}

func middleware_slach_command(evt *Event, client *Client) {
	//do nothing
}

func TestSocketmodeHandler_Handle(t *testing.T) {
	type args struct {
		evt      Event
		evt_type EventType
	}
	tests := []struct {
		name string
		args args
		want string //what is the name of the function we want to be called
	}{
		{
			name: "Event Match registered function",
			args: args{
				evt: Event{
					Type: EventTypeConnecting,
				},
				evt_type: EventTypeConnecting,
			},
			want: "github.com/slack-go/slack/socketmode.middleware",
		}, {
			name: "Event do not registered function",
			args: args{
				evt: Event{
					Type: EventTypeConnected,
				},
				evt_type: EventTypeConnecting,
			},
			want: "github.com/slack-go/slack/socketmode.defaultmiddleware",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := init_SocketmodeHandler()

			c := make(chan string)

			r.Handle(tt.args.evt_type, testing_wrapper(c, middleware))
			r.HandleDefault(testing_wrapper(c, defaultmiddleware))

			r.dispatcher(tt.args.evt)

			got := <-c

			if got != tt.want {
				t.Fatalf("middleware was not called for EventTy(\"%v\"), got %v", tt.args.evt_type, got)
			}
		})
	}
}

func TestSocketmodeHandler_HandleInteraction(t *testing.T) {
	type args struct {
		evt      Event
		register func(*SocketmodeHandler, chan<- string)
	}
	tests := []struct {
		name string
		args args
		want string //what is the name of the function we want to be called
	}{
		{
			name: "Event Match registered function",
			args: args{
				evt: Event{
					Type: EventTypeInteractive,
					Data: slack.InteractionCallback{
						Type: slack.InteractionTypeBlockActions,
					},
				},
				register: func(r *SocketmodeHandler, c chan<- string) {
					r.HandleInteraction(slack.InteractionTypeBlockActions, testing_wrapper(c, middleware_interaction))
				},
			},
			want: "github.com/slack-go/slack/socketmode.middleware_interaction",
		}, {
			name: "Event do not Match any registered function",
			args: args{
				evt: Event{
					Type: EventTypeInteractive,
					Data: slack.InteractionCallback{
						Type: slack.InteractionTypeBlockActions,
					},
				},
				register: func(r *SocketmodeHandler, c chan<- string) {
					r.HandleInteraction(slack.InteractionTypeBlockSuggestion, testing_wrapper(c, middleware_interaction))
				},
			},
			want: "github.com/slack-go/slack/socketmode.defaultmiddleware",
		}, {
			name: "Event with invalid data is handled by default middleware",
			args: args{
				evt: Event{
					Type: EventTypeInteractive,
					Data: map[string]string{
						"brokendata": "test",
					},
				},
				register: func(r *SocketmodeHandler, c chan<- string) {
					r.HandleInteraction(slack.InteractionTypeBlockActions, testing_wrapper(c, middleware_interaction))
				},
			},
			want: "github.com/slack-go/slack/socketmode.defaultmiddleware",
		}, {
			name: "Event is handled as EventTypeInteractive",
			args: args{
				evt: Event{
					Type: EventTypeInteractive,
					Data: slack.InteractionCallback{
						Type: slack.InteractionTypeBlockActions,
					},
				},
				register: func(r *SocketmodeHandler, c chan<- string) {
					r.Handle(EventTypeInteractive, testing_wrapper(c, middleware))
				},
			},
			want: "github.com/slack-go/slack/socketmode.middleware",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := init_SocketmodeHandler()

			c := make(chan string)

			tt.args.register(r, c)
			r.HandleDefault(testing_wrapper(c, defaultmiddleware))

			r.dispatcher(tt.args.evt)

			got := <-c

			if got != tt.want {
				t.Fatalf("%s was not called for EventTy(\"%v\"), got %v", tt.want, tt.args.evt.Type, got)
			}
		})
	}
}

func TestSocketmodeHandler_HandleEvents(t *testing.T) {
	type args struct {
		evt      Event
		register func(*SocketmodeHandler, chan<- string)
	}
	tests := []struct {
		name string
		args args
		want string //what is the name of the function we want to be called
	}{
		{
			name: "Event Match registered function",
			args: args{
				evt: Event{
					Type: EventTypeEventsAPI,
					Data: slackevents.EventsAPIEvent{
						Type: "event_callback",
						InnerEvent: slackevents.EventsAPIInnerEvent{
							Type: string(slackevents.AppMention),
						},
					},
				},
				register: func(r *SocketmodeHandler, c chan<- string) {
					r.HandleEvents(slackevents.AppMention, testing_wrapper(c, middleware_eventapi))
				},
			},
			want: "github.com/slack-go/slack/socketmode.middleware_eventapi",
		}, {
			name: "Event do not Match any registered function",
			args: args{
				evt: Event{
					Type: EventTypeEventsAPI,
					Data: slackevents.EventsAPIEvent{
						Type: "event_callback",
						InnerEvent: slackevents.EventsAPIInnerEvent{
							Type: string(slackevents.MemberJoinedChannel),
						},
					},
				},
				register: func(r *SocketmodeHandler, c chan<- string) {
					r.HandleEvents(slackevents.AppMention, testing_wrapper(c, middleware_eventapi))
				},
			},
			want: "github.com/slack-go/slack/socketmode.defaultmiddleware",
		}, {
			name: "Event with invalid data is handled by default middleware",
			args: args{
				evt: Event{
					Type: EventTypeEventsAPI,
					Data: map[string]string{
						"brokendata": "test",
					},
				},
				register: func(r *SocketmodeHandler, c chan<- string) {
					r.HandleEvents(slackevents.AppMention, testing_wrapper(c, middleware_eventapi))
				},
			},
			want: "github.com/slack-go/slack/socketmode.defaultmiddleware",
		}, {
			name: "Event is handled as EventTypeInteractive",
			args: args{
				evt: Event{
					Type: EventTypeEventsAPI,
					Data: slackevents.EventsAPIEvent{
						Type: "event_callback",
						InnerEvent: slackevents.EventsAPIInnerEvent{
							Type: string(slackevents.MemberJoinedChannel),
						},
					},
				},
				register: func(r *SocketmodeHandler, c chan<- string) {
					r.Handle(EventTypeEventsAPI, testing_wrapper(c, middleware))
				},
			},
			want: "github.com/slack-go/slack/socketmode.middleware",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := init_SocketmodeHandler()

			c := make(chan string)

			tt.args.register(r, c)
			r.HandleDefault(testing_wrapper(c, defaultmiddleware))

			r.dispatcher(tt.args.evt)

			got := <-c

			if got != tt.want {
				t.Fatalf("%s was not called for EventTy(\"%v\"), got %v", tt.want, tt.args.evt.Type, got)
			}
		})
	}
}

func TestSocketmodeHandler_HandleInteractionBlockAction(t *testing.T) {
	type args struct {
		evt      Event
		register func(*SocketmodeHandler, chan<- string)
	}
	tests := []struct {
		name string
		args args
		want string //what is the name of the function we want to be called
	}{
		{
			name: "Event Match registered function",
			args: args{
				evt: Event{
					Type: EventTypeInteractive,
					Data: slack.InteractionCallback{
						Type: slack.InteractionTypeBlockActions,
						ActionCallback: slack.ActionCallbacks{
							BlockActions: []*slack.BlockAction{
								{
									ActionID: "add_note",
									Text: slack.TextBlockObject{
										Type: "plain_text",
										Text: "Add a Stickie",
									},
								},
							},
						},
					},
				},
				register: func(r *SocketmodeHandler, c chan<- string) {
					r.HandleInteractionBlockAction("add_note", testing_wrapper(c, middleware_interaction_block_action))
				},
			},
			want: "github.com/slack-go/slack/socketmode.middleware_interaction_block_action",
		}, {
			name: "Event do not Match any registered function",
			args: args{
				evt: Event{
					Type: EventTypeInteractive,
					Data: slack.InteractionCallback{
						Type: slack.InteractionTypeBlockActions,
					},
				},
				register: func(r *SocketmodeHandler, c chan<- string) {
					r.HandleInteractionBlockAction("add_note", testing_wrapper(c, middleware_interaction_block_action))
				},
			},
			want: "github.com/slack-go/slack/socketmode.defaultmiddleware",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := init_SocketmodeHandler()

			c := make(chan string)

			tt.args.register(r, c)
			r.HandleDefault(testing_wrapper(c, defaultmiddleware))

			r.dispatcher(tt.args.evt)

			got := <-c

			if got != tt.want {
				t.Fatalf("%s was not called for EventTy(\"%v\"), got %v", tt.want, tt.args.evt.Type, got)
			}
		})
	}
}

func TestSocketmodeHandler_HandleSlashCommand(t *testing.T) {
	type args struct {
		evt      Event
		register func(*SocketmodeHandler, chan<- string)
	}
	tests := []struct {
		name string
		args args
		want string //what is the name of the function we want to be called
	}{
		{
			name: "Event Match registered function",
			args: args{
				evt: Event{
					Type: EventTypeSlashCommand,
					Data: slack.SlashCommand{
						Command: "/rocket",
						Text:    "key=value",
					},
				},
				register: func(r *SocketmodeHandler, c chan<- string) {
					r.HandleSlashCommand("/rocket", testing_wrapper(c, middleware_slach_command))
				},
			},
			want: "github.com/slack-go/slack/socketmode.middleware_slach_command",
		}, {
			name: "Event do not Match any registered function",
			args: args{
				evt: Event{
					Type: EventTypeSlashCommand,
					Data: slack.SlashCommand{
						Command: "/broken_rocket",
						Text:    "key=value",
					},
				},
				register: func(r *SocketmodeHandler, c chan<- string) {
					r.HandleSlashCommand("/rocket", testing_wrapper(c, middleware_slach_command))
				},
			},
			want: "github.com/slack-go/slack/socketmode.defaultmiddleware",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := init_SocketmodeHandler()

			c := make(chan string)

			tt.args.register(r, c)
			r.HandleDefault(testing_wrapper(c, defaultmiddleware))

			r.dispatcher(tt.args.evt)

			got := <-c

			if got != tt.want {
				t.Fatalf("%s was not called for EventTy(\"%v\"), got %v", tt.want, tt.args.evt.Type, got)
			}
		})
	}
}

func TestSocketmodeHandler_Handle_errors(t *testing.T) {
	type args struct {
		register func(*SocketmodeHandler, chan<- string)
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Attempt to register empty command",
			args: args{
				register: func(r *SocketmodeHandler, c chan<- string) {
					r.HandleSlashCommand("", testing_wrapper(c, middleware_slach_command))
				},
			},
		}, {
			name: "Attempt to register nil handler",
			args: args{
				register: func(r *SocketmodeHandler, c chan<- string) {
					r.HandleSlashCommand("/command", nil)
				},
			},
		}, {
			name: "Attempt to register duplicate command",
			args: args{
				register: func(r *SocketmodeHandler, c chan<- string) {
					r.HandleSlashCommand("/command", testing_wrapper(c, middleware_slach_command))
					r.HandleSlashCommand("/command", testing_wrapper(c, middleware_slach_command))
				},
			},
		}, {
			name: "Attempt to register empty Block ActionID",
			args: args{
				register: func(r *SocketmodeHandler, c chan<- string) {
					r.HandleInteractionBlockAction("", testing_wrapper(c, middleware_interaction_block_action))
				},
			},
		}, {
			name: "Attempt to register nil handler",
			args: args{
				register: func(r *SocketmodeHandler, c chan<- string) {
					r.HandleInteractionBlockAction("action_id", nil)
				},
			},
		}, {
			name: "Attempt to register duplicate Block ActionID",
			args: args{
				register: func(r *SocketmodeHandler, c chan<- string) {
					r.HandleInteractionBlockAction("action_id", testing_wrapper(c, middleware_interaction_block_action))
					r.HandleInteractionBlockAction("action_id", testing_wrapper(c, middleware_interaction_block_action))
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := init_SocketmodeHandler()

			c := make(chan string)

			defer func() { recover() }()

			tt.args.register(r, c)

			t.Errorf("should have panicked")

		})
	}
}
