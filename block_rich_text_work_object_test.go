package slack

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRichTextWorkObjectMentionRoundTrip(t *testing.T) {
	for _, tt := range []struct {
		name    string
		payload string
	}{
		{
			name: "minimal mention",
			payload: `[{"type":"rich_text","elements":[{"type":"rich_text_section","elements":[
				{"type":"work_object_mention","entity_id":"E123ABC456","app_id":"A123ABC456","text":"Project document","url":"https://example.com/document"}
			]}]}]`,
		},
		{
			name: "live app-scoped entity ID",
			payload: `[{"type":"rich_text","block_id":"reference_work-object-minimal","elements":[{"type":"rich_text_section","elements":[
				{"type":"work_object_mention","entity_id":"W:A0ABXHWB1H6:slack-go-work-object-mention-docs","app_id":"A0ABXHWB1H6","text":"Work Object ID check","url":"https://docs.slack.dev/reference/block-kit/block-elements/work-object-mention-element/"}
			]}]}]`,
		},
		{
			name: "live optional fields and explicit false",
			payload: `[{"type":"rich_text","block_id":"reference_work-object-full","elements":[{"type":"rich_text_section","elements":[
				{"type":"work_object_mention","entity_id":"W:A0ABXHWB1H6:slack-go-work-object-mention-docs","app_id":"A0ABXHWB1H6","text":"Work Object full fields","url":"https://docs.slack.dev/reference/block-kit/block-elements/work-object-mention-element/","icon_url":"https://a.slack-edge.com/80588/img/plugins/app/bot_48.png","full_size_preview_enabled":false,"product_name":"Slack Developer Docs","style":{"bold":true,"italic":true,"underline":true}}
			]}]}]`,
		},
		{
			name: "all optional fields",
			payload: `[{"type":"rich_text","elements":[{"type":"rich_text_section","elements":[
				{"type":"work_object_mention","entity_id":"E123ABC456","app_id":"A123ABC456","text":"Project document","url":"https://example.com/document","icon_url":"https://example.com/icon.png","full_size_preview_enabled":true,"product_name":"Documents","style":{"bold":true,"italic":true,"strike":true,"highlight":true,"client_highlight":true,"underline":true,"unlink":true}}
			]}]}]`,
		},
		{
			name: "quote explicit false and empty style",
			payload: `[{"type":"rich_text","elements":[{"type":"rich_text_quote","elements":[
				{"type":"work_object_mention","entity_id":"E123ABC456","app_id":"A123ABC456","text":"Project document","url":"https://example.com/document","full_size_preview_enabled":false,"style":{}}
			]}]}]`,
		},
		{
			name: "list mention",
			payload: `[{"type":"rich_text","elements":[{"type":"rich_text_list","style":"bullet","indent":0,"border":0,"offset":0,"elements":[{"type":"rich_text_section","elements":[
				{"type":"work_object_mention","entity_id":"E123ABC456","app_id":"A123ABC456","text":"Project document","url":"https://example.com/document"}
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

func TestRichTextWorkObjectMentionConstructor(t *testing.T) {
	for _, tt := range []struct {
		name  string
		style *RichTextSectionTextStyle
		want  string
	}{
		{"minimal", nil, `{"type":"work_object_mention","entity_id":"E123ABC456","app_id":"A123ABC456","text":"Project document","url":"https://example.com/document"}`},
		{"styled", &RichTextSectionTextStyle{Bold: true}, `{"type":"work_object_mention","entity_id":"E123ABC456","app_id":"A123ABC456","text":"Project document","url":"https://example.com/document","style":{"bold":true}}`},
	} {
		t.Run(tt.name, func(t *testing.T) {
			element := NewRichTextSectionWorkObjectMentionElement("E123ABC456", "A123ABC456", "Project document", "https://example.com/document", tt.style)
			assert.Equal(t, RTSEWorkObjectMention, element.RichTextSectionElementType())
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
