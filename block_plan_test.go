package slack

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewPlanBlock(t *testing.T) {
	block := NewPlanBlock("My Plan", PlanBlockOptionBlockID("plan-block-1"))

	assert.Equal(t, MBTPlan, block.BlockType())
	assert.Equal(t, "plan", string(block.Type))
	assert.Equal(t, "plan-block-1", block.ID())
	assert.Equal(t, "My Plan", block.Title)
}

func TestPlanBlockWithTasks(t *testing.T) {
	task1 := NewTaskCardBlock("task-1", "First task").WithStatus(TaskCardStatusComplete)
	task2 := NewTaskCardBlock("task-2", "Second task").WithStatus(TaskCardStatusInProgress)

	block := NewPlanBlock("My Plan").WithTasks(task1, task2)

	require.Len(t, block.Tasks, 2)
	assert.Equal(t, "task-1", block.Tasks[0].TaskID)
	assert.Equal(t, "First task", block.Tasks[0].Title)
	assert.Equal(t, TaskCardStatusComplete, block.Tasks[0].Status)
	assert.Equal(t, "task-2", block.Tasks[1].TaskID)
	assert.Equal(t, TaskCardStatusInProgress, block.Tasks[1].Status)
}

func TestPlanBlockJSONRoundTrip(t *testing.T) {
	payload := `{
		"type": "plan",
		"block_id": "plan-1",
		"title": "Research Plan",
		"tasks": [
			{
				"type": "task_card",
				"task_id": "task-1",
				"title": "Search the web",
				"status": "complete",
				"sources": [
					{
						"type": "url",
						"url": "https://example.com",
						"text": "Example"
					}
				]
			},
			{
				"type": "task_card",
				"task_id": "task-2",
				"title": "Analyze results",
				"status": "in_progress",
				"details": {
					"type": "rich_text",
					"elements": [
						{
							"type": "rich_text_section",
							"elements": [
								{
									"type": "text",
									"text": "Processing data..."
								}
							]
						}
					]
				}
			}
		]
	}`

	var block PlanBlock
	err := json.Unmarshal([]byte(payload), &block)
	require.NoError(t, err)

	assert.Equal(t, MBTPlan, block.BlockType())
	assert.Equal(t, "plan-1", block.ID())
	assert.Equal(t, "Research Plan", block.Title)
	require.Len(t, block.Tasks, 2)
	assert.Equal(t, "task-1", block.Tasks[0].TaskID)
	assert.Equal(t, TaskCardStatusComplete, block.Tasks[0].Status)
	require.Len(t, block.Tasks[0].Sources, 1)
	assert.Equal(t, "task-2", block.Tasks[1].TaskID)
	require.NotNil(t, block.Tasks[1].Details)

	marshalled, err := json.Marshal(block)
	require.NoError(t, err)

	var expected, actual map[string]interface{}
	err = json.Unmarshal([]byte(payload), &expected)
	require.NoError(t, err)
	err = json.Unmarshal(marshalled, &actual)
	require.NoError(t, err)

	assert.Equal(t, expected, actual)
}

func TestPlanBlockUnmarshalViaBlocks(t *testing.T) {
	payload := `[
		{
			"type": "plan",
			"title": "Agent Plan",
			"tasks": [
				{
					"type": "task_card",
					"task_id": "t1",
					"title": "Step 1",
					"status": "pending"
				}
			]
		}
	]`

	var blocks Blocks
	err := json.Unmarshal([]byte(payload), &blocks)
	require.NoError(t, err)
	require.Len(t, blocks.BlockSet, 1)

	plan, ok := blocks.BlockSet[0].(*PlanBlock)
	require.True(t, ok, "expected *PlanBlock, got %T", blocks.BlockSet[0])
	assert.Equal(t, MBTPlan, plan.BlockType())
	assert.Equal(t, "Agent Plan", plan.Title)
	require.Len(t, plan.Tasks, 1)
	assert.Equal(t, "t1", plan.Tasks[0].TaskID)
	assert.Equal(t, TaskCardStatusPending, plan.Tasks[0].Status)
}
