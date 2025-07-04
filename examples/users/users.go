package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/slack-go/slack"
)

func main() {
	userID := flag.String("user", "", "User ID (required)")
	flag.Parse()

	// Get token from environment variable
	userToken := os.Getenv("SLACK_USER_TOKEN")
	if userToken == "" {
		fmt.Fprintf(os.Stderr, "SLACK_USER_TOKEN environment variable is required\n")
		os.Exit(1)
	}

	// Get user ID from flag
	if *userID == "" {
		fmt.Println("User ID is required: use -user flag")
		os.Exit(1)
	}

	api := slack.New(userToken)
	user, err := api.GetUserInfo(*userID)
	if err != nil {
		fmt.Printf("%s\n", err)
		return
	}
	fmt.Printf("ID: %s, Fullname: %s, Email: %s\n", user.ID, user.Profile.RealName, user.Profile.Email)
}
