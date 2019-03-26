package slack

// ContextBlock defines data that is used to display message context, which can
// include both images and text.
//
// More Information: https://api.slack.com/reference/messaging/blocks#actions
type ContextBlock struct {
	Type     MessageBlockType `json:"type"`
	BlockID  string           `json:"block_id,omitempty"`
	Elements ContextElements  `json:"elements"`
}

// blockType returns the type of the block
func (s ContextBlock) blockType() MessageBlockType {
	return s.Type
}

type ContextElements struct {
	ImageElements []*ImageBlockElement
	TextObjects   []*TextBlockObject
}

// NewContextBlock returns a new instance of a context block
func NewContextBlock(blockID string, elements ContextElements) *ContextBlock {
	return &ContextBlock{
		Type:     mbtContext,
		BlockID:  blockID,
		Elements: elements,
	}
}
