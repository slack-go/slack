package main

import (
	"context"
	"fmt"

	"github.com/slack-go/slack"
)

func main() {
	api := slack.New("YOUR_TOKEN")
	params := slack.GetConversationHistoryParameters{
		ChannelID: "C0123456789",
	}
	messages, err := api.GetConversationHistoryContext(context.Background(), &params)
	if err != nil {
		fmt.Printf("%s\n", err)
		return
	}
	for _, message := range messages.Messages {
		fmt.Printf("Message: %s\n", message.Attachments[0].Color)
	}
}
