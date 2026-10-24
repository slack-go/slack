package slack

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRichTextTagRoundTrip(t *testing.T) {
	for _, tt := range []struct {
		name    string
		payload string
	}{
		{
			name: "minimal tag",
			payload: `[{"type":"rich_text","elements":[{"type":"rich_text_section","elements":[
				{"type":"tag","text":"In progress"}
			]}]}]`,
		},
		{
			name: "all optional fields",
			payload: `[{"type":"rich_text","elements":[{"type":"rich_text_section","elements":[
				{"type":"tag","text":"In progress","color":"blue","style":{"bold":true,"italic":true,"strike":true,"highlight":true,"client_highlight":true,"underline":true,"unlink":true}}
			]}]}]`,
		},
		{
			name: "quote with empty style",
			payload: `[{"type":"rich_text","elements":[{"type":"rich_text_quote","elements":[
				{"type":"tag","text":"In progress","color":"green","style":{}}
			]}]}]`,
		},
		{
			name: "list tag",
			payload: `[{"type":"rich_text","elements":[{"type":"rich_text_list","style":"bullet","indent":0,"border":0,"offset":0,"elements":[{"type":"rich_text_section","elements":[
				{"type":"tag","text":"In progress","color":"orange"}
			]}]}]}]`,
		},
		{
			name: "future color",
			payload: `[{"type":"rich_text","elements":[{"type":"rich_text_section","elements":[
				{"type":"tag","text":"Future tag","color":"future_color"}
			]}]}]`,
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

func TestRichTextTagConstructor(t *testing.T) {
	for _, tt := range []struct {
		name  string
		style *RichTextSectionTextStyle
		want  string
	}{
		{"minimal", nil, `{"type":"tag","text":"In progress"}`},
		{"styled", &RichTextSectionTextStyle{Bold: true}, `{"type":"tag","text":"In progress","style":{"bold":true}}`},
	} {
		t.Run(tt.name, func(t *testing.T) {
			element := NewRichTextSectionTagElement("In progress", tt.style)
			assert.Equal(t, RTSETag, element.RichTextSectionElementType())
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

func TestRichTextTagColors(t *testing.T) {
	for _, tt := range []struct {
		color RichTextTagColor
		want  string
	}{
		{RichTextTagColorGray, "gray"},
		{RichTextTagColorBrown, "brown"},
		{RichTextTagColorPurple, "purple"},
		{RichTextTagColorIndigo, "indigo"},
		{RichTextTagColorBlue, "blue"},
		{RichTextTagColorGreen, "green"},
		{RichTextTagColorYellow, "yellow"},
		{RichTextTagColorOrange, "orange"},
		{RichTextTagColorRed, "red"},
	} {
		t.Run(tt.want, func(t *testing.T) {
			element := NewRichTextSectionTagElement("In progress", nil)
			element.Color = tt.color
			encoded, err := json.Marshal(element)
			require.NoError(t, err)
			assert.JSONEq(t, `{"type":"tag","text":"In progress","color":"`+tt.want+`"}`, string(encoded))
		})
	}
}
