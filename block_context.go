package slack

import "strings"

// ContextBlock defines data that is used to display message context, which can
// include both images and text.
//
// More Information: https://api.slack.com/reference/messaging/blocks#actions
type ContextBlock struct {
	Type     string        `json:"type"`
	BlockID  string        `json:"block_id,omitempty"`
	Elements []BlockObject `json:"elements"`
}

// ValidateBlock ensures that the type set to the block is found in the list of
// valid slack block.
func (s *ContextBlock) ValidateBlock() bool {
	return isStringInSlice(validBlockList, strings.ToLower(s.Type))
}

// NewContextBlock returns a newinstance of a context block
func NewContextBlock(blockID string, elements ...BlockObject) *ContextBlock {
	return &ContextBlock{
		Type:     "context",
		BlockID:  blockID,
		Elements: elements,
	}
}
