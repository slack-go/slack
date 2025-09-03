package slack

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
)

// InteractionType type of interactions
type InteractionType string

// ActionType type represents the type of action (attachment, block, etc.)
type ActionType string

// action is an interface that should be implemented by all callback action types
type action interface {
	actionType() ActionType
}

// Types of interactions that can be received.
const (
	InteractionTypeDialogCancellation = InteractionType("dialog_cancellation")
	InteractionTypeDialogSubmission   = InteractionType("dialog_submission")
	InteractionTypeDialogSuggestion   = InteractionType("dialog_suggestion")
	InteractionTypeInteractionMessage = InteractionType("interactive_message")
	InteractionTypeMessageAction      = InteractionType("message_action")
	InteractionTypeBlockActions       = InteractionType("block_actions")
	InteractionTypeBlockSuggestion    = InteractionType("block_suggestion")
	InteractionTypeViewSubmission     = InteractionType("view_submission")
	InteractionTypeViewClosed         = InteractionType("view_closed")
	InteractionTypeShortcut           = InteractionType("shortcut")
	InteractionTypeWorkflowStepEdit   = InteractionType("workflow_step_edit")
)

// InteractionCallback is sent from slack when a user interactions with a button or dialog.
type InteractionCallback struct {
	Type                InteractionType `json:"type" form:"type"` // continue adding form tag
	Token               string          `json:"token" form:"token"`
	CallbackID          string          `json:"callback_id" form:"callback_id"`
	ResponseURL         string          `json:"response_url" form:"response_url"`
	TriggerID           string          `json:"trigger_id" form:"trigger_id"`
	ActionTs            string          `json:"action_ts" form:"action_ts"`
	Team                Team            `json:"team" form:"team"`
	Channel             Channel         `json:"channel" form:"channel"`
	User                User            `json:"user" form:"user"`
	OriginalMessage     Message         `json:"original_message" form:"original_message"`
	Message             Message         `json:"message" form:"message"`
	Name                string          `json:"name" form:"name"`
	Value               string          `json:"value" form:"value"`
	MessageTs           string          `json:"message_ts" form:"message_ts"`
	AttachmentID        string          `json:"attachment_id" form:"attachment_id"`
	ActionCallback      ActionCallbacks `json:"actions" form:"actions"`
	View                View            `json:"view" form:"view"`
	ActionID            string          `json:"action_id" form:"action_id"`
	APIAppID            string          `json:"api_app_id" form:"api_app_id"`
	BlockID             string          `json:"block_id" form:"block_id"`
	Container           Container       `json:"container" form:"container"`
	Enterprise          Enterprise      `json:"enterprise" form:"enterprise"`
	IsEnterpriseInstall bool            `json:"is_enterprise_install" form:"is_enterprise_install"`
	DialogSubmissionCallback
	ViewSubmissionCallback
	ViewClosedCallback

	// FIXME(kanata2): just workaround for backward-compatibility.
	// See also https://github.com/slack-go/slack/issues/816
	RawState json.RawMessage `json:"state,omitempty"`

	// BlockActionState stands for the `state` field in block_actions type.
	// NOTE: InteractionCallback.State has a role for the state of dialog_submission type,
	// so we cannot use this field for backward-compatibility for now.
	BlockActionState *BlockActionStates `json:"-"`
}

type BlockActionStates struct {
	Values map[string]map[string]BlockAction `json:"values" form:"values"`
}

// InteractionCallbackParse parses the HTTP form value "payload" from r, unmarshals
// it as JSON into an InteractionCallback, and returns the result.
// It returns an error if the payload is missing or cannot be decoded.
//
// See https://github.com/slack-go/slack/issues/660 for context.
func InteractionCallbackParse(r *http.Request) (InteractionCallback, error) {
	payload := r.FormValue("payload")
	if len(payload) == 0 {
		return InteractionCallback{}, errors.New("payload is empty")
	}

	var ic InteractionCallback
	if err := json.Unmarshal([]byte(payload), &ic); err != nil {
		return InteractionCallback{}, err
	}
	return ic, nil
}

func (ic *InteractionCallback) MarshalJSON() ([]byte, error) {
	type alias InteractionCallback
	tmp := alias(*ic)
	if tmp.Type == InteractionTypeBlockActions {
		if tmp.BlockActionState == nil {
			tmp.RawState = []byte(`{}`)
		} else {
			state, err := json.Marshal(tmp.BlockActionState.Values)
			if err != nil {
				return nil, err
			}
			tmp.RawState = []byte(`{"values":` + string(state) + `}`)
		}
	} else if ic.Type == InteractionTypeDialogSubmission {
		tmp.RawState = []byte(tmp.State)
	}
	// Use pointer for go1.7
	return json.Marshal(&tmp)
}

func (ic *InteractionCallback) UnmarshalJSON(b []byte) error {
	type alias InteractionCallback
	tmp := struct {
		Type InteractionType `json:"type"`
		*alias
	}{
		alias: (*alias)(ic),
	}
	if err := json.Unmarshal(b, &tmp); err != nil {
		return err
	}
	*ic = InteractionCallback(*tmp.alias)
	ic.Type = tmp.Type
	if ic.Type == InteractionTypeBlockActions {
		if len(ic.RawState) > 0 {
			err := json.Unmarshal(ic.RawState, &ic.BlockActionState)
			if err != nil {
				return err
			}
		}
	} else if ic.Type == InteractionTypeDialogSubmission {
		ic.State = string(ic.RawState)
	}
	return nil
}

type Container struct {
	Type         string      `json:"type" form:"type"`
	ViewID       string      `json:"view_id" form:"view_id"`
	MessageTs    string      `json:"message_ts" form:"message_ts"`
	ThreadTs     string      `json:"thread_ts,omitempty" form:"thread_ts"`
	AttachmentID json.Number `json:"attachment_id" form:"attachment_id"`
	ChannelID    string      `json:"channel_id" form:"channel_id"`
	IsEphemeral  bool        `json:"is_ephemeral" form:"is_ephemeral"`
	IsAppUnfurl  bool        `json:"is_app_unfurl" form:"is_app_unfurl"`
}

type Enterprise struct {
	ID   string `json:"id" form:"id"`
	Name string `json:"name" form:"name"`
}

// ActionCallback is a convenience struct defined to allow dynamic unmarshalling of
// the "actions" value in Slack's JSON response, which varies depending on block type
type ActionCallbacks struct {
	AttachmentActions []*AttachmentAction
	BlockActions      []*BlockAction
}

// MarshalJSON implements the Marshaller interface in order to combine both
// action callback types back into a single array, like how the api responds.
// This makes Marshaling and Unmarshaling an InteractionCallback symmetrical
func (a ActionCallbacks) MarshalJSON() ([]byte, error) {
	count := 0
	length := len(a.AttachmentActions) + len(a.BlockActions)
	buffer := bytes.NewBufferString("[")

	f := func(obj interface{}) error {
		js, err := json.Marshal(obj)
		if err != nil {
			return err
		}
		_, err = buffer.Write(js)
		if err != nil {
			return err
		}

		count++
		if count < length {
			_, err = buffer.WriteString(",")
			return err
		}
		return nil
	}

	for _, act := range a.AttachmentActions {
		err := f(act)
		if err != nil {
			return nil, err
		}
	}
	for _, blk := range a.BlockActions {
		err := f(blk)
		if err != nil {
			return nil, err
		}
	}
	buffer.WriteString("]")
	return buffer.Bytes(), nil
}

// UnmarshalJSON implements the Marshaller interface in order to delegate
// marshalling and allow for proper type assertion when decoding the response
func (a *ActionCallbacks) UnmarshalJSON(data []byte) error {
	var raw []json.RawMessage
	err := json.Unmarshal(data, &raw)
	if err != nil {
		return err
	}

	for _, r := range raw {
		var obj map[string]interface{}
		err := json.Unmarshal(r, &obj)
		if err != nil {
			return err
		}

		if _, ok := obj["block_id"].(string); ok {
			action, err := unmarshalAction(r, &BlockAction{})
			if err != nil {
				return err
			}

			a.BlockActions = append(a.BlockActions, action.(*BlockAction))
			continue
		}

		action, err := unmarshalAction(r, &AttachmentAction{})
		if err != nil {
			return err
		}
		a.AttachmentActions = append(a.AttachmentActions, action.(*AttachmentAction))
	}

	return nil
}

func unmarshalAction(r json.RawMessage, callbackAction action) (action, error) {
	err := json.Unmarshal(r, callbackAction)
	if err != nil {
		return nil, err
	}
	return callbackAction, nil
}
