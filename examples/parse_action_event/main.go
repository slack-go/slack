// This example shows how to handle both Events API and Interactions API
// using HTTP endpoints. It listens for messages (Events API), replies with
// a button (Block Kit), and handles the button click (Interactions API).
//
// This is also a migration guide for users of slackevents.ParseActionEvent,
// which is deprecated because it cannot parse block_actions payloads. The
// correct approach is to use slack.InteractionCallback directly:
//
//	// Before (broken for block_actions):
//	action, err := slackevents.ParseActionEvent(payload, slackevents.OptionNoVerifyToken())
//
//	// After (handles all interaction types):
//	var ic slack.InteractionCallback
//	err := json.Unmarshal([]byte(payload), &ic)
//	// Use ic.ActionCallback.BlockActions for block actions
//	// Use ic.ActionCallback.AttachmentActions for legacy attachment actions
//
// Note: block_actions are delivered to your Interactivity Request URL, not
// your Events API Request URL. These are two separate endpoints in your
// Slack app configuration.
//
// Setup:
//  1. export SLACK_BOT_TOKEN=xoxb-...
//  2. export SLACK_SIGNING_SECRET=...
//  3. go run ./examples/parse_action_event
//  4. Expose port 3000 with ngrok: ngrok http 3000
//  5. In your Slack app config:
//     - Events API Request URL:    https://<ngrok>/events
//     - Interactivity Request URL: https://<ngrok>/interactions
//     - Subscribe to bot events: message.channels (so the bot can post)
//     - Bot scopes: chat:write, channels:read
//  6. Invite the bot to a channel, then send any message — the bot will
//     reply with a message containing a block action button.
//  7. Click the button and watch the logs.
package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"

	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
)

func main() {
	botToken := os.Getenv("SLACK_BOT_TOKEN")
	signingSecret := os.Getenv("SLACK_SIGNING_SECRET")
	if botToken == "" || signingSecret == "" {
		fmt.Fprintln(os.Stderr, "Set SLACK_BOT_TOKEN and SLACK_SIGNING_SECRET")
		os.Exit(1)
	}

	api := slack.New(botToken)

	// Events API endpoint — receives subscription events (message, app_mention, etc.)
	http.HandleFunc("/events", func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		sv, err := slack.NewSecretsVerifier(r.Header, signingSecret)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if _, err := sv.Write(body); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if err := sv.Ensure(); err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		eventsAPIEvent, err := slackevents.ParseEvent(json.RawMessage(body), slackevents.OptionNoVerifyToken())
		if err != nil {
			fmt.Printf("[EVENTS] parse error: %v\n", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		fmt.Printf("[EVENTS] received type=%q\n", eventsAPIEvent.Type)

		switch eventsAPIEvent.Type {
		case slackevents.URLVerification:
			var cr *slackevents.ChallengeResponse
			if err := json.Unmarshal(body, &cr); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "text/plain")
			w.Write([]byte(cr.Challenge))

		case slackevents.CallbackEvent:
			innerEvent := eventsAPIEvent.InnerEvent
			fmt.Printf("[EVENTS] inner type=%q\n", innerEvent.Type)

			switch ev := innerEvent.Data.(type) {
			case *slackevents.MessageEvent:
				// Ignore bot messages to avoid loops.
				if ev.BotID != "" {
					return
				}
				fmt.Printf("[EVENTS] message from user %s: %q\n", ev.User, ev.Text)

				// Reply with a message containing a block action button.
				_, _, err := api.PostMessage(ev.Channel, slack.MsgOptionBlocks(
					slack.NewSectionBlock(
						slack.NewTextBlockObject("mrkdwn", "Click the button to test block_actions delivery:", false, false),
						nil, nil,
					),
					slack.NewActionBlock("test_actions_block",
						slack.NewButtonBlockElement("test_button", "clicked",
							slack.NewTextBlockObject("plain_text", "Click me", false, false),
						),
					),
				))
				if err != nil {
					fmt.Printf("[EVENTS] PostMessage error: %v\n", err)
				}
			}
		}
	})

	// Interactions endpoint — receives interactive callbacks (block_actions,
	// interactive_message, view_submission, etc.).
	//
	// If you were previously using slackevents.ParseActionEvent, this is the
	// correct replacement: unmarshal into slack.InteractionCallback and switch
	// on ic.Type.
	http.HandleFunc("/interactions", func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		sv, err := slack.NewSecretsVerifier(r.Header, signingSecret)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if _, err := sv.Write(body); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if err := sv.Ensure(); err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		// Interactions come as form-encoded with a "payload" field.
		payload, err := url.QueryUnescape(string(body))
		if err != nil {
			fmt.Printf("[INTERACTIONS] unescape error: %v\n", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		// Strip the "payload=" prefix.
		const prefix = "payload="
		if len(payload) > len(prefix) {
			payload = payload[len(prefix):]
		}

		var ic slack.InteractionCallback
		if err := json.Unmarshal([]byte(payload), &ic); err != nil {
			fmt.Printf("[INTERACTIONS] parse error: %v\n", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		fmt.Printf("[INTERACTIONS] type=%q callback_id=%q\n", ic.Type, ic.CallbackID)
		switch ic.Type {
		case slack.InteractionTypeBlockActions:
			for i, a := range ic.ActionCallback.BlockActions {
				fmt.Printf("[INTERACTIONS] block_action[%d]: action_id=%q block_id=%q type=%q value=%q\n",
					i, a.ActionID, a.BlockID, a.Type, a.Value)
			}
		case slack.InteractionTypeInteractionMessage:
			for i, a := range ic.ActionCallback.AttachmentActions {
				fmt.Printf("[INTERACTIONS] attachment_action[%d]: name=%q type=%q value=%q\n",
					i, a.Name, a.Type, a.Value)
			}
		}
	})

	fmt.Println("[INFO] Listening on :3000")
	fmt.Println("[INFO] Events API:    http://localhost:3000/events")
	fmt.Println("[INFO] Interactions:  http://localhost:3000/interactions")
	http.ListenAndServe(":3000", nil)
}
