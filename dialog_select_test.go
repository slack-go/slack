package slack

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOptionsFromArray(t *testing.T) {
	options := []string{"opt 1"}
	expectedOptions := optionsFromArray(options)
	assert.Equal(t, len(options), len(expectedOptions))

	firstOption := expectedOptions[0]
	assert.Equal(t, "opt 1", firstOption.Label)
	assert.Equal(t, "opt 1", firstOption.Value)
}

func TestOptionsFromMap(t *testing.T) {
	options := make(map[string]string)
	options["key"] = "myValue"

	selectOptions := optionsFromMap(options)
	assert.Equal(t, 1, len(selectOptions))

	firstOption := selectOptions[0]
	assert.Equal(t, "key", firstOption.Label)
	assert.Equal(t, "myValue", firstOption.Value)
}

func TestStaticSelectFromArray(t *testing.T) {
	name := "static select"
	label := "Static Select Label"
	options := []string{"opt 1", "opt 2", "opt 3"}
	expectedOptions := optionsFromArray(options)

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
	expectedOptions := optionsFromMap(optionsMap)

	selectInput := NewStaticSelectDialogInput(name, label, expectedOptions)
	assert.Equal(t, name, selectInput.Name)
	assert.Equal(t, label, selectInput.Label)
	assert.Equal(t, expectedOptions, selectInput.Options)
}

func TestStaticGroupedSelect(t *testing.T) {
	group1 := make(map[string]string)
	group1["G1_O1"] = "First (1)"
	group1["G1_O2"] = "Second (1)"

	group2 := make(map[string]string)
	group2["G2_O1"] = "First (2)"
	group2["G2_O2"] = "Second (2)"
	group2["G2_O3"] = "Third (2)"

	groups := make(map[string]map[string]string)
	groups["Group 1"] = group1
	groups["Group 2"] = group2

	groupSelect := NewGroupedSelectDialogInput("groupSelect", "User Label", groups)
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
