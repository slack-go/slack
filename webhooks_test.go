package slack

import (
	"testing"
	"net/http"
	"encoding/json"
	"reflect"
)

func TestPostWebhook_OK(t *testing.T) {
	once.Do(startServer)

	var receivedPayload WebhookMessage

	http.HandleFunc("/webhook", func(rw http.ResponseWriter, r *http.Request) {
		rw.Header().Set("Content-Type", "application/json")

		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&receivedPayload)
		if err != nil {
			t.Errorf("Request contained invalid JSON, %s", err)
		}

		response := []byte(`{}`)
		rw.Write(response)
	})

	url := "http://" + serverAddr + "/webhook"

	payload := &WebhookMessage{
		Text: "Test Text",
		Attachments: []Attachment{
			{
				Text: "Foo",
			},
		},
	}

	err := PostWebhook(url, payload)

	if err != nil {
		t.Errorf("Expected not to receive error: %s", err)
	}

	if !reflect.DeepEqual(payload, &receivedPayload) {
		t.Errorf("Payload did not match\nwant: %#v\n got: %#v", payload, receivedPayload)
	}
}

func TestPostWebhook_NotOK(t *testing.T) {
	once.Do(startServer)

	http.HandleFunc("/webhook2", func(rw http.ResponseWriter, r *http.Request) {
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte("500 - Something bad happened!"))
	})

	url := "http://" + serverAddr + "/webhook2"

	err := PostWebhook(url, &WebhookMessage{})

	if err == nil {
		t.Errorf("Expected to receive error")
	}
}
