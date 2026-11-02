package slack

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewActionBlock(t *testing.T) {
	approveBtnTxt := NewTextBlockObject("plain_text", "Approve", false, false)
	approveBtn := NewButtonBlockElement("", "click_me_123", approveBtnTxt)
	actionBlock := NewActionBlock("test", approveBtn)

	assert.Equal(t, MBTAction, actionBlock.BlockType())
	assert.Equal(t, "actions", string(actionBlock.Type))
	assert.Equal(t, "test", actionBlock.BlockID)
	assert.Equal(t, "test", actionBlock.ID())
	assert.Len(t, actionBlock.Elements.ElementSet, 1)
}
