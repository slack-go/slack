package main

import (
	"context"
	"fmt"
	"os"

	"github.com/slack-go/slack"
)

func main() {
	token, ok := os.LookupEnv("SLACK_BOT_TOKEN")
	if !ok {
		fmt.Println("Missing SLACK_BOT_TOKEN in environment")
		os.Exit(1)
	}
	channelID := "CXXXXXXXX" // Replace with your channel ID

	api := slack.New(token, slack.OptionDebug(true))
	ctx := context.Background()

	files := []slack.UploadFileParameters{
		{
			Title:    "Batman Example",
			Filename: "batman.txt",
			File:     "batman.txt",
			FileSize: 39,
		},
		{
			Title:    "Superman Example",
			Filename: "superman.txt",
			File:     "superman.txt",
			FileSize: 37,
		},
	}

	uploads := []*slack.GetUploadURLExternalResponse{}
	filesToComplete := []slack.FileSummary{}

	for _, file := range files {
		u, err := api.GetUploadURLExternalContext(ctx, slack.GetUploadURLExternalParameters{
			AltTxt:   "An alt text for superheroes",
			FileName: file.Filename,
			FileSize: file.FileSize,
		})
		if err != nil {
			fmt.Printf("%s\n", err)
			return
		}
		uploads = append(uploads, u)
	}

	for i, file := range files {
		fmt.Printf("Uploading file %s to %s\n", file.Filename, uploads[i].UploadURL)

		err := api.UploadToURL(ctx, slack.UploadToURLParameters{
			UploadURL: uploads[i].UploadURL,
			Filename:  file.Filename,
			File:      file.File,
		})
		if err != nil {
			fmt.Printf("%s\n", err)
			return
		}
		filesToComplete = append(filesToComplete, slack.FileSummary{
			ID:    uploads[i].FileID,
			Title: file.Title,
		})
	}

	c, err := api.CompleteUploadExternalContext(ctx, slack.CompleteUploadExternalParameters{
		Files:   filesToComplete,
		Channel: channelID,
	})

	if err != nil {
		fmt.Printf("%s\n", err)
		return
	}

	fmt.Printf("Files uploaded successfully: %+v\n", c.Files)
}
