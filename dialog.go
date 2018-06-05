package slack

import (
	"context"
	"encoding/json"
	"net/url"
)

// Dialog as in Slack dialogs
// 	https://api.slack.com/dialogs#option_element_attributes#top-level_dialog_attributes
type Dialog struct {
	TriggerID      string          `json:"trigger_id,omitempty"`
	CallbackID     string          `json:"callback_id"`
	NotifyOnCancel bool            `json:"notify_on_cancel"`
	Title          string          `json:"title"`
	SubmitLabel    string          `json:"submit_label,omitempty"`
	Elements       []DialogElement `json:"elements"`
}

// DialogElement Abstract interface for Elements
type DialogElement interface{}

// DialogInput for dialogs input type text or menu
type DialogInput struct {
	Type        InputType `json:"type"`
	Label       string    `json:"label"`
	Name        string    `json:"name"`
	Placeholder string    `json:"placeholder"`
	Optional    bool      `json:"optional"`
}

// InputType is the type of the dialog input type
type InputType string

const (
	// InputTypeText textfield input
	InputTypeText InputType = "text"
	// InputTypeTextArea textarea input
	InputTypeTextArea InputType = "textarea"
	// InputTypeSelect textfield input
	InputTypeSelect InputType = "select"
)

// DialogSubmitCallback to parse the response back from the Dialog
type DialogSubmitCallback struct {
	Type       string            `json:"type"`
	Submission map[string]string `json:"submission"`
	CallbackID string            `json:"callback_id"`

	Team        Team    `json:"team"`
	Channel     Channel `json:"channel"`
	User        User    `json:"user"`
	ActionTs    string  `json:"action_ts"`
	Token       string  `json:"token"`
	ResponseURL string  `json:"response_url"`
}

// DialogOpenResponse response from `dialog.open`
type DialogOpenResponse struct {
	Ok                     bool             `json:"ok"`
	Error                  string           `json:"error"`
	DialogResponseMetadata ResponseMetadata `json:"response_metadata"`
}

// DialogResponseMetadata lists the error messages
type DialogResponseMetadata struct {
	Messages []string `json:"messages"`
}

// OpenDialog posts the `dialog` to slack's `Dialog.open` endpoint
func (api *Client) OpenDialog(dialog Dialog) (err error) {
	return api.OpenDialogContext(context.Background(), dialog)
}

// OpenDialogContext opens a dialog window where the triggerId originated from with a custom context
func (api *Client) OpenDialogContext(ctx context.Context, dialog Dialog) (err error) {
	dialogjson, err := json.Marshal(dialog)
	if err != nil {
		return
	}

	values := url.Values{
		"token":      {api.token},
		"trigger_id": {dialog.TriggerID},
		"dialog":     {string(dialogjson)},
	}

	response := &DialogOpenResponse{}
	return postForm(ctx, api.httpclient, "dialog.open", values, response, api.debug)
}
