package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/slack-go/slack"
)

func main() {
	debug := flag.Bool("debug", false, "Show JSON output")
	channelID := flag.String("channel", "", "Channel ID (required)")
	flag.Parse()

	// Get token from environment variable
	apiToken := os.Getenv("SLACK_USER_TOKEN")
	if apiToken == "" {
		fmt.Println("SLACK_USER_TOKEN environment variable is required")
		os.Exit(1)
	}

	api := slack.New(apiToken, slack.OptionDebug(*debug))

	var (
		postAsUserName string
		postAsUserID   string
	)

	// Find the user to post as.
	authTest, err := api.AuthTest()
	if err != nil {
		fmt.Printf("Error getting channels: %s\n", err)
		return
	}

	// Post as the authenticated user.
	postAsUserName = authTest.User
	postAsUserID = authTest.UserID

	fmt.Printf("Posting as %s (%s) in channel %s\n", postAsUserName, postAsUserID, *channelID)

	// Post a message.
	_, timestamp, err := api.PostMessage(*channelID, slack.MsgOptionText("Is this any good?", false))
	if err != nil {
		fmt.Printf("Error posting message: %s\n", err)
		return
	}

	// // Grab a reference to the message.
	msgRef := slack.NewRefToMessage(*channelID, timestamp)

	fmt.Printf("Adding reaction to message with reference %v\n", msgRef)

	// React with :+1:
	if err = api.AddReaction("+1", msgRef); err != nil {
		fmt.Printf("Error adding reaction: %s\n", err)
		return
	}

	// React with :-1:
	if err = api.AddReaction("cry", msgRef); err != nil {
		fmt.Printf("Error adding reaction: %s\n", err)
		return
	}

	// Get all reactions on the message.
	msgReactionsResp, err := api.GetReactions(msgRef, slack.NewGetReactionsParameters())
	if err != nil {
		fmt.Printf("Error getting reactions: %s\n", err)
		return
	}
	fmt.Printf("\n")
	fmt.Printf("%d reactions to message...\n", len(msgReactionsResp.Reactions))
	for _, r := range msgReactionsResp.Reactions {
		fmt.Printf("  %d users say %s in channel %s\n", r.Count, r.Name, msgReactionsResp.Item.Channel)
	}

	// List all of the users reactions.
	listParams := slack.NewListReactionsParameters()
	fmt.Printf("Listing reactions with params: User=%q, TeamID=%q, Count=%d, Page=%d, Full=%v\n",
		listParams.User, listParams.TeamID, listParams.Count, listParams.Page, listParams.Full)
	listReactions, _, err := api.ListReactions(listParams)
	if err != nil {
		fmt.Printf("Error listing reactions: %v\n", err)
		if slackErr, ok := err.(slack.SlackErrorResponse); ok {
			fmt.Printf("  ResponseMetadata.Messages: %v\n", slackErr.ResponseMetadata.Messages)
		}
		return
	}
	fmt.Printf("\n")
	fmt.Printf("All reactions by %s...\n", authTest.User)
	for _, item := range listReactions {
		fmt.Printf("%d on a %s...\n", len(item.Reactions), item.Type)
		for _, r := range item.Reactions {
			fmt.Printf("  %s (along with %d others)\n", r.Name, r.Count-1)
		}
	}

	// Remove the :cry: reaction.
	err = api.RemoveReaction("cry", msgRef)
	if err != nil {
		fmt.Printf("Error remove reaction: %s\n", err)
		return
	}

	// Get all reactions on the message.
	msgReactionsResp, err = api.GetReactions(msgRef, slack.NewGetReactionsParameters())
	if err != nil {
		fmt.Printf("Error getting reactions: %s\n", err)
		return
	}
	fmt.Printf("\n")
	fmt.Printf("%d reactions to message after removing cry...\n", len(msgReactionsResp.Reactions))
	for _, r := range msgReactionsResp.Reactions {
		fmt.Printf("  %d users say %s\n", r.Count, r.Name)
	}
}
