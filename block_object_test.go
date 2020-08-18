package slack

import (
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

func TestNewOptionBlockObject(t *testing.T) {

	valTextObj := NewTextBlockObject("plain_text", "testText", false, false)
	optObj := NewOptionBlockObject("testOpt", valTextObj)

	assert.Equal(t, optObj.Text.Text, "testText")
	assert.Equal(t, optObj.Value, "testOpt")

}

func TestNewOptionGroupBlockElement(t *testing.T) {

	labelObj := NewTextBlockObject("plain_text", "testLabel", false, false)
	valTextObj := NewTextBlockObject("plain_text", "testText", false, false)
	optObj := NewOptionBlockObject("testOpt", valTextObj)

	optGroup := NewOptionGroupBlockElement(labelObj, optObj)

	assert.Equal(t, optGroup.Label.Text, "testLabel")
	assert.Len(t, optGroup.Options, 1, "Options should contain one element")

}
