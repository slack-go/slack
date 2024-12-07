package slack

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewHeaderBlock(t *testing.T) {
	textInfo := NewTextBlockObject("plain_text", "This is quite the header", false, false)
	headerBlock := NewHeaderBlock(textInfo, HeaderBlockOptionBlockID("test_block"))

	assert.Equal(t, headerBlock.BlockType(), MBTHeader)
	assert.Equal(t, string(headerBlock.Type), "header")
	assert.Equal(t, headerBlock.ID(), "test_block")
	assert.Equal(t, headerBlock.BlockID, "test_block")
	assert.Equal(t, headerBlock.Text.Type, "plain_text")
	assert.Contains(t, headerBlock.Text.Text, "quite the header")
}
