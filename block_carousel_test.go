package slack

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewCarouselBlock(t *testing.T) {
	cardA := NewCardBlock().WithTitle(NewTextBlockObject("mrkdwn", "A", false, false))
	cardB := NewCardBlock().WithTitle(NewTextBlockObject("mrkdwn", "B", false, false))

	block := NewCarouselBlock(cardA, cardB).WithBlockID("carousel-1")

	assert.Equal(t, MBTCarousel, block.BlockType())
	assert.Equal(t, "carousel", string(block.Type))
	assert.Equal(t, "carousel-1", block.ID())
	require.Len(t, block.Elements, 2)
	assert.Equal(t, cardA, block.Elements[0])
	assert.Equal(t, cardB, block.Elements[1])

	cardC := NewCardBlock().WithTitle(NewTextBlockObject("mrkdwn", "C", false, false))
	block.AddCard(cardC)
	require.Len(t, block.Elements, 3)
	assert.Equal(t, cardC, block.Elements[2])
}

func TestCarouselBlockJSONRoundTrip(t *testing.T) {
	payload := `{
		"type": "carousel",
		"block_id": "carousel-1",
		"elements": [
			{
				"type": "card",
				"title": {"type": "mrkdwn", "text": "MDR"},
				"body": {"type": "mrkdwn", "text": "Macrodata Refinement"}
			},
			{
				"type": "card",
				"title": {"type": "mrkdwn", "text": "O&D"},
				"body": {"type": "mrkdwn", "text": "Optics and Design"}
			}
		]
	}`

	var block CarouselBlock
	err := json.Unmarshal([]byte(payload), &block)
	require.NoError(t, err)

	assert.Equal(t, MBTCarousel, block.BlockType())
	assert.Equal(t, "carousel-1", block.ID())
	require.Len(t, block.Elements, 2)
	assert.Equal(t, "MDR", block.Elements[0].Title.Text)
	assert.Equal(t, "O&D", block.Elements[1].Title.Text)

	marshalled, err := json.Marshal(block)
	require.NoError(t, err)

	var expected, actual map[string]any
	require.NoError(t, json.Unmarshal([]byte(payload), &expected))
	require.NoError(t, json.Unmarshal(marshalled, &actual))

	assert.Equal(t, expected, actual)
}

func TestCarouselBlockUnmarshalViaBlocks(t *testing.T) {
	payload := `[
		{
			"type": "carousel",
			"elements": [
				{"type": "card", "title": {"type": "mrkdwn", "text": "Only"}}
			]
		}
	]`

	var blocks Blocks
	require.NoError(t, json.Unmarshal([]byte(payload), &blocks))
	require.Len(t, blocks.BlockSet, 1)

	carousel, ok := blocks.BlockSet[0].(*CarouselBlock)
	require.True(t, ok, "expected *CarouselBlock, got %T", blocks.BlockSet[0])
	assert.Equal(t, MBTCarousel, carousel.BlockType())
	require.Len(t, carousel.Elements, 1)
	assert.Equal(t, "Only", carousel.Elements[0].Title.Text)
}
