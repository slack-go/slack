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
	"io/ioutil"
	"net/http"

	"github.com/slack-go/slack"
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

// This was taken from the slash example
// https://github.com/slack-go/slack/blob/master/examples/slash/slash.go
func verifySigningSecret(r *http.Request) error {
	signingSecret := "YOUR_SIGNING_SECRET_HERE"
	verifier, err := slack.NewSecretsVerifier(r.Header, signingSecret)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	// Need to use r.Body again when unmarshalling SlashCommand and InteractionCallback
	r.Body = ioutil.NopCloser(bytes.NewBuffer(body))

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
	case "/humboldttest":
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

	// Note there might be a better way to get this info, but I figured this structure out from looking at the json response
	firstName := i.View.State.Values["First Name"]["firstName"].Value
	lastName := i.View.State.Values["Last Name"]["lastName"].Value

	msg := fmt.Sprintf("Hello %s %s, nice to meet you!", firstName, lastName)

	api := slack.New("YOUR_TOKEN_HERE")
	_, _, err = api.PostMessage(i.User.ID,
		slack.MsgOptionText(msg, false),
		slack.MsgOptionAttachments())
	if err != nil {
		fmt.Printf(err.Error())
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
}

func main() {
	http.HandleFunc("/slash", handleSlash)
	http.HandleFunc("/modal", handleModal)
	http.ListenAndServe(":4390", nil)
}
