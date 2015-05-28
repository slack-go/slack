package main

import (
	"fmt"

	"github.com/abourget/slack"
)

func main() {
	api := slack.New("YOUR TOKEN HERE")
	api.SetDebug(true)

	rtm := api.NewRTM()
	go rtm.ManageConnection()

	for {
		select {
		case msg := <-rtm.IncomingEvents:
			fmt.Print("Event Received: ")
			switch ev := msg.Data.(type) {
			case *slack.HelloEvent:
				// Ignore hello

			case *slack.ConnectedEvent:
				fmt.Println("Infos:", ev.Info)
				fmt.Println("Connection counter:", ev.ConnectionCount)
				rtm.SendMessage(rtm.NewOutgoingMessage("Hello world", "#general"))

			case *slack.MessageEvent:
				fmt.Printf("Message: %v\n", ev)

			case *slack.PresenceChangeEvent:
				fmt.Printf("Presence Change: %v\n", ev)

			case *slack.LatencyReport:
				fmt.Printf("Current latency: %v\n", ev.Value)

			case *slack.SlackWSError:
				fmt.Printf("Error: %d - %s\n", ev.Code, ev.Msg)

			default:
				// Ignore other events..
				//fmt.Printf("Unexpected: %v\n", msg.Data)
			}
		}
	}
}
