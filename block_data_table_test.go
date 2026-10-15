package slack

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewDataTableBlock(t *testing.T) {
	block := NewDataTableBlock("A Fabulous Table",
		DataTableBlockOptionBlockID("dt-1"),
	).
		WithPageSize(10).
		WithRowHeaderColumnIndex(1)

	assert.Equal(t, MBTDataTable, block.BlockType())
	assert.Equal(t, "data_table", string(block.Type))
	assert.Equal(t, "dt-1", block.ID())
	assert.Equal(t, "A Fabulous Table", block.Caption)
	assert.Equal(t, 10, block.PageSize)
	assert.Equal(t, 1, block.RowHeaderColumnIndex)
	assert.Empty(t, block.Rows)
}

func TestNewDataTableBlockWithNilOption(t *testing.T) {
	assert.NotPanics(t, func() {
		NewDataTableBlock("caption", nil)
	}, "should not panic when nil option passed")
}

func TestDataTableBlockAddRow(t *testing.T) {
	block := NewDataTableBlock("caption")
	block.AddRow(NewDataTableRawTextCell("Name"), NewDataTableRawTextCell("Score"))
	block.AddRow(
		NewDataTableRawTextCell("Helly"),
		NewDataTableRawNumberCell(42).WithText("forty-two"),
	)

	require.Len(t, block.Rows, 2)
	require.Len(t, block.Rows[0], 2)
	require.Len(t, block.Rows[1], 2)

	assert.Equal(t, DataTableCellRawText, block.Rows[0][0].DataTableCellType())
	assert.Equal(t, DataTableCellRawNumber, block.Rows[1][1].DataTableCellType())

	num, ok := block.Rows[1][1].(*DataTableRawNumberCell)
	require.True(t, ok)
	assert.Equal(t, float64(42), num.Value)
	assert.Equal(t, "forty-two", num.Text)
}

func TestDataTableBlockJSONRoundTrip(t *testing.T) {
	payload := `{
		"type": "data_table",
		"block_id": "dt-1",
		"caption": "A Fabulous Table",
		"page_size": 5,
		"rows": [
			[
				{"type": "raw_text", "text": "Name"},
				{"type": "raw_text", "text": "Department"},
				{"type": "raw_text", "text": "Badge"}
			],
			[
				{"type": "raw_text", "text": "Data Refinement Department"},
				{"type": "raw_text", "text": "MDR"},
				{
					"type": "rich_text",
					"elements": [
						{
							"type": "rich_text_section",
							"elements": [
								{"type": "text", "text": "Blue", "style": {"bold": true}}
							]
						}
					]
				}
			],
			[
				{"type": "raw_text", "text": "Wellness Department"},
				{"type": "raw_number", "value": 7, "text": "seven"},
				{
					"type": "rich_text",
					"elements": [
						{
							"type": "rich_text_section",
							"elements": [
								{"type": "text", "text": "Limited", "style": {"bold": true}}
							]
						}
					]
				}
			]
		]
	}`

	var block DataTableBlock
	require.NoError(t, json.Unmarshal([]byte(payload), &block))

	assert.Equal(t, MBTDataTable, block.BlockType())
	assert.Equal(t, "dt-1", block.ID())
	assert.Equal(t, "A Fabulous Table", block.Caption)
	assert.Equal(t, 5, block.PageSize)
	require.Len(t, block.Rows, 3)
	require.Len(t, block.Rows[0], 3)

	header, ok := block.Rows[0][0].(*DataTableRawTextCell)
	require.True(t, ok)
	assert.Equal(t, "Name", header.Text)

	num, ok := block.Rows[2][1].(*DataTableRawNumberCell)
	require.True(t, ok)
	assert.Equal(t, float64(7), num.Value)
	assert.Equal(t, "seven", num.Text)

	rich, ok := block.Rows[1][2].(*DataTableRichTextCell)
	require.True(t, ok)
	require.Len(t, rich.Elements, 1)
	section, ok := rich.Elements[0].(*RichTextSection)
	require.True(t, ok)
	require.Len(t, section.Elements, 1)
	text, ok := section.Elements[0].(*RichTextSectionTextElement)
	require.True(t, ok)
	assert.Equal(t, "Blue", text.Text)
	require.NotNil(t, text.Style)
	assert.True(t, text.Style.Bold)

	marshalled, err := json.Marshal(block)
	require.NoError(t, err)

	var expected, actual map[string]any
	require.NoError(t, json.Unmarshal([]byte(payload), &expected))
	require.NoError(t, json.Unmarshal(marshalled, &actual))
	assert.Equal(t, expected, actual)
}

func TestDataTableBlockUnmarshalViaBlocks(t *testing.T) {
	payload := `[
		{
			"type": "data_table",
			"caption": "Tiny Table",
			"rows": [
				[{"type": "raw_text", "text": "Col"}],
				[{"type": "raw_number", "value": 1}]
			]
		}
	]`

	var blocks Blocks
	require.NoError(t, json.Unmarshal([]byte(payload), &blocks))
	require.Len(t, blocks.BlockSet, 1)

	dt, ok := blocks.BlockSet[0].(*DataTableBlock)
	require.True(t, ok, "expected *DataTableBlock, got %T", blocks.BlockSet[0])
	assert.Equal(t, MBTDataTable, dt.BlockType())
	assert.Equal(t, "Tiny Table", dt.Caption)
	require.Len(t, dt.Rows, 2)

	num, ok := dt.Rows[1][0].(*DataTableRawNumberCell)
	require.True(t, ok)
	assert.Equal(t, float64(1), num.Value)
	assert.Empty(t, num.Text)
}

func TestDataTableBlockUnknownCellType(t *testing.T) {
	payload := `{
		"type": "data_table",
		"caption": "Bad Cell",
		"rows": [
			[{"type": "raw_text", "text": "ok"}],
			[{"type": "mystery_cell", "text": "huh"}]
		]
	}`

	var block DataTableBlock
	err := json.Unmarshal([]byte(payload), &block)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "mystery_cell")
}

func TestDataTableActionCellJSON(t *testing.T) {
	button := NewButtonBlockElement("mark_done", "task_123", &TextBlockObject{Type: PlainTextType, Text: "Mark done"})
	tests := []struct {
		name     string
		fallback DataTableCell
		payload  string
	}{
		{
			name: "no fallback",
			payload: `{"type":"action_cell","element":{
				"type":"button","text":{"type":"plain_text","text":"Mark done"},
				"action_id":"mark_done","value":"task_123"
			}}`,
		},
		{
			name:     "text fallback",
			fallback: NewDataTableRawTextCell("Open"),
			payload: `{"type":"action_cell","element":{
				"type":"button","text":{"type":"plain_text","text":"Mark done"},
				"action_id":"mark_done","value":"task_123"
			},"fallback":{"type":"raw_text","text":"Open"}}`,
		},
		{
			name:     "number fallback",
			fallback: NewDataTableRawNumberCell(0).WithText("No points"),
			payload: `{"type":"action_cell","element":{
				"type":"button","text":{"type":"plain_text","text":"Mark done"},
				"action_id":"mark_done","value":"task_123"
			},"fallback":{"type":"raw_number","value":0,"text":"No points"}}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cell := NewDataTableActionCell(button)
			assert.Nil(t, cell.Fallback)
			assert.Same(t, cell, cell.WithFallback(tt.fallback))
			assert.Equal(t, DataTableCellAction, cell.DataTableCellType())

			encoded, err := json.Marshal(cell)
			require.NoError(t, err)
			assert.JSONEq(t, tt.payload, string(encoded))

			var decoded DataTableActionCell
			require.NoError(t, json.Unmarshal([]byte(tt.payload), &decoded))
			assert.Equal(t, *cell, decoded)
		})
	}
}

func TestDataTableActionCellUnmarshalClearsFallback(t *testing.T) {
	for _, payload := range []string{
		`{"type":"action_cell","element":{"type":"button","text":{"type":"plain_text","text":"Done"}}}`,
		`{"type":"action_cell","element":{"type":"button","text":{"type":"plain_text","text":"Done"}},"fallback":null}`,
	} {
		cell := NewDataTableActionCell(nil).WithFallback(NewDataTableRawTextCell("Open"))
		require.NoError(t, json.Unmarshal([]byte(payload), cell))
		assert.Nil(t, cell.Fallback)
	}
}

func TestDataTableActionCellInvalidFallback(t *testing.T) {
	tests := []struct {
		name     string
		fallback string
	}{
		{"rich text", `{"type":"rich_text","elements":[]}`},
		{"action cell", `{"type":"action_cell"}`},
		{"unknown type", `{"type":"mystery_cell"}`},
		{"invalid shape", `"Open"`},
		{"invalid text", `{"type":"raw_text","text":42}`},
		{"invalid number", `{"type":"raw_number","value":"zero"}`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			payload := `[{"type":"data_table","caption":"Tasks","rows":[
				[{"type":"raw_text","text":"Action"}],
				[{"type":"action_cell","element":{
					"type":"button","text":{"type":"plain_text","text":"Mark done"}
				},"fallback":` + tt.fallback + `}]
			]}]`
			var blocks Blocks
			require.Error(t, json.Unmarshal([]byte(payload), &blocks))
		})
	}
}
