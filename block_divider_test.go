package slack

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewDividerBlock(t *testing.T) {
	dividerBlock := NewDividerBlock()

	assert.Equal(t, MBTDivider, dividerBlock.BlockType())
	assert.Equal(t, "divider", string(dividerBlock.Type))
	assert.Empty(t, dividerBlock.BlockID)
	assert.Empty(t, dividerBlock.ID())
}
