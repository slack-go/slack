package slack

import (
	"encoding/json"
	"errors"
	"fmt"
)

// EventsAPIEventType is the type of EventsAPI event recieved.
type EventsAPIEventType string

const (
	// CallbackEventType is the "outer" event of an EventsAPI event.
	CallbackEventType EventsAPIEventType = "event_callback"
	// URLVerification is an event used when configuring your EventsAPI app
	URLVerification EventsAPIEventType = "url_verification"
)

// EventsAPIEvent is an EventsAPI catch-all
type EventsAPIEvent struct {
	Type EventsAPIEventType `json:"type"`
}

// EventsAPIURLVerificationEvent recieved when configuring a EventsAPI driven app
type EventsAPIURLVerificationEvent struct {
	Token     string             `json:"token"`
	Challenge string             `json:"challenge"`
	Type      EventsAPIEventType `json:"type"`
}

// EventsAPICallbackEvent is the
type EventsAPICallbackEvent struct {
	Type        EventsAPIEventType `json:"type"`
	Token       string             `json:"token"`
	TeamID      string             `json:"team_id"`
	APIAppID    string             `json:"api_app_id"`
	Event       Event              `json:"event"`
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

// ParseEvent parses an EventsAPI event string.
func (api *Client) ParseEvent(event string) (EventsAPICallbackEvent, error) {
	var e EventsAPIEvent
	var cE EventsAPICallbackEvent

	err := json.Unmarshal([]byte(event), &e)
	// Currenlty only supporting callback events
	if e.Type != CallbackEventType || err != nil {
		if err == nil {
			err = errors.New("not implemented")
		}
		return EventsAPICallbackEvent{}, err
	}
	err = json.Unmarshal([]byte(event), &cE)
	if err != nil {
		api.Debugf("ParseEvent Error, could not unmarshall event: %s\n", event)
		err := fmt.Errorf("ParseEvent Error, could not unmarshall event: %s", event)
		return EventsAPICallbackEvent{}, err
	}
	typeStr := cE.Event.Type
	_, exists := EventMapping[cE.Event.Type]
	if !exists {
		api.Debugf("ParseEvent Error, received unmapped event %q: %s\n", typeStr, cE.Event)
		err := fmt.Errorf("ParseEvent Error: Received unmapped event %q: %s", typeStr, cE.Event)
		return EventsAPICallbackEvent{}, err
	}
	return cE, nil
}
