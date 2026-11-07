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
	require.NoError(t, json.Unmarshal([]byte(`[{"type":"plain_text","text":"replacement"}]`), &contextBlock.ContextElements))
	require.Len(t, contextBlock.ContextElements.Elements, 1)
	assert.Equal(t, "replacement", contextBlock.ContextElements.Elements[0].(*TextBlockObject).Text)
	require.NoError(t, json.Unmarshal([]byte(`[]`), &contextBlock.ContextElements))
	assert.Empty(t, contextBlock.ContextElements.Elements)
}

func TestContextElementsUnmarshalReuse(t *testing.T) {
	const original = `[{"type":"image","image_url":"https://example.com/old.png","alt_text":"old image"},{"type":"plain_text","text":"old text"}]`
	for _, input := range []string{
		`[{"type":"image","image_url":"https://example.com/new.png","alt_text":"new image"}]`,
		`[{"type":"plain_text","text":"replacement","emoji":true}]`,
		`[{"type":"mrkdwn","text":"*replacement*","verbatim":true}]`,
		`[]`,
		`null`,
	} {
		t.Run(input, func(t *testing.T) {
			var reused, fresh ContextElements
			require.NoError(t, json.Unmarshal([]byte(original), &reused))
			require.NoError(t, json.Unmarshal([]byte(input), &reused))
			require.NoError(t, json.Unmarshal([]byte(input), &fresh))
			require.Equal(t, fresh, reused)
			encoded, err := json.Marshal(reused)
			require.NoError(t, err)
			if input == `[]` || input == `null` {
				require.Empty(t, reused.Elements)
				require.JSONEq(t, `null`, string(encoded))
			} else {
				require.JSONEq(t, input, string(encoded))
			}
		})
	}
}

func TestContextElementsUnmarshalLegacyNull(t *testing.T) {
	elements := ContextElements{Elements: []MixedElement{NewTextBlockObject(PlainTextType, "keep", false, false)}}
	previous := elements
	require.NoError(t, json.Unmarshal([]byte(`{"elements":null}`), &elements))
	require.Equal(t, previous, elements)
}

func TestContextElementsUnmarshalInvalidInput(t *testing.T) {
	for _, input := range []string{
		`[`, `42`, `[42]`, `[{"type":42}]`, `[{"type":"unsupported"}]`,
		`[{"type":"plain_text","text":42}]`,
		`[{"type":"image","image_url":42}]`,
		`[{"type":"plain_text","text":"valid first element"},{"type":"image","alt_text":42}]`,
	} {
		t.Run(input, func(t *testing.T) {
			elements := ContextElements{Elements: []MixedElement{NewTextBlockObject(PlainTextType, "old", false, false)}}
			require.Error(t, json.Unmarshal([]byte(input), &elements))
		})
	}
}
