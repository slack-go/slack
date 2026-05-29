package slack

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewAlertBlock(t *testing.T) {
	text := NewTextBlockObject("mrkdwn", "The work is mysterious and important.", false, false)
	block := NewAlertBlock(text,
		AlertBlockOptionLevel(AlertLevelInfo),
		AlertBlockOptionBlockID("alert-1"),
	)

	assert.Equal(t, MBTAlert, block.BlockType())
	assert.Equal(t, "alert", string(block.Type))
	assert.Equal(t, "alert-1", block.ID())
	assert.Equal(t, AlertLevelInfo, block.Level)
	assert.Equal(t, text, block.Text)
}

func TestNewAlertBlockWithNilOption(t *testing.T) {
	text := NewTextBlockObject("plain_text", "hi", false, false)
	assert.NotPanics(t, func() {
		NewAlertBlock(text, nil)
	}, "should not panic when nil option passed")
}

func TestAlertBlockJSONRoundTrip(t *testing.T) {
	payload := `{
		"type": "alert",
		"text": {
			"type": "mrkdwn",
			"text": "The work is mysterious and important."
		},
		"level": "info",
		"block_id": "alert-1"
	}`

	var block AlertBlock
	err := json.Unmarshal([]byte(payload), &block)
	require.NoError(t, err)

	assert.Equal(t, MBTAlert, block.BlockType())
	assert.Equal(t, "alert-1", block.ID())
	assert.Equal(t, AlertLevelInfo, block.Level)
	require.NotNil(t, block.Text)
	assert.Equal(t, "mrkdwn", block.Text.Type)

	marshalled, err := json.Marshal(block)
	require.NoError(t, err)

	var expected, actual map[string]any
	require.NoError(t, json.Unmarshal([]byte(payload), &expected))
	require.NoError(t, json.Unmarshal(marshalled, &actual))

	assert.Equal(t, expected, actual)
}

func TestAlertBlockUnmarshalViaBlocks(t *testing.T) {
	payload := `[
		{
			"type": "alert",
			"text": {"type": "plain_text", "text": "Heads up"},
			"level": "warning"
		}
	]`

	var blocks Blocks
	require.NoError(t, json.Unmarshal([]byte(payload), &blocks))
	require.Len(t, blocks.BlockSet, 1)

	alert, ok := blocks.BlockSet[0].(*AlertBlock)
	require.True(t, ok, "expected *AlertBlock, got %T", blocks.BlockSet[0])
	assert.Equal(t, MBTAlert, alert.BlockType())
	assert.Equal(t, AlertLevelWarning, alert.Level)
}
