package slack

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewInputBlock(t *testing.T) {
	label := NewTextBlockObject("plain_text", "Input", false, false)
	inputElement := NewPlainTextInputBlockElement(nil, "input_123")

	inputBlock := NewInputBlock("test", label, inputElement)
	assert.Equal(t, string(inputBlock.Type), "input")
	assert.Equal(t, inputBlock.BlockID, "test")
	assert.Equal(t, string(inputBlock.Element.ElementType()), "plain_text_input")
}
