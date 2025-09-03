package slack

// Comment contains all the information relative to a comment
type Comment struct {
	ID        string   `json:"id,omitempty" form:"id"`
	Created   JSONTime `json:"created,omitempty" form:"created"`
	Timestamp JSONTime `json:"timestamp,omitempty" form:"timestamp"`
	User      string   `json:"user,omitempty" form:"user"`
	Comment   string   `json:"comment,omitempty" form:"comment"`
}
