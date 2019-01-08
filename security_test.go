package slack

import (
	"io"
	"log"
	"net/http"
	"testing"
)

const (
	validSigningSecret   = "e6b19c573432dcc6b075501d51b51bb8"
	invalidSigningSecret = "e6b19c573432dcc6b075501d51b51boo"
	validBody            = `{"token":"aF5ynEYQH0dFN9imlgcADxDB","team_id":"XXXXXXXXX","api_app_id":"YYYYYYYYY","event":{"type":"app_mention","user":"AAAAAAAAA","text":"<@EEEEEEEEE> hello world","client_msg_id":"477cc591-ch73-a14z-4db8-g0cd76321bec","ts":"1531431954.000073","channel":"TTTTTTTTT","event_ts":"1531431954.000073"},"type":"event_callback","event_id":"TvBP7LRED7","event_time":1531431954,"authed_users":["EEEEEEEEE"]}`
	invalidBody          = `{"token":"12345678abcdefghlmnopqrs","team_id":"XXXXXXXXX","api_app_id":"YYYYYYYYY","event":{"type":"app_mention","user":"AAAAAAAAA","text":"<@EEEEEEEEE> hello world","client_msg_id":"477cc591-ch73-a14z-4db8-g0cd76321bec","ts":"1531431954.000073","channel":"TTTTTTTTT","event_ts":"1531431954.000073"},"type":"event_callback","event_id":"TvBP7LRED7","event_time":1531431954,"authed_users":["EEEEEEEEE"]}`
)

func newHeader(valid bool) http.Header {
	h := http.Header{}
	if valid {
		h.Set("X-Slack-Signature", "v0=adada4ed31709aef585c2580ca3267678c6a8eaeb7e0c1aca3ee57b656886b2c")
		h.Set("X-Slack-Request-Timestamp", "1531431954")
	} else {
		h.Set("X-Slack-Signature", "")
	}
	return h
}

func TestExpiredTimestamp(t *testing.T) {
	_, err := NewSecretsVerifier(newHeader(true), "abcdefg12345")
	if err == nil {
		t.Fatal("expected an error but got none")
	}
}

func TestUnsafeSignatureVerifier(t *testing.T) {
	tests := []struct {
		title         string
		header        http.Header
		signingSecret string
		expectError   bool
	}{
		{
			title:         "Testing with acceptable params",
			header:        newHeader(true),
			signingSecret: "abcdefg12345",
			expectError:   false,
		},
		{
			title:         "Testing with unacceptable params",
			header:        newHeader(false),
			signingSecret: "abcdefg12345",
			expectError:   true,
		},
	}

	for _, test := range tests {
		_, err := unsafeSignatureVerifier(test.header, test.signingSecret)

		if !test.expectError && err != nil {
			log.Fatalf("%s: Unexpected error: %s in test", test.title, err)
		} else if test.expectError == true && err == nil {
			log.Fatalf("Expected error but got none")
		}
	}
}

func TestEnsure(t *testing.T) {
	tests := []struct {
		title         string
		header        http.Header
		signingSecret string
		body          string
		expectError   bool
	}{
		{
			title:         "Testing with acceptable signing secret and valid body",
			header:        newHeader(true),
			signingSecret: validSigningSecret,
			body:          validBody,
			expectError:   false,
		},
		{
			title:         "Testing with unacceptable signing secret and valid body",
			header:        newHeader(true),
			signingSecret: invalidSigningSecret,
			body:          validBody,
			expectError:   true,
		},
		{
			title:         "Testing with acceptable signing secret and invalid body",
			header:        newHeader(true),
			signingSecret: validSigningSecret,
			body:          invalidBody,
			expectError:   true,
		},
	}

	for _, test := range tests {
		sv, err := unsafeSignatureVerifier(test.header, test.signingSecret)
		if err != nil {
			t.Fatalf("unexpected error: %s", err)
		}
		io.WriteString(&sv, test.body)

		err = sv.Ensure()

		if !test.expectError && err != nil {
			log.Fatalf("%s: Unexpected error: %s in test", test.title, err)
		} else if test.expectError == true && err == nil {
			log.Fatalf("Expected error but got none")
		}
	}

}
