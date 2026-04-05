// This example connects via RTM and prints Slack Call/Huddle room events
// (sh_room_join, sh_room_leave). Start a call or huddle in a channel where
// the bot is present to see the events.
//
// To run:
//
//	export SLACK_BOT_TOKEN=xoxb-...
//	go run examples/rtm_call_events/rtm_call_events.go
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
		fmt.Fprintln(os.Stderr, "SLACK_BOT_TOKEN environment variable is required")
		os.Exit(1)
	}

	api := slack.New(
		token,
		slack.OptionDebug(true),
		slack.OptionLog(log.New(os.Stdout, "rtm: ", log.Lshortfile|log.LstdFlags)),
	)

	rtm := api.NewRTM()
	go rtm.ManageConnection()

	fmt.Println("Listening for call/huddle events... start a call in a channel where this bot is present.")

	for msg := range rtm.IncomingEvents {
		switch ev := msg.Data.(type) {
		case *slack.ConnectedEvent:
			fmt.Printf("Connected as %s\n", ev.Info.User.Name)

		case *slack.SHRoomJoinEvent:
			fmt.Printf("User %s joined call in room %s (channels: %v, participants: %v)\n",
				ev.User, ev.Room.ID, ev.Room.Channels, ev.Room.Participants)

		case *slack.SHRoomLeaveEvent:
			fmt.Printf("User %s left call in room %s (remaining: %v)\n",
				ev.User, ev.Room.ID, ev.Room.Participants)

		case *slack.SHRoomUpdateEvent:
			name := "<unnamed>"
			if ev.Room.Name != nil {
				name = *ev.Room.Name
			}
			fmt.Printf("Room %s updated: %q (family: %s, participants: %v)\n",
				ev.Room.ID, name, ev.Room.CallFamily, ev.Room.Participants)

		case *slack.RTMError:
			fmt.Printf("RTM Error: %s\n", ev.Error())

		case *slack.InvalidAuthEvent:
			fmt.Fprintln(os.Stderr, "Invalid credentials")
			return
		}
	}
}
