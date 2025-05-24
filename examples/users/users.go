package main

import (
	"fmt"
	"os"

	"github.com/slack-go/slack"
)

func main() {
	userToken := os.Getenv("SLACK_USER_TOKEN")
	if userToken == "" {
		fmt.Fprintf(os.Stderr, "SLACK_USER_TOKEN must be set.\n")
		os.Exit(1)
	}

	api := slack.New(userToken)
	user, err := api.GetUserInfo("U023BECGF")
	if err != nil {
		fmt.Printf("%s\n", err)
		return
	}
	fmt.Printf("ID: %s, Fullname: %s, Email: %s\n", user.ID, user.Profile.RealName, user.Profile.Email)
}
