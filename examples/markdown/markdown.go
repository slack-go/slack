package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/slack-go/slack"
)

func main() {
	channelID := flag.String("channel", "", "Channel ID (required)")
	flag.Parse()

	token := os.Getenv("SLACK_BOT_TOKEN")
	if token == "" {
		fmt.Println("SLACK_BOT_TOKEN environment variable is required")
		os.Exit(1)
	}

	if *channelID == "" {
		fmt.Println("Channel ID is required: use -channel flag")
		os.Exit(1)
	}

	api := slack.New(token)

	// Slack uses its own markdown-like syntax called mrkdwn.
	// See https://api.slack.com/reference/surfaces/formatting for the full spec.
	mrkdwnText := `*Bold text* and _italic text_ and ~strikethrough~

Inline ` + "`code`" + ` and a code block:
` + "```" + `
func main() {
    fmt.Println("Hello, Slack!")
}
` + "```" + `

A link: <https://api.slack.com|Slack API docs>

> A blockquote for emphasis

And a list:
• First item
• Second item
• Third item`

	section := slack.NewSectionBlock(
		slack.NewTextBlockObject("mrkdwn", mrkdwnText, false, false),
		nil,
		nil,
	)

	respChannelID, timestamp, err := api.PostMessage(
		*channelID,
		slack.MsgOptionText("Markdown formatting example (fallback text)", false),
		slack.MsgOptionBlocks(section),
	)
	if err != nil {
		fmt.Printf("Error sending message: %s\n", err)
		os.Exit(1)
	}

	fmt.Printf("Message successfully sent to channel %s at %s\n", respChannelID, timestamp)
}
