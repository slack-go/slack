package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/slack-go/slack"
)

func main() {
	var (
		debug bool
		team  string
	)

	// Get token from environment variable
	apiToken := os.Getenv("SLACK_USER_TOKEN")
	if apiToken == "" {
		fmt.Println("SLACK_USER_TOKEN environment variable is required")
		os.Exit(1)
	}

	flag.BoolVar(&debug, "debug", false, "Show JSON output")
	flag.StringVar(&team, "team", "", "Team ID (required for Enterprise Grid)")
	flag.Parse()

	api := slack.New(apiToken, slack.OptionDebug(debug))

	// Get all stars for the user.
	params := slack.NewStarsParameters()
	params.TeamID = team

	starredItems, _, err := api.GetStarred(params)
	if err != nil {
		fmt.Printf("Error getting stars: %v\n", err)
		return
	}
	for _, s := range starredItems {
		var desc string
		switch s.Type {
		case slack.TYPE_MESSAGE:
			desc = s.Message.Text
		case slack.TYPE_FILE:
			desc = s.File.Name
		case slack.TYPE_FILE_COMMENT:
			desc = s.File.Name + " - " + s.Comment.Comment
		case slack.TYPE_CHANNEL, slack.TYPE_IM, slack.TYPE_GROUP:
			desc = s.Channel
		}
		fmt.Printf("Starred %s: %s\n", s.Type, desc)
	}
}
