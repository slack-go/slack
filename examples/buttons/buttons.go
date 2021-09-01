package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/slack-go/slack"
)

func main() {
	var token, channel string
	var ok bool
	token, ok = os.LookupEnv("SLACK_TOKEN")
	if !ok {
		fmt.Println("Missing SLACK_TOKEN in environment")
		os.Exit(1)
	}
	channel, ok = os.LookupEnv("SLACK_CHANNEL")
	if !ok {
		fmt.Println("Missing SLACK_CHANNEL in environment")
		os.Exit(1)
	}
	api := slack.New(token)
	attachment := slack.Attachment{
		Pretext:    "pretext",
		Fallback:   "We don't currently support your client",
		CallbackID: "accept_or_reject",
		Color:      "#3AA3E3",
		Actions: []slack.AttachmentAction{
			slack.AttachmentAction{
				Name:  "accept",
				Text:  "Accept",
				Type:  "button",
				Value: "accept",
			},
			slack.AttachmentAction{
				Name:  "reject",
				Text:  "Reject",
				Type:  "button",
				Value: "reject",
				Style: "danger",
			},
		},
	}

	message := slack.MsgOptionAttachments(attachment)
	channelID, timestamp, err := api.PostMessage(channel, slack.MsgOptionText("", false), message)
	if err != nil {
		fmt.Printf("Could not send message: %v", err)
	}
	fmt.Printf("Message with buttons sucessfully sent to channel %s at %s", channelID, timestamp)
	http.HandleFunc("/actions", actionHandler)
	http.ListenAndServe(":3000", nil)
}

func actionHandler(w http.ResponseWriter, r *http.Request) {
	var payload slack.InteractionCallback
	err := json.Unmarshal([]byte(r.FormValue("payload")), &payload)
	if err != nil {
		fmt.Printf("Could not parse action response JSON: %v", err)
	}
	fmt.Printf("Message button pressed by user %s with value %s", payload.User.Name, payload.ActionCallback.AttachmentActions[0].Value)
}
