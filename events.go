package slack

import (
	"encoding/json"
	"fmt"
	"reflect"
)

// EventsAPIEventType is the type of EventsAPI event recieved.
type EventsAPIEventType string

const (
	// CallbackEvent is the "outer" event of an EventsAPI event.
	CallbackEvent EventsAPIEventType = "event_callback"
	// URLVerification is an event used when configuring your EventsAPI app
	URLVerification EventsAPIEventType = "url_verification"
)

// EventsAPIEventMap maps Event API events to their corresponding struct
// implementations. The structs should be instances of the unmarshalling
// target for the matching event type.
var EventsAPIEventMap = map[EventsAPIEventType]interface{}{
	CallbackEvent:   EventsAPICallbackEvent{},
	URLVerification: EventsAPIURLVerificationEvent{},
}

// EventsAPIEvent is the base EventsAPIEvent
type EventsAPIEvent struct {
	Type       EventsAPIEventType `json:"type"`
	Data       interface{}
	InnerEvent EventsAPIInnerEvent
}

// EventsAPIInnerEvent the inner event of a EventsAPI event_callback Event.
type EventsAPIInnerEvent struct {
	Type string `json:"type"`
	Data interface{}
}

// EventsAPIURLVerificationEvent recieved when configuring a EventsAPI driven app
type EventsAPIURLVerificationEvent struct {
	Token     string             `json:"token"`
	Challenge string             `json:"challenge"`
	Type      EventsAPIEventType `json:"type"`
}

// EventsAPICallbackEvent is the main EventsAPI event.
type EventsAPICallbackEvent struct {
	Type        EventsAPIEventType `json:"type"`
	Token       string             `json:"token"`
	TeamID      string             `json:"team_id"`
	APIAppID    string             `json:"api_app_id"`
	InnerEvent  *json.RawMessage   `json:"event"`
	AuthedUsers []string           `json:"authed_users"`
	EventID     string             `json:"event_id"`
	EventTime   int                `json:"event_time"`
}

// AppMentionEvent is an EventsAPI subscribable event.
type AppMentionEvent struct {
	Type           string `json:"type"`
	User           string `json:"user"`
	Text           string `json:"text"`
	TimeStamp      string `json:"ts"`
	Channel        string `json:"channel"`
	EventTimeStamp string `json:"event_ts"`
}

// ParseOuterEvent parses the outter event of a EventsAPI event.
func ParseOuterEvent(rawEvent json.RawMessage) EventsAPIEvent {
	e := &Event{}
	err := json.Unmarshal(rawEvent, e)
	if err != nil {
		return EventsAPIEvent{
			"unmarshalling_error",
			&UnmarshallingErrorEvent{},
			EventsAPIInnerEvent{},
		}
	}
	if e.Type == string(CallbackEvent) {
		cbEvent := &EventsAPICallbackEvent{}
		err = json.Unmarshal(rawEvent, cbEvent)
		if err != nil {
			fmt.Println(err)
		}
		iE := &Event{}
		innerE := cbEvent.InnerEvent
		err = json.Unmarshal(*innerE, iE)
		fmt.Println(iE.Type)
		v, exists := EventMapping[iE.Type]
		if !exists {
			fmt.Println("lol")
		}
		t := reflect.TypeOf(v)
		recvEvent := reflect.New(t).Interface()
		err = json.Unmarshal(*innerE, recvEvent)
		if err != nil {
			return EventsAPIEvent{
				"unmarshalling_error",
				&UnmarshallingErrorEvent{err},
				EventsAPIInnerEvent{},
			}
		}
		return EventsAPIEvent{
			EventsAPIEventType(e.Type),
			cbEvent,
			EventsAPIInnerEvent{iE.Type, recvEvent},
		}
	}
	// must be a urlverification event
	urlVerificationEvent := &EventsAPIURLVerificationEvent{}
	err = json.Unmarshal(rawEvent, urlVerificationEvent)
	// handle error
	return EventsAPIEvent{
		EventsAPIEventType(e.Type),
		urlVerificationEvent,
		EventsAPIInnerEvent{},
	}
}

// func ParseInnerEvent(e EventsAPIInnerEvent) {
// 	v, exists := EventMapping[e.Type]
// 	if !exists {
// 		fmt.Println("does not exist")
// 	}
// 	t := reflect.TypeOf(v)
// 	reflect.NewAt(t, &e)
// }

// ParseEvent parses an EventsAPI event string and returns the inner event
// func (api *Client) ParseInnerEvent(event json.RawMessage) (EventsAPIInnerEvent, error) {
// 	var e EventsAPIEvent
// 	err := json.Unmarshal(event, &e)
// 	if err != nil {
// 		fmt.Println("unmarshalling err")
// 	}
// 	v, exists := EventMapping[e.Type]
// 	if !exists {
// 		fmt.Println()
// 	}
// }
