package slack

// @NOTE: Blocks are in beta and subject to change.

// More Information: https://api.slack.com/block-kit

var (
	// validBlockList contains a list of
	validBlockList = []string{
		"section",
		"divider",
		"image",
		"actions",
		"context",
	}
)

// Block defines an interface all block types should implement
// to ensure consistency between blocks.
type Block interface {
	ValidateBlock() bool
}

// NewBlockMessage creates a new Message that contains one or more blocks to be displayed
func NewBlockMessage(blocks ...Block) Message {
	return Message{
		Msg: Msg{
			Blocks: blocks,
		},
	}
}

// isStringInSlice is a helper function used in validating the block structs to
// verify a valid type has been used.
func isStringInSlice(a []string, x string) bool {
	for _, n := range a {
		if x == n {
			return true
		}
	}
	return false
}
