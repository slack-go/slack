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

	tableBlock.AddRow(&RichTextBlock{
		Type: MBTRichText,
		Elements: []RichTextElement{
			&RichTextSection{
				Type: RTESection,
				Elements: []RichTextSectionElement{
					&RichTextSectionTextElement{Type: RTSEText, Text: "Col1"},
				},
			},
		},
	}, &RichTextBlock{
		Type: MBTRichText,
		Elements: []RichTextElement{
			&RichTextSection{
				Type: RTESection,
				Elements: []RichTextSectionElement{
					&RichTextSectionTextElement{Type: RTSEText, Text: "Col2"},
				},
			},
		},
	})

	tableBlock.AddRow(&RichTextBlock{
		Type: MBTRichText,
		Elements: []RichTextElement{
			&RichTextSection{
				Type: RTESection,
				Elements: []RichTextSectionElement{
					&RichTextSectionTextElement{Type: RTSEText, Text: "Val1"},
				},
			},
		},
	}, &RichTextBlock{
		Type: MBTRichText,
		Elements: []RichTextElement{
			&RichTextSection{
				Type: RTESection,
				Elements: []RichTextSectionElement{
					&RichTextSectionTextElement{Type: RTSEText, Text: "Val2"},
				},
			},
		},
	})

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
