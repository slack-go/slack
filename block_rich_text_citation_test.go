package slack

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRichTextCitationRoundTrip(t *testing.T) {
	details := []struct {
		name string
		json string
	}{
		{"file", `{"citation_type":"file","file_id":"F123","descriptor":"Report"}`},
		{"file minimal", `{"citation_type":"file","file_id":"F123"}`},
		{"file empty descriptor", `{"citation_type":"file","file_id":"F123","descriptor":""}`},
		{"external", `{"citation_type":"external","app_name":"Example","app_icon_url":"https://example.com/icon.png"}`},
		{"web", `{"citation_type":"web","display_name":"Example","title":"Title","snippet":"Snippet"}`},
		{"web minimal", `{"citation_type":"web","display_name":"Example"}`},
		{"web empty strings", `{"citation_type":"web","display_name":"Example","title":"","snippet":""}`},
		{"message", `{"citation_type":"message","channel":"C123","message_ts":"1720710212.123456"}`},
		{"memory", `{"citation_type":"memory","memory_id":"memory123"}`},
		{"unknown", `{"citation_type":"future_source","record_id":"R123","extra":{"counter":9007199254740993,"enabled":false}}`},
	}
	containers := []struct{ name, format string }{
		{"section", `{"type":"rich_text_section","elements":[%s]}`},
		{"quote", `{"type":"rich_text_quote","elements":[%s]}`},
		{"list", `{"type":"rich_text_list","style":"bullet","indent":0,"border":0,"offset":0,"elements":[{"type":"rich_text_section","elements":[%s]}]}`},
	}
	for _, detail := range details {
		for _, container := range containers {
			t.Run(detail.name+"/"+container.name, func(t *testing.T) {
				citation := `{"type":"citation","url":"https://example.com/source","text":"Source","index":1.5,"details":` + detail.json + `,"from_llm":false,"is_slack_url":false}`
				payload := `[{"type":"rich_text","elements":[` + fmt.Sprintf(container.format, citation) + `]}]`
				var blocks Blocks
				require.NoError(t, json.Unmarshal([]byte(payload), &blocks))
				require.Len(t, blocks.BlockSet, 1)
				block, ok := blocks.BlockSet[0].(*RichTextBlock)
				require.True(t, ok)
				for _, element := range block.Elements {
					assertKnownRichTextFieldElements(t, element)
				}
				encoded, err := json.Marshal(blocks)
				require.NoError(t, err)
				assert.JSONEq(t, payload, string(encoded))
			})
		}
	}
}

func TestRichTextCitationConstructors(t *testing.T) {
	tests := []struct {
		name    string
		details RichTextCitationDetails
		want    string
	}{
		{"file", NewRichTextCitationFileDetails("F123"), `{"citation_type":"file","file_id":"F123"}`},
		{"external", NewRichTextCitationExternalDetails("Example", "https://example.com/icon.png"), `{"citation_type":"external","app_name":"Example","app_icon_url":"https://example.com/icon.png"}`},
		{"web", NewRichTextCitationWebDetails("Example"), `{"citation_type":"web","display_name":"Example"}`},
		{"message", NewRichTextCitationMessageDetails("C123", "1720710212.123456"), `{"citation_type":"message","channel":"C123","message_ts":"1720710212.123456"}`},
		{"memory", NewRichTextCitationMemoryDetails("memory123"), `{"citation_type":"memory","memory_id":"memory123"}`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			citation := NewRichTextSectionCitationElement("https://example.com", "Source", 0, tt.details)
			assert.Equal(t, RTSECitation, citation.RichTextSectionElementType())
			encoded, err := json.Marshal(citation)
			require.NoError(t, err)
			assert.JSONEq(t, `{"type":"citation","url":"https://example.com","text":"Source","index":0,"details":`+tt.want+`}`, string(encoded))
			var decoded RichTextSectionCitationElement
			require.NoError(t, json.Unmarshal(encoded, &decoded))
			require.IsType(t, tt.details, decoded.Details)
			assert.Equal(t, citation, &decoded)
			assert.Equal(t, tt.details.CitationType(), decoded.Details.CitationType())
		})
	}
}

func TestRichTextCitationUnknownDetails(t *testing.T) {
	tests := []struct {
		name, raw string
		kind      RichTextCitationType
	}{
		{"future source", `{"citation_type":"future_source","counter":9007199254740993,"extra":[true,null,{"text":"kept"}]}`, "future_source"},
		{"missing discriminator", `{"extra":{"key":"value"}}`, ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			payload := `{"type":"citation","url":"https://example.com","text":"Source","index":1,"details":` + tt.raw + `}`
			var citation RichTextSectionCitationElement
			require.NoError(t, json.Unmarshal([]byte(payload), &citation))
			unknown, ok := citation.Details.(*RichTextCitationUnknownDetails)
			require.True(t, ok)
			assert.Equal(t, tt.kind, unknown.CitationType())
			assert.Equal(t, json.RawMessage(tt.raw), unknown.Raw)
			encoded, err := json.Marshal(citation)
			require.NoError(t, err)
			var envelope map[string]json.RawMessage
			require.NoError(t, json.Unmarshal(encoded, &envelope))
			assert.Equal(t, json.RawMessage(tt.raw), envelope["details"])
		})
	}
}

func TestRichTextCitationUnmarshalResetsFields(t *testing.T) {
	var citation RichTextSectionCitationElement
	require.NoError(t, json.Unmarshal([]byte(`{"type":"citation","url":"https://example.com","text":"Old","index":1,"details":{"citation_type":"web","display_name":"Example"},"from_llm":true,"is_slack_url":true}`), &citation))
	require.IsType(t, &RichTextCitationWebDetails{}, citation.Details)
	require.NotNil(t, citation.FromLLM)
	assert.True(t, *citation.FromLLM)
	require.NotNil(t, citation.IsSlackURL)
	assert.True(t, *citation.IsSlackURL)
	payload := `{"type":"citation","url":"https://example.org","text":"New","index":0,"details":null}`
	require.NoError(t, json.Unmarshal([]byte(payload), &citation))
	assert.Nil(t, citation.Details)
	assert.Nil(t, citation.FromLLM)
	assert.Nil(t, citation.IsSlackURL)
	encoded, err := json.Marshal(citation)
	require.NoError(t, err)
	assert.JSONEq(t, payload, string(encoded))
}

func TestRichTextCitationMalformedDetails(t *testing.T) {
	tests := []struct {
		name, details string
	}{
		{"non-string discriminator", `{"citation_type":123,"extra":false}`},
		{"array", `[{"citation_type":"web","display_name":"Example"}]`},
		{"string", `"opaque details"`},
		{"known source with invalid field", `{"citation_type":"message","channel":"C123","message_ts":123}`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			payload := `[{"type":"rich_text","elements":[{"type":"rich_text_section","elements":[{"type":"citation","url":"https://example.com","text":"Source","index":1,"details":` + tt.details + `}]}]}]`
			var blocks Blocks
			err := json.Unmarshal([]byte(payload), &blocks)
			var typeError *json.UnmarshalTypeError
			require.ErrorAs(t, err, &typeError)
		})
	}
}
