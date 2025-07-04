package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/slack-go/slack"
)

func main() {
	channelID := flag.String("channel", "", "Channel ID (required)")
	flag.Parse()

	// Get token from environment variable
	token := os.Getenv("SLACK_BOT_TOKEN")
	if token == "" {
		fmt.Println("SLACK_BOT_TOKEN environment variable is required")
		os.Exit(1)
	}

	// Get channel ID from flag
	if *channelID == "" {
		fmt.Println("Channel ID is required: use -channel flag")
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
	respChannelID, timestamp, err := api.PostMessage(*channelID, slack.MsgOptionText("", false), message)
	if err != nil {
		fmt.Printf("Could not send message: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Message with buttons successfully sent to channel %s at %s", respChannelID, timestamp)
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
