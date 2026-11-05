package slack

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestContextBlockMarshalJSON(t *testing.T) {
	block := NewContextBlock("context", NewTextBlockObject(PlainTextType, "hello", true, false))
	for _, b := range []Block{block, *block} {
		data, err := json.Marshal(Blocks{BlockSet: []Block{b}})
		require.NoError(t, err)
		assert.JSONEq(t, `[{"type":"context","block_id":"context","elements":[{"type":"plain_text","text":"hello","emoji":true}]}]`, string(data))
	}
}

func TestNewContextBlock(t *testing.T) {
	locationPinImage := NewImageBlockElement("https://api.slack.com/img/blocks/bkb_template_images/tripAgentLocationMarker.png", "Location Pin Icon")
	textExample := NewTextBlockObject("plain_text", "Location: Central Business District", true, false)
	elements := []MixedElement{locationPinImage, textExample}
	contextBlock := NewContextBlock("test", elements...)

	assert.Equal(t, MBTContext, contextBlock.BlockType())
	assert.Equal(t, "context", string(contextBlock.Type))
	assert.Equal(t, "test", contextBlock.BlockID)
	assert.Equal(t, "test", contextBlock.ID())
	assert.Len(t, contextBlock.ContextElements.Elements, 2)
}
