package main

import (
	"fmt"
	"os"

	"github.com/slack-go/slack"
)

func main() {
	userToken := os.Getenv("SLACK_USER_TOKEN")
	if userToken == "" {
		fmt.Fprintf(os.Stderr, "SLACK_USER_TOKEN must be set.\n")
		os.Exit(1)
	}

	api := slack.New(userToken)
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
		info, err := api.GetConversationInfo(&slack.GetConversationInfoInput{
			ChannelID:         channel.ID,
			IncludeNumMembers: true,
			IncludeLocale:     true,
		})
		if err != nil {
			fmt.Printf("Error getting info for channel %s: %s\n", channel.ID, err)
			continue
		}
		fmt.Printf("Channel: %s\n", channel.ID)
		if info.Properties != nil {
			fmt.Printf("Canvas: %+v\n", info.Properties.Canvas)
			fmt.Printf("Tabs: %+v\n", info.Properties.Tabs)
		}
	}
}
