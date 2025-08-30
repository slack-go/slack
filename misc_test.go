package slack

import (
	"context"
	"encoding/json"
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

func TestSlackResponseErrorsMarshaling(t *testing.T) {
	tests := []struct {
		name     string
		errors   SlackResponseErrors
		expected string
	}{
		{
			name: "AppsManifestCreateResponseError",
			errors: SlackResponseErrors{
				AppsManifestCreateResponseError: &AppsManifestCreateResponseError{
					Message: "Interactivity requires Socket Mode enabled",
					Pointer: "/settings/interactivity",
				},
			},
			expected: `{"message":"Interactivity requires Socket Mode enabled","pointer":"/settings/interactivity"}`,
		},
		{
			name: "ConversationsInviteResponseError",
			errors: SlackResponseErrors{
				ConversationsInviteResponseError: &ConversationsInviteResponseError{
					Error: "invalid_user",
					Ok:    false,
					User:  "U12345678",
				},
			},
			expected: `{"error":"invalid_user","ok":false,"user":"U12345678"}`,
		},
		{
			name: "StringError",
			errors: SlackResponseErrors{
				Message: func() *string { s := "failed to match all allowed schemas"; return &s }(),
			},
			expected: `"failed to match all allowed schemas"`,
		},
		{
			name:     "EmptyErrors",
			errors:   SlackResponseErrors{},
			expected: `null`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := json.Marshal(tt.errors)
			if err != nil {
				t.Fatalf("Marshal failed: %v", err)
			}
			if string(data) != tt.expected {
				t.Errorf("got %s; want %s", string(data), tt.expected)
			}
		})
	}
}

func TestSlackResponseErrorsUnmarshaling(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected SlackResponseErrors
	}{
		{
			name:  "AppsManifestCreateResponseError",
			input: `{"pointer":"/settings/interactivity","message":"Interactivity requires Socket Mode enabled"}`,
			expected: SlackResponseErrors{
				AppsManifestCreateResponseError: &AppsManifestCreateResponseError{
					Pointer: "/settings/interactivity",
					Message: "Interactivity requires Socket Mode enabled",
				},
			},
		},
		{
			name:  "ConversationsInviteResponseError",
			input: `{"error":"invalid_user","ok":false,"user":"U12345678"}`,
			expected: SlackResponseErrors{
				ConversationsInviteResponseError: &ConversationsInviteResponseError{
					Error: "invalid_user",
					Ok:    false,
					User:  "U12345678",
				},
			},
		},
		{
			name:     "NullInput",
			input:    `null`,
			expected: SlackResponseErrors{},
		},
		{
			name:  "StringError",
			input: `"failed to match all allowed schemas [json-pointer:\\/blocks\\/3\\/text]"`,
			expected: SlackResponseErrors{
				Message: func() *string {
					s := "failed to match all allowed schemas [json-pointer:\\/blocks\\/3\\/text]"
					return &s
				}(),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var errors SlackResponseErrors
			err := json.Unmarshal([]byte(tt.input), &errors)
			if err != nil {
				t.Fatalf("Unmarshal failed: %v", err)
			}

			if tt.expected.AppsManifestCreateResponseError != nil {
				if errors.AppsManifestCreateResponseError == nil {
					t.Error("expected AppsManifestCreateResponseError, got nil")
				} else if *errors.AppsManifestCreateResponseError != *tt.expected.AppsManifestCreateResponseError {
					t.Errorf("got %+v; want %+v", *errors.AppsManifestCreateResponseError, *tt.expected.AppsManifestCreateResponseError)
				}
			}

			if tt.expected.ConversationsInviteResponseError != nil {
				if errors.ConversationsInviteResponseError == nil {
					t.Error("expected ConversationsInviteResponseError, got nil")
				} else if *errors.ConversationsInviteResponseError != *tt.expected.ConversationsInviteResponseError {
					t.Errorf("got %+v; want %+v", *errors.ConversationsInviteResponseError, *tt.expected.ConversationsInviteResponseError)
				}
			}

			if tt.expected.Message != nil {
				if errors.Message == nil {
					t.Error("expected Message, got nil")
				} else if *errors.Message != *tt.expected.Message {
					t.Errorf("got %+v; want %+v", *errors.Message, *tt.expected.Message)
				}
			}
		})
	}
}

func TestSlackResponseErrorsUnmarshalingUnknownStructure(t *testing.T) {
	input := `{"unknown_field":"value","other_field":123}`
	var errors SlackResponseErrors
	err := json.Unmarshal([]byte(input), &errors)
	if err == nil {
		t.Error("expected error for unknown structure, got nil")
	}
	expectedError := "unknown error structure: " + input
	if err.Error() != expectedError {
		t.Errorf("got error %q; want %q", err.Error(), expectedError)
	}
}

func TestSlackResponseWithErrors(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected SlackResponse
	}{
		{
			name:  "ResponseWithAppsManifestCreateResponseError",
			input: `{"ok":false,"error":"invalid_manifest","errors":[{"pointer":"/settings/interactivity","message":"Interactivity requires Socket Mode enabled"}]}`,
			expected: SlackResponse{
				Ok:    false,
				Error: "invalid_manifest",
				Errors: []SlackResponseErrors{
					{
						AppsManifestCreateResponseError: &AppsManifestCreateResponseError{
							Pointer: "/settings/interactivity",
							Message: "Interactivity requires Socket Mode enabled",
						},
					},
				},
			},
		},
		{
			name:  "ResponseWithoutErrors",
			input: `{"ok":true}`,
			expected: SlackResponse{
				Ok: true,
			},
		},
		{
			name:  "ResponseWithStringErrors",
			input: `{"ok":false,"error":"invalid_blocks","errors":["failed to match all allowed schemas [json-pointer:\\/blocks\\/3\\/text]","invalid additional property: emoji [json-pointer:\\/blocks\\/3\\/text]"]}`,
			expected: SlackResponse{
				Ok:    false,
				Error: "invalid_blocks",
				Errors: []SlackResponseErrors{
					{
						Message: func() *string {
							s := "failed to match all allowed schemas [json-pointer:\\/blocks\\/3\\/text]"
							return &s
						}(),
					},
					{
						Message: func() *string {
							s := "invalid additional property: emoji [json-pointer:\\/blocks\\/3\\/text]"
							return &s
						}(),
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var response SlackResponse
			err := json.Unmarshal([]byte(tt.input), &response)
			if err != nil {
				t.Fatalf("Unmarshal failed: %v", err)
			}

			if response.Ok != tt.expected.Ok {
				t.Errorf("got Ok=%v; want Ok=%v", response.Ok, tt.expected.Ok)
			}
			if response.Error != tt.expected.Error {
				t.Errorf("got Error=%q; want Error=%q", response.Error, tt.expected.Error)
			}

			if tt.expected.Errors == nil {
				if response.Errors != nil {
					t.Error("expected nil Errors, got non-nil")
				}
			} else {
				if response.Errors == nil {
					t.Error("expected non-nil Errors, got nil")
					return
				}
				if len(tt.expected.Errors) == 0 {
					t.Errorf("got Errors=%v; want Errors=%v", response.Errors, tt.expected.Errors)
				}
			}
		})
	}
}

func TestKickUserFromConversationSlackResponseWithErrors(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected KickUserFromConversationSlackResponse
	}{
		{
			name:  "ResponseWithEmptyErrorsObject",
			input: `{"ok":true,"errors":{}}`,
			expected: KickUserFromConversationSlackResponse{
				Ok:     true,
				Errors: nil,
			},
		},
		{
			name:  "ResponseWithErrorsArray",
			input: `{"ok":false,"error":"some_error","errors":[]}`,
			expected: KickUserFromConversationSlackResponse{
				Ok:     false,
				Error:  "some_error",
				Errors: []SlackResponseErrors{},
			},
		},
		{
			name:  "ResponseWithoutErrors",
			input: `{"ok":true}`,
			expected: KickUserFromConversationSlackResponse{
				Ok: true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var response KickUserFromConversationSlackResponse
			err := json.Unmarshal([]byte(tt.input), &response)
			if err != nil {
				t.Fatalf("Unmarshal failed: %v", err)
			}

			if response.Ok != tt.expected.Ok {
				t.Errorf("got Ok=%v; want Ok=%v", response.Ok, tt.expected.Ok)
			}
			if response.Error != tt.expected.Error {
				t.Errorf("got Error=%q; want Error=%q", response.Error, tt.expected.Error)
			}

			if tt.expected.Errors == nil {
				if response.Errors != nil {
					t.Error("expected nil Errors, got non-nil")
				}
			} else {
				if response.Errors == nil {
					t.Error("expected non-nil Errors, got nil")
					return
				}
				if len(tt.expected.Errors) != len(response.Errors) {
					t.Errorf("got Errors length=%d; want length=%d", len(response.Errors), len(tt.expected.Errors))
				}
			}
		})
	}
}
