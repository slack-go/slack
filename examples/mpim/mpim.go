package main

import (
	"fmt"

	"github.com/nlopes/slack"
)

func main() {
	api := slack.New("YOUR_TOKEN_HERE")
	//api.SetDebug(true)
	groups, err := api.GetMPIMChannels()
	if err != nil {
		fmt.Printf("%s\n", err)
		return
	}
	for _, group := range groups {
		fmt.Printf("ID: %s, Name: %s, MPIM: %t\n", group.ID, group.Name, group.IsMPIM)
	}
}
