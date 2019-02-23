package slack

import "strings"

// SectionBlock defines a new block of type section
//
// More Information: https://api.slack.com/reference/messaging/blocks#section
type SectionBlock struct {
	Type      string             `json:"type"`
	Text      *TextBlockObject   `json:"text,omitempty"`
	BlockID   string             `json:"block_id,omitempty"`
	Fields    []*TextBlockObject `json:"fields,omitempty"`
	Accessory BlockElement       `json:"accessory,omitempty"`
}

// ValidateBlock ensures that the type set to the block is found in the list of
// valid slack block.
func (s *SectionBlock) ValidateBlock() bool {
	return isStringInSlice(validBlockList, strings.ToLower(s.Type))

}

// NewSectionBlock returns a new instance of a section block to be rendered
func NewSectionBlock(textObj *TextBlockObject, fields []*TextBlockObject, accessory BlockElement) *SectionBlock {
	return &SectionBlock{
		Type:      "section",
		Text:      textObj,
		Fields:    fields,
		Accessory: accessory,
	}
}
