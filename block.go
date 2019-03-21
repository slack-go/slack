package slack

import (
	"encoding/json"
)

// @NOTE: Blocks are in beta and subject to change.

// More Information: https://api.slack.com/block-kit

// MessageBlockType defines a named string type to define each block type
// as a constant for use within the package.
type MessageBlockType string
type MessageElementType string
type MessageObjectType string

const (
	mbtSection MessageBlockType = "section"
	mbtDivider MessageBlockType = "divider"
	mbtImage   MessageBlockType = "image"
	mbtAction  MessageBlockType = "actions"
	mbtContext MessageBlockType = "context"

	metImage      MessageElementType = "image"
	metButton     MessageElementType = "button"
	metOverflow   MessageElementType = "overflow"
	metDatepicker MessageElementType = "datepicker"
	metSelect     MessageElementType = "static_select"

	motImage        MessageObjectType = "image"
	motConfirmation MessageObjectType = "confirmation"
	motOption       MessageObjectType = "option"
	motOptionGroup  MessageObjectType = "option_group"
)

// Block defines an interface all block types should implement
// to ensure consistency between blocks.
type Block interface {
	blockType() MessageBlockType
}

// Blocks is a convenience struct defined to allow dynamic unmarshalling of
// the "blocks" value in Slack's JSON response, which varies depending on block type
type Blocks struct {
	BlockSet []Block `json:"blocks"`
}

// UnmarshalJSON implements the Unmarshaller interface for Blocks, so that any JSON
// unmarshalling is delegated and proper type determination can be made before unmarshal
func (b *Blocks) UnmarshalJSON(data []byte) error {
	var raw []json.RawMessage
	err := json.Unmarshal(data, &raw)
	if err != nil {
		return err
	}

	var blocks Blocks
	for _, r := range raw {
		var obj map[string]interface{}
		err := json.Unmarshal(r, &obj)
		if err != nil {
			return err
		}

		blockType := ""
		if t, ok := obj["type"].(string); ok {
			blockType = t
		}

		var block Block
		switch blockType {
		case "actions":
			block = &ActionBlock{}
		case "context":
			block = &ContextBlock{}
		case "divider":
			block = &DividerBlock{}
		case "image":
			block = &ImageBlock{}
		case "section":
			block = &SectionBlock{}
		}

		err = json.Unmarshal(r, block)
		if err != nil {
			return err
		}

		blocks.BlockSet = append(blocks.BlockSet, block)
	}

	*b = blocks
	return nil
}

// NewBlockMessage creates a new Message that contains one or more blocks to be displayed
func NewBlockMessage(blocks ...Block) Message {
	return Message{
		Msg: Msg{
			Blocks: blocks,
		},
	}
}

// AddBlockMessage appends a block to the end of the existing list of blocks
func AddBlockMessage(message Message, newBlk Block) Message {
	message.Msg.Blocks = append(message.Msg.Blocks, newBlk)
	return message
}
