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

	assert.Equal(t, MBTContext, contextBlock.BlockType())
	assert.Equal(t, "context", string(contextBlock.Type))
	assert.Equal(t, "test", contextBlock.BlockID)
	assert.Equal(t, "test", contextBlock.ID())
	assert.Equal(t, 2, len(contextBlock.ContextElements.Elements))
}
