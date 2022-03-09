package main

import (
	"fmt"
	"os"
	"time"

	"github.com/slack-go/slack"
)

func main() {
	// Build Slack client with proper auth
	token, ok := os.LookupEnv("SLACK_TOKEN")
	if !ok {
		fmt.Println("Missing SLACK_TOKEN in environment")
		os.Exit(1)
	}
	cookie, ok := os.LookupEnv("SLACK_COOKIE")
	if !ok {
		fmt.Println("Missing SLACK_COOKIE in environment")
		os.Exit(1)
	}
	api := slack.NewWithCookie(token, cookie)

	// Periodically query for current user presence in the background
	go func() {
		for {
			pres, err := api.GetUserPresence("")
			if err != nil {
				fmt.Printf("%s\n", err)
				continue
			}
			fmt.Printf("Presence: %s, Online: %t, AutoAway: %t, ManualAway: %t, ConnectionCount: %d, LastActivity: %v\n", pres.Presence, pres.Online, pres.AutoAway, pres.ManualAway, pres.ConnectionCount, pres.LastActivity)
			time.Sleep(10 * time.Second)
		}
	}()

	// Initiate RTM connection to Slack in the background
	rtm := api.NewRTM()
	go rtm.ManageConnection()

	// Handle incoming RTM events
	for msg := range rtm.IncomingEvents {
		switch ev := msg.Data.(type) {

		case *slack.ConnectedEvent:
			fmt.Printf("[ConnectedEvent] Info: %v; ConnectionCount: %d\n", ev.Info, ev.ConnectionCount)

		case *slack.RTMError:
			fmt.Printf("[RTMError] Error: %s\n", ev.Error())

		case *slack.InvalidAuthEvent:
			fmt.Printf("[InvalidAuthEvent] Invalid credentials\n")
			return

		default:
			// Ignore other events..
			// fmt.Printf("Unexpected: %v\n", msg.Data)
		}
	}
}
