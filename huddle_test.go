package slack

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHuddleRoomUnmarshal(t *testing.T) {
	jsonData := []byte(`{
		"id": "R0AAZKMD88M",
		"name": "",
		"media_server": "",
		"created_by": "U031L4VDD",
		"date_start": 1769466090,
		"date_end": 0,
		"participants": [],
		"participant_history": ["U031L4VDD"],
		"participants_events": {
			"U031L4VDD": {
				"user_team": {},
				"joined": true,
				"camera_on": false,
				"camera_off": false,
				"screenshare_on": false,
				"screenshare_off": false
			}
		},
		"participants_camera_on": [],
		"participants_camera_off": [],
		"participants_screenshare_on": [],
		"participants_screenshare_off": [],
		"canvas_thread_ts": "1769466090.342109",
		"thread_root_ts": "1769466090.342109",
		"channels": ["D0QJK3LDA"],
		"is_dm_call": true,
		"was_rejected": false,
		"was_missed": false,
		"was_accepted": false,
		"has_ended": false,
		"background_id": "GRADIENT_03",
		"canvas_background": "GRADIENT_03",
		"is_prewarmed": false,
		"is_scheduled": false,
		"recording": {
			"can_record_summary": "unavailable"
		},
		"locale": "en-US",
		"attached_file_ids": [],
		"media_backend_type": "free_willy",
		"display_id": "",
		"external_unique_id": "7069679a-3cb6-4622-a900-51b8d2ff2713",
		"app_id": "A00",
		"call_family": "huddle",
		"pending_invitees": {},
		"last_invite_status_by_user": {},
		"knocks": {},
		"huddle_link": "https://app.slack.com/huddle/T031L4VD9/D0QJK3LDA"
	}`)

	var room HuddleRoom
	err := json.Unmarshal(jsonData, &room)
	require.NoError(t, err)

	assert.Equal(t, "R0AAZKMD88M", room.ID)
	assert.Equal(t, "U031L4VDD", room.CreatedBy)
	assert.Equal(t, int64(1769466090), room.DateStart)
	assert.Equal(t, int64(0), room.DateEnd)
	assert.True(t, room.IsDMCall)
	assert.False(t, room.HasEnded)
	assert.Equal(t, "GRADIENT_03", room.BackgroundID)
	assert.Equal(t, "GRADIENT_03", room.CanvasBackground)
	assert.Equal(t, "free_willy", room.MediaBackendType)
	assert.Equal(t, "huddle", room.CallFamily)
	assert.Equal(t, "https://app.slack.com/huddle/T031L4VD9/D0QJK3LDA", room.HuddleLink)
	assert.Equal(t, "1769466090.342109", room.CanvasThreadTs)
	assert.Equal(t, "1769466090.342109", room.ThreadRootTs)
	assert.Equal(t, "en-US", room.Locale)

	require.Len(t, room.Channels, 1)
	assert.Equal(t, "D0QJK3LDA", room.Channels[0])

	require.Len(t, room.ParticipantHistory, 1)
	assert.Equal(t, "U031L4VDD", room.ParticipantHistory[0])

	require.NotNil(t, room.Recording)
	assert.Equal(t, "unavailable", room.Recording.CanRecordSummary)

	require.Contains(t, room.ParticipantsEvents, "U031L4VDD")
	pe := room.ParticipantsEvents["U031L4VDD"]
	assert.True(t, pe.Joined)
	assert.False(t, pe.CameraOn)
	assert.False(t, pe.ScreenshareOn)
}

func TestHuddleRoomWithActiveParticipants(t *testing.T) {
	jsonData := []byte(`{
		"id": "R123",
		"date_start": 1769466090,
		"date_end": 0,
		"participants": ["U001", "U002"],
		"participant_history": ["U001", "U002", "U003"],
		"participants_events": {
			"U001": {"joined": true, "camera_on": true, "camera_off": false, "screenshare_on": false, "screenshare_off": false},
			"U002": {"joined": true, "camera_on": false, "camera_off": false, "screenshare_on": true, "screenshare_off": false}
		},
		"participants_camera_on": ["U001"],
		"participants_screenshare_on": ["U002"],
		"is_dm_call": false,
		"was_rejected": false,
		"was_missed": false,
		"was_accepted": true,
		"has_ended": false,
		"call_family": "huddle"
	}`)

	var room HuddleRoom
	err := json.Unmarshal(jsonData, &room)
	require.NoError(t, err)

	assert.Equal(t, "R123", room.ID)
	require.Len(t, room.Participants, 2)
	assert.Equal(t, "U001", room.Participants[0])
	assert.Equal(t, "U002", room.Participants[1])

	require.Len(t, room.ParticipantHistory, 3)

	require.Len(t, room.ParticipantsCameraOn, 1)
	assert.Equal(t, "U001", room.ParticipantsCameraOn[0])

	require.Len(t, room.ParticipantsScreenshareOn, 1)
	assert.Equal(t, "U002", room.ParticipantsScreenshareOn[0])

	assert.True(t, room.WasAccepted)
	assert.False(t, room.IsDMCall)

	// Check participant events
	u1 := room.ParticipantsEvents["U001"]
	assert.True(t, u1.Joined)
	assert.True(t, u1.CameraOn)

	u2 := room.ParticipantsEvents["U002"]
	assert.True(t, u2.Joined)
	assert.True(t, u2.ScreenshareOn)
}

func TestHuddleRoomEndedState(t *testing.T) {
	jsonData := []byte(`{
		"id": "R456",
		"date_start": 1769454922,
		"date_end": 1769455026,
		"participants": [],
		"participant_history": ["U031L4VDD"],
		"is_dm_call": false,
		"was_rejected": false,
		"was_missed": false,
		"was_accepted": false,
		"has_ended": true,
		"call_family": "huddle"
	}`)

	var room HuddleRoom
	err := json.Unmarshal(jsonData, &room)
	require.NoError(t, err)

	assert.Equal(t, "R456", room.ID)
	assert.Equal(t, int64(1769454922), room.DateStart)
	assert.Equal(t, int64(1769455026), room.DateEnd)
	assert.True(t, room.HasEnded)
	assert.Empty(t, room.Participants)
	require.Len(t, room.ParticipantHistory, 1)
}
