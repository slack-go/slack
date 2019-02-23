package slack

import "strings"

// DividerBlock for displaying a divider line between blocks (similar to <hr> tag in html)
//
// More Information: https://api.slack.com/reference/messaging/blocks#divider
type DividerBlock struct {
	Type    string `json:"type"`
	BlockID string `json:"block_id,omitempty"`
}

// ValidateBlock ensures that the type set to the block is found in the list of
// valid slack block.
func (s *DividerBlock) ValidateBlock() bool {
	return isStringInSlice(validBlockList, strings.ToLower(s.Type))

}

// NewDividerBlock returns a new instance of a divider block
func NewDividerBlock() *DividerBlock {

	return &DividerBlock{
		Type: "divider",
	}

}
