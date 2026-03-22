package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/slack-go/slack"
)

func main() {
	action := flag.String("action", "", "Action to perform: list, add, remove, set (required)")
	channel := flag.String("channel", "", "Channel ID (required for add, remove, set)")
	channels := flag.String("channels", "", "Comma-separated channel IDs (required for list)")
	triggers := flag.String("triggers", "", "Comma-separated trigger IDs (required for add, remove, set)")
	flag.Parse()

	token := os.Getenv("SLACK_BOT_TOKEN")
	if token == "" {
		fmt.Fprintf(os.Stderr, "SLACK_BOT_TOKEN environment variable is required\n")
		os.Exit(1)
	}

	if *action == "" {
		fmt.Fprintf(os.Stderr, "Error: -action flag is required (list, add, remove, set)\n")
		os.Exit(1)
	}

	api := slack.New(token)

	switch *action {
	case "list":
		listFeatured(api, *channels)
	case "add":
		addFeatured(api, *channel, *triggers)
	case "remove":
		removeFeatured(api, *channel, *triggers)
	case "set":
		setFeatured(api, *channel, *triggers)
	default:
		fmt.Fprintf(os.Stderr, "Error: unknown action %q (must be list, add, remove, set)\n", *action)
		os.Exit(1)
	}
}

func splitCSV(s string) []string {
	if s == "" {
		return nil
	}
	parts := strings.Split(s, ",")
	for i := range parts {
		parts[i] = strings.TrimSpace(parts[i])
	}
	return parts
}

func listFeatured(api *slack.Client, channelsFlag string) {
	channelIDs := splitCSV(channelsFlag)
	if len(channelIDs) == 0 {
		fmt.Fprintf(os.Stderr, "Error: -channels flag is required for list action\n")
		os.Exit(1)
	}

	resp, err := api.WorkflowsFeaturedList(context.Background(), &slack.WorkflowsFeaturedListInput{
		ChannelIDs: channelIDs,
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error listing featured workflows: %s\n", err)
		os.Exit(1)
	}

	for _, fw := range resp.FeaturedWorkflows {
		fmt.Printf("Channel %s:\n", fw.ChannelID)
		if len(fw.Triggers) == 0 {
			fmt.Println("  (no featured workflows)")
			continue
		}
		for _, t := range fw.Triggers {
			fmt.Printf("  - %s (ID: %s)\n", t.Title, t.ID)
		}
	}
}

func addFeatured(api *slack.Client, channelID, triggersFlag string) {
	if channelID == "" {
		fmt.Fprintf(os.Stderr, "Error: -channel flag is required for add action\n")
		os.Exit(1)
	}
	triggerIDs := splitCSV(triggersFlag)
	if len(triggerIDs) == 0 {
		fmt.Fprintf(os.Stderr, "Error: -triggers flag is required for add action\n")
		os.Exit(1)
	}

	err := api.WorkflowsFeaturedAdd(context.Background(), &slack.WorkflowsFeaturedAddInput{
		ChannelID:  channelID,
		TriggerIDs: triggerIDs,
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error adding featured workflows: %s\n", err)
		os.Exit(1)
	}

	fmt.Println("Featured workflows added successfully")
}

func removeFeatured(api *slack.Client, channelID, triggersFlag string) {
	if channelID == "" {
		fmt.Fprintf(os.Stderr, "Error: -channel flag is required for remove action\n")
		os.Exit(1)
	}
	triggerIDs := splitCSV(triggersFlag)
	if len(triggerIDs) == 0 {
		fmt.Fprintf(os.Stderr, "Error: -triggers flag is required for remove action\n")
		os.Exit(1)
	}

	err := api.WorkflowsFeaturedRemove(context.Background(), &slack.WorkflowsFeaturedRemoveInput{
		ChannelID:  channelID,
		TriggerIDs: triggerIDs,
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error removing featured workflows: %s\n", err)
		os.Exit(1)
	}

	fmt.Println("Featured workflows removed successfully")
}

func setFeatured(api *slack.Client, channelID, triggersFlag string) {
	if channelID == "" {
		fmt.Fprintf(os.Stderr, "Error: -channel flag is required for set action\n")
		os.Exit(1)
	}
	triggerIDs := splitCSV(triggersFlag)
	if len(triggerIDs) == 0 {
		fmt.Fprintf(os.Stderr, "Error: -triggers flag is required for set action\n")
		os.Exit(1)
	}

	err := api.WorkflowsFeaturedSet(context.Background(), &slack.WorkflowsFeaturedSetInput{
		ChannelID:  channelID,
		TriggerIDs: triggerIDs,
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error setting featured workflows: %s\n", err)
		os.Exit(1)
	}

	fmt.Println("Featured workflows set successfully")
}
