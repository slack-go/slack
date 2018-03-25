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

// ParseEvent parses an EventsAPI event.
func (api *Client) ParseEvent(event string) (EventPOSTRequest, error) {
	var e EventPOSTRequest
	err := json.Unmarshal([]byte(event), &e)
	if err != nil {
		// panic? handle response
		fmt.Println(err)
		return EventPOSTRequest{}, err
	}
	return e, nil
}
