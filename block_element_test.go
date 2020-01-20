package slack

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewImageBlockElement(t *testing.T) {

	imageElement := NewImageBlockElement("https://api.slack.com/img/blocks/bkb_template_images/tripAgentLocationMarker.png", "Location Pin Icon")

	assert.Equal(t, string(imageElement.Type), "image")
	assert.Contains(t, imageElement.ImageURL, "tripAgentLocationMarker")
	assert.Equal(t, imageElement.AltText, "Location Pin Icon")

}

func TestNewButtonBlockElement(t *testing.T) {

	btnTxt := NewTextBlockObject("plain_text", "Next 2 Results", false, false)
	btnElement := NewButtonBlockElement("test", "click_me_123", btnTxt)

	assert.Equal(t, string(btnElement.Type), "button")
	assert.Equal(t, btnElement.ActionID, "test")
	assert.Equal(t, btnElement.Value, "click_me_123")
	assert.Equal(t, btnElement.Text.Text, "Next 2 Results")

}

func TestNewOptionsSelectBlockElement(t *testing.T) {

	testOptionText := NewTextBlockObject("plain_text", "Option One", false, false)
	testOption := NewOptionBlockObject("test", testOptionText)

	option := NewOptionsSelectBlockElement("static_select", nil, "test", testOption)
	assert.Equal(t, option.Type, "static_select")
	assert.Equal(t, len(option.Options), 1)
	assert.Nil(t, option.OptionGroups)

}

func TestNewOptionsGroupSelectBlockElement(t *testing.T) {

	testOptionText := NewTextBlockObject("plain_text", "Option One", false, false)
	testOption := NewOptionBlockObject("test", testOptionText)
	testLabel := NewTextBlockObject("plain_text", "Test Label", false, false)
	testGroupOption := NewOptionGroupBlockElement(testLabel, testOption)

	optGroup := NewOptionsGroupSelectBlockElement("static_select", nil, "test", testGroupOption)

	assert.Equal(t, string(optGroup.Type), "static_select")
	assert.Equal(t, optGroup.ActionID, "test")
	assert.Equal(t, len(optGroup.OptionGroups), 1)

}

func TestNewOverflowBlockElement(t *testing.T) {

	// Build Text Objects associated with each option
	overflowOptionTextOne := NewTextBlockObject("plain_text", "Option One", false, false)
	overflowOptionTextTwo := NewTextBlockObject("plain_text", "Option Two", false, false)
	overflowOptionTextThree := NewTextBlockObject("plain_text", "Option Three", false, false)

	// Build each option, providing a value for the option
	overflowOptionOne := NewOptionBlockObject("value-0", overflowOptionTextOne)
	overflowOptionTwo := NewOptionBlockObject("value-1", overflowOptionTextTwo)
	overflowOptionThree := NewOptionBlockObject("value-2", overflowOptionTextThree)

	// Build overflow section
	overflowElement := NewOverflowBlockElement("test", overflowOptionOne, overflowOptionTwo, overflowOptionThree)

	assert.Equal(t, string(overflowElement.Type), "overflow")
	assert.Equal(t, overflowElement.ActionID, "test")
	assert.Equal(t, len(overflowElement.Options), 3)

}

func TestNewDatePickerBlockElement(t *testing.T) {

	datepickerElement := NewDatePickerBlockElement("test")

	assert.Equal(t, string(datepickerElement.Type), "datepicker")
	assert.Equal(t, datepickerElement.ActionID, "test")

}

func TestNewPlainTextInputBlockElement(t *testing.T) {
	inputElement := NewPlainTextInputBlockElement(nil, "test")

	assert.Equal(t, string(inputElement.Type), "plain_text_input")
	assert.Equal(t, inputElement.ActionID, "test")
}
