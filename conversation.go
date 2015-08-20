package slack

// Conversation is the foundation for IM and BaseGroupConversation
type Conversation struct {
	ID                 string   `json:"id"`
	Created            JSONTime `json:"created"`
	IsOpen             bool     `json:"is_open"`
	LastRead           string   `json:"last_read,omitempty"`
	Latest             *Message `json:"latest,omitempty"`
	UnreadCount        int      `json:"unread_count,omitempty"`
	UnreadCountDisplay int      `json:"unread_count_display,omitempty"`
}

// GroupConversation is the foundation for Group and Channel
type GroupConversation struct {
	Conversation
	Name       string   `json:"name"`
	Creator    string   `json:"creator"`
	IsChannel  bool     `json:"is_channel"`
	IsGroup    bool     `json:"is_group"`
	IsGeneral  bool     `json:"is_general"`
	IsMember   bool     `json:"is_member"`
	IsArchived bool     `json:"is_archived"`
	Members    []string `json:"members"`
	NumMembers int      `json:"num_members,omitempty"`
	Topic      Topic    `json:"topic"`
	Purpose    Purpose  `json:"purpose"`
}

// Topic contains information about the topic
type Topic struct {
	Value   string   `json:"value"`
	Creator string   `json:"creator"`
	LastSet JSONTime `json:"last_set"`
}

// Purpose contains information about the purpose
type Purpose struct {
	Value   string   `json:"value"`
	Creator string   `json:"creator"`
	LastSet JSONTime `json:"last_set"`
}
