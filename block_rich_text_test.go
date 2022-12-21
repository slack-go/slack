package slack

import (
	"encoding/json"
	"testing"

	"github.com/go-test/deep"
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
)

func TestRichTextBlock_UnmarshalJSON(t *testing.T) {
	cases := []struct {
		raw      []byte
		expected RichTextBlock
		err      error
	}{
		{
			[]byte(`{"elements":[{"type":"rich_text_unknown"},{"type":"rich_text_section"}]}`),
			RichTextBlock{
				Elements: []RichTextElement{
					&RichTextUnknown{Type: RTEUnknown, Raw: `{"type":"rich_text_unknown"}`},
					&RichTextSection{Type: RTESection, Elements: []RichTextSectionElement{}},
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
	}
	for _, tc := range cases {
		var actual RichTextBlock
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

func TestRichTextSection_UnmarshalJSON(t *testing.T) {
	cases := []struct {
		raw      []byte
		expected RichTextSection
		err      error
	}{
		{
			[]byte(`{"elements":[{"type":"unknown","value":10},{"type":"text","text":"hi"},{"type":"date","timestamp":1636961629}]}`),
			RichTextSection{
				Type: RTESection,
				Elements: []RichTextSectionElement{
					&RichTextSectionUnknownElement{Type: RTSEUnknown, Raw: `{"type":"unknown","value":10}`},
					&RichTextSectionTextElement{Type: RTSEText, Text: "hi"},
					&RichTextSectionDateElement{Type: RTSEDate, Timestamp: JSONTime(1636961629)},
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
