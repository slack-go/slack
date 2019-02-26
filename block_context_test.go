package slack

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewContextBlock(t *testing.T) {

	locationPinImage := NewImageBlockObject("https://api.slack.com/img/blocks/bkb_template_images/tripAgentLocationMarker.png", "Location Pin Icon")
	textExample := NewTextBlockObject("plain_text", "Location: Central Business District", true, false)

	actionBlock := NewContextBlock("test", locationPinImage, textExample)
	assert.Equal(t, string(actionBlock.Type), "context")
	assert.Equal(t, actionBlock.BlockID, "test")
	assert.Equal(t, len(actionBlock.Elements), 2)

}
