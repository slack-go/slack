// Preview with go run ./examples/layout, or post with
// SLACK_BOT_TOKEN set and go run ./examples/layout -channel C_CHANNEL.
package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/slack-go/slack"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run() error {
	channel := flag.String("channel", "", "Channel ID to post to; omit to print the message JSON")
	flag.Parse()
	section := func(text string) *slack.SectionBlock {
		return slack.NewSectionBlock(slack.NewTextBlockObject(slack.MarkdownType, text, false, false), nil, nil)
	}
	stage := section("*Stage*\nNegotiation")
	amount := section("*Amount*\n$125,000")
	closeDate := section("*Close*\n2026-12-15")
	actions := slack.NewActionBlock("deal-actions", slack.NewButtonBlockElement("view-deal", "deal-123",
		slack.NewTextBlockObject(slack.PlainTextType, "View deal", false, false)).WithURL("https://example.com/deals/123"))
	layout := slack.NewLayoutBlock("Deal summary",
		slack.NewLayoutRow(slack.NewLayoutItem(stage), slack.NewLayoutItem(amount), slack.NewLayoutItem(closeDate)),
		slack.NewLayoutRow(slack.NewLayoutItem(section("*Next steps*\nSchedule sponsor meeting")).WithWidth(slack.LayoutWidthWide),
			slack.NewLayoutItem(actions).WithWidth(slack.LayoutWidthNarrow).WithHorizontalAlignment(slack.LayoutHorizontalAlignmentEnd)),
	).WithBlockID("deal-summary").WithLabelVisible(true).WithGridLines(slack.LayoutGridLinesRows).
		WithFallback(section("*Deal summary*\nStage: Negotiation · Amount: $125,000 · Close: 2026-12-15\nNext steps: Schedule sponsor meeting\n<https://example.com/deals/123|View deal>"))
	if err := layout.Validate(); err != nil {
		return err
	}
	if *channel == "" {
		return json.NewEncoder(os.Stdout).Encode(slack.NewBlockMessage(layout))
	}
	token := os.Getenv("SLACK_BOT_TOKEN")
	if token == "" {
		return fmt.Errorf("SLACK_BOT_TOKEN is required when posting with -channel")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	postedChannel, timestamp, err := slack.New(token).PostMessageContext(ctx, *channel,
		slack.MsgOptionText("Deal summary: Negotiation, $125,000, closing 2026-12-15. Next steps: Schedule sponsor meeting. View deal: https://example.com/deals/123", false),
		slack.MsgOptionBlocks(layout))
	if err != nil {
		return err
	}
	fmt.Printf("Posted layout to %s at %s\n", postedChannel, timestamp)
	return nil
}
