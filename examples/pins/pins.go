package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/slack-go/slack"
)

// This example creates and archives a channel named testpinning if the channel does not exist.
func main() {
	debug := flag.Bool("debug", false, "Show JSON output")
	flag.Parse()

	// Get token from environment variable
	apiToken := os.Getenv("SLACK_BOT_TOKEN")
	if apiToken == "" {
		fmt.Println("SLACK_BOT_TOKEN environment variable is required")
		os.Exit(1)
	}

	api := slack.New(apiToken, slack.OptionDebug(*debug))

	var (
		postAsUserName  string
		postAsUserID    string
		postToChannelID string
	)

	// Find the user to post as.
	authTest, err := api.AuthTest()
	if err != nil {
		fmt.Printf("Error getting channels: %s\n", err)
		return
	}

	channelName := "testpinning"

	// Post as the authenticated user.
	postAsUserName = authTest.User
	postAsUserID = authTest.UserID

	// Create a temporary channel
	channel, err := api.CreateConversation(slack.CreateConversationParams{ChannelName: channelName})

	if err != nil {
		if err.Error() != "name_taken" {
			fmt.Printf("Error creating test channel: %s\n", err)
			return
		}
		channels, err := api.GetAllConversations()
		if err != nil {
			fmt.Printf("Could not retrieve channels: %s\n", err)
			return
		}
		for i := range channels {
			if channels[i].Name == channelName {
				channel = &channels[i]
				break
			}
		}
		if channel == nil {
			fmt.Printf("Channel %q was not found\n", channelName)
			return
		}
		if channel.IsArchived {
			if err := api.UnArchiveConversation(channel.ID); err != nil {
				fmt.Printf("Could not unarchive %s: %s\n", channel.ID, err)
				return
			}
		}
	} else {
		defer func() {
			if err := api.ArchiveConversation(channel.ID); err != nil {
				fmt.Printf("Error archiving test channel: %s\n", err)
			}
		}()
	}
	postToChannelID = channel.ID

	fmt.Printf("Posting as %s (%s) in channel %s\n", postAsUserName, postAsUserID, postToChannelID)

	// Post a message.
	channelID, timestamp, err := api.PostMessage(postToChannelID, slack.MsgOptionText("Is this any good?", false))
	if err != nil {
		fmt.Printf("Error posting message: %s\n", err)
		return
	}

	// Grab a reference to the message.
	msgRef := slack.NewRefToMessage(channelID, timestamp)

	// Add message pin to channel
	if err = api.AddPin(channelID, msgRef); err != nil {
		fmt.Printf("Error adding pin: %s\n", err)
		return
	}

	// List all of the users pins.
	listPins, _, err := api.ListPins(channelID)
	if err != nil {
		fmt.Printf("Error listing pins: %s\n", err)
		return
	}
	fmt.Printf("\n")
	fmt.Printf("All pins by %s...\n", authTest.User)
	for _, item := range listPins {
		fmt.Printf(" > Item type: %s\n", item.Type)
	}

	// Remove the pin.
	err = api.RemovePin(channelID, msgRef)
	if err != nil {
		fmt.Printf("Error remove pin: %s\n", err)
		return
	}
}
