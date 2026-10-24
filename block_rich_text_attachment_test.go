package slack

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRichTextAttachmentMentionRoundTrip(t *testing.T) {
	for _, tt := range []struct {
		name    string
		payload string
	}{
		{
			name: "minimal mention",
			payload: `[{"type":"rich_text","elements":[{"type":"rich_text_section","elements":[
				{"type":"attachment_mention","url":"https://example.com/attachment"}
			]}]}]`,
		},
		{
			name: "live assigned message location",
			payload: `[{"type":"rich_text","block_id":"reference_attachment-minimal","elements":[{"type":"rich_text_section","elements":[
				{"type":"attachment_mention","url":"https://docs.slack.dev/reference/block-kit/block-elements/attachment-mention-element/","channel_id":"C090051UMUP","ts":"1792874180.038299"}
			]}]}]`,
		},
		{
			name: "all optional fields",
			payload: `[{"type":"rich_text","elements":[{"type":"rich_text_section","elements":[
				{"type":"attachment_mention","url":"https://example.com/attachment","text":"Project document","app_id":"A123ABC456","entity_id":"document-123","icon_url":"https://example.com/icon.png","channel_id":"C123ABC456","ts":"1720710212.123456","full_size_preview_enabled":true,"icon_name":"document","reference_object_type":"record","product_name":"Documents","style":{"bold":true,"italic":true,"strike":true,"highlight":true,"client_highlight":true,"underline":true,"unlink":true}}
			]}]}]`,
		},
		{
			name: "quote explicit false and empty style",
			payload: `[{"type":"rich_text","elements":[{"type":"rich_text_quote","elements":[
				{"type":"attachment_mention","url":"https://example.com/attachment","full_size_preview_enabled":false,"style":{}}
			]}]}]`,
		},
		{
			name: "list mention",
			payload: `[{"type":"rich_text","elements":[{"type":"rich_text_list","style":"bullet","indent":0,"border":0,"offset":0,"elements":[{"type":"rich_text_section","elements":[
				{"type":"attachment_mention","url":"https://example.com/attachment","app_id":"A123ABC456","channel_id":"C123ABC456","ts":"1720710212.123456"}
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

func TestRichTextAttachmentMentionConstructor(t *testing.T) {
	for _, tt := range []struct {
		name  string
		style *RichTextSectionTextStyle
		want  string
	}{
		{"minimal", nil, `{"type":"attachment_mention","url":"https://example.com/attachment"}`},
		{"styled", &RichTextSectionTextStyle{Bold: true}, `{"type":"attachment_mention","url":"https://example.com/attachment","style":{"bold":true}}`},
	} {
		t.Run(tt.name, func(t *testing.T) {
			element := NewRichTextSectionAttachmentMentionElement("https://example.com/attachment", tt.style)
			assert.Equal(t, RTSEAttachmentMention, element.RichTextSectionElementType())
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
