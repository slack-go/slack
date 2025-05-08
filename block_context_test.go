package slack

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewContextBlock(t *testing.T) {
	locationPinImage := NewImageBlockElement("https://api.slack.com/img/blocks/bkb_template_images/tripAgentLocationMarker.png", "Location Pin Icon")
	textExample := NewTextBlockObject("plain_text", "Location: Central Business District", true, false)
	elements := []MixedElement{locationPinImage, textExample}
	contextBlock := NewContextBlock("test", elements...)

	assert.Equal(t, contextBlock.BlockType(), MBTContext)
	assert.Equal(t, string(contextBlock.Type), "context")
	assert.Equal(t, contextBlock.BlockID, "test")
	assert.Equal(t, contextBlock.ID(), "test")
	assert.Equal(t, len(contextBlock.ContextElements.Elements), 2)
}
