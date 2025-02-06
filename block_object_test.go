package slack

import (
	"errors"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewImageBlockObject(t *testing.T) {
	imageObject := NewImageBlockElement("https://api.slack.com/img/blocks/bkb_template_images/beagle.png", "Beagle")

	assert.Equal(t, string(imageObject.Type), "image")
	assert.Equal(t, imageObject.AltText, "Beagle")
	assert.Contains(t, imageObject.ImageURL, "beagle.png")
}

func TestNewTextBlockObject(t *testing.T) {
	textObject := NewTextBlockObject("plain_text", "test", true, false)

	assert.Equal(t, textObject.Type, "plain_text")
	assert.Equal(t, textObject.Text, "test")
	assert.True(t, textObject.Emoji, "Emoji property should be true")
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
	tests := []struct {
		input    TextBlockObject
		expected error
	}{
		{
			input: TextBlockObject{
				Type:     "plain_text",
				Text:     "testText",
				Emoji:    false,
				Verbatim: false,
			},
			expected: nil,
		},
		{
			input: TextBlockObject{
				Type:     "mrkdwn",
				Text:     "testText",
				Emoji:    false,
				Verbatim: false,
			},
			expected: nil,
		},
		{
			input: TextBlockObject{
				Type:     "invalid",
				Text:     "testText",
				Emoji:    false,
				Verbatim: false,
			},
			expected: errors.New("type must be either of plain_text or mrkdwn"),
		},
		{
			input: TextBlockObject{
				Type:     "mrkdwn",
				Text:     "testText",
				Emoji:    true,
				Verbatim: false,
			},
			expected: errors.New("emoji cannot be true in mrkdown"),
		},
		{
			input: TextBlockObject{
				Type:     "mrkdwn",
				Text:     "",
				Emoji:    false,
				Verbatim: false,
			},
			expected: errors.New("text must have a minimum length of 1"),
		},
		{
			input: TextBlockObject{
				Type:     "mrkdwn",
				Text:     strings.Repeat("a", 3001),
				Emoji:    false,
				Verbatim: false,
			},
			expected: errors.New("text cannot be longer than 3000 characters"),
		},
	}

	for _, test := range tests {
		err := test.input.Validate()
		assert.Equal(t, err, test.expected)
	}
}
