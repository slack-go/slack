package slack

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUnknownBlockRoundTrip(t *testing.T) {
	input := `[{"type":"some_future_block","block_id":"fb1","custom_field":"value","nested":{"key":"val"}}]`

	var blocks Blocks
	err := json.Unmarshal([]byte(input), &blocks)
	require.NoError(t, err)
	require.Len(t, blocks.BlockSet, 1)

	assert.Equal(t, MessageBlockType("some_future_block"), blocks.BlockSet[0].BlockType())
	assert.Equal(t, "fb1", blocks.BlockSet[0].ID())

	output, err := json.Marshal(blocks)
	require.NoError(t, err)
	assert.JSONEq(t, input, string(output))
}
