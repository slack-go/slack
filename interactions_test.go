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
			}]
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
			"blocks": [{
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
			}],
			"state": {
				"values": {
					"multi-line": {
						"ml-value": {
							"type": "plain_text_input",
							"value": "No onions"
						}
					}
				}
			}
		}
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
						&PlainTextInputBlockElement{
							Type:      "plain_text_input",
							ActionID:  "ml-value",
							Multiline: true,
						},
					),
				},
			},
			State: &ViewState{
				Values: map[string]map[string]BlockAction{
					"multi-line": map[string]BlockAction{
						"ml-value": BlockAction{
							Type:  "plain_text_input",
							Value: "No onions",
						},
					},
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
			AttachmentActions: []*AttachmentAction{{Value: "value"}},
			BlockActions:      []*BlockAction{{ActionID: "id123"}},
		},
		View: View{
			Type:  VTModal,
			Title: NewTextBlockObject("plain_text", "title", false, false),
			Blocks: Blocks{
				BlockSet: []Block{NewDividerBlock()},
			},
		},
		DialogSubmissionCallback: DialogSubmissionCallback{State: "dsstate"},
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
