package main

import (
	"flag"
	"fmt"
	"log"
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
	api := slack.New(
		token,
		slack.OptionDebug(true),
		slack.OptionLog(log.New(os.Stdout, "slack-bot: ", log.Lshortfile|log.LstdFlags)),
	)

	rtm := api.NewRTM()
	go rtm.ManageConnection()

	for msg := range rtm.IncomingEvents {
		fmt.Print("Event Received: ")
		switch ev := msg.Data.(type) {
		case *slack.HelloEvent:
			// Ignore hello

		case *slack.ConnectedEvent:
			fmt.Println("Infos:", ev.Info)
			fmt.Println("Connection counter:", ev.ConnectionCount)
			// Send message to provided channel ID
			rtm.SendMessage(rtm.NewOutgoingMessage("Hello world", *channelID))

		case *slack.MessageEvent:
			fmt.Printf("Message: %v\n", ev)

		case *slack.PresenceChangeEvent:
			fmt.Printf("Presence Change: %v\n", ev)

		case *slack.LatencyReport:
			fmt.Printf("Current latency: %v\n", ev.Value)

		case *slack.DesktopNotificationEvent:
			fmt.Printf("Desktop Notification: %v\n", ev)

		case *slack.RTMError:
			fmt.Printf("Error: %s\n", ev.Error())

		case *slack.InvalidAuthEvent:
			fmt.Printf("Invalid credentials")
			return

		default:

			// Ignore other events..
			// fmt.Printf("Unexpected: %v\n", msg.Data)
		}
	}
}
