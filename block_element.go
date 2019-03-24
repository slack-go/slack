package slack

import (
	"encoding/json"
)

// https://api.slack.com/reference/messaging/block-elements

const (
	metImage      MessageElementType = "image"
	metButton     MessageElementType = "button"
	metOverflow   MessageElementType = "overflow"
	metDatepicker MessageElementType = "datepicker"
	metSelect     MessageElementType = "static_select"

	mixedElementImage mixedElementType = "mixed_image"
	mixedElementText  mixedElementType = "mixed_text"
)

type MessageElementType string
type mixedElementType string

// BlockElement defines an interface that all block element types should implement.
type BlockElement interface {
	elementType() MessageElementType
}

type mixedElement interface {
	mixedElementType() mixedElementType
}

type Accessory struct {
	ImageElement      *ImageBlockElement
	ButtonElement     *ButtonBlockElement
	OverflowElement   *OverflowBlockElement
	DatePickerElement *DatePickerBlockElement
	SelectElement     *SelectBlockElement
}

// NewAccessory returns a new Accessory for a given block element
func NewAccessory(element BlockElement) *Accessory {
	switch element.(type) {
	case *ImageBlockElement:
		return &Accessory{ImageElement: element.(*ImageBlockElement)}
	case *ButtonBlockElement:
		return &Accessory{ButtonElement: element.(*ButtonBlockElement)}
	case *OverflowBlockElement:
		return &Accessory{OverflowElement: element.(*OverflowBlockElement)}
	case *DatePickerBlockElement:
		return &Accessory{DatePickerElement: element.(*DatePickerBlockElement)}
	case *SelectBlockElement:
		return &Accessory{SelectElement: element.(*SelectBlockElement)}
	}

	return nil
}

// the "elements" value in Slack's JSON response, which varies depending on BlockElement type
type BlockElements struct {
	ImageElements      []*ImageBlockElement
	ButtonElements     []*ButtonBlockElement
	OverflowElements   []*OverflowBlockElement
	DatePickerElements []*DatePickerBlockElement
	SelectElements     []*SelectBlockElement
}

// ImageBlockElement An element to insert an image - this element can be used
// in section and context blocks only. If you want a block with only an image
// in it, you're looking for the image block.
//
// More Information: https://api.slack.com/reference/messaging/block-elements#image
type ImageBlockElement struct {
	Type     MessageElementType `json:"type"`
	ImageURL string             `json:"image_url"`
	AltText  string             `json:"alt_text"`
}

func (s ImageBlockElement) elementType() MessageElementType {
	return s.Type
}

func (s ImageBlockElement) mixedElementType() mixedElementType {
	return mixedElementImage
}

// NewImageBlockElement returns a new instance of an image block element
func NewImageBlockElement(imageURL, altText string) *ImageBlockElement {
	return &ImageBlockElement{
		Type:     metImage,
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
	Type     MessageElementType       `json:"type,omitempty"`
	Text     *TextBlockObject         `json:"text"`
	ActionID string                   `json:"action_id,omitempty"`
	URL      string                   `json:"url,omitempty"`
	Value    string                   `json:"value,omitempty"`
	Confirm  *ConfirmationBlockObject `json:"confirm,omitempty"`
}

func (s ButtonBlockElement) elementType() MessageElementType {
	return s.Type
}

// NewButtonBlockElement returns an instance of a new button element to be used within a block
func NewButtonBlockElement(actionID, value string, text *TextBlockObject) *ButtonBlockElement {
	return &ButtonBlockElement{
		Type:     metButton,
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

func (s SelectBlockElement) elementType() MessageElementType {
	return MessageElementType(s.Type)
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
	Type     MessageElementType       `json:"type"`
	ActionID string                   `json:"action_id,omitempty"`
	Options  []*OptionBlockObject     `json:"options"`
	Confirm  *ConfirmationBlockObject `json:"confirm,omitempty"`
}

func (s OverflowBlockElement) elementType() MessageElementType {
	return s.Type
}

// NewOverflowBlockElement returns an instance of a new Overflow Block Element
func NewOverflowBlockElement(actionID string, options ...*OptionBlockObject) *OverflowBlockElement {
	return &OverflowBlockElement{
		Type:     metOverflow,
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
	Type        MessageElementType       `json:"type"`
	ActionID    string                   `json:"action_id"`
	Placeholder *TextBlockObject         `json:"placeholder,omitempty"`
	InitialDate string                   `json:"initial_date,omitempty"`
	Confirm     *ConfirmationBlockObject `json:"confirm,omitempty"`
}

func (s DatePickerBlockElement) elementType() MessageElementType {
	return s.Type
}

// NewDatePickerBlockElement returns an instance of a date picker element
func NewDatePickerBlockElement(actionID string) *DatePickerBlockElement {
	return &DatePickerBlockElement{
		Type:     metDatepicker,
		ActionID: actionID,
	}
}
func (e *BlockElements) appendToBlockElements(appendElements []BlockElement) {
	for _, element := range appendElements {
		switch element.(type) {
		case *ImageBlockElement:
			e.ImageElements = append(e.ImageElements, element.(*ImageBlockElement))
		case *ButtonBlockElement:
			e.ButtonElements = append(e.ButtonElements, element.(*ButtonBlockElement))
		case *OverflowBlockElement:
			e.OverflowElements = append(e.OverflowElements, element.(*OverflowBlockElement))
		case *DatePickerBlockElement:
			e.DatePickerElements = append(e.DatePickerElements, element.(*DatePickerBlockElement))
		case *SelectBlockElement:
			e.SelectElements = append(e.SelectElements, element.(*SelectBlockElement))
		}
	}
}

func (e *BlockElements) MarshalJSON() ([]byte, error) {
	bytes, err := json.Marshal(toBlockElementSlice(e))
	if err != nil {
		return nil, err
	}

	return bytes, nil
}

func toBlockElementSlice(elements *BlockElements) []BlockElement {
	var slice []BlockElement
	for _, element := range elements.ImageElements {
		slice = append(slice, element)
	}
	for _, element := range elements.ButtonElements {
		slice = append(slice, element)
	}
	for _, element := range elements.OverflowElements {
		slice = append(slice, element)
	}
	for _, element := range elements.DatePickerElements {
		slice = append(slice, element)
	}
	for _, element := range elements.SelectElements {
		slice = append(slice, element)
	}

	return slice
}

func (a *Accessory) MarshalJSON() ([]byte, error) {
	bytes, err := json.Marshal(toBlockElement(a))
	if err != nil {
		return nil, err
	}

	return bytes, nil
}

func toBlockElement(element *Accessory) BlockElement {
	if element.ImageElement != nil {
		return element.ImageElement
	}
	if element.ButtonElement != nil {
		return element.ButtonElement
	}
	if element.OverflowElement != nil {
		return element.OverflowElement
	}
	if element.DatePickerElement != nil {
		return element.DatePickerElement
	}
	if element.SelectElement != nil {
		return element.SelectElement
	}

	return nil
}

// UnmarshalJSON implements the Unmarshaller interface for BlockElements, so that any JSON
// unmarshalling is delegated and proper type determination can be made before unmarshal
func (b *BlockElements) UnmarshalJSON(data []byte) error {
	var raw []json.RawMessage
	err := json.Unmarshal(data, &raw)
	if err != nil {
		return err
	}

	for _, r := range raw {
		var obj map[string]interface{}
		err := json.Unmarshal(r, &obj)
		if err != nil {
			return err
		}

		var blockElementType string
		if t, ok := obj["type"].(string); ok {
			blockElementType = t
		}

		switch blockElementType {
		case "image":
			element, err := unmarshalBlockElement(r, &ImageBlockElement{})
			if err != nil {
				return err
			}
			b.ImageElements = append(b.ImageElements, element.(*ImageBlockElement))
		case "button":
			element, err := unmarshalBlockElement(r, &ButtonBlockElement{})
			if err != nil {
				return err
			}
			b.ButtonElements = append(b.ButtonElements, element.(*ButtonBlockElement))
		case "overflow":
			element, err := unmarshalBlockElement(r, &OverflowBlockElement{})
			if err != nil {
				return err
			}
			b.OverflowElements = append(b.OverflowElements, element.(*OverflowBlockElement))
		case "datepicker":
			element, err := unmarshalBlockElement(r, &DatePickerBlockElement{})
			if err != nil {
				return err
			}
			b.DatePickerElements = append(b.DatePickerElements, element.(*DatePickerBlockElement))
		case "static_select":
			element, err := unmarshalBlockElement(r, &SelectBlockElement{})
			if err != nil {
				return err
			}
			b.SelectElements = append(b.SelectElements, element.(*SelectBlockElement))
		}
	}

	return nil
}

func (e *ContextElements) MarshalJSON() ([]byte, error) {
	bytes, err := json.Marshal(toMixedElements(e))
	if err != nil {
		return nil, err
	}

	return bytes, nil
}

func toMixedElements(elements *ContextElements) []mixedElement {
	var slice []mixedElement
	for _, element := range elements.ImageElements {
		slice = append(slice, element)
	}
	for _, element := range elements.TextObjects {
		slice = append(slice, element)
	}

	return slice
}

// UnmarshalJSON implements the Unmarshaller interface for Accessory, so that any JSON
// unmarshalling is delegated and proper type determination can be made before unmarshal
func (a *Accessory) UnmarshalJSON(data []byte) error {
	var r json.RawMessage
	err := json.Unmarshal(data, &r)
	if err != nil {
		return err
	}

	var obj map[string]interface{}
	err = json.Unmarshal(r, &obj)
	if err != nil {
		return err
	}

	var blockElementType string
	if t, ok := obj["type"].(string); ok {
		blockElementType = t
	}

	switch blockElementType {
	case "image":
		element, err := unmarshalBlockElement(r, &ImageBlockElement{})
		if err != nil {
			return err
		}
		a.ImageElement = element.(*ImageBlockElement)
	case "button":
		element, err := unmarshalBlockElement(r, &ButtonBlockElement{})
		if err != nil {
			return err
		}
		a.ButtonElement = element.(*ButtonBlockElement)
	case "overflow":
		element, err := unmarshalBlockElement(r, &OverflowBlockElement{})
		if err != nil {
			return err
		}
		a.OverflowElement = element.(*OverflowBlockElement)
	case "datepicker":
		element, err := unmarshalBlockElement(r, &DatePickerBlockElement{})
		if err != nil {
			return err
		}
		a.DatePickerElement = element.(*DatePickerBlockElement)
	case "static_select":
		element, err := unmarshalBlockElement(r, &SelectBlockElement{})
		if err != nil {
			return err
		}
		a.SelectElement = element.(*SelectBlockElement)
	}

	return nil
}
func unmarshalBlockElement(r json.RawMessage, element BlockElement) (BlockElement, error) {
	err := json.Unmarshal(r, element)
	if err != nil {
		return nil, err
	}
	return element, nil
}
