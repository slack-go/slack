package slack

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
	"time"
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

func TestNewSecretsVerifierRejectsEmptySigningSecret(t *testing.T) {
	_, err := NewSecretsVerifier(newHeader(true), "")
	if !errors.Is(err, ErrInvalidConfiguration) {
		t.Fatalf("expected ErrInvalidConfiguration, got %v", err)
	}
}

func TestSignatureVerifierWithoutTimestampCheckRejectsEmptySigningSecret(t *testing.T) {
	_, err := signatureVerifierWithoutTimestampCheck(newHeader(true), "")
	if !errors.Is(err, ErrInvalidConfiguration) {
		t.Fatalf("expected ErrInvalidConfiguration, got %v", err)
	}
}

func TestSignatureVerifierWithoutTimestampCheck(t *testing.T) {
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
		_, err := signatureVerifierWithoutTimestampCheck(test.header, test.signingSecret)

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
		sv, err := signatureVerifierWithoutTimestampCheck(test.header, test.signingSecret)
		if err != nil {
			t.Fatalf("unexpected error: %s", err)
		}
		io.WriteString(&sv, test.body)

		err = sv.Ensure()

		if !test.expectError && err != nil {
			log.Fatalf("%s: Unexpected error: %s in test", test.title, err)
		} else if test.expectError && !errors.Is(err, ErrInvalidSignature) {
			t.Fatalf("%s: expected ErrInvalidSignature, got %v", test.title, err)
		}
	}

}

func TestEnsureDoesNotDiscloseSignatures(t *testing.T) {
	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	mac := hmac.New(sha256.New, []byte(validSigningSecret))
	io.WriteString(mac, "v0:"+timestamp+":"+validBody)
	computedSignature := hex.EncodeToString(mac.Sum(nil))
	suppliedSignature := strings.Repeat("0", sha256.Size*2)

	for _, test := range []struct {
		name      string
		withDebug bool
		debug     bool
	}{
		{name: "without debug logger"},
		{name: "debug disabled", withDebug: true},
		{name: "debug enabled", withDebug: true, debug: true},
	} {
		t.Run(test.name, func(t *testing.T) {
			header := http.Header{}
			header.Set(hTimestamp, timestamp)
			header.Set(hSignature, "v0="+suppliedSignature)
			verifier, err := NewSecretsVerifier(header, validSigningSecret)
			if err != nil {
				t.Fatal(err)
			}

			var debugOutput bytes.Buffer
			if test.withDebug {
				verifier.WithDebug(New("", OptionDebug(test.debug), OptionLog(log.New(&debugOutput, "", 0))))
			}
			if _, err := io.WriteString(&verifier, validBody); err != nil {
				t.Fatal(err)
			}
			err = verifier.Ensure()
			if err == nil {
				t.Fatal("expected an invalid signature error")
			}
			if !errors.Is(err, ErrInvalidSignature) {
				t.Errorf("expected ErrInvalidSignature, got %v", err)
			}

			// Echoing the error must not give a caller a signature they can replay.
			response := httptest.NewRecorder()
			http.Error(response, err.Error(), http.StatusUnauthorized)
			for name, output := range map[string]string{
				"error response": response.Body.String(),
				"debug log":      debugOutput.String(),
			} {
				for _, sensitive := range []string{computedSignature, suppliedSignature, validSigningSecret, validBody} {
					if strings.Contains(output, sensitive) {
						t.Errorf("%s disclosed sensitive verification data", name)
					}
				}
			}
			if logged := debugOutput.Len() > 0; logged != test.debug {
				t.Errorf("debug output present = %t, want %t", logged, test.debug)
			}
			if test.debug {
				want := "Invalid signature: received prefix=00000000 (32 bytes), computed prefix=" + computedSignature[:8] + " (32 bytes)\n"
				if got := debugOutput.String(); got != want {
					t.Errorf("unexpected debug output: got %q, want %q", got, want)
				}
			}
		})
	}
}

func TestNewSecretsVerifierSignatureLength(t *testing.T) {
	for _, test := range []struct {
		size    int
		wantErr error
	}{
		{size: 0, wantErr: ErrInvalidSignature},
		{size: 1, wantErr: ErrInvalidSignature},
		{size: 4, wantErr: ErrInvalidSignature},
		{size: 31, wantErr: ErrInvalidSignature},
		{size: 32},
		{size: 33, wantErr: ErrInvalidSignature},
	} {
		t.Run(strconv.Itoa(test.size)+" bytes", func(t *testing.T) {
			header := newHeader(true)
			header.Set(hSignature, "v0="+strings.Repeat("ab", test.size))
			header.Set(hTimestamp, strconv.FormatInt(time.Now().Unix(), 10))
			_, err := NewSecretsVerifier(header, validSigningSecret)
			if !errors.Is(err, test.wantErr) {
				t.Fatalf("expected %v, got %v", test.wantErr, err)
			}
		})
	}
}

func TestNewSecretsVerifierSignatureVersion(t *testing.T) {
	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	mac := hmac.New(sha256.New, []byte(validSigningSecret))
	io.WriteString(mac, "v0:"+timestamp+":"+validBody)
	digest := hex.EncodeToString(mac.Sum(nil))

	for _, test := range []struct {
		name      string
		signature string
		wantErr   error
	}{
		{name: "v0", signature: "v0=" + digest},
		{name: "missing version", signature: digest, wantErr: ErrInvalidSignature},
		{name: "unsupported version", signature: "v1=" + digest, wantErr: ErrInvalidSignature},
		{name: "uppercase version", signature: "V0=" + digest, wantErr: ErrInvalidSignature},
		{name: "wrong separator", signature: "v0:" + digest, wantErr: ErrInvalidSignature},
	} {
		t.Run(test.name, func(t *testing.T) {
			header := http.Header{}
			header.Set(hTimestamp, timestamp)
			header.Set(hSignature, test.signature)
			verifier, err := NewSecretsVerifier(header, validSigningSecret)
			if !errors.Is(err, test.wantErr) {
				t.Fatalf("got %v, want %v", err, test.wantErr)
			}
			if err == nil {
				if _, err := io.WriteString(&verifier, validBody); err != nil {
					t.Fatal(err)
				}
				if err := verifier.Ensure(); err != nil {
					t.Fatalf("valid v0 signature failed verification: %v", err)
				}
			}
		})
	}
}
