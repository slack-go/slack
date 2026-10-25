package slack

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRichTextSalesforceDataFieldRoundTrip(t *testing.T) {
	for _, tt := range []struct {
		name    string
		payload string
	}{
		{
			name: "minimal field",
			payload: `[{"type":"rich_text","elements":[{"type":"rich_text_section","elements":[
				{"type":"salesforce_data_field","salesforce_record_id":"001ABC456DEF789"}
			]}]}]`,
		},
		{
			name: "live unavailable field",
			payload: `[{"type":"rich_text","block_id":"reference_salesforce-minimal","elements":[{"type":"rich_text_section","elements":[
				{"type":"text","text":"Salesforce schema probe (no record validation): "},
				{"type":"salesforce_data_field","salesforce_record_id":"001ABC456DEF789","salesforce_field_label":"Unavailable","has_withheld_salesforce_field_names":true}
			]}]}]`,
		},
		{
			name: "live unavailable field with false label visibility",
			payload: `[{"type":"rich_text","block_id":"reference_salesforce-without-label","elements":[{"type":"rich_text_section","elements":[
				{"type":"text","text":"Salesforce schema probe (no record validation): "},
				{"type":"salesforce_data_field","salesforce_record_id":"001ABC456DEF789","salesforce_field_label":"Unavailable","salesforce_include_field_label":false,"style":{"bold":true,"italic":true,"underline":true},"has_withheld_salesforce_field_names":true}
			]}]}]`,
		},
		{
			name: "explicit false withheld flag",
			payload: `[{"type":"rich_text","elements":[{"type":"rich_text_section","elements":[
				{"type":"salesforce_data_field","salesforce_record_id":"001ABC456DEF789","has_withheld_salesforce_field_names":false}
			]}]}]`,
		},
		{
			name: "all optional fields",
			payload: `[{"type":"rich_text","elements":[{"type":"rich_text_section","elements":[
				{"type":"salesforce_data_field","salesforce_record_id":"001ABC456DEF789","salesforce_field_label":"Account Name","salesforce_field_api_name":"Name","salesforce_include_field_label":true,"style":{"bold":true,"italic":true,"strike":true,"highlight":true,"client_highlight":true,"underline":true,"unlink":true}}
			]}]}]`,
		},
		{
			name: "quote explicit false and empty style",
			payload: `[{"type":"rich_text","elements":[{"type":"rich_text_quote","elements":[
				{"type":"salesforce_data_field","salesforce_record_id":"001ABC456DEF789","salesforce_include_field_label":false,"style":{}}
			]}]}]`,
		},
		{
			name: "list field",
			payload: `[{"type":"rich_text","elements":[{"type":"rich_text_list","style":"bullet","indent":0,"border":0,"offset":0,"elements":[{"type":"rich_text_section","elements":[
				{"type":"salesforce_data_field","salesforce_record_id":"001ABC456DEF789","salesforce_field_api_name":"Name"}
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

func TestRichTextSalesforceDataFieldConstructor(t *testing.T) {
	for _, tt := range []struct {
		name  string
		style *RichTextSectionTextStyle
		want  string
	}{
		{"minimal", nil, `{"type":"salesforce_data_field","salesforce_record_id":"001ABC456DEF789"}`},
		{"styled", &RichTextSectionTextStyle{Bold: true}, `{"type":"salesforce_data_field","salesforce_record_id":"001ABC456DEF789","style":{"bold":true}}`},
	} {
		t.Run(tt.name, func(t *testing.T) {
			element := NewRichTextSectionSalesforceDataFieldElement("001ABC456DEF789", tt.style)
			assert.Equal(t, RTSESalesforceDataField, element.RichTextSectionElementType())
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
