package slack

type fileActionEvent struct {
	Type           string         `json:"type"`
	EventTimestamp JSONTimeString `json:"event_ts"`
	File           File           `json:"file"`
	// FileID is used for FileDeletedEvent
	FileID string `json:"file_id,omitempty"`
}

type FileCreatedEvent fileActionEvent
type FileSharedEvent fileActionEvent
type FilePublicEvent fileActionEvent
type FileUnsharedEvent fileActionEvent
type FileChangeEvent fileActionEvent
type FileDeletedEvent fileActionEvent
type FilePrivateEvent fileActionEvent

type FileCommentAddedEvent struct {
	fileActionEvent
	Comment Comment `json:"comment"`
}

type FileCommentEditedEvent struct {
	fileActionEvent
	Comment Comment `json:"comment"`
}

type FileCommentDeletedEvent struct {
	fileActionEvent
	Comment string `json:"comment"`
}
