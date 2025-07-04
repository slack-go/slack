package main

import (
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/slack-go/slack"
)

func main() {
	channelID := flag.String("channel", "", "Channel ID (required)")
	userID := flag.String("user", "", "User ID to invite (required)")

	flag.Parse()

	// Get token from environment variable
	userToken := os.Getenv("SLACK_USER_TOKEN")
	if userToken == "" {
		fmt.Fprintf(os.Stderr, "SLACK_USER_TOKEN environment variable is required\n")
		os.Exit(1)
	}

	// Get channel ID from flag
	if *channelID == "" {
		fmt.Println("Channel ID is required: use -channel flag")
		os.Exit(1)
	}

	// Get user ID from flag
	if *userID == "" {
		fmt.Println("User ID is required: use -user flag")
		os.Exit(1)
	}

	api := slack.New(userToken)
	_, err := api.InviteUsersToConversation(*channelID, *userID)
	if err != nil {
		var errorResponse slack.SlackErrorResponse
		if errors.As(err, &errorResponse) {
			for _, e := range errorResponse.Errors {
				if e.ConversationsInviteResponseError != nil {
					fmt.Fprintf(os.Stderr, "error inviting user (%s) to conversation: %s\n", e.ConversationsInviteResponseError.User, e.ConversationsInviteResponseError.Error)
				}
			}
		} else {
			fmt.Fprintf(os.Stderr, "error inviting user to conversation: %s\n", err.Error())
		}
		os.Exit(1)
	}
	fmt.Println("User invited successfully to the conversation.")
}
