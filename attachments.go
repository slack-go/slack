package slack

// AttachmentField contains information for an attachment field
// An Attachment can contain multiple of these
type AttachmentField struct {
	Title string `json:"title"`
	Value string `json:"value"`
	Short bool   `json:"short"`
}

// AttachmentActionConfirm
type AttachmentActionConfirm struct {
	Title       string `json:"title"`
	Text        string `json:"text"`
	OkText      string `json:"ok_text"`
	DismissText string `json:"dismiss_text"`
}

// AttachmentAction is used for creating a button in an outgoing slack message
type AttachmentAction struct {
	Name       string `json:"name"`
	Text       string `json:"text"`
	Type       string `json:"type"`
	Value      string `json:"value,omitempty"`
	Style      string `json:"style,omitempty"`

	Confirm AttachmentActionConfirm `json:"confirm,omitempty"`
}

// Attachment contains all the information for an attachment
type Attachment struct {
	Color    string `json:"color,omitempty"`
	Fallback string `json:"fallback"`

	AuthorName    string `json:"author_name,omitempty"`
	AuthorSubname string `json:"author_subname,omitempty"`
	AuthorLink    string `json:"author_link,omitempty"`
	AuthorIcon    string `json:"author_icon,omitempty"`

	Title     string `json:"title,omitempty"`
	TitleLink string `json:"title_link,omitempty"`
	Pretext   string `json:"pretext,omitempty"`
	Text      string `json:"text"`

	ImageURL string `json:"image_url,omitempty"`
	ThumbURL string `json:"thumb_url,omitempty"`

	Fields     []AttachmentField `json:"fields,omitempty"`
	MarkdownIn []string          `json:"mrkdwn_in,omitempty"`

	Actions []AttachmentAction `json:"actions,omitempty"`
	CallbackID string `json:"callback_id,omitempty"`
}
