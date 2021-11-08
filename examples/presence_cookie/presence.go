package main

import (
	"fmt"
	"os"

	"github.com/slack-go/slack"
)

func main() {
	token := os.Getenv("SLACK_TOKEN")
	cookie := os.Getenv("SLACK_COOKIE")

	api := slack.NewWithCookie(token, cookie)
	pres, err := api.GetUserPresence("")
	if err != nil {
		fmt.Printf("%s\n", err)
		return
	}
	fmt.Printf("Presence: %s, Online: %t, AutoAway: %t, ManualAway: %t, ConnectionCount: %d, LastActivity: %v\n", pres.Presence, pres.Online, pres.AutoAway, pres.ManualAway, pres.ConnectionCount, pres.LastActivity)
}
