package slack

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

const (
	IDExampleSelectInput = "ae9642ae-a9ef-4394-904b-a5c7a83bf4a6"
	IDSelectOptionBlock  = "832bf7af-22ea-4acb-82e3-a0cc3722052b"
)

func TestNewConfigurationModalRequest(t *testing.T) {
	blocks := configModalBlocks()
	privateMetaData := "An optional string that will be sent to your app in view_submission and block_actions events. Max length of 3000 characters."
	externalID := "c4baf441-fbc1-4131-b349-7c8df0ae7df6"

	result := NewConfigurationModalRequest(blocks, privateMetaData, externalID)

	if result.ModalViewRequest.Title != nil {
		t.Fail()
	}
	if result.PrivateMetadata != privateMetaData {
		t.Fail()
	}
	if result.ExternalID != externalID {
		t.Fail()
	}
}

func TestGetInitialOptionFromWorkflowStepInput(t *testing.T) {
	options, testOption := createOptionBlockObjects()
	selection := createSelection(options)

	scenarios := []struct {
		options        []*OptionBlockObject
		inputs         *WorkflowStepInputs
		expectedResult *OptionBlockObject
		expectedFlag   bool
	}{
		{
			options:        options,
			inputs:         createWorkflowStepInputs1(),
			expectedResult: &OptionBlockObject{},
			expectedFlag:   false,
		},
		{
			options:        []*OptionBlockObject{},
			inputs:         createWorkflowStepInputs4(testOption.Value),
			expectedResult: &OptionBlockObject{},
			expectedFlag:   false,
		},
		{
			options:        options,
			inputs:         createWorkflowStepInputs2(),
			expectedResult: &OptionBlockObject{},
			expectedFlag:   false,
		},
		{
			options:        options,
			inputs:         createWorkflowStepInputs3(),
			expectedResult: &OptionBlockObject{},
			expectedFlag:   false,
		},
		{
			options:        options,
			inputs:         createWorkflowStepInputs4(testOption.Value),
			expectedResult: testOption,
			expectedFlag:   true,
		},
	}

	for _, scenario := range scenarios {
		result, ok := GetInitialOptionFromWorkflowStepInput(selection, scenario.inputs, scenario.options)
		if ok != scenario.expectedFlag {
			t.Fail()
		}

		if !cmp.Equal(result, scenario.expectedResult) {
			t.Fail()
		}
	}
}

func createOptionBlockObjects() ([]*OptionBlockObject, *OptionBlockObject) {
	var options []*OptionBlockObject
	options = append(
		options,
		NewOptionBlockObject("one", NewTextBlockObject("plain_text", "One", false, false), nil),
	)

	option2 := NewOptionBlockObject("two", NewTextBlockObject("plain_text", "Two", false, false), nil)
	options = append(
		options,
		option2,
	)

	options = append(
		options,
		NewOptionBlockObject("three", NewTextBlockObject("plain_text", "Three", false, false), nil),
	)

	return options, option2
}

func createSelection(options []*OptionBlockObject) *SelectBlockElement {
	return NewOptionsSelectBlockElement(
		"static_select",
		NewTextBlockObject("plain_text", "your choice", false, false),
		IDExampleSelectInput,
		options...,
	)
}

func configModalBlocks() Blocks {
	headerText := NewTextBlockObject("mrkdwn", "Hello World!\nThis is your workflow step app configuration view", false, false)
	headerSection := NewSectionBlock(headerText, nil, nil)

	options, _ := createOptionBlockObjects()

	selection := createSelection(options)

	inputBlock := NewInputBlock(
		IDSelectOptionBlock,
		NewTextBlockObject("plain_text", "Select an option", false, false),
		NewTextBlockObject("plain_text", "Hint", false, false),
		selection,
	)

	blocks := Blocks{
		BlockSet: []Block{
			headerSection,
			inputBlock,
		},
	}

	return blocks
}

func createWorkflowStepInputs1() *WorkflowStepInputs {
	return &WorkflowStepInputs{}
}

func createWorkflowStepInputs2() *WorkflowStepInputs {
	return &WorkflowStepInputs{
		"test": WorkflowStepInputElement{
			Value:                   "random-string",
			SkipVariableReplacement: false,
		},
		"123-test": WorkflowStepInputElement{
			Value:                   "another-string",
			SkipVariableReplacement: false,
		},
	}
}

func createWorkflowStepInputs3() *WorkflowStepInputs {
	return &WorkflowStepInputs{
		"test": WorkflowStepInputElement{
			Value:                   "random-string",
			SkipVariableReplacement: false,
		},
		"123-test": WorkflowStepInputElement{
			Value:                   "another-string",
			SkipVariableReplacement: false,
		},
		IDExampleSelectInput: WorkflowStepInputElement{
			Value:                   "lorem-ipsum",
			SkipVariableReplacement: true,
		},
	}
}

func createWorkflowStepInputs4(optionValue string) *WorkflowStepInputs {
	return &WorkflowStepInputs{
		"test": WorkflowStepInputElement{
			Value:                   "random-string",
			SkipVariableReplacement: false,
		},
		"123-test": WorkflowStepInputElement{
			Value:                   "another-string",
			SkipVariableReplacement: false,
		},
		IDExampleSelectInput: WorkflowStepInputElement{
			Value:                   optionValue,
			SkipVariableReplacement: false,
		},
	}
}
