package slack

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewActionBlock(t *testing.T) {
	approveBtnTxt := NewTextBlockObject("plain_text", "Approve", false, false)
	approveBtn := NewButtonBlockElement("", "click_me_123", approveBtnTxt)
	actionBlock := NewActionBlock("test", approveBtn)

	assert.Equal(t, actionBlock.BlockType(), MBTAction)
	assert.Equal(t, string(actionBlock.Type), "actions")
	assert.Equal(t, actionBlock.BlockID, "test")
	assert.Equal(t, actionBlock.ID(), "test")
	assert.Equal(t, len(actionBlock.Elements.ElementSet), 1)
}
