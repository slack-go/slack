package slack

// SlackMetadata https://api.slack.com/reference/metadata
type SlackMetadata struct {
	EventType    string      `json:"event_type"`
	EventPayload interface{} `json:"event_payload"`
}
