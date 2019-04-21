package slack

import (
	"encoding/json"

	"github.com/pkg/errors"
)

// UnmarshalJSON implements the Unmarshaller interface for Blocks, so that any JSON
// unmarshalling is delegated and proper type determination can be made before unmarshal
func (b *Blocks) UnmarshalJSON(data []byte) error {
	var raw []json.RawMessage

	if string(data) == "{\"blocks\":null}" {
		return nil
	}

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

		var blockType string
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
		default:
			return errors.New("unsupported block type")
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

// UnmarshalJSON implements the Unmarshaller interface for BlockElements, so that any JSON
// unmarshalling is delegated and proper type determination can be made before unmarshal
func (b *BlockElements) UnmarshalJSON(data []byte) error {
	var raw []json.RawMessage

	if string(data) == "{\"elements\":null}" {
		return nil
	}

	err := json.Unmarshal(data, &raw)
	if err != nil {
		return err
	}

	var blockElements BlockElements
	for _, r := range raw {
		var obj map[string]interface{}
		err := json.Unmarshal(r, &obj)
		if err != nil {
			return err
		}

		var blockElementType string
		if t, ok := obj["type"].(string); ok {
			blockElementType = t
		}

		var blockElement BlockElement
		switch blockElementType {
		case "image":
			blockElement = &ImageBlockElement{}
		case "button":
			blockElement = &ButtonBlockElement{}
		case "overflow":
			blockElement = &OverflowBlockElement{}
		case "datepicker":
			blockElement = &DatePickerBlockElement{}
		case "static_select":
			blockElement = &SelectBlockElement{}
		default:
			return errors.New("unsupported block element type")
		}

		err = json.Unmarshal(r, blockElement)
		if err != nil {
			return err
		}

		blockElements.BlockElementSet = append(blockElements.BlockElementSet, blockElement)
	}

	*b = blockElements
	return nil
}

// MarshalJSON implements the Marshaller interface for Accessory so that any JSON
// marshalling is delegated and proper type determination can be made before marshal
func (a *Accessory) MarshalJSON() ([]byte, error) {
	bytes, err := json.Marshal(toBlockElement(a))
	if err != nil {
		return nil, err
	}

	return bytes, nil
}

// UnmarshalJSON implements the Unmarshaller interface for Accessory, so that any JSON
// unmarshalling is delegated and proper type determination can be made before unmarshal
func (a *Accessory) UnmarshalJSON(data []byte) error {
	var r json.RawMessage

	if string(data) == "{\"accessory\":null}" {
		return nil
	}

	err := json.Unmarshal(data, &r)
	if err != nil {
		return err
	}

	var obj map[string]interface{}
	err = json.Unmarshal(r, &obj)
	if err != nil {
		return err
	}

	var blockElementType string
	if t, ok := obj["type"].(string); ok {
		blockElementType = t
	}

	switch blockElementType {
	case "image":
		element, err := unmarshalBlockElement(r, &ImageBlockElement{})
		if err != nil {
			return err
		}
		a.ImageElement = element.(*ImageBlockElement)
	case "button":
		element, err := unmarshalBlockElement(r, &ButtonBlockElement{})
		if err != nil {
			return err
		}
		a.ButtonElement = element.(*ButtonBlockElement)
	case "overflow":
		element, err := unmarshalBlockElement(r, &OverflowBlockElement{})
		if err != nil {
			return err
		}
		a.OverflowElement = element.(*OverflowBlockElement)
	case "datepicker":
		element, err := unmarshalBlockElement(r, &DatePickerBlockElement{})
		if err != nil {
			return err
		}
		a.DatePickerElement = element.(*DatePickerBlockElement)
	case "static_select":
		element, err := unmarshalBlockElement(r, &SelectBlockElement{})
		if err != nil {
			return err
		}
		a.SelectElement = element.(*SelectBlockElement)
	}

	return nil
}

func unmarshalBlockElement(r json.RawMessage, element BlockElement) (BlockElement, error) {
	err := json.Unmarshal(r, element)
	if err != nil {
		return nil, err
	}
	return element, nil
}

func toBlockElement(element *Accessory) BlockElement {
	if element.ImageElement != nil {
		return element.ImageElement
	}
	if element.ButtonElement != nil {
		return element.ButtonElement
	}
	if element.OverflowElement != nil {
		return element.OverflowElement
	}
	if element.DatePickerElement != nil {
		return element.DatePickerElement
	}
	if element.SelectElement != nil {
		return element.SelectElement
	}

	return nil
}

// UnmarshalJSON implements the Unmarshaller interface for ContextElements, so that any JSON
// unmarshalling is delegated and proper type determination can be made before unmarshal
func (e *ContextElements) UnmarshalJSON(data []byte) error {
	var raw []json.RawMessage

	if string(data) == "{\"elements\":null}" {
		return nil
	}

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

			e.ContextElementSet = append(e.ContextElementSet, elem.(*TextBlockObject))
		case "image":
			elem, err := unmarshalBlockElement(r, &ImageBlockElement{})
			if err != nil {
				return err
			}

			e.ContextElementSet = append(e.ContextElementSet, elem.(*ImageBlockElement))
		default:
			return errors.New("unsupported context element type")
		}
	}

	return nil
}
