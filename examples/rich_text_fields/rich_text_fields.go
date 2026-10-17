// Verify rich-text fields against Slack responses. See README.md for setup.
package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"maps"
	"net/http"
	"net/url"
	"os"
	"reflect"
	"slices"
	"strings"
	"time"

	"github.com/slack-go/slack"
)

type fieldCase struct {
	name    string
	element slack.RichTextSectionElement
	missing string
}

type settings struct {
	action, channel, timestamp, permalink, selected, user, tab, emoji, emojiURL string
	debug, printJSON, allowBroadcast                                            bool
}

var addedFields = map[string][]string{
	"broadcast": {"style"},
	"color":     {"style"},
	"date":      {"style", "timezone"},
	"channel":   {"tab_id", "from_llm"},
	"user":      {"from_llm"},
	"link":      {"unsafe", "from_llm", "is_slack_url", "truncated"},
	"emoji":     {"url"},
}

func main() {
	if err := run(); err != nil {
		printError(err)
		os.Exit(1)
	}
}

func run() error {
	var cfg settings
	flag.StringVar(&cfg.action, "action", "post", "post test messages or read one existing message")
	flag.StringVar(&cfg.channel, "channel", "", "Channel ID")
	flag.StringVar(&cfg.timestamp, "ts", "", "Message timestamp for -action read")
	flag.StringVar(&cfg.permalink, "url", "", "Message permalink for -action read")
	flag.StringVar(&cfg.selected, "case", "all", "Case name, all, or list")
	flag.StringVar(&cfg.user, "user", "", "User ID for user-from-llm")
	flag.StringVar(&cfg.tab, "tab", "", "Real channel tab ID for channel-tab")
	flag.StringVar(&cfg.emoji, "emoji", "", "Custom emoji name for emoji-url")
	flag.StringVar(&cfg.emojiURL, "emoji-url", "", "Public custom emoji asset URL for emoji-url")
	flag.BoolVar(&cfg.debug, "debug", false, "Enable Slack client debug logs")
	flag.BoolVar(&cfg.printJSON, "json", false, "Print outgoing blocks without connecting to Slack")
	flag.BoolVar(&cfg.allowBroadcast, "allow-broadcast", false, "Include broadcast-style in all; this posts an @here mention")
	flag.Parse()

	if cfg.selected == "list" {
		for _, c := range cases(cfg) {
			fmt.Println(c.name)
		}
		return nil
	}
	if cfg.action != "post" && cfg.action != "read" {
		return errors.New("-action must be post or read")
	}
	if cfg.permalink != "" {
		u, err := url.Parse(cfg.permalink)
		if err != nil {
			return err
		}
		parts := strings.Split(strings.Trim(u.Path, "/"), "/")
		if len(parts) != 3 || parts[0] != "archives" || !strings.HasPrefix(parts[2], "p") || len(parts[2]) < 8 {
			return errors.New("-url must be a Slack message permalink ending in /archives/CHANNEL/pTIMESTAMP")
		}
		digits := strings.TrimPrefix(parts[2], "p")
		cfg.channel, cfg.timestamp = parts[1], digits[:len(digits)-6]+"."+digits[len(digits)-6:]
	}
	if cfg.channel == "" {
		return errors.New("-channel or -url is required")
	}
	if cfg.action == "read" && cfg.timestamp == "" {
		return errors.New("-ts or -url is required for -action read")
	}

	logger := &responseLogger{client: &http.Client{Timeout: 30 * time.Second}}
	api := slack.New(os.Getenv("SLACK_BOT_TOKEN"), slack.OptionDebug(cfg.debug), slack.OptionHTTPClient(logger), slack.OptionRetry(2))
	if !cfg.printJSON && !strings.HasPrefix(os.Getenv("SLACK_BOT_TOKEN"), "xoxb-") {
		return errors.New("SLACK_BOT_TOKEN must contain a bot token (xoxb-)")
	}
	ctx := context.Background()
	if cfg.action == "read" {
		if cfg.printJSON {
			return errors.New("-json is only available with -action post")
		}
		_, err := api.GetConversationHistoryContext(ctx, &slack.GetConversationHistoryParameters{
			ChannelID: cfg.channel, Latest: cfg.timestamp, Oldest: cfg.timestamp, Inclusive: true, Limit: 1,
		})
		return errors.Join(err, errors.Join(logger.checkErrors...))
	}

	var failures []error
	matched := false
	for _, c := range cases(cfg) {
		if cfg.selected != "all" && cfg.selected != c.name {
			continue
		}
		matched = true
		if c.name == "broadcast-style" && cfg.selected == "all" && !cfg.allowBroadcast && !cfg.printJSON {
			log.Print("SKIP broadcast-style: use -allow-broadcast or -case broadcast-style to post an @here mention")
			continue
		}
		if c.missing != "" {
			log.Printf("SKIP %s: supply %s", c.name, c.missing)
			if cfg.selected != "all" {
				return fmt.Errorf("%s requires %s", c.name, c.missing)
			}
			continue
		}
		block := slack.NewRichTextBlock("field_"+c.name, slack.NewRichTextSection(c.element))
		encoded, err := json.Marshal(block)
		if err != nil {
			return err
		}
		log.Printf("%s outgoing block: %s", c.name, encoded)
		if cfg.printJSON {
			continue
		}
		logger.expected = c.element
		_, ts, err := api.PostMessageContext(ctx, cfg.channel,
			slack.MsgOptionText("Rich-text field check: "+c.name, false),
			slack.MsgOptionBlocks(block), slack.MsgOptionDisableLinkUnfurl(), slack.MsgOptionDisableMediaUnfurl())
		if err != nil {
			printError(fmt.Errorf("%s: %w", c.name, err))
			failures = append(failures, err)
			continue
		}
		log.Printf("ACCEPTED %s: channel=%s ts=%s. Check the returned fields and Slack UI separately.", c.name, cfg.channel, ts)
	}
	if !matched {
		return errors.New("unknown -case; use -case list")
	}
	return errors.Join(errors.Join(failures...), errors.Join(logger.checkErrors...))
}

func cases(cfg settings) []fieldCase {
	style := &slack.RichTextSectionTextStyle{Bold: true, Italic: true, Underline: true}
	missing := func(value, flagName string) string {
		if value == "" {
			return flagName
		}
		return ""
	}
	return []fieldCase{
		{"broadcast-style", &slack.RichTextSectionBroadcastElement{Type: slack.RTSEBroadcast, Range: "here", Style: style}, ""},
		{"color-style", &slack.RichTextSectionColorElement{Type: slack.RTSEColor, Value: "#36C5F0", Style: style}, ""},
		{"date-style", &slack.RichTextSectionDateElement{Type: slack.RTSEDate, Timestamp: 1720710212, Format: "{date_short} {time}", Style: style}, ""},
		{"date-timezone", &slack.RichTextSectionDateElement{Type: slack.RTSEDate, Timestamp: 1720710212, Format: "{date_short} {time}", Timezone: "Europe/Madrid"}, ""},
		{"channel-tab", &slack.RichTextSectionChannelElement{Type: slack.RTSEChannel, ChannelID: cfg.channel, TabID: cfg.tab}, missing(cfg.tab, "-tab")},
		{"channel-from-llm", &slack.RichTextSectionChannelElement{Type: slack.RTSEChannel, ChannelID: cfg.channel, FromLLM: new(true)}, ""},
		{"user-from-llm", &slack.RichTextSectionUserElement{Type: slack.RTSEUser, UserID: cfg.user, FromLLM: new(true)}, missing(cfg.user, "-user")},
		{"link-unsafe", &slack.RichTextSectionLinkElement{Type: slack.RTSELink, URL: "https://example.com", Unsafe: new(true)}, ""},
		{"link-from-llm", &slack.RichTextSectionLinkElement{Type: slack.RTSELink, URL: "https://example.com", FromLLM: new(true)}, ""},
		{"link-is-slack-url", &slack.RichTextSectionLinkElement{Type: slack.RTSELink, URL: "https://slack.com", IsSlackURL: new(true)}, ""},
		{"link-truncated", &slack.RichTextSectionLinkElement{Type: slack.RTSELink, URL: "https://example.com", Truncated: new(true)}, ""},
		{"emoji-url", &slack.RichTextSectionEmojiElement{Type: slack.RTSEEmoji, Name: cfg.emoji, URL: cfg.emojiURL}, missing(cfg.emoji, "-emoji") + missing(cfg.emojiURL, " -emoji-url")},
	}
}

type responseLogger struct {
	client      *http.Client
	expected    slack.RichTextSectionElement
	checkErrors []error
}

func (c *responseLogger) Do(req *http.Request) (*http.Response, error) {
	response, err := c.client.Do(req)
	if err != nil {
		return response, err
	}
	body, err := io.ReadAll(response.Body)
	response.Body.Close()
	if err != nil {
		return nil, err
	}
	log.Printf("raw %s response (%d): %s", req.URL.Path, response.StatusCode, body)
	response.Body = io.NopCloser(bytes.NewReader(body))
	var envelope struct {
		OK       bool              `json:"ok"`
		Message  json.RawMessage   `json:"message"`
		Messages []json.RawMessage `json:"messages"`
	}
	if err := json.Unmarshal(body, &envelope); err != nil || !envelope.OK {
		return response, nil
	}
	if len(envelope.Message) != 0 {
		envelope.Messages = append(envelope.Messages, envelope.Message)
	}
	if len(envelope.Messages) == 0 {
		log.Print("No message returned; there are no fields to verify.")
	}
	for _, message := range envelope.Messages {
		if err := inspect(message, c.expected); err != nil {
			log.Printf("DECODE CHECK FAILED: %v", err)
			c.checkErrors = append(c.checkErrors, err)
		}
	}
	return response, nil
}

func inspect(raw json.RawMessage, expected slack.RichTextSectionElement) error {
	var original struct {
		Blocks json.RawMessage `json:"blocks"`
	}
	if err := json.Unmarshal(raw, &original); err != nil {
		return err
	}
	if len(original.Blocks) == 0 {
		log.Print("No blocks returned; there are no fields to verify.")
		return nil
	}
	var blocks slack.Blocks
	if err := json.Unmarshal(original.Blocks, &blocks); err != nil {
		return err
	}
	encoded, err := json.Marshal(blocks)
	if err != nil {
		return err
	}
	before, after := collect(original.Blocks), collect(encoded)
	var failures []error
	for _, path := range sortedKeys(before) {
		log.Printf("DECODE %s: preserved=%t raw=%s decoded=%s", path, reflect.DeepEqual(before[path], after[path]), jsonValue(before[path]), jsonValue(after[path]))
		if !reflect.DeepEqual(before[path], after[path]) {
			failures = append(failures, fmt.Errorf("library changed returned field %s", path))
		}
	}
	if len(before) == 0 {
		log.Print("Slack returned none of the added rich-text fields. This does not verify their decoding.")
	}
	if expected != nil {
		wanted, err := json.Marshal([]any{map[string]any{"type": "rich_text", "elements": []any{map[string]any{"type": "rich_text_section", "elements": []any{expected}}}}})
		if err != nil {
			return err
		}
		sent := collect(wanted)
		for _, path := range sortedKeys(sent) {
			value := sent[path]
			actual, present := before[path]
			log.Printf("RETURNED %s: present=%t unchanged=%t sent=%s returned=%s", path, present, reflect.DeepEqual(value, actual), jsonValue(value), jsonValue(actual))
		}
	}
	return errors.Join(failures...)
}

func collect(raw []byte) map[string]any {
	var data any
	if err := json.Unmarshal(raw, &data); err != nil {
		return nil
	}
	result := make(map[string]any)
	var visit func(any, string)
	visit = func(value any, path string) {
		switch v := value.(type) {
		case []any:
			for i, item := range v {
				visit(item, fmt.Sprintf("%s[%d]", path, i))
			}
		case map[string]any:
			kind, _ := v["type"].(string)
			for _, field := range addedFields[kind] {
				if fieldValue, ok := v[field]; ok {
					result[path+"."+field] = fieldValue
				}
			}
			if elements, ok := v["elements"]; ok {
				visit(elements, path+".elements")
			}
		}
	}
	visit(data, "blocks")
	return result
}

func sortedKeys(values map[string]any) []string {
	return slices.Sorted(maps.Keys(values))
}

func jsonValue(value any) string {
	encoded, _ := json.Marshal(value)
	return string(encoded)
}

func printError(err error) {
	log.Print(err)
	if slackErr, ok := errors.AsType[slack.SlackErrorResponse](err); ok {
		log.Printf("Slack response messages: %v", slackErr.ResponseMetadata.Messages)
	}
}
