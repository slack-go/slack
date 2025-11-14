package slack

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewImageBlockElement(t *testing.T) {
	imageElement := NewImageBlockElement("https://api.slack.com/img/blocks/bkb_template_images/tripAgentLocationMarker.png", "Location Pin Icon")

	assert.Equal(t, string(imageElement.Type), "image")
	assert.Contains(t, *imageElement.ImageURL, "tripAgentLocationMarker")
	assert.Equal(t, imageElement.AltText, "Location Pin Icon")
}

func TestNewImageBlockElementSlackFile(t *testing.T) {
	slackFile := &SlackFileObject{URL: "https://api.slack.com/img/blocks/bkb_template_images/tripAgentLocationMarker.png"}
	imageElement := NewImageBlockElementSlackFile(slackFile, "Location Pin Icon")

	assert.Equal(t, string(imageElement.Type), "image")
	assert.Contains(t, imageElement.SlackFile.URL, "tripAgentLocationMarker")
	assert.Equal(t, imageElement.AltText, "Location Pin Icon")
	assert.Nil(t, imageElement.ImageURL, "ImageURL should be nil when SlackFile is provided")
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

func TestNewFeedbackButton(t *testing.T) {
	btnText := NewTextBlockObject("plain_text", "Good", false, false)
	feedbackButton := NewFeedbackButton(btnText, "positive_feedback")

	assert.Equal(t, feedbackButton.Text.Text, "Good")
	assert.Equal(t, feedbackButton.Value, "positive_feedback")
	assert.Equal(t, feedbackButton.AccessibilityLabel, "")

	feedbackButton.WithAccessibilityLabel("Mark as good")
	assert.Equal(t, feedbackButton.AccessibilityLabel, "Mark as good")
}

func TestNewFeedbackButtonsBlockElement(t *testing.T) {
	positiveBtnText := NewTextBlockObject("plain_text", "👍", false, false)
	negativeBtnText := NewTextBlockObject("plain_text", "👎", false, false)
	positiveBtn := NewFeedbackButton(positiveBtnText, "positive")
	negativeBtn := NewFeedbackButton(negativeBtnText, "negative")

	feedbackElement := NewFeedbackButtonsBlockElement("feedback_1", positiveBtn, negativeBtn)

	assert.Equal(t, string(feedbackElement.Type), "feedback_buttons")
	assert.Equal(t, feedbackElement.ActionID, "feedback_1")
	assert.Equal(t, feedbackElement.PositiveButton.Value, "positive")
	assert.Equal(t, feedbackElement.NegativeButton.Value, "negative")
}

func TestFeedbackButtonsFluentMethods(t *testing.T) {
	positiveBtnText := NewTextBlockObject("plain_text", "Good", false, false)
	negativeBtnText := NewTextBlockObject("plain_text", "Bad", false, false)
	positiveBtn := NewFeedbackButton(positiveBtnText, "pos")
	negativeBtn := NewFeedbackButton(negativeBtnText, "neg")

	feedbackElement := NewFeedbackButtonsBlockElement("feedback_1", positiveBtn, negativeBtn)

	newPositiveText := NewTextBlockObject("plain_text", "Excellent", false, false)
	newPositiveBtn := NewFeedbackButton(newPositiveText, "excellent")
	feedbackElement.WithPositiveButton(newPositiveBtn)
	assert.Equal(t, feedbackElement.PositiveButton.Value, "excellent")

	newNegativeText := NewTextBlockObject("plain_text", "Poor", false, false)
	newNegativeBtn := NewFeedbackButton(newNegativeText, "poor")
	feedbackElement.WithNegativeButton(newNegativeBtn)
	assert.Equal(t, feedbackElement.NegativeButton.Value, "poor")
}

func TestFeedbackButtonsJSONMarshalling(t *testing.T) {
	positiveBtnText := NewTextBlockObject("plain_text", "Good", false, false)
	negativeBtnText := NewTextBlockObject("plain_text", "Bad", false, false)
	positiveBtn := NewFeedbackButton(positiveBtnText, "positive_feedback")
	negativeBtn := NewFeedbackButton(negativeBtnText, "negative_feedback")
	feedbackElement := NewFeedbackButtonsBlockElement("feedback_buttons_1", positiveBtn, negativeBtn)

	data, err := json.Marshal(feedbackElement)
	assert.NoError(t, err)
	assert.NotNil(t, data)

	var unmarshalled FeedbackButtonsBlockElement
	err = json.Unmarshal(data, &unmarshalled)
	assert.NoError(t, err)
	assert.Equal(t, "feedback_buttons", string(unmarshalled.Type))
	assert.Equal(t, "feedback_buttons_1", unmarshalled.ActionID)
	assert.Equal(t, "positive_feedback", unmarshalled.PositiveButton.Value)
	assert.Equal(t, "negative_feedback", unmarshalled.NegativeButton.Value)
}

func TestNewIconButtonBlockElement(t *testing.T) {
	btnText := NewTextBlockObject("plain_text", "Delete", false, false)
	iconButton := NewIconButtonBlockElement("trash", btnText, "delete_action")

	assert.Equal(t, string(iconButton.Type), "icon_button")
	assert.Equal(t, iconButton.Icon, "trash")
	assert.Equal(t, iconButton.Text.Text, "Delete")
	assert.Equal(t, iconButton.ActionID, "delete_action")
}

func TestIconButtonFluentMethods(t *testing.T) {
	btnText := NewTextBlockObject("plain_text", "Delete", false, false)
	iconButton := NewIconButtonBlockElement("trash", btnText, "delete_action")

	iconButton.WithValue("item_123")
	assert.Equal(t, iconButton.Value, "item_123")

	iconButton.WithAccessibilityLabel("Delete this item")
	assert.Equal(t, iconButton.AccessibilityLabel, "Delete this item")

	iconButton.WithVisibleToUserIDs([]string{"U123", "U456"})
	assert.Equal(t, len(iconButton.VisibleToUserIDs), 2)
	assert.Contains(t, iconButton.VisibleToUserIDs, "U123")

	titleText := NewTextBlockObject("plain_text", "Are you sure?", false, false)
	messageText := NewTextBlockObject("plain_text", "This will delete the item", false, false)
	confirmText := NewTextBlockObject("plain_text", "Yes", false, false)
	denyText := NewTextBlockObject("plain_text", "No", false, false)
	confirmObj := NewConfirmationBlockObject(titleText, messageText, confirmText, denyText)
	iconButton.WithConfirm(confirmObj)
	assert.NotNil(t, iconButton.Confirm)
	assert.Equal(t, iconButton.Confirm.Title.Text, "Are you sure?")
}

func TestIconButtonJSONMarshalling(t *testing.T) {
	btnText := NewTextBlockObject("plain_text", "Delete", false, false)
	iconButton := NewIconButtonBlockElement("trash", btnText, "delete_button_1")
	iconButton.WithValue("delete_item")

	data, err := json.Marshal(iconButton)
	assert.NoError(t, err)
	assert.NotNil(t, data)

	var unmarshalled IconButtonBlockElement
	err = json.Unmarshal(data, &unmarshalled)
	assert.NoError(t, err)
	assert.Equal(t, "icon_button", string(unmarshalled.Type))
	assert.Equal(t, "trash", unmarshalled.Icon)
	assert.Equal(t, "delete_button_1", unmarshalled.ActionID)
	assert.Equal(t, "delete_item", unmarshalled.Value)
}

func TestNewWorkflowButtonBlockElement(t *testing.T) {
	btnText := NewTextBlockObject("plain_text", "Run Workflow", false, false)
	workflow := &Workflow{
		Trigger: &WorkflowTrigger{
			URL: "https://slack.com/shortcuts/Ft123456/xyz123",
			CustomizableInputParameters: []CustomizableInputParameter{
				{Name: "input_param_a", Value: "Value for input param A"},
				{Name: "input_param_b", Value: "Value for input param B"},
			},
		},
	}
	workflowButton := NewWorkflowButtonBlockElement(btnText, workflow, "workflow_action_1")

	assert.Equal(t, string(workflowButton.Type), "workflow_button")
	assert.Equal(t, "workflow_action_1", workflowButton.ActionID)
	assert.Equal(t, "Run Workflow", workflowButton.Text.Text)
	assert.NotNil(t, workflowButton.Workflow)
	assert.Equal(t, "https://slack.com/shortcuts/Ft123456/xyz123", workflowButton.Workflow.Trigger.URL)
	assert.Equal(t, 2, len(workflowButton.Workflow.Trigger.CustomizableInputParameters))
	assert.Equal(t, "input_param_a", workflowButton.Workflow.Trigger.CustomizableInputParameters[0].Name)
	assert.Equal(t, "Value for input param A", workflowButton.Workflow.Trigger.CustomizableInputParameters[0].Value)
}

func TestWorkflowButtonFluentMethods(t *testing.T) {
	btnText := NewTextBlockObject("plain_text", "Execute", false, false)
	workflow := &Workflow{
		Trigger: &WorkflowTrigger{
			URL: "https://slack.com/shortcuts/Ft123456/xyz123",
		},
	}
	workflowButton := NewWorkflowButtonBlockElement(btnText, workflow, "workflow_1")

	// Test WithStyle
	workflowButton.WithStyle(StylePrimary)
	assert.Equal(t, StylePrimary, workflowButton.Style)

	workflowButton.WithStyle(StyleDanger)
	assert.Equal(t, StyleDanger, workflowButton.Style)

	// Test WithAccessibilityLabel
	workflowButton.WithAccessibilityLabel("This button triggers an important workflow")
	assert.Equal(t, "This button triggers an important workflow", workflowButton.AccessibilityLabel)

	// Test method chaining
	chainedButton := NewWorkflowButtonBlockElement(btnText, workflow, "workflow_2").
		WithStyle(StylePrimary).
		WithAccessibilityLabel("Chained accessibility label")

	assert.Equal(t, StylePrimary, chainedButton.Style)
	assert.Equal(t, "Chained accessibility label", chainedButton.AccessibilityLabel)
}

func TestWorkflowButtonJSONMarshalling(t *testing.T) {
	btnText := NewTextBlockObject("plain_text", "Start Process", false, false)
	workflow := &Workflow{
		Trigger: &WorkflowTrigger{
			URL: "https://slack.com/shortcuts/Ft123456/abc789",
			CustomizableInputParameters: []CustomizableInputParameter{
				{Name: "user_id", Value: "U123456"},
				{Name: "channel_id", Value: "C789012"},
			},
		},
	}
	workflowButton := NewWorkflowButtonBlockElement(btnText, workflow, "start_workflow").
		WithStyle(StylePrimary).
		WithAccessibilityLabel("Start the approval process")

	jsonData, err := json.Marshal(workflowButton)
	assert.NoError(t, err)

	var unmarshalled WorkflowButtonBlockElement
	err = json.Unmarshal(jsonData, &unmarshalled)
	assert.NoError(t, err)

	assert.Equal(t, "workflow_button", string(unmarshalled.Type))
	assert.Equal(t, "start_workflow", unmarshalled.ActionID)
	assert.Equal(t, "Start Process", unmarshalled.Text.Text)
	assert.Equal(t, StylePrimary, unmarshalled.Style)
	assert.Equal(t, "Start the approval process", unmarshalled.AccessibilityLabel)
	assert.NotNil(t, unmarshalled.Workflow)
	assert.Equal(t, "https://slack.com/shortcuts/Ft123456/abc789", unmarshalled.Workflow.Trigger.URL)
	assert.Equal(t, 2, len(unmarshalled.Workflow.Trigger.CustomizableInputParameters))
	assert.Equal(t, "user_id", unmarshalled.Workflow.Trigger.CustomizableInputParameters[0].Name)
	assert.Equal(t, "U123456", unmarshalled.Workflow.Trigger.CustomizableInputParameters[0].Value)
}

func TestWorkflowButtonMinimalConfiguration(t *testing.T) {
	// Test with minimal required fields only
	btnText := NewTextBlockObject("plain_text", "Simple Workflow", false, false)
	workflow := &Workflow{
		Trigger: &WorkflowTrigger{
			URL: "https://slack.com/shortcuts/Ft123456/minimal",
		},
	}
	workflowButton := NewWorkflowButtonBlockElement(btnText, workflow, "minimal_workflow")

	// Verify no optional fields are set
	assert.Equal(t, Style(""), workflowButton.Style)
	assert.Equal(t, "", workflowButton.AccessibilityLabel)
	assert.Nil(t, workflowButton.Workflow.Trigger.CustomizableInputParameters)

	// Ensure it marshals correctly without optional fields
	jsonData, err := json.Marshal(workflowButton)
	assert.NoError(t, err)

	// Check that optional fields are omitted from JSON
	jsonStr := string(jsonData)
	assert.NotContains(t, jsonStr, "style")
	assert.NotContains(t, jsonStr, "accessibility_label")
	assert.NotContains(t, jsonStr, "customizable_input_parameters")
}
