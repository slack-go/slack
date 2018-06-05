package slack

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
)

const (
	triggerID      = "trigger_xyz"
	callbackID     = "callback_xyz"
	notifyOnCancel = false
	title          = "Dialog_title"
	submitLabel    = "Send"
	token          = "xoxa-123-123-123-213"
)

func _mocDialog() *Dialog {
	triggerID := triggerID
	callbackID := callbackID
	notifyOnCancel := notifyOnCancel
	title := title
	submitLabel := submitLabel

	return &Dialog{
		TriggerID:      triggerID,
		CallbackID:     callbackID,
		NotifyOnCancel: notifyOnCancel,
		Title:          title,
		SubmitLabel:    submitLabel,
	}
}

func TestDialogCreate(t *testing.T) {
	dialog := _mocDialog()
	if dialog == nil {
		t.Errorf("Should be able to construct a dialog")
		t.Fail()
	}
}

func ExampleDialog() {
	dialog := _mocDialog()
	fmt.Println(*dialog)
	// Output:
	// {trigger_xyz callback_xyz false Dialog_title Send []}
}

// This tests GET request with passing in a parameter.
func TestDialogOpen(t *testing.T) {
	api := New(token)
	api.httpclient = &mocHTTPClient{
		content: SlackResponse{
			Ok: false,
		},
	}
	dialog := _mocDialog()
	err := api.OpenDialog(*dialog)
	if err != nil {
		t.Errorf("Failed to open Dialog %v", dialog)
	}
}

type mocHTTPClient struct {
	content interface{}
}

func (moc *mocHTTPClient) Do(req *http.Request) (*http.Response, error) {
	req.ParseForm()
	fmt.Println(req.URL, req.PostForm)
	contents, _ := json.Marshal(moc.content)
	bodyReader := ioutil.NopCloser(bytes.NewReader(contents))
	response := &http.Response{
		StatusCode: http.StatusOK,
		Body:       bodyReader,
	}
	return response, nil
}
