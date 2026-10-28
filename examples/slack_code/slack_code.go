// Exercise Slack Code's channel, view, and canvas APIs. See README.md for setup.
package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/slack-go/slack"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		var slackErr slack.SlackErrorResponse
		if errors.As(err, &slackErr) && len(slackErr.ResponseMetadata.Messages) > 0 {
			fmt.Fprintf(os.Stderr, "Slack response messages: %v\n", slackErr.ResponseMetadata.Messages)
		}
		os.Exit(1)
	}
}

func run() error {
	action := flag.String("action", "", "Action: create, info, properties, commands, set-view, list-views, remove-view, create-canvas, get-canvas, set-canvas-content, archive")
	channelID := flag.String("channel", "", "Code channel ID (required except for create and create-canvas)")
	name := flag.String("name", "", "Channel/canvas name (required for create and create-canvas), or view tab name")
	originChannel := flag.String("origin-channel", "", "Source channel containing the message that started this task (create)")
	originTS := flag.String("origin-ts", "", "Source message timestamp; use with -origin-channel (create)")
	viewKey := flag.String("view-key", "", "Stable view key (set-view, remove-view)")
	viewID := flag.String("view-id", "", "View ID; use exactly one of -view-id or -view-key (remove-view)")
	viewType := flag.String("type", slack.AgentViewTypeHTML, "View type: html, diff, block_kit, canvas; documented pull_request may be rejected by Slack (set-view)")
	contentFile := flag.String("content-file", "", "File containing HTML, a unified diff, a JSON blocks array, or canvas markdown")
	canvasID := flag.String("canvas", "", "Canvas ID (canvas views, get-canvas, set-canvas-content)")
	access := flag.String("access", "comment", "Canvas view access: read, write, comment (set-view)")
	prURL := flag.String("pr-url", "", "Pull request URL (pull_request views)")
	resourceDomain := flag.String("resource-domain", "", "HTTPS origin for HTML scripts/styles/images (set-view)")
	resourceURL := flag.String("resource-url", "", "External session resource URL (properties)")
	summaryTS := flag.String("summary-ts", "", "Timestamp of an existing summary message (properties, archive)")
	summaryThreadTS := flag.String("summary-thread-ts", "", "Parent thread timestamp of a summary reply; requires -summary-ts (properties)")
	includeResolved := flag.Bool("include-resolved", false, "Include resolved canvas comment threads (get-canvas)")
	clear := flag.Bool("clear", false, "Clear the agent's registered commands (commands)")
	timeout := flag.Duration("timeout", 30*time.Second, "Timeout for the API call")
	flag.Parse()
	if flag.NArg() != 0 {
		return errors.New("use flags for all arguments; see -help")
	}
	if *timeout <= 0 {
		return errors.New("-timeout must be positive")
	}
	switch *action {
	case "create", "create-canvas":
		if *name == "" {
			return errors.New("-name is required for create and create-canvas")
		}
	case "info", "properties", "commands", "set-view", "list-views", "remove-view", "get-canvas", "set-canvas-content", "archive":
		if *channelID == "" {
			return errors.New("-channel is required for this action")
		}
	default:
		return errors.New("a valid -action is required; see -help")
	}
	if *action == "create" && (*originChannel == "") != (*originTS == "") {
		return errors.New("-origin-channel and -origin-ts must be provided together")
	}
	if *action == "remove-view" && (*viewKey == "") == (*viewID == "") {
		return errors.New("remove-view requires exactly one of -view-key or -view-id")
	}
	if (*action == "get-canvas" || *action == "set-canvas-content") && *canvasID == "" {
		return errors.New("-canvas is required for get-canvas and set-canvas-content")
	}
	if *summaryThreadTS != "" && (*action != "properties" || *summaryTS == "") {
		return errors.New("-summary-thread-ts requires -action properties and -summary-ts")
	}

	var content []byte
	needsContent := *action == "create-canvas" || *action == "set-canvas-content" ||
		(*action == "set-view" && (*viewType == slack.AgentViewTypeHTML || *viewType == slack.AgentViewTypeDiff || *viewType == slack.AgentViewTypeBlockKit))
	if needsContent {
		if *contentFile == "" {
			return errors.New("-content-file is required for this action and view type")
		}
		var err error
		content, err = os.ReadFile(*contentFile)
		if err != nil {
			return fmt.Errorf("read -content-file: %w", err)
		}
	}

	var viewOptions []slack.SetAgentConversationViewOption
	if *action == "set-view" {
		viewOptions = append(viewOptions, slack.SetAgentConversationViewOptionType(*viewType))
		if *viewKey != "" {
			viewOptions = append(viewOptions, slack.SetAgentConversationViewOptionViewKey(*viewKey))
		}
		if *name != "" {
			viewOptions = append(viewOptions, slack.SetAgentConversationViewOptionName(*name))
		}
		if (*viewType == slack.AgentViewTypeHTML || *viewType == slack.AgentViewTypeBlockKit || *viewType == slack.AgentViewTypeCanvas) && *viewKey == "" {
			return errors.New("-view-key is required for html, block_kit, and canvas views")
		}
		switch *viewType {
		case slack.AgentViewTypeHTML, slack.AgentViewTypeDiff:
			viewOptions = append(viewOptions, slack.SetAgentConversationViewOptionContent(string(content)))
			if *viewType == slack.AgentViewTypeHTML && *resourceDomain != "" {
				viewOptions = append(viewOptions, slack.SetAgentConversationViewOptionCSP(slack.AgentViewCSP{
					ResourceDomains: []string{*resourceDomain},
				}))
			}
		case slack.AgentViewTypeBlockKit:
			var blocks slack.Blocks
			if err := json.Unmarshal(content, &blocks); err != nil {
				return fmt.Errorf("decode JSON blocks array: %w", err)
			}
			if len(blocks.BlockSet) == 0 {
				return errors.New("the JSON blocks array must contain at least one block")
			}
			viewOptions = append(viewOptions, slack.SetAgentConversationViewOptionBlocks(blocks.BlockSet...))
		case slack.AgentViewTypeCanvas:
			if *canvasID == "" {
				return errors.New("-canvas is required for a canvas view")
			}
			switch *access {
			case "read", "write", "comment":
			default:
				return errors.New("-access must be read, write, or comment")
			}
			viewOptions = append(viewOptions, slack.SetAgentConversationViewOptionCanvasID(*canvasID),
				slack.SetAgentConversationViewOptionAccessLevel(*access))
		case slack.AgentViewTypePullRequest:
			if *prURL == "" {
				return errors.New("-pr-url is required for a pull_request view")
			}
			viewOptions = append(viewOptions, slack.SetAgentConversationViewOptionPRURL(*prURL))
		default:
			return errors.New("-type must be html, diff, block_kit, canvas, or pull_request")
		}
	}

	token := os.Getenv("SLACK_BOT_TOKEN")
	if token == "" {
		return errors.New("SLACK_BOT_TOKEN environment variable is required")
	}
	api := slack.New(token)
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()
	ctx, cancel := context.WithTimeout(ctx, *timeout)
	defer cancel()

	var response any
	var err error
	switch *action {
	case "create":
		var options []slack.CreateAgentConversationOption
		if *originChannel != "" {
			options = append(options, slack.CreateAgentConversationOptionOriginChannelID(*originChannel),
				slack.CreateAgentConversationOptionOriginMessageTimestamp(*originTS))
		}
		response, err = api.CreateAgentConversationContext(ctx, *name, options...)
	case "info":
		var channel *slack.Channel
		channel, err = api.GetConversationInfoContext(ctx, &slack.GetConversationInfoInput{ChannelID: *channelID})
		if err == nil {
			response = map[string]any{
				"channel_id": channel.ID, "name": channel.Name,
				"is_archived": channel.IsArchived, "properties": channel.Properties,
			}
		}
	case "properties":
		properties := slack.AgentCodeChannelProperties{
			ContextBarItems: []slack.AgentContextBarItem{{
				Key: "library", Label: "slack-go/slack", Icon: "code", URL: "https://github.com/slack-go/slack",
			}},
		}
		if *summaryTS != "" {
			properties.SummaryMessage = &slack.AgentSummaryMessage{MessageTS: *summaryTS, ThreadTS: *summaryThreadTS}
		}
		options := []slack.SetAgentConversationPropertiesOption{slack.SetAgentConversationPropertiesOptionCodeChannel(properties)}
		if *resourceURL != "" {
			options = append(options, slack.SetAgentConversationPropertiesOptionAgentResource(slack.AgentResourceProperties{URL: resourceURL}))
		}
		response, err = api.SetAgentConversationPropertiesContext(ctx, *channelID, options...)
	case "commands":
		commands := []slack.AgentCommand{
			{Name: "review-status", Description: "Show the current review status"},
			{Name: "summarize", Description: "Summarize the session"},
		}
		if *clear {
			commands = nil
		}
		response, err = api.SetAgentConversationCommandsContext(ctx, *channelID, commands)
	case "set-view":
		response, err = setView(ctx, api, *channelID, *viewType, viewOptions...)
	case "list-views":
		response, err = api.ListAgentConversationViewsContext(ctx, *channelID)
	case "remove-view":
		option := slack.RemoveAgentConversationViewOptionViewKey(*viewKey)
		if *viewID != "" {
			option = slack.RemoveAgentConversationViewOptionViewID(*viewID)
		}
		response, err = api.RemoveAgentConversationViewContext(ctx, *channelID, option)
	case "create-canvas":
		var options []slack.CreateCanvasOption
		if *channelID != "" {
			options = append(options, slack.CreateCanvasOptionChannelID(*channelID))
		}
		var id string
		id, err = api.CreateCanvasContext(ctx, *name, slack.DocumentContent{Type: "markdown", Markdown: string(content)}, options...)
		response = map[string]string{"canvas_id": id}
	case "get-canvas":
		response, err = api.GetAgentConversationCanvasContext(ctx, *channelID, *canvasID,
			slack.GetAgentConversationCanvasOptionIncludeResolved(*includeResolved))
	case "set-canvas-content":
		response, err = api.SetAgentConversationCanvasContentContext(ctx, *channelID, *canvasID, string(content))
	case "archive":
		err = api.ArchiveAgentConversationContext(ctx, *channelID,
			slack.ArchiveAgentConversationOptionSummaryMessageTimestamp(*summaryTS))
		response = map[string]any{"channel_id": *channelID, "archived": true}
	}
	if err != nil {
		return fmt.Errorf("%s: %w", *action, err)
	}
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	return encoder.Encode(response)
}

func setView(ctx context.Context, api *slack.Client, channelID, viewType string, options ...slack.SetAgentConversationViewOption) (*slack.AgentConversationSetViewResponse, error) {
	response, err := api.SetAgentConversationViewContext(ctx, channelID, options...)
	var slackErr slack.SlackErrorResponse
	if viewType == slack.AgentViewTypePullRequest && errors.As(err, &slackErr) && slackErr.Err == "invalid_view_type" {
		return response, fmt.Errorf("%w; Slack rejected the documented pull_request type; publish the PR's unified diff with -type diff -content-file instead", err)
	}
	return response, err
}
