package slack

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewCallBlock(t *testing.T) {
	callBlock := NewCallBlock("ACallID")

	assert.Equal(t, MBTCall, callBlock.BlockType())
	assert.Equal(t, "call", string(callBlock.Type))
	assert.Equal(t, "ACallID", callBlock.CallID)
	assert.Equal(t, "", callBlock.BlockID)
	assert.Equal(t, "", callBlock.ID())
}

func TestNewCallBlockWithBlockID(t *testing.T) {
	callBlock := NewCallBlock("ACallID", CallBlockOptionBlockID("block-123"))

	assert.Equal(t, MBTCall, callBlock.BlockType())
	assert.Equal(t, "ACallID", callBlock.CallID)
	assert.Equal(t, "block-123", callBlock.BlockID)
	assert.Equal(t, "block-123", callBlock.ID())
}

func TestCallBlockJSONRoundTrip(t *testing.T) {
	// Create block with call data using the V1 structure
	callBlock := NewCallBlock("R123", CallBlockOptionBlockID("block-1"))
	callBlock.Call = &CallBlockData{
		V1: &CallBlockDataV1{
			ID:      "R123",
			Name:    "Team Standup",
			JoinURL: "https://example.com/call",
		},
		MediaBackendType: "platform_call",
	}

	// Marshal to JSON
	data, err := json.Marshal(callBlock)
	require.NoError(t, err)

	// Verify expected JSON structure
	var jsonMap map[string]any
	err = json.Unmarshal(data, &jsonMap)
	require.NoError(t, err)
	assert.Equal(t, "call", jsonMap["type"])
	assert.Equal(t, "R123", jsonMap["call_id"])
	assert.Equal(t, "block-1", jsonMap["block_id"])

	// Verify nested structure
	callData, ok := jsonMap["call"].(map[string]any)
	require.True(t, ok, "call should be a map")
	assert.Equal(t, "platform_call", callData["media_backend_type"])
	v1Data, ok := callData["v1"].(map[string]any)
	require.True(t, ok, "v1 should be a map")
	assert.Equal(t, "R123", v1Data["id"])

	// Unmarshal back
	var parsed CallBlock
	err = json.Unmarshal(data, &parsed)
	require.NoError(t, err)

	assert.Equal(t, MBTCall, parsed.Type)
	assert.Equal(t, callBlock.CallID, parsed.CallID)
	assert.Equal(t, callBlock.BlockID, parsed.BlockID)
	require.NotNil(t, parsed.Call.V1)
	assert.Equal(t, callBlock.Call.V1.ID, parsed.Call.V1.ID)
	assert.Equal(t, callBlock.Call.V1.Name, parsed.Call.V1.Name)
	assert.Equal(t, callBlock.Call.V1.JoinURL, parsed.Call.V1.JoinURL)
	assert.Equal(t, callBlock.Call.MediaBackendType, parsed.Call.MediaBackendType)
}

func TestCallBlockInBlocks(t *testing.T) {
	// Test that call block can be unmarshalled as part of a Blocks collection
	// Using the actual Slack structure with v1 wrapper
	jsonData := []byte(`[
		{
			"type": "call",
			"block_id": "call-block-1",
			"call_id": "R123456",
			"api_decoration_available": false,
			"call": {
				"v1": {
					"id": "R123456",
					"name": "Team Standup",
					"join_url": "https://example.com/join/123",
					"desktop_app_join_url": "slack://call/123",
					"date_start": 1769457524,
					"date_end": 0,
					"is_dm_call": false,
					"was_rejected": false,
					"was_missed": false,
					"was_accepted": false,
					"has_ended": false
				},
				"media_backend_type": "platform_call"
			}
		}
	]`)

	var blocks Blocks
	err := json.Unmarshal(jsonData, &blocks)
	require.NoError(t, err)
	require.Len(t, blocks.BlockSet, 1)

	assert.Equal(t, MBTCall, blocks.BlockSet[0].BlockType())
	assert.Equal(t, "call-block-1", blocks.BlockSet[0].ID())

	callBlock, ok := blocks.BlockSet[0].(*CallBlock)
	require.True(t, ok, "expected *CallBlock, got %T", blocks.BlockSet[0])
	assert.Equal(t, "R123456", callBlock.CallID)
	assert.False(t, callBlock.APIDecorationAvailable)
	assert.Equal(t, "platform_call", callBlock.Call.MediaBackendType)
	require.NotNil(t, callBlock.Call.V1)
	assert.Equal(t, "R123456", callBlock.Call.V1.ID)
	assert.Equal(t, "Team Standup", callBlock.Call.V1.Name)
	assert.Equal(t, "https://example.com/join/123", callBlock.Call.V1.JoinURL)
	assert.Equal(t, "slack://call/123", callBlock.Call.V1.DesktopAppJoinURL)
	assert.Equal(t, int64(1769457524), callBlock.Call.V1.DateStart)
	assert.False(t, callBlock.Call.V1.HasEnded)
}

func TestCallBlockWithParticipants(t *testing.T) {
	jsonData := []byte(`{
		"type": "call",
		"call_id": "R789",
		"call": {
			"v1": {
				"id": "R789",
				"name": "Design Review",
				"date_start": 0,
				"date_end": 0,
				"active_participants": [
					{"slack_id": "U123", "display_name": "Alice"},
					{"slack_id": "U456", "display_name": "Bob"}
				],
				"all_participants": [
					{"slack_id": "U123", "display_name": "Alice"},
					{"slack_id": "U456", "display_name": "Bob"},
					{"slack_id": "U789", "display_name": "Charlie"}
				],
				"is_dm_call": false,
				"was_rejected": false,
				"was_missed": false,
				"was_accepted": false,
				"has_ended": false
			}
		}
	}`)

	var callBlock CallBlock
	err := json.Unmarshal(jsonData, &callBlock)
	require.NoError(t, err)

	assert.Equal(t, "R789", callBlock.CallID)
	require.NotNil(t, callBlock.Call.V1)
	assert.Equal(t, "Design Review", callBlock.Call.V1.Name)
	require.Len(t, callBlock.Call.V1.ActiveParticipants, 2)
	assert.Equal(t, "U123", callBlock.Call.V1.ActiveParticipants[0].SlackID)
	assert.Equal(t, "Alice", callBlock.Call.V1.ActiveParticipants[0].DisplayName)
	assert.Equal(t, "U456", callBlock.Call.V1.ActiveParticipants[1].SlackID)
	assert.Equal(t, "Bob", callBlock.Call.V1.ActiveParticipants[1].DisplayName)

	require.Len(t, callBlock.Call.V1.AllParticipants, 3)
	assert.Equal(t, "U789", callBlock.Call.V1.AllParticipants[2].SlackID)
	assert.Equal(t, "Charlie", callBlock.Call.V1.AllParticipants[2].DisplayName)
}

func TestCallBlockWithAppIcons(t *testing.T) {
	// Test parsing of app icon URLs as seen in real Zoom integration
	jsonData := []byte(`{
		"type": "call",
		"call_id": "R0ABF31RWGH",
		"block_id": "+cgoe",
		"api_decoration_available": false,
		"call": {
			"v1": {
				"id": "R0ABF31RWGH",
				"app_id": "A5GE9BMQC",
				"app_icon_urls": {
					"image_32": "https://example.com/icon_32.png",
					"image_48": "https://example.com/icon_48.png",
					"image_72": "https://example.com/icon_72.png",
					"image_192": "https://example.com/icon_192.png"
				},
				"date_start": 1769457524,
				"date_end": 0,
				"display_id": "863-5835-0956",
				"join_url": "https://zoom.us/j/123",
				"desktop_app_join_url": "zoommtg://zoom.us/join?confno=123",
				"name": "Zoom meeting started by user",
				"created_by": "U0ABF1CJPG9",
				"channels": ["C0AACTVQ2EB"],
				"is_dm_call": false,
				"was_rejected": false,
				"was_missed": false,
				"was_accepted": false,
				"has_ended": false
			},
			"media_backend_type": "platform_call"
		}
	}`)

	var callBlock CallBlock
	err := json.Unmarshal(jsonData, &callBlock)
	require.NoError(t, err)

	assert.Equal(t, "R0ABF31RWGH", callBlock.CallID)
	assert.Equal(t, "+cgoe", callBlock.BlockID)
	assert.False(t, callBlock.APIDecorationAvailable)
	assert.Equal(t, "platform_call", callBlock.Call.MediaBackendType)

	require.NotNil(t, callBlock.Call.V1)
	v1 := callBlock.Call.V1
	assert.Equal(t, "R0ABF31RWGH", v1.ID)
	assert.Equal(t, "A5GE9BMQC", v1.AppID)
	assert.Equal(t, "863-5835-0956", v1.DisplayID)
	assert.Equal(t, "Zoom meeting started by user", v1.Name)
	assert.Equal(t, "U0ABF1CJPG9", v1.CreatedBy)
	assert.Equal(t, int64(1769457524), v1.DateStart)
	assert.Equal(t, int64(0), v1.DateEnd)
	assert.False(t, v1.HasEnded)
	require.Len(t, v1.Channels, 1)
	assert.Equal(t, "C0AACTVQ2EB", v1.Channels[0])

	require.NotNil(t, v1.AppIconURLs)
	assert.Equal(t, "https://example.com/icon_32.png", v1.AppIconURLs.Image32)
	assert.Equal(t, "https://example.com/icon_48.png", v1.AppIconURLs.Image48)
	assert.Equal(t, "https://example.com/icon_72.png", v1.AppIconURLs.Image72)
	assert.Equal(t, "https://example.com/icon_192.png", v1.AppIconURLs.Image192)
}
