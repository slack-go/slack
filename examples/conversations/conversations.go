package main

import (
	"fmt"
	"os"

	"github.com/slack-go/slack"
)

func main() {
	botToken := os.Getenv("SLACK_BOT_TOKEN")
	if botToken == "" {
		fmt.Fprintf(os.Stderr, "SLACK_BOT_TOKEN must be set.\n")
		os.Exit(1)
	}

	api := slack.New(botToken)
	params := slack.GetConversationsParameters{
		ExcludeArchived: true,
		Limit:           100,
	}
	channels, _, err := api.GetConversations(&params)
	if err != nil {
		fmt.Printf("%s\n", err)
		return
	}
	for _, channel := range channels {
		fmt.Printf("Channel: %v\n", channel)
	}
}
