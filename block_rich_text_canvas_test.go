package slack

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRichTextCanvasRoundTrip(t *testing.T) {
	tests := []struct {
		name    string
		payload string
	}{
		{
			name: "minimal canvas",
			payload: `[{"type":"rich_text","elements":[{"type":"rich_text_section","elements":[
				{"type":"canvas","file_id":"F123ABC456"}
			]}]}]`,
		},
		{
			name: "all optional fields",
			payload: `[{"type":"rich_text","elements":[{"type":"rich_text_section","elements":[
				{"type":"canvas","file_id":"F123ABC456","label":"Project plan","hide_title":true,"section_id":"temp:C:abc123","text":"Plan","url":"https://example.slack.com/docs/T123/F123ABC456","is_skill_invocation":true,"style":{"bold":true,"italic":true,"strike":true,"highlight":true,"client_highlight":true,"underline":true,"unlink":true}}
			]}]}]`,
		},
		{
			name: "quote explicit false and empty style",
			payload: `[{"type":"rich_text","elements":[{"type":"rich_text_quote","elements":[
				{"type":"canvas","file_id":"F123ABC456","hide_title":false,"is_skill_invocation":false,"style":{}}
			]}]}]`,
		},
		{
			name: "live empty label",
			payload: `[{"type":"rich_text","elements":[{"type":"rich_text_section","elements":[
				{"type":"canvas","file_id":"F0CBEFSB1TJ","label":""}
			]}]}]`,
		},
		{
			name: "list canvas",
			payload: `[{"type":"rich_text","elements":[{"type":"rich_text_list","style":"bullet","indent":0,"border":0,"offset":0,"elements":[{"type":"rich_text_section","elements":[
				{"type":"canvas","file_id":"F123ABC456"}
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

func TestRichTextCanvasConstructor(t *testing.T) {
	for _, tt := range []struct {
		name  string
		style *RichTextSectionTextStyle
		want  string
	}{
		{"minimal", nil, `{"type":"canvas","file_id":"F123ABC456"}`},
		{"styled", &RichTextSectionTextStyle{Bold: true}, `{"type":"canvas","file_id":"F123ABC456","style":{"bold":true}}`},
	} {
		t.Run(tt.name, func(t *testing.T) {
			element := NewRichTextSectionCanvasElement("F123ABC456", tt.style)
			assert.Equal(t, RTSECanvas, element.RichTextSectionElementType())
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

func TestRichTextCanvasReferencesRoundTrip(t *testing.T) {
	tests := []struct {
		name    string
		payload string
	}{
		{
			name: "minimal references",
			payload: `[{"type":"rich_text","elements":[{"type":"rich_text_section","elements":[
				{"type":"canvas_user_mention","user_id":"U123ABC456"},
				{"type":"canvas_message_unfurl","root_message_channel":"C123ABC456","root_message_ts":"1720710212.123456"}
			]}]}]`,
		},
		{
			name: "all optional fields",
			payload: `[{"type":"rich_text","elements":[{"type":"rich_text_section","elements":[
				{"type":"canvas_user_mention","user_id":"U123ABC456","thread_id":"KLE9AAS5goj","style":{"bold":true,"italic":true,"strike":true,"highlight":true,"client_highlight":true,"underline":true,"unlink":true}},
				{"type":"canvas_message_unfurl","root_message_channel":"C123ABC456","root_message_ts":"1720710212.123456","style":{"bold":true,"italic":true,"strike":true,"highlight":true,"client_highlight":true,"underline":true}}
			]}]}]`,
		},
		{
			name: "quote with empty styles",
			payload: `[{"type":"rich_text","elements":[{"type":"rich_text_quote","elements":[
				{"type":"canvas_user_mention","user_id":"U123ABC456","style":{}},
				{"type":"canvas_message_unfurl","root_message_channel":"C123ABC456","root_message_ts":"1720710212.123456","style":{}}
			]}]}]`,
		},
		{
			name: "list references",
			payload: `[{"type":"rich_text","elements":[{"type":"rich_text_list","style":"bullet","indent":0,"border":0,"offset":0,"elements":[{"type":"rich_text_section","elements":[
				{"type":"canvas_user_mention","user_id":"U123ABC456","thread_id":"KLE9AAS5goj"},
				{"type":"canvas_message_unfurl","root_message_channel":"C123ABC456","root_message_ts":"1720710212.123456"}
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

func TestRichTextCanvasReferenceConstructors(t *testing.T) {
	style := &RichTextSectionTextStyle{Bold: true}
	for _, tt := range []struct {
		name    string
		element RichTextSectionElement
		want    string
	}{
		{"user", NewRichTextSectionCanvasUserMentionElement("U123ABC456", nil), `{"type":"canvas_user_mention","user_id":"U123ABC456"}`},
		{"styled user", NewRichTextSectionCanvasUserMentionElement("U123ABC456", style), `{"type":"canvas_user_mention","user_id":"U123ABC456","style":{"bold":true}}`},
		{"message", NewRichTextSectionCanvasMessageUnfurlElement("C123ABC456", "1720710212.123456", nil), `{"type":"canvas_message_unfurl","root_message_channel":"C123ABC456","root_message_ts":"1720710212.123456"}`},
		{"styled message", NewRichTextSectionCanvasMessageUnfurlElement("C123ABC456", "1720710212.123456", style), `{"type":"canvas_message_unfurl","root_message_channel":"C123ABC456","root_message_ts":"1720710212.123456","style":{"bold":true}}`},
	} {
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
