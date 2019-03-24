package slack

import (
	"encoding/json"
)

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

// NewContextBlock returns a newinstance of a context block
func NewContextBlock(blockID string, elements ContextElements) *ContextBlock {
	return &ContextBlock{
		Type:     mbtContext,
		BlockID:  blockID,
		Elements: elements,
	}
}

// UnmarshalJSON implements the Unmarshaller interface for ContextElements, so that any JSON
// unmarshalling is delegated and proper type determination can be made before unmarshal
func (e *ContextElements) UnmarshalJSON(data []byte) error {
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

		contextElementType := ""
		if t, ok := obj["type"].(string); ok {
			contextElementType = t
		}

		switch contextElementType {
		case PlainTextType, MarkdownType:
			elem, err := unmarshalBlockObject(r, &TextBlockObject{})
			if err != nil {
				return err
			}
			e.TextObjects = append(e.TextObjects, elem.(*TextBlockObject))
		case "image":
			elem, err := unmarshalBlockElement(r, &ImageBlockElement{})
			if err != nil {
				return err
			}
			e.ImageElements = append(e.ImageElements, elem.(*ImageBlockElement))
		}
	}

	return nil
}
