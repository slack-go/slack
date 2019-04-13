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
