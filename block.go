package slack

// @NOTE: Blocks are in beta and subject to change.

// More Information: https://api.slack.com/block-kit

// MessageBlockType defines a named string type to define each block type
// as a constant for use within the package.
type MessageBlockType string

const (
	MBTSection MessageBlockType = "section"
	MBTDivider MessageBlockType = "divider"
	MBTImage   MessageBlockType = "image"
	MBTAction  MessageBlockType = "actions"
	MBTContext MessageBlockType = "context"
)

// Block defines an interface all block types should implement
// to ensure consistency between blocks.
type Block interface {
	BlockType() MessageBlockType
}

// Blocks is a convenience struct defined to allow dynamic unmarshalling of
// the "blocks" value in Slack's JSON response, which varies depending on block type
type Blocks struct {
	BlockSet []Block `json:"blocks"`
}

// BlockAction is the action callback sent when a block is interacted with
type BlockAction struct {
	ActionID string          `json:"action_id"`
	BlockID  string          `json:"block_id"`
	Text     TextBlockObject `json:"text"`
	Value    string          `json:"value"`
	Type     actionType      `json:"type"`
	ActionTs string          `json:"action_ts"`
}

// actionType returns the type of the action
func (b BlockAction) actionType() actionType {
	return b.Type
}

// NewBlockMessage creates a new Message that contains one or more blocks to be displayed
func NewBlockMessage(blocks ...Block) Message {
	return Message{
		Msg: Msg{
			Blocks: Blocks{
				BlockSet: blocks,
			},
		},
	}
}

// AddBlockMessage appends a block to the end of the existing list of blocks
func AddBlockMessage(message Message, newBlk Block) Message {
	message.Msg.Blocks.BlockSet = append(message.Msg.Blocks.BlockSet, newBlk)
	return message
}
