package slack

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewBlockMessage(t *testing.T) {

	dividerBlock := NewDividerBlock()
	blockMessage := NewBlockMessage(dividerBlock)

	assert.Len(t, blockMessage.Msg.Blocks.BlockSet, 1)

}
