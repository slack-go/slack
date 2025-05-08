package main

import (
	"context"
	"fmt"

	"github.com/slack-go/slack"
)

func main() {
	api := slack.New("YOUR_TOKEN")
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
