package slack

import (
	"encoding/json"
	"net/http"
	"reflect"
	"testing"
	"time"
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
	if scerr, ok := err.(StatusCodeError); !ok {
		t.Errorf("Expected error of type StatusCodeError, got %#v", err)
	} else if scerr.Code != http.StatusInternalServerError {
		t.Errorf("Expected %d, got %d", http.StatusInternalServerError, scerr.Code)
	}
}

func TestPostWebhook_RateLimited(t *testing.T) {
	once.Do(startServer)

	http.HandleFunc("/webhook3", func(rw http.ResponseWriter, r *http.Request) {
		rw.Header().Add("Retry-After", "60")
		rw.WriteHeader(http.StatusTooManyRequests)
	})

	url := "http://" + serverAddr + "/webhook3"

	err := PostWebhook(url, &WebhookMessage{})

	if err == nil {
		t.Errorf("Expected to receive error")
	}
	if rlerr, ok := err.(*RateLimitedError); !ok {
		t.Errorf("Expected error of type RateLimitedError, got %#v", err)
	} else if rlerr.RetryAfter != 60*time.Second {
		t.Errorf("Expected retry after %s, got %s", 60*time.Second, rlerr.RetryAfter)
	}
}
