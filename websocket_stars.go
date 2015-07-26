package slack

type starEvent struct {
	Type           string         `json:"type"`
	User           string         `json:"user"`
	Item           StarredItem    `json:"item"`
	EventTimestamp JSONTimeString `json:"event_ts"`
}
type StarAddedEvent starEvent
type StarRemovedEvent starEvent
