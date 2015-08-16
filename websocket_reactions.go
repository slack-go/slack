package slack

type reactionEvent struct {
	Type           string         `json:"type"`
	User           string         `json:"user"`
	Item           ReactedItem    `json:"item"`
	Reaction       string         `json:"reaction"`
	EventTimestamp JSONTimeString `json:"event_ts"`
}

// ReactionAddedEvent represents the Reaction added event
type ReactionAddedEvent reactionEvent

// ReactionRemovedEvent represents the Reaction removed event
type ReactionRemovedEvent reactionEvent
