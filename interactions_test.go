package slack

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	dialogSubmissionCallback = `{
		"type": "dialog_submission",
		"submission": {
			"name": "Sigourney Dreamweaver",
			"email": "sigdre@example.com",
			"phone": "+1 800-555-1212",
			"meal": "burrito",
			"comment": "No sour cream please",
			"team_channel": "C0LFFBKPB",
			"who_should_sing": "U0MJRG1AL"
		},
		"callback_id": "employee_offsite_1138b",
		"team": {
			"id": "T1ABCD2E12",
			"domain": "coverbands"
		},
		"user": {
			"id": "W12A3BCDEF",
			"name": "dreamweaver"
		},
		"channel": {
			"id": "C1AB2C3DE",
			"name": "coverthon-1999"
		},
		"action_ts": "936893340.702759",
		"token": "M1AqUUw3FqayAbqNtsGMch72",
		"response_url": "https://hooks.slack.com/app/T012AB0A1/123456789/JpmK0yzoZDeRiqfeduTBYXWQ"
	}`
	actionCallback     = `{}`
	viewClosedCallback = `{
		"type": "view_closed",
		"team": {
			"id": "T1ABCD2E12",
			"domain": "coverbands"
		},
		"user": {
			"id": "W12A3BCDEF",
			"name": "dreamweaver"
		},
		"view": {
			"type": "modal",
			"title": {
				"type": "plain_text",
				"text": "launch project"
			},
			"blocks": [{
				"type": "section",
				"text": {
				  "text": "*Sally* has requested you set the deadline for the Nano launch project",
				  "type": "mrkdwn"
				},
				"accessory": {
				  "type": "datepicker",
				  "action_id": "datepicker123",
				  "initial_date": "1990-04-28",
				  "placeholder": {
					"type": "plain_text",
					"text": "Select a date"
				  }
				}
			}],
			"app_installed_team_id": "T1ABCD2E12"
		},
		"api_app_id": "A123ABC",
		"is_cleared": false
	}`
	viewSubmissionCallback = `{
		"type": "view_submission",
		"team": {
			"id": "T1ABCD2E12",
			"domain": "coverbands"
		},
		"user": {
			"id": "W12A3BCDEF",
			"name": "dreamweaver"
		},
		"channel": {
			"id": "C1AB2C3DE",
			"name": "coverthon-1999"
		},
		"view": {
			"type": "modal",
			"title": {
				"type": "plain_text",
				"text": "meal choice"
			},
			"blocks": [
				{
					"type": "input",
					"block_id": "multi-line",
					"label": {
						"type": "plain_text",
						"text": "dietary restrictions"
					},
					"element": {
						"type": "plain_text_input",
						"multiline": true,
						"action_id": "ml-value"
					}
				},
				{
					"type": "input",
					"block_id": "target_channel",
					"label": {
						"type": "plain_text",
						"text": "Select a channel to post the result on"
					},
					"element": {
						"type": "conversations_select",
						"action_id": "target_select",
						"default_to_current_conversation": true,
						"response_url_enabled": true
					}
				}
			],
			"state": {
				"values": {
					"multi-line": {
						"ml-value": {
							"type": "plain_text_input",
							"value": "No onions"
						}
					},
					"target_channel": {
						"target_select": {
							"type": "conversations_select",
							"value": "C1AB2C3DE"
						}
					}
				}
			},
			"app_installed_team_id": "T1ABCD2E12"
		},
		"hash": "156663117.cd33ad1f",
		"response_urls": [
			{
				"block_id": "target_channel",
				"action_id": "target_select",
				"channel_id": "C1AB2C3DE",
				"response_url": "https:\/\/hooks.slack.com\/app\/ABC12312\/1234567890\/A100B100C100d100"
			}
		]
	}`
)

func assertInteractionCallback(t *testing.T, callback InteractionCallback, encoded string) {
	var decoded InteractionCallback
	assert.Nil(t, json.Unmarshal([]byte(encoded), &decoded))
	assert.Equal(t, decoded, callback)
}

func TestDialogCallback(t *testing.T) {
	expected := InteractionCallback{
		Type:        InteractionTypeDialogSubmission,
		Token:       "M1AqUUw3FqayAbqNtsGMch72",
		CallbackID:  "employee_offsite_1138b",
		ResponseURL: "https://hooks.slack.com/app/T012AB0A1/123456789/JpmK0yzoZDeRiqfeduTBYXWQ",
		ActionTs:    "936893340.702759",
		Team:        Team{ID: "T1ABCD2E12", Name: "", Domain: "coverbands"},
		Channel: Channel{
			GroupConversation: GroupConversation{
				Conversation: Conversation{
					ID: "C1AB2C3DE",
				},
				Name: "coverthon-1999",
			},
		},
		User: User{
			ID:   "W12A3BCDEF",
			Name: "dreamweaver",
		},
		DialogSubmissionCallback: DialogSubmissionCallback{
			Submission: map[string]string{
				"team_channel":    "C0LFFBKPB",
				"who_should_sing": "U0MJRG1AL",
				"name":            "Sigourney Dreamweaver",
				"email":           "sigdre@example.com",
				"phone":           "+1 800-555-1212",
				"meal":            "burrito",
				"comment":         "No sour cream please",
			},
		},
	}
	assertInteractionCallback(t, expected, dialogSubmissionCallback)
}

func TestActionCallback(t *testing.T) {
	assertInteractionCallback(t, InteractionCallback{}, actionCallback)
}

func TestViewClosedck(t *testing.T) {
	expected := InteractionCallback{
		Type: InteractionTypeViewClosed,
		Team: Team{ID: "T1ABCD2E12", Name: "", Domain: "coverbands"},
		User: User{
			ID:   "W12A3BCDEF",
			Name: "dreamweaver",
		},
		View: View{
			Type:  VTModal,
			Title: NewTextBlockObject("plain_text", "launch project", false, false),
			Blocks: Blocks{
				BlockSet: []Block{
					NewSectionBlock(
						NewTextBlockObject("mrkdwn", "*Sally* has requested you set the deadline for the Nano launch project", false, false),
						nil,
						NewAccessory(&DatePickerBlockElement{
							Type:        METDatepicker,
							ActionID:    "datepicker123",
							InitialDate: "1990-04-28",
							Placeholder: NewTextBlockObject("plain_text", "Select a date", false, false),
						}),
					),
				},
			},
			AppInstalledTeamID: "T1ABCD2E12",
		},
		APIAppID: "A123ABC",
	}
	assertInteractionCallback(t, expected, viewClosedCallback)
}

func TestViewSubmissionCallback(t *testing.T) {
	expected := InteractionCallback{
		Type: InteractionTypeViewSubmission,
		Team: Team{ID: "T1ABCD2E12", Name: "", Domain: "coverbands"},
		Channel: Channel{
			GroupConversation: GroupConversation{
				Conversation: Conversation{
					ID: "C1AB2C3DE",
				},
				Name: "coverthon-1999",
			},
		},
		User: User{
			ID:   "W12A3BCDEF",
			Name: "dreamweaver",
		},
		View: View{
			Type:  VTModal,
			Title: NewTextBlockObject("plain_text", "meal choice", false, false),
			Blocks: Blocks{
				BlockSet: []Block{
					NewInputBlock(
						"multi-line",
						NewTextBlockObject(
							"plain_text",
							"dietary restrictions",
							false,
							false,
						),
						nil,
						&PlainTextInputBlockElement{
							Type:      "plain_text_input",
							ActionID:  "ml-value",
							Multiline: true,
						},
					),
					NewInputBlock(
						"target_channel",
						NewTextBlockObject(
							"plain_text",
							"Select a channel to post the result on",
							false,
							false,
						),
						nil,
						&SelectBlockElement{
							Type:                         "conversations_select",
							ActionID:                     "target_select",
							DefaultToCurrentConversation: true,
							ResponseURLEnabled:           true,
						},
					),
				},
			},
			State: &ViewState{
				Values: map[string]map[string]BlockAction{
					"multi-line": {
						"ml-value": {
							Type:  "plain_text_input",
							Value: "No onions",
						},
					},
					"target_channel": {
						"target_select": {
							Type:  "conversations_select",
							Value: "C1AB2C3DE",
						},
					},
				},
			},
			AppInstalledTeamID: "T1ABCD2E12",
		},
		ViewSubmissionCallback: ViewSubmissionCallback{
			Hash: "156663117.cd33ad1f",
			ResponseURLs: []ViewSubmissionCallbackResponseURL{
				{
					BlockID:     "target_channel",
					ActionID:    "target_select",
					ChannelID:   "C1AB2C3DE",
					ResponseURL: "https://hooks.slack.com/app/ABC12312/1234567890/A100B100C100d100",
				},
			},
		},
	}
	assertInteractionCallback(t, expected, viewSubmissionCallback)
}

func TestInteractionCallbackJSONMarshalAndUnmarshal(t *testing.T) {
	cb := &InteractionCallback{
		Type:        InteractionTypeBlockActions,
		Token:       "token",
		CallbackID:  "",
		ResponseURL: "responseURL",
		TriggerID:   "triggerID",
		ActionTs:    "actionTS",
		Team: Team{
			ID:   "teamid",
			Name: "teamname",
		},
		Channel: Channel{
			GroupConversation: GroupConversation{
				Name:         "channelname",
				Conversation: Conversation{ID: "channelid"},
			},
		},
		User: User{
			ID:      "userid",
			Name:    "username",
			Profile: UserProfile{RealName: "userrealname"},
		},
		OriginalMessage: Message{
			Msg: Msg{
				Text:      "ogmsg text",
				Timestamp: "ogmsg ts",
			},
		},
		Message: Message{
			Msg: Msg{
				Text:      "text",
				Timestamp: "ts",
			},
		},
		Name:         "name",
		Value:        "value",
		MessageTs:    "messageTs",
		AttachmentID: "attachmentID",
		ActionCallback: ActionCallbacks{
			AttachmentActions: []*AttachmentAction{
				{Value: "value"},
				{Value: "value2"},
			},
			BlockActions: []*BlockAction{
				{ActionID: "id123"},
				{ActionID: "id456"},
			},
		},
		View: View{
			Type:  VTModal,
			Title: NewTextBlockObject("plain_text", "title", false, false),
			Blocks: Blocks{
				BlockSet: []Block{NewDividerBlock()},
			},
		},
		DialogSubmissionCallback: DialogSubmissionCallback{State: ""},
		RawState:                 json.RawMessage(`{}`),
	}

	cbJSONBytes, err := json.Marshal(cb)
	assert.NoError(t, err)

	jsonCB := new(InteractionCallback)
	err = json.Unmarshal(cbJSONBytes, jsonCB)
	assert.NoError(t, err)

	assert.Equal(t, cb.Type, jsonCB.Type)
	assert.Equal(t, cb.Token, jsonCB.Token)
	assert.Equal(t, cb.CallbackID, jsonCB.CallbackID)
	assert.Equal(t, cb.ResponseURL, jsonCB.ResponseURL)
	assert.Equal(t, cb.TriggerID, jsonCB.TriggerID)
	assert.Equal(t, cb.ActionTs, jsonCB.ActionTs)
	assert.Equal(t, cb.Team.ID, jsonCB.Team.ID)
	assert.Equal(t, cb.Team.Name, jsonCB.Team.Name)
	assert.Equal(t, cb.Channel.ID, jsonCB.Channel.ID)
	assert.Equal(t, cb.Channel.Name, jsonCB.Channel.Name)
	assert.Equal(t, cb.Channel.Created, jsonCB.Channel.Created)
	assert.Equal(t, cb.User.ID, jsonCB.User.ID)
	assert.Equal(t, cb.User.Name, jsonCB.User.Name)
	assert.Equal(t, cb.User.Profile.RealName, jsonCB.User.Profile.RealName)
	assert.Equal(t, cb.OriginalMessage.Text, jsonCB.OriginalMessage.Text)
	assert.Equal(t, cb.OriginalMessage.Timestamp,
		jsonCB.OriginalMessage.Timestamp)
	assert.Equal(t, cb.Message.Text, jsonCB.Message.Text)
	assert.Equal(t, cb.Message.Timestamp, jsonCB.Message.Timestamp)
	assert.Equal(t, cb.Name, jsonCB.Name)
	assert.Equal(t, cb.Value, jsonCB.Value)
	assert.Equal(t, cb.MessageTs, jsonCB.MessageTs)
	assert.Equal(t, cb.AttachmentID, jsonCB.AttachmentID)
	assert.Equal(t, len(cb.ActionCallback.AttachmentActions),
		len(jsonCB.ActionCallback.AttachmentActions))
	assert.Equal(t, len(cb.ActionCallback.BlockActions),
		len(jsonCB.ActionCallback.BlockActions))
	assert.Equal(t, cb.View.Type, jsonCB.View.Type)
	assert.Equal(t, cb.View.Title, jsonCB.View.Title)
	assert.Equal(t, cb.View.Blocks, jsonCB.View.Blocks)
	assert.Equal(t, cb.DialogSubmissionCallback.State,
		jsonCB.DialogSubmissionCallback.State)
}

func TestInteractionCallback_InteractionTypeBlockActions_Unmarshal(t *testing.T) {
	raw := []byte(`{
		"type": "block_actions",
		"actions": [
			{
				"type": "multi_conversations_select",
				"action_id": "multi_convos",
				"block_id": "test123",
				"selected_conversations": ["G12345"]
			}
		],
		"container": {
			"type": "view",
			"view_id": "V12345"
		},
		"state": {
			"values": {
				"section_block_id": {
					"multi_convos": {
						"type": "multi_conversations_select",
						"selected_conversations": ["G12345"]
					}
				},
				"other_block_id": {
					"other_action_id": {
						"type": "plain_text_input",
						"value": "test123"
					}
				}
			}
		}
	}`)
	var cb InteractionCallback
	assert.NoError(t, json.Unmarshal(raw, &cb))
	assert.Equal(t, cb.State, "")
	assert.Equal(t,
		cb.BlockActionState.Values["section_block_id"]["multi_convos"].actionType(),
		ActionType(MultiOptTypeConversations))
	assert.Equal(t,
		cb.BlockActionState.Values["section_block_id"]["multi_convos"].SelectedConversations,
		[]string{"G12345"})
}

func TestInteractionCallback_Container_Marshal_And_Unmarshal(t *testing.T) {
	// Contrived - you generally won't see all of the fields set in a single message
	raw := []byte(
		`
		{
			"container": {
				"type": "message",
				"view_id": "viewID",
				"message_ts": "messageTS",
				"attachment_id": "123",
				"channel_id": "channelID",
				"is_ephemeral": false,
				"is_app_unfurl": false
			}
		}
		`)

	expected := &InteractionCallback{
		Container: Container{
			Type:         "message",
			ViewID:       "viewID",
			MessageTs:    "messageTS",
			AttachmentID: "123",
			ChannelID:    "channelID",
			IsEphemeral:  false,
			IsAppUnfurl:  false,
		},
		RawState: json.RawMessage(`{}`),
	}

	actual := new(InteractionCallback)
	err := json.Unmarshal(raw, actual)
	assert.NoError(t, err)
	assert.Equal(t, expected.Container, actual.Container)

	expectedJSON := []byte(`{"type":"message","view_id":"viewID","message_ts":"messageTS","attachment_id":123,"channel_id":"channelID","is_ephemeral":false,"is_app_unfurl":false}`)
	actualJSON, err := json.Marshal(actual.Container)
	assert.NoError(t, err)
	assert.Equal(t, expectedJSON, actualJSON)
}

func TestInteractionCallback_In_Thread_Container_Marshal_And_Unmarshal(t *testing.T) {
	// Contrived - you generally won't see all of the fields set in a single message
	raw := []byte(
		`
		{
			"container": {
				"type": "message",
				"view_id": "viewID",
				"message_ts": "messageTS",
				"thread_ts": "threadTS",
				"attachment_id": "123",
				"channel_id": "channelID",
				"is_ephemeral": false,
				"is_app_unfurl": false
			}
		}
		`)

	expected := &InteractionCallback{
		Container: Container{
			Type:         "message",
			ViewID:       "viewID",
			MessageTs:    "messageTS",
			ThreadTs:     "threadTS",
			AttachmentID: "123",
			ChannelID:    "channelID",
			IsEphemeral:  false,
			IsAppUnfurl:  false,
		},
		RawState: json.RawMessage(`{}`),
	}

	actual := new(InteractionCallback)
	err := json.Unmarshal(raw, actual)
	assert.NoError(t, err)
	assert.Equal(t, expected.Container, actual.Container)

	expectedJSON := []byte(`{"type":"message","view_id":"viewID","message_ts":"messageTS","thread_ts":"threadTS","attachment_id":123,"channel_id":"channelID","is_ephemeral":false,"is_app_unfurl":false}`)
	actualJSON, err := json.Marshal(actual.Container)
	assert.NoError(t, err)
	assert.Equal(t, expectedJSON, actualJSON)
}
