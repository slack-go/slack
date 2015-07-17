package slack

// ItemRef is a reference to a message of any type. One of FileID,
// CommentId, or the combination of ChannelId and Timestamp must be
// specified.
type ItemRef struct {
	ChannelId string `json:"channel"`
	Timestamp string `json:"timestamp"`
	FileId    string `json:"file"`
	CommentId string `json:"file_comment"`
}

// NewRefToMessage initializes a reference to to a message.
func NewRefToMessage(channelID, timestamp string) ItemRef {
	return ItemRef{ChannelId: channelID, Timestamp: timestamp}
}

// NewRefToFile initializes a reference to a file.
func NewRefToFile(fileID string) ItemRef {
	return ItemRef{FileId: fileID}
}

// NewRefToComment initializes a reference to a file comment.
func NewRefToComment(commentID string) ItemRef {
	return ItemRef{CommentId: commentID}
}
