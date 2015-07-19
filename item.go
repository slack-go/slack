package slack

const (
	TYPE_MESSAGE      = "message"
	TYPE_FILE         = "file"
	TYPE_FILE_COMMENT = "file_comment"
	TYPE_CHANNEL      = "channel"
	TYPE_IM           = "im"
	TYPE_GROUP        = "group"
)

// Item is any type of slack message - message, file, or file comment.
type Item struct {
	Type    string   `json:"type"`
	Channel string   `json:"channel,omitempty"`
	Message *Message `json:"message,omitempty"`
	File    *File    `json:"file,omitempty"`
	Comment *Comment `json:"comment,omitempty"`
}

// NewMessageItem turns a message on a channel into a typed message struct.
func NewMessageItem(ch string, m *Message) Item {
	return Item{Type: TYPE_MESSAGE, Channel: ch, Message: m}
}

// NewFileItem turns a file into a typed file struct.
func NewFileItem(f *File) Item {
	return Item{Type: TYPE_FILE, File: f}
}

// NewFileCommentItem turns a file and comment into a typed file_comment struct.
func NewFileCommentItem(f *File, c *Comment) Item {
	return Item{Type: TYPE_FILE_COMMENT, File: f, Comment: c}
}

// NewChannelItem turns a channel id into a typed channel struct.
func NewChannelItem(ch string) Item {
	return Item{Type: TYPE_CHANNEL, Channel: ch}
}

// NewIMItem turns a channel id into a typed im struct.
func NewIMItem(ch string) Item {
	return Item{Type: TYPE_IM, Channel: ch}
}

// NewGroupItem turns a channel id into a typed group struct.
func NewGroupItem(ch string) Item {
	return Item{Type: TYPE_GROUP, Channel: ch}
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
