package slack

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRichTextWorkflowMentionRoundTrip(t *testing.T) {
	for _, tt := range []struct {
		name    string
		payload string
	}{
		{
			name: "minimal mention",
			payload: `[{"type":"rich_text","elements":[{"type":"rich_text_section","elements":[
				{"type":"workflow_mention","workflow_id":"Wf123ABC456","function_trigger_id":"Ft123ABC456","text":"Run workflow"}
			]}]}]`,
		},
		{
			name: "live assigned message location",
			payload: `[{"type":"rich_text","block_id":"reference_workflow-minimal","elements":[{"type":"rich_text_section","elements":[
				{"type":"workflow_mention","workflow_id":"Wf0ANH6YFP6D","function_trigger_id":"Ft0AMNRX2S9M","text":"Yahoo workflow (minimal)","channel_id":"C090051UMUP","ts":"1792927030.141719"}
			]}]}]`,
		},
		{
			// The captured trigger URL's suffix is replaced with a fixture value.
			name: "live full mention",
			payload: `[{"type":"rich_text","block_id":"reference_workflow-full","elements":[{"type":"rich_text_section","elements":[
				{"type":"workflow_mention","workflow_id":"Wf0ANH6YFP6D","function_trigger_id":"Ft0AMNRX2S9M","text":"Yahoo workflow (all fields)","url":"https://slack.com/shortcuts/Ft0AMNRX2S9M/test-trigger","channel_id":"C090051UMUP","ts":"1792927243.818989","style":{"bold":true,"italic":true,"underline":true}}
			]}]}]`,
		},
		{
			name: "all optional fields",
			payload: `[{"type":"rich_text","elements":[{"type":"rich_text_section","elements":[
				{"type":"workflow_mention","workflow_id":"Wf123ABC456","function_trigger_id":"Ft123ABC456","text":"Run workflow","url":"https://slack.com/shortcuts/Ft123ABC456/trigger","channel_id":"C123ABC456","ts":"1720710212.123456","style":{"bold":true,"italic":true,"strike":true,"highlight":true,"client_highlight":true,"underline":true,"unlink":true}}
			]}]}]`,
		},
		{
			name: "quote with empty style",
			payload: `[{"type":"rich_text","elements":[{"type":"rich_text_quote","elements":[
				{"type":"workflow_mention","workflow_id":"Wf123ABC456","function_trigger_id":"Ft123ABC456","text":"Run workflow","style":{}}
			]}]}]`,
		},
		{
			name: "list mention",
			payload: `[{"type":"rich_text","elements":[{"type":"rich_text_list","style":"bullet","indent":0,"border":0,"offset":0,"elements":[{"type":"rich_text_section","elements":[
				{"type":"workflow_mention","workflow_id":"Wf123ABC456","function_trigger_id":"Ft123ABC456","text":"Run workflow"}
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

func TestRichTextWorkflowMentionConstructor(t *testing.T) {
	for _, tt := range []struct {
		name  string
		style *RichTextSectionTextStyle
		want  string
	}{
		{"minimal", nil, `{"type":"workflow_mention","workflow_id":"Wf123ABC456","function_trigger_id":"Ft123ABC456","text":"Run workflow"}`},
		{"styled", &RichTextSectionTextStyle{Bold: true}, `{"type":"workflow_mention","workflow_id":"Wf123ABC456","function_trigger_id":"Ft123ABC456","text":"Run workflow","style":{"bold":true}}`},
	} {
		t.Run(tt.name, func(t *testing.T) {
			element := NewRichTextSectionWorkflowMentionElement("Wf123ABC456", "Ft123ABC456", "Run workflow", tt.style)
			assert.Equal(t, RTSEWorkflowMention, element.RichTextSectionElementType())
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
