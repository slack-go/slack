package slack

// ContextBlock defines data that is used to display message context, which can
// include both images and text.
//
// More Information: https://api.slack.com/reference/messaging/blocks#actions
type ContextBlock struct {
	Type            MessageBlockType `json:"type"`
	BlockID         string           `json:"block_id,omitempty"`
	ContextElements ContextElements  `json:"elements"`
}

// BlockType returns the type of the block
func (s ContextBlock) BlockType() MessageBlockType {
	return s.Type
}

type ContextElements struct {
	Elements []MixedElement
}

// NewContextElements is a convenience method for generating ContextElements
func NewContextElements(elements []MixedElement) ContextElements {
	return ContextElements{
		Elements: elements,
	}
}

// NewContextBlock returns a new instance of a context block
func NewContextBlock(blockID string, elements ContextElements) *ContextBlock {
	return &ContextBlock{
		Type:     MbtContext,
		BlockID:  blockID,
		Elements: elements,
	}
}
