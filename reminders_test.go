package slack

import (
	"bytes"
	"io"
	"net/http"
	"reflect"
	"testing"
)

type remindersHandler struct {
	gotParams map[string]string
}

func newRemindersHandler() *remindersHandler {
	return &remindersHandler{
		gotParams: make(map[string]string),
	}
}

func (rh *remindersHandler) accumulateFormValue(k string, r *http.Request) {
	if v := r.FormValue(k); v != "" {
		rh.gotParams[k] = v
	}
}

func (rh *remindersHandler) handler(w http.ResponseWriter, r *http.Request) {
	rh.accumulateFormValue("channel", r)
	rh.accumulateFormValue("user", r)
	rh.accumulateFormValue("text", r)
	rh.accumulateFormValue("time", r)
	rh.accumulateFormValue("reminder", r)
	w.Header().Set("Content-Type", "application/json")
	if rh.gotParams["text"] == "trigger-error" || rh.gotParams["reminder"] == "trigger-error" {
		w.Write([]byte(`{ "ok": false, "error": "oh no" }`))
	} else {
		w.Write([]byte(`{ "ok": true }`))
	}
}

func TestSlack_AddReminder(t *testing.T) {
	once.Do(startServer)
	api := New("testing-token", OptionAPIURL("http://"+serverAddr+"/"))
	tests := []struct {
		chanID     string
		userID     string
		text       string
		time       string
		wantParams map[string]string
		expectErr  bool
	}{
		{
			"someChannelID",
			"",
			"hello world",
			"tomorrow at 9am",
			map[string]string{
				"text":    "hello world",
				"time":    "tomorrow at 9am",
				"channel": "someChannelID",
			},
			false,
		},
		{
			"someChannelID",
			"",
			"trigger-error",
			"tomorrow at 9am",
			map[string]string{
				"text":    "trigger-error",
				"time":    "tomorrow at 9am",
				"channel": "someChannelID",
			},
			true,
		},
		{
			"",
			"someUserID",
			"hello world",
			"tomorrow at 9am",
			map[string]string{
				"text": "hello world",
				"time": "tomorrow at 9am",
				"user": "someUserID",
			},
			false,
		},
		{
			"",
			"someUserID",
			"trigger-error",
			"tomorrow at 9am",
			map[string]string{
				"text": "trigger-error",
				"time": "tomorrow at 9am",
				"user": "someUserID",
			},
			true,
		},
	}
	var rh *remindersHandler
	http.HandleFunc("/reminders.add", func(w http.ResponseWriter, r *http.Request) { rh.handler(w, r) })
	for i, test := range tests {
		rh = newRemindersHandler()
		var err error
		if test.chanID != "" {
			_, err = api.AddChannelReminder(test.chanID, test.text, test.time)
		} else {
			_, err = api.AddUserReminder(test.userID, test.text, test.time)
		}
		if test.expectErr == false && err != nil {
			t.Fatalf("%d: Unexpected error: %s", i, err)
		} else if test.expectErr == true && err == nil {
			t.Fatalf("%d: Expected error but got none!", i)
		}
		if !reflect.DeepEqual(rh.gotParams, test.wantParams) {
			t.Errorf("%d: Got params %#v, want %#v", i, rh.gotParams, test.wantParams)
		}
	}
}

func TestSlack_DeleteReminder(t *testing.T) {
	once.Do(startServer)
	api := New("testing-token", OptionAPIURL("http://"+serverAddr+"/"))
	tests := []struct {
		reminder   string
		wantParams map[string]string
		expectErr  bool
	}{
		{
			"foo",
			map[string]string{
				"reminder": "foo",
			},
			false,
		},
		{
			"trigger-error",
			map[string]string{
				"reminder": "trigger-error",
			},
			true,
		},
	}
	var rh *remindersHandler
	http.HandleFunc("/reminders.delete", func(w http.ResponseWriter, r *http.Request) { rh.handler(w, r) })
	for i, test := range tests {
		rh = newRemindersHandler()
		err := api.DeleteReminder(test.reminder)
		if test.expectErr == false && err != nil {
			t.Fatalf("%d: Unexpected error: %s", i, err)
		} else if test.expectErr == true && err == nil {
			t.Fatalf("%d: Expected error but got none!", i)
		}
		if !reflect.DeepEqual(rh.gotParams, test.wantParams) {
			t.Errorf("%d: Got params %#v, want %#v", i, rh.gotParams, test.wantParams)
		}
	}
}

type mockRemindersListHTTPClient struct{}

func (m *mockRemindersListHTTPClient) Do(*http.Request) (*http.Response, error) {
	responseString := `{
		"ok": true,
		"reminders": [
	        {
				"id": "Rm12345678",
				"creator": "U18888888",
				"user": "U18888888",
				"text": "eat a banana",
				"recurring": false,
				"time": 1458678068,
				"complete_ts": 0
			},
			{
				"id": "Gm12345678",
				"creator": "U18888888",
				"user": "U18888888",
				"text": "drink some water",
				"recurring": false,
				"time": 1458678090,
				"complete_ts": 0
			}
		]
	}`

	return &http.Response{StatusCode: 200, Body: io.NopCloser(bytes.NewBufferString(responseString))}, nil
}

func TestSlack_ListReminders(t *testing.T) {
	expectedIDs := []string{"Rm12345678", "Gm12345678"}

	once.Do(startServer)
	api := &Client{
		endpoint:   "http://" + serverAddr + "/",
		token:      "testing-token",
		httpclient: &mockRemindersListHTTPClient{},
	}

	reminders, err := api.ListReminders()
	if err != nil {
		t.Fatalf("Unexpected error: %s", err)
	}

	for i := range reminders {
		if reminders[i].ID != expectedIDs[i] {
			t.Fatalf("List Reminders data wasn't correctly populated: wanted %v, got %v", expectedIDs[i], reminders[i].ID)
		}
	}
}
