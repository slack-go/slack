package main

import (
	"context"
	"flag"
	"fmt"
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
	params := slack.GetConversationHistoryParameters{
		ChannelID: *channelID,
	}
	messages, err := api.GetConversationHistoryContext(context.Background(), &params)
	if err != nil {
		fmt.Printf("%s\n", err)
		return
	}
	for _, message := range messages.Messages {
		if len(message.Attachments) > 0 {
			fmt.Printf("Message: %s\n", message.Attachments[0].Color)
		} else {
			fmt.Printf("Message: %s\n", message.Text)
		}
	}
}
