package main

import (
	"fmt"

	"github.com/essentialkaos/slack"
)

func main() {
	api := slack.New("YOUR_TOKEN_HERE")
	channels, err := api.GetChannels(false)

	if err != nil {
		fmt.Printf("%s\n", err)
		return
	}

	for _, channel := range channels {
		fmt.Println(channel.Name)
		// channel is of type conversation & groupConversation
		// see all available methods in `conversation.go`
	}
}
