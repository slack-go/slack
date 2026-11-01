// Exercise the Enterprise admin Slackbot Skills API. See README.md for setup.
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
	"strings"
	"time"

	"github.com/slack-go/slack"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		if slackErr, ok := errors.AsType[slack.SlackErrorResponse](err); ok && slackErr.Err == "token_revoked" {
			fmt.Fprintln(os.Stderr, "SLACK_ADMIN_TOKEN was rejected as revoked. Reauthorize the org-level app installation and export its current org admin/owner user token, then retry.")
		}
		os.Exit(1)
	}
}

func run() error {
	action := flag.String("action", "", "Method: list, publish, set-auto-enabled, unpublish, workflow")
	skills := flag.String("skills", "", "Comma-separated skill IDs for publish/unpublish (1–50)")
	skill := flag.String("skill", "", "Single canvas skill ID for set-auto-enabled")
	permission := flag.String("permission", "", "Deployment audience: everyone, named_entities, no_one")
	groups := flag.String("usergroups", "", "Comma-separated group IDs; only for named_entities")
	category := flag.String("category", "", "Inventory category filter")
	creators := flag.String("creators", "", "Comma-separated creator U or W user IDs")
	published := flag.String("publish-status", "all", "all, published, unpublished")
	enabled := flag.String("auto-enabled-status", "all", "all, enabled, disabled")
	sort := flag.String("sort", "name", "name, usage_count, date_updated")
	direction := flag.String("sort-dir", "asc", "asc, desc")
	cursor := flag.String("cursor", "", "Inventory continuation cursor")
	limit := flag.Int("limit", 100, "Inventory page size: 1–200")
	all := flag.Bool("all", false, "Follow every inventory page, including empty pages")
	flag.Parse()
	provided := map[string]bool{}
	flag.Visit(func(f *flag.Flag) { provided[f.Name] = true })
	if flag.NArg() != 0 {
		return errors.New("use flags for all arguments; see -help")
	}
	if *action != "list" && *action != "publish" && *action != "set-auto-enabled" && *action != "unpublish" && *action != "workflow" {
		return errors.New("unknown -action; see -help")
	}
	if *action == "workflow" && (len(splitIDs(*groups)) < 1 || len(splitIDs(*groups)) > 100) {
		return errors.New("workflow requires 1–100 -usergroups for temporary deployment")
	}
	if (*action == "publish" || *action == "unpublish") && (len(splitIDs(*skills)) < 1 || len(splitIDs(*skills)) > 50) {
		return errors.New("publish/unpublish require 1–50 comma-separated -skills")
	}
	if *action == "set-auto-enabled" {
		if *skill == "" {
			return errors.New("set-auto-enabled requires -skill")
		}
		switch *permission {
		case "named_entities":
			if len(splitIDs(*groups)) < 1 || len(splitIDs(*groups)) > 100 {
				return errors.New("named_entities requires 1–100 -usergroups")
			}
		case "everyone", "no_one":
			if provided["usergroups"] {
				return errors.New("omit -usergroups for everyone/no_one")
			}
		default:
			return errors.New("-permission must be everyone, named_entities, or no_one")
		}
	}
	token := os.Getenv("SLACK_ADMIN_TOKEN")
	if token == "" {
		return errors.New("SLACK_ADMIN_TOKEN is required: org admin/owner user token from an org-level app installation")
	}
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()
	api := slack.New(token, slack.OptionHTTPClient(&http.Client{Timeout: 30 * time.Second}))
	if *action == "workflow" {
		return runWorkflow(ctx, api, splitIDs(*groups))
	}
	if *action == "publish" {
		response, err := api.AdminSlackbotSkillsPublishContext(ctx, splitIDs(*skills))
		return output(response, err)
	}
	if *action == "unpublish" {
		response, err := api.AdminSlackbotSkillsUnpublishContext(ctx, splitIDs(*skills))
		return output(response, err)
	}
	if *action == "set-auto-enabled" {
		options := []slack.AdminSlackbotSkillsSetAutoEnabledOption{}
		if provided["usergroups"] {
			options = append(options, slack.AdminSlackbotSkillsSetAutoEnabledOptionUsergroupIDs(splitIDs(*groups)...))
		}
		response, err := api.AdminSlackbotSkillsSetAutoEnabledContext(ctx, *skill, *permission, options...)
		return output(response, err)
	}
	options := []slack.AdminSlackbotSkillsListOption{}
	if provided["category"] {
		options = append(options, slack.AdminSlackbotSkillsListOptionCategory(*category))
	}
	if provided["creators"] {
		options = append(options, slack.AdminSlackbotSkillsListOptionCreatorIDs(splitIDs(*creators)...))
	}
	if provided["publish-status"] {
		options = append(options, slack.AdminSlackbotSkillsListOptionPublishStatus(*published))
	}
	if provided["auto-enabled-status"] {
		options = append(options, slack.AdminSlackbotSkillsListOptionAutoEnabledStatus(*enabled))
	}
	if provided["sort"] {
		options = append(options, slack.AdminSlackbotSkillsListOptionSort(*sort))
	}
	if provided["sort-dir"] {
		options = append(options, slack.AdminSlackbotSkillsListOptionSortDir(*direction))
	}
	if provided["limit"] {
		options = append(options, slack.AdminSlackbotSkillsListOptionLimit(*limit))
	}
	for {
		pageOptions := append([]slack.AdminSlackbotSkillsListOption{}, options...)
		if *cursor != "" || provided["cursor"] {
			pageOptions = append(pageOptions, slack.AdminSlackbotSkillsListOptionCursor(*cursor))
		}
		response, err := api.AdminSlackbotSkillsListContext(ctx, pageOptions...)
		if err := output(response, err); err != nil {
			return err
		}
		if !*all || response.ResponseMetadata.Cursor == "" {
			return nil
		}
		*cursor = response.ResponseMetadata.Cursor
	}
}

func splitIDs(value string) []string {
	ids := []string{}
	for id := range strings.SplitSeq(value, ",") {
		if id = strings.TrimSpace(id); id != "" {
			ids = append(ids, id)
		}
	}
	return ids
}

func output(response any, callErr error) error {
	if response != nil {
		encoder := json.NewEncoder(os.Stdout)
		encoder.SetIndent("", "  ")
		if err := encoder.Encode(response); err != nil {
			return errors.Join(callErr, err)
		}
	}
	return callErr
}
