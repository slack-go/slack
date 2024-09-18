package main

import (
	"fmt"

	"github.com/slack-go/slack"
)

// An example how to open a modal with different kinds of input fields
func main() {

	// Create a ModalViewRequest with a header and two inputs
	titleText := slack.NewTextBlockObject(slack.PlainTextType, "Create channel demo", false, false)
	closeText := slack.NewTextBlockObject(slack.PlainTextType, "Close", false, false)
	submitText := slack.NewTextBlockObject(slack.PlainTextType, "Submit", false, false)

	contextText := slack.NewTextBlockObject(slack.MarkdownType, "This app demonstrates the use of different fields", false, false)
	contextBlock := slack.NewContextBlock("context", contextText)

	// Only the inputs in input blocks will be included in view_submission’s view.state.values: https://slack.dev/java-slack-sdk/guides/modals
	// This means the inputs will not be interactive either because they do not trigger block_actions messages: https://api.slack.com/surfaces/modals/using#interactions
	channelNameText := slack.NewTextBlockObject(slack.PlainTextType, "Channel Name", false, false)
	channelNameHint := slack.NewTextBlockObject(slack.PlainTextType, "Channel names may only contain lowercase letters, numbers, hyphens, and underscores, and must be 80 characters or less", false, false)
	channelPlaceholder := slack.NewTextBlockObject(slack.PlainTextType, "New channel name", false, false)
	channelNameElement := slack.NewPlainTextInputBlockElement(channelPlaceholder, "channel_name")
	// Slack channel names can be maximum 80 characters: https://api.slack.com/methods/conversations.create
	channelNameElement.MaxLength = 80
	channelNameBlock := slack.NewInputBlock("channel_name", channelNameText, channelNameHint, channelNameElement)

	// Provide a static list of users to choose from, those provided now are just made up user IDs
	// Get user IDs by right clicking on them in Slack, select "Copy link", and inspect the last part of the link
	// The user ID should start with "U" followed by 8 random characters
	memberOptions := createOptionBlockObjects([]string{"U9911MMAA", "U2233KKNN", "U00112233"}, true)
	inviteeText := slack.NewTextBlockObject(slack.PlainTextType, "Invitee from static list", false, false)
	inviteeOption := slack.NewOptionsSelectBlockElement(slack.OptTypeStatic, nil, "invitee", memberOptions...)
	inviteeBlock := slack.NewInputBlock("invitee", inviteeText, nil, inviteeOption)

	// Section with users select - this input will not be included in the view_submission's view.state.values,
	// but instead be sent as a "block_actions" request
	additionalInviteeText := slack.NewTextBlockObject(slack.PlainTextType, "Invitee from complete list of users", false, false)
	additionalInviteeHintText := slack.NewTextBlockObject(slack.PlainTextType, "", false, false)
	additionalInviteeOption := slack.NewOptionsSelectBlockElement(slack.OptTypeUser, additionalInviteeText, "user")
	additionalInviteeSection := slack.NewSectionBlock(additionalInviteeText, nil, slack.NewAccessory(additionalInviteeOption))

	// Input with users select - this input will be included in the view_submission's view.state.values
	// It can be fetched as for example "payload.View.State.Values["user"]["user"].SelectedUser"
	additionalInviteeBlock := slack.NewInputBlock("user", additionalInviteeText, additionalInviteeHintText, additionalInviteeOption)

	checkboxTxt := slack.NewTextBlockObject(slack.PlainTextType, "Checkbox", false, false)
	checkboxOptions := createOptionBlockObjects([]string{"option 1", "option 2", "option 3"}, false)
	checkboxOptionsBlock := slack.NewCheckboxGroupsBlockElement("chkbox", checkboxOptions...)
	checkboxBlock := slack.NewInputBlock("chkbox", checkboxTxt, nil, checkboxOptionsBlock)

	summaryText := slack.NewTextBlockObject(slack.PlainTextType, "Summary", false, false)
	summaryHint := slack.NewTextBlockObject(slack.PlainTextType, "Summary Hint", false, false)
	summaryPlaceholder := slack.NewTextBlockObject(slack.PlainTextType, "Summary of reason for creating channel", false, false)
	summaryElement := slack.NewPlainTextInputBlockElement(summaryPlaceholder, "summary")
	// Just set an arbitrary max length to avoid too prose summary
	summaryElement.MaxLength = 200
	summaryElement.Multiline = true
	summaryBlock := slack.NewInputBlock("summary", summaryText, summaryHint, summaryElement)

	blocks := slack.Blocks{
		BlockSet: []slack.Block{
			contextBlock,
			channelNameBlock,
			inviteeBlock,
			additionalInviteeSection,
			additionalInviteeBlock,
			checkboxBlock,
			summaryBlock,
		},
	}

	var modalRequest slack.ModalViewRequest
	modalRequest.Type = slack.ViewType("modal")
	modalRequest.Title = titleText
	modalRequest.Close = closeText
	modalRequest.Submit = submitText
	modalRequest.Blocks = blocks
	modalRequest.CallbackID = "create_channel"

	api := slack.New("YOUR_BOT_TOKEN_HERE")

	// Using a trigger ID you can open a modal
	// The trigger ID is provided through certain events and interactions
	// More information can be found here: https://api.slack.com/interactivity/handling#modal_responses
	_, err := api.OpenView("YOUR_TRIGGERID_HERE", modalRequest)
	if err != nil {
		fmt.Printf("Error opening view: %s", err)
	}
}

// createOptionBlockObjects - utility function for generating option block objects
func createOptionBlockObjects(options []string, users bool) []*slack.OptionBlockObject {
	optionBlockObjects := make([]*slack.OptionBlockObject, 0, len(options))
	var text string
	for _, o := range options {
		if users {
			text = fmt.Sprintf("<@%s>", o)
		} else {
			text = o
		}
		optionText := slack.NewTextBlockObject(slack.PlainTextType, text, false, false)
		optionBlockObjects = append(optionBlockObjects, slack.NewOptionBlockObject(o, optionText, nil))
	}
	return optionBlockObjects
}
