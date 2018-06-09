package slack

// TextInputSubtype Accepts email, number, tel, or url. In some form factors, optimized input is provided for this subtype.
type TextInputSubtype string

const (
	// EmailTextInputSubtype email keyboard
	EmailTextInputSubtype TextInputSubtype = "email"
	// NumberTextInputSubtype numeric keyboard
	NumberTextInputSubtype TextInputSubtype = "number"
	// TelTextInputSubtype Phone keyboard
	TelTextInputSubtype TextInputSubtype = "tel"
	// URLTextInputSubtype Phone keyboard
	URLTextInputSubtype TextInputSubtype = "url"
)

// TextInputElement subtype of DialogInput
//	https://api.slack.com/dialogs#option_element_attributes#text_element_attributes
type TextInputElement struct {
	DialogInput
	MaxLength int              `json:"max_length,omitempty"`
	MinLength int              `json:"min_length,omitempty"`
	Hint      string           `json:"hint,omitempty"`
	Subtype   TextInputSubtype `json:"subtype"`
	Value     string           `json:"value"`
}

// NewTextInput constructor for a `text` input
func NewTextInput(name, label, text string) *TextInputElement {
	return &TextInputElement{
		DialogInput: DialogInput{
			Type:  InputTypeText,
			Name:  name,
			Label: label,
		},
		Value: text,
	}
}

// NewTextAreaInput constructor for a `textarea` input
func NewTextAreaInput(name, label, text string) *TextInputElement {
	return &TextInputElement{
		DialogInput: DialogInput{
			Type:  InputTypeTextArea,
			Name:  name,
			Label: label,
		},
		Value: text,
	}
}
