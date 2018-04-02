package slack

import (
	"encoding/json"
	"fmt"
	"reflect"
)

const (
	// CallbackEvent is the "outer" event of an EventsAPI event.
	CallbackEvent = "event_callback"
	// URLVerification is an event used when configuring your EventsAPI app
	URLVerification = "url_verification"
)

// EventsAPIEventMap maps Event API events to their corresponding struct
// implementations. The structs should be instances of the unmarshalling
// target for the matching event type.
var EventsAPIEventMap = map[string]interface{}{
	CallbackEvent:   EventsAPICallbackEvent{},
	URLVerification: EventsAPIURLVerificationEvent{},
}

// EventsAPIEvent is the base EventsAPIEvent
type EventsAPIEvent struct {
	Type       string `json:"type"`
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
	Token     string `json:"token"`
	Challenge string `json:"challenge"`
	Type      string `json:"type"`
}

// EventsAPICallbackEvent is the main EventsAPI event.
type EventsAPICallbackEvent struct {
	Type        string           `json:"type"`
	Token       string           `json:"token"`
	TeamID      string           `json:"team_id"`
	APIAppID    string           `json:"api_app_id"`
	InnerEvent  *json.RawMessage `json:"event"`
	AuthedUsers []string         `json:"authed_users"`
	EventID     string           `json:"event_id"`
	EventTime   int              `json:"event_time"`
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

func parseOuterEvent(rawE json.RawMessage) EventsAPIEvent {
	e := &Event{}
	err := json.Unmarshal(rawE, e)
	if err != nil {
		return EventsAPIEvent{
			"unmarshalling_error",
			&UnmarshallingErrorEvent{},
			EventsAPIInnerEvent{},
		}
	}
	if e.Type == CallbackEvent {
		cbEvent := &EventsAPICallbackEvent{}
		err = json.Unmarshal(rawE, cbEvent)
		if err != nil {
			fmt.Println(err)
		}
		return EventsAPIEvent{
			e.Type,
			cbEvent,
			EventsAPIInnerEvent{},
		}
	}
	urlVE := &EventsAPIURLVerificationEvent{}
	err = json.Unmarshal(rawE, urlVE)
	if err != nil {
		fmt.Println("lol")
	}
	return EventsAPIEvent{
		e.Type,
		urlVE,
		EventsAPIInnerEvent{},
	}
}

func parseInnerEvent(e *EventsAPICallbackEvent) EventsAPIEvent {
	iE := &Event{}
	rawInnerJSON := e.InnerEvent
	err := json.Unmarshal(*rawInnerJSON, iE)
	if err != nil {
		fmt.Println("FAIL!")
	}
	v, exists := EventMapping[iE.Type]
	if !exists {
		fmt.Println("lol!")
	}
	t := reflect.TypeOf(v)
	recvEvent := reflect.New(t).Interface()
	err = json.Unmarshal(*rawInnerJSON, recvEvent)
	if err != nil {
		return EventsAPIEvent{
			"unmarshalling_error",
			&UnmarshallingErrorEvent{err},
			EventsAPIInnerEvent{},
		}
	}
	return EventsAPIEvent{
		e.Type,
		e,
		EventsAPIInnerEvent{iE.Type, recvEvent},
	}
}

// ParseEventsAPIEvent parses the outter event of a EventsAPI event.
func ParseEventsAPIEvent(rawEvent json.RawMessage) EventsAPIEvent {
	e := parseOuterEvent(rawEvent)
	if e.Type == CallbackEvent {
		cbEvent := e.Data.(*EventsAPICallbackEvent)
		return parseInnerEvent(cbEvent)
	}
	urlVerificationEvent := &EventsAPIURLVerificationEvent{}
	err := json.Unmarshal(rawEvent, urlVerificationEvent)
	if err != nil {
		fmt.Println("ho")
	}
	return EventsAPIEvent{
		e.Type,
		urlVerificationEvent,
		EventsAPIInnerEvent{},
	}
}
