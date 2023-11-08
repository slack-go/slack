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
				Indent:	  2,
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
