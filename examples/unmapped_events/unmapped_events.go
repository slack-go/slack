// This example connects to Slack via RTM and listens for events.
//
// Before the fix for #1544, the events "apps_uninstalled", "activity", and
// "badge_counts_updated" were not mapped and would appear as
// UnmarshallingErrorEvent with an "Received unmapped event" message.
//
// Usage:
//
//	export SLACK_BOT_TOKEN=xoxb-...
//	go run ./examples/unmapped_events/
//
// Then interact with your workspace (open channels, browse around) and watch
// for unmapped event errors in the output. The "activity" and
// "badge_counts_updated" events tend to appear during normal workspace usage.
// The "apps_uninstalled" event appears when an app is removed.
package main

import (
	"fmt"
	"log"
	"os"

	"github.com/slack-go/slack"
)

func main() {
	token := os.Getenv("SLACK_BOT_TOKEN")
	if token == "" {
		fmt.Println("SLACK_BOT_TOKEN environment variable is required")
		os.Exit(1)
	}

	api := slack.New(
		token,
		slack.OptionDebug(true),
		slack.OptionLog(log.New(os.Stdout, "slack-bot: ", log.Lshortfile|log.LstdFlags)),
	)

	rtm := api.NewRTM()
	go rtm.ManageConnection()

	for msg := range rtm.IncomingEvents {
		switch ev := msg.Data.(type) {
		case *slack.ConnectedEvent:
			fmt.Printf("Connected: %s (connection count: %d)\n", ev.Info.User.ID, ev.ConnectionCount)

		case *slack.AppsUninstalledEvent:
			fmt.Printf("Apps uninstalled: %+v\n", ev)

		case *slack.ActivityEvent:
			fmt.Printf("Activity: subtype=%s key=%s\n", ev.SubType, ev.Key)

		case *slack.BadgeCountsUpdatedEvent:
			fmt.Printf("Badge counts updated: %+v\n", ev)

		case *slack.UnmarshallingErrorEvent:
			// Before the fix, apps_uninstalled/activity/badge_counts_updated
			// end up here as unmapped events.
			fmt.Printf("UNMAPPED EVENT ERROR: %v\n", ev.ErrorObj)

		case *slack.InvalidAuthEvent:
			fmt.Println("Invalid credentials")
			return

		case *slack.DisconnectedEvent:
			if ev.Intentional {
				return
			}

		default:
			fmt.Printf("Event: type=%s data=%T\n", msg.Type, ev)
		}
	}
}
