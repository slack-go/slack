package slack

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWorkflowButtonBlockElementUnmarshal(t *testing.T) {
	workflowButtonJSON := `{
		"type": "workflow_button",
		"text": {"type": "plain_text", "text": "Run Workflow"},
		"workflow": {"trigger": {"url": "https://slack.com/shortcuts/Ft0123ABC/xyz"}},
		"action_id": "start_workflow"
	}`

	t.Run("BlockElements", func(t *testing.T) {
		elementsJSON := fmt.Sprintf("[%s]", workflowButtonJSON)
		var elements BlockElements
		err := json.Unmarshal([]byte(elementsJSON), &elements)
		require.NoError(t, err)
		require.Len(t, elements.ElementSet, 1)
		assert.IsType(t, &WorkflowButtonBlockElement{}, elements.ElementSet[0])

		wb := elements.ElementSet[0].(*WorkflowButtonBlockElement)
		assert.Equal(t, METWorkflowButton, wb.Type)
		assert.Equal(t, "Run Workflow", wb.Text.Text)
		assert.Equal(t, "start_workflow", wb.ActionID)
		assert.Equal(t, "https://slack.com/shortcuts/Ft0123ABC/xyz", wb.Workflow.Trigger.URL)
	})

	t.Run("InputBlock", func(t *testing.T) {
		inputJSON := fmt.Sprintf(`{
			"type": "input",
			"label": {"type": "plain_text", "text": "Workflow"},
			"element": %s
		}`, workflowButtonJSON)
		var input InputBlock
		err := json.Unmarshal([]byte(inputJSON), &input)
		require.NoError(t, err)
		require.NotNil(t, input.Element)
		assert.IsType(t, &WorkflowButtonBlockElement{}, input.Element)
	})

	t.Run("Accessory", func(t *testing.T) {
		var accessory Accessory
		err := json.Unmarshal([]byte(workflowButtonJSON), &accessory)
		require.NoError(t, err)
		require.NotNil(t, accessory.WorkflowButtonElement)
		assert.Equal(t, METWorkflowButton, accessory.WorkflowButtonElement.Type)
	})
}

// TestAllBlockElementTypesUnmarshal ensures every known MET* element type can be
// unmarshalled through BlockElements.UnmarshalJSON without error. This test acts
// as a safety net: when a new element type constant is added to block_element.go
// but its case is not added to the switch statement, this test will fail.
func TestAllBlockElementTypesUnmarshal(t *testing.T) {
	allTypes := []string{
		string(METCheckboxGroups),
		string(METImage),
		string(METButton),
		string(METOverflow),
		string(METDatepicker),
		string(METTimepicker),
		string(METDatetimepicker),
		string(METPlainTextInput),
		string(METRadioButtons),
		string(METRichTextInput),
		string(METEmailTextInput),
		string(METURLTextInput),
		string(METNumber),
		string(METFileInput),
		string(METFeedbackButtons),
		string(METIconButton),
		string(METWorkflowButton),
		OptTypeStatic,
		OptTypeExternal,
		OptTypeUser,
		OptTypeConversations,
		OptTypeChannels,
		MultiOptTypeStatic,
		MultiOptTypeExternal,
		MultiOptTypeUser,
		MultiOptTypeConversations,
		MultiOptTypeChannels,
	}

	for _, typ := range allTypes {
		t.Run(typ, func(t *testing.T) {
			elemJSON := fmt.Sprintf(`[{"type": "%s"}]`, typ)
			var elements BlockElements
			err := json.Unmarshal([]byte(elemJSON), &elements)
			require.NoError(t, err, "BlockElements.UnmarshalJSON should handle type %q", typ)
			require.Len(t, elements.ElementSet, 1)
		})
	}
}

// TestAllAccessoryTypesRoundTrip ensures every Accessory field can survive a
// marshal→unmarshal round trip. When a new field is added to the Accessory struct
// and wired into NewAccessory but not into Accessory.UnmarshalJSON, this test
// will fail.
func TestAllAccessoryTypesRoundTrip(t *testing.T) {
	text := NewTextBlockObject("plain_text", "test", false, false)
	workflow := &Workflow{Trigger: &WorkflowTrigger{URL: "https://example.com"}}

	cases := []struct {
		name    string
		element BlockElement
		check   func(t *testing.T, a *Accessory)
	}{
		{"image", &ImageBlockElement{Type: METImage}, func(t *testing.T, a *Accessory) { assert.NotNil(t, a.ImageElement) }},
		{"button", &ButtonBlockElement{Type: METButton}, func(t *testing.T, a *Accessory) { assert.NotNil(t, a.ButtonElement) }},
		{"overflow", &OverflowBlockElement{Type: METOverflow}, func(t *testing.T, a *Accessory) { assert.NotNil(t, a.OverflowElement) }},
		{"datepicker", &DatePickerBlockElement{Type: METDatepicker}, func(t *testing.T, a *Accessory) { assert.NotNil(t, a.DatePickerElement) }},
		{"timepicker", &TimePickerBlockElement{Type: METTimepicker}, func(t *testing.T, a *Accessory) { assert.NotNil(t, a.TimePickerElement) }},
		{"plain_text_input", &PlainTextInputBlockElement{Type: METPlainTextInput}, func(t *testing.T, a *Accessory) { assert.NotNil(t, a.PlainTextInputElement) }},
		{"rich_text_input", &RichTextInputBlockElement{Type: METRichTextInput}, func(t *testing.T, a *Accessory) { assert.NotNil(t, a.RichTextInputElement) }},
		{"radio_buttons", &RadioButtonsBlockElement{Type: METRadioButtons}, func(t *testing.T, a *Accessory) { assert.NotNil(t, a.RadioButtonsElement) }},
		{"static_select", &SelectBlockElement{Type: OptTypeStatic}, func(t *testing.T, a *Accessory) { assert.NotNil(t, a.SelectElement) }},
		{"multi_static_select", &MultiSelectBlockElement{Type: MultiOptTypeStatic}, func(t *testing.T, a *Accessory) { assert.NotNil(t, a.MultiSelectElement) }},
		{"checkboxes", &CheckboxGroupsBlockElement{Type: METCheckboxGroups}, func(t *testing.T, a *Accessory) { assert.NotNil(t, a.CheckboxGroupsBlockElement) }},
		{"workflow_button", NewWorkflowButtonBlockElement(text, workflow, "action1"), func(t *testing.T, a *Accessory) { assert.NotNil(t, a.WorkflowButtonElement) }},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			accessory := NewAccessory(tc.element)

			data, err := json.Marshal(accessory)
			require.NoError(t, err, "Marshal failed for %s", tc.name)

			var unmarshalled Accessory
			err = json.Unmarshal(data, &unmarshalled)
			require.NoError(t, err, "Unmarshal failed for %s", tc.name)

			tc.check(t, &unmarshalled)
		})
	}
}
