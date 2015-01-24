package main

import (
	"fmt"

	"github.com/nlopes/slack"
)

func main() {
	api := slack.New("YOUR_TOKEN_HERE")
	user, err := api.GetUserInfo("U023BECGF")
	if err != nil {
		fmt.Printf("%s\n", err)
		return
	}
	fmt.Printf("Id: %s, Fullname: %s, Email: %s\n", user.Id, user.Profile.RealName, user.Profile.Email)
}
