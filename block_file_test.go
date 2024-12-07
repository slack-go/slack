package slack

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewFileBlock(t *testing.T) {
	fileBlock := NewFileBlock("test", "external_id", "source")

	assert.Equal(t, fileBlock.BlockType(), MBTFile)
	assert.Equal(t, string(fileBlock.Type), "file")
	assert.Equal(t, fileBlock.BlockID, "test")
	assert.Equal(t, fileBlock.ID(), "test")
	assert.Equal(t, fileBlock.ExternalID, "external_id")
	assert.Equal(t, fileBlock.Source, "source")
}
