package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/slack-go/slack"
)

// This example demonstrates how to use the assistant.search.context API
// to search across a Slack workspace. This API is designed for AI/LLM
// consumption and returns messages, files, and channels matching a query.
//
// Usage: go run assistant_search.go <query>
//
// This example uses a user token (xoxp-...), which can call the API directly.
// Bot tokens (xoxb-...) require an action_token received from message events;
// see the ai_apps example for that pattern.
//
// Required scopes (user token): search:read.public, search:read.private,
//
//	search:read.im, search:read.mpim, search:read.files, search:read.users
//
// See https://docs.slack.dev/reference/methods/assistant.search.context
func main() {
	token := os.Getenv("SLACK_USER_TOKEN")
	if token == "" {
		fmt.Println("SLACK_USER_TOKEN environment variable is required (xoxp-... token)")
		os.Exit(1)
	}

	if len(os.Args) < 2 {
		fmt.Println("Usage: go run assistant_search.go <query>")
		os.Exit(1)
	}

	query := strings.Join(os.Args[1:], " ")
	api := slack.New(token)

	// Search across public channels for messages, files, and channels
	resp, err := api.SearchAssistantContext(slack.AssistantSearchContextParameters{
		Query:        query,
		ChannelTypes: []string{"public_channel"},
		ContentTypes: []string{"messages", "files", "channels"},
		Limit:        10,
	})
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		os.Exit(1)
	}

	printResults(resp)

	// Paginate using cursor
	if resp.ResponseMetadata.NextCursor != "" {
		fmt.Printf("\nNext page cursor: %s\n", resp.ResponseMetadata.NextCursor)
	}

	// Advanced search: keyword-only, sorted by time, with context messages
	advanced, err := api.SearchAssistantContext(slack.AssistantSearchContextParameters{
		Query:                  query,
		ChannelTypes:           []string{"public_channel", "private_channel"},
		ContentTypes:           []string{"messages", "files"},
		Sort:                   "timestamp",
		SortDir:                "desc",
		IncludeContextMessages: true,
		Highlight:              true,
		DisableSemanticSearch:  true,
		Limit:                  5,
	})
	if err != nil {
		fmt.Printf("Advanced search error: %s\n", err)
		os.Exit(1)
	}

	fmt.Println("\n--- Advanced Search (keyword-only, with context) ---")
	printResults(advanced)
}

func printResults(resp *slack.AssistantSearchContextResponse) {
	fmt.Printf("=== Messages (%d) ===\n", len(resp.Results.Messages))
	for _, msg := range resp.Results.Messages {
		bot := ""
		if msg.IsAuthorBot {
			bot = " [bot]"
		}
		fmt.Printf("  %s (%s)%s in #%s:\n    %s\n    %s\n",
			msg.AuthorName, msg.AuthorUserID, bot,
			msg.ChannelName, msg.Content, msg.Permalink)

		if msg.ContextMessages != nil {
			for _, before := range msg.ContextMessages.Before {
				fmt.Printf("    [before] %s: %s\n", before.AuthorUserID, before.Content)
			}
			for _, after := range msg.ContextMessages.After {
				fmt.Printf("    [after]  %s: %s\n", after.AuthorUserID, after.Content)
			}
		}
	}

	fmt.Printf("\n=== Files (%d) ===\n", len(resp.Results.Files))
	for _, f := range resp.Results.Files {
		fmt.Printf("  %s (%s) by %s\n    %s\n",
			f.Title, f.FileType, f.AuthorName, f.Permalink)
	}

	fmt.Printf("\n=== Channels (%d) ===\n", len(resp.Results.Channels))
	for _, ch := range resp.Results.Channels {
		fmt.Printf("  #%s — %s\n    %s\n",
			ch.Name, ch.Purpose, ch.Permalink)
	}
}
