package main

import (
	"context"
	"fmt"
	"os"

	"github.com/slack-go/slack"
)

func main() {
	token := os.Getenv("SLACK_BOT_TOKEN")
	if token == "" {
		fmt.Println("SLACK_BOT_TOKEN environment variable is required")
		os.Exit(1)
	}
	api := slack.New(token, slack.OptionDebug(true))

	ctx := context.Background()

	// Upload a file
	params := slack.UploadFileParameters{
		Title:    "Batman Example",
		Filename: "example.txt",
		File:     "example.txt",
		FileSize: 38,
	}
	file, err := api.UploadFileContext(ctx, params)
	if err != nil {
		fmt.Printf("%s\n", err)
		return
	}
	fmt.Printf("ID: %s, title: %s\n", file.ID, file.Title)

	err = api.DeleteFile(file.ID)
	if err != nil {
		fmt.Printf("%s\n", err)
		return
	}
	fmt.Printf("File %s deleted successfully.\n", file.ID)
}
