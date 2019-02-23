package slack

// https://api.slack.com/reference/messaging/block-elements

// BlockElement defines an interface that all block element types should
// implement.
type BlockElement interface {
	ValidateElement() bool
}

// ImageBlockElement An element to insert an image - this element can be used
// in section and context blocks only. If you want a block with only an image
// in it, you're looking for the image block.
//
// More Information: https://api.slack.com/reference/messaging/block-elements#image
type ImageBlockElement struct {
	Type     string `json:"type"`
	ImageURL string `json:"image_url"`
	AltText  string `json:"alt_text"`
}

// ValidateElement performs validation checks to ensure the element is valid
func (s ImageBlockElement) ValidateElement() bool {
	return true
}

// NewImageBlockElement returns a new instance of an image block element
func NewImageBlockElement(imageURL, altText string) *ImageBlockElement {
	return &ImageBlockElement{
		Type:     "image",
		ImageURL: imageURL,
		AltText:  altText,
	}
}

// ButtonBlockElement defines an interactive element that inserts a button. The
// button can be a trigger for anything from opening a simple link to starting
// a complex workflow.
//
// More Information: https://api.slack.com/reference/messaging/block-elements#button
type ButtonBlockElement struct {
	Type     string                   `json:"type,omitempty"`
	Text     *TextBlockObject         `json:"text"`
	ActionID string                   `json:"action_id,omitempty"`
	URL      string                   `json:"url,omitempty"`
	Value    string                   `json:"value,omitempty"`
	Confirm  *ConfirmationBlockObject `json:"confirm,omitempty"`
}

// ValidateElement performs validation checks to ensure the element is valid
func (s ButtonBlockElement) ValidateElement() bool {
	return true
}

// NewButtonBlockElement returns an instance of a new button element to be used within a block
func NewButtonBlockElement(actionID, value string, text *TextBlockObject) *ButtonBlockElement {
	return &ButtonBlockElement{
		Type:     "button",
		ActionID: actionID,
		Text:     text,
		Value:    value,
	}
}

// SelectBlockElement defines the simplest form of select menu, with a static list
// of options passed in when defining the element.
//
// More Information: https://api.slack.com/reference/messaging/block-elements#select
type SelectBlockElement struct {
	Type          string                    `json:"type,omitempty"`
	Placeholder   *TextBlockObject          `json:"placeholder,omitempty"`
	ActionID      string                    `json:"action_id,omitempty"`
	Options       []*OptionBlockObject      `json:"options,omitempty"`
	OptionGroups  []*OptionGroupBlockObject `json:"option_groups,omitempty"`
	InitialOption *OptionBlockObject        `json:"initial_option,omitempty"`
	Confirm       *ConfirmationBlockObject  `json:"confirm,omitempty"`
}

// ValidateElement performs validation checks to ensure the element is valid
func (s SelectBlockElement) ValidateElement() bool {
	return true
}

// NewOptionsSelectBlockElement returns a new instance of SelectBlockElement for use with
// the Options object only.
func NewOptionsSelectBlockElement(optType string, placeholder *TextBlockObject, actionID string, options ...*OptionBlockObject) *SelectBlockElement {
	return &SelectBlockElement{
		Type:        optType,
		Placeholder: placeholder,
		ActionID:    actionID,
		Options:     options,
	}
}

// NewOptionsGroupSelectBlockElement returns a new instance of SelectBlockElement for use with
// the Options object only.
func NewOptionsGroupSelectBlockElement(
	optType string,
	placeholder *TextBlockObject,
	actionID string,
	optGroups ...*OptionGroupBlockObject,
) *SelectBlockElement {
	return &SelectBlockElement{
		Type:         optType,
		Placeholder:  placeholder,
		ActionID:     actionID,
		OptionGroups: optGroups,
	}
}

// OverflowBlockElement defines the fields needed to use an overflow element.
// And Overflow Element is like a cross between a button and a select menu -
// when a user clicks on this overflow button, they will be presented with a
// list of options to choose from.
//
// More Information: https://api.slack.com/reference/messaging/block-elements#overflow
type OverflowBlockElement struct {
	Type     string                   `json:"type"`
	ActionID string                   `json:"action_id,omitempty"`
	Options  []*OptionBlockObject     `json:"options"`
	Confirm  *ConfirmationBlockObject `json:"confirm,omitempty"`
}

// ValidateElement performs validation checks to ensure the element is valid
func (s OverflowBlockElement) ValidateElement() bool {
	return true
}

// NewOverflowBlockElement returns an instance of a new Overflow Block Element
func NewOverflowBlockElement(actionID string, options ...*OptionBlockObject) *OverflowBlockElement {
	return &OverflowBlockElement{
		Type:     "overflow",
		ActionID: actionID,
		Options:  options,
	}
}

// DatePickerBlockElement defines an element which lets users easily select a
// date from a calendar style UI. Date picker elements can be used inside of
// section and actions blocks.
//
// More Information: https://api.slack.com/reference/messaging/block-elements#datepicker
type DatePickerBlockElement struct {
	Type        string                   `json:"type"`
	ActionID    string                   `json:"action_id"`
	Placeholder *TextBlockObject         `json:"placeholder,omitempty"`
	InitialDate string                   `json:"initial_date,omitempty"`
	Confirm     *ConfirmationBlockObject `json:"confirm,omitempty"`
}

// ValidateElement performs validation checks to ensure the element is valid
func (s DatePickerBlockElement) ValidateElement() bool {
	return true
}

// NewDatePickerBlockElement returns an instance of a date picker element
func NewDatePickerBlockElement(actionID string) *DatePickerBlockElement {
	return &DatePickerBlockElement{
		Type:     "datepicker",
		ActionID: actionID,
	}
}
