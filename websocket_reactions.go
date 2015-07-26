package slack

type reactionEvent struct {
	Type           string         `json:"type"`
	UserId         string         `json:"user"`
	Item           ReactedItem    `json:"item"`
	Reaction       string         `json:"reaction"`
	EventTimestamp JSONTimeString `json:"event_ts"`
}
type ReactionAddedEvent reactionEvent
type ReactionRemovedEvent reactionEvent
