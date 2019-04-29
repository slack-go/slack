package main

import (
	"fmt"
	"strconv"
	"time"

	"github.com/nlopes/slack"
)

func scheduleMessageExample() {
	api := slack.New("YOUR_TOKEN_HERE")
	attachment := slack.Attachment{
		Pretext: "some pretext",
		Text:    "some text",
	}

	// Schedule message for 15 minutes from now
	sendAt := time.Now().Add(15 * time.Minute).UTC().Unix()
	sendAtString := strconv.FormatInt(sendAt, 10)

	channelID, timestamp, err := api.PostScheduledMessage("CHANNEL_ID", sendAtString, slack.MsgOptionText("Some text", false), slack.MsgOptionAttachments(attachment))
	if err != nil {
		fmt.Printf("%s\n", err)
		return
	}
	fmt.Printf("Message successfully sent to channel %s at %s", channelID, timestamp)
}
