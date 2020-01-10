package slack

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewDividerBlock(t *testing.T) {

	dividerBlock := NewDividerBlock()
	assert.Equal(t, string(dividerBlock.Type), "divider")

}
