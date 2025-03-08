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

func TestWithStyleForButtonElement(t *testing.T) {
	// these values are irrelevant in this test
	btnTxt := NewTextBlockObject("plain_text", "Next 2 Results", false, false)
	btnElement := NewButtonBlockElement("test", "click_me_123", btnTxt)

	btnElement.WithStyle(StyleDefault)
	assert.Equal(t, btnElement.Style, Style(""))
	btnElement.WithStyle(StylePrimary)
	assert.Equal(t, btnElement.Style, Style("primary"))
	btnElement.WithStyle(StyleDanger)
	assert.Equal(t, btnElement.Style, Style("danger"))
}

func TestWithURLForButtonElement(t *testing.T) {
	btnTxt := NewTextBlockObject("plain_text", "Next 2 Results", false, false)
	btnElement := NewButtonBlockElement("test", "click_me_123", btnTxt)

	btnElement.WithURL("https://foo.bar")
	assert.Equal(t, btnElement.URL, "https://foo.bar")
}

func TestNewOptionsSelectBlockElement(t *testing.T) {
	testOptionText := NewTextBlockObject("plain_text", "Option One", false, false)
	testOption := NewOptionBlockObject("test", testOptionText, nil)

	option := NewOptionsSelectBlockElement("static_select", nil, "test", testOption)
	assert.Equal(t, option.Type, "static_select")
	assert.Equal(t, len(option.Options), 1)
	assert.Nil(t, option.OptionGroups)
}

func TestNewOptionsGroupSelectBlockElement(t *testing.T) {
	testOptionText := NewTextBlockObject("plain_text", "Option One", false, false)
	testOption := NewOptionBlockObject("test", testOptionText, nil)
	testLabel := NewTextBlockObject("plain_text", "Test Label", false, false)
	testGroupOption := NewOptionGroupBlockElement(testLabel, testOption)

	optGroup := NewOptionsGroupSelectBlockElement("static_select", nil, "test", testGroupOption)

	assert.Equal(t, optGroup.Type, "static_select")
	assert.Equal(t, optGroup.ActionID, "test")
	assert.Equal(t, len(optGroup.OptionGroups), 1)
}

func TestNewOptionsMultiSelectBlockElement(t *testing.T) {
	testOptionText := NewTextBlockObject("plain_text", "Option One", false, false)
	testDescriptionText := NewTextBlockObject("plain_text", "Description One", false, false)
	testOption := NewOptionBlockObject("test", testOptionText, testDescriptionText)

	option := NewOptionsMultiSelectBlockElement("static_select", nil, "test", testOption)
	assert.Equal(t, option.Type, "static_select")
	assert.Equal(t, len(option.Options), 1)
	assert.Nil(t, option.OptionGroups)
}

func TestNewOptionsGroupMultiSelectBlockElement(t *testing.T) {

	testOptionText := NewTextBlockObject("plain_text", "Option One", false, false)
	testOption := NewOptionBlockObject("test", testOptionText, nil)
	testLabel := NewTextBlockObject("plain_text", "Test Label", false, false)
	testGroupOption := NewOptionGroupBlockElement(testLabel, testOption)

	optGroup := NewOptionsGroupMultiSelectBlockElement("static_select", nil, "test", testGroupOption)

	assert.Equal(t, optGroup.Type, "static_select")
	assert.Equal(t, optGroup.ActionID, "test")
	assert.Equal(t, len(optGroup.OptionGroups), 1)

}
func TestNewOverflowBlockElement(t *testing.T) {

	// Build Text Objects associated with each option
	overflowOptionTextOne := NewTextBlockObject("plain_text", "Option One", false, false)
	overflowOptionTextTwo := NewTextBlockObject("plain_text", "Option Two", false, false)
	overflowOptionTextThree := NewTextBlockObject("plain_text", "Option Three", false, false)

	// Build each option, providing a value for the option
	overflowOptionOne := NewOptionBlockObject("value-0", overflowOptionTextOne, nil)
	overflowOptionTwo := NewOptionBlockObject("value-1", overflowOptionTextTwo, nil)
	overflowOptionThree := NewOptionBlockObject("value-2", overflowOptionTextThree, nil)

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

func TestNewTimePickerBlockElement(t *testing.T) {
	timepickerElement := NewTimePickerBlockElement("test")
	assert.Equal(t, string(timepickerElement.Type), "timepicker")
	assert.Equal(t, timepickerElement.ActionID, "test")
}

func TestNewDateTimePickerBlockElement(t *testing.T) {
	datetimepickerElement := NewDateTimePickerBlockElement("test")
	assert.Equal(t, string(datetimepickerElement.Type), "datetimepicker")
	assert.Equal(t, datetimepickerElement.ActionID, "test")
}

func TestNewPlainTextInputBlockElement(t *testing.T) {

	plainTextInputElement := NewPlainTextInputBlockElement(nil, "test")

	assert.Equal(t, string(plainTextInputElement.Type), "plain_text_input")
	assert.Equal(t, plainTextInputElement.ActionID, "test")

}

func TestNewRichTextInputBlockElement(t *testing.T) {
	richTextInputElement := NewRichTextInputBlockElement(nil, "test")
	assert.Equal(t, string(richTextInputElement.Type), "rich_text_input")
	assert.Equal(t, richTextInputElement.ActionID, "test")
}

func TestNewEmailTextInputBlockElement(t *testing.T) {
	emailTextInputElement := NewEmailTextInputBlockElement(nil, "example@example.com")

	assert.Equal(t, string(emailTextInputElement.Type), "email_text_input")
	assert.Equal(t, emailTextInputElement.ActionID, "example@example.com")
}

func TestNewURLTextInputBlockElement(t *testing.T) {
	urlTextInputElement := NewURLTextInputBlockElement(nil, "www.example.com")

	assert.Equal(t, string(urlTextInputElement.Type), "url_text_input")
	assert.Equal(t, urlTextInputElement.ActionID, "www.example.com")
}

func TestNewCheckboxGroupsBlockElement(t *testing.T) {
	// Build Text Objects associated with each option
	checkBoxOptionTextOne := NewTextBlockObject("plain_text", "Check One", false, false)
	checkBoxOptionTextTwo := NewTextBlockObject("plain_text", "Check Two", false, false)
	checkBoxOptionTextThree := NewTextBlockObject("plain_text", "Check Three", false, false)

	checkBoxDescriptionTextOne := NewTextBlockObject("plain_text", "Description One", false, false)
	checkBoxDescriptionTextTwo := NewTextBlockObject("plain_text", "Description Two", false, false)
	checkBoxDescriptionTextThree := NewTextBlockObject("plain_text", "Description Three", false, false)

	// Build each option, providing a value for the option
	checkBoxOptionOne := NewOptionBlockObject("value-0", checkBoxOptionTextOne, checkBoxDescriptionTextOne)
	checkBoxOptionTwo := NewOptionBlockObject("value-1", checkBoxOptionTextTwo, checkBoxDescriptionTextTwo)
	checkBoxOptionThree := NewOptionBlockObject("value-2", checkBoxOptionTextThree, checkBoxDescriptionTextThree)

	// Build checkbox-group element
	checkBoxGroupElement := NewCheckboxGroupsBlockElement("test", checkBoxOptionOne, checkBoxOptionTwo, checkBoxOptionThree)

	assert.Equal(t, string(checkBoxGroupElement.Type), "checkboxes")
	assert.Equal(t, checkBoxGroupElement.ActionID, "test")
	assert.Equal(t, len(checkBoxGroupElement.Options), 3)
}

func TestNewRadioButtonsBlockElement(t *testing.T) {

	// Build Text Objects associated with each option
	radioButtonsOptionTextOne := NewTextBlockObject("plain_text", "Option One", false, false)
	radioButtonsOptionTextTwo := NewTextBlockObject("plain_text", "Option Two", false, false)
	radioButtonsOptionTextThree := NewTextBlockObject("plain_text", "Option Three", false, false)

	// Build each option, providing a value for the option
	radioButtonsOptionOne := NewOptionBlockObject("value-0", radioButtonsOptionTextOne, nil)
	radioButtonsOptionTwo := NewOptionBlockObject("value-1", radioButtonsOptionTextTwo, nil)
	radioButtonsOptionThree := NewOptionBlockObject("value-2", radioButtonsOptionTextThree, nil)

	// Build radio button element
	radioButtonsElement := NewRadioButtonsBlockElement("test", radioButtonsOptionOne, radioButtonsOptionTwo, radioButtonsOptionThree)

	assert.Equal(t, string(radioButtonsElement.Type), "radio_buttons")
	assert.Equal(t, radioButtonsElement.ActionID, "test")
	assert.Equal(t, len(radioButtonsElement.Options), 3)

}

func TestNewNumberInputBlockElement(t *testing.T) {

	numberInputElement := NewNumberInputBlockElement(nil, "test", true)

	assert.Equal(t, string(numberInputElement.Type), "number_input")
	assert.Equal(t, numberInputElement.ActionID, "test")
	assert.Equal(t, numberInputElement.IsDecimalAllowed, true)

}

func TestNewFileInputBlockElement(t *testing.T) {

	fileInputElement := NewFileInputBlockElement("test")

	assert.Equal(t, string(fileInputElement.Type), "file_input")
	assert.Equal(t, fileInputElement.ActionID, "test")

	fileInputElement.WithFileTypes("jpg", "png")
	assert.Equal(t, len(fileInputElement.FileTypes), 2)
	assert.Contains(t, fileInputElement.FileTypes, "jpg")
	assert.Contains(t, fileInputElement.FileTypes, "png")

	fileInputElement.WithMaxFiles(10)
	assert.Equal(t, fileInputElement.MaxFiles, 10)
}
