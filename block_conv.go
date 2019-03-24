package slack

import "encoding/json"

// Conv/JSON encoding logic for Blocks

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

// Conv/JSON encoding logic for BlockElements

func (e *BlockElements) MarshalJSON() ([]byte, error) {
	bytes, err := json.Marshal(toBlockElementSlice(e))
	if err != nil {
		return nil, err
	}

	return bytes, nil
}

// UnmarshalJSON implements the Unmarshaller interface for BlockElements, so that any JSON
// unmarshalling is delegated and proper type determination can be made before unmarshal
func (b *BlockElements) UnmarshalJSON(data []byte) error {
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
			b.ImageElements = append(b.ImageElements, element.(*ImageBlockElement))
		case "button":
			element, err := unmarshalBlockElement(r, &ButtonBlockElement{})
			if err != nil {
				return err
			}
			b.ButtonElements = append(b.ButtonElements, element.(*ButtonBlockElement))
		case "overflow":
			element, err := unmarshalBlockElement(r, &OverflowBlockElement{})
			if err != nil {
				return err
			}
			b.OverflowElements = append(b.OverflowElements, element.(*OverflowBlockElement))
		case "datepicker":
			element, err := unmarshalBlockElement(r, &DatePickerBlockElement{})
			if err != nil {
				return err
			}
			b.DatePickerElements = append(b.DatePickerElements, element.(*DatePickerBlockElement))
		case "static_select":
			element, err := unmarshalBlockElement(r, &SelectBlockElement{})
			if err != nil {
				return err
			}
			b.SelectElements = append(b.SelectElements, element.(*SelectBlockElement))
		}
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

func (e *BlockElements) appendToBlockElements(appendElements []BlockElement) {
	for _, element := range appendElements {
		switch element.(type) {
		case *ImageBlockElement:
			e.ImageElements = append(e.ImageElements, element.(*ImageBlockElement))
		case *ButtonBlockElement:
			e.ButtonElements = append(e.ButtonElements, element.(*ButtonBlockElement))
		case *OverflowBlockElement:
			e.OverflowElements = append(e.OverflowElements, element.(*OverflowBlockElement))
		case *DatePickerBlockElement:
			e.DatePickerElements = append(e.DatePickerElements, element.(*DatePickerBlockElement))
		case *SelectBlockElement:
			e.SelectElements = append(e.SelectElements, element.(*SelectBlockElement))
		}
	}
}

func toBlockElementSlice(elements *BlockElements) []BlockElement {
	var slice []BlockElement
	for _, element := range elements.ImageElements {
		slice = append(slice, element)
	}
	for _, element := range elements.ButtonElements {
		slice = append(slice, element)
	}
	for _, element := range elements.OverflowElements {
		slice = append(slice, element)
	}
	for _, element := range elements.DatePickerElements {
		slice = append(slice, element)
	}
	for _, element := range elements.SelectElements {
		slice = append(slice, element)
	}

	return slice
}

// Conv/JSON encoding related logic for Accessory

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

// Conv/JSON encoding related logic for ContextElements

func (e *ContextElements) MarshalJSON() ([]byte, error) {
	bytes, err := json.Marshal(toMixedElements(e))
	if err != nil {
		return nil, err
	}

	return bytes, nil
}

func toMixedElements(elements *ContextElements) []mixedElement {
	var slice []mixedElement
	for _, element := range elements.ImageElements {
		slice = append(slice, element)
	}
	for _, element := range elements.TextObjects {
		slice = append(slice, element)
	}

	return slice
}
