package slack

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
)

type ChallengeHandlerOptions struct {
	// Set to 'true' if you want to process events synchronously. This means the supplied
	// http.Handler is responsible for sending a 200 back to Slack within 3 seconds of recieving
	// the request. By default, your http.Handler will be called in a new goroutine with a mocked
	// ResponseWriter and Request object.
	SyncProcessing bool
}

type challengeEvent struct {
	Token     string `json:"token"`
	Challenge []byte `json:"challenge"`
	Type      string `json:"type"`
}

const challengeEventType = "url_verification"

// fakeResponseWriter is mocking the http.ResponseWriter interface
// so we can send slack a 200 as soon as possible, and pass off
// this mock to the user-supplied HTTP handler. It's got a lower memory
// overhead than the httptest Recorder.
//
// We can make this optional for advanced users who require handling
// the response in all cases.
//
// Slack suggests this approach as a best practice:
// https://api.slack.com/events-api#tips
type fakeResponseWriter struct {
	headers http.Header
}

// Wraps an HTTP Handler to respond to Slack's Events API challenge. If the handler gets
// a JSON request that isn't the challenge event, it will forward the request to the underlying
// handler for processing.
func WithChallengeHandler(h http.Handler, opts ChallengeHandlerOptions) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		r.Body = ioutil.NopCloser(bytes.NewBuffer(body))

		// If we have a valid payload from Slack, attempt to parse out the challenge
		// event and respond to it.
		if len(body) > 0 {
			evt := challengeEvent{}
			err = json.Unmarshal(body, &evt)
			// If the body is invalid JSON, we should return a BadRequest,
			// but we shouldn't complain about bad types, etc. since this may
			// be a conflicting payload from Slack
			if _, ok := err.(*json.SyntaxError); ok {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			if err == nil && strings.TrimSpace(evt.Type) == challengeEventType {
				w.Header().Set("Content-Type", "text/plain")
				// XXX: We can likely ignore this error, if we fail to respond
				//      slack _should_ retry to send the request.
				_, _ = w.Write(evt.Challenge)
				return
			}
		}

		// If we didn't intercept the challenge event, pass it off to the underlying
		// HTTP handler. The sync handler should be called inline with the real ResponseWriter
		// while the async handler should write a 200 OK and pass the event off to the handler
		if opts.SyncProcessing {
			h.ServeHTTP(w, r)
		} else {
			w.WriteHeader(http.StatusOK)
			go h.ServeHTTP(newResponseWriter(), r)
		}
	})
}

func newResponseWriter() *fakeResponseWriter {
	return &fakeResponseWriter{
		headers: make(http.Header),
	}
}

func (r *fakeResponseWriter) Header() http.Header {
	return r.headers
}

func (r *fakeResponseWriter) Write(body []byte) (int, error) {
	return len(body), nil
}

func (r *fakeResponseWriter) WriteHeader(status int) {}
