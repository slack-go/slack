package main

import (
	"fmt"

	"github.com/nlopes/slack"
)

func main() {
	api := slack.New("YOUR_TOKEN_HERE")
	// If you set debugging, it will log all requests to the console
	// Useful when encountering issues
	// api.SetDebug(true)
	usergroups, err := api.GetUserGroups(false, true, false)
	if err != nil {
		fmt.Printf("%s\n", err)
		return
	}
	for _, usergroup := range usergroups {
		fmt.Printf("ID: %s, Name: %s, Handle: %s, Description: %s\n", usergroup.ID, usergroup.Name, usergroup.Handle, usergroup.Description)
	}
}
