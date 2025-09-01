package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
	"github.com/slack-go/slack/socketmode"
)

// This is an example of AI apps https://api.slack.com/docs/apps/ai
// Developing and using AI apps requires a paid plan or joining the Developer Program and provision a sandbox with access to all Slack features for free.
//
// This example also have calling https://api.slack.com/docs/apps/data-access-api sample code.
// This API is currently in a limited access stage. You may be able to obtain a token and call the API, but to get a valid response, you must be enrolled in the program. Contact Customer Experience at feedback@slack.com to request to be added.
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

	// Handle a specific event from EventsAPI
	socketmodeHandler.HandleEvents(slackevents.AssistantThreadStarted, middlewareAssistantThreadStartedEvent)
	socketmodeHandler.HandleEvents(slackevents.AppMention, middlewareAppMentionEvent)

	socketmodeHandler.RunEventLoop()
}

func middlewareAssistantThreadStartedEvent(evt *socketmode.Event, client *socketmode.Client) {
	ctx := context.Background()
	fmt.Printf("assistant thread started event: %+v\n", evt)

	eventsAPIEvent, ok := evt.Data.(slackevents.EventsAPIEvent)
	if !ok {
		fmt.Printf("Ignored: %+v\n", evt)
		return
	}

	innerEvent, ok := eventsAPIEvent.InnerEvent.Data.(*slackevents.AssistantThreadStartedEvent)
	if !ok {
		fmt.Printf("Can't get inner event: %+v\n", evt)
		return
	}

	fmt.Printf("Inner event: %+v\n", innerEvent)

	params := slack.AssistantThreadsSetSuggestedPromptsParameters{
		Title:     "Welcome. What can I do for you?",
		ChannelID: innerEvent.AssistantThread.ChannelID,
		ThreadTS:  innerEvent.AssistantThread.ThreadTimeStamp,
		Prompts: []slack.AssistantThreadsPrompt{
			{
				Title:   "Generate ideas",
				Message: "Pretend you are a marketing associate and you need new ideas for an enterprise productivity feature. Generate 10 ideas for a new feature launch.",
			},
			{
				Title:   "Explain what SLACK stands for",
				Message: "What does SLACK stand for?",
			},
		},
	}

	if err := client.SetAssistantThreadsSuggestedPromptsContext(ctx, params); err != nil {
		fmt.Printf("Can't SetAssistantThreadsSuggestedPromptsContext: %v\n", err)
		return
	}
}

func middlewareAppMentionEvent(evt *socketmode.Event, client *socketmode.Client) {
	ctx := context.Background()
	eventsAPIEvent, ok := evt.Data.(slackevents.EventsAPIEvent)
	if !ok {
		fmt.Printf("Ignored %+v\n", evt)
		return
	}
	client.Ack(*evt.Request)

	ev, ok := eventsAPIEvent.InnerEvent.Data.(*slackevents.AppMentionEvent)
	if !ok {
		fmt.Printf("Ignored %+v\n", evt)
		return
	}
	fmt.Printf("We have been mentioned in %v\n", ev.Channel)
	if err := searchAssistantContext(ctx, client, ev.AssistantThread, ev.Text, ev.Channel, ev.TimeStamp); err != nil {
		fmt.Printf("Failed searchAssistantContext: %+v\n", err)
		return
	}
}

// This is sample of calling https://api.slack.com/docs/apps/data-access-api
// This API is currently in a limited access stage. You may be able to obtain a token and call the API, but to get a valid response, you must be enrolled in the program. Contact Customer Experience at feedback@slack.com to request to be added.
func searchAssistantContext(ctx context.Context, client *socketmode.Client, at *slackevents.AssistantThreadActionToken, text, channel, ts string) error {
	// Assistant thread message handling
	if at != nil {
		fmt.Printf("Assistant thread message received - text: %s, channel: %s",
			text, channel)

		// Call Data Access API for context search
		resp, err := client.SearchAssistantContextContext(ctx, slack.AssistantSearchContextParameters{
			Query:        text,
			ActionToken:  at.ActionToken,
			ChannelTypes: []string{"public_channel"},
			ContentTypes: []string{"messages"},
			Limit:        10,
		})
		if err != nil {
			return err
		}
		fmt.Printf("SearchAssistantContextContext response: %+v\n", resp)

		if len(resp.Results.Messages) > 0 {
			_, _, err = client.Client.PostMessage(channel, slack.MsgOptionText("Hello! I searched your query.text:\n: "+text+"\nSearch first result:\n"+resp.Results.Messages[0].Content, false), slack.MsgOptionTS(ts))
			if err != nil {
				return err
			}
		}
	}
	return nil
}
