package slack

import (
	"encoding/json"
	"fmt"
)

// EventPOSTRequest is the expected shape of an EventsAPI event
type EventPOSTRequest struct {
	Token       string `json:"token"`
	TeamID      string `json:"team_id"`
	APIAppID    string `json:"api_app_id"`
	Event       `json:"event"`
	Type        string   `json:"type"`
	AuthedUsers []string `json:"authed_users"`
	EventID     string   `json:"event_id"`
	EventTime   int      `json:"event_time"`
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
func (api *Client) ParseEvent(event string) (EventPOSTRequest, error) {
	var e EventPOSTRequest
	err := json.Unmarshal([]byte(event), &e)
	if err != nil {
		api.Debugf("ParseEvent Error, could not unmarshall event: %s\n", event)
		err := fmt.Errorf("ParseEvent Error, could not unmarshall event: %s", event)
		return EventPOSTRequest{}, err
	}
	typeStr := e.Event.Type
	_, exists := EventMapping[e.Event.Type]
	if !exists {
		api.Debugf("ParseEvent Error, received unmapped event %q: %s\n", typeStr, e.Event)
		err := fmt.Errorf("ParseEvent Error: Received unmapped event %q: %s", typeStr, e.Event)
		return EventPOSTRequest{}, err
	}
	return e, nil
}
