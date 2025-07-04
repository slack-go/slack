package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/slack-go/slack"
)

func main() {
	userID := flag.String("user", "", "User ID for billing info (optional)")
	flag.Parse()

	// Get token from environment variable
	token := os.Getenv("SLACK_BOT_TOKEN")
	if token == "" {
		fmt.Println("SLACK_BOT_TOKEN environment variable is required")
		os.Exit(1)
	}

	api := slack.New(token)

	if *userID != "" {
		// Example for single user
		billingActive, err := api.GetBillableInfo(slack.GetBillableInfoParams{User: *userID})
		if err != nil {
			fmt.Printf("%s\n", err)
			return
		}
		fmt.Printf("ID: %s, BillingActive: %v\n\n\n", *userID, billingActive[*userID])
	} else {
		// Example for team. Note: passing empty TeamID just uses the current user team.
		billingActiveForTeam, err := api.GetBillableInfo(slack.GetBillableInfoParams{})
		if err != nil {
			fmt.Printf("%s\n", err)
			return
		}
		for id, value := range billingActiveForTeam {
			fmt.Printf("ID: %v, BillingActive: %v\n", id, value)
		}
	}
}
