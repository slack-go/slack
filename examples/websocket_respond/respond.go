package main

import (
	"fmt"
	"strings"

	"github.com/slack-go/slack"
)

func main() {
	api := slack.New(
		"YOUR-TOKEN-HERE",
	)

	rtm := api.NewRTM()
	go rtm.ManageConnection()

	for msg := range rtm.IncomingEvents {
		switch ev := msg.Data.(type) {
		case *slack.MessageEvent:
			msg := ev.Msg

			if msg.SubType != "" {
				break // We're only handling normal messages.
			}

			// Create a response object.
			resp := rtm.NewOutgoingMessage(fmt.Sprintf("echo %s", msg.Text), msg.Channel)

			// Respond in thread if not a direct message.
			if !strings.HasPrefix(msg.Channel, "D") {
				resp.ThreadTimestamp = msg.Timestamp
			}

			// Respond in same thread if message came from a thread.
			if msg.ThreadTimestamp != "" {
				resp.ThreadTimestamp = msg.ThreadTimestamp
			}

			rtm.SendMessage(resp)

		case *slack.ConnectedEvent:
			fmt.Println("Connected to Slack")

		case *slack.InvalidAuthEvent:
			fmt.Println("Invalid token")
			return
		}
	}
}
