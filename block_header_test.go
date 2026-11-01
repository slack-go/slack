package slack

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewHeaderBlock(t *testing.T) {
	textInfo := NewTextBlockObject("plain_text", "This is quite the header", false, false)
	headerBlock := NewHeaderBlock(textInfo, HeaderBlockOptionBlockID("test_block"))

	assert.Equal(t, MBTHeader, headerBlock.BlockType())
	assert.Equal(t, "header", string(headerBlock.Type))
	assert.Equal(t, "test_block", headerBlock.ID())
	assert.Equal(t, "test_block", headerBlock.BlockID)
	assert.Equal(t, "plain_text", headerBlock.Text.Type)
	assert.Contains(t, headerBlock.Text.Text, "quite the header")
}

func TestNewHeaderBlockWithLevel(t *testing.T) {
	textInfo := NewTextBlockObject("plain_text", "This is quite the header", false, false)
	headerBlock := NewHeaderBlock(textInfo, HeaderBlockOptionLevel(2))

	assert.Equal(t, MBTHeader, headerBlock.BlockType())
	assert.Equal(t, 2, headerBlock.Level)
}

// TestNewHeaderBlockWithNilOption reproduces issue #1236: passing nil as an
// option to NewHeaderBlock causes a nil pointer dereference panic.
func TestNewHeaderBlockWithNilOption(t *testing.T) {
	textInfo := NewTextBlockObject("plain_text", "Header text", false, false)

	assert.NotPanics(t, func() {
		NewHeaderBlock(textInfo, nil)
	}, "NewHeaderBlock should not panic when nil is passed as an option")
}
