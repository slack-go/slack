package slack

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewMarkdownBlock(t *testing.T) {
	markdownBlock := NewMarkdownBlock("test", "*asfd*")

	assert.Equal(t, MBTMarkdown, markdownBlock.BlockType())
	assert.Equal(t, "markdown", string(markdownBlock.Type))
	assert.Equal(t, "test", markdownBlock.ID())
	assert.Equal(t, "test", markdownBlock.BlockID)
	assert.Equal(t, "*asfd*", markdownBlock.Text)
}
