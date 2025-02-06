package slack

import (
	"encoding/json"
	"testing"

	"github.com/go-test/deep"
	"github.com/stretchr/testify/assert"
)

const (
	dummyPayload = `{
  "type":"rich_text",
  "block_id":"FaYCD",
  "elements": [
    {
      "type":"rich_text_section",
      "elements": [
        {
          "type":"channel",
          "channel_id":"C012345678"
        },
        {
          "type":"text",
          "text":"dummy_text"
        }
      ]
    }
  ]
}`

	richTextQuotePayload = `{
		"type": "rich_text",
		"block_id": "G7G",
		"elements": [
			{
				"type": "rich_text_section",
				"elements": [
					{
						"type": "text",
						"text": "Holy moly\n\n"
					}
				]
			},
			{
				"type": "rich_text_preformatted",
				"elements": [
					{
						"type": "text",
						"text": "Preformatted\n\n"
					}
				],
				"border": 2
			},
			{
				"type": "rich_text_quote",
				"elements": [
					{
						"type": "text",
						"text": "Quote\n\n"
					}
				]
			},
			{
				"type": "rich_text_quote",
				"elements": [
					{
						"type": "text",
						"text": "Another quote"
					}
				]
			},
			{
				"type": "rich_text_preformatted",
				"elements": [
					{
						"type": "text",
						"text": "Another preformatted\n\n"
					}
				],
				"border": 42
			}
		]
	}`
)

func TestNewRichTextBlock(t *testing.T) {
	richTextBlock := NewRichTextBlock("test_block")

	assert.Equal(t, richTextBlock.BlockType(), MBTRichText)
	assert.Equal(t, string(richTextBlock.Type), "rich_text")
	assert.Equal(t, richTextBlock.BlockID, "test_block")
	assert.Equal(t, richTextBlock.ID(), "test_block")
}

func TestRichTextBlock_UnmarshalJSON(t *testing.T) {
	cases := []struct {
		raw      []byte
		expected RichTextBlock
		err      error
	}{
		{
			[]byte(`{"elements":[{"type":"rich_text_unknown"},{"type":"rich_text_section"},{"type":"rich_text_list"}]}`),
			RichTextBlock{
				Elements: []RichTextElement{
					&RichTextUnknown{Type: RTEUnknown, Raw: `{"type":"rich_text_unknown"}`},
					&RichTextSection{Type: RTESection, Elements: []RichTextSectionElement{}},
					&RichTextList{Type: RTEList, Elements: []RichTextElement{}},
				},
			},
			nil,
		},
		{
			[]byte(`{"type": "rich_text","block_id":"blk","elements":[]}`),
			RichTextBlock{
				Type:     MBTRichText,
				BlockID:  "blk",
				Elements: []RichTextElement{},
			},
			nil,
		},
		{
			[]byte(dummyPayload),
			RichTextBlock{
				Type:    MBTRichText,
				BlockID: "FaYCD",
				Elements: []RichTextElement{
					&RichTextSection{
						Type: RTESection,
						Elements: []RichTextSectionElement{
							&RichTextSectionChannelElement{Type: RTSEChannel, ChannelID: "C012345678"},
							&RichTextSectionTextElement{Type: RTSEText, Text: "dummy_text"},
						},
					},
				},
			},
			nil,
		},
		{
			[]byte(richTextQuotePayload),
			RichTextBlock{
				Type:    MBTRichText,
				BlockID: "G7G",
				Elements: []RichTextElement{
					&RichTextSection{Type: RTESection, Elements: []RichTextSectionElement{&RichTextSectionTextElement{Type: RTSEText, Text: "Holy moly\n\n"}}},
					&RichTextPreformatted{RichTextSection: RichTextSection{Type: RTEPreformatted, Elements: []RichTextSectionElement{&RichTextSectionTextElement{Type: RTSEText, Text: "Preformatted\n\n"}}}, Border: 2},
					&RichTextQuote{Type: RTEQuote, Elements: []RichTextSectionElement{&RichTextSectionTextElement{Type: RTSEText, Text: "Quote\n\n"}}},
					&RichTextQuote{Type: RTEQuote, Elements: []RichTextSectionElement{&RichTextSectionTextElement{Type: RTSEText, Text: "Another quote"}}},
					&RichTextPreformatted{RichTextSection: RichTextSection{Type: RTEPreformatted, Elements: []RichTextSectionElement{&RichTextSectionTextElement{Type: RTSEText, Text: "Another preformatted\n\n"}}}, Border: 42},
				},
			},
			nil,
		},
	}
	for _, tc := range cases {
		var actual RichTextBlock
		err := json.Unmarshal(tc.raw, &actual)
		if err != nil {
			if tc.err == nil {
				t.Errorf("unexpected error: %s", err)
			}
			t.Errorf("expected error is %v, but got %v", tc.err, err)
		}
		if tc.err != nil {
			t.Errorf("expected to raise an error %v", tc.err)
		}
		if diff := deep.Equal(actual, tc.expected); diff != nil {
			t.Errorf("actual value does not match expected one\n%s", diff)
		}
	}
}

func TestRichTextSection_UnmarshalJSON(t *testing.T) {
	cases := []struct {
		raw      []byte
		expected RichTextSection
		err      error
	}{
		{
			[]byte(`{"elements":[{"type":"unknown","value":10},{"type":"text","text":"hi"},{"type":"date","timestamp":1636961629,"format":"{date_short_pretty}"},{"type":"date","timestamp":1636961629,"format":"{date_short_pretty}","url":"https://example.com","fallback":"default"}]}`),
			RichTextSection{
				Type: RTESection,
				Elements: []RichTextSectionElement{
					&RichTextSectionUnknownElement{Type: RTSEUnknown, Raw: `{"type":"unknown","value":10}`},
					&RichTextSectionTextElement{Type: RTSEText, Text: "hi"},
					&RichTextSectionDateElement{Type: RTSEDate, Timestamp: JSONTime(1636961629), Format: "{date_short_pretty}"},
					&RichTextSectionDateElement{Type: RTSEDate, Timestamp: JSONTime(1636961629), Format: "{date_short_pretty}", URL: strp("https://example.com"), Fallback: strp("default")},
				},
			},
			nil,
		},
		{
			[]byte(`{"type": "rich_text_section","elements":[]}`),
			RichTextSection{
				Type:     RTESection,
				Elements: []RichTextSectionElement{},
			},
			nil,
		},
		{
			[]byte(`{"type": "rich_text_section","elements":[{"type": "emoji","name": "+1"}]}`),
			RichTextSection{
				Type: RTESection,
				Elements: []RichTextSectionElement{
					&RichTextSectionEmojiElement{Type: RTSEEmoji, Name: "+1"},
				},
			},
			nil,
		},
		{
			[]byte(`{"type": "rich_text_section","elements":[{"type": "emoji","name": "+1","unicode": "1f44d-1f3fb","skin_tone": 2}]}`),
			RichTextSection{
				Type: RTESection,
				Elements: []RichTextSectionElement{
					&RichTextSectionEmojiElement{Type: RTSEEmoji, Name: "+1", Unicode: "1f44d-1f3fb", SkinTone: 2},
				},
			},
			nil,
		},
	}
	for _, tc := range cases {
		var actual RichTextSection
		err := json.Unmarshal(tc.raw, &actual)
		if err != nil {
			if tc.err == nil {
				t.Errorf("unexpected error: %s", err)
			}
			t.Errorf("expected error is %s, but got %s", tc.err, err)
		}
		if tc.err != nil {
			t.Errorf("expected to raise an error %s", tc.err)
		}
		if diff := deep.Equal(actual, tc.expected); diff != nil {
			t.Errorf("actual value does not match expected one\n%s", diff)
		}
	}
}

func TestRichTextList_UnmarshalJSON(t *testing.T) {
	cases := []struct {
		raw      []byte
		expected RichTextList
		err      error
	}{
		{
			[]byte(`{"style":"ordered","elements":[{"type":"rich_text_unknown","value":10},{"type":"rich_text_section","elements":[{"type":"text","text":"hi"}]}]}`),
			RichTextList{
				Type:  RTEList,
				Style: RTEListOrdered,
				Elements: []RichTextElement{
					&RichTextUnknown{Type: RTEUnknown, Raw: `{"type":"rich_text_unknown","value":10}`},
					&RichTextSection{
						Type: RTESection,
						Elements: []RichTextSectionElement{
							&RichTextSectionTextElement{Type: RTSEText, Text: "hi"},
						},
					},
				},
			},
			nil,
		},
		{
			[]byte(`{"style":"ordered","elements":[{"type":"rich_text_list","style":"bullet","elements":[{"type":"rich_text_section","elements":[{"type":"text","text":"hi"}]}]}]}`),
			RichTextList{
				Type:  RTEList,
				Style: RTEListOrdered,
				Elements: []RichTextElement{
					&RichTextList{
						Type:  RTEList,
						Style: RTEListBullet,
						Elements: []RichTextElement{
							&RichTextSection{
								Type: RTESection,
								Elements: []RichTextSectionElement{
									&RichTextSectionTextElement{Type: RTSEText, Text: "hi"},
								},
							},
						},
					},
				},
			},
			nil,
		},
		{
			[]byte(`{"type": "rich_text_list","elements":[]}`),
			RichTextList{
				Type:     RTEList,
				Elements: []RichTextElement{},
			},
			nil,
		},
		{
			[]byte(`{"type": "rich_text_list","elements":[],"indent":2}`),
			RichTextList{
				Type:     RTEList,
				Indent:   2,
				Elements: []RichTextElement{},
			},
			nil,
		},
	}
	for _, tc := range cases {
		var actual RichTextList
		err := json.Unmarshal(tc.raw, &actual)
		if err != nil {
			if tc.err == nil {
				t.Errorf("unexpected error: %s", err)
			}
			t.Errorf("expected error is %s, but got %s", tc.err, err)
		}
		if tc.err != nil {
			t.Errorf("expected to raise an error %s", tc.err)
		}
		if diff := deep.Equal(actual, tc.expected); diff != nil {
			t.Errorf("actual value does not match expected one\n%s", diff)
		}
	}
}

func TestRichTextQuote_Marshal(t *testing.T) {
	t.Run("rich_text_section", func(t *testing.T) {
		const rawRSE = "{\"type\":\"rich_text_section\",\"elements\":[{\"type\":\"text\",\"text\":\"Some Text\"},{\"type\":\"emoji\",\"name\":\"+1\"},{\"type\":\"emoji\",\"name\":\"+1\",\"skin_tone\":2}]}"

		var got RichTextSection
		if err := json.Unmarshal([]byte(rawRSE), &got); err != nil {
			t.Fatal(err)
		}
		want := RichTextSection{
			Type: RTESection,
			Elements: []RichTextSectionElement{
				&RichTextSectionTextElement{Type: RTSEText, Text: "Some Text"},
				&RichTextSectionEmojiElement{Type: RTSEEmoji, Name: "+1"},
				&RichTextSectionEmojiElement{Type: RTSEEmoji, Name: "+1", SkinTone: 2},
			},
		}

		if diff := deep.Equal(got, want); diff != nil {
			t.Errorf("actual value does not match expected one\n%s", diff)
		}
		b, err := json.Marshal(got)
		if err != nil {
			t.Fatal(err)
		}
		if diff := deep.Equal(string(b), rawRSE); diff != nil {
			t.Errorf("actual value does not match expected one\n%s", diff)
		}
	})
	t.Run("rich_text_quote", func(t *testing.T) {
		const rawRTS = "{\"type\":\"rich_text_quote\",\"elements\":[{\"type\":\"text\",\"text\":\"Some text\"}]}"

		var got RichTextQuote
		if err := json.Unmarshal([]byte(rawRTS), &got); err != nil {
			t.Fatal(err)
		}
		want := RichTextQuote{
			Type: RTEQuote,
			Elements: []RichTextSectionElement{
				&RichTextSectionTextElement{Type: RTSEText, Text: "Some text"},
			},
		}
		if diff := deep.Equal(got, want); diff != nil {
			t.Errorf("actual value does not match expected one\n%s", diff)
		}
		b, err := json.Marshal(got)
		if err != nil {
			t.Fatal(err)
		}
		if diff := deep.Equal(string(b), rawRTS); diff != nil {
			t.Errorf("actual value does not match expected one\n%s", diff)
		}
	})
	t.Run("rich_text_preformatted", func(t *testing.T) {
		const rawRTP = "{\"type\":\"rich_text_preformatted\",\"elements\":[{\"type\":\"text\",\"text\":\"Some other text\"}],\"border\":2}"
		want := RichTextPreformatted{
			RichTextSection: RichTextSection{
				Type:     RTEPreformatted,
				Elements: []RichTextSectionElement{&RichTextSectionTextElement{Type: RTSEText, Text: "Some other text"}},
			},
			Border: 2,
		}
		var got RichTextPreformatted
		if err := json.Unmarshal([]byte(rawRTP), &got); err != nil {
			t.Fatal(err)
		}
		if diff := deep.Equal(got, want); diff != nil {
			t.Errorf("actual value does not match expected one\n%s", diff)
		}
		b, err := json.Marshal(got)
		if err != nil {
			t.Fatal(err)
		}
		if diff := deep.Equal(string(b), rawRTP); diff != nil {
			t.Errorf("actual value does not match expected one\n%s", diff)
		}
	})
}

func strp(in string) *string { return &in }
