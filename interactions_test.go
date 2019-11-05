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
	actionCallback = `{}`
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

func TestInteractionCallbackJSONMarshalAndUnmarshal(t *testing.T) {
	cb := &InteractionCallback{
		Type:        InteractionTypeBlockActions,
		Token:       "token",
		CallbackID:  "",
		ResponseURL: "responseURL",
		TriggerID:   "triggerID",
		ActionTs:    "actionTS",
		Team:        Team{ID: "teamid", Name: "teamname"},
		Channel: Channel{GroupConversation: GroupConversation{
			Name: "channelname", Conversation: Conversation{ID: "channelid"}}},
		User: User{ID: "userid", Name: "username",
			Profile: UserProfile{RealName: "userrealname"}},
		OriginalMessage: Message{Msg: Msg{Text: "ogmsg text",
			Timestamp: "ogmsg ts"}},
		Message:      Message{Msg: Msg{Text: "text", Timestamp: "ts"}},
		Name:         "name",
		Value:        "value",
		MessageTs:    "messageTs",
		AttachmentID: "attachmentID",
		ActionCallback: ActionCallbacks{
			AttachmentActions: []*AttachmentAction{{Value: "value"}},
			BlockActions:      []*BlockAction{{ActionID: "id123"}},
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
	assert.Equal(t, cb.DialogSubmissionCallback.State,
		jsonCB.DialogSubmissionCallback.State)
}
