package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/slack-go/slack"
)

// WARNING: This example is destructive in the sense that it create a channel called testpinning
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
		channels        []slack.Channel
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
		// If the channel exists, that means we just need to unarchive it
		if err.Error() == "name_taken" {
			err = nil
			params := &slack.GetConversationsParameters{ExcludeArchived: false}
			if channels, _, err = api.GetConversations(params); err != nil {
				fmt.Println("Could not retrieve channels")
				return
			}
			for _, archivedChannel := range channels {
				if archivedChannel.Name == channelName {
					if archivedChannel.IsArchived {
						err = api.UnArchiveConversation(archivedChannel.ID)
						if err != nil {
							fmt.Printf("Could not unarchive %s: %s\n", archivedChannel.ID, err)
							return
						}
					}
					channel = &archivedChannel
					break
				}
			}
		}
		if err != nil {
			fmt.Printf("Error setting test channel for pinning: %s\n", err)
			return
		}
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

	if err = api.UnArchiveConversation(channelID); err != nil {
		fmt.Printf("Error archiving channel: %s\n", err)
		return
	}

}
