package slack

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestImageURLForNewImageBlock(t *testing.T) {
	imageText := NewTextBlockObject("plain_text", "Location", false, false)
	imageBlock := NewImageBlock("https://api.slack.com/img/blocks/bkb_template_images/tripAgentLocationMarker.png", "Marker", "test", imageText)

	assert.Equal(t, MBTImage, imageBlock.BlockType())
	assert.Equal(t, "image", string(imageBlock.Type))
	assert.Equal(t, "plain_text", imageBlock.Title.Type)
	assert.Equal(t, "test", imageBlock.ID())
	assert.Equal(t, "test", imageBlock.BlockID)
	assert.Contains(t, imageBlock.Title.Text, "Location")
	assert.Contains(t, imageBlock.ImageURL, "tripAgentLocationMarker.png")
}

func TestSlackFileForNewImageBlock(t *testing.T) {
	imageText := NewTextBlockObject("plain_text", "Location", false, false)
	slackFile := &SlackFileObject{URL: "https://api.slack.com/img/blocks/bkb_template_images/tripAgentLocationMarker.png"}
	imageBlock := NewImageBlockSlackFile(slackFile, "Marker", "test", imageText)

	assert.Equal(t, MBTImage, imageBlock.BlockType())
	assert.Equal(t, "image", string(imageBlock.Type))
	assert.Equal(t, "plain_text", imageBlock.Title.Type)
	assert.Equal(t, "test", imageBlock.ID())
	assert.Equal(t, "test", imageBlock.BlockID)
	assert.Contains(t, imageBlock.Title.Text, "Location")
	assert.Contains(t, imageBlock.SlackFile.URL, "tripAgentLocationMarker.png")
}
