package slack

import (
	"encoding/json"
)

// @NOTE: Blocks are in beta and subject to change.

// More Information: https://api.slack.com/block-kit

// MessageBlockType defines a named string type to define each block type
// as a constant for use within the package.
type MessageBlockType string

const (
	mbtSection MessageBlockType = "section"
	mbtDivider MessageBlockType = "divider"
	mbtImage   MessageBlockType = "image"
	mbtAction  MessageBlockType = "actions"
	mbtContext MessageBlockType = "context"
)

// Block defines an interface all block types should implement
// to ensure consistency between blocks.
type Block interface {
	blockType() MessageBlockType
}

// Blocks is a convenience struct defined to allow dynamic unmarshalling of
// the "blocks" value in Slack's JSON response, which varies depending on block type
type Blocks struct {
	ActionBlocks  []*ActionBlock
	ContextBlocks []*ContextBlock
	DividerBlocks []*DividerBlock
	ImageBlocks   []*ImageBlock
	SectionBlocks []*SectionBlock
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

// actionType returns the type of the block action
func (b BlockAction) actionType() actionType {
	return b.Type
}

// NewBlockMessage creates a new Message that contains one or more blocks to be displayed
func NewBlockMessage(blocks ...Block) Message {
	b := Blocks{}
	b.appendToBlocks(blocks)
	return Message{
		Msg: Msg{
			Blocks: b,
		},
	}
}

// AddBlockMessage appends a block to the end of the existing list of blocks
func AddBlockMessage(message Message, newBlk Block) Message {
	message.Msg.Blocks.appendToBlocks([]Block{newBlk})
	return message
}

func (b *Blocks) appendToBlocks(appendBlocks []Block) {
	for _, block := range appendBlocks {
		switch block.(type) {
		case *ActionBlock:
			b.ActionBlocks = append(b.ActionBlocks, block.(*ActionBlock))
		case *ContextBlock:
			b.ContextBlocks = append(b.ContextBlocks, block.(*ContextBlock))
		case *DividerBlock:
			b.DividerBlocks = append(b.DividerBlocks, block.(*DividerBlock))
		case *ImageBlock:
			b.ImageBlocks = append(b.ImageBlocks, block.(*ImageBlock))
		case *SectionBlock:
			b.SectionBlocks = append(b.SectionBlocks, block.(*SectionBlock))
		}
	}
}

// UnmarshalJSON implements the Unmarshaller interface for Blocks, so that any JSON
// unmarshalling is delegated and proper type determination can be made before unmarshal
func (b *Blocks) UnmarshalJSON(data []byte) error {
	var raw []json.RawMessage
	err := json.Unmarshal(data, &raw)
	if err != nil {
		return err
	}

	for _, r := range raw {
		var obj map[string]interface{}
		err := json.Unmarshal(r, &obj)
		if err != nil {
			return err
		}

		var blockType string
		if t, ok := obj["type"].(string); ok {
			blockType = t
		}

		switch blockType {
		case "actions":
			block, err := unmarshalBlock(r, &ActionBlock{})
			if err != nil {
				return err
			}
			b.ActionBlocks = append(b.ActionBlocks, block.(*ActionBlock))
		case "context":
			block, err := unmarshalBlock(r, &ContextBlock{})
			if err != nil {
				return err
			}
			b.ContextBlocks = append(b.ContextBlocks, block.(*ContextBlock))
		case "divider":
			block, err := unmarshalBlock(r, &DividerBlock{})
			if err != nil {
				return err
			}
			b.DividerBlocks = append(b.DividerBlocks, block.(*DividerBlock))
		case "image":
			block, err := unmarshalBlock(r, &ImageBlock{})
			if err != nil {
				return err
			}
			b.ImageBlocks = append(b.ImageBlocks, block.(*ImageBlock))
		case "section":
			block, err := unmarshalBlock(r, &SectionBlock{})
			if err != nil {
				return err
			}
			b.SectionBlocks = append(b.SectionBlocks, block.(*SectionBlock))
		}
	}

	return nil
}

func unmarshalBlock(r json.RawMessage, block Block) (Block, error) {
	err := json.Unmarshal(r, block)
	if err != nil {
		return nil, err
	}
	return block, nil
}
