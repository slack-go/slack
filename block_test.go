package slack

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewBlockMessage(t *testing.T) {

	dividerBlock := NewDividerBlock()
	blockMessage := NewBlockMessage(dividerBlock)

	assert.Equal(t, 1, len(blockMessage.Msg.Blocks.BlockSet))

}
