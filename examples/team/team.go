package main

import (
	"fmt"

	"github.com/slack-go/slack"
)

func main() {
	api := slack.New("YOUR_TOKEN_HERE")
	//Example for single user
	billingActive, err := api.GetBillableInfo(slack.GetBillableInfoParams{User: "U023BECGF"})
	if err != nil {
		fmt.Printf("%s\n", err)
		return
	}
	fmt.Printf("ID: U023BECGF, BillingActive: %v\n\n\n", billingActive["U023BECGF"])

	//Example for team. Note: passing empty TeamID just uses the current user team.
	billingActiveForTeam, _ := api.GetBillableInfo(slack.GetBillableInfoParams{})
	for id, value := range billingActiveForTeam {
		fmt.Printf("ID: %v, BillingActive: %v\n", id, value)
	}
}
