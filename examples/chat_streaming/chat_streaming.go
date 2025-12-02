package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
	"github.com/slack-go/slack/socketmode"
)

// This example demonstrates using Slack's chat streaming API.
// It listens for app mentions and streams a response back in real-time.
//
// Required environment variables:
// - SLACK_APP_TOKEN: Your Slack app token (starts with xapp-)
// - SLACK_BOT_TOKEN: Your Slack bot token (starts with xoxb-)
//
// Required Slack app scopes:
// - app_mentions:read
// - chat:write
//
// Required Event Subscriptions:
// - app_mention

func main() {
	appToken := os.Getenv("SLACK_APP_TOKEN")
	if appToken == "" {
		panic("SLACK_APP_TOKEN must be set.\n")
	}

	if !strings.HasPrefix(appToken, "xapp-") {
		panic("SLACK_APP_TOKEN must have the prefix \"xapp-\".")
	}

	botToken := os.Getenv("SLACK_BOT_TOKEN")
	if botToken == "" {
		panic("SLACK_BOT_TOKEN must be set.\n")
	}

	if !strings.HasPrefix(botToken, "xoxb-") {
		panic("SLACK_BOT_TOKEN must have the prefix \"xoxb-\".")
	}

	api := slack.New(
		botToken,
		slack.OptionDebug(true),
		slack.OptionLog(log.New(os.Stdout, "api: ", log.Lshortfile|log.LstdFlags)),
		slack.OptionAppLevelToken(appToken),
	)

	client := socketmode.New(
		api,
		socketmode.OptionDebug(true),
		socketmode.OptionLog(log.New(os.Stdout, "socketmode: ", log.Lshortfile|log.LstdFlags)),
	)

	socketmodeHandler := socketmode.NewSocketmodeHandler(client)

	// Handle app mentions
	socketmodeHandler.HandleEvents(slackevents.AppMention, func(evt *socketmode.Event, clt *socketmode.Client) {
		handleAppMention(evt, clt)
	})

	log.Println("Starting chat streaming bot...")
	socketmodeHandler.RunEventLoop()
}

func handleAppMention(evt *socketmode.Event, client *socketmode.Client) {
	eventsAPIEvent, ok := evt.Data.(slackevents.EventsAPIEvent)
	if !ok {
		fmt.Printf("Ignored: %+v\n", evt)
		return
	}

	// Acknowledge the event
	client.Ack(*evt.Request)

	ev, ok := eventsAPIEvent.InnerEvent.Data.(*slackevents.AppMentionEvent)
	if !ok {
		fmt.Printf("Ignored: %+v\n", evt)
		return
	}

	log.Printf("Received app mention in channel %s: %s", ev.Channel, ev.Text)

	// Start the stream
	channel, streamTS, err := client.Client.StartStream(
		ev.Channel,
		slack.MsgOptionTS(ev.TimeStamp), // Reply in thread
	)
	if err != nil {
		log.Printf("Failed to start stream: %v", err)
		return
	}

	log.Printf("Started stream in channel %s with timestamp %s", channel, streamTS)

	// Simulate a streaming response by breaking up a message into chunks
	response := "Here's a streaming response! " +
		"This example demonstrates how to use Slack's chat streaming API. " +
		"The streaming API consists of three methods: " +
		"StartStream to begin streaming, " +
		"AppendStream to add content incrementally, " +
		"and StopStream to finish the stream. " +
		"This is perfect for AI-powered apps that generate responses progressively."

	// Stream the response in chunks
	if err := streamResponse(&client.Client, channel, streamTS, response); err != nil {
		log.Printf("Error during streaming: %v", err)

		// Try to stop stream with error message
		_, _, _ = client.Client.StopStream(
			channel,
			streamTS,
			slack.MsgOptionMarkdownText("\n\n_Error occurred while streaming response._"),
		)
		return
	}

	// Create feedback buttons
	thumbsUpText := slack.NewTextBlockObject(slack.PlainTextType, "👍 Helpful", true, false)
	thumbsDownText := slack.NewTextBlockObject(slack.PlainTextType, "👎 Not Helpful", true, false)

	feedbackButtons := slack.NewActionBlock(
		"feedback",
		slack.NewButtonBlockElement("thumbs_up", "thumbs_up", thumbsUpText),
		slack.NewButtonBlockElement("thumbs_down", "thumbs_down", thumbsDownText),
	)

	// Stop the stream with feedback buttons
	_, _, err = client.Client.StopStream(
		channel,
		streamTS,
		slack.MsgOptionBlocks(feedbackButtons),
	)
	if err != nil {
		log.Printf("Failed to stop stream: %v", err)
		return
	}

	log.Printf("Successfully completed streaming response in channel %s", channel)
}

// streamResponse demonstrates buffering and streaming text chunks to Slack
func streamResponse(api *slack.Client, channel, streamTS, response string) error {
	const (
		chunkSize  = 5  // Characters to simulate per "chunk"
		bufferSize = 20 // Send to Slack when buffer reaches this size
		delayMS    = 50 // Milliseconds between chunks (simulates generation)
	)

	buffer := strings.Builder{}
	words := strings.Split(response, " ")

	for i, word := range words {
		// Add word and space to buffer
		if i > 0 {
			buffer.WriteString(" ")
		}
		buffer.WriteString(word)

		// Simulate streaming delay
		time.Sleep(time.Duration(delayMS) * time.Millisecond)

		// Send to Slack when buffer reaches threshold or at the end
		if buffer.Len() >= bufferSize || i == len(words)-1 {
			_, _, err := api.AppendStream(
				channel,
				streamTS,
				slack.MsgOptionMarkdownText(buffer.String()),
			)
			if err != nil {
				return fmt.Errorf("failed to append to stream: %w", err)
			}

			log.Printf("Appended %d characters to stream", buffer.Len())
			buffer.Reset()
		}
	}

	return nil
}
