package slack

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewFileBlock(t *testing.T) {
	fileBlock := NewFileBlock("test", "external_id", "source")

	assert.Equal(t, MBTFile, fileBlock.BlockType())
	assert.Equal(t, "file", string(fileBlock.Type))
	assert.Equal(t, "test", fileBlock.BlockID)
	assert.Equal(t, "test", fileBlock.ID())
	assert.Equal(t, "external_id", fileBlock.ExternalID)
	assert.Equal(t, "source", fileBlock.Source)
}
