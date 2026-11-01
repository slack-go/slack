package slack

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewVideoBlock(t *testing.T) {
	videoTitle := NewTextBlockObject("plain_text", "VideoTitle", false, false)
	videoBlock := NewVideoBlock(
		"https://example.com/example.mp4",
		"https://example.com/thumbnail.png",
		"alternative text", "blockID", videoTitle)

	assert.Equal(t, MBTVideo, videoBlock.Type)
	assert.Equal(t, "video", string(videoBlock.Type))
	assert.Equal(t, "plain_text", videoBlock.Title.Type)
	assert.Equal(t, "blockID", videoBlock.BlockID)
	assert.Equal(t, "blockID", videoBlock.ID())
	assert.Contains(t, videoBlock.Title.Text, "VideoTitle")
	assert.Contains(t, videoBlock.VideoURL, "example.mp4")
}
