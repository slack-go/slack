package slack

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewContextActionsBlock(t *testing.T) {
	positiveBtnText := NewTextBlockObject("plain_text", "Good", false, false)
	negativeBtnText := NewTextBlockObject("plain_text", "Bad", false, false)
	positiveBtn := NewFeedbackButton(positiveBtnText, "positive_feedback")
	negativeBtn := NewFeedbackButton(negativeBtnText, "negative_feedback")
	feedbackElement := NewFeedbackButtonsBlockElement("feedback_1", positiveBtn, negativeBtn)

	contextActionsBlock := NewContextActionsBlock("test_block", feedbackElement)

	assert.Equal(t, contextActionsBlock.BlockType(), MBTContextActions)
	assert.Equal(t, string(contextActionsBlock.Type), "context_actions")
	assert.Equal(t, contextActionsBlock.BlockID, "test_block")
	assert.Equal(t, contextActionsBlock.ID(), "test_block")
	assert.Equal(t, len(contextActionsBlock.Elements.ElementSet), 1)
}

func TestContextActionsBlockWithIconButton(t *testing.T) {
	deleteText := NewTextBlockObject("plain_text", "Delete", false, false)
	iconButton := NewIconButtonBlockElement("trash", deleteText, "delete_action")

	contextActionsBlock := NewContextActionsBlock("icon_block", iconButton)

	assert.Equal(t, contextActionsBlock.BlockType(), MBTContextActions)
	assert.Equal(t, string(contextActionsBlock.Type), "context_actions")
	assert.Equal(t, contextActionsBlock.BlockID, "icon_block")
	assert.Equal(t, len(contextActionsBlock.Elements.ElementSet), 1)
}

func TestContextActionsBlockWithMultipleElements(t *testing.T) {
	// Create feedback buttons
	positiveBtnText := NewTextBlockObject("plain_text", "👍", false, false)
	negativeBtnText := NewTextBlockObject("plain_text", "👎", false, false)
	positiveBtn := NewFeedbackButton(positiveBtnText, "positive")
	negativeBtn := NewFeedbackButton(negativeBtnText, "negative")
	feedbackElement := NewFeedbackButtonsBlockElement("feedback_1", positiveBtn, negativeBtn)

	// Create icon button
	deleteText := NewTextBlockObject("plain_text", "Delete", false, false)
	iconButton := NewIconButtonBlockElement("trash", deleteText, "delete_action")

	contextActionsBlock := NewContextActionsBlock("multi_block", feedbackElement, iconButton)

	assert.Equal(t, contextActionsBlock.BlockType(), MBTContextActions)
	assert.Equal(t, len(contextActionsBlock.Elements.ElementSet), 2)
}

func TestContextActionsBlockJSONMarshalling(t *testing.T) {
	positiveBtnText := NewTextBlockObject("plain_text", "Good", false, false)
	negativeBtnText := NewTextBlockObject("plain_text", "Bad", false, false)
	positiveBtn := NewFeedbackButton(positiveBtnText, "positive_feedback")
	negativeBtn := NewFeedbackButton(negativeBtnText, "negative_feedback")
	feedbackElement := NewFeedbackButtonsBlockElement("feedback_buttons_1", positiveBtn, negativeBtn)

	contextActionsBlock := NewContextActionsBlock("test_block", feedbackElement)

	// Marshal to JSON
	data, err := json.Marshal(contextActionsBlock)
	assert.NoError(t, err)
	assert.NotNil(t, data)

	// Unmarshal back
	var unmarshalled ContextActionsBlock
	err = json.Unmarshal(data, &unmarshalled)
	assert.NoError(t, err)
	assert.Equal(t, "context_actions", string(unmarshalled.Type))
	assert.Equal(t, "test_block", unmarshalled.BlockID)
	assert.Equal(t, 1, len(unmarshalled.Elements.ElementSet))
}

func TestContextActionsBlockUnmarshalJSON(t *testing.T) {
	jsonData := []byte(`{
		"type": "context_actions",
		"block_id": "test_block",
		"elements": [
			{
				"type": "feedback_buttons",
				"action_id": "feedback_buttons_1",
				"positive_button": {
					"text": {
						"type": "plain_text",
						"text": "Good"
					},
					"value": "positive_feedback"
				},
				"negative_button": {
					"text": {
						"type": "plain_text",
						"text": "Bad"
					},
					"value": "negative_feedback"
				}
			}
		]
	}`)

	var block ContextActionsBlock
	err := json.Unmarshal(jsonData, &block)
	assert.NoError(t, err)
	assert.Equal(t, "context_actions", string(block.Type))
	assert.Equal(t, "test_block", block.BlockID)
	assert.Equal(t, 1, len(block.Elements.ElementSet))
}

func TestContextActionsBlockInBlocks(t *testing.T) {
	// Test that context_actions block can be unmarshalled as part of a Blocks collection
	jsonData := []byte(`[
		{
			"type": "context_actions",
			"block_id": "actions_block",
			"elements": [
				{
					"type": "icon_button",
					"icon": "trash",
					"text": {
						"type": "plain_text",
						"text": "Delete"
					},
					"action_id": "delete_button_1",
					"value": "delete_item"
				}
			]
		}
	]`)

	var blocks Blocks
	err := json.Unmarshal(jsonData, &blocks)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(blocks.BlockSet))
	assert.Equal(t, MBTContextActions, blocks.BlockSet[0].BlockType())
}
