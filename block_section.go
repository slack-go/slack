package slack

// SectionBlock defines a new block of type section
//
// More Information: https://api.slack.com/reference/messaging/blocks#section
type SectionBlock struct {
	Type      MessageBlockType   `json:"type"`
	Text      *TextBlockObject   `json:"text,omitempty"`
	BlockID   string             `json:"block_id,omitempty"`
	Fields    []*TextBlockObject `json:"fields,omitempty"`
	Accessory *Accessory         `json:"accessory,omitempty"`
}

// BlockType returns the type of the block
func (s SectionBlock) BlockType() MessageBlockType {
	return s.Type
}

// SectionBlockOption allows configuration of options for a new section block
type SectionBlockOption func(*SectionBlock)

func SectionBlockOptionBlockID(blockID string) SectionBlockOption {
	return func(block *SectionBlock) {
		block.BlockID = blockID
	}
}

func SectionBlockOptionAccessory(accessory *Accessory) SectionBlockOption {
	return func(block *SectionBlock) {
		block.Accessory = accessory
	}
}

func SectionBlockOptionFields(fields []*TextBlockObject) SectionBlockOption {
	return func(block *SectionBlock) {
		block.Fields = fields
	}
}

// NewSectionBlock returns a new instance of a section block to be rendered
func NewSectionBlock(textObj *TextBlockObject, options ...SectionBlockOption) *SectionBlock {
	block := SectionBlock{
		Type: MBTSection,
		Text: textObj,
	}

	for _, option := range options {
		option(&block)
	}

	return &block
}
