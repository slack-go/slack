package slack

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFileImageMetadataRoundTrip(t *testing.T) {
	// Shape and integer dimensions captured from files.info after image processing.
	const payload = `{"alt_txt":"blue left half and orange right half","user_team":"E123","file_access":"visible","has_more_shares":false,"thumb_800":"https://files.slack.com/example_800.png","thumb_800_w":800,"thumb_800_h":600,"thumb_tiny":"AwAkADCKiiiukxCr"}`
	var file File
	require.NoError(t, json.Unmarshal([]byte(payload), &file))
	encoded, err := json.Marshal(file)
	require.NoError(t, err)
	var want, got map[string]json.RawMessage
	require.NoError(t, json.Unmarshal([]byte(payload), &want))
	require.NoError(t, json.Unmarshal(encoded, &got))
	for key, value := range want {
		t.Run(key, func(t *testing.T) {
			require.Contains(t, got, key)
			require.JSONEq(t, string(value), string(got[key]))
		})
	}
}

func TestRemoteFileMetadataRoundTrip(t *testing.T) {
	// Captured from files.info and files.remote.info for a newly shared remote file.
	const payload = `{"external_id":"slack-go-file-fields-test","external_url":"https://docs.slack.dev/reference/objects/file-object/","app_id":"A123","app_provides_file_work_objects":true,"user_team":"E123","file_access":"visible","has_more_shares":false}`
	for _, tc := range []struct {
		name  string
		value any
	}{
		{name: "File", value: &File{}},
		{name: "RemoteFile", value: &RemoteFile{}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			require.NoError(t, json.Unmarshal([]byte(payload), tc.value))
			encoded, err := json.Marshal(tc.value)
			require.NoError(t, err)
			var want, got map[string]json.RawMessage
			require.NoError(t, json.Unmarshal([]byte(payload), &want))
			require.NoError(t, json.Unmarshal(encoded, &got))
			for key, value := range want {
				t.Run(key, func(t *testing.T) {
					actual, exists := got[key]
					require.True(t, exists, "missing %s", key)
					require.JSONEq(t, string(value), string(actual))
				})
			}
		})
	}
}

func TestFileCanvasMetadataRoundTrip(t *testing.T) {
	// files.info returned fractional read time, millisecond last_read, and empty editors.
	const payload = `{
		"url_static_preview":"https://example.slack.com/docs/T123/F123/mobile",
		"quip_thread_id":"KLE9AAS5goj",
		"updated":1792482756,
		"update_notification":1792482756,
		"canvas_readtime":0.17181854676837,
		"access":"owner",
		"private_channels_with_file_access_count":0,
		"dm_mpdm_users_with_file_access":[{"user_id":"U123","access":"write"}],
		"teams_shared_with":["E123"],
		"is_restricted_sharing_enabled":false,
		"excluded_from_slack_ai":false,
		"org_or_workspace_access":"none",
		"title_blocks":[{"type":"rich_text","block_id":"temp:canvas-title","elements":[{"type":"rich_text_section","elements":[{"type":"text","text":"File fields: canvas metadata"}]}]}],
		"canvas_creator_id":"U123",
		"team_pref_version_history_enabled":true,
		"canvas_printing_enabled":true,
		"is_ai_suggested":false,
		"editors":[],
		"last_read":1792482757209,
		"show_badge":false
	}`
	var file File
	require.NoError(t, json.Unmarshal([]byte(payload), &file))
	require.NotNil(t, file.TitleBlocks)
	require.Len(t, file.TitleBlocks.BlockSet, 1)
	require.IsType(t, &RichTextBlock{}, file.TitleBlocks.BlockSet[0])
	require.Equal(t, []FileUserAccess{{UserID: "U123", Access: "write"}}, file.DMMPDMUsersWithFileAccess)
	encoded, err := json.Marshal(file)
	require.NoError(t, err)
	var want, got map[string]json.RawMessage
	require.NoError(t, json.Unmarshal([]byte(payload), &want))
	require.NoError(t, json.Unmarshal(encoded, &got))
	for key, value := range want {
		t.Run(key, func(t *testing.T) {
			actual, exists := got[key]
			require.True(t, exists, "missing %s", key)
			require.JSONEq(t, string(value), string(actual))
		})
	}
}

func TestFileListMetadataPreservation(t *testing.T) {
	const payload = `{
	  "list_metadata": {
	    "schema": [
	      {
	        "id": "ColTask",
	        "name": "Task",
	        "key": "task",
	        "type": "text",
	        "is_primary_column": true
	      },
	      {
	        "id": "ColDone",
	        "name": "Done",
	        "key": "done",
	        "type": "checkbox",
	        "is_primary_column": false
	      }
	    ],
	    "views": [
	      {
	        "id": "ViewRecord",
	        "name": "Record",
	        "type": "record",
	        "is_locked": false,
	        "position": "1792485941",
	        "columns": [
	          {
	            "visible": true,
	            "key": "task",
	            "id": "ColTask",
	            "position": "5000000000"
	          },
	          {
	            "visible": true,
	            "key": "done",
	            "id": "ColDone",
	            "position": "5000000001"
	          }
	        ],
	        "date_created": 1792485941,
	        "created_by": "U123",
	        "stick_column_left": false
	      },
	      {
	        "id": "ViewAll",
	        "name": "All items",
	        "type": "table",
	        "is_locked": false,
	        "position": "1792485942",
	        "columns": [
	          {
	            "visible": true,
	            "key": "task",
	            "id": "ColTask",
	            "position": "5000000000"
	          },
	          {
	            "visible": true,
	            "key": "done",
	            "id": "ColDone",
	            "position": "5000000001"
	          }
	        ],
	        "date_created": 1792485942,
	        "created_by": "U123",
	        "stick_column_left": false,
	        "is_all_items_view": true,
	        "default_view_key": "all_items",
	        "show_completed_items": true
	      }
	    ],
	    "integrations": [],
	    "icon": "",
	    "description": "",
	    "description_blocks": [],
	    "is_trial": false,
	    "subtask_schema": [
	      {
	        "id": "ColTask",
	        "name": "Task",
	        "key": "task",
	        "type": "text",
	        "is_primary_column": true
	      }
	    ],
	    "creation_source": {
	      "type": "blank"
	    },
	    "default_view": ""
	  },
	  "list_limits": {
	    "over_row_maximum": false,
	    "row_count_limit": 5000,
	    "row_count": 1,
	    "archived_row_count": 0,
	    "over_column_maximum": false,
	    "column_count": 2,
	    "column_count_limit": 30,
	    "over_view_maximum": false,
	    "view_count": 2,
	    "view_count_limit": 50,
	    "max_attachments_per_cell": 10
	  },
	  "list_csv_download_url": "https://files.slack.com/example/list.csv",
	  "last_editor": "U123",
	  "edit_timestamp": 1792485942
	}`
	var expected map[string]json.RawMessage
	require.NoError(t, json.Unmarshal([]byte(payload), &expected))
	var metadata map[string]any
	require.NoError(t, json.Unmarshal(expected["list_metadata"], &metadata))
	// Existing List types marshal these defaults even when Slack omits them.
	metadata["todo_mode"] = false
	metadata["creation_source"].(map[string]any)["reference_id"] = ""
	metadata["views"].([]any)[0].(map[string]any)["is_all_items_view"] = false
	encodedMetadata, err := json.Marshal(metadata)
	expected["list_metadata"] = encodedMetadata
	require.NoError(t, err)
	for _, tc := range []struct {
		name  string
		value any
	}{
		{name: "File", value: &File{}},
		{name: "List", value: &List{}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			require.NoError(t, json.Unmarshal([]byte(payload), tc.value))
			encoded, err := json.Marshal(tc.value)
			require.NoError(t, err)
			var got map[string]json.RawMessage
			require.NoError(t, json.Unmarshal(encoded, &got))
			for key, value := range expected {
				t.Run(key, func(t *testing.T) {
					actual, exists := got[key]
					require.True(t, exists, "missing %s", key)
					require.JSONEq(t, string(value), string(actual))
				})
			}
		})
	}
}

func TestFileGIFMetadataRoundTrip(t *testing.T) {
	const payload = `{"thumb_480_gif":"https://files.slack.com/example_480.gif","deanimate":"https://files.slack.com/example_deanimate.png"}`
	var file File
	require.NoError(t, json.Unmarshal([]byte(payload), &file))
	encoded, err := json.Marshal(file)
	require.NoError(t, err)
	var want, got map[string]json.RawMessage
	require.NoError(t, json.Unmarshal([]byte(payload), &want))
	require.NoError(t, json.Unmarshal(encoded, &got))
	for key, value := range want {
		t.Run(key, func(t *testing.T) {
			actual, exists := got[key]
			require.True(t, exists, "missing %s", key)
			require.JSONEq(t, string(value), string(actual))
		})
	}
}

func TestFileAudioMetadataRoundTrip(t *testing.T) {
	const payload = `{
	  "aac": "https://files.slack.com/example_audio.mp4",
	  "audio_wave_samples": [85, 93, 100, 98, 89, 86, 93, 100, 98, 89, 86, 93, 100, 98, 89, 86, 93, 100, 98, 89, 86, 93, 100, 98, 89, 86, 93, 100, 98, 89, 86, 93, 100, 98, 89, 86, 93, 100, 98, 89, 86, 93, 100, 98, 89, 86, 93, 100, 98, 89, 86, 93, 100, 98, 89, 86, 93, 100, 98, 89, 86, 93, 100, 98, 89, 86, 93, 100, 98, 89, 86, 93, 100, 98, 89, 86, 93, 100, 98, 89, 86, 93, 100, 98, 89, 86, 93, 100, 98, 89, 86, 93, 100, 98, 89, 86, 93, 100, 98, 75],
	  "is_transcription_region_supported": true,
	  "transcription": {"status": "none"}
	}`
	var file File
	require.NoError(t, json.Unmarshal([]byte(payload), &file))
	encoded, err := json.Marshal(file)
	require.NoError(t, err)
	var want, got map[string]json.RawMessage
	require.NoError(t, json.Unmarshal([]byte(payload), &want))
	require.NoError(t, json.Unmarshal(encoded, &got))
	for key, value := range want {
		t.Run(key, func(t *testing.T) {
			actual, exists := got[key]
			require.True(t, exists, "missing %s", key)
			require.JSONEq(t, string(value), string(actual))
		})
	}
}

func TestFilePreviewTruncation(t *testing.T) {
	// Slack returned true for the long snippet and false for the uploaded .eml.
	for _, payload := range []string{`{"preview_is_truncated":true}`, `{"preview_is_truncated":false}`} {
		var file File
		require.NoError(t, json.Unmarshal([]byte(payload), &file))
		encoded, err := json.Marshal(file)
		require.NoError(t, err)
		var want, got map[string]json.RawMessage
		require.NoError(t, json.Unmarshal([]byte(payload), &want))
		require.NoError(t, json.Unmarshal(encoded, &got))
		actual, exists := got["preview_is_truncated"]
		require.True(t, exists, "missing preview_is_truncated")
		require.JSONEq(t, string(want["preview_is_truncated"]), string(actual))
	}
}

func TestFilePDFMetadataRoundTrip(t *testing.T) {
	const payload = `{"thumb_pdf":"https://files.slack.com/example_thumb_pdf.png","thumb_pdf_h":1286,"thumb_pdf_w":909}`
	var file File
	require.NoError(t, json.Unmarshal([]byte(payload), &file))
	encoded, err := json.Marshal(file)
	require.NoError(t, err)
	var want, got map[string]json.RawMessage
	require.NoError(t, json.Unmarshal([]byte(payload), &want))
	require.NoError(t, json.Unmarshal(encoded, &got))
	for key, value := range want {
		t.Run(key, func(t *testing.T) {
			actual, exists := got[key]
			require.True(t, exists, "missing %s", key)
			require.JSONEq(t, string(value), string(actual))
		})
	}
}
