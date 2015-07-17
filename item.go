package slack

const (
	TYPE_MESSAGE = "message"
	TYPE_FILE    = "file"
	TYPE_COMMENT = "comment"
)

// Item is any type of slack message - message, file, or comment.
type Item struct {
	Type    string
	Message *Message
	File    *File
	Comment *Comment
}

// NewMessageItem turns a message into a typed message struct.
func NewMessageItem(m *Message) Item {
	return Item{Type: TYPE_MESSAGE, Message: m}
}

// NewFileItem turns a file into a typed file struct.
func NewFileItem(f *File) Item {
	return Item{Type: TYPE_FILE, File: f}
}

// NewCommentItem turns a comment into a typed comment struct.
func NewCommentItem(c *Comment) Item {
	return Item{Type: TYPE_COMMENT, Comment: c}
}

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
