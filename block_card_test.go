package slack

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewCardBlock(t *testing.T) {
	title := NewTextBlockObject("mrkdwn", "Card title", false, false)
	subtitle := NewTextBlockObject("mrkdwn", "Card subtitle", false, false)
	body := NewTextBlockObject("mrkdwn", "Card body text.", false, false)

	iconURL := "https://example.com/icon.png"
	icon := &ImageBlockElement{Type: METImage, ImageURL: &iconURL, AltText: "icon"}
	heroURL := "https://example.com/hero.png"
	hero := &ImageBlockElement{Type: METImage, ImageURL: &heroURL, AltText: "hero"}

	btnText := NewTextBlockObject("plain_text", "Open", false, false)
	btn := NewButtonBlockElement("open_action", "go", btnText)

	block := NewCardBlock(CardBlockOptionBlockID("card-1")).
		WithTitle(title).
		WithSubtitle(subtitle).
		WithBody(body).
		WithIcon(icon).
		WithHeroImage(hero).
		WithActions(btn)

	assert.Equal(t, MBTCard, block.BlockType())
	assert.Equal(t, "card", string(block.Type))
	assert.Equal(t, "card-1", block.ID())
	assert.Equal(t, title, block.Title)
	assert.Equal(t, subtitle, block.Subtitle)
	assert.Equal(t, body, block.Body)
	assert.Equal(t, icon, block.Icon)
	assert.Equal(t, hero, block.HeroImage)
	require.NotNil(t, block.Actions)
	require.Len(t, block.Actions.ElementSet, 1)
}

func TestCardBlockJSONRoundTrip(t *testing.T) {
	payload := `{
		"type": "card",
		"block_id": "card-1",
		"hero_image": {
			"type": "image",
			"image_url": "https://example.com/hero.png",
			"alt_text": "hero"
		},
		"icon": {
			"type": "image",
			"image_url": "https://example.com/icon.png",
			"alt_text": "icon"
		},
		"title": {"type": "mrkdwn", "text": "Lumon Industries"},
		"subtitle": {"type": "mrkdwn", "text": "Macrodata Refinement"},
		"body": {"type": "mrkdwn", "text": "The work is mysterious and important."},
		"actions": [
			{
				"type": "button",
				"text": {"type": "plain_text", "text": "Enter"},
				"action_id": "enter",
				"value": "mdr"
			}
		]
	}`

	var block CardBlock
	err := json.Unmarshal([]byte(payload), &block)
	require.NoError(t, err)

	assert.Equal(t, MBTCard, block.BlockType())
	assert.Equal(t, "card-1", block.ID())
	require.NotNil(t, block.HeroImage)
	require.NotNil(t, block.Icon)
	require.NotNil(t, block.Title)
	require.NotNil(t, block.Subtitle)
	require.NotNil(t, block.Body)
	require.NotNil(t, block.Actions)
	require.Len(t, block.Actions.ElementSet, 1)

	btn, ok := block.Actions.ElementSet[0].(*ButtonBlockElement)
	require.True(t, ok, "expected *ButtonBlockElement, got %T", block.Actions.ElementSet[0])
	assert.Equal(t, "enter", btn.ActionID)

	marshalled, err := json.Marshal(block)
	require.NoError(t, err)

	var expected, actual map[string]interface{}
	require.NoError(t, json.Unmarshal([]byte(payload), &expected))
	require.NoError(t, json.Unmarshal(marshalled, &actual))

	assert.Equal(t, expected, actual)
}

func TestCardBlockUnmarshalViaBlocks(t *testing.T) {
	payload := `[
		{
			"type": "card",
			"title": {"type": "mrkdwn", "text": "Hello"}
		}
	]`

	var blocks Blocks
	require.NoError(t, json.Unmarshal([]byte(payload), &blocks))
	require.Len(t, blocks.BlockSet, 1)

	card, ok := blocks.BlockSet[0].(*CardBlock)
	require.True(t, ok, "expected *CardBlock, got %T", blocks.BlockSet[0])
	assert.Equal(t, MBTCard, card.BlockType())
	require.NotNil(t, card.Title)
	assert.Equal(t, "Hello", card.Title.Text)
}
