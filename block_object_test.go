package slack

import (
	"encoding/json"
	"errors"
	"strings"
	"testing"

	"github.com/go-test/deep"
	"github.com/stretchr/testify/assert"
)

func TestNewImageBlockObject(t *testing.T) {
	imageObject := NewImageBlockElement("https://api.slack.com/img/blocks/bkb_template_images/beagle.png", "Beagle")

	assert.Equal(t, string(imageObject.Type), "image")
	assert.Equal(t, imageObject.AltText, "Beagle")
	assert.Contains(t, *imageObject.ImageURL, "beagle.png")
}

func TestNewTextBlockObject(t *testing.T) {
	textObject := NewTextBlockObject("plain_text", "test", true, false)

	assert.Equal(t, textObject.Type, "plain_text")
	assert.Equal(t, textObject.Text, "test")
	assert.True(t, *textObject.Emoji, "Emoji property should be true")
	assert.False(t, textObject.Verbatim, "Verbatim should be false")
}

func TestNewConfirmationBlockObject(t *testing.T) {
	titleObj := NewTextBlockObject("plain_text", "testTitle", false, false)
	textObj := NewTextBlockObject("plain_text", "testText", false, false)
	confirmObj := NewTextBlockObject("plain_text", "testConfirm", false, false)

	confirmation := NewConfirmationBlockObject(titleObj, textObj, confirmObj, nil)

	assert.Equal(t, confirmation.Title.Text, "testTitle")
	assert.Equal(t, confirmation.Text.Text, "testText")
	assert.Equal(t, confirmation.Confirm.Text, "testConfirm")
	assert.Nil(t, confirmation.Deny, "Deny should be nil")
}

func TestWithStyleForConfirmation(t *testing.T) {
	// these values are irrelevant in this test
	titleObj := NewTextBlockObject("plain_text", "testTitle", false, false)
	textObj := NewTextBlockObject("plain_text", "testText", false, false)
	confirmObj := NewTextBlockObject("plain_text", "testConfirm", false, false)
	confirmation := NewConfirmationBlockObject(titleObj, textObj, confirmObj, nil)

	confirmation.WithStyle(StyleDefault)
	assert.Equal(t, confirmation.Style, Style(""))
	confirmation.WithStyle(StylePrimary)
	assert.Equal(t, confirmation.Style, Style("primary"))
	confirmation.WithStyle(StyleDanger)
	assert.Equal(t, confirmation.Style, Style("danger"))
}

func TestNewOptionBlockObject(t *testing.T) {
	valTextObj := NewTextBlockObject("plain_text", "testText", false, false)
	valDescriptionObj := NewTextBlockObject("plain_text", "testDescription", false, false)
	optObj := NewOptionBlockObject("testOpt", valTextObj, valDescriptionObj)

	assert.Equal(t, optObj.Text.Text, "testText")
	assert.Equal(t, optObj.Description.Text, "testDescription")
	assert.Equal(t, optObj.Value, "testOpt")
}

func TestNewOptionGroupBlockElement(t *testing.T) {
	labelObj := NewTextBlockObject("plain_text", "testLabel", false, false)
	valTextObj := NewTextBlockObject("plain_text", "testText", false, false)
	optObj := NewOptionBlockObject("testOpt", valTextObj, nil)

	optGroup := NewOptionGroupBlockElement(labelObj, optObj)

	assert.Equal(t, optGroup.Label.Text, "testLabel")
	assert.Len(t, optGroup.Options, 1, "Options should contain one element")
}

func TestValidateTextBlockObject(t *testing.T) {
	emojiTrue := new(bool)
	emojiFalse := new(bool)
	*emojiTrue = true
	*emojiFalse = false

	tests := []struct {
		input    TextBlockObject
		expected error
	}{
		{
			input: TextBlockObject{
				Type:     "plain_text",
				Text:     "testText",
				Emoji:    emojiFalse,
				Verbatim: false,
			},
			expected: nil,
		},
		{
			input: TextBlockObject{
				Type:     "plain_text",
				Text:     "testText",
				Emoji:    emojiTrue,
				Verbatim: false,
			},
			expected: nil,
		},
		{
			input: TextBlockObject{
				Type:     "plain_text",
				Text:     "testText",
				Emoji:    nil,
				Verbatim: false,
			},
			expected: nil,
		},
		{
			input: TextBlockObject{
				Type:     "mrkdwn",
				Text:     "testText",
				Emoji:    nil,
				Verbatim: false,
			},
			expected: nil,
		},
		{
			input: TextBlockObject{
				Type:     "invalid",
				Text:     "testText",
				Emoji:    emojiFalse,
				Verbatim: false,
			},
			expected: errors.New("type must be either of plain_text or mrkdwn"),
		},
		{
			input: TextBlockObject{
				Type:     "mrkdwn",
				Text:     "testText",
				Emoji:    emojiTrue,
				Verbatim: false,
			},
			expected: errors.New("emoji cannot be set for mrkdwn type"),
		},
		{
			input: TextBlockObject{
				Type:     "mrkdwn",
				Text:     "testText",
				Emoji:    emojiFalse,
				Verbatim: false,
			},
			expected: errors.New("emoji cannot be set for mrkdwn type"),
		},
		{
			input: TextBlockObject{
				Type:     "mrkdwn",
				Text:     "",
				Emoji:    nil,
				Verbatim: false,
			},
			expected: errors.New("text must have a minimum length of 1"),
		},
		{
			input: TextBlockObject{
				Type:     "mrkdwn",
				Text:     strings.Repeat("a", 3001),
				Emoji:    nil,
				Verbatim: false,
			},
			expected: errors.New("text cannot be longer than 3000 characters"),
		},
	}

	for _, test := range tests {
		err := test.input.Validate()
		assert.Equal(t, test.expected, err)
	}
}

func TestTextBlockObject_UnmarshalJSON(t *testing.T) {
	emojiTrue := new(bool)
	emojiFalse := new(bool)
	*emojiTrue = true
	*emojiFalse = false

	cases := []struct {
		raw      []byte
		expected TextBlockObject
		err      error
	}{
		{
			[]byte(`{"type":"plain_text","text":"testText"}`),
			TextBlockObject{
				Type:     "plain_text",
				Text:     "testText",
				Emoji:    nil,
				Verbatim: false,
			},
			nil,
		},
		{
			[]byte(`{"type":"plain_text","text":":+1:","emoji":true}`),
			TextBlockObject{
				Type:     "plain_text",
				Text:     ":+1:",
				Emoji:    emojiTrue,
				Verbatim: false,
			},
			nil,
		},
		{
			[]byte(`{"type":"plain_text","text":"No emojis allowed :(","emoji":false}`),
			TextBlockObject{
				Type:     "plain_text",
				Text:     "No emojis allowed :(",
				Emoji:    emojiFalse,
				Verbatim: false,
			},
			nil,
		},
		{
			[]byte(`{"type":"mrkdwn","text":"testText"}`),
			TextBlockObject{
				Type:     "mrkdwn",
				Text:     "testText",
				Emoji:    nil,
				Verbatim: false,
			},
			nil,
		},
		{
			[]byte(`{"type":"mrkdwn","text":"No emojis allowed :(","emoji":false}`),
			TextBlockObject{
				Type:     "mrkdwn",
				Text:     "No emojis allowed :(",
				Emoji:    emojiFalse,
				Verbatim: false,
			},
			nil,
		},
	}
	for _, tc := range cases {
		var actual TextBlockObject
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
