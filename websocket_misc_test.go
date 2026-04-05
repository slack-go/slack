package slack

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSHRoomJoinEventUnmarshal(t *testing.T) {
	raw := `{
		"type": "sh_room_join",
		"room": {
			"id": "R01XXXBW",
			"name": null,
			"media_server": "",
			"created_by": "U12334",
			"date_start": 1607089008,
			"date_end": 0,
			"participants": ["U12334", "U56789"],
			"participant_history": ["U12334", "U56789"],
			"participants_camera_on": [],
			"participants_camera_off": [],
			"participants_screenshare_on": [],
			"participants_screenshare_off": [],
			"channels": ["C12334"],
			"is_dm_call": false,
			"was_rejected": false,
			"was_missed": false,
			"was_accepted": false,
			"has_ended": false,
			"media_backend_type": "free_willy",
			"external_unique_id": "8c92471f-test",
			"app_id": "A00"
		},
		"user": "U12334",
		"event_ts": "1607089059.080900",
		"ts": "1607089059.080900"
	}`

	var ev SHRoomJoinEvent
	err := json.Unmarshal([]byte(raw), &ev)
	require.NoError(t, err)

	assert.Equal(t, "sh_room_join", ev.Type)
	assert.Equal(t, "U12334", ev.User)
	assert.Equal(t, "R01XXXBW", ev.Room.ID)
	assert.Nil(t, ev.Room.Name)
	assert.Equal(t, "U12334", ev.Room.CreatedBy)
	assert.Equal(t, int64(1607089008), ev.Room.DateStart)
	assert.Equal(t, []string{"U12334", "U56789"}, ev.Room.Participants)
	assert.Equal(t, []string{"C12334"}, ev.Room.Channels)
	assert.False(t, ev.Room.IsDMCall)
	assert.Equal(t, "free_willy", ev.Room.MediaBackendType)
	assert.Equal(t, "1607089059.080900", ev.EventTS)
}

func TestSHRoomLeaveEventUnmarshal(t *testing.T) {
	raw := `{
		"type": "sh_room_leave",
		"room": {
			"id": "R01XXXBW",
			"name": null,
			"media_server": "",
			"created_by": "U12334",
			"date_start": 1607089008,
			"date_end": 0,
			"participants": ["U12334"],
			"participant_history": ["U12334", "U56789"],
			"participants_camera_on": [],
			"participants_camera_off": [],
			"participants_screenshare_on": [],
			"participants_screenshare_off": [],
			"channels": ["C12334"],
			"is_dm_call": false,
			"was_rejected": false,
			"was_missed": false,
			"was_accepted": false,
			"has_ended": false,
			"media_backend_type": "free_willy",
			"external_unique_id": "8c92471f-test",
			"app_id": "A00"
		},
		"user": "U56789",
		"event_ts": "1607091086.081500",
		"ts": "1607091086.081500"
	}`

	var ev SHRoomLeaveEvent
	err := json.Unmarshal([]byte(raw), &ev)
	require.NoError(t, err)

	assert.Equal(t, "sh_room_leave", ev.Type)
	assert.Equal(t, "U56789", ev.User)
	assert.Equal(t, "R01XXXBW", ev.Room.ID)
	assert.Equal(t, []string{"U12334"}, ev.Room.Participants)
	assert.Equal(t, "1607091086.081500", ev.EventTS)
}

func TestSHRoomUpdateEventUnmarshal(t *testing.T) {
	raw := `{
		"type": "sh_room_update",
		"room": {
			"id": "R0AQSG0Q859",
			"name": "A sort of topic",
			"media_server": "",
			"created_by": "U031L4VDD",
			"date_start": 1775402709,
			"date_end": 0,
			"participants": ["U031L4VDD"],
			"participant_history": ["U031L4VDD"],
			"participants_events": {"U031L4VDD": {"joined": true, "camera_on": false}},
			"participants_camera_on": [],
			"participants_camera_off": [],
			"participants_screenshare_on": [],
			"participants_screenshare_off": [],
			"canvas_thread_ts": "1775402709.576349",
			"thread_root_ts": "1775402709.576349",
			"channels": ["C031L4VDP"],
			"is_dm_call": false,
			"was_rejected": false,
			"was_missed": false,
			"was_accepted": false,
			"has_ended": false,
			"background_id": "GRADIENT_02",
			"canvas_background": "GRADIENT_02",
			"is_prewarmed": true,
			"is_scheduled": false,
			"recording": {"can_record_summary": "unavailable"},
			"locale": "en-US",
			"attached_file_ids": [],
			"media_backend_type": "free_willy",
			"display_id": "",
			"external_unique_id": "755e016f-aae1-4d4f-abcc-952b1b872713",
			"app_id": "A00",
			"call_family": "huddle",
			"huddle_link": "https://app.slack.com/huddle/T031L4VD9/C031L4VDP"
		},
		"user": "U031L4VDD",
		"huddle": {"channel_id": "C031L4VDP"},
		"event_ts": "1775402785.000200",
		"ts": "1775402785.000200"
	}`

	var ev SHRoomUpdateEvent
	err := json.Unmarshal([]byte(raw), &ev)
	require.NoError(t, err)

	assert.Equal(t, "sh_room_update", ev.Type)
	assert.Equal(t, "U031L4VDD", ev.User)
	assert.Equal(t, "R0AQSG0Q859", ev.Room.ID)
	assert.NotNil(t, ev.Room.Name)
	assert.Equal(t, "A sort of topic", *ev.Room.Name)
	assert.Equal(t, "huddle", ev.Room.CallFamily)
	assert.True(t, ev.Room.IsPrewarmed)
	assert.Equal(t, "1775402709.576349", ev.Room.CanvasThreadTS)
	assert.Equal(t, "GRADIENT_02", ev.Room.BackgroundID)
	assert.Equal(t, "en-US", ev.Room.Locale)
	assert.NotNil(t, ev.Room.Recording)
	assert.Equal(t, "unavailable", ev.Room.Recording.CanRecordSummary)
	assert.Equal(t, "https://app.slack.com/huddle/T031L4VD9/C031L4VDP", ev.Room.HuddleLink)
	assert.NotNil(t, ev.Huddle)
	assert.Equal(t, "C031L4VDP", ev.Huddle.ChannelID)
	assert.NotNil(t, ev.Room.ParticipantsEvents)
	assert.Contains(t, ev.Room.ParticipantsEvents, "U031L4VDD")
}
