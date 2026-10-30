// Exercise individual Slackbot Skills methods or their complete lifecycle.
package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/slack-go/slack"
)

const exampleInstructions = "# Release review\n\nReview the supplied release plan. Identify the three largest risks and suggest one concrete follow-up for each. Ask for the plan when it is missing."

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		var slackErr slack.SlackErrorResponse
		if errors.As(err, &slackErr) && slackErr.Err == "token_revoked" {
			fmt.Fprintln(os.Stderr, "SLACK_USER_TOKEN was rejected as revoked. Reauthorize the app and export its current User OAuth Token with skills:read and skills:write, then retry.")
		}
		os.Exit(1)
	}
}

func run() error {
	action := flag.String("action", "", "Method: create, info, update, list, library-add, library-remove, delete, workflow, mcp-workflow")
	skill := flag.String("skill", "", "Skill ID (or builtin: reference for info)")
	name := flag.String("name", "", "Skill display name")
	instructions := flag.String("instructions", exampleInstructions, "Markdown instructions; create defaults to a built-in example, update changes only when supplied")
	canvas := flag.String("canvas", "", "Owned canvas to convert instead of Markdown")
	description := flag.String("description", "", "Description; an explicit empty value clears it on update")
	category := flag.String("category", "", "Skill category")
	prompt := flag.String("example-prompt", "", "Example prompt; may contain {skill}")
	mcpFile := flag.String("mcp-servers-file", "", "JSON array using Slack's MCP request schema; [] clears")
	mcpJSON := flag.String("mcp-servers", "", "Inline JSON array of MCP declarations; [] clears")
	mcpApp := flag.String("mcp-app", "", "Installed MCP app ID to verify (mcp-workflow)")
	filter := flag.String("filter", "library", "List: library, owned, catalog")
	cursor := flag.String("cursor", "", "List continuation cursor")
	limit := flag.Int("limit", 100, "List page size")
	all := flag.Bool("all", false, "List all pages, including empty pages with continuation")
	flag.Parse()
	provided := map[string]bool{}
	flag.Visit(func(f *flag.Flag) { provided[f.Name] = true })
	if flag.NArg() != 0 {
		return errors.New("use flags for all arguments; see -help")
	}
	switch *action {
	case "create", "list", "workflow":
	case "mcp-workflow":
		if *mcpApp == "" {
			return errors.New("mcp-workflow requires -mcp-app with an installed, usable MCP app ID")
		}
	case "info", "update", "library-add", "library-remove", "delete":
		if *skill == "" {
			return errors.New("this action requires -skill")
		}
	default:
		return errors.New("unknown -action; see -help")
	}
	if *action == "create" && *name == "" {
		return errors.New("create requires -name")
	}
	if provided["instructions"] && provided["canvas"] {
		return errors.New("provide only one of -instructions or -canvas")
	}
	if provided["mcp-servers"] && provided["mcp-servers-file"] {
		return errors.New("provide only one of -mcp-servers or -mcp-servers-file")
	}
	var servers []json.RawMessage
	if provided["mcp-servers-file"] || provided["mcp-servers"] {
		body := []byte(*mcpJSON)
		if provided["mcp-servers-file"] {
			var err error
			body, err = os.ReadFile(*mcpFile)
			if err != nil {
				return err
			}
		}
		if err := json.Unmarshal(body, &servers); err != nil {
			return err
		}
		if servers == nil {
			return errors.New("MCP declarations must be a JSON array, not null")
		}
	}
	token := os.Getenv("SLACK_USER_TOKEN")
	if token == "" {
		return errors.New("SLACK_USER_TOKEN is required (user token with skills:read and skills:write)")
	}
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()
	api := slack.New(token, slack.OptionHTTPClient(&http.Client{Timeout: 30 * time.Second}))
	if *action == "mcp-workflow" {
		return runMCPWorkflow(ctx, api, *mcpApp, servers)
	}
	if *action == "workflow" {
		return runWorkflow(ctx, api)
	}
	if *action == "library-add" {
		return output(nil, api.AddSlackbotSkillToLibraryContext(ctx, *skill))
	}
	if *action == "library-remove" {
		return output(nil, api.RemoveSlackbotSkillFromLibraryContext(ctx, *skill))
	}
	if *action == "delete" {
		return output(nil, api.DeleteSlackbotSkillContext(ctx, *skill))
	}
	if *action == "list" {
		options := []slack.ListSlackbotSkillsOption{}
		if provided["filter"] {
			options = append(options, slack.ListSlackbotSkillsOptionFilter(*filter))
		}
		if provided["category"] {
			options = append(options, slack.ListSlackbotSkillsOptionCategory(*category))
		}
		if provided["limit"] {
			options = append(options, slack.ListSlackbotSkillsOptionLimit(*limit))
		}
		for {
			pageOptions := append([]slack.ListSlackbotSkillsOption{}, options...)
			if *cursor != "" || provided["cursor"] {
				pageOptions = append(pageOptions, slack.ListSlackbotSkillsOptionCursor(*cursor))
			}
			response, err := api.ListSlackbotSkillsContext(ctx, pageOptions...)
			if err := output(response, err); err != nil {
				return err
			}
			if !*all || response.ResponseMetadata.Cursor == "" {
				return nil
			}
			*cursor = response.ResponseMetadata.Cursor
		}
	}
	if *action == "info" {
		response, err := api.GetSlackbotSkillInfoContext(ctx, *skill)
		return output(response, err)
	}
	if *action == "update" {
		options := []slack.UpdateSlackbotSkillOption{}
		if provided["name"] {
			options = append(options, slack.UpdateSlackbotSkillOptionDisplayName(*name))
		}
		if provided["instructions"] {
			options = append(options, slack.UpdateSlackbotSkillOptionInstructions(*instructions))
		}
		if provided["description"] {
			options = append(options, slack.UpdateSlackbotSkillOptionDescription(*description))
		}
		if provided["category"] {
			options = append(options, slack.UpdateSlackbotSkillOptionCategory(*category))
		}
		if provided["example-prompt"] {
			options = append(options, slack.UpdateSlackbotSkillOptionExamplePrompt(*prompt))
		}
		if provided["mcp-servers-file"] || provided["mcp-servers"] {
			options = append(options, slack.UpdateSlackbotSkillOptionMCPServers(servers...))
		}
		if len(options) == 0 {
			return errors.New("update requires at least one field flag")
		}
		response, err := api.UpdateSlackbotSkillContext(ctx, *skill, options...)
		return output(response, err)
	}
	options := []slack.CreateSlackbotSkillOption{}
	if !provided["canvas"] {
		options = append(options, slack.CreateSlackbotSkillOptionInstructions(*instructions))
	}
	if provided["canvas"] {
		options = append(options, slack.CreateSlackbotSkillOptionCanvasID(*canvas))
	}
	if provided["description"] {
		options = append(options, slack.CreateSlackbotSkillOptionDescription(*description))
	}
	if provided["category"] {
		options = append(options, slack.CreateSlackbotSkillOptionCategory(*category))
	}
	if provided["example-prompt"] {
		options = append(options, slack.CreateSlackbotSkillOptionExamplePrompt(*prompt))
	}
	if provided["mcp-servers-file"] || provided["mcp-servers"] {
		options = append(options, slack.CreateSlackbotSkillOptionMCPServers(servers...))
	}
	response, err := api.CreateSlackbotSkillContext(ctx, *name, options...)
	return output(response, err)
}

func output(response any, callErr error) error {
	if response == nil && callErr == nil {
		response = map[string]bool{"ok": true}
	}
	if response != nil {
		encoder := json.NewEncoder(os.Stdout)
		encoder.SetIndent("", "  ")
		if err := encoder.Encode(response); err != nil {
			return errors.Join(callErr, err)
		}
	}
	return callErr
}
