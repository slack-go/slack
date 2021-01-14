package main

import (
	"fmt"
	"log"
	"os"

	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
)

// You more than likely want your "Bot User OAuth Access Token" which starts with "xoxb-"
var api = slack.New("TOKEN")

func main() {
	api := slack.New(
		"YOUR TOKEN HERE",
		slack.OptionDebug(true),
		slack.OptionLog(log.New(os.Stdout, "slack-bot: ", log.Lshortfile|log.LstdFlags)),
	)

	client := api.NewSocketModeClient()
	go client.ManageConnection()

	for evt := range client.IncomingEvents {
		eventsAPIEvent, ok := evt.Data.(slackevents.EventsAPIEvent)
		if !ok {
			fmt.Printf("Ignored %v\n")

			continue
		}

		fmt.Printf("Event Received: %+v", eventsAPIEvent)

		switch evt.Type {
		case slackevents.CallbackEvent:
			innerEvent := eventsAPIEvent.InnerEvent
			switch ev := innerEvent.Data.(type) {
			case *slackevents.AppMentionEvent:
				api.PostMessage(ev.Channel, slack.MsgOptionText("Yes, hello.", false))
			}
		case slackevents.MemberJoinedChannel:
			ev := eventsAPIEvent.Data.(*slackevents.MemberJoinedChannelEvent)

			fmt.Printf("user %q joined to channel %q", ev.User, ev.Channel)
		}
	}
}
