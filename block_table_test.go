package slack

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewTableBlock(t *testing.T) {
	testPayload := `{
		"type":"table",
		"block_id":"test1",
		"rows": [
			[{
				"type":"rich_text",
				"elements": [
					{
						"type":"rich_text_section",
						"elements": [
							{
								"type":"text",
								"text":"Col1"
							}
						]
					}
				]
			},
			{
				"type":"rich_text",
				"elements": [
					{
						"type":"rich_text_section",
						"elements": [
							{
								"type":"text",
								"text":"Col2"
							}
						]
					}
				]
			}],
			[{
				"type":"rich_text",
				"elements": [
					{
						"type":"rich_text_section",
						"elements": [
							{
								"type":"text",
								"text":"Val1"
							}
						]
					}
				]
			},
			{
				"type":"rich_text",
				"elements": [
					{
						"type":"rich_text_section",
						"elements": [
							{
								"type":"text",
								"text":"Val2"
							}
						]
					}
				]
			}]
		]
	}`

	tableBlock := NewTableBlock("test1")

	tableBlock.AddRow(NewTableRichTextCell(
		&RichTextSection{
			Type: RTESection,
			Elements: []RichTextSectionElement{
				&RichTextSectionTextElement{Type: RTSEText, Text: "Col1"},
			},
		},
	), NewTableRichTextCell(
		&RichTextSection{
			Type: RTESection,
			Elements: []RichTextSectionElement{
				&RichTextSectionTextElement{Type: RTSEText, Text: "Col2"},
			},
		},
	))

	tableBlock.AddRow(NewTableRichTextCell(
		&RichTextSection{
			Type: RTESection,
			Elements: []RichTextSectionElement{
				&RichTextSectionTextElement{Type: RTSEText, Text: "Val1"},
			},
		},
	), NewTableRichTextCell(
		&RichTextSection{
			Type: RTESection,
			Elements: []RichTextSectionElement{
				&RichTextSectionTextElement{Type: RTSEText, Text: "Val2"},
			},
		},
	))

	assert.Equal(t, tableBlock.BlockType(), MBTTable)
	assert.Equal(t, string(tableBlock.Type), "table")
	assert.Equal(t, tableBlock.BlockID, "test1")
	assert.Equal(t, tableBlock.ID(), "test1")
	assert.Equal(t, len(tableBlock.Rows), 2)
	assert.Equal(t, len(tableBlock.ColumnSettings), 0)

	// Check if marshalled payload matches expected JSON
	marshalled, err := json.Marshal(tableBlock)
	assert.NoError(t, err)

	var expected, actual map[string]interface{}
	err = json.Unmarshal([]byte(testPayload), &expected)
	assert.NoError(t, err)
	err = json.Unmarshal(marshalled, &actual)
	assert.NoError(t, err)

	assert.Equal(t, expected, actual)
}

// TestTableBlockHeterogeneousCells covers the case reported in issue #1558: tables
// pasted from a spreadsheet deliver header cells as rich_text and data cells as
// raw_text/raw_number, with null for empty cells. All cell text must be preserved
// across an unmarshal/marshal round trip.
func TestTableBlockHeterogeneousCells(t *testing.T) {
	payload := `{
		"type":"table",
		"block_id":"t1",
		"rows":[
			[
				{"type":"rich_text","elements":[{"type":"rich_text_section","elements":[{"type":"text","text":"Name","style":{"bold":true}}]}]},
				{"type":"rich_text","elements":[{"type":"rich_text_section","elements":[{"type":"text","text":"Score","style":{"bold":true}}]}]}
			],
			[
				{"type":"raw_text","text":"Alice"},
				{"type":"raw_number","value":42}
			],
			[
				{"type":"raw_text","text":"Bob"},
				null
			]
		]
	}`

	var tb TableBlock
	err := json.Unmarshal([]byte(payload), &tb)
	assert.NoError(t, err)
	assert.Equal(t, MBTTable, tb.BlockType())
	assert.Len(t, tb.Rows, 3)

	// Header row: rich_text cells.
	header, ok := tb.Rows[0][0].(*TableRichTextCell)
	assert.True(t, ok)
	assert.Equal(t, TableCellRichText, header.TableCellType())
	assert.Len(t, header.Elements, 1)

	// Data row: raw_text + raw_number must keep their values.
	rawText, ok := tb.Rows[1][0].(*TableRawTextCell)
	assert.True(t, ok)
	assert.Equal(t, "Alice", rawText.Text)

	rawNumber, ok := tb.Rows[1][1].(*TableRawNumberCell)
	assert.True(t, ok)
	assert.Equal(t, float64(42), rawNumber.Value)

	// Empty cell decodes as nil.
	assert.Nil(t, tb.Rows[2][1])

	// Round trip preserves the payload.
	marshalled, err := json.Marshal(&tb)
	assert.NoError(t, err)

	var expected, actual map[string]interface{}
	assert.NoError(t, json.Unmarshal([]byte(payload), &expected))
	assert.NoError(t, json.Unmarshal(marshalled, &actual))
	assert.Equal(t, expected, actual)
}

func TestTableBlockUnsupportedCellType(t *testing.T) {
	payload := `{"type":"table","rows":[[{"type":"bogus","text":"x"}]]}`
	var tb TableBlock
	err := json.Unmarshal([]byte(payload), &tb)
	assert.Error(t, err)
}
