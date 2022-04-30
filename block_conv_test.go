package slack

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBlocksJSONMarshalAndUnmarshal(t *testing.T) {
	input := Blocks{
		BlockSet: []Block{
			TextBlockObject{
				Text: "test",
			},
		},
	}

	actualFromValue, err := json.Marshal(input)
	assert.NoError(t, err)
	actualFromPtr, err := json.Marshal(&input)
	assert.NoError(t, err)

	assert.Equal(t, actualFromValue, actualFromPtr)
}

func TestBlockElementsElementsJSONMarshalAndUnmarshal(t *testing.T) {
	input := BlockElements{
		ElementSet: []BlockElement{
			ImageBlockElement{
				ImageURL: "test",
			},
		},
	}

	actualFromValue, err := json.Marshal(input)
	assert.NoError(t, err)
	actualFromPtr, err := json.Marshal(&input)
	assert.NoError(t, err)

	assert.Equal(t, actualFromValue, actualFromPtr)
}

func TestAccessoryJSONMarshalAndUnmarshal(t *testing.T) {
	input := Accessory{
		ImageElement: &ImageBlockElement{
			ImageURL: "test",
		},
	}

	actualFromValue, err := json.Marshal(input)
	assert.NoError(t, err)
	actualFromPtr, err := json.Marshal(&input)
	assert.NoError(t, err)

	assert.Equal(t, actualFromValue, actualFromPtr)
}

func TestContextElementsJSONMarshalAndUnmarshal(t *testing.T) {
	input := ContextElements{
		Elements: []MixedElement{
			ImageBlockElement{
				ImageURL: "test",
			},
		},
	}

	actualFromValue, err := json.Marshal(input)
	assert.NoError(t, err)
	actualFromPtr, err := json.Marshal(&input)
	assert.NoError(t, err)

	assert.Equal(t, actualFromValue, actualFromPtr)
}
