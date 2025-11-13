package slack

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBlockFromJSON(t *testing.T) {
	tests := []struct {
		name      string
		json      string
		wantType  MessageBlockType
		wantError bool
	}{
		{
			name:     "valid single divider block",
			json:     `{"type": "divider"}`,
			wantType: MBTDivider,
		},
		{
			name:     "valid section block",
			json:     `{"type": "section", "text": {"type": "mrkdwn", "text": "Hello"}}`,
			wantType: MBTSection,
		},
		{
			name:     "valid array with single block",
			json:     `[{"type": "divider"}]`,
			wantType: MBTDivider,
		},
		{
			name:     "valid array with multiple blocks (takes first)",
			json:     `[{"type": "divider"}, {"type": "section", "text": {"type": "plain_text", "text": "Hi"}}]`,
			wantType: MBTDivider,
		},
		{
			name:      "invalid JSON syntax",
			json:      `{"type": "divider"`,
			wantError: true,
		},
		{
			name:      "empty JSON object",
			json:      `{}`,
			wantError: true, // Cannot determine block type without wrapping in array
		},
		{
			name:      "empty array",
			json:      `[]`,
			wantError: true,
		},
		{
			name:      "null",
			json:      `null`,
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			block, err := BlockFromJSON(tt.json)

			if tt.wantError {
				assert.Error(t, err)
				assert.Nil(t, block)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, block)
				assert.Equal(t, tt.wantType, block.BlockType())
			}
		})
	}
}

func TestMustBlockFromJSON(t *testing.T) {
	t.Run("valid JSON does not panic", func(t *testing.T) {
		assert.NotPanics(t, func() {
			block := MustBlockFromJSON(`{"type": "divider"}`)
			assert.NotNil(t, block)
			assert.Equal(t, MBTDivider, block.BlockType())
		})
	})

	t.Run("invalid JSON panics", func(t *testing.T) {
		assert.Panics(t, func() {
			MustBlockFromJSON(`invalid json`)
		})
	})

	t.Run("empty array panics", func(t *testing.T) {
		assert.Panics(t, func() {
			MustBlockFromJSON(`[]`)
		})
	})
}

func TestRawJSONBlockRoundTrip(t *testing.T) {
	t.Run("simple block preserves all fields", func(t *testing.T) {
		originalJSON := `{"type": "section", "text": {"type": "mrkdwn", "text": "Hello World"}, "block_id": "section1"}`

		block, err := BlockFromJSON(originalJSON)
		assert.NoError(t, err)
		assert.Equal(t, MBTSection, block.BlockType())
		assert.Equal(t, "section1", block.ID())

		// Marshal back to JSON
		marshalled, err := json.Marshal(block)
		assert.NoError(t, err)

		// Unmarshal both to compare (ignoring whitespace differences)
		var original, result map[string]interface{}
		assert.NoError(t, json.Unmarshal([]byte(originalJSON), &original))
		assert.NoError(t, json.Unmarshal(marshalled, &result))

		assert.Equal(t, original, result, "Round-trip should preserve all fields")
	})

	t.Run("complex block with nested elements preserves everything", func(t *testing.T) {
		// Using a complex context_actions block as an example
		originalJSON := `{
			"type": "context_actions",
			"block_id": "feedback_block",
			"elements": [
				{
					"type": "feedback_buttons",
					"action_id": "ai_feedback",
					"positive_button": {
						"text": {"type": "plain_text", "text": "👍"},
						"value": "positive"
					},
					"negative_button": {
						"text": {"type": "plain_text", "text": "👎"},
						"value": "negative"
					}
				},
				{
					"type": "icon_button",
					"icon": "trash",
					"text": {"type": "plain_text", "text": "Delete"},
					"action_id": "delete_action",
					"value": "delete_response"
				}
			]
		}`

		block, err := BlockFromJSON(originalJSON)
		assert.NoError(t, err)
		assert.Equal(t, MessageBlockType("context_actions"), block.BlockType())
		assert.Equal(t, "feedback_block", block.ID())

		// Marshal back to JSON
		marshalled, err := json.Marshal(block)
		assert.NoError(t, err)

		// Unmarshal both to compare
		var original, result map[string]interface{}
		assert.NoError(t, json.Unmarshal([]byte(originalJSON), &original))
		assert.NoError(t, json.Unmarshal(marshalled, &result))

		assert.Equal(t, original, result, "Complex block should preserve all nested fields")

		// Specifically verify elements array is preserved
		resultElements, ok := result["elements"].([]interface{})
		assert.True(t, ok, "elements should be an array")
		assert.Equal(t, 2, len(resultElements), "should have 2 elements")
	})
}

func TestRawJSONBlockInDeepStructure(t *testing.T) {
	t.Run("RawJSONBlock in Message with mixed blocks", func(t *testing.T) {
		// Create a regular divider block
		divider := NewDividerBlock()

		// Create a RawJSONBlock with complex content
		contextActionsJSON := `{
			"type": "context_actions",
			"block_id": "feedback_block",
			"elements": [
				{
					"type": "feedback_buttons",
					"action_id": "ai_feedback",
					"positive_button": {
						"text": {"type": "plain_text", "text": "👍"},
						"value": "positive"
					},
					"negative_button": {
						"text": {"type": "plain_text", "text": "👎"},
						"value": "negative"
					}
				}
			]
		}`
		rawBlock := MustBlockFromJSON(contextActionsJSON)

		// Create a regular section block
		sectionText := NewTextBlockObject("mrkdwn", "*Regular Section*", false, false)
		section := NewSectionBlock(sectionText, nil, nil)

		// Create a Message with all three blocks
		msg := NewBlockMessage(divider, rawBlock, section)

		// Marshal the entire message
		marshalled, err := json.Marshal(msg)
		assert.NoError(t, err)

		// Unmarshal to verify structure
		var result struct {
			Blocks []json.RawMessage `json:"blocks"`
		}
		assert.NoError(t, json.Unmarshal(marshalled, &result))
		assert.Equal(t, 3, len(result.Blocks), "should have 3 blocks")

		// Verify the middle block (our RawJSONBlock) has all fields preserved
		var contextActionsBlock map[string]interface{}
		assert.NoError(t, json.Unmarshal(result.Blocks[1], &contextActionsBlock))

		assert.Equal(t, "context_actions", contextActionsBlock["type"])
		assert.Equal(t, "feedback_block", contextActionsBlock["block_id"])

		// Verify elements array is present and intact
		elements, ok := contextActionsBlock["elements"].([]interface{})
		assert.True(t, ok, "elements should be an array")
		assert.Equal(t, 1, len(elements), "should have 1 element")

		// Verify the nested feedback_buttons element
		firstElement, ok := elements[0].(map[string]interface{})
		assert.True(t, ok, "first element should be an object")
		assert.Equal(t, "feedback_buttons", firstElement["type"])
		assert.Equal(t, "ai_feedback", firstElement["action_id"])

		// Verify the nested buttons exist
		assert.NotNil(t, firstElement["positive_button"], "should have positive_button")
		assert.NotNil(t, firstElement["negative_button"], "should have negative_button")
	})
}
