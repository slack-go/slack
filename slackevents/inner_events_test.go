package slackevents

import (
	"encoding/json"
	"testing"
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

func TestTeamDomainChangeEvent(t *testing.T) {
	rawE := []byte(`
	{
		"type": "team_domain_change",
		"url": "https://my.slack.com",
		"domain": "my"
		"team_id": "T1234"
	}	
	`)
	err := json.Unmarshal(rawE, &TeamDomainChangeEvent{})
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
	err := json.Unmarshal(rawE, &MemberJoinedChannelEvent{})
	if err != nil {
		t.Error(err)
	}
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

func TestSharedChannelReceivedEvent(t *testing.T) {
	rawE := []byte(`
	{
		"type": "shared_channel_invite_received",
		"invite": {
		  "id": "I028YDERZSQ",
		  "date_created": 1626876000,
		  "date_invalid": 1628085600,
		  "inviting_team": {
			"id": "T12345678",
			"name": "Corgis",
			"icon": {...},
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
	`)
	err := json.Unmarshal(rawE, &SharedChannelInviteReceivedEvent{})
	if err != nil {
		t.Error(err)
	}
}

func TestSharedChannelAcceptedEvent(t *testing.T) {
	rawE := []byte(`
	{
		"type": "shared_channel_invite_accepted",
		"approval_required": false,
		"invite": {
		  "id": "I028YDERZSQ",
		  "date_created": 1626876000,
		  "date_invalid": 1628085600,
		  "inviting_team": {
			"id": "T12345678",
			"name": "Corgis",
			"icon": {...},
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
			"icon": {...},
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
	`)
	err := json.Unmarshal(rawE, &SharedChannelInviteAcceptedEvent{})
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

func TestChannelShared(t *testing.T) {
	rawE := []byte(`
	{    
		"type": "channel_shared",
		"connected_team_id": "E163Q94DX",
		"channel": "C123ABC456",
		"event_ts": "1561064063.001100"
	}
`)
	err := json.Unmarshal(rawE, &ChannelSharedEvent{})
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
