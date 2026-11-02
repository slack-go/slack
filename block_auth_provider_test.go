package slack

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAuthProviderBlock(t *testing.T) {
	block := NewAuthProviderBlock("example", NewTextBlockObject(PlainTextType, "Example", false, false), false, nil, AuthProviderBlockOptionBlockID("auth"))
	assert.Equal(t, MBTAuthProvider, block.BlockType())
	assert.Equal(t, "auth", block.ID())
	encoded, err := json.Marshal(block)
	require.NoError(t, err)
	assert.JSONEq(t, `{"type":"auth_provider","block_id":"auth","value":"example","name":{"type":"plain_text","text":"Example","emoji":false},"is_connected":false}`, string(encoded))
}

func TestAuthProviderBlockRoundTrip(t *testing.T) {
	const payload = `[{"type":"auth_provider","block_id":"auth","value":"example","name":{"type":"plain_text","text":"Example"},"is_connected":true,"description":{"type":"plain_text","text":"Connect your account"},"icon_url":"https://example.com/icon.png","auth_url":"https://example.com/auth","disconnect_url":"https://example.com/disconnect","connected_identity":{"type":"plain_text","text":"alice"},"connected_context":{"type":"plain_text","text":"Example org"},"accessory":{"type":"overflow","action_id":"menu","options":[{"text":{"type":"plain_text","text":"Disconnect"},"value":"disconnect"}]},"fallback":[{"type":"section","text":{"type":"mrkdwn","text":"Connect to Example"}},{"type":"future_block","value":false}]}]`
	var blocks Blocks
	require.NoError(t, json.Unmarshal([]byte(payload), &blocks))
	require.Len(t, blocks.BlockSet, 1)
	block, ok := blocks.BlockSet[0].(*AuthProviderBlock)
	require.True(t, ok)
	require.NotNil(t, block.Fallback)
	require.Len(t, block.Fallback.BlockSet, 2)
	assert.IsType(t, &SectionBlock{}, block.Fallback.BlockSet[0])
	assert.IsType(t, &UnknownBlock{}, block.Fallback.BlockSet[1])
	encoded, err := json.Marshal(blocks)
	require.NoError(t, err)
	assert.JSONEq(t, payload, string(encoded))
	for _, fallback := range []string{"", `,"fallback":[]`} {
		raw := `[{"type":"auth_provider","value":"example","name":{"type":"plain_text","text":"Example"},"is_connected":false` + fallback + `}]`
		require.NoError(t, json.Unmarshal([]byte(raw), &blocks))
		encoded, err := json.Marshal(blocks)
		require.NoError(t, err)
		assert.JSONEq(t, raw, string(encoded))
	}
}
