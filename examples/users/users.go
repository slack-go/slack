package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/slack-go/slack"
)

func main() {
	userID := flag.String("user", "", "User ID to fetch info for")
	list := flag.Bool("list", false, "List all users")
	teamID := flag.String("team", "", "Team ID (required for Enterprise Grid)")
	flag.Parse()

	userToken := os.Getenv("SLACK_USER_TOKEN")
	if userToken == "" {
		fmt.Fprintf(os.Stderr, "SLACK_USER_TOKEN environment variable is required\n")
		os.Exit(1)
	}

	if *userID == "" && !*list {
		fmt.Fprintf(os.Stderr, "Use -user <ID> to fetch a user or -list to list all users\n")
		os.Exit(1)
	}

	api := slack.New(userToken)

	if *list {
		var opts []slack.GetUsersOption
		if *teamID != "" {
			opts = append(opts, slack.GetUsersOptionTeamID(*teamID))
		}
		listUsers(api, opts...)
		return
	}

	user, err := api.GetUserInfo(*userID)
	if err != nil {
		fmt.Printf("%s\n", err)
		return
	}
	b, err := json.MarshalIndent(user, "", "  ")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(b))
	fmt.Printf("ID: %s, Fullname: %s, Email: %s\n", user.ID, user.Profile.RealName, user.Profile.Email)
}

func listUsers(api *slack.Client, opts ...slack.GetUsersOption) {
	users, err := api.GetUsers(opts...)
	if err != nil {
		log.Fatal(err)
	}

	for _, user := range users {
		guestType := ""
		if user.IsUltraRestricted {
			guestType = " (single-channel guest)"
		} else if user.IsRestricted {
			guestType = " (multi-channel guest)"
		}

		fmt.Printf("%-12s %-30s %s%s\n", user.ID, user.Profile.RealName, user.Profile.Email, guestType)

		if (user.IsRestricted || user.IsUltraRestricted) && user.Profile.GuestInvitedBy != "" {
			fmt.Printf("  invited_by: %s\n", user.Profile.GuestInvitedBy)
		}
	}
}
