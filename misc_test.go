package slack

import (
	"context"
	"log"
	"net/http"
	"net/url"
	"sync"
	"testing"

	"github.com/slack-go/slack/slackutilsx"
)

var (
	parseResponseOnce sync.Once
)

func parseResponseHandler(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	token := r.FormValue("token")
	log.Println(token)
	if token == "" {
		rw.Write([]byte(`{"ok":false,"error":"not_authed"}`))
		return
	}
	if token != validToken {
		rw.Write([]byte(`{"ok":false,"error":"invalid_auth"}`))
		return
	}
	response := []byte(`{"ok": true}`)
	rw.Write(response)
}

func setParseResponseHandler() {
	http.HandleFunc("/parseResponse", parseResponseHandler)
}

func TestParseResponse(t *testing.T) {
	parseResponseOnce.Do(setParseResponseHandler)
	once.Do(startServer)
	APIURL := "http://" + serverAddr + "/"
	values := url.Values{
		"token": {validToken},
	}

	responsePartial := &SlackResponse{}
	err := postForm(context.Background(), http.DefaultClient, APIURL+"parseResponse", values, responsePartial, discard{})
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}
}

func TestParseResponseNoToken(t *testing.T) {
	parseResponseOnce.Do(setParseResponseHandler)
	once.Do(startServer)
	APIURL := "http://" + serverAddr + "/"
	values := url.Values{}

	responsePartial := &SlackResponse{}
	err := postForm(context.Background(), http.DefaultClient, APIURL+"parseResponse", values, responsePartial, discard{})
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
		return
	}
	if responsePartial.Ok {
		t.Errorf("Unexpected error: %s", err)
	} else if responsePartial.Error != "not_authed" {
		t.Errorf("got %v; want %v", responsePartial.Error, "not_authed")
	}
}

func TestParseResponseInvalidToken(t *testing.T) {
	parseResponseOnce.Do(setParseResponseHandler)
	once.Do(startServer)
	APIURL := "http://" + serverAddr + "/"
	values := url.Values{
		"token": {"whatever"},
	}
	responsePartial := &SlackResponse{}
	err := postForm(context.Background(), http.DefaultClient, APIURL+"parseResponse", values, responsePartial, discard{})
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
		return
	}
	if responsePartial.Ok {
		t.Errorf("Unexpected error: %s", err)
	} else if responsePartial.Error != "invalid_auth" {
		t.Errorf("got %v; want %v", responsePartial.Error, "invalid_auth")
	}
}

func TestRetryable(t *testing.T) {
	for _, e := range []error{
		&RateLimitedError{},
		StatusCodeError{Code: http.StatusInternalServerError},
		StatusCodeError{Code: http.StatusTooManyRequests},
	} {
		r, ok := e.(slackutilsx.Retryable)
		if !ok {
			t.Errorf("expected %#v to implement Retryable", e)
		}
		if !r.Retryable() {
			t.Errorf("expected %#v to be Retryable", e)
		}
	}
}
