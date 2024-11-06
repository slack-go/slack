package slackevents

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAssistantThreadStartedEvent(t *testing.T) {

	rawE := []byte(`
		{
			"type": "assistant_thread_started",
			"assistant_thread": {
				"user_id": "U123ABC456",
				"context": { 
					"channel_id": "C123ABC456", 
					"team_id": "T07XY8FPJ5C", 
					"enterprise_id": "E480293PS82"
					},
				"channel_id": "D123ABC456",
				"thread_ts": "1729999327.187299"

			},
			"event_ts": "1715873754.429808"
		}
	`)

	err := json.Unmarshal(rawE, &AssistantThreadStartedEvent{})
	if err != nil {
		t.Error(err)
	}

}

func TestAssistantThreadContextChangedEvent(t *testing.T) {

	rawE := []byte(`
		{
			"type": "assistant_thread_context_changed",
			"assistant_thread": {
				"user_id": "U123ABC456",
				"context": { 
					"channel_id": "C123ABC456", 
					"team_id": "T07XY8FPJ5C", 
					"enterprise_id": "E480293PS82"
					},
				"channel_id": "D123ABC456",
				"thread_ts": "1729999327.187299"
			},
			"event_ts": "17298244.022142"
		}
	`)

	err := json.Unmarshal(rawE, &AssistantThreadContextChangedEvent{})
	if err != nil {
		t.Error(err)
	}

}

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
				"metadata": {
					"event_type": "example",
					"event_payload": {
						"key": "value"
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

func TestDndUpdatedUserEvent(t *testing.T) {
	rawE := []byte(`
		{
    		"type": "dnd_updated_user",
    		"user": "U1234",
    		"dnd_status": {
        		"dnd_enabled": true,
        		"next_dnd_start_ts": 1450387800,
        		"next_dnd_end_ts": 1450423800
    		}
		}
	`)

	var e DndUpdatedUserEvent
	if err := json.Unmarshal(rawE, &e); err != nil {
		t.Fatal(err)
	}
	if e.Type != "dnd_updated_user" {
		t.Errorf("type should be dnd_updated_user, was %s", e.Type)
	}
	if e.User != "U1234" {
		t.Errorf("user should be U1234, was %s", e.User)
	}
	if !e.DndStatus.DndEnabled {
		t.Errorf("dnd_enabled should be true, was %v", e.DndStatus.DndEnabled)
	}
	if e.DndStatus.NextDndStartTs != 1450387800 {
		t.Errorf("next_dnd_start_ts should be 1450387800, was %d", e.DndStatus.NextDndStartTs)
	}
	if e.DndStatus.NextDndEndTs != 1450423800 {
		t.Errorf("next_dnd_end_ts should be 1450423800, was %d", e.DndStatus.NextDndEndTs)
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

func TestGroupHistoryChangedEvent(t *testing.T) {
	rawE := []byte(`
		{
    		"type": "group_history_changed",
    		"latest": "1358877455.000010",
    		"ts": "1361482916.000003",
    		"event_ts": "1361482916.000004"
		}
	`)

	var e GroupHistoryChangedEvent
	if err := json.Unmarshal(rawE, &e); err != nil {
		t.Fatal(err)
	}
	if e.Type != "group_history_changed" {
		t.Errorf("type should be group_history_changed, was %s", e.Type)
	}
	if e.Latest != "1358877455.000010" {
		t.Errorf("latest should be 1358877455.000010, was %s", e.Latest)
	}
	if e.Ts != "1361482916.000003" {
		t.Errorf("ts should be 1361482916.000003, was %s", e.Ts)
	}
}

func TestGroupOpenEvent(t *testing.T) {
	rawE := []byte(`
		{
    		"type": "group_open",
    		"user": "U024BE7LH",
    		"channel": "G024BE91L"
		}
	`)

	var e GroupOpenEvent
	if err := json.Unmarshal(rawE, &e); err != nil {
		t.Fatal(err)
	}
	if e.Type != "group_open" {
		t.Errorf("type should be group_open, was %s", e.Type)
	}
	if e.User != "U024BE7LH" {
		t.Errorf("user should be U024BE7LH, was %s", e.User)
	}
	if e.Channel != "G024BE91L" {
		t.Errorf("channel should be G024BE91L, was %s", e.Channel)
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
			"date_previous_update": 1446670362,
			"date_update": 1624473600,
			"added_users": ["U1234567890"],
			"added_users_count": "3",
			"removed_users": ["U0987654321"],
			"removed_users_count": "1"
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
			"domain": "newdomain",
			"team_id": "T1234"
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
	if e.TeamID != "T1234" {
		t.Errorf("team_id should be 'T1234', was %s", e.TeamID)
	}
}

func TestTeamRenameEvent(t *testing.T) {
	rawE := []byte(`
		{
			"type": "team_rename",
			"name": "new_team_name",
			"team_id": "T1234"
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
	if e.TeamID != "T1234" {
		t.Errorf("team_id should be 'T1234', was %s", e.TeamID)
	}
}

func TestUserChangeEvent(t *testing.T) {
	jsonStr := `{
		"user": {
			"id": "U1234567",
			"team_id": "T1234567",
			"name": "some-user",
			"deleted": false,
			"color": "4bbe2e",
			"real_name": "Some User",
			"tz": "America/Los_Angeles",
			"tz_label": "Pacific Daylight Time",
			"tz_offset": -25200,
			"profile": {
				"title": "",
				"phone": "",
				"skype": "",
				"real_name": "Some User",
				"real_name_normalized": "Some User",
				"display_name": "",
				"display_name_normalized": "",
				"fields": {},
				"status_text": "riding a train",
				"status_emoji": ":mountain_railway:",
				"status_emoji_display_info": [],
				"status_expiration": 0,
				"avatar_hash": "g12345678910",
				"first_name": "Some",
				"last_name": "User",
				"image_24": "https://secure.gravatar.com/avatar/cb0c2b2ca5e8de16be31a55a734d0f31.jpg?s=24&d=https%3A%2F%2Fdev.slack.com%2Fdev-cdn%2Fv1648136338%2Fimg%2Favatars%2Fuser_shapes%2Fava_0001-24.png",
				"image_32": "https://secure.gravatar.com/avatar/cb0c2b2ca5e8de16be31a55a734d0f31.jpg?s=32&d=https%3A%2F%2Fdev.slack.com%2Fdev-cdn%2Fv1648136338%2Fimg%2Favatars%2Fuser_shapes%2Fava_0001-32.png",
				"image_48": "https://secure.gravatar.com/avatar/cb0c2b2ca5e8de16be31a55a734d0f31.jpg?s=48&d=https%3A%2F%2Fdev.slack.com%2Fdev-cdn%2Fv1648136338%2Fimg%2Favatars%2Fuser_shapes%2Fava_0001-48.png",
				"image_72": "https://secure.gravatar.com/avatar/cb0c2b2ca5e8de16be31a55a734d0f31.jpg?s=72&d=https%3A%2F%2Fdev.slack.com%2Fdev-cdn%2Fv1648136338%2Fimg%2Favatars%2Fuser_shapes%2Fava_0001-72.png",
				"image_192": "https://secure.gravatar.com/avatar/cb0c2b2ca5e8de16be31a55a734d0f31.jpg?s=192&d=https%3A%2F%2Fdev.slack.com%2Fdev-cdn%2Fv1648136338%2Fimg%2Favatars%2Fuser_shapes%2Fava_0001-192.png",
				"image_512": "https://secure.gravatar.com/avatar/cb0c2b2ca5e8de16be31a55a734d0f31.jpg?s=512&d=https%3A%2F%2Fdev.slack.com%2Fdev-cdn%2Fv1648136338%2Fimg%2Favatars%2Fuser_shapes%2Fava_0001-512.png",
				"status_text_canonical": "",
				"team": "T1234567"
			},
			"is_admin": false,
			"is_owner": false,
			"is_primary_owner": false,
			"is_restricted": false,
			"is_ultra_restricted": false,
			"is_bot": false,
			"is_app_user": false,
			"updated": 1648596421,
			"is_email_confirmed": true,
			"who_can_share_contact_card": "EVERYONE",
			"locale": "en-US"
		},
		"cache_ts": 1648596421,
		"type": "user_change",
		"event_ts": "1648596712.000001"
	}`

	var event UserChangeEvent
	if err := json.Unmarshal([]byte(jsonStr), &event); err != nil {
		t.Errorf("Failed to unmarshal UserChangeEvent: %v", err)
	}

	if event.Type != "user_change" {
		t.Errorf("Expected type to be 'user_change', got %s", event.Type)
	}

	if event.User.ID != "U1234567" {
		t.Errorf("Expected user ID to be 'U1234567', got %s", event.User.ID)
	}

	if event.User.Profile.StatusText != "riding a train" {
		t.Errorf("Expected status text to be 'riding a train', got %s", event.User.Profile.StatusText)
	}

	if event.User.Profile.StatusEmoji != ":mountain_railway:" {
		t.Errorf("Expected status emoji to be ':mountain_railway:', got %s", event.User.Profile.StatusEmoji)
	}

	if event.CacheTS != 1648596421 {
		t.Errorf("Expected cache_ts to be 1648596421, got %d", event.CacheTS)
	}

	if event.EventTS != "1648596712.000001" {
		t.Errorf("Expected event_ts to be '1648596712.000001', got %s", event.EventTS)
	}
}

func TestAppDeletedEvent(t *testing.T) {
	jsonStr := `{
		"type": "app_deleted",
		"app_id": "A015CA1LGHG",
		"app_name": "my-admin-app",
		"app_owner_id": "U013B64J7MSZ",
		"team_id": "E073D7H7BBE",
		"team_domain": "ACME Enterprises",
		"event_ts": "1700001891.279278"
	}`

	var event AppDeletedEvent
	if err := json.Unmarshal([]byte(jsonStr), &event); err != nil {
		t.Errorf("Failed to unmarshal AppDeletedEvent: %v", err)
	}

	if event.Type != "app_deleted" {
		t.Errorf("Expected type to be 'app_deleted', got %s", event.Type)
	}

	if event.AppName != "my-admin-app" {
		t.Errorf("app_name should be 'my-admin-app', was %s", event.AppName)
	}

	if event.AppOwnerID != "U013B64J7MSZ" {
		t.Errorf("app_owner_id should be 'U013B64J7MSZ', was %s", event.AppOwnerID)
	}

	if event.TeamID != "E073D7H7BBE" {
		t.Errorf("team_id should be 'E073D7H7BBE', was %s", event.TeamID)
	}
}

func TestAppInstalledEvent(t *testing.T) {
	jsonStr := `{
		"type": "app_installed",
		"app_id": "A015CA1LGHG",
		"app_name": "my-admin-app",
		"app_owner_id": "U013B64J7MSZ",
		"user_id": "U013B64J7SZ",
		"team_id": "E073D7H7BBE",
		"team_domain": "ACME Enterprises",
		"event_ts": "1700001891.279278"
	}`

	var event AppInstalledEvent
	if err := json.Unmarshal([]byte(jsonStr), &event); err != nil {
		t.Errorf("Failed to unmarshal AppInstalledEvent: %v", err)
	}

	if event.Type != "app_installed" {
		t.Errorf("Expected type to be 'app_installed', got %s", event.Type)
	}

	if event.AppName != "my-admin-app" {
		t.Errorf("app_name should be 'my-admin-app', was %s", event.AppName)
	}

	if event.AppOwnerID != "U013B64J7MSZ" {
		t.Errorf("app_owner_id should be 'U013B64J7MSZ', was %s", event.AppOwnerID)
	}

	if event.TeamID != "E073D7H7BBE" {
		t.Errorf("team_id should be 'E073D7H7BBE', was %s", event.TeamID)
	}
}

func TestAppRequestedEvent(t *testing.T) {
	jsonStr := `{
		"type": "app_requested",
		"app_request": {
			"id": "1234",
			"app": {
				"id": "A5678",
				"name": "Brent's app",
				"description": "They're good apps, Bront.",
				"help_url": "brontsapp.com",
				"privacy_policy_url": "brontsapp.com",
				"app_homepage_url": "brontsapp.com",
				"app_directory_url": "https://slack.slack.com/apps/A102ARD7Y",
				"is_app_directory_approved": true,
				"is_internal": false,
				"additional_info": "none"
			},
			"previous_resolution": {
				"status": "approved",
				"scopes": [{
					"name": "app_requested",
					"description": "allows this app to listen for app install requests",
					"is_sensitive": false,
					"token_type": "user"
				}]
			},
			"user": {
				"id": "U1234",
				"name": "Bront",
				"email": "bront@brent.com"
			},
			"team": {
				"id": "T1234",
				"name": "Brant App Team",
				"domain": "brantappteam"
			},
			"enterprise": null,
			"scopes": [{
				"name": "app_requested",
				"description": "allows this app to listen for app install requests",
				"is_sensitive": false,
				"token_type": "user"
			}],
			"message": "none"
		}
	}`

	var event AppRequestedEvent
	if err := json.Unmarshal([]byte(jsonStr), &event); err != nil {
		t.Errorf("Failed to unmarshal AppRequestedEvent: %v", err)
	}

	if event.Type != "app_requested" {
		t.Errorf("Expected type to be 'app_requested', got %s", event.Type)
	}

	if event.AppRequest.ID != "1234" {
		t.Errorf("app_request.id should be '1234', was %s", event.AppRequest.ID)
	}

	if event.AppRequest.App.ID != "A5678" {
		t.Fail()
	}

	if event.AppRequest.User.ID != "U1234" {
		t.Errorf("app_request.user.id should be 'U1234', was %s", event.AppRequest.User.ID)
	}

	if event.AppRequest.Team.ID != "T1234" {
		t.Fail()
	}
}

func TestAppUninstalledTeamEvent(t *testing.T) {
	jsonStr := `{
		"type": "app_uninstalled_team",
		"app_id": "A015CA1LGHG",
		"app_name": "my-admin-app",
		"app_owner_id": "U013B64J7MSZ",
		"user_id": "U013B64J7SZ",
		"team_id": "E073D7H7BBE",
		"team_domain": "ACME Enterprises",
		"event_ts": "1700001891.279278"
	}`

	var event AppUninstalledTeamEvent
	if err := json.Unmarshal([]byte(jsonStr), &event); err != nil {
		t.Errorf("Failed to unmarshal AppUninstalledTeamEvent: %v", err)
	}

	if event.Type != "app_uninstalled_team" {
		t.Errorf("Expected type to be 'app_uninstalled_team', got %s", event.Type)
	}

	if event.AppName != "my-admin-app" {
		t.Errorf("app_name should be 'my-admin-app', was %s", event.AppName)
	}

	if event.AppOwnerID != "U013B64J7MSZ" {
		t.Errorf("app_owner_id should be 'U013B64J7MSZ', was %s", event.AppOwnerID)
	}

	if event.TeamID != "E073D7H7BBE" {
		t.Errorf("team_id should be 'E073D7H7BBE', was %s", event.TeamID)
	}
}

func TestCallRejectedEvent(t *testing.T) {
	jsonStr := `{
		"token": "12345FVmRUzNDOAu12345h",
		"team_id": "T123ABC456",
		"api_app_id": "BBBU04BB4",
		"event": {
			"type": "call_rejected",
			"call_id": "R123ABC456",
			"user_id": "U123ABC456",
			"channel_id": "D123ABC456",
			"external_unique_id": "123-456-7890"
		},
		"type": "event_callback",
		"event_id": "Ev123ABC456",
		"event_time": 1563448153,
		"authed_users": ["U123ABC456"]
	}`

	var event CallRejectedEvent
	if err := json.Unmarshal([]byte(jsonStr), &event); err != nil {
		t.Errorf("Failed to unmarshal CallRejectedEvent: %v", err)
	}

	if event.Event.Type != "call_rejected" {
		t.Errorf("Expected event type to be 'call_rejected', got %s", event.Event.Type)
	}
	if event.TeamID != "T123ABC456" {
		t.Errorf("Expected team_id to be 'T123ABC456', got %s", event.TeamID)
	}
	if event.Event.CallID != "R123ABC456" {
		t.Fail()
	}

}

func TestChannelSharedEvent(t *testing.T) {
	jsonStr := `{
		"type": "channel_shared",
		"connected_team_id": "E163Q94DX",
		"channel": "C123ABC456",
		"event_ts": "1561064063.001100"
	}`

	var event ChannelSharedEvent
	if err := json.Unmarshal([]byte(jsonStr), &event); err != nil {
		t.Errorf("Failed to unmarshal ChannelSharedEvent: %v", err)
	}

	if event.Type != "channel_shared" {
		t.Errorf("Expected type to be 'channel_shared', got %s", event.Type)
	}

	if event.ConnectedTeamID != "E163Q94DX" {
		t.Errorf("Expected connected_team_id to be 'E163Q94DX', got %s", event.ConnectedTeamID)
	}

	if event.Channel != "C123ABC456" {
		t.Fail()
	}
}

func TestFileCreatedEvent(t *testing.T) {
	jsonStr := `{
		"type": "file_created",
		"file_id": "F2147483862",
		"file": {
			"id": "F2147483862"
		}
	}`

	var event FileCreatedEvent
	if err := json.Unmarshal([]byte(jsonStr), &event); err != nil {
		t.Errorf("Failed to unmarshal FileCreatedEvent: %v", err)
	}

	if event.Type != "file_created" {
		t.Errorf("Expected type to be 'file_created', got %s", event.Type)
	}
	if event.FileID != "F2147483862" {
		t.Errorf("Expected file_id to be 'F2147483862', got %s", event.FileID)
	}
}

func TestFilePublicEvent(t *testing.T) {
	jsonStr := `{
		"type": "file_public",
		"file_id": "F2147483862",
		"file": {
			"id": "F2147483862"
		}
	}`

	var event FilePublicEvent
	if err := json.Unmarshal([]byte(jsonStr), &event); err != nil {
		t.Errorf("Failed to unmarshal FilePublicEvent: %v", err)
	}

	if event.Type != "file_public" {
		t.Errorf("Expected type to be 'file_public', got %s", event.Type)
	}

	if event.FileID != "F2147483862" {
		t.Errorf("Expected file_id to be 'F2147483862', got %s", event.FileID)
	}
}

func TestFunctionExecutedEvent(t *testing.T) {
	jsonStr := `{
		"type": "function_executed",
		"function": {
			"id": "Fn123456789O",
			"callback_id": "sample_function",
			"title": "Sample function",
			"description": "Runs sample function",
			"type": "app",
			"input_parameters": [
				{
					"type": "slack#/types/user_id",
					"name": "user_id",
					"description": "Message recipient",
					"title": "User",
					"is_required": true
				}
			],
			"output_parameters": [
				{
					"type": "slack#/types/user_id",
					"name": "user_id",
					"description": "User that completed the function",
					"title": "Greeting",
					"is_required": true
				}
			],
			"app_id": "AP123456789",
			"date_created": 1694727597,
			"date_updated": 1698947481,
			"date_deleted": 0
		},
		"inputs": { "user_id": "USER12345678" },
		"function_execution_id": "Fx1234567O9L",
		"workflow_execution_id": "WxABC123DEF0",
		"event_ts": "1698958075.998738",
		"bot_access_token": "abcd-1325532282098-1322446258629-6123648410839-527a1cab3979cad288c9e20330d212cf"
	}`

	var event FunctionExecutedEvent
	if err := json.Unmarshal([]byte(jsonStr), &event); err != nil {
		t.Errorf("Failed to unmarshal FunctionExecutedEvent: %v", err)
	}

	if event.Type != "function_executed" {
		t.Errorf("Expected type to be 'function_executed', got %s", event.Type)
	}

	if event.Function.ID != "Fn123456789O" {
		t.Errorf("Expected function.id to be 'Fn123456789O', got %s", event.Function.ID)
	}

	if event.FunctionExecutionID != "Fx1234567O9L" {
		t.Fail()
	}
}

func TestInviteRequestedEvent(t *testing.T) {
	jsonStr := `{
		"type": "invite_requested",
		"invite_request": {
			"id": "12345",
			"email": "bront@puppies.com",
			"date_created": 123455,
			"requester_ids": ["U123ABC456"],
			"channel_ids": ["C123ABC456"],
			"invite_type": "full_member",
			"real_name": "Brent",
			"date_expire": 123456,
			"request_reason": "They're good dogs, Brant",
			"team": {
				"id": "T12345",
				"name": "Puppy ratings workspace incorporated",
				"domain": "puppiesrus"
			}
		}
	}`

	var event InviteRequestedEvent
	if err := json.Unmarshal([]byte(jsonStr), &event); err != nil {
		t.Errorf("Failed to unmarshal InviteRequestedEvent: %v", err)
	}

	if event.Type != "invite_requested" {
		t.Errorf("Expected type to be 'invite_requested', got %s", event.Type)
	}

	if event.InviteRequest.ID != "12345" {
		t.Errorf("invite_request.id should be '12345', was %s", event.InviteRequest.ID)
	}

	if event.InviteRequest.Email != "bront@puppies.com" {
		t.Fail()
	}
}

func TestSharedChannelInviteRequested_UnmarshalJSON(t *testing.T) {
	jsonData := `
	{
		"actor": {
			"id": "U012345ABCD",
			"name": "primary-owner",
			"is_bot": false,
			"team_id": "E0123456ABC",
			"timezone": "",
			"real_name": "primary-owner",
			"display_name": ""
		},
		"channel_id": "C0123ABCDEF",
		"event_type": "slack#/events/shared_channel_invite_requested",
		"channel_name": "our-channel",
		"channel_type": "public",
		"target_users": [
			{
				"email": "user@some-corp.com",
				"invite_id": "I0123456ABC"
			}
		],
		"teams_in_channel": [
			{
				"id": "E0123456ABC",
				"icon": {
					"image_34": "https://slack.com/some-corp/v123/img/abc_0123.png",
					"image_default": true
				},
				"name": "some_enterprise",
				"domain": "someenterprise",
				"is_verified": false,
				"date_created": 1637947110,
				"avatar_base_url": "https://slack.com/some-corp/",
				"requires_sponsorship": false
			},
			{
				"id": "T012345ABCD",
				"icon": {
					"image_34": "https://slack.com/another-corp/v456/img/def_4567.png",
					"image_default": true
				},
				"name": "another_enterprise",
				"domain": "anotherenterprise",
				"is_verified": false,
				"date_created": 1645550933,
				"avatar_base_url": "https://slack.com/another-corp/",
				"requires_sponsorship": false
			}
		],
		"is_external_limited": true,
		"channel_date_created": 1718725442,
		"channel_message_latest_counted_timestamp": 1718745614025449
	}`

	var event SharedChannelInviteRequestedEvent
	err := json.Unmarshal([]byte(jsonData), &event)
	if err != nil {
		t.Fatalf("Failed to unmarshal JSON: %v", err)
	}

	if event.Actor.ID != "U012345ABCD" {
		t.Errorf("Expected Actor.ID to be 'U012345ABCD', got '%s'", event.Actor.ID)
	}
	if event.ChannelID != "C0123ABCDEF" {
		t.Errorf("Expected ChannelID to be 'C0123ABCDEF', got '%s'", event.ChannelID)
	}
	if len(event.TargetUsers) != 1 || event.TargetUsers[0].Email != "user@some-corp.com" {
		t.Errorf("Expected one TargetUser with Email 'user@some-corp.com', got '%v'", event.TargetUsers)
	}
	if len(event.TeamsInChannel) != 2 || event.TeamsInChannel[1].Name != "another_enterprise" {
		t.Errorf("Expected second team to have name 'another_enterprise', got '%v'", event.TeamsInChannel)
	}
}
