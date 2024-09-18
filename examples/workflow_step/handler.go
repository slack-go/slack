package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
)

const (
	IDSelectOptionBlock  = "select-option-block"
	IDExampleSelectInput = "example-select-input"
)

func handleMyWorkflowStep(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// see: https://github.com/slack-go/slack/blob/master/examples/eventsapi/events.go
	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	eventsAPIEvent, err := slackevents.ParseEvent(json.RawMessage(body), slackevents.OptionNoVerifyToken())
	if err != nil {
		log.Printf("[ERROR] Failed on parsing event: %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// see: https://api.slack.com/apis/connections/events-api#subscriptions
	if eventsAPIEvent.Type == slackevents.URLVerification {
		var r *slackevents.ChallengeResponse
		err := json.Unmarshal([]byte(body), &r)
		if err != nil {
			log.Printf("[ERROR] Failed to decode json message on event url_verification: %s", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "text")
		w.Write([]byte(r.Challenge))
		return
	}

	// see: https://api.slack.com/apis/connections/events-api#receiving_events
	if eventsAPIEvent.Type == slackevents.CallbackEvent {
		innerEvent := eventsAPIEvent.InnerEvent

		switch ev := innerEvent.Data.(type) {

		// see: https://api.slack.com/events/workflow_step_execute
		case *slackevents.WorkflowStepExecuteEvent:
			if ev.CallbackID == MyExampleWorkflowStepCallbackID {
				go doHeavyLoad(ev.WorkflowStep)

				w.WriteHeader(http.StatusOK)
				return
			}
			w.WriteHeader(http.StatusBadRequest)
			log.Printf("[WARN] unknown callbackID: %s", ev.CallbackID)
			return

		default:
			w.WriteHeader(http.StatusBadRequest)
			log.Printf("[WARN] unknown inner event type: %s", eventsAPIEvent.InnerEvent.Type)
			return
		}
	}

	w.WriteHeader(http.StatusBadRequest)
	log.Printf("[WARN] unknown event type: %s", eventsAPIEvent.Type)
}

func handleInteraction(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	jsonStr, err := url.QueryUnescape(string(body)[8:])
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var message slack.InteractionCallback
	if err := json.Unmarshal([]byte(jsonStr), &message); err != nil {
		log.Printf("[ERROR] Failed to decode json message from slack: %s", jsonStr)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	switch message.Type {
	case slack.InteractionTypeWorkflowStepEdit:
		// https://api.slack.com/workflows/steps#handle_config_view
		err := replyWithConfigurationView(message, "", "")
		if err != nil {
			log.Printf("[ERROR] Failed to open configuration modal in slack: %s", err.Error())
		}

	case slack.InteractionTypeViewSubmission:
		// https://api.slack.com/workflows/steps#handle_view_submission

		// process user inputs
		// this is just for demonstration, so we print it to console only
		blockAction := message.View.State.Values
		selectedOption := blockAction[IDSelectOptionBlock][IDExampleSelectInput].SelectedOption.Value
		log.Println(fmt.Sprintf("user selected: %s", selectedOption))

		in := &slack.WorkflowStepInputs{
			IDExampleSelectInput: slack.WorkflowStepInputElement{
				Value:                   selectedOption,
				SkipVariableReplacement: false,
			},
		}

		err := saveUserSettingsForWorkflowStep(message.WorkflowStep.WorkflowStepEditID, in, nil)
		if err != nil {
			log.Printf("[ERROR] Failed on doing a POST request to workflows.updateStep: %s", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
		}

	default:
		log.Printf("[WARN] unknown message type: %s", message.Type)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func replyWithConfigurationView(message slack.InteractionCallback, privateMetaData string, externalID string) error {
	headerText := slack.NewTextBlockObject("mrkdwn", "Hello World!\nThis is your workflow step app configuration view", false, false)
	headerSection := slack.NewSectionBlock(headerText, nil, nil)

	options := []*slack.OptionBlockObject{}
	options = append(
		options,
		slack.NewOptionBlockObject("one", slack.NewTextBlockObject("plain_text", "One", false, false), nil),
	)

	options = append(
		options,
		slack.NewOptionBlockObject("two", slack.NewTextBlockObject("plain_text", "Two", false, false), nil),
	)

	options = append(
		options,
		slack.NewOptionBlockObject("three", slack.NewTextBlockObject("plain_text", "Three", false, false), nil),
	)

	selection := slack.NewOptionsSelectBlockElement(
		"static_select",
		slack.NewTextBlockObject("plain_text", "your choice", false, false),
		IDExampleSelectInput,
		options...,
	)

	// preselect option, if workflow step input is defined
	initialOption, ok := slack.GetInitialOptionFromWorkflowStepInput(selection, message.WorkflowStep.Inputs, options)
	if ok {
		selection.InitialOption = initialOption
	}

	inputBlock := slack.NewInputBlock(
		IDSelectOptionBlock,
		slack.NewTextBlockObject("plain_text", "Select an option", false, false),
		selection,
	)

	blocks := slack.Blocks{
		BlockSet: []slack.Block{
			headerSection,
			inputBlock,
		},
	}

	cmr := slack.NewConfigurationModalRequest(blocks, privateMetaData, externalID)
	_, err := appCtx.slack.OpenView(message.TriggerID, cmr.ModalViewRequest)
	return err
}

func saveUserSettingsForWorkflowStep(workflowStepEditID string, inputs *slack.WorkflowStepInputs, outputs *[]slack.WorkflowStepOutput) error {
	return appCtx.slack.SaveWorkflowStepConfiguration(workflowStepEditID, inputs, outputs)
}

func doHeavyLoad(workflowStep slackevents.EventWorkflowStep) {
	// process user configuration e.g. inputs
	log.Printf("Inputs:")
	for name, input := range *workflowStep.Inputs {
		log.Printf(fmt.Sprintf("%s: %s", name, input.Value))
	}

	// do heavy load
	time.Sleep(10 * time.Second)
	log.Println("Done")
}
