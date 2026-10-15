// Exercise the Slack Lists API in a real workspace. See README.md for setup
// and commands. Raw responses are printed before the library decodes them.
package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
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
		if slackErr, ok := errors.AsType[slack.SlackErrorResponse](err); ok {
			fmt.Fprintf(os.Stderr, "Slack response messages: %v\n", slackErr.ResponseMetadata.Messages)
		}
		os.Exit(1)
	}
}

func run() error {
	action := flag.String("action", "", "Action: create, rename, add, list, info, update, delete, delete-multiple, grant, revoke, export, download")
	listID := flag.String("list", "", "List ID (required except for create)")
	name := flag.String("name", "", "List name (create, rename)")
	itemID := flag.String("item", "", "Item ID (info, update, delete)")
	itemIDs := flag.String("items", "", "Comma-separated item IDs (delete-multiple)")
	columnID := flag.String("column", "", "Column ID (add, update)")
	text := flag.String("text", "", "Text cell value (add, update)")
	checked := flag.Bool("checked", false, "Checkbox cell value (update); use -checked=false to clear")
	limit := flag.Int("limit", 100, "Page size (list)")
	cursor := flag.String("cursor", "", "Pagination cursor (list)")
	all := flag.Bool("all", false, "Fetch all pages (list)")
	archived := flag.Bool("archived", false, "Read or export archived items (list, export)")
	userIDs := flag.String("user", "", "Comma-separated user IDs (grant, revoke)")
	channelIDs := flag.String("channel", "", "Comma-separated channel IDs (grant, revoke)")
	access := flag.String("access", "read", "Access level: read, write, owner (grant)")
	jobID := flag.String("job", "", "Export job ID (download)")
	format := flag.String("format", "csv", "Export format: csv, json (export, download)")
	threads := flag.Bool("threads", false, "Include item threads in JSON exports (export, download)")
	attachments := flag.Bool("attachments", false, "Include attachment metadata in JSON exports (export, download)")
	debug := flag.Bool("debug", false, "Enable Slack client debug logs")
	flag.Parse()
	provided := make(map[string]bool)
	flag.Visit(func(f *flag.Flag) { provided[f.Name] = true })
	if flag.NArg() != 0 {
		return errors.New("use flags for all arguments; see -help")
	}
	users, channels := splitIDs(*userIDs), splitIDs(*channelIDs)

	switch *action {
	case "create", "rename":
		if *name == "" {
			return errors.New("-name is required for create and rename")
		}
	case "add", "update":
		if *columnID == "" {
			return errors.New("-column is required for add and update")
		}
		if *action == "add" && (!provided["text"] || provided["checked"]) {
			return errors.New("add requires -text for a text column")
		}
		if *action == "update" && provided["text"] == provided["checked"] {
			return errors.New("update requires exactly one of -text or -checked")
		}
	case "list":
		if *limit < 1 {
			return errors.New("-limit must be positive")
		}
	case "info", "delete":
	case "delete-multiple":
		if len(splitIDs(*itemIDs)) == 0 {
			return errors.New("-items is required for delete-multiple")
		}
	case "grant", "revoke":
		if (len(users) == 0) == (len(channels) == 0) {
			return errors.New("provide exactly one of -user or -channel")
		}
		if *action == "grant" {
			switch slack.ListAccessLevel(*access) {
			case slack.ListAccessRead, slack.ListAccessWrite:
			case slack.ListAccessOwner:
				if len(channels) != 0 {
					return errors.New("owner access requires -user")
				}
			default:
				return errors.New("-access must be read, write, or owner")
			}
		}
	case "export", "download":
		if *format != "csv" && *format != "json" {
			return errors.New("-format must be csv or json")
		}
		if *format != "json" && (*threads || *attachments) {
			return errors.New("-threads and -attachments require -format json")
		}
		if *action == "download" && *jobID == "" {
			return errors.New("-job is required for download")
		}
	default:
		return errors.New("a valid -action is required; see -help")
	}
	if *action != "create" && *listID == "" {
		return errors.New("-list is required for this action")
	}
	if (*action == "info" || *action == "update" || *action == "delete") && *itemID == "" {
		return errors.New("-item is required for this action")
	}

	token := os.Getenv("SLACK_BOT_TOKEN")
	if token == "" {
		token = os.Getenv("SLACK_USER_TOKEN")
	}
	if token == "" {
		return errors.New("set SLACK_BOT_TOKEN or SLACK_USER_TOKEN")
	}
	api := slack.New(token,
		slack.OptionDebug(*debug),
		slack.OptionHTTPClient(responseLogger{client: &http.Client{Timeout: 30 * time.Second}}),
	)
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()
	fmt.Println("Compare raw and decoded output. A value present in raw JSON but missing from decoded output indicates a decoding mismatch.")

	switch *action {
	case "create":
		primary := true
		response, err := api.CreateListContext(ctx, *name,
			slack.CreateListOptionDescriptionBlocks(richText("Tasks created with slack-go")),
			slack.CreateListOptionSchema([]slack.ListColumn{
				{Key: "task", Name: "Task", Type: slack.ListColumnText, IsPrimaryColumn: &primary},
				{Key: "done", Name: "Done", Type: slack.ListColumnCheckbox},
			}),
		)
		if err != nil {
			return err
		}
		fmt.Printf("List ID: %s\n", response.ListID)
		for _, column := range response.ListMetadata.Schema {
			fmt.Printf("Column %s (%s): %s\n", column.Key, column.Type, column.ID)
		}
		return show("Decoded create response", response)
	case "rename":
		if err := api.UpdateListContext(ctx, *listID, slack.UpdateListOptionName(*name)); err != nil {
			return err
		}
	case "add":
		response, err := api.CreateListItemContext(ctx, *listID,
			slack.CreateListItemOptionInitialFields([]slack.ListItemField{
				{ColumnID: *columnID, RichText: richText(*text)},
			}),
		)
		if err != nil {
			return err
		}
		fmt.Printf("Item ID: %s\n", response.Item.ID)
		return show("Decoded item", response.Item)
	case "list":
		return listItems(ctx, api, *listID, *cursor, *limit, *archived, *all)
	case "info":
		response, err := api.GetListItemInfoContext(ctx, *listID, *itemID,
			slack.GetListItemInfoOptionIncludeIsSubscribed(true),
		)
		if err != nil {
			return err
		}
		fmt.Printf("List %s: %s; subtasks: %d\n", response.List.ID, response.List.Title, len(response.Subtasks))
		return show("Decoded item and subtasks", struct {
			Record   slack.ListItem   `json:"record"`
			Subtasks []slack.ListItem `json:"subtasks"`
		}{response.Record, response.Subtasks})
	case "update":
		field := slack.ListItemField{ColumnID: *columnID}
		if provided["checked"] {
			field.Checkbox = checked
		} else {
			field.RichText = richText(*text)
		}
		if err := api.UpdateListItemsContext(ctx, *listID, []slack.ListItemCell{
			{RowID: *itemID, ListItemField: field},
		}); err != nil {
			return err
		}
		response, err := api.GetListItemInfoContext(ctx, *listID, *itemID)
		if err != nil {
			return fmt.Errorf("update succeeded, but reading the item failed: %w", err)
		}
		return show("Decoded item after update", response.Record)
	case "delete":
		if err := api.DeleteListItemContext(ctx, *listID, *itemID); err != nil {
			return err
		}
	case "delete-multiple":
		if err := api.DeleteListItemsContext(ctx, *listID, splitIDs(*itemIDs)); err != nil {
			return err
		}
	case "grant":
		option := slack.SetListAccessOptionUserIDs(users)
		if len(channels) != 0 {
			option = slack.SetListAccessOptionChannelIDs(channels)
		}
		if err := api.SetListAccessContext(ctx, *listID, slack.ListAccessLevel(*access), option); err != nil {
			return err
		}
	case "revoke":
		option := slack.DeleteListAccessOptionUserIDs(users)
		if len(channels) != 0 {
			option = slack.DeleteListAccessOptionChannelIDs(channels)
		}
		if err := api.DeleteListAccessContext(ctx, *listID, option); err != nil {
			return err
		}
	case "export":
		options := []slack.StartListDownloadOption{
			slack.StartListDownloadOptionIncludeArchived(*archived),
			slack.StartListDownloadOptionFormat(slack.ListDownloadFormat(*format)),
		}
		if *format == "json" {
			options = append(options, slack.StartListDownloadOptionIncludeThreads(*threads),
				slack.StartListDownloadOptionIncludeAttachments(*attachments))
		}
		response, err := api.StartListDownloadContext(ctx, *listID, options...)
		if err != nil {
			return err
		}
		fmt.Printf("Job ID: %s\nRun -action download with this -job and the same -format, -threads, and -attachments flags.\n", response.JobID)
		return nil
	case "download":
		options := []slack.GetListDownloadOption{
			slack.GetListDownloadOptionFormat(slack.ListDownloadFormat(*format)),
		}
		if *format == "json" {
			options = append(options, slack.GetListDownloadOptionIncludeThreads(*threads),
				slack.GetListDownloadOptionIncludeAttachments(*attachments))
		}
		response, err := api.GetListDownloadContext(ctx, *listID, *jobID, options...)
		if err != nil {
			return err
		}
		fmt.Printf("Export status: %s\nDownload URL: %s\n", response.Status, response.DownloadURL)
		return nil
	}
	fmt.Printf("%s succeeded for list %s.\n", *action, *listID)
	return nil
}

func richText(text string) []slack.RichTextBlock {
	return []slack.RichTextBlock{*slack.NewRichTextBlock("",
		slack.NewRichTextSection(slack.NewRichTextSectionTextElement(text, nil)),
	)}
}

func listItems(ctx context.Context, api *slack.Client, listID, cursor string, limit int, archived, all bool) error {
	for {
		response, err := api.GetListItemsContext(ctx, listID,
			slack.GetListItemsOptionLimit(limit),
			slack.GetListItemsOptionCursor(cursor),
			slack.GetListItemsOptionArchived(archived),
			slack.GetListItemsOptionIncludeList(true),
		)
		if err != nil {
			return err
		}
		if response.List != nil {
			fmt.Printf("List %s: %s; total rows: %d\n", response.List.ID, response.List.Title, response.List.ListLimits.RowCount)
			for _, column := range response.List.ListMetadata.Schema {
				fmt.Printf("Column %s (%s): %s\n", column.Key, column.Type, column.ID)
			}
		}
		if err := show("Decoded items", response.Items); err != nil {
			return err
		}
		next := response.ResponseMetadata.Cursor
		fmt.Printf("Next cursor: %q\n", next)
		if !all || next == "" {
			return nil
		}
		if next == cursor {
			return errors.New("pagination cursor did not advance")
		}
		cursor = next
	}
}

func splitIDs(value string) []string {
	var ids []string
	for part := range strings.SplitSeq(value, ",") {
		if id := strings.TrimSpace(part); id != "" {
			ids = append(ids, id)
		}
	}
	return ids
}

func show(label string, value any) error {
	fmt.Println(label + ":")
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	return encoder.Encode(value)
}

type responseLogger struct {
	client *http.Client
}

func (c responseLogger) Do(request *http.Request) (*http.Response, error) {
	response, err := c.client.Do(request)
	if err != nil {
		return response, err
	}
	body, err := io.ReadAll(response.Body)
	response.Body.Close()
	if err != nil {
		return nil, err
	}
	fmt.Fprintf(os.Stderr, "Raw %s (%d): %s\n", request.URL.Path, response.StatusCode, body)
	response.Body = io.NopCloser(bytes.NewReader(body))
	return response, nil
}
