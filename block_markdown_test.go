package slack

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewMarkdownBlock(t *testing.T) {
	markdownBlock := NewMarkdownBlock("test", "*asfd*")

	assert.Equal(t, markdownBlock.BlockType(), MBTMarkdown)
	assert.Equal(t, string(markdownBlock.Type), "markdown")
	assert.Equal(t, markdownBlock.ID(), "test")
	assert.Equal(t, markdownBlock.BlockID, "test")
	assert.Equal(t, markdownBlock.Text, "*asfd*")
}
