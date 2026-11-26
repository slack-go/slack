// Set SLACK_BOT_TOKEN and SLACK_SIGNING_SECRET, configure the app's Interactivity
// Request URL to point to /actions, and run go run ./examples/buttons -channel C123.
package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
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
	signingSecret := os.Getenv("SLACK_SIGNING_SECRET")
	if signingSecret == "" {
		fmt.Fprintln(os.Stderr, "SLACK_SIGNING_SECRET environment variable is required")
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
	http.HandleFunc("POST /actions", func(w http.ResponseWriter, r *http.Request) {
		actionHandler(w, r, signingSecret)
	})
	http.ListenAndServe(":3000", nil)
}

func actionHandler(w http.ResponseWriter, r *http.Request, signingSecret string) {
	verifier, err := slack.NewSecretsVerifier(r.Header, signingSecret)
	if err != nil {
		http.Error(w, "invalid Slack signature", http.StatusUnauthorized)
		return
	}

	// Verify the original bytes before parsing the form.
	body, err := io.ReadAll(io.TeeReader(http.MaxBytesReader(w, r.Body, 1<<20), &verifier))
	if err != nil {
		http.Error(w, "could not read request body", http.StatusBadRequest)
		return
	}
	if err := verifier.Ensure(); err != nil {
		http.Error(w, "invalid Slack signature", http.StatusUnauthorized)
		return
	}

	r.Body = io.NopCloser(bytes.NewReader(body))
	payload, err := slack.InteractionCallbackParse(r)
	actions := payload.ActionCallback.AttachmentActions
	if err != nil || len(actions) == 0 || actions[0] == nil || actions[0].Type != "button" {
		http.Error(w, "expected a button action", http.StatusBadRequest)
		return
	}
	fmt.Printf("Message button pressed by user %s with value %s\n", payload.User.Name, actions[0].Value)
}
