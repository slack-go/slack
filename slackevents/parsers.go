package slackevents

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/ACollectionOfAtoms/slack"
)

// eventsMap checks both slack.EventsMapping and
// and slackevents.EventsAPIInnerEventMapping. If the event
// exists, returns the the unmarshalled struct instance of
// target for the matching event type.
// TODO: Consider moving all events into its own package?
func eventsMap(t string) (interface{}, bool) {
	// Must parse EventsAPI FIRST as both RTM and EventsAPI
	// have a type: "Message" event.
	// TODO: Handle these cases more explicitly.
	v, exists := EventsAPIInnerEventMapping[t]
	if exists {
		return v, exists
	}
	v, exists = slack.EventMapping[t]
	if exists {
		return v, exists
	}
	return v, exists
}

func parseOuterEvent(rawE json.RawMessage) EventsAPIEvent {
	e := &slack.Event{}
	err := json.Unmarshal(rawE, e)
	if err != nil {
		return EventsAPIEvent{
			"unmarshalling_error",
			&slack.UnmarshallingErrorEvent{err},
			EventsAPIInnerEvent{},
		}
	}
	if e.Type == CallbackEvent {
		cbEvent := &EventsAPICallbackEvent{}
		err = json.Unmarshal(rawE, cbEvent)
		if err != nil {
			return EventsAPIEvent{
				"unmarshalling_error",
				&slack.UnmarshallingErrorEvent{err},
				EventsAPIInnerEvent{},
			}
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
		return EventsAPIEvent{
			"unmarshalling_error",
			&slack.UnmarshallingErrorEvent{err},
			EventsAPIInnerEvent{},
		}
	}
	return EventsAPIEvent{
		e.Type,
		urlVE,
		EventsAPIInnerEvent{},
	}
}

func parseInnerEvent(e *EventsAPICallbackEvent) (EventsAPIEvent, error) {
	iE := &slack.Event{}
	rawInnerJSON := e.InnerEvent
	err := json.Unmarshal(*rawInnerJSON, iE)
	if err != nil {
		return EventsAPIEvent{
			"unmarshalling_error",
			&slack.UnmarshallingErrorEvent{err},
			EventsAPIInnerEvent{},
		}, err
	}
	v, exists := eventsMap(iE.Type)
	if !exists {
		return EventsAPIEvent{
			iE.Type,
			nil,
			EventsAPIInnerEvent{},
		}, fmt.Errorf("Inner Event does not exist! %s", iE.Type)
	}
	t := reflect.TypeOf(v)
	recvEvent := reflect.New(t).Interface()
	err = json.Unmarshal(*rawInnerJSON, recvEvent)
	if err != nil {
		return EventsAPIEvent{
			"unmarshalling_error",
			&slack.UnmarshallingErrorEvent{err},
			EventsAPIInnerEvent{},
		}, err
	}
	return EventsAPIEvent{
		e.Type,
		e,
		EventsAPIInnerEvent{iE.Type, recvEvent},
	}, nil
}

// ParseEventsAPIEvent parses the outter and inner events (if applicable) of an events
// api event returning a EventsAPIEvent type. If the event is a url_verification event,
// the inner event is empty.
func ParseEventsAPIEvent(rawEvent json.RawMessage) (EventsAPIEvent, error) {
	e := parseOuterEvent(rawEvent)
	if e.Type == CallbackEvent {
		cbEvent := e.Data.(*EventsAPICallbackEvent)
		innerEvent, err := parseInnerEvent(cbEvent)
		if err != nil {
			err := fmt.Errorf("EventsAPI Error parsing inner event: %s, %s", innerEvent.Type, err)
			return EventsAPIEvent{
				"unmarshalling_error",
				&slack.UnmarshallingErrorEvent{err},
				EventsAPIInnerEvent{},
			}, err
		}
		return innerEvent, nil
	}
	urlVerificationEvent := &EventsAPIURLVerificationEvent{}
	err := json.Unmarshal(rawEvent, urlVerificationEvent)
	if err != nil {
		return EventsAPIEvent{
			"unmarshalling_error",
			&slack.UnmarshallingErrorEvent{err},
			EventsAPIInnerEvent{},
		}, err
	}
	return EventsAPIEvent{
		e.Type,
		urlVerificationEvent,
		EventsAPIInnerEvent{},
	}, nil
}
