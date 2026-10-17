// Run this example with the app configuration and manual checks in README.md.
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
	"net/http"
	"os"
	"os/signal"
	"reflect"
	"regexp"
	"slices"
	"strings"
	"time"

	"github.com/slack-go/slack"
	"github.com/slack-go/slack/socketmode"
)

const (
	buttonAction   = "field_check_button"
	agentAction    = "field_check_agent"
	checkboxAction = "field_check_checkboxes"
	launchAction   = "field_check_open_modal"
	selectAction   = "field_check_conversations"
	inputAction    = "field_check_rich_text"
	modalCallback  = "block_element_fields"
)

func main() {
	channel := flag.String("channel", "", "Channel ID in which to post the test message (required)")
	user := flag.String("user", "", "User ID allowed to see the test button and checkboxes")
	which := flag.String("case", "all", "Case: all, button, agent, checkboxes, or modal")
	debug := flag.Bool("debug", false, "Enable additional library diagnostics")
	flag.Parse()
	if err := run(*channel, *user, *which, *debug); err != nil {
		printError(err)
		os.Exit(1)
	}
}

func run(channel, user, which string, debug bool) (runErr error) {
	if channel == "" {
		return errors.New("-channel is required")
	}
	if !slices.Contains([]string{"all", "button", "agent", "checkboxes", "modal"}, which) {
		return errors.New("-case must be all, button, agent, checkboxes, or modal")
	}
	if user == "" && slices.Contains([]string{"all", "button", "checkboxes"}, which) {
		return errors.New("-user is required for the all, button, and checkboxes cases")
	}
	botToken, appToken := os.Getenv("SLACK_BOT_TOKEN"), os.Getenv("SLACK_APP_TOKEN")
	if !strings.HasPrefix(botToken, "xoxb-") || !strings.HasPrefix(appToken, "xapp-") {
		return errors.New("SLACK_BOT_TOKEN (xoxb-) and SLACK_APP_TOKEN (xapp-) are required")
	}
	logger := payloadLogger{debug: debug, secrets: []string{botToken, appToken}}
	failures := make(chan error, 1)
	defer func() {
		select {
		case err := <-failures:
			runErr = errors.Join(runErr, err)
		default:
		}
	}()
	api := slack.New(botToken,
		slack.OptionAppLevelToken(appToken),
		slack.OptionDebug(debug),
		slack.OptionLog(logger),
		slack.OptionHTTPClient(&http.Client{Timeout: 20 * time.Second, Transport: responseLogger{logger: logger, failures: failures}}),
	)
	// Socket Mode has no raw-message callback. Its logger receives each envelope before parsing.
	client := socketmode.New(api, socketmode.OptionDebug(true), socketmode.OptionLog(logger))
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()
	done := make(chan error, 1)
	go func() { done <- client.RunContext(ctx) }()
	posted := false
	for {
		select {
		case <-ctx.Done():
			return nil
		case err := <-failures:
			return err
		case err := <-done:
			if errors.Is(err, context.Canceled) {
				return nil
			}
			return err
		case event, ok := <-client.Events:
			if !ok {
				return errors.New("socket mode event channel closed")
			}
			log.Printf("Socket Mode event: %s", event.Type)
			if event.Request != nil && event.Request.EnvelopeID != "" {
				if err := client.Ack(*event.Request); err != nil {
					printError(fmt.Errorf("acknowledge %s: %w", event.Type, err))
				}
			}
			switch event.Type {
			case socketmode.EventTypeInvalidAuth:
				return errors.New("socket mode authentication failed")
			case socketmode.EventTypeConnected:
				if posted {
					continue
				}
				blocks := messageBlocks(which, user)
				printJSON("outgoing blocks", slack.Blocks{BlockSet: blocks})
				_, ts, err := api.PostMessageContext(ctx, channel,
					slack.MsgOptionText("Block element field verification", false), slack.MsgOptionBlocks(blocks...))
				if err != nil {
					return fmt.Errorf("post %s case (use -case to isolate rejected elements): %w", which, err)
				}
				posted = true
				log.Printf("Slack accepted the message: channel=%s ts=%s. Field retention and UI behavior are separate checks. Follow README.md; Ctrl-C stops listening.", channel, ts)
			case socketmode.EventTypeInteractive:
				callback, ok := event.Data.(slack.InteractionCallback)
				if !ok {
					log.Printf("Unexpected interactive payload type: %T", event.Data)
					continue
				}
				if event.Request != nil {
					if err := inspectResponse("interaction", event.Request.Payload); err != nil {
						return err
					}
				}
				if err := handleInteraction(ctx, api, callback, channel); err != nil {
					printError(err)
				}
			}
		}
	}
}

func plain(text string) *slack.TextBlockObject {
	return slack.NewTextBlockObject(slack.PlainTextType, text, false, false)
}

func messageBlocks(which, user string) []slack.Block {
	blocks := []slack.Block{slack.NewSectionBlock(plain("Block element field verification. See the example README for each check."), nil, nil)}
	if which == "all" || which == "button" {
		button := slack.NewButtonBlockElement(buttonAction, "clicked", plain("Visible test button"))
		button.AccessibilityLabel = "Verify the button accessibility label"
		button.VisibleToUserIDs = []string{user}
		blocks = append(blocks, slack.NewActionBlock("button_fields", button))
	}
	if which == "all" || which == "agent" {
		button := slack.NewButtonBlockElement(agentAction, "fallback", plain("Test Slackbot prompt"))
		button.AgentPrompt = "Reply with the words: Block element fields verified. This is a test; do not take any other action."
		button.AgentPromptDisplay = "Verify Block Kit prompt fields"
		blocks = append(blocks, slack.NewActionBlock("agent_fields", button))
	}
	if which == "all" || which == "checkboxes" {
		checkboxes := slack.NewCheckboxGroupsBlockElement(checkboxAction,
			slack.NewOptionBlockObject("checked", plain("Visible test checkbox"), nil))
		checkboxes.VisibleToUserIDs = []string{user}
		blocks = append(blocks, slack.NewActionBlock("checkbox_fields", checkboxes))
	}
	if which == "all" || which == "modal" {
		blocks = append(blocks, slack.NewActionBlock("modal_launch",
			slack.NewButtonBlockElement(launchAction, "open", plain("Open input tests"))))
	}
	return blocks
}

func handleInteraction(ctx context.Context, api *slack.Client, callback slack.InteractionCallback, channel string) error {
	log.Printf("Interaction: type=%s user=%s", callback.Type, callback.User.ID)
	for _, action := range callback.ActionCallback.BlockActions {
		log.Printf("Action: %s", action.ActionID)
		if action.ActionID == agentAction {
			log.Print("Agent button used its standard interaction fallback. This does not verify the Slackbot AI path.")
		}
		if action.ActionID == launchAction {
			view := modal()
			printJSON("outgoing modal", view)
			if _, err := api.OpenViewContext(ctx, callback.TriggerID, view); err != nil {
				return fmt.Errorf("open modal: %w", err)
			}
			log.Printf("Slack accepted the modal. Before changing the selector, check that %s is selected. Type more than four lines and check scrolling.", channel)
		}
	}
	if callback.Type == slack.InteractionTypeViewSubmission && callback.View.CallbackID == modalCallback {
		if callback.View.State == nil {
			return errors.New("modal submission has no view.state")
		}
		selected := callback.View.State.Values["conversation_input"][selectAction].SelectedConversations
		log.Printf("Submitted conversations=%v; includes launch channel=%t. This checks the default only if you left the selector unchanged.", selected, slices.Contains(selected, channel))
		printJSON("submitted rich text", callback.View.State.Values["rich_text_input"][inputAction].RichTextValue)
	}
	return nil
}

func modal() slack.ModalViewRequest {
	conversations := slack.NewOptionsMultiSelectBlockElement(slack.MultiOptTypeConversations, plain("Current conversation"), selectAction)
	conversations.DefaultToCurrentConversation = true
	input := slack.NewRichTextInputBlockElement(plain("Enter more than four lines"), inputAction)
	input.MinLines, input.MaxLines = 2, 4
	return slack.ModalViewRequest{
		Type: slack.VTModal, Title: plain("Block element fields"), Submit: plain("Submit"), Close: plain("Cancel"), CallbackID: modalCallback,
		Blocks: slack.Blocks{BlockSet: []slack.Block{
			slack.NewInputBlock("conversation_input", plain("Leave the default selected"), nil, conversations).WithOptional(true),
			slack.NewInputBlock("rich_text_input", plain("Two to four visible lines"), nil, input).WithOptional(true),
		}},
	}
}

type responseLogger struct {
	logger   payloadLogger
	failures chan<- error
}

func (t responseLogger) RoundTrip(req *http.Request) (*http.Response, error) {
	resp, err := http.DefaultTransport.RoundTrip(req)
	if err != nil {
		return nil, err
	}
	// apps.connections.open returns a credential-bearing WebSocket URL.
	if req.URL.Path != "/api/chat.postMessage" && req.URL.Path != "/api/views.open" {
		return resp, nil
	}
	body, err := io.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return nil, err
	}
	resp.Body = io.NopCloser(bytes.NewReader(body))
	log.Printf("raw %s response (secrets redacted): %s", req.URL.Path, t.logger.redact(string(body)))
	if err := inspectResponse(req.URL.Path, body); err != nil {
		// Keep the response intact so a diagnostic failure cannot look like a failed API call.
		select {
		case t.failures <- err:
		default:
		}
	}
	return resp, nil
}

var secretJSON = regexp.MustCompile(`("(?:token|response_url|trigger_id)"\s*:\s*)"(?:\\.|[^"\\])*"`)
var websocketURL = regexp.MustCompile(`wss://[^\s"<>]+`)

type payloadLogger struct {
	debug   bool
	secrets []string
}

func (l payloadLogger) redact(message string) string {
	for _, secret := range l.secrets {
		message = strings.ReplaceAll(message, secret, "[REDACTED]")
	}
	message = secretJSON.ReplaceAllString(message, `${1}"[REDACTED]"`)
	return websocketURL.ReplaceAllString(message, "[REDACTED_WEBSOCKET_URL]")
}

func (l payloadLogger) Output(_ int, message string) error {
	if strings.HasPrefix(message, "Received WebSocket message:") {
		log.Printf("raw Socket Mode envelope (secrets redacted): %s", l.redact(strings.TrimPrefix(message, "Received WebSocket message: ")))
	} else if l.debug {
		log.Print(l.redact(message))
	}
	return nil
}

func inspectResponse(source string, body []byte) error {
	var envelope struct {
		Message struct{ Blocks json.RawMessage } `json:"message"`
		View    struct{ Blocks json.RawMessage } `json:"view"`
	}
	if err := json.Unmarshal(body, &envelope); err != nil {
		return fmt.Errorf("cannot inspect %s: %w", source, err)
	}
	var failures []error
	for _, raw := range []json.RawMessage{envelope.Message.Blocks, envelope.View.Blocks} {
		if len(raw) == 0 {
			continue
		}
		var blocks slack.Blocks
		if err := json.Unmarshal(raw, &blocks); err != nil {
			failures = append(failures, fmt.Errorf("decode blocks from %s: %w", source, err))
			continue
		}
		decoded, err := json.Marshal(blocks)
		if err != nil {
			failures = append(failures, fmt.Errorf("encode blocks from %s: %w", source, err))
			continue
		}
		var before, after any
		if err := json.Unmarshal(raw, &before); err != nil {
			failures = append(failures, fmt.Errorf("cannot inspect raw blocks from %s: %w", source, err))
			continue
		}
		if err := json.Unmarshal(decoded, &after); err != nil {
			failures = append(failures, fmt.Errorf("cannot inspect decoded blocks from %s: %w", source, err))
			continue
		}
		checks := []struct {
			action string
			fields []string
		}{
			{buttonAction, []string{"accessibility_label", "visible_to_user_ids"}},
			{agentAction, []string{"agent_prompt", "agent_prompt_display"}},
			{checkboxAction, []string{"visible_to_user_ids"}},
			{selectAction, []string{"default_to_current_conversation"}},
			{inputAction, []string{"min_lines", "max_lines"}},
		}
		for _, check := range checks {
			sent, got := findElement(before, check.action), findElement(after, check.action)
			if sent == nil {
				continue
			}
			for _, field := range check.fields {
				value, present := sent[field]
				if !present {
					log.Printf("%s %s.%s: Slack omitted the field; decoder preservation is unverified", source, check.action, field)
				} else {
					preserved := reflect.DeepEqual(value, got[field])
					log.Printf("%s %s.%s: Slack returned %v; Go preserved=%t", source, check.action, field, value, preserved)
					if !preserved {
						failures = append(failures, fmt.Errorf("%s %s.%s: Go did not preserve Slack's returned value", source, check.action, field))
					}
				}
			}
		}
	}
	return errors.Join(failures...)
}

func findElement(value any, action string) map[string]any {
	switch value := value.(type) {
	case map[string]any:
		if value["action_id"] == action {
			return value
		}
		for _, child := range value {
			if found := findElement(child, action); found != nil {
				return found
			}
		}
	case []any:
		for _, child := range value {
			if found := findElement(child, action); found != nil {
				return found
			}
		}
	}
	return nil
}

func printJSON(label string, value any) {
	body, err := json.Marshal(value)
	if err != nil {
		log.Printf("%s: %v", label, err)
		return
	}
	log.Printf("%s: %s", label, body)
}

func printError(err error) {
	log.Print(err)
	if response, ok := errors.AsType[slack.SlackErrorResponse](err); ok {
		for _, message := range response.ResponseMetadata.Messages {
			log.Printf("Slack response_metadata.messages: %s", message)
		}
	}
}
