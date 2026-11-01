package slack

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewImageBlockElement(t *testing.T) {
	imageElement := NewImageBlockElement("https://api.slack.com/img/blocks/bkb_template_images/tripAgentLocationMarker.png", "Location Pin Icon")

	assert.Equal(t, "image", string(imageElement.Type))
	assert.Contains(t, *imageElement.ImageURL, "tripAgentLocationMarker")
	assert.Equal(t, "Location Pin Icon", imageElement.AltText)
}

func TestNewImageBlockElementSlackFile(t *testing.T) {
	slackFile := &SlackFileObject{URL: "https://api.slack.com/img/blocks/bkb_template_images/tripAgentLocationMarker.png"}
	imageElement := NewImageBlockElementSlackFile(slackFile, "Location Pin Icon")

	assert.Equal(t, "image", string(imageElement.Type))
	assert.Contains(t, imageElement.SlackFile.URL, "tripAgentLocationMarker")
	assert.Equal(t, "Location Pin Icon", imageElement.AltText)
	assert.Nil(t, imageElement.ImageURL, "ImageURL should be nil when SlackFile is provided")
}

func TestNewButtonBlockElement(t *testing.T) {
	btnTxt := NewTextBlockObject("plain_text", "Next 2 Results", false, false)
	btnElement := NewButtonBlockElement("test", "click_me_123", btnTxt)

	assert.Equal(t, "button", string(btnElement.Type))
	assert.Equal(t, "test", btnElement.ActionID)
	assert.Equal(t, "click_me_123", btnElement.Value)
	assert.Equal(t, "Next 2 Results", btnElement.Text.Text)
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
	assert.Equal(t, "https://foo.bar", btnElement.URL)
}

func TestNewOptionsSelectBlockElement(t *testing.T) {
	testOptionText := NewTextBlockObject("plain_text", "Option One", false, false)
	testOption := NewOptionBlockObject("test", testOptionText, nil)

	option := NewOptionsSelectBlockElement("static_select", nil, "test", testOption)
	assert.Equal(t, "static_select", option.Type)
	assert.Equal(t, 1, len(option.Options))
	assert.Nil(t, option.OptionGroups)
}

func TestNewOptionsGroupSelectBlockElement(t *testing.T) {
	testOptionText := NewTextBlockObject("plain_text", "Option One", false, false)
	testOption := NewOptionBlockObject("test", testOptionText, nil)
	testLabel := NewTextBlockObject("plain_text", "Test Label", false, false)
	testGroupOption := NewOptionGroupBlockElement(testLabel, testOption)

	optGroup := NewOptionsGroupSelectBlockElement("static_select", nil, "test", testGroupOption)

	assert.Equal(t, "static_select", optGroup.Type)
	assert.Equal(t, "test", optGroup.ActionID)
	assert.Equal(t, 1, len(optGroup.OptionGroups))
}

func TestNewOptionsMultiSelectBlockElement(t *testing.T) {
	testOptionText := NewTextBlockObject("plain_text", "Option One", false, false)
	testDescriptionText := NewTextBlockObject("plain_text", "Description One", false, false)
	testOption := NewOptionBlockObject("test", testOptionText, testDescriptionText)

	option := NewOptionsMultiSelectBlockElement("static_select", nil, "test", testOption)
	assert.Equal(t, "static_select", option.Type)
	assert.Equal(t, 1, len(option.Options))
	assert.Nil(t, option.OptionGroups)
}

func TestNewOptionsGroupMultiSelectBlockElement(t *testing.T) {

	testOptionText := NewTextBlockObject("plain_text", "Option One", false, false)
	testOption := NewOptionBlockObject("test", testOptionText, nil)
	testLabel := NewTextBlockObject("plain_text", "Test Label", false, false)
	testGroupOption := NewOptionGroupBlockElement(testLabel, testOption)

	optGroup := NewOptionsGroupMultiSelectBlockElement("static_select", nil, "test", testGroupOption)

	assert.Equal(t, "static_select", optGroup.Type)
	assert.Equal(t, "test", optGroup.ActionID)
	assert.Equal(t, 1, len(optGroup.OptionGroups))

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

	assert.Equal(t, "overflow", string(overflowElement.Type))
	assert.Equal(t, "test", overflowElement.ActionID)
	assert.Equal(t, 3, len(overflowElement.Options))

}

func TestNewDatePickerBlockElement(t *testing.T) {

	datepickerElement := NewDatePickerBlockElement("test")

	assert.Equal(t, "datepicker", string(datepickerElement.Type))
	assert.Equal(t, "test", datepickerElement.ActionID)

}

func TestNewTimePickerBlockElement(t *testing.T) {
	timepickerElement := NewTimePickerBlockElement("test")
	assert.Equal(t, "timepicker", string(timepickerElement.Type))
	assert.Equal(t, "test", timepickerElement.ActionID)
}

func TestNewDateTimePickerBlockElement(t *testing.T) {
	datetimepickerElement := NewDateTimePickerBlockElement("test")
	assert.Equal(t, "datetimepicker", string(datetimepickerElement.Type))
	assert.Equal(t, "test", datetimepickerElement.ActionID)
}

func TestNewPlainTextInputBlockElement(t *testing.T) {

	plainTextInputElement := NewPlainTextInputBlockElement(nil, "test")

	assert.Equal(t, "plain_text_input", string(plainTextInputElement.Type))
	assert.Equal(t, "test", plainTextInputElement.ActionID)

}

func TestNewRichTextInputBlockElement(t *testing.T) {
	richTextInputElement := NewRichTextInputBlockElement(nil, "test")
	assert.Equal(t, "rich_text_input", string(richTextInputElement.Type))
	assert.Equal(t, "test", richTextInputElement.ActionID)
}

func TestNewEmailTextInputBlockElement(t *testing.T) {
	emailTextInputElement := NewEmailTextInputBlockElement(nil, "example@example.com")

	assert.Equal(t, "email_text_input", string(emailTextInputElement.Type))
	assert.Equal(t, "example@example.com", emailTextInputElement.ActionID)
}

func TestNewURLTextInputBlockElement(t *testing.T) {
	urlTextInputElement := NewURLTextInputBlockElement(nil, "www.example.com")

	assert.Equal(t, "url_text_input", string(urlTextInputElement.Type))
	assert.Equal(t, "www.example.com", urlTextInputElement.ActionID)
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

	assert.Equal(t, "checkboxes", string(checkBoxGroupElement.Type))
	assert.Equal(t, "test", checkBoxGroupElement.ActionID)
	assert.Equal(t, 3, len(checkBoxGroupElement.Options))
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

	assert.Equal(t, "radio_buttons", string(radioButtonsElement.Type))
	assert.Equal(t, "test", radioButtonsElement.ActionID)
	assert.Equal(t, 3, len(radioButtonsElement.Options))

}

func TestNewNumberInputBlockElement(t *testing.T) {

	numberInputElement := NewNumberInputBlockElement(nil, "test", true)

	assert.Equal(t, "number_input", string(numberInputElement.Type))
	assert.Equal(t, "test", numberInputElement.ActionID)
	assert.Equal(t, true, numberInputElement.IsDecimalAllowed)

}

func TestNewFileInputBlockElement(t *testing.T) {

	fileInputElement := NewFileInputBlockElement("test")

	assert.Equal(t, "file_input", string(fileInputElement.Type))
	assert.Equal(t, "test", fileInputElement.ActionID)

	fileInputElement.WithFileTypes("jpg", "png")
	assert.Equal(t, 2, len(fileInputElement.FileTypes))
	assert.Contains(t, fileInputElement.FileTypes, "jpg")
	assert.Contains(t, fileInputElement.FileTypes, "png")

	fileInputElement.WithMaxFiles(10)
	assert.Equal(t, 10, fileInputElement.MaxFiles)
}

func TestNewFeedbackButton(t *testing.T) {
	btnText := NewTextBlockObject("plain_text", "Good", false, false)
	feedbackButton := NewFeedbackButton(btnText, "positive_feedback")

	assert.Equal(t, "Good", feedbackButton.Text.Text)
	assert.Equal(t, "positive_feedback", feedbackButton.Value)
	assert.Equal(t, "", feedbackButton.AccessibilityLabel)

	feedbackButton.WithAccessibilityLabel("Mark as good")
	assert.Equal(t, "Mark as good", feedbackButton.AccessibilityLabel)
}

func TestNewFeedbackButtonsBlockElement(t *testing.T) {
	positiveBtnText := NewTextBlockObject("plain_text", "👍", false, false)
	negativeBtnText := NewTextBlockObject("plain_text", "👎", false, false)
	positiveBtn := NewFeedbackButton(positiveBtnText, "positive")
	negativeBtn := NewFeedbackButton(negativeBtnText, "negative")

	feedbackElement := NewFeedbackButtonsBlockElement("feedback_1", positiveBtn, negativeBtn)

	assert.Equal(t, "feedback_buttons", string(feedbackElement.Type))
	assert.Equal(t, "feedback_1", feedbackElement.ActionID)
	assert.Equal(t, "positive", feedbackElement.PositiveButton.Value)
	assert.Equal(t, "negative", feedbackElement.NegativeButton.Value)
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
	assert.Equal(t, "excellent", feedbackElement.PositiveButton.Value)

	newNegativeText := NewTextBlockObject("plain_text", "Poor", false, false)
	newNegativeBtn := NewFeedbackButton(newNegativeText, "poor")
	feedbackElement.WithNegativeButton(newNegativeBtn)
	assert.Equal(t, "poor", feedbackElement.NegativeButton.Value)
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

	assert.Equal(t, "icon_button", string(iconButton.Type))
	assert.Equal(t, "trash", iconButton.Icon)
	assert.Equal(t, "Delete", iconButton.Text.Text)
	assert.Equal(t, "delete_action", iconButton.ActionID)
}

func TestIconButtonFluentMethods(t *testing.T) {
	btnText := NewTextBlockObject("plain_text", "Delete", false, false)
	iconButton := NewIconButtonBlockElement("trash", btnText, "delete_action")

	iconButton.WithValue("item_123")
	assert.Equal(t, "item_123", iconButton.Value)

	iconButton.WithAccessibilityLabel("Delete this item")
	assert.Equal(t, "Delete this item", iconButton.AccessibilityLabel)

	iconButton.WithVisibleToUserIDs([]string{"U123", "U456"})
	assert.Equal(t, 2, len(iconButton.VisibleToUserIDs))
	assert.Contains(t, iconButton.VisibleToUserIDs, "U123")

	titleText := NewTextBlockObject("plain_text", "Are you sure?", false, false)
	messageText := NewTextBlockObject("plain_text", "This will delete the item", false, false)
	confirmText := NewTextBlockObject("plain_text", "Yes", false, false)
	denyText := NewTextBlockObject("plain_text", "No", false, false)
	confirmObj := NewConfirmationBlockObject(titleText, messageText, confirmText, denyText)
	iconButton.WithConfirm(confirmObj)
	require.NotNil(t, iconButton.Confirm)
	assert.Equal(t, "Are you sure?", iconButton.Confirm.Title.Text)
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

	assert.Equal(t, "workflow_button", string(workflowButton.Type))
	assert.Equal(t, "workflow_action_1", workflowButton.ActionID)
	assert.Equal(t, "Run Workflow", workflowButton.Text.Text)
	require.NotNil(t, workflowButton.Workflow)
	assert.Equal(t, "https://slack.com/shortcuts/Ft123456/xyz123", workflowButton.Workflow.Trigger.URL)
	require.Len(t, workflowButton.Workflow.Trigger.CustomizableInputParameters, 2)
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
	require.NotNil(t, unmarshalled.Workflow)
	assert.Equal(t, "https://slack.com/shortcuts/Ft123456/abc789", unmarshalled.Workflow.Trigger.URL)
	require.Len(t, unmarshalled.Workflow.Trigger.CustomizableInputParameters, 2)
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
	workflowButton := NewWorkflowButtonBlockElement(btnText, workflow, "")

	// Verify no optional fields are set
	assert.Equal(t, Style(""), workflowButton.Style)
	assert.Equal(t, "", workflowButton.AccessibilityLabel)
	assert.Nil(t, workflowButton.Workflow.Trigger.CustomizableInputParameters)

	// Ensure it marshals correctly without optional fields
	jsonData, err := json.Marshal(workflowButton)
	assert.NoError(t, err)

	// Check that optional fields are omitted from JSON
	jsonStr := string(jsonData)
	assert.NotContains(t, jsonStr, "action_id")
	assert.NotContains(t, jsonStr, "style")
	assert.NotContains(t, jsonStr, "accessibility_label")
	assert.NotContains(t, jsonStr, "customizable_input_parameters")
}
