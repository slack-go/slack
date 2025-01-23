package slack

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewDividerBlock(t *testing.T) {
	dividerBlock := NewDividerBlock()

	assert.Equal(t, dividerBlock.BlockType(), MBTDivider)
	assert.Equal(t, string(dividerBlock.Type), "divider")
	assert.Equal(t, dividerBlock.BlockID, "")
	assert.Equal(t, dividerBlock.ID(), "")
}
