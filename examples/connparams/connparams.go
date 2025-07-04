package main

import (
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"

	"github.com/slack-go/slack"
)

func main() {
	channelID := flag.String("channel", "", "Channel ID (required)")
	userIDs := flag.String("users", "", "Comma-separated user IDs for presence monitoring (required)")
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

	// Get user IDs from flag
	if *userIDs == "" {
		fmt.Println("User IDs are required: use -users flag (comma-separated)")
		os.Exit(1)
	}

	// Parse comma-separated user IDs
	userIDList := strings.Split(*userIDs, ",")
	for i, userID := range userIDList {
		userIDList[i] = strings.TrimSpace(userID)
	}

	api := slack.New(
		token,
		slack.OptionDebug(true),
		slack.OptionLog(log.New(os.Stdout, "slack-bot: ", log.Lshortfile|log.LstdFlags)),
	)

	// turn on the batch_presence_aware option
	rtm := api.NewRTM(slack.RTMOptionConnParams(url.Values{
		"batch_presence_aware": {"1"},
	}))
	go rtm.ManageConnection()

	for msg := range rtm.IncomingEvents {
		fmt.Print("Event Received: ")
		switch ev := msg.Data.(type) {
		case *slack.HelloEvent:
			// Subscribe to user presence using provided user IDs
			rtm.SendMessage(rtm.NewSubscribeUserPresence(userIDList))

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
