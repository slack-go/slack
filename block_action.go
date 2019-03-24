package slack

// ActionBlock defines data that is used to hold interactive elements.
//
// More Information: https://api.slack.com/reference/messaging/blocks#actions
type ActionBlock struct {
	Type     MessageBlockType `json:"type"`
	BlockID  string           `json:"block_id,omitempty"`
	Elements BlockElements    `json:"elements"`
}

// blockType returns the type of the block
func (s ActionBlock) blockType() MessageBlockType {
	return s.Type
}

// NewActionBlock returns a new instance of an Action Block
func NewActionBlock(blockID string, elements ...BlockElement) *ActionBlock {
	e := BlockElements{}
	e.appendToBlockElements(elements)
	return &ActionBlock{
		Type:     mbtAction,
		BlockID:  blockID,
		Elements: e,
	}
}
