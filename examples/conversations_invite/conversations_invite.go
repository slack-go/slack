package main

import (
	"errors"
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
	_, err := api.InviteUsersToConversation(
		"C1234567890", // Replace with the actual channel ID you want to invite the user to
		"U1234567890", // Replace with the actual user ID you want to invite
	)
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
