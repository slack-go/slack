package main

import (
	"fmt"
	"os"

	"github.com/slack-go/slack"
)

func main() {
	api := slack.New("YOUR_TOKEN_HERE")
	r, err := os.Open("slack-go.png")
	if err != nil {
		fmt.Printf("%s\n", err)
		return
	}
	defer r.Close()
	remotefile, err := api.AddRemoteFile(slack.RemoteFileParameters{
		ExternalID:            "slack-go",
		ExternalURL:           "https://github.com/slack-go/slack",
		Title:                 "slack-go",
		Filetype:              "go",
		IndexableFileContents: "golang, slack",
		// PreviewImage:          "slack-go.png",
		PreviewImageReader: r,
	})
	if err != nil {
		fmt.Printf("add remote file failed: %s\n", err)
		return
	}
	fmt.Printf("remote file: %v\n", remotefile)

	_, err = api.ShareRemoteFile([]string{"CPB8DC1CM"}, remotefile.ExternalID, "")
	if err != nil {
		fmt.Printf("share remote file failed: %s\n", err)
		return
	}
	fmt.Printf("share remote file %s successfully.\n", remotefile.Name)

	remotefiles, err := api.ListRemoteFiles(slack.ListRemoteFilesParameters{
		Channel: "YOUR_CHANNEL_HERE",
	})
	if err != nil {
		fmt.Printf("list remote files failed: %s\n", err)
		return
	}
	fmt.Printf("remote files: %v\n", remotefiles)

	remotefile, err = api.UpdateRemoteFile(remotefile.ID, slack.RemoteFileParameters{
		ExternalID:            "slack-go",
		ExternalURL:           "https://github.com/slack-go/slack",
		Title:                 "slack-go",
		Filetype:              "go",
		IndexableFileContents: "golang, slack, github",
	})
	if err != nil {
		fmt.Printf("update remote file failed: %s\n", err)
		return
	}
	fmt.Printf("remote file: %v\n", remotefile)

	info, err := api.GetRemoteFileInfo(remotefile.ExternalID, "")
	if err != nil {
		fmt.Printf("get remote file info failed: %s\n", err)
		return
	}
	fmt.Printf("remote file info: %v\n", info)

	err = api.RemoveRemoteFile(remotefile.ExternalID, "")
	if err != nil {
		fmt.Printf("remove remote file failed: %s\n", err)
		return
	}
	fmt.Printf("remote file %s deleted successfully.\n", remotefile.Name)
}
