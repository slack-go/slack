package slack

// ItemRef is a reference to a message of any type. One of FileID,
// FileCommentId, or the combination of ChannelId and Timestamp must be
// specified.
type ItemRef struct {
	ChannelId     string `json:"channel"`
	Timestamp     string `json:"timestamp"`
	FileId        string `json:"file"`
	FileCommentId string `json:"file_comment"`
}

// NewRefToMessage initializes a reference to to a message.
func NewRefToMessage(channelID, timestamp string) ItemRef {
	return ItemRef{ChannelId: channelID, Timestamp: timestamp}
}

// NewRefToFile initializes a reference to a file.
func NewRefToFile(fileID string) ItemRef {
	return ItemRef{FileId: fileID}
}

// NewRefToFileComment initializes a reference to a file comment.
func NewRefToFileComment(fileCommentID string) ItemRef {
	return ItemRef{FileCommentId: fileCommentID}
}
