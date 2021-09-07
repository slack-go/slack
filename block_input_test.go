package slack

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewInputBlock(t *testing.T) {
	label := NewTextBlockObject("plain_text", "label", false, false)
	element := NewDatePickerBlockElement("action_id")
	hint := NewTextBlockObject("plain_text", "hint", false, false)
	inputBlock := NewInputBlock("test", label, hint, element)
	assert.Equal(t, string(inputBlock.Type), "input")
	assert.Equal(t, inputBlock.BlockID, "test")
	assert.Equal(t, inputBlock.Label, label)
	assert.Equal(t, inputBlock.Element, element)
}
