package slack

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRichTextReferencesRoundTrip(t *testing.T) {
	tests := []struct {
		name    string
		payload string
	}{
		{
			name: "minimal references",
			payload: `[{"type":"rich_text","elements":[{"type":"rich_text_section","elements":[
				{"type":"file","file_id":"F123ABC456"},
				{"type":"message_mention","channel_id":"C123ABC456","message_ts":"1720710212.123456"}
			]}]}]`,
		},
		{
			name: "all optional fields",
			payload: `[{"type":"rich_text","elements":[{"type":"rich_text_section","elements":[
				{"type":"file","file_id":"F123ABC456","text":"Report","url":"https://example.com/report","is_skill_invocation":true,"style":{"bold":true,"italic":true,"strike":true,"highlight":true,"client_highlight":true,"underline":true}},
				{"type":"message_mention","channel_id":"C123ABC456","message_ts":"1720710212.123456","author_id":"U123ABC456","thread_ts":"1720710200.000001","text":"Discussion","url":"https://example.slack.com/archives/C123ABC456/p1720710212123456","style":{"bold":true,"italic":true,"strike":true,"highlight":true,"client_highlight":true,"underline":true,"unlink":true}}
			]}]}]`,
		},
		{
			name: "quote explicit false and empty style",
			payload: `[{"type":"rich_text","elements":[{"type":"rich_text_quote","elements":[
				{"type":"file","file_id":"F123ABC456","is_skill_invocation":false,"style":{}},
				{"type":"message_mention","channel_id":"C123ABC456","message_ts":"1720710212.123456","style":{}}
			]}]}]`,
		},
		{
			name: "list references",
			payload: `[{"type":"rich_text","elements":[{"type":"rich_text_list","style":"bullet","indent":0,"border":0,"offset":0,"elements":[{"type":"rich_text_section","elements":[
				{"type":"file","file_id":"F123ABC456"},
				{"type":"message_mention","channel_id":"C123ABC456","message_ts":"1720710212.123456"}
			]}]}]}]`,
		},
	}
	for _, tt := range tests {
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

func TestRichTextReferenceConstructors(t *testing.T) {
	style := &RichTextSectionTextStyle{Bold: true}
	tests := []struct {
		name    string
		element RichTextSectionElement
		want    string
	}{
		{"file", NewRichTextSectionFileElement("F123ABC456", nil), `{"type":"file","file_id":"F123ABC456"}`},
		{"styled file", NewRichTextSectionFileElement("F123ABC456", style), `{"type":"file","file_id":"F123ABC456","style":{"bold":true}}`},
		{"message", NewRichTextSectionMessageMentionElement("C123ABC456", "1720710212.123456", nil), `{"type":"message_mention","channel_id":"C123ABC456","message_ts":"1720710212.123456"}`},
		{"styled message", NewRichTextSectionMessageMentionElement("C123ABC456", "1720710212.123456", style), `{"type":"message_mention","channel_id":"C123ABC456","message_ts":"1720710212.123456","style":{"bold":true}}`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			encoded, err := json.Marshal(tt.element)
			require.NoError(t, err)
			assert.JSONEq(t, tt.want, string(encoded))
			var section RichTextSection
			require.NoError(t, json.Unmarshal([]byte(`{"type":"rich_text_section","elements":[`+tt.want+`]}`), &section))
			require.Len(t, section.Elements, 1)
			assert.Equal(t, tt.element, section.Elements[0])
			assert.Equal(t, tt.element.RichTextSectionElementType(), section.Elements[0].RichTextSectionElementType())
		})
	}
}
