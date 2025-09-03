package slack

import "encoding/json"

// AttachmentField contains information for an attachment field
// An Attachment can contain multiple of these
type AttachmentField struct {
	Title string `json:"title" form:"title"`
	Value string `json:"value" form:"value"`
	Short bool   `json:"short" from:"short"`
}

// AttachmentAction is a button or menu to be included in the attachment. Required when
// using message buttons or menus and otherwise not useful. A maximum of 5 actions may be
// provided per attachment.
type AttachmentAction struct {
	Name            string                        `json:"name" form:"name"`                                   // Required.
	Text            string                        `json:"text" from:"text"`                                   // Required.
	Style           string                        `json:"style,omitempty" form:"style"`                       // Optional. Allowed values: "default", "primary", "danger".
	Type            ActionType                    `json:"type" form:"type"`                                   // Required. Must be set to "button" or "select".
	Value           string                        `json:"value,omitempty" form:"value"`                       // Optional.
	DataSource      string                        `json:"data_source,omitempty" form:"data_source"`           // Optional.
	MinQueryLength  int                           `json:"min_query_length,omitempty" form:"min_query_length"` // Optional. Default value is 1.
	Options         []AttachmentActionOption      `json:"options,omitempty" form:"options"`                   // Optional. Maximum of 100 options can be provided in each menu.
	SelectedOptions []AttachmentActionOption      `json:"selected_options,omitempty" form:"selected_options"` // Optional. The first element of this array will be set as the pre-selected option for this menu.
	OptionGroups    []AttachmentActionOptionGroup `json:"option_groups,omitempty" form:"option_groups"`       // Optional.
	Confirm         *ConfirmationField            `json:"confirm,omitempty" form:"confirm"`                   // Optional.
	URL             string                        `json:"url,omitempty" form:"url"`                           // Optional.
}

// actionType returns the type of the action
func (a AttachmentAction) actionType() ActionType {
	return a.Type
}

// AttachmentActionOption the individual option to appear in action menu.
type AttachmentActionOption struct {
	Text        string `json:"text" form:"text"`                         // Required.
	Value       string `json:"value" form:"value"`                       // Required.
	Description string `json:"description,omitempty" form:"description"` // Optional. Up to 30 characters.
}

// AttachmentActionOptionGroup is a semi-hierarchal way to list available options to appear in action menu.
type AttachmentActionOptionGroup struct {
	Text    string                   `json:"text" form:"text"`       // Required.
	Options []AttachmentActionOption `json:"options" form:"options"` // Required.
}

// AttachmentActionCallback is sent from Slack when a user clicks a button in an interactive message (aka AttachmentAction)
// DEPRECATED: use InteractionCallback
type AttachmentActionCallback InteractionCallback

// ConfirmationField are used to ask users to confirm actions
type ConfirmationField struct {
	Title       string `json:"title,omitempty" form:"title"`               // Optional.
	Text        string `json:"text" form:"text"`                           // Required.
	OkText      string `json:"ok_text,omitempty" form:"ok_text"`           // Optional. Defaults to "Okay"
	DismissText string `json:"dismiss_text,omitempty" form:"dismiss_text"` // Optional. Defaults to "Cancel"
}

// Attachment contains all the information for an attachment
type Attachment struct {
	Color    string `json:"color,omitempty" form:"color"`
	Fallback string `json:"fallback,omitempty" form:"fallback"`

	CallbackID string `json:"callback_id,omitempty" form:"callback_id"`
	ID         int    `json:"id,omitempty" form:"id"`

	AuthorID      string `json:"author_id,omitempty" form:"author_id"`
	AuthorName    string `json:"author_name,omitempty" form:"author_name"`
	AuthorSubname string `json:"author_subname,omitempty" form:"author_subname"`
	AuthorLink    string `json:"author_link,omitempty" form:"author_link"`
	AuthorIcon    string `json:"author_icon,omitempty" form:"author_icon"`

	Title     string `json:"title,omitempty" form:"title"`
	TitleLink string `json:"title_link,omitempty" form:"title_link"`
	Pretext   string `json:"pretext,omitempty" form:"pretext"`
	Text      string `json:"text,omitempty" form:"text"`

	ImageURL string `json:"image_url,omitempty" form:"image_url"`
	ThumbURL string `json:"thumb_url,omitempty" form:"thumb_url"`

	ServiceName string `json:"service_name,omitempty" form:"service_name"`
	ServiceIcon string `json:"service_icon,omitempty" form:"service_icon"`
	FromURL     string `json:"from_url,omitempty" form:"from_url"`
	OriginalURL string `json:"original_url,omitempty" form:"original_url"`

	Fields     []AttachmentField  `json:"fields,omitempty" form:"fields"`
	Actions    []AttachmentAction `json:"actions,omitempty" form:"actions"`
	MarkdownIn []string           `json:"mrkdwn_in,omitempty" form:"mrkdwn_in"`

	Blocks Blocks `json:"blocks,omitempty" form:"blocks"`

	Footer     string `json:"footer,omitempty" form:"footer"`
	FooterIcon string `json:"footer_icon,omitempty" form:"footer_icon"`

	Ts json.Number `json:"ts,omitempty" form:"ts"`
}
