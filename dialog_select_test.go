package slack

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func selectOptionsFromArray(options ...string) []DialogSelectOption {
	selectOptions := make([]DialogSelectOption, len(options))
	for idx, value := range options {
		selectOptions[idx] = DialogSelectOption{
			Label: value,
			Value: value,
		}
	}
	return selectOptions
}

func selectOptionsFromMap(options map[string]string) []DialogSelectOption {
	selectOptions := make([]DialogSelectOption, len(options))
	idx := 0
	var option DialogSelectOption
	for key, value := range options {
		option = DialogSelectOption{
			Label: key,
			Value: value,
		}
		selectOptions[idx] = option
		idx++
	}
	return selectOptions
}

func TestSelectOptionsFromArray(t *testing.T) {
	options := []string{"opt 1"}
	expectedOptions := selectOptionsFromArray(options...)
	assert.Equal(t, len(options), len(expectedOptions))

	firstOption := expectedOptions[0]
	assert.Equal(t, "opt 1", firstOption.Label)
	assert.Equal(t, "opt 1", firstOption.Value)
}

func TestOptionsFromMap(t *testing.T) {
	options := make(map[string]string)
	options["key"] = "myValue"

	selectOptions := selectOptionsFromMap(options)
	assert.Equal(t, 1, len(selectOptions))

	firstOption := selectOptions[0]
	assert.Equal(t, "key", firstOption.Label)
	assert.Equal(t, "myValue", firstOption.Value)
}

func TestStaticSelectFromArray(t *testing.T) {
	name := "static select"
	label := "Static Select Label"
	expectedOptions := selectOptionsFromArray("opt 1", "opt 2", "opt 3")

	selectInput := NewStaticSelectDialogInput(name, label, expectedOptions)
	assert.Equal(t, name, selectInput.Name)
	assert.Equal(t, label, selectInput.Label)
	assert.Equal(t, expectedOptions, selectInput.Options)
}

func TestStaticSelectFromDictionary(t *testing.T) {
	name := "static select"
	label := "Static Select Label"

	optionsMap := make(map[string]string)
	optionsMap["option_1"] = "First"
	optionsMap["option_2"] = "Second"
	optionsMap["option_3"] = "Third"
	expectedOptions := selectOptionsFromMap(optionsMap)

	selectInput := NewStaticSelectDialogInput(name, label, expectedOptions)
	assert.Equal(t, name, selectInput.Name)
	assert.Equal(t, label, selectInput.Label)
	assert.Equal(t, expectedOptions, selectInput.Options)
}

func TestNewDialogOptionGroup(t *testing.T) {
	expectedOptions := selectOptionsFromArray("option_1", "option_2")

	label := "GroupLabel"
	optionGroup := NewDialogOptionGroup(label, expectedOptions...)

	assert.Equal(t, label, optionGroup.Label)
	assert.Equal(t, expectedOptions, optionGroup.Options)

}

func TestStaticGroupedSelect(t *testing.T) {

	groupOpt1 := NewDialogOptionGroup("group1", selectOptionsFromArray("G1_01", "G1_02")...)
	groupOpt2 := NewDialogOptionGroup("group2", selectOptionsFromArray("G2_01", "G2_02", "G2_03")...)

	options := []DialogOptionGroup{groupOpt1, groupOpt2}

	groupSelect := NewGroupedSelectDialogInput("groupSelect", "User Label", options)
	assert.Equal(t, InputTypeSelect, groupSelect.Type)
	assert.Equal(t, "groupSelect", groupSelect.Name)
	assert.Equal(t, "User Label", groupSelect.Label)
	assert.Nil(t, groupSelect.Options)
	assert.NotNil(t, groupSelect.OptionGroups)
	assert.Equal(t, 2, len(groupSelect.OptionGroups))
}

func TestConversationSelect(t *testing.T) {
	convoSelect := NewConversationsSelect("", "")
	assert.Equal(t, InputTypeSelect, convoSelect.Type)
	assert.Equal(t, DialogDataSourceConversations, convoSelect.DataSource)
}

func TestChannelSelect(t *testing.T) {
	convoSelect := NewChannelsSelect("", "")
	assert.Equal(t, InputTypeSelect, convoSelect.Type)
	assert.Equal(t, DialogDataSourceChannels, convoSelect.DataSource)
}

func TestUserSelect(t *testing.T) {
	convoSelect := NewUsersSelect("", "")
	assert.Equal(t, InputTypeSelect, convoSelect.Type)
	assert.Equal(t, DialogDataSourceUsers, convoSelect.DataSource)
}
