package slackevents

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAppMention(t *testing.T) {
	rawE := []byte(`
			{
				"type": "app_mention",
				"user": "U061F7AUR",
				"text": "<@U0LAN0Z89> is it everything a river should be?",
				"ts": "1515449522.000016",
				"thread_ts": "1515449522.000016",
				"channel": "C0LAN2Q65",
				"event_ts": "1515449522000016",
				"source_team": "T3MQV36V7",
				"user_team": "T3MQV36V7",
				"blah": "test"
		}
	`)
	err := json.Unmarshal(rawE, &AppMentionEvent{})
	if err != nil {
		t.Error(err)
	}
}

func TestAppUninstalled(t *testing.T) {
	rawE := []byte(`
		{
			"type": "app_uninstalled"
		}
	`)
	err := json.Unmarshal(rawE, &AppUninstalledEvent{})
	if err != nil {
		t.Error(err)
	}
}

func TestFileChangeEvent(t *testing.T) {
	rawE := []byte(`
		{
			"type": "file_change",
			"file_id": "F1234567890",
			"file": {
				"id": "F1234567890"
			}
		}
	`)

	var e FileChangeEvent
	if err := json.Unmarshal(rawE, &e); err != nil {
		t.Fatal(err)
	}
	if e.Type != "file_change" {
		t.Errorf("type should be file_change, was %s", e.Type)
	}
	if e.FileID != "F1234567890" {
		t.Errorf("file ID should be F1234567890, was %s", e.FileID)
	}
	if e.File.ID != "F1234567890" {
		t.Errorf("file.id should be F1234567890, was %s", e.File.ID)
	}
}

func TestFileDeletedEvent(t *testing.T) {
	rawE := []byte(`
		{
			"type": "file_deleted",
			"file_id": "F1234567890",
			"event_ts": "1234567890.123456"
		}
	`)

	var e FileDeletedEvent
	if err := json.Unmarshal(rawE, &e); err != nil {
		t.Fatal(err)
	}
	if e.Type != "file_deleted" {
		t.Errorf("type should be file_deleted, was %s", e.Type)
	}
	if e.FileID != "F1234567890" {
		t.Errorf("file ID should be F1234567890, was %s", e.FileID)
	}
	if e.EventTimestamp != "1234567890.123456" {
		t.Errorf("event timestamp should be 1234567890.123456, was %s", e.EventTimestamp)
	}
}

func TestFileSharedEvent(t *testing.T) {
	rawE := []byte(`
		{
			"type": "file_shared",
			"channel_id": "C1234567890",
			"file_id": "F1234567890",
			"user_id": "U11235813",
			"file": {
				"id": "F1234567890"
			},
			"event_ts": "1234567890.123456"
		}
	`)

	var e FileSharedEvent
	if err := json.Unmarshal(rawE, &e); err != nil {
		t.Fatal(err)
	}
	if e.Type != "file_shared" {
		t.Errorf("type should be file_shared, was %s", e.Type)
	}
	if e.ChannelID != "C1234567890" {
		t.Errorf("channel ID should be C1234567890, was %s", e.ChannelID)
	}
	if e.FileID != "F1234567890" {
		t.Errorf("file ID should be F1234567890, was %s", e.FileID)
	}
	if e.UserID != "U11235813" {
		t.Errorf("user ID should be U11235813, was %s", e.UserID)
	}
	if e.File.ID != "F1234567890" {
		t.Errorf("file.id should be F1234567890, was %s", e.File.ID)
	}
	if e.EventTimestamp != "1234567890.123456" {
		t.Errorf("event timestamp should be 1234567890.123456, was %s", e.EventTimestamp)
	}
}

func TestFileUnsharedEvent(t *testing.T) {
	rawE := []byte(`
		{
			"type": "file_unshared",
			"file_id": "F1234567890",
			"file": {
				"id": "F1234567890"
			}
		}
	`)

	var e FileUnsharedEvent
	if err := json.Unmarshal(rawE, &e); err != nil {
		t.Fatal(err)
	}
	if e.Type != "file_unshared" {
		t.Errorf("type should be file_shared, was %s", e.Type)
	}
	if e.FileID != "F1234567890" {
		t.Errorf("file ID should be F1234567890, was %s", e.FileID)
	}
	if e.File.ID != "F1234567890" {
		t.Errorf("file.id should be F1234567890, was %s", e.File.ID)
	}
}

func TestGridMigrationFinishedEvent(t *testing.T) {
	rawE := []byte(`
			{
				"type": "grid_migration_finished",
				"enterprise_id": "EXXXXXXXX"
			}
	`)
	err := json.Unmarshal(rawE, &GridMigrationFinishedEvent{})
	if err != nil {
		t.Error(err)
	}
}

func TestGridMigrationStartedEvent(t *testing.T) {
	rawE := []byte(`
			{
				"token": "XXYYZZ",
				"team_id": "TXXXXXXXX",
				"api_app_id": "AXXXXXXXXX",
				"event": {
						"type": "grid_migration_started",
						"enterprise_id": "EXXXXXXXX"
				},
				"type": "event_callback",
				"event_id": "EvXXXXXXXX",
				"event_time": 1234567890
		}
	`)
	err := json.Unmarshal(rawE, &GridMigrationStartedEvent{})
	if err != nil {
		t.Error(err)
	}
}

func TestLinkSharedEvent(t *testing.T) {
	rawE := []byte(`
			{
				"type": "link_shared",
				"channel": "Cxxxxxx",
				"user": "Uxxxxxxx",
				"message_ts": "123456789.9875",
				"thread_ts": "123456789.9876",
				"links":
						[
								{
										"domain": "example.com",
										"url": "https://example.com/12345"
								},
								{
										"domain": "example.com",
										"url": "https://example.com/67890"
								},
								{
										"domain": "another-example.com",
										"url": "https://yet.another-example.com/v/abcde"
								}
						]
		}
	`)
	err := json.Unmarshal(rawE, &LinkSharedEvent{})
	if err != nil {
		t.Error(err)
	}
}

func TestLinkSharedEvent_struct(t *testing.T) {
	e := LinkSharedEvent{
		Type:             "link_shared",
		User:             "Uxxxxxxx",
		TimeStamp:        "123456789.9876",
		Channel:          "Cxxxxxx",
		MessageTimeStamp: "123456789.9875",
		ThreadTimeStamp:  "123456789.9876",
		Links: []SharedLinks{
			{Domain: "example.com", URL: "https://example.com/12345"},
			{Domain: "example.com", URL: "https://example.com/67890"},
			{Domain: "another-example.com", URL: "https://yet.another-example.com/v/abcde"},
		},
		EventTimestamp: "123456789.9876",
	}
	rawE, err := json.Marshal(e)
	if err != nil {
		t.Error(err)
	}
	expected := `{"type":"link_shared","user":"Uxxxxxxx","ts":"123456789.9876","channel":"Cxxxxxx",` +
		`"message_ts":"123456789.9875","thread_ts":"123456789.9876","links":[{"domain":"example.com",` +
		`"url":"https://example.com/12345"},{"domain":"example.com","url":"https://example.com/67890"},` +
		`{"domain":"another-example.com","url":"https://yet.another-example.com/v/abcde"}],"event_ts":"123456789.9876"}`
	if string(rawE) != expected {
		t.Errorf("expected %s, but got %s", expected, string(rawE))
	}
}

func TestLinkSharedComposerEvent(t *testing.T) {
	rawE := []byte(`
			{
				"type": "link_shared",
				"channel": "COMPOSER",
				"is_bot_user_member": true,
				"user": "Uxxxxxxx",
				"message_ts": "Uxxxxxxx-909b5454-75f8-4ac4-b325-1b40e230bbd8-gryl3kb80b3wm49ihzoo35fyqoq08n2y",
				"unfurl_id": "Uxxxxxxx-909b5454-75f8-4ac4-b325-1b40e230bbd8-gryl3kb80b3wm49ihzoo35fyqoq08n2y",
				"source": "composer",
				"links": [
					{
						"domain": "example.com",
						"url": "https://example.com/12345"
					},
					{
						"domain": "example.com",
						"url": "https://example.com/67890"
					},
					{
						"domain": "another-example.com",
						"url": "https://yet.another-example.com/v/abcde"
					}
				]
			}
	`)
	err := json.Unmarshal(rawE, &LinkSharedEvent{})
	if err != nil {
		t.Error(err)
	}
}

func TestMessageEvent(t *testing.T) {
	rawE := []byte(`
			{
				"client_msg_id": "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee",
				"type": "message",
				"channel": "G024BE91L",
				"user": "U2147483697",
				"text": "Live long and prospect.",
				"ts": "1355517523.000005",
				"event_ts": "1355517523.000005",
				"channel_type": "channel",
				"source_team": "T3MQV36V7",
				"user_team": "T3MQV36V7",
				"message": {
					"text": "To infinity and beyond.",
					"edited": {
						"user": "U2147483697",
						"ts": "1355517524.000000"
					}
				},
				"previous_message": {
					"text": "Live long and prospect."
				}
		}
	`)
	err := json.Unmarshal(rawE, &MessageEvent{})
	if err != nil {
		t.Error(err)
	}
}

func TestBotMessageEvent(t *testing.T) {
	rawE := []byte(`
			{
				"type": "message",
				"subtype": "bot_message",
				"ts": "1358877455.000010",
				"text": "Pushing is the answer",
				"bot_id": "BB12033",
				"username": "github",
				"icons": {}
		}
	`)
	err := json.Unmarshal(rawE, &MessageEvent{})
	if err != nil {
		t.Error(err)
	}
}

func TestThreadBroadcastEvent(t *testing.T) {
	rawE := []byte(`
			{
				"type": "message",
				"subtype": "thread_broadcast",
				"channel": "G024BE91L",
				"user": "U2147483697",
				"text": "Live long and prospect.",
				"ts": "1355517523.000005",
				"event_ts": "1355517523.000005",
				"channel_type": "channel",
				"source_team": "T3MQV36V7",
				"user_team": "T3MQV36V7",
				"message": {
					"text": "To infinity and beyond.",
					"root": {
						"text": "To infinity and beyond.",
						"ts": "1355517523.000005"
					},
					"edited": {
						"user": "U2147483697",
						"ts": "1355517524.000000"
					}
				},
				"previous_message": {
					"text": "Live long and prospect."
				}
		}
	`)

	var me MessageEvent
	if err := json.Unmarshal(rawE, &me); err != nil {
		t.Error(err)
	}

	if me.Root != nil {
		t.Error("me.Root should be nil")
	}

	if me.Message.Root == nil {
		t.Fatal("me.Message.Root is nil")
	}

	if me.Message.Root.TimeStamp != "1355517523.000005" {
		t.Errorf("me.Message.Root.TimeStamp = %q, want %q", me.Root.TimeStamp, "1355517523.000005")
	}
}

func TestMemberJoinedChannelEvent(t *testing.T) {
	rawE := []byte(`
		{
			"type": "member_joined_channel",
			"user": "W06GH7XHN",
			"channel": "C0698JE0H",
			"channel_type": "C",
			"team": "T024BE7LD",
			"inviter": "U123456789"
		}
	`)
	evt := MemberJoinedChannelEvent{}
	err := json.Unmarshal(rawE, &evt)
	if err != nil {
		t.Error(err)
	}

	expected := MemberJoinedChannelEvent{
		Type:        "member_joined_channel",
		User:        "W06GH7XHN",
		Channel:     "C0698JE0H",
		ChannelType: "C",
		Team:        "T024BE7LD",
		Inviter:     "U123456789",
	}

	assert.Equal(t, expected, evt)
}

func TestMemberLeftChannelEvent(t *testing.T) {
	rawE := []byte(`
			{
				"type": "member_left_channel",
				"user": "W06GH7XHN",
				"channel": "C0698JE0H",
				"channel_type": "C",
				"team": "T024BE7LD"
		}
	`)
	err := json.Unmarshal(rawE, &MemberLeftChannelEvent{})
	if err != nil {
		t.Error(err)
	}
}

func TestPinAdded(t *testing.T) {
	rawE := []byte(`
			{
				"type": "pin_added",
				"user": "U061F7AUR",
				"item": {
					"type": "message",
					"channel":"C0LAN2Q65",
					"message":{
						"type":"message",
						"user":"U061F7AUR",
						"text": "<@U0LAN0Z89> is it everything a river should be?",
						"ts":"1539904112.000100",
						"pinned_to":["C0LAN2Q65"],
						"replace_original":false,
						"delete_original":false
					}
				},
				"channel_id":"C0LAN2Q65",
				"event_ts": "1515449522000016"
		}
	`)
	err := json.Unmarshal(rawE, &PinAddedEvent{})
	if err != nil {
		t.Error(err)
	}
}

func TestPinRemoved(t *testing.T) {
	rawE := []byte(`
			{
				"type": "pin_removed",
				"user": "U061F7AUR",
				"item": {
					"type": "message",
					"channel":"C0LAN2Q65",
					"message":{
						"type":"message",
						"user":"U061F7AUR",
						"text": "<@U0LAN0Z89> is it everything a river should be?",
						"ts":"1539904112.000100",
						"pinned_to":["C0LAN2Q65"],
						"replace_original":false,
						"delete_original":false
					}
				},
				"channel_id":"C0LAN2Q65",
				"event_ts": "1515449522000016"
		}
	`)
	err := json.Unmarshal(rawE, &PinRemovedEvent{})
	if err != nil {
		t.Error(err)
	}
}

func TestTokensRevoked(t *testing.T) {
	rawE := []byte(`
	{
		"type": "tokens_revoked",
		"tokens": {
				"oauth": [
						"OUXXXXXXXX"
				],
				"bot": [
						"BUXXXXXXXX"
				]
		}
	}
`)
	tre := TokensRevokedEvent{}
	err := json.Unmarshal(rawE, &tre)
	if err != nil {
		t.Error(err)
	}

	if tre.Type != "tokens_revoked" {
		t.Fail()
	}

	if len(tre.Tokens.Bot) != 1 || tre.Tokens.Bot[0] != "BUXXXXXXXX" {
		t.Fail()
	}

	if len(tre.Tokens.Oauth) != 1 || tre.Tokens.Oauth[0] != "OUXXXXXXXX" {
		t.Fail()
	}
}

func TestEmojiChanged(t *testing.T) {
	var (
		ece EmojiChangedEvent
		err error
	)

	// custom emoji added event
	rawAddE := []byte(`
	{
		"type": "emoji_changed",
		"subtype": "add",
		"name": "picard_facepalm",
		"value": "https://my.slack.com/emoji/picard_facepalm/db8e287430eaa459.gif",
		"event_ts" : "1361482916.000004"
	}
`)
	ece = EmojiChangedEvent{}
	err = json.Unmarshal(rawAddE, &ece)
	if err != nil {
		t.Error(err)
	}
	if ece.Subtype != "add" {
		t.Fail()
	}
	if ece.Name != "picard_facepalm" {
		t.Fail()
	}

	// emoji removed event
	rawRemoveE := []byte(`
	{
		"type": "emoji_changed",
		"subtype": "remove",
		"names": ["picard_facepalm"],
		"event_ts" : "1361482916.000004"
	}
`)
	ece = EmojiChangedEvent{}
	err = json.Unmarshal(rawRemoveE, &ece)
	if err != nil {
		t.Error(err)
	}
	if ece.Subtype != "remove" {
		t.Fail()
	}
	if len(ece.Names) != 1 {
		t.Fail()
	}
	if ece.Names[0] != "picard_facepalm" {
		t.Fail()
	}

	// custom emoji rename event
	rawRenameE := []byte(`
	{
		"type": "emoji_changed",
		"subtype": "rename",
		"old_name": "grin",
		"new_name": "cheese-grin",
		"value": "https://my.slack.com/emoji/picard_facepalm/db8e287430eaa459.gif",
		"event_ts" : "1361482916.000004"
	}
`)
	ece = EmojiChangedEvent{}
	err = json.Unmarshal(rawRenameE, &ece)
	if err != nil {
		t.Error(err)
	}
	if ece.Subtype != "rename" {
		t.Fail()
	}
	if ece.OldName != "grin" {
		t.Fail()
	}
	if ece.NewName != "cheese-grin" {
		t.Fail()
	}
}

func TestWorkflowStepExecute(t *testing.T) {
	// see: https://api.slack.com/events/workflow_step_execute
	rawE := []byte(`
	{
		"type":"workflow_step_execute",
		"callback_id":"open_ticket",
		"workflow_step":{
			"workflow_step_execute_id":"1036669284371.19077474947.c94bcf942e047298d21f89faf24f1326",
			"workflow_id":"123456789012345678",
			"workflow_instance_id":"987654321098765432",
			"step_id":"12a345bc-1a23-4567-8b90-1234a567b8c9",
			"inputs":{
				"example-select-input":{
					"value": "value-two",
					"skip_variable_replacement": false
				}
			},
			"outputs":[
			]
		},
		"event_ts":"1643290847.766536"
	}
	`)

	wse := WorkflowStepExecuteEvent{}
	err := json.Unmarshal(rawE, &wse)
	if err != nil {
		t.Error(err)
	}

	if wse.Type != "workflow_step_execute" {
		t.Fail()
	}
	if wse.CallbackID != "open_ticket" {
		t.Fail()
	}
	if wse.WorkflowStep.WorkflowStepExecuteID != "1036669284371.19077474947.c94bcf942e047298d21f89faf24f1326" {
		t.Fail()
	}
	if wse.WorkflowStep.WorkflowID != "123456789012345678" {
		t.Fail()
	}
	if wse.WorkflowStep.WorkflowInstanceID != "987654321098765432" {
		t.Fail()
	}
	if wse.WorkflowStep.StepID != "12a345bc-1a23-4567-8b90-1234a567b8c9" {
		t.Fail()
	}
	if len(*wse.WorkflowStep.Inputs) == 0 {
		t.Fail()
	}
	if inputElement, ok := (*wse.WorkflowStep.Inputs)["example-select-input"]; ok {
		if inputElement.Value != "value-two" {
			t.Fail()
		}
		if inputElement.SkipVariableReplacement != false {
			t.Fail()
		}
	}
}

func TestMessageMetadataPosted(t *testing.T) {
	rawE := []byte(`
	{
		"type":"message_metadata_posted",
		"app_id":"APPXXX",
		"bot_id":"BOTXXX",	
		"user_id":"USERXXX",	
		"team_id":"TEAMXXX",	
		"channel_id":"CHANNELXXX",	
		"metadata":{
			"event_type":"type",
			"event_payload":{"key": "value"}
		},
		"message_ts":"1660398079.756349",
		"event_ts":"1660398079.756349"
	}
	`)

	mmp := MessageMetadataPostedEvent{}
	err := json.Unmarshal(rawE, &mmp)
	if err != nil {
		t.Error(err)
	}

	if mmp.Type != "message_metadata_posted" {
		t.Fail()
	}
	if mmp.AppId != "APPXXX" {
		t.Fail()
	}
	if mmp.BotId != "BOTXXX" {
		t.Fail()
	}
	if mmp.UserId != "USERXXX" {
		t.Fail()
	}
	if mmp.TeamId != "TEAMXXX" {
		t.Fail()
	}
	if mmp.ChannelId != "CHANNELXXX" {
		t.Fail()
	}
	if mmp.Metadata.EventType != "type" {
		t.Fail()
	}
	payload := mmp.Metadata.EventPayload
	if len(payload) <= 0 {
		t.Fail()
	}
	if mmp.EventTimestamp != "1660398079.756349" {
		t.Fail()
	}
	if mmp.MessageTimestamp != "1660398079.756349" {
		t.Fail()
	}
}

func TestMessageMetadataUpdated(t *testing.T) {
	rawE := []byte(`
	{
		"type":"message_metadata_updated",
		"channel_id":"CHANNELXXX",	
		"event_ts":"1660398079.756349",
		"previous_metadata":{
			"event_type":"type1",
			"event_payload":{"key1": "value1"}
		},
		"app_id":"APPXXX",
		"bot_id":"BOTXXX",	
		"user_id":"USERXXX",	
		"team_id":"TEAMXXX",	
		"message_ts":"1660398079.756349",
		"metadata":{
			"event_type":"type2",
			"event_payload":{"key2": "value2"}
		}
	}
	`)

	mmp := MessageMetadataUpdatedEvent{}
	err := json.Unmarshal(rawE, &mmp)
	if err != nil {
		t.Error(err)
	}

	if mmp.Type != "message_metadata_updated" {
		t.Fail()
	}
	if mmp.ChannelId != "CHANNELXXX" {
		t.Fail()
	}
	if mmp.EventTimestamp != "1660398079.756349" {
		t.Fail()
	}
	if mmp.PreviousMetadata.EventType != "type1" {
		t.Fail()
	}
	payload := mmp.PreviousMetadata.EventPayload
	if len(payload) <= 0 {
		t.Fail()
	}
	if mmp.AppId != "APPXXX" {
		t.Fail()
	}
	if mmp.BotId != "BOTXXX" {
		t.Fail()
	}
	if mmp.UserId != "USERXXX" {
		t.Fail()
	}
	if mmp.TeamId != "TEAMXXX" {
		t.Fail()
	}
	if mmp.MessageTimestamp != "1660398079.756349" {
		t.Fail()
	}
	if mmp.Metadata.EventType != "type2" {
		t.Fail()
	}
	payload = mmp.Metadata.EventPayload
	if len(payload) <= 0 {
		t.Fail()
	}
}

func TestMessageMetadataDeleted(t *testing.T) {
	rawE := []byte(`
	{
		"type":"message_metadata_deleted",
		"channel_id":"CHANNELXXX",	
		"event_ts":"1660398079.756349",
		"previous_metadata":{
			"event_type":"type",
			"event_payload":{"key": "value"}
		},
		"app_id":"APPXXX",
		"bot_id":"BOTXXX",	
		"user_id":"USERXXX",	
		"team_id":"TEAMXXX",	
		"message_ts":"1660398079.756349",
		"deleted_ts":"1660398079.756349"
	}
	`)

	mmp := MessageMetadataDeletedEvent{}
	err := json.Unmarshal(rawE, &mmp)
	if err != nil {
		t.Error(err)
	}

	if mmp.Type != "message_metadata_deleted" {
		t.Fail()
	}
	if mmp.ChannelId != "CHANNELXXX" {
		t.Fail()
	}
	if mmp.EventTimestamp != "1660398079.756349" {
		t.Fail()
	}
	if mmp.PreviousMetadata.EventType != "type" {
		t.Fail()
	}
	payload := mmp.PreviousMetadata.EventPayload
	if len(payload) <= 0 {
		t.Fail()
	}
	if mmp.AppId != "APPXXX" {
		t.Fail()
	}
	if mmp.BotId != "BOTXXX" {
		t.Fail()
	}
	if mmp.UserId != "USERXXX" {
		t.Fail()
	}
	if mmp.TeamId != "TEAMXXX" {
		t.Fail()
	}
	if mmp.MessageTimestamp != "1660398079.756349" {
		t.Fail()
	}
	if mmp.DeletedTimestamp != "1660398079.756349" {
		t.Fail()
	}
}

func TestUserProfileChanged(t *testing.T) {
	rawE := []byte(`
	{
		"token": "whatever",
		"team_id": "whatever",
		"api_app_id": "whatever",
		"event": {
			"user": {
				"id": "whatever",
				"team_id": "whatever",
				"name": "whatever",
				"deleted": true,
				"profile": {
					"title": "",
					"phone": "",
					"skype": "",
					"real_name": "whatever",
					"real_name_normalized": "whatever",
					"display_name": "",
					"display_name_normalized": "",
					"fields": {},
					"status_text": "",
					"status_emoji": "",
					"status_emoji_display_info": [],
					"status_expiration": 0,
					"avatar_hash": "whatever",
					"api_app_id": "whatever",
					"always_active": true,
					"bot_id": "whatever",
					"first_name": "whatever",
					"last_name": "",
					"image_24": "https://secure.gravatar.com/avatar/whatever.jpg",
					"image_32": "https://secure.gravatar.com/avatar/whatever.jpg",
					"image_48": "https://secure.gravatar.com/avatar/whatever.jpg",
					"image_72": "https://secure.gravatar.com/avatar/whatever.jpg",
					"image_192": "https://secure.gravatar.com/avatar/whatever.jpg",
					"image_512": "https://secure.gravatar.com/avatar/whatever.jpg",
					"status_text_canonical": "",
					"team": "whatever"
				},
				"is_bot": true,
				"is_app_user": false,
				"updated": 1678984254
			},
			"cache_ts": 1678984254,
			"type": "user_profile_changed",
			"event_ts": "1678984255.006500"
		},
		"type": "event_callback",
		"event_id": "whatever",
		"event_time": 1678984255,
		"authorizations": [
			{
				"enterprise_id": null,
				"team_id": "whatever",
				"user_id": "whatever",
				"is_bot": false,
				"is_enterprise_install": false
			}
		],
		"is_ext_shared_channel": false
	}
	`)

	evt := &EventsAPICallbackEvent{}
	err := json.Unmarshal(rawE, &evt)
	if err != nil {
		t.Error(err)
	}

	if evt.Type != "event_callback" {
		t.Fail()
	}

	parsedEvent, err := parseInnerEvent(evt)
	if err != nil {
		t.Error(err)
	}

	if parsedEvent.InnerEvent.Type != "user_profile_changed" {
		t.Fail()
	}

	actual, ok := parsedEvent.InnerEvent.Data.(*UserProfileChangedEvent)
	if !ok {
		t.Fail()
	}

	if actual.User.Name != "whatever" {
		t.Fail()
	}
}

func TestSharedChannelInvite(t *testing.T) {
	rawE := []byte(`
	{
		"token": "whatever",
		"team_id": "whatever",
		"api_app_id": "whatever",
		"event": {
			"type": "shared_channel_invite_received",
			"invite": {
				"id": "I028YDERZSQ",
				"date_created": 1626876000,
				"date_invalid": 1628085600,
				"inviting_team": {
					"id": "T12345678",
					"name": "Corgis",
					"icon": {},
					"is_verified": false,
					"domain": "corgis",
					"date_created": 1480946400
				},
				"inviting_user": {
					"id": "U12345678",
					"team_id": "T12345678",
					"name": "crus",
					"updated": 1608081902,
					"profile": {
						"real_name": "Corgis Rus",
						"display_name": "Corgis Rus",
						"real_name_normalized": "Corgis Rus",
						"display_name_normalized": "Corgis Rus",
						"team": "T12345678",
						"avatar_hash": "gcfh83a4c72k",
						"email": "corgisrus@slack-corp.com",
						"image_24": "https://placekitten.com/24/24",
						"image_32": "https://placekitten.com/32/32",
						"image_48": "https://placekitten.com/48/48",
						"image_72": "https://placekitten.com/72/72",
						"image_192": "https://placekitten.com/192/192",
						"image_512": "https://placekitten.com/512/512"
					}
				},
				"recipient_user_id": "U87654321"
			},
			"channel": {
				"id": "C12345678",
				"is_private": false,
				"is_im": false,
				"name": "test-slack-connect"
			},
			"event_ts": "1626876010.000100"
		}
	}
	`)

	evt := &EventsAPICallbackEvent{}
	err := json.Unmarshal(rawE, evt)
	if err != nil {
		t.Fatal(err)
	}

	parsedEvent, err := parseInnerEvent(evt)
	if err != nil {
		t.Fatal(err)
	}

	actual, ok := parsedEvent.InnerEvent.Data.(*SharedChannelInviteReceivedEvent)
	if !ok {
		t.Fail()
	}

	if actual.Invite.ID != "I028YDERZSQ" {
		t.Fail()
	}

	if actual.Invite.InvitingTeam.ID != "T12345678" {
		t.Fail()
	}

	if actual.Invite.InvitingUser.ID != "U12345678" {
		t.Fail()
	}

	if actual.Invite.RecipientUserID != "U87654321" {
		t.Fail()
	}

	if actual.Channel.ID != "C12345678" {
		t.Fail()
	}

	if parsedEvent.InnerEvent.Type != "shared_channel_invite_received" {
		t.Fail()
	}

}

// Test that the shared_channel_invite_accepted event can be unmarshalled
func TestSharedChannelAccepted(t *testing.T) {
	rawE := []byte(`
	{
		"token": "whatever",
		"team_id": "whatever",
		"api_app_id": "whatever",
		"event": {
			"type": "shared_channel_invite_accepted",
			"approval_required": false,
			"invite": {
				"id": "I028YDERZSQ",
				"date_created": 1626876000,
				"date_invalid": 1628085600,
				"inviting_team": {
					"id": "T12345678",
					"name": "Corgis",
					"icon": {
						"image_default": true,
						"image_34": "https://a.slack-edge.com/80588/img/avatars-teams/ava_0011-34.png",
						"image_44": "https://a.slack-edge.com/80588/img/avatars-teams/ava_0011-44.png",
						"image_68": "https://a.slack-edge.com/80588/img/avatars-teams/ava_0011-68.png",
						"image_88": "https://a.slack-edge.com/80588/img/avatars-teams/ava_0011-88.png",
						"image_102": "https://a.slack-edge.com/80588/img/avatars-teams/ava_0011-102.png",
						"image_230": "https://a.slack-edge.com/80588/img/avatars-teams/ava_0011-230.png",
						"image_132": "https://a.slack-edge.com/80588/img/avatars-teams/ava_0011-132.png"
					  },
					"is_verified": false,
					"domain": "corgis",
					"date_created": 1480946400
				},
				"inviting_user": {
					"id": "U12345678",
					"team_id": "T12345678",
					"name": "crus",
					"updated": 1608081902,
					"profile": {
						"real_name": "Corgis Rus",
						"display_name": "Corgis Rus",
						"real_name_normalized": "Corgis Rus",
						"display_name_normalized": "Corgis Rus",
						"team": "T12345678",
						"avatar_hash": "gcfh83a4c72k",
						"email": "corgisrus@slack-corp.com",
						"image_24": "https://placekitten.com/24/24",
						"image_32": "https://placekitten.com/32/32",
						"image_48": "https://placekitten.com/48/48",
						"image_72": "https://placekitten.com/72/72",
						"image_192": "https://placekitten.com/192/192",
						"image_512": "https://placekitten.com/512/512"
					}
				},
				"recipient_email": "golden@doodle.com",
				"recipient_user_id": "U87654321"
			},
			"channel": {
				"id": "C12345678",
				"is_private": false,
				"is_im": false,
				"name": "test-slack-connect"
			},
			"teams_in_channel": [
				{
				"id": "T12345678",
				"name": "Corgis",
				"icon": {
					"image_default": true,
					"image_34": "https://a.slack-edge.com/80588/img/avatars-teams/ava_0011-34.png",
					"image_44": "https://a.slack-edge.com/80588/img/avatars-teams/ava_0011-44.png",
					"image_68": "https://a.slack-edge.com/80588/img/avatars-teams/ava_0011-68.png",
					"image_88": "https://a.slack-edge.com/80588/img/avatars-teams/ava_0011-88.png",
					"image_102": "https://a.slack-edge.com/80588/img/avatars-teams/ava_0011-102.png",
					"image_230": "https://a.slack-edge.com/80588/img/avatars-teams/ava_0011-230.png",
					"image_132": "https://a.slack-edge.com/80588/img/avatars-teams/ava_0011-132.png"
				  },
				"is_verified": false,
				"domain": "corgis",
				"date_created": 1626789600
				}
			],
			"accepting_user": {
				"id": "U87654321",
				"team_id": "T87654321",
				"name": "golden",
				"updated": 1624406113,
				"profile": {
					"real_name": "Golden Doodle",
					"display_name": "Golden",
					"real_name_normalized": "Golden Doodle",
					"display_name_normalized": "Golden",
					"team": "T87654321",
					"avatar_hash": "g717728b118x",
					"email": "golden@doodle.com",
					"image_24": "https://placekitten.com/24/24",
					"image_32": "https://placekitten.com/32/32",
					"image_48": "https://placekitten.com/48/48",
					"image_72": "https://placekitten.com/72/72",
					"image_192": "https://placekitten.com/192/192",
					"image_512": "https://placekitten.com/512/512"
				}
			},
			"event_ts": "1626877800.000000"
		}
	}
	`)

	evt := &EventsAPICallbackEvent{}
	err := json.Unmarshal(rawE, evt)
	if err != nil {
		t.Fatal(err)
	}

	parsedEvent, err := parseInnerEvent(evt)
	if err != nil {
		t.Fatal(err)
	}

	actual, ok := parsedEvent.InnerEvent.Data.(*SharedChannelInviteAcceptedEvent)
	if !ok {
		t.Fail()
	}

	if actual.Invite.ID != "I028YDERZSQ" {
		t.Fail()
	}

	if actual.Invite.InvitingTeam.ID != "T12345678" {
		t.Fail()
	}

	if actual.Invite.InvitingUser.ID != "U12345678" {
		t.Fail()
	}

	if actual.Invite.RecipientUserID != "U87654321" {
		t.Fail()
	}

	if actual.Channel.ID != "C12345678" {
		t.Fail()
	}

	if actual.Channel.Name != "test-slack-connect" {
		t.Fail()
		fmt.Println(actual.Channel.Name + ", does not match the test name.")
	}

	if actual.AcceptingUser.ID != "U87654321" {
		t.Fail()
	}

	if actual.AcceptingUser.Profile.RealName != "Golden Doodle" {
		t.Fail()
	}

	if parsedEvent.InnerEvent.Type != "shared_channel_invite_accepted" {
		t.Fail()
	}

}

// Test that the shared_channel_invite_declined event can be unmarshalled
func TestSharedChannelApproved(t *testing.T) {
	rawE := []byte(`
	{
		"token": "whatever",
		"team_id": "whatever",
		"api_app_id": "whatever",
		"event": {
			"type": "shared_channel_invite_approved",
			"invite": {
				"id": "I01354X80CA",
				"date_created": 1626876000,
				"date_invalid": 1628085600,
				"inviting_team": {
					"id": "T12345678",
					"name": "Corgis",
					"icon": {
						"image_default": true,
						"image_34": "https://a.slack-edge.com/80588/img/avatars-teams/ava_0011-34.png",
						"image_44": "https://a.slack-edge.com/80588/img/avatars-teams/ava_0011-44.png",
						"image_68": "https://a.slack-edge.com/80588/img/avatars-teams/ava_0011-68.png",
						"image_88": "https://a.slack-edge.com/80588/img/avatars-teams/ava_0011-88.png",
						"image_102": "https://a.slack-edge.com/80588/img/avatars-teams/ava_0011-102.png",
						"image_230": "https://a.slack-edge.com/80588/img/avatars-teams/ava_0011-230.png",
						"image_132": "https://a.slack-edge.com/80588/img/avatars-teams/ava_0011-132.png"
					  },
					"is_verified": false,
					"domain": "corgis",
					"date_created": 1480946400
				},
				"inviting_user": {
					"id": "U12345678",
					"team_id": "T12345678",
					"name": "crus",
					"updated": 1608081902,
					"profile": {
						"real_name": "Corgis Rus",
						"display_name": "Corgis Rus",
						"real_name_normalized": "Corgis Rus",
						"display_name_normalized": "Corgis Rus",
						"team": "T12345678",
						"avatar_hash": "gcfh83a4c72k",
						"email": "corgisrus@slack-corp.com",
						"image_24": "https://placekitten.com/24/24",
						"image_32": "https://placekitten.com/32/32",
						"image_48": "https://placekitten.com/48/48",
						"image_72": "https://placekitten.com/72/72",
						"image_192": "https://placekitten.com/192/192",
						"image_512": "https://placekitten.com/512/512"
					}
				},
				"recipient_email": "golden@doodle.com",
				"recipient_user_id": "U87654321"
			},
			"channel": {
				"id": "C12345678",
				"is_private": false,
				"is_im": false,
				"name": "test-slack-connect"
			},
			"approving_team_id": "T87654321",
			"teams_in_channel": [
				{
				"id": "T12345678",
				"name": "Corgis",
				"icon": {
					"image_default": true,
					"image_34": "https://a.slack-edge.com/80588/img/avatars-teams/ava_0011-34.png",
					"image_44": "https://a.slack-edge.com/80588/img/avatars-teams/ava_0011-44.png",
					"image_68": "https://a.slack-edge.com/80588/img/avatars-teams/ava_0011-68.png",
					"image_88": "https://a.slack-edge.com/80588/img/avatars-teams/ava_0011-88.png",
					"image_102": "https://a.slack-edge.com/80588/img/avatars-teams/ava_0011-102.png",
					"image_230": "https://a.slack-edge.com/80588/img/avatars-teams/ava_0011-230.png",
					"image_132": "https://a.slack-edge.com/80588/img/avatars-teams/ava_0011-132.png"
				  },
				"is_verified": false,
				"domain": "corgis",
				"date_created": 1626789600
				}
			],
			"approving_user": {
				"id": "U012A3CDE",
				"team_id": "T87654321",
				"name": "spengler",
				"updated": 1624406532,
				"profile": {
					"real_name": "Egon Spengler",
					"display_name": "Egon",
					"real_name_normalized": "Egon Spengler",
					"display_name_normalized": "Egon",
					"team": "T87654321",
					"avatar_hash": "g216425b1681",
					"email": "spengler@ghostbusters.example.com",
					"image_24": "https://placekitten.com/24/24",
					"image_32": "https://placekitten.com/32/32",
					"image_48": "https://placekitten.com/48/48",
					"image_72": "https://placekitten.com/72/72",
					"image_192": "https://placekitten.com/192/192",
					"image_512": "https://placekitten.com/512/512"
				}
			},
			"event_ts": "1626881400.000000"
		}
	}
	`)

	evt := &EventsAPICallbackEvent{}
	err := json.Unmarshal(rawE, evt)
	if err != nil {
		t.Fatal(err)
	}

	parsedEvent, err := parseInnerEvent(evt)
	if err != nil {
		t.Fatal(err)
	}

	actual, ok := parsedEvent.InnerEvent.Data.(*SharedChannelInviteApprovedEvent)
	if !ok {
		t.Fail()
	}

	if actual.Invite.ID != "I01354X80CA" {
		t.Fail()
	}

	if actual.Invite.InvitingTeam.ID != "T12345678" {
		t.Fail()
	}

	if actual.Invite.InvitingUser.ID != "U12345678" {
		t.Fail()
	}

	if actual.Invite.RecipientUserID != "U87654321" {
		t.Fail()
	}

	if actual.Channel.ID != "C12345678" {
		t.Fail()
	}

	if actual.ApprovingTeamID != "T87654321" {
		t.Fail()
	}

	if actual.ApprovingUser.Name != "spengler" {
		t.Fail()
	}

	if actual.ApprovingUser.Profile.RealName != "Egon Spengler" {
		t.Fail()
	}

	if actual.TeamsInChannel[0].ID != "T12345678" {
		t.Fail()
	}

	if parsedEvent.InnerEvent.Type != "shared_channel_invite_approved" {
		t.Fail()
	}

}

func TestSharedChannelDeclined(t *testing.T) {
	rawE := []byte(`
	{
		"token": "whatever",
		"team_id": "whatever",
		"api_app_id": "whatever",
		"event": {
			"type": "shared_channel_invite_declined",
			"invite": {
				"id": "I01354X80CA",
				"date_created": 1626876000,
				"date_invalid": 1628085600,
				"inviting_team": {
					"id": "T12345678",
					"name": "Corgis",
					"icon": {},
					"is_verified": false,
					"domain": "corgis",
					"date_created": 1480946400
				},
				"inviting_user": {
					"id": "U12345678",
					"team_id": "T12345678",
					"name": "crus",
					"updated": 1608081902,
					"profile": {
						"real_name": "Corgis Rus",
						"display_name": "Corgis Rus",
						"real_name_normalized": "Corgis Rus",
						"display_name_normalized": "Corgis Rus",
						"team": "T12345678",
						"avatar_hash": "gcfh83a4c72k",
						"email": "corgisrus@slack-corp.com",
						"image_24": "https://placekitten.com/24/24",
						"image_32": "https://placekitten.com/32/32",
						"image_48": "https://placekitten.com/48/48",
						"image_72": "https://placekitten.com/72/72",
						"image_192": "https://placekitten.com/192/192",
						"image_512": "https://placekitten.com/512/512"
					}
				},
				"recipient_email": "golden@doodle.com"
			},
			"channel": {
				"id": "C12345678",
				"is_private": false,
				"is_im": false,
				"name": "test-slack-connect"
			},
			"declining_team_id": "T87654321",
			"teams_in_channel": [
				{
					"id": "T12345678",
					"name": "Corgis",
					"icon": {},
					"is_verified": false,
					"domain": "corgis",
					"date_created": 1626789600
				}
			],
			"declining_user": {
				"id": "U012A3CDE",
				"team_id": "T87654321",
				"name": "spengler",
				"updated": 1624406532,
					"profile": {
					"real_name": "Egon Spengler",
					"display_name": "Egon",
					"real_name_normalized": "Egon Spengler",
					"display_name_normalized": "Egon",
					"team": "T87654321",
					"avatar_hash": "g216425b1681",
					"email": "spengler@ghostbusters.example.com",
					"image_24": "https://placekitten.com/24/24",
					"image_32": "https://placekitten.com/32/32",
					"image_48": "https://placekitten.com/48/48",
					"image_72": "https://placekitten.com/72/72",
					"image_192": "https://placekitten.com/192/192",
					"image_512": "https://placekitten.com/512/512"
				}
			},
			"event_ts": "1626881400.000000"
		}
	}
	`)

	evt := &EventsAPICallbackEvent{}
	err := json.Unmarshal(rawE, evt)
	if err != nil {
		t.Fatal(err)
	}

	parsedEvent, err := parseInnerEvent(evt)
	if err != nil {
		t.Fatal(err)
	}

	actual, ok := parsedEvent.InnerEvent.Data.(*SharedChannelInviteDeclinedEvent)
	if !ok {
		t.Fail()
	}

	if actual.Invite.ID != "I01354X80CA" {
		t.Fail()
	}

	if actual.Invite.InvitingTeam.ID != "T12345678" {
		t.Fail()
	}

	if actual.Invite.InvitingUser.ID != "U12345678" {
		t.Fail()
	}

	if actual.Invite.RecipientEmail != "golden@doodle.com" {
		t.Fail()
	}

	if actual.Channel.ID != "C12345678" {
		t.Fail()
	}

	if actual.DecliningTeamID != "T87654321" {
		t.Fail()
	}

	if actual.DecliningUser.Name != "spengler" {
		t.Fail()
	}

	if actual.DecliningUser.Profile.RealName != "Egon Spengler" {
		t.Fail()
	}

	if actual.TeamsInChannel[0].ID != "T12345678" {
		t.Fail()
	}

	if actual.EventTs != "1626881400.000000" {
		t.Fail()
	}

	if parsedEvent.InnerEvent.Type != "shared_channel_invite_declined" {
		t.Fail()
	}

}

func TestChannelHistoryChangedEvent(t *testing.T) {
	rawE := []byte(`
		{
			"type": "channel_history_changed",
			"latest": "1358877455.000010",
			"ts": "1358877455.000008",
			"event_ts": "1358877455.000011"
		}
	`)

	var e ChannelHistoryChangedEvent
	if err := json.Unmarshal(rawE, &e); err != nil {
		t.Fatal(err)
	}
	if e.Type != "channel_history_changed" {
		t.Errorf("type should be channel_history_changed, was %s", e.Type)
	}
	if e.Latest != "1358877455.000010" {
		t.Errorf("latest should be 1358877455.000010, was %s", e.Latest)
	}
	if e.Ts != "1358877455.000008" {
		t.Errorf("ts should be 1358877455.000008, was %s", e.Ts)
	}
	if e.EventTs != "1358877455.000011" {
		t.Errorf("event_ts should be 1358877455.000011, was %s", e.EventTs)
	}
}

func TestDndUpdatedEvent(t *testing.T) {
	rawE := []byte(`
		{
			"type": "dnd_updated",
			"user": "U1234567890",
			"dnd_status": {
				"dnd_enabled": true,
				"next_dnd_start_ts": 1624473600,
				"next_dnd_end_ts": 1624516800,
				"snooze_enabled": false
			}
		}
	`)

	var e DndUpdatedEvent
	if err := json.Unmarshal(rawE, &e); err != nil {
		t.Fatal(err)
	}
	if e.Type != "dnd_updated" {
		t.Errorf("type should be dnd_updated, was %s", e.Type)
	}
	if e.User != "U1234567890" {
		t.Errorf("user should be U1234567890, was %s", e.User)
	}
	if !e.DndStatus.DndEnabled {
		t.Errorf("dnd_enabled should be true, was %v", e.DndStatus.DndEnabled)
	}
	if e.DndStatus.NextDndStartTs != 1624473600 {
		t.Errorf("next_dnd_start_ts should be 1624473600, was %d", e.DndStatus.NextDndStartTs)
	}
	if e.DndStatus.NextDndEndTs != 1624516800 {
		t.Errorf("next_dnd_end_ts should be 1624516800, was %d", e.DndStatus.NextDndEndTs)
	}
	if e.DndStatus.SnoozeEnabled {
		t.Errorf("snooze_enabled should be false, was %v", e.DndStatus.SnoozeEnabled)
	}
}

func TestEmailDomainChangedEvent(t *testing.T) {
	rawE := []byte(`
		{
			"type": "email_domain_changed",
			"email_domain": "example.com",
			"event_ts": "1234567890.123456"
		}
	`)

	var e EmailDomainChangedEvent
	if err := json.Unmarshal(rawE, &e); err != nil {
		t.Fatal(err)
	}
	if e.Type != "email_domain_changed" {
		t.Errorf("type should be email_domain_changed, was %s", e.Type)
	}
	if e.EmailDomain != "example.com" {
		t.Errorf("email_domain should be example.com, was %s", e.EmailDomain)
	}
	if e.EventTs != "1234567890.123456" {
		t.Errorf("event_ts should be 1234567890.123456, was %s", e.EventTs)
	}
}

func TestFileCommentAddedEvent(t *testing.T) {
	rawE := []byte(`
		{
			"type": "file_comment_added",
			"comment": {
				"id": "Fc1234567890",
				"created": 1624473600,
				"timestamp": 1624473600,
				"user": "U1234567890",
				"comment": "This is a comment"
			},
			"file": {
				"id": "F1234567890"
			}
		}
	`)

	var e FileCommentAddedEvent
	if err := json.Unmarshal(rawE, &e); err != nil {
		t.Fatal(err)
	}
	if e.Type != "file_comment_added" {
		t.Errorf("type should be file_comment_added, was %s", e.Type)
	}
	if e.Comment.ID != "Fc1234567890" {
		t.Errorf("comment.id should be Fc1234567890, was %s", e.Comment.ID)
	}
	if e.Comment.Created != 1624473600 {
		t.Errorf("comment.created should be 1624473600, was %d", e.Comment.Created)
	}
	if e.Comment.Timestamp != 1624473600 {
		t.Errorf("comment.timestamp should be 1624473600, was %d", e.Comment.Timestamp)
	}
	if e.Comment.User != "U1234567890" {
		t.Errorf("comment.user should be U1234567890, was %s", e.Comment.User)
	}
	if e.Comment.Comment != "This is a comment" {
		t.Errorf("comment.comment should be 'This is a comment', was %s", e.Comment.Comment)
	}
	if e.File.ID != "F1234567890" {
		t.Errorf("file.id should be F1234567890, was %s", e.File.ID)
	}
}

func TestGroupCloseEvent(t *testing.T) {
	rawE := []byte(`
		{
			"type": "group_close",
			"user": "U1234567890",
			"channel": "G1234567890"
		}
	`)

	var e GroupCloseEvent
	if err := json.Unmarshal(rawE, &e); err != nil {
		t.Fatal(err)
	}
	if e.Type != "group_close" {
		t.Errorf("type should be group_close, was %s", e.Type)
	}
	if e.User != "U1234567890" {
		t.Errorf("user should be U1234567890, was %s", e.User)
	}
	if e.Channel != "G1234567890" {
		t.Errorf("channel should be G1234567890, was %s", e.Channel)
	}
}

func TestImCloseEvent(t *testing.T) {
	rawE := []byte(`
		{
			"type": "im_close",
			"user": "U1234567890",
			"channel": "D1234567890"
		}
	`)

	var e ImCloseEvent
	if err := json.Unmarshal(rawE, &e); err != nil {
		t.Fatal(err)
	}
	if e.Type != "im_close" {
		t.Errorf("type should be im_close, was %s", e.Type)
	}
	if e.User != "U1234567890" {
		t.Errorf("user should be U1234567890, was %s", e.User)
	}
	if e.Channel != "D1234567890" {
		t.Errorf("channel should be D1234567890, was %s", e.Channel)
	}
}

func TestImCreatedEvent(t *testing.T) {
	rawE := []byte(`
		{
			"type": "im_created",
			"user": "U1234567890",
			"channel": {
				"id": "C12345678"
			}
		}
	`)

	var e ImCreatedEvent
	if err := json.Unmarshal(rawE, &e); err != nil {
		t.Fatal(err)
	}
	if e.Type != "im_created" {
		t.Errorf("type should be im_created, was %s", e.Type)
	}
	if e.User != "U1234567890" {
		t.Errorf("user should be U1234567890, was %s", e.User)
	}
	if e.Channel.ID != "C12345678" {
		t.Errorf("channel.id should be C12345678, was %s", e.Channel.ID)
	}
}

func TestImHistoryChangedEvent(t *testing.T) {
	rawE := []byte(`
		{
    		"type": "im_history_changed",
    		"latest": "1358877455.000010",
    		"ts": "1361482916.000003",
    		"event_ts": "1361482916.000004"
		}
	`)

	var e ImHistoryChangedEvent
	if err := json.Unmarshal(rawE, &e); err != nil {
		t.Fatal(err)
	}
	if e.Type != "im_history_changed" {
		t.Errorf("type should be im_created, was %s", e.Type)
	}
	if e.Latest != "1358877455.000010" {
		t.Errorf("latest should be 1358877455.000010, was %s", e.Latest)
	}
	if e.Ts != "1361482916.000003" {
		t.Errorf("ts should be 1361482916.000003, was %s", e.Ts)
	}
	if e.EventTs != "1361482916.000004" {
		t.Errorf("event_ts should be 1361482916.000004, was %s", e.EventTs)
	}
}

func TestImOpenEvent(t *testing.T) {
	rawE := []byte(`
		{
			"type": "im_open",
			"user": "U1234567890",
			"channel": "D1234567890"
		}
	`)

	var e ImOpenEvent
	if err := json.Unmarshal(rawE, &e); err != nil {
		t.Fatal(err)
	}
	if e.Type != "im_open" {
		t.Errorf("type should be im_open, was %s", e.Type)
	}
	if e.User != "U1234567890" {
		t.Errorf("user should be U1234567890, was %s", e.User)
	}
	if e.Channel != "D1234567890" {
		t.Errorf("channel should be D1234567890, was %s", e.Channel)
	}
}

func TestSubteamCreatedEvent(t *testing.T) {
	rawE := []byte(`
		{
			"type": "subteam_created",
			"subteam": {
				"id": "S1234567890",
				"team_id": "T1234567890",
				"is_usergroup": true,
				"name": "subteam",
				"description": "A test subteam",
				"handle": "subteam_handle",
				"is_external": false,
				"date_create": 1624473600,
				"date_update": 1624473600,
				"date_delete": 0,
				"auto_type": "auto",
				"created_by": "U1234567890",
				"updated_by": "U1234567890",
				"deleted_by": "",
				"prefs": {
					"channels": ["C1234567890"],
					"groups": ["G1234567890"]
				},
				"users": ["U1234567890"],
				"user_count": 1
			}
		}
	`)

	var e SubteamCreatedEvent
	if err := json.Unmarshal(rawE, &e); err != nil {
		t.Fatal(err)
	}
	if e.Type != "subteam_created" {
		t.Errorf("type should be subteam_created, was %s", e.Type)
	}
	if e.Subteam.ID != "S1234567890" {
		t.Errorf("subteam.id should be S1234567890, was %s", e.Subteam.ID)
	}
	if e.Subteam.TeamID != "T1234567890" {
		t.Errorf("subteam.team_id should be T1234567890, was %s", e.Subteam.TeamID)
	}
	if !e.Subteam.IsUsergroup {
		t.Errorf("subteam.is_usergroup should be true, was %v", e.Subteam.IsUsergroup)
	}
	if e.Subteam.Name != "subteam" {
		t.Errorf("subteam.name should be subteam, was %s", e.Subteam.Name)
	}
	if e.Subteam.Description != "A test subteam" {
		t.Errorf("subteam.description should be 'A test subteam', was %s", e.Subteam.Description)
	}
	if e.Subteam.Handle != "subteam_handle" {
		t.Errorf("subteam.handle should be subteam_handle, was %s", e.Subteam.Handle)
	}
	if e.Subteam.IsExternal {
		t.Errorf("subteam.is_external should be false, was %v", e.Subteam.IsExternal)
	}
	if e.Subteam.DateCreate != 1624473600 {
		t.Errorf("subteam.date_create should be 1624473600, was %d", e.Subteam.DateCreate)
	}
	if e.Subteam.DateUpdate != 1624473600 {
		t.Errorf("subteam.date_update should be 1624473600, was %d", e.Subteam.DateUpdate)
	}
	if e.Subteam.DateDelete != 0 {
		t.Errorf("subteam.date_delete should be 0, was %d", e.Subteam.DateDelete)
	}
	if e.Subteam.AutoType != "auto" {
		t.Errorf("subteam.auto_type should be auto, was %s", e.Subteam.AutoType)
	}
	if e.Subteam.CreatedBy != "U1234567890" {
		t.Errorf("subteam.created_by should be U1234567890, was %s", e.Subteam.CreatedBy)
	}
	if e.Subteam.UpdatedBy != "U1234567890" {
		t.Errorf("subteam.updated_by should be U1234567890, was %s", e.Subteam.UpdatedBy)
	}
	if e.Subteam.DeletedBy != "" {
		t.Errorf("subteam.deleted_by should be empty, was %s", e.Subteam.DeletedBy)
	}
	if len(e.Subteam.Prefs.Channels) != 1 || e.Subteam.Prefs.Channels[0] != "C1234567890" {
		t.Errorf("subteam.prefs.channels should contain C1234567890, was %v", e.Subteam.Prefs.Channels)
	}
	if len(e.Subteam.Prefs.Groups) != 1 || e.Subteam.Prefs.Groups[0] != "G1234567890" {
		t.Errorf("subteam.prefs.groups should contain G1234567890, was %v", e.Subteam.Prefs.Groups)
	}
	if len(e.Subteam.Users) != 1 || e.Subteam.Users[0] != "U1234567890" {
		t.Errorf("subteam.users should contain U1234567890, was %v", e.Subteam.Users)
	}
	if e.Subteam.UserCount != 1 {
		t.Errorf("subteam.user_count should be 1, was %d", e.Subteam.UserCount)
	}
}

func TestSubteamMembersChangedEvent(t *testing.T) {
	rawE := []byte(`
		{
			"type": "subteam_members_changed",
			"subteam_id": "S1234567890",
			"team_id": "T1234567890",
			"date_update": 1624473600,
			"added_users": ["U1234567890"],
			"removed_users": ["U0987654321"]
		}
	`)

	var e SubteamMembersChangedEvent
	if err := json.Unmarshal(rawE, &e); err != nil {
		t.Fatal(err)
	}
	if e.Type != "subteam_members_changed" {
		t.Errorf("type should be subteam_members_changed, was %s", e.Type)
	}
	if e.SubteamID != "S1234567890" {
		t.Errorf("subteam_id should be S1234567890, was %s", e.SubteamID)
	}
	if e.TeamID != "T1234567890" {
		t.Errorf("team_id should be T1234567890, was %s", e.TeamID)
	}
	if e.DateUpdate != 1624473600 {
		t.Errorf("date_update should be 1624473600, was %d", e.DateUpdate)
	}
	if len(e.AddedUsers) != 1 || e.AddedUsers[0] != "U1234567890" {
		t.Errorf("subteam.users should contain U1234567890, was %v", e.AddedUsers)
	}
	if len(e.RemovedUsers) != 1 || e.RemovedUsers[0] != "U0987654321" {
		t.Errorf("subteam.users should contain U0987654321, was %v", e.RemovedUsers)
	}
}

func TestSubteamSelfAddedEvent(t *testing.T) {
	rawE := []byte(`
		{
			"type": "subteam_self_added",
			"subteam_id": "S1234567890"
		}
	`)

	var e SubteamSelfAddedEvent
	if err := json.Unmarshal(rawE, &e); err != nil {
		t.Fatal(err)
	}
	if e.Type != "subteam_self_added" {
		t.Errorf("type should be subteam_self_added, was %s", e.Type)
	}
	if e.SubteamID != "S1234567890" {
		t.Errorf("subteam_id should be S1234567890, was %s", e.SubteamID)
	}
}

func TestSubteamSelfRemovedEvent(t *testing.T) {
	rawE := []byte(`
		{
			"type": "subteam_self_removed",
			"subteam_id": "S1234567890"
		}
	`)

	var e SubteamSelfRemovedEvent
	if err := json.Unmarshal(rawE, &e); err != nil {
		t.Fatal(err)
	}
	if e.Type != "subteam_self_removed" {
		t.Errorf("type should be subteam_self_removed, was %s", e.Type)
	}
	if e.SubteamID != "S1234567890" {
		t.Errorf("subteam_id should be S1234567890, was %s", e.SubteamID)
	}
}

func TestSubteamUpdatedEvent(t *testing.T) {
	rawE := []byte(`
		{
			"type": "subteam_updated",
			"subteam": {
				"id": "S1234567890",
				"team_id": "T1234567890",
				"is_usergroup": true,
				"name": "updated_subteam",
				"description": "An updated test subteam",
				"handle": "updated_subteam_handle",
				"is_external": false,
				"date_create": 1624473600,
				"date_update": 1624473600,
				"date_delete": 0,
				"auto_type": "auto",
				"created_by": "U1234567890",
				"updated_by": "U1234567890",
				"deleted_by": "",
				"prefs": {
					"channels": ["C1234567890"],
					"groups": ["G1234567890"]
				},
				"users": ["U1234567890"],
				"user_count": 1
			}
		}
	`)

	var e SubteamUpdatedEvent
	if err := json.Unmarshal(rawE, &e); err != nil {
		t.Fatal(err)
	}
	if e.Type != "subteam_updated" {
		t.Errorf("type should be subteam_updated, was %s", e.Type)
	}
	if e.Subteam.ID != "S1234567890" {
		t.Errorf("subteam.id should be S1234567890, was %s", e.Subteam.ID)
	}
	if e.Subteam.TeamID != "T1234567890" {
		t.Errorf("subteam.team_id should be T1234567890, was %s", e.Subteam.TeamID)
	}
	if !e.Subteam.IsUsergroup {
		t.Errorf("subteam.is_usergroup should be true, was %v", e.Subteam.IsUsergroup)
	}
	if e.Subteam.Name != "updated_subteam" {
		t.Errorf("subteam.name should be updated_subteam, was %s", e.Subteam.Name)
	}
	if e.Subteam.Description != "An updated test subteam" {
		t.Errorf("subteam.description should be 'An updated test subteam', was %s", e.Subteam.Description)
	}
	if e.Subteam.Handle != "updated_subteam_handle" {
		t.Errorf("subteam.handle should be updated_subteam_handle, was %s", e.Subteam.Handle)
	}
	if e.Subteam.IsExternal {
		t.Errorf("subteam.is_external should be false, was %v", e.Subteam.IsExternal)
	}
	if e.Subteam.DateCreate != 1624473600 {
		t.Errorf("subteam.date_create should be 1624473600, was %d", e.Subteam.DateCreate)
	}
	if e.Subteam.DateUpdate != 1624473600 {
		t.Errorf("subteam.date_update should be 1624473600, was %d", e.Subteam.DateUpdate)
	}
	if e.Subteam.DateDelete != 0 {
		t.Errorf("subteam.date_delete should be 0, was %d", e.Subteam.DateDelete)
	}
	if e.Subteam.AutoType != "auto" {
		t.Errorf("subteam.auto_type should be auto, was %s", e.Subteam.AutoType)
	}
	if e.Subteam.CreatedBy != "U1234567890" {
		t.Errorf("subteam.created_by should be U1234567890, was %s", e.Subteam.CreatedBy)
	}
	if e.Subteam.UpdatedBy != "U1234567890" {
		t.Errorf("subteam.updated_by should be U1234567890, was %s", e.Subteam.UpdatedBy)
	}
	if e.Subteam.DeletedBy != "" {
		t.Errorf("subteam.deleted_by should be empty, was %s", e.Subteam.DeletedBy)
	}
	if len(e.Subteam.Prefs.Channels) != 1 || e.Subteam.Prefs.Channels[0] != "C1234567890" {
		t.Errorf("subteam.prefs.channels should contain C1234567890, was %v", e.Subteam.Prefs.Channels)
	}
	if len(e.Subteam.Prefs.Groups) != 1 || e.Subteam.Prefs.Groups[0] != "G1234567890" {
		t.Errorf("subteam.prefs.groups should contain G1234567890, was %v", e.Subteam.Prefs.Groups)
	}
	if len(e.Subteam.Users) != 1 || e.Subteam.Users[0] != "U1234567890" {
		t.Errorf("subteam.users should contain U1234567890, was %v", e.Subteam.Users)
	}
	if e.Subteam.UserCount != 1 {
		t.Errorf("subteam.user_count should be 1, was %d", e.Subteam.UserCount)
	}
}

func TestTeamDomainChangeEvent(t *testing.T) {
	rawE := []byte(`
		{
			"type": "team_domain_change",
			"url": "https://newdomain.slack.com",
			"domain": "newdomain"
		}
	`)

	var e TeamDomainChangeEvent
	if err := json.Unmarshal(rawE, &e); err != nil {
		t.Fatal(err)
	}
	if e.Type != "team_domain_change" {
		t.Errorf("type should be team_domain_change, was %s", e.Type)
	}
	if e.URL != "https://newdomain.slack.com" {
		t.Errorf("url should be https://newdomain.slack.com, was %s", e.URL)
	}
	if e.Domain != "newdomain" {
		t.Errorf("domain should be newdomain, was %s", e.Domain)
	}
}

func TestTeamRenameEvent(t *testing.T) {
	rawE := []byte(`
		{
			"type": "team_rename",
			"name": "new_team_name"
		}
	`)

	var e TeamRenameEvent
	if err := json.Unmarshal(rawE, &e); err != nil {
		t.Fatal(err)
	}
	if e.Type != "team_rename" {
		t.Errorf("type should be team_rename, was %s", e.Type)
	}
	if e.Name != "new_team_name" {
		t.Errorf("name should be new_team_name, was %s", e.Name)
	}
}

func TestUserChangeEvent(t *testing.T) {
	rawE := []byte(`
		{
			"type": "user_change",
			"user": {
				"id": "U1234567890",
				"team_id": "T1234567890",
				"name": "test_user",
				"deleted": false,
				"profile": {
					"avatar_hash": "123456abcdef",
					"real_name": "Test User",
					"display_name": "TestUser",
					"real_name_normalized": "Test User",
					"display_name_normalized": "TestUser",
					"email": "testuser@example.com",
					"image_24": "https://example.com/image_24.jpg",
					"image_32": "https://example.com/image_32.jpg",
					"image_48": "https://example.com/image_48.jpg",
					"image_72": "https://example.com/image_72.jpg",
					"image_192": "https://example.com/image_192.jpg",
					"image_512": "https://example.com/image_512.jpg"
				},
				"is_bot": false,
				"updated": 1624473600
			}
		}
	`)

	var e UserChangeEvent
	if err := json.Unmarshal(rawE, &e); err != nil {
		t.Fatal(err)
	}
	if e.Type != "user_change" {
		t.Errorf("type should be user_change, was %s", e.Type)
	}
	if e.User.ID != "U1234567890" {
		t.Errorf("user.id should be U1234567890, was %s", e.User.ID)
	}
	if e.User.TeamID != "T1234567890" {
		t.Errorf("user.team_id should be T1234567890, was %s", e.User.TeamID)
	}
	if e.User.Name != "test_user" {
		t.Errorf("user.name should be test_user, was %s", e.User.Name)
	}
	if e.User.Deleted {
		t.Errorf("user.deleted should be false, was %v", e.User.Deleted)
	}
	if e.User.Profile.AvatarHash != "123456abcdef" {
		t.Errorf("user.profile.avatar_hash should be 123456abcdef, was %s", e.User.Profile.AvatarHash)
	}
	if e.User.Profile.RealName != "Test User" {
		t.Errorf("user.profile.real_name should be 'Test User', was %s", e.User.Profile.RealName)
	}
	if e.User.Profile.DisplayName != "TestUser" {
		t.Errorf("user.profile.display_name should be 'TestUser', was %s", e.User.Profile.DisplayName)
	}
	if e.User.Profile.RealNameNormalized != "Test User" {
		t.Errorf("user.profile.real_name_normalized should be 'Test User', was %s", e.User.Profile.RealNameNormalized)
	}
	if e.User.Profile.DisplayNameNormalized != "TestUser" {
		t.Errorf("user.profile.display_name_normalized should be 'TestUser', was %s", e.User.Profile.DisplayNameNormalized)
	}
	if e.User.Profile.Email != "testuser@example.com" {
		t.Errorf("user.profile.email should be 'testuser@example.com', was %s", e.User.Profile.Email)
	}
	if e.User.Profile.Image24 != "https://example.com/image_24.jpg" {
		t.Errorf("user.profile.image_24 should be 'https://example.com/image_24.jpg', was %s", e.User.Profile.Image24)
	}
	if e.User.Profile.Image32 != "https://example.com/image_32.jpg" {
		t.Errorf("user.profile.image_32 should be 'https://example.com/image_32.jpg', was %s", e.User.Profile.Image32)
	}
	if e.User.Profile.Image48 != "https://example.com/image_48.jpg" {
		t.Errorf("user.profile.image_48 should be 'https://example.com/image_48.jpg', was %s", e.User.Profile.Image48)
	}
	if e.User.Profile.Image72 != "https://example.com/image_72.jpg" {
		t.Errorf("user.profile.image_72 should be 'https://example.com/image_72.jpg', was %s", e.User.Profile.Image72)
	}
	if e.User.Profile.Image192 != "https://example.com/image_192.jpg" {
		t.Errorf("user.profile.image_192 should be 'https://example.com/image_192.jpg', was %s", e.User.Profile.Image192)
	}
	if e.User.Profile.Image512 != "https://example.com/image_512.jpg" {
		t.Errorf("user.profile.image_512 should be 'https://example.com/image_512.jpg', was %s", e.User.Profile.Image512)
	}
	if e.User.IsBot {
		t.Errorf("user.is_bot should be false, was %v", e.User.IsBot)
	}
	if e.User.Updated != 1624473600 {
		t.Errorf("user.updated should be 1624473600, was %d", e.User.Updated)
	}
}
