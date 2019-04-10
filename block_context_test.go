package slack

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewContextBlock(t *testing.T) {

	locationPinImage := NewImageBlockElement("https://api.slack.com/img/blocks/bkb_template_images/tripAgentLocationMarker.png", "Location Pin Icon")
	textExample := NewTextBlockObject("plain_text", "Location: Central Business District", true, false)

	contextElements := ContextElements{
		ContextElementSet: []MixedElement{locationPinImage, textExample},
	}

	contextBlock := NewContextBlock("test", contextElements)
	assert.Equal(t, string(contextBlock.Type), "context")
	assert.Equal(t, contextBlock.BlockID, "test")
	assert.Equal(t, len(contextBlock.Elements.ContextElementSet), 1)

}
