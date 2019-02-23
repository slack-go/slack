package slack

import "strings"

// ActionBlock defines data that is used to hold interactive elements.
//
// More Information: https://api.slack.com/reference/messaging/blocks#actions
type ActionBlock struct {
	Type     string         `json:"type"`
	BlockID  string         `json:"block_id,omitempty"`
	Elements []BlockElement `json:"elements"`
}

// ValidateBlock ensures that the type set to the block is found in the list of
// valid slack block.
func (s *ActionBlock) ValidateBlock() bool {
	return isStringInSlice(validBlockList, strings.ToLower(s.Type))
}

// NewActionBlock returns a new instance of an Action Block
func NewActionBlock(blockID string, elements ...BlockElement) *ActionBlock {
	return &ActionBlock{
		Type:     "actions",
		BlockID:  blockID,
		Elements: elements,
	}
}
