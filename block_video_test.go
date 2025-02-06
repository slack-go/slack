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

	assert.Equal(t, videoBlock.Type, MBTVideo)
	assert.Equal(t, string(videoBlock.Type), "video")
	assert.Equal(t, videoBlock.Title.Type, "plain_text")
	assert.Equal(t, videoBlock.BlockID, "blockID")
	assert.Equal(t, videoBlock.ID(), "blockID")
	assert.Contains(t, videoBlock.Title.Text, "VideoTitle")
	assert.Contains(t, videoBlock.VideoURL, "example.mp4")
}
