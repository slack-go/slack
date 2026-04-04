// This example demonstrates Socket Mode's behavior with large Ack payloads.
//
// Slack's Socket Mode silently drops WebSocket responses that are 20KB or
// larger. The library detects this and returns an error from Ack().
//
// To run:
//
//	export SLACK_APP_TOKEN=xapp-...
//	export SLACK_BOT_TOKEN=xoxb-...
//	go run examples/socketmode/large_payload_ack/large_payload_ack.go
//
// Then use your slash command with a byte count as the argument:
//
//	/your-command 1000     -> small payload, works
//	/your-command 19000    -> near the 20KB limit, works
//	/your-command 21000    -> over the limit, Ack() returns an error
//	/your-command          -> defaults to 100 bytes
package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/slack-go/slack"
	"github.com/slack-go/slack/socketmode"
)

func main() {
	appToken := os.Getenv("SLACK_APP_TOKEN")
	if appToken == "" {
		fmt.Fprintf(os.Stderr, "SLACK_APP_TOKEN environment variable is required\n")
		os.Exit(1)
	}

	if !strings.HasPrefix(appToken, "xapp-") {
		fmt.Fprintf(os.Stderr, "SLACK_APP_TOKEN must have the prefix \"xapp-\"\n")
		os.Exit(1)
	}

	botToken := os.Getenv("SLACK_BOT_TOKEN")
	if botToken == "" {
		fmt.Fprintf(os.Stderr, "SLACK_BOT_TOKEN environment variable is required\n")
		os.Exit(1)
	}

	if !strings.HasPrefix(botToken, "xoxb-") {
		fmt.Fprintf(os.Stderr, "SLACK_BOT_TOKEN must have the prefix \"xoxb-\"\n")
		os.Exit(1)
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

	go func() {
		for evt := range client.Events {
			switch evt.Type {
			case socketmode.EventTypeConnecting:
				fmt.Println("Connecting to Slack with Socket Mode...")
			case socketmode.EventTypeConnectionError:
				fmt.Println("Connection failed. Retrying later...")
			case socketmode.EventTypeConnected:
				fmt.Println("Connected to Slack with Socket Mode.")
			case socketmode.EventTypeSlashCommand:
				cmd, ok := evt.Data.(slack.SlashCommand)
				if !ok {
					fmt.Printf("Ignored %+v\n", evt)
					continue
				}

				size := 100
				if cmd.Text != "" {
					if n, err := strconv.Atoi(strings.TrimSpace(cmd.Text)); err == nil && n > 0 {
						size = n
					}
				}

				payload := map[string]any{
					"text": fmt.Sprintf("[%d bytes] %s", size, strings.Repeat("x", size)),
				}

				if err := client.Ack(*evt.Request, payload); err != nil {
					fmt.Printf("Ack() error: %v\n", err)
					fmt.Println("Use the Web API (e.g. chat.PostMessage) for large payloads.")

					// Ack without payload so Slack knows we received the event.
					client.Ack(*evt.Request)
				}

			default:
			}
		}
	}()

	client.Run()
}
