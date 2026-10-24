package slack

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRichTextListRecordRoundTrip(t *testing.T) {
	for _, tt := range []struct {
		name    string
		payload string
	}{
		{
			name: "minimal record",
			payload: `[{"type":"rich_text","elements":[{"type":"rich_text_section","elements":[
				{"type":"list_record","file_id":"F123ABC456"}
			]}]}]`,
		},
		{
			name: "all optional fields",
			payload: `[{"type":"rich_text","elements":[{"type":"rich_text_section","elements":[
				{"type":"list_record","file_id":"F123ABC456","record_id":"Rec123ABC456","view_id":"View123ABC456","text":"Project task","url":"https://example.slack.com/lists/T123/F123ABC456","style":{"bold":true,"italic":true,"strike":true,"highlight":true,"client_highlight":true,"underline":true,"unlink":true}}
			]}]}]`,
		},
		{
			name: "quote with empty style",
			payload: `[{"type":"rich_text","elements":[{"type":"rich_text_quote","elements":[
				{"type":"list_record","file_id":"F123ABC456","record_id":"Rec123ABC456","style":{}}
			]}]}]`,
		},
		{
			name: "list record",
			payload: `[{"type":"rich_text","elements":[{"type":"rich_text_list","style":"bullet","indent":0,"border":0,"offset":0,"elements":[{"type":"rich_text_section","elements":[
				{"type":"list_record","file_id":"F123ABC456","record_id":"Rec123ABC456","view_id":"View123ABC456"}
			]}]}]}]`,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			var blocks Blocks
			require.NoError(t, json.Unmarshal([]byte(tt.payload), &blocks))
			require.Len(t, blocks.BlockSet, 1)
			block, ok := blocks.BlockSet[0].(*RichTextBlock)
			require.True(t, ok)
			for _, element := range block.Elements {
				assertKnownRichTextFieldElements(t, element)
			}
			encoded, err := json.Marshal(blocks)
			require.NoError(t, err)
			assert.JSONEq(t, tt.payload, string(encoded))
		})
	}
}

func TestRichTextListRecordConstructor(t *testing.T) {
	for _, tt := range []struct {
		name  string
		style *RichTextSectionTextStyle
		want  string
	}{
		{"minimal", nil, `{"type":"list_record","file_id":"F123ABC456"}`},
		{"styled", &RichTextSectionTextStyle{Bold: true}, `{"type":"list_record","file_id":"F123ABC456","style":{"bold":true}}`},
	} {
		t.Run(tt.name, func(t *testing.T) {
			element := NewRichTextSectionListRecordElement("F123ABC456", tt.style)
			assert.Equal(t, RTSEListRecord, element.RichTextSectionElementType())
			encoded, err := json.Marshal(element)
			require.NoError(t, err)
			assert.JSONEq(t, tt.want, string(encoded))
			var section RichTextSection
			require.NoError(t, json.Unmarshal([]byte(`{"type":"rich_text_section","elements":[`+tt.want+`]}`), &section))
			require.Len(t, section.Elements, 1)
			assert.Equal(t, element, section.Elements[0])
		})
	}
}
