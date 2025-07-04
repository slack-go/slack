package main

import (
	"context"
	"fmt"
	"os"

	"github.com/slack-go/slack"
)

func main() {
	// Get token from environment variable
	token := os.Getenv("SLACK_BOT_TOKEN")
	if token == "" {
		fmt.Println("SLACK_BOT_TOKEN environment variable is required")
		os.Exit(1)
	}

	api := slack.New(token)
	params := slack.RemoteFileParameters{
		Title:       "My File",
		ExternalID:  "my-file-123",
		ExternalURL: "https://raw.githubusercontent.com/slack-go/slack/master/README.md",
	}
	file, err := api.AddRemoteFileContext(context.Background(), params)
	if err != nil {
		fmt.Printf("%s\n", err)
		return
	}
	fmt.Printf("Name: %s, URL: %s\n", file.Name, file.URLPrivate)

	err = api.DeleteFileContext(context.Background(), file.ID)
	if err != nil {
		fmt.Printf("%s\n", err)
		return
	}
	fmt.Printf("File %s deleted successfully.\n", file.Name)
}
