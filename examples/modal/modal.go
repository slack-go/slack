// Modal example - How to respond to a slash command with an interactive modal and parse the response
// The flow of this example:
// 1. User trigers your app with a slash command (e.g. /modaltest) that will send a request to http://URL/slash and respond with a request to open a modal
// 2. User fills out fields first and last name in modal and hits submit
// 3. This will send a request to http://URL/modal and send a greeting message to the user

// Note: Within your slack app you will need to enable and provide a URL for "Interactivity & Shortcuts" and "Slash Commands"
// Note: Be sure to update YOUR_SIGNING_SECRET_HERE and YOUR_TOKEN_HERE
// You can use ngrok to test this example: https://api.slack.com/tutorials/tunneling-with-ngrok
// Helpful slack documentation to learn more: https://api.slack.com/interactivity/handling

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/slack-go/slack"
	"time"
)

func generateModalRequest() slack.ModalViewRequest {
	// Create a ModalViewRequest with a header and two inputs
	titleText := slack.NewTextBlockObject("plain_text", "My App", false, false)
	closeText := slack.NewTextBlockObject("plain_text", "Close", false, false)
	submitText := slack.NewTextBlockObject("plain_text", "Submit", false, false)

	headerText := slack.NewTextBlockObject("mrkdwn", "Please enter your name", false, false)
	headerSection := slack.NewSectionBlock(headerText, nil, nil)

	firstNameText := slack.NewTextBlockObject("plain_text", "First Name", false, false)
	firstNameHint := slack.NewTextBlockObject("plain_text", "First Name Hint", false, false)
	firstNamePlaceholder := slack.NewTextBlockObject("plain_text", "Enter your first name", false, false)
	firstNameElement := slack.NewPlainTextInputBlockElement(firstNamePlaceholder, "firstName")
	// Notice that blockID is a unique identifier for a block
	firstName := slack.NewInputBlock("First Name", firstNameText, firstNameHint, firstNameElement)

	lastNameText := slack.NewTextBlockObject("plain_text", "Last Name", false, false)
	lastNameHint := slack.NewTextBlockObject("plain_text", "Last Name Hint", false, false)
	lastNamePlaceholder := slack.NewTextBlockObject("plain_text", "Enter your first name", false, false)
	lastNameElement := slack.NewPlainTextInputBlockElement(lastNamePlaceholder, "lastName")
	lastName := slack.NewInputBlock("Last Name", lastNameText, lastNameHint, lastNameElement)

	blocks := slack.Blocks{
		BlockSet: []slack.Block{
			headerSection,
			firstName,
			lastName,
		},
	}

	var modalRequest slack.ModalViewRequest
	modalRequest.Type = slack.ViewType("modal")
	modalRequest.Title = titleText
	modalRequest.Close = closeText
	modalRequest.Submit = submitText
	modalRequest.Blocks = blocks
	return modalRequest
}

func updateModal() slack.ModalViewRequest {
	// Create a ModalViewRequest with a header and two inputs
	titleText := slack.NewTextBlockObject("plain_text", "My App", false, false)
	closeText := slack.NewTextBlockObject("plain_text", "Close", false, false)
	submitText := slack.NewTextBlockObject("plain_text", "Submit", false, false)

	headerText := slack.NewTextBlockObject("mrkdwn", "Modal updated!", false, false)
	headerSection := slack.NewSectionBlock(headerText, nil, nil)

	blocks := slack.Blocks{
		BlockSet: []slack.Block{
			headerSection,
		},
	}

	var modalRequest slack.ModalViewRequest
	modalRequest.Type = slack.ViewType("modal")
	modalRequest.Title = titleText
	modalRequest.Close = closeText
	modalRequest.Submit = submitText
	modalRequest.Blocks = blocks
	return modalRequest
}

// This was taken from the slash example
// https://github.com/slack-go/slack/blob/master/examples/slash/slash.go
func verifySigningSecret(r *http.Request) error {
	signingSecret := "YOUR_SIGNING_SECRET_HERE"
	verifier, err := slack.NewSecretsVerifier(r.Header, signingSecret)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	// Need to use r.Body again when unmarshalling SlashCommand and InteractionCallback
	r.Body = io.NopCloser(bytes.NewBuffer(body))

	verifier.Write(body)
	if err = verifier.Ensure(); err != nil {
		fmt.Println(err.Error())
		return err
	}

	return nil
}

func handleSlash(w http.ResponseWriter, r *http.Request) {

	err := verifySigningSecret(r)
	if err != nil {
		fmt.Printf(err.Error())
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	s, err := slack.SlashCommandParse(r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println(err.Error())
		return
	}

	switch s.Command {
	case "/slash":
		api := slack.New("YOUR_TOKEN_HERE")
		modalRequest := generateModalRequest()
		_, err = api.OpenView(s.TriggerID, modalRequest)
		if err != nil {
			fmt.Printf("Error opening view: %s", err)
		}
	default:
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func handleModal(w http.ResponseWriter, r *http.Request) {

	err := verifySigningSecret(r)
	if err != nil {
		fmt.Printf(err.Error())
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	var i slack.InteractionCallback
	err = json.Unmarshal([]byte(r.FormValue("payload")), &i)
	if err != nil {
		fmt.Printf(err.Error())
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	api := slack.New("YOUR_TOKEN_HERE")
	if err != nil {
		fmt.Printf(err.Error())
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// update modal sample
	switch i.Type {
	//update when interaction type is view_submission
	case slack.InteractionTypeViewSubmission:
		//you can use any modal you want to show to users just like creating modal.
		updateModal := updateModal()
		// You must set one of external_id or view_id and you can use hash for avoiding race condition.
		// More details: https://api.slack.com/surfaces/modals/using#updating_apis
		_, err := api.UpdateView(updateModal, "", i.View.Hash, i.View.ID)
		// Wait for a few seconds to see result this code is necesarry due to slack server modal is going to be closed after the update
		time.Sleep(time.Second * 2)
		if err != nil {
			fmt.Printf("Error updating view: %s", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}

func main() {
	http.HandleFunc("/slash", handleSlash)
	http.HandleFunc("/modal", handleModal)
	http.ListenAndServe(":4390", nil)
}
