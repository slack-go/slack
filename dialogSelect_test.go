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
}

func TestStaticExternalDataSourceSelect(t *testing.T) {
}

func TestConversationSelect(t *testing.T) {
}

func TestChannelSelect(t *testing.T) {
}

func TestUserSelect(t *testing.T) {
}
