package slack

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewTaskCardBlock(t *testing.T) {
	block := NewTaskCardBlock("task-1", "Search the web", TaskCardBlockOptionBlockID("block-1"))

	assert.Equal(t, MBTTaskCard, block.BlockType())
	assert.Equal(t, "task_card", string(block.Type))
	assert.Equal(t, "block-1", block.ID())
	assert.Equal(t, "task-1", block.TaskID)
	assert.Equal(t, "Search the web", block.Title)
}

func TestTaskCardBlockChainableMethods(t *testing.T) {
	details := &RichTextBlock{
		Type: MBTRichText,
		Elements: []RichTextElement{
			&RichTextSection{
				Type: RTESection,
				Elements: []RichTextSectionElement{
					&RichTextSectionTextElement{Type: RTSEText, Text: "Searching..."},
				},
			},
		},
	}

	output := &RichTextBlock{
		Type: MBTRichText,
		Elements: []RichTextElement{
			&RichTextSection{
				Type: RTESection,
				Elements: []RichTextSectionElement{
					&RichTextSectionTextElement{Type: RTSEText, Text: "Found 3 results"},
				},
			},
		},
	}

	block := NewTaskCardBlock("task-1", "Search the web").
		WithStatus(TaskCardStatusComplete).
		WithDetails(details).
		WithOutput(output).
		WithSources(
			NewTaskCardSource("https://example.com", "Example"),
			NewTaskCardSource("https://other.com", "Other"),
		)

	assert.Equal(t, TaskCardStatusComplete, block.Status)
	assert.Equal(t, details, block.Details)
	assert.Equal(t, output, block.Output)
	assert.Len(t, block.Sources, 2)
	assert.Equal(t, "url", block.Sources[0].Type)
	assert.Equal(t, "https://example.com", block.Sources[0].URL)
	assert.Equal(t, "Example", block.Sources[0].Text)
}

func TestTaskCardBlockJSONRoundTrip(t *testing.T) {
	payload := `{
		"type": "task_card",
		"block_id": "block-1",
		"task_id": "task-1",
		"title": "Search the web",
		"status": "in_progress",
		"details": {
			"type": "rich_text",
			"elements": [
				{
					"type": "rich_text_section",
					"elements": [
						{
							"type": "text",
							"text": "Searching for results"
						}
					]
				}
			]
		},
		"output": {
			"type": "rich_text",
			"elements": [
				{
					"type": "rich_text_section",
					"elements": [
						{
							"type": "text",
							"text": "Found 3 results"
						}
					]
				}
			]
		},
		"sources": [
			{
				"type": "url",
				"url": "https://example.com",
				"text": "Example"
			}
		]
	}`

	var block TaskCardBlock
	err := json.Unmarshal([]byte(payload), &block)
	require.NoError(t, err)

	assert.Equal(t, MBTTaskCard, block.BlockType())
	assert.Equal(t, "block-1", block.ID())
	assert.Equal(t, "task-1", block.TaskID)
	assert.Equal(t, "Search the web", block.Title)
	assert.Equal(t, TaskCardStatusInProgress, block.Status)
	require.NotNil(t, block.Details)
	require.NotNil(t, block.Output)
	require.Len(t, block.Sources, 1)
	assert.Equal(t, "https://example.com", block.Sources[0].URL)

	marshalled, err := json.Marshal(block)
	require.NoError(t, err)

	var expected, actual map[string]interface{}
	err = json.Unmarshal([]byte(payload), &expected)
	require.NoError(t, err)
	err = json.Unmarshal(marshalled, &actual)
	require.NoError(t, err)

	assert.Equal(t, expected, actual)
}

func TestTaskCardBlockUnmarshalViaBlocks(t *testing.T) {
	payload := `[
		{
			"type": "task_card",
			"task_id": "task-1",
			"title": "Analyze data",
			"status": "complete"
		}
	]`

	var blocks Blocks
	err := json.Unmarshal([]byte(payload), &blocks)
	require.NoError(t, err)
	require.Len(t, blocks.BlockSet, 1)

	taskCard, ok := blocks.BlockSet[0].(*TaskCardBlock)
	require.True(t, ok, "expected *TaskCardBlock, got %T", blocks.BlockSet[0])
	assert.Equal(t, MBTTaskCard, taskCard.BlockType())
	assert.Equal(t, "task-1", taskCard.TaskID)
	assert.Equal(t, "Analyze data", taskCard.Title)
	assert.Equal(t, TaskCardStatusComplete, taskCard.Status)
}
