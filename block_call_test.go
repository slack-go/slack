package slack

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewCallBlock(t *testing.T) {
	callBlock := NewCallBlock("ACallID")

	assert.Equal(t, callBlock.BlockType(), MBTCall)
	assert.Equal(t, string(callBlock.Type), "call")
	assert.Equal(t, callBlock.CallID, "ACallID")
	assert.Equal(t, callBlock.BlockID, "")
	assert.Equal(t, callBlock.ID(), "")
}
