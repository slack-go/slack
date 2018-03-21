package slack

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

var challengeResponse = []byte("challenge")
var challengeBody, _ = json.Marshal(challengeEvent{
	Token:     "1234",
	Challenge: challengeResponse,
	Type:      challengeEventType,
})

var eventBody = []byte("{ \"foo\": \"bar\"}")

func TestChallengeHandlerResponses(t *testing.T) {
	var teapotBody = []byte("I'm a teapot")
	var emptyBody = []byte("")
	teapotHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusTeapot)
		w.Write(teapotBody)
	})

	cases := []struct {
		Options      ChallengeHandlerOptions
		RequestBody  []byte
		ExpectedCode int
		ExpectedBody []byte
	}{
		// Tests that the sync handler will rely on the underlying handler
		// for event processing / responses
		{
			Options:      ChallengeHandlerOptions{SyncProcessing: true},
			RequestBody:  eventBody,
			ExpectedCode: http.StatusTeapot,
			ExpectedBody: teapotBody,
		},
		// Tests that the sync handler will rely on the challenge handler
		// for challenge responses
		{
			Options:      ChallengeHandlerOptions{SyncProcessing: true},
			RequestBody:  challengeBody,
			ExpectedCode: http.StatusOK,
			ExpectedBody: challengeResponse,
		},
		// Tests that the sync handler will return 400 on invalid json
		{
			Options:      ChallengeHandlerOptions{SyncProcessing: true},
			RequestBody:  []byte("c"),
			ExpectedCode: http.StatusBadRequest,
			ExpectedBody: []byte("invalid character 'c' looking for beginning of value\n"),
		},
		// Tests that the async handler will rely on the challenge handler
		// for responses
		{
			Options:      ChallengeHandlerOptions{},
			RequestBody:  eventBody,
			ExpectedCode: http.StatusOK,
			ExpectedBody: emptyBody,
		},
		// Tests that the async handler will rely on the challenge handler
		// for challenge responses
		{
			Options:      ChallengeHandlerOptions{},
			RequestBody:  challengeBody,
			ExpectedCode: http.StatusOK,
			ExpectedBody: challengeResponse,
		},
		// Tests that the async handler will return 400 on invalid json
		{
			Options:      ChallengeHandlerOptions{},
			RequestBody:  []byte("c"),
			ExpectedCode: http.StatusBadRequest,
			ExpectedBody: []byte("invalid character 'c' looking for beginning of value\n"),
		},
	}

	for i, test := range cases {
		wrappedHandler := WithChallengeHandler(teapotHandler, test.Options)

		resp := httptest.NewRecorder()
		req, err := http.NewRequest("POST", "/webhook", bytes.NewBuffer(test.RequestBody))
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Content-Type", "application/json")

		wrappedHandler.ServeHTTP(resp, req)
		if resp.Code != test.ExpectedCode {
			t.Errorf("[Case %d]: Server returned status code %d, expected %d", i, resp.Code, test.ExpectedCode)
		}

		if resp.Body.String() != string(test.ExpectedBody) {
			t.Errorf("[Case %d] Server returned body %q, expected %q", i, resp.Body, string(test.ExpectedBody))
		}
	}
}

// The default challenge handler runs each http handler in its own Goroutine
// and responds 200 back to Slack for you.
func ExampleWithChallengeHandler(t *testing.T) {
	type SlackEvent struct {
		Type           string          `json:"type"`
		EventTimestamp string          `json:"event_ts"`
		Timestamp      string          `json:"ts"`
		Event          json.RawMessage `json:"event"`
	}

	webhookHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			return
		}
		msg := SlackEvent{}
		json.Unmarshal(body, &msg)
		fmt.Printf("Got Event: %s", msg.Type)
	})

	http.Handle("/webhook", WithChallengeHandler(webhookHandler, ChallengeHandlerOptions{}))
	http.ListenAndServe(":8000", nil)
}

// If you want to respond to the Events API yourself, you can configure
// the webhook handler to use synchronous processing. Slack has a 3 second timeout in
// which you must respond, therefor you should respond to the API prior to processing
// the event.
func ExampleWithChallengeHandler_withSyncProcessing(t *testing.T) {
	type SlackEvent struct {
		Type           string          `json:"type"`
		EventTimestamp string          `json:"event_ts"`
		Timestamp      string          `json:"ts"`
		Event          json.RawMessage `json:"event"`
	}

	webhookHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		go func() {
			msg := SlackEvent{}
			json.Unmarshal(body, &msg)
			fmt.Printf("Got Event: %s", msg.Type)
		}()
		w.WriteHeader(http.StatusOK)
	})

	http.Handle("/webhook", WithChallengeHandler(webhookHandler, ChallengeHandlerOptions{
		SyncProcessing: true,
	}))
	http.ListenAndServe(":8000", nil)
}
