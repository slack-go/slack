package slack

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRichTextElementFieldsRoundTrip(t *testing.T) {
	tests := []struct {
		name    string
		payload string
	}{
		{
			name: "section styles and metadata",
			payload: `[{"type":"rich_text","elements":[{"type":"rich_text_section","elements":[
				{"type":"broadcast","range":"here","style":{"bold":true,"underline":true,"unlink":true}},
				{"type":"color","value":"#F405B3","style":{"italic":true,"highlight":true}},
				{"type":"date","timestamp":1720710212,"format":"{date_num} at {time}","timezone":"Europe/Madrid","style":{"strike":true,"client_highlight":true}},
				{"type":"channel","channel_id":"C123ABC456","tab_id":"tab123","from_llm":true},
				{"type":"user","user_id":"U123ABC456","from_llm":true},
				{"type":"link","url":"https://docs.slack.dev","unsafe":true,"from_llm":true,"is_slack_url":true,"truncated":true},
				{"type":"emoji","name":"custom_emoji","url":"https://emoji.slack-edge.com/T123/custom_emoji/example.png"}
			]}]}]`,
		},
		{
			name: "quote explicit false metadata",
			payload: `[{"type":"rich_text","elements":[{"type":"rich_text_quote","elements":[
				{"type":"channel","channel_id":"C123ABC456","from_llm":false},
				{"type":"user","user_id":"U123ABC456","from_llm":false},
				{"type":"link","url":"https://docs.slack.dev","unsafe":false,"from_llm":false,"is_slack_url":false,"truncated":false}
			]}]}]`,
		},
		{
			name: "list explicit empty styles",
			payload: `[{"type":"rich_text","elements":[{
				"type":"rich_text_list","style":"bullet","indent":0,"border":0,"offset":0,
				"elements":[{"type":"rich_text_section","elements":[
					{"type":"broadcast","range":"here","style":{}},
					{"type":"color","value":"#F405B3","style":{}},
					{"type":"date","timestamp":1720710212,"format":"{date_num}","style":{}}
				]}]
			}]}]`,
		},
		{
			name: "preformatted link metadata",
			payload: `[{"type":"rich_text","elements":[{"type":"rich_text_preformatted","border":2,"elements":[
				{"type":"link","url":"https://docs.slack.dev","unsafe":false,"from_llm":true,"is_slack_url":false,"truncated":true}
			]}]}]`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var blocks Blocks
			require.NoError(t, json.Unmarshal([]byte(tt.payload), &blocks))
			require.Len(t, blocks.BlockSet, 1)
			require.IsType(t, &RichTextBlock{}, blocks.BlockSet[0])
			block := blocks.BlockSet[0].(*RichTextBlock)
			for _, element := range block.Elements {
				assertKnownRichTextFieldElements(t, element)
			}

			encoded, err := json.Marshal(blocks)
			require.NoError(t, err)
			assert.JSONEq(t, tt.payload, string(encoded))
		})
	}
}

func assertKnownRichTextFieldElements(t *testing.T, element RichTextElement) {
	t.Helper()
	var elements []RichTextSectionElement
	switch element := element.(type) {
	case *RichTextSection:
		elements = element.Elements
	case *RichTextQuote:
		elements = element.Elements
	case *RichTextPreformatted:
		elements = element.Elements
	case *RichTextList:
		for _, child := range element.Elements {
			assertKnownRichTextFieldElements(t, child)
		}
	default:
		t.Fatalf("unexpected rich text container %T", element)
	}
	for _, element := range elements {
		_, unknown := element.(*RichTextSectionUnknownElement)
		assert.False(t, unknown, "element must decode into a concrete type")
	}
}

func TestRichTextElementOptionalFieldsOmitted(t *testing.T) {
	tests := []struct {
		name    string
		element RichTextSectionElement
		want    string
	}{
		{"broadcast", NewRichTextSectionBroadcastElement("here"), `{"type":"broadcast","range":"here"}`},
		{"color", NewRichTextSectionColorElement("#F405B3"), `{"type":"color","value":"#F405B3"}`},
		{"date", NewRichTextSectionDateElement(1720710212, "{date_num}", nil, nil), `{"type":"date","timestamp":1720710212,"format":"{date_num}"}`},
		{"channel", NewRichTextSectionChannelElement("C123ABC456", nil), `{"type":"channel","channel_id":"C123ABC456"}`},
		{"user", NewRichTextSectionUserElement("U123ABC456", nil), `{"type":"user","user_id":"U123ABC456"}`},
		{"link", NewRichTextSectionLinkElement("https://docs.slack.dev", "", nil), `{"type":"link","url":"https://docs.slack.dev"}`},
		{"emoji", NewRichTextSectionEmojiElement("wave", 0, nil), `{"type":"emoji","name":"wave"}`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			encoded, err := json.Marshal(tt.element)
			require.NoError(t, err)
			assert.JSONEq(t, tt.want, string(encoded))
		})
	}
}
