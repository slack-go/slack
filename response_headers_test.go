package slack

import (
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestResponseHeaders(t *testing.T) {
	t.Run("AuthTest captures headers", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.Header().Set("X-Oauth-Scopes", "users:read,channels:read")
			w.Header().Set("X-Accepted-Oauth-Scopes", "users:read")
			w.Write([]byte(`{"ok":true,"url":"https://example.slack.com","team":"T","user":"U","team_id":"T1","user_id":"U1"}`))
		}))
		defer ts.Close()

		api := New("test-token", OptionAPIURL(ts.URL+"/"))
		resp, err := api.AuthTestContext(t.Context())
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got := resp.Header.Get("X-Oauth-Scopes"); got != "users:read,channels:read" {
			t.Fatalf("expected X-OAuth-Scopes=users:read,channels:read, got %q", got)
		}
		if got := resp.Header.Get("X-Accepted-Oauth-Scopes"); got != "users:read" {
			t.Fatalf("expected X-Accepted-OAuth-Scopes=users:read, got %q", got)
		}
	})

	t.Run("SlackResponse embeds headers via callback", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.Header().Set("X-Oauth-Scopes", "admin")
			w.Write([]byte(`{"ok":true,"url":"https://example.slack.com","team":"T","user":"U","team_id":"T1","user_id":"U1"}`))
		}))
		defer ts.Close()

		var gotHeaders http.Header
		api := New("test-token",
			OptionAPIURL(ts.URL+"/"),
			OptionOnResponseHeaders(func(path string, headers http.Header) {
				gotHeaders = headers
			}),
		)

		_, err := api.AuthTestContext(t.Context())
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if gotHeaders == nil {
			t.Fatal("expected headers from SlackResponse, got nil")
		}
		if got := gotHeaders.Get("X-Oauth-Scopes"); got != "admin" {
			t.Fatalf("expected X-OAuth-Scopes=admin, got %q", got)
		}
	})
}

func TestResponseCallbacksForChatMultipartAndAudit(t *testing.T) {
	photo, _ := createUserPhoto(t)
	tests := []struct {
		method   string
		response string
		call     func(*Client) error
	}{
		{"chat.postMessage", `{"ok":true,"channel":"C123","ts":"1.0","warning":"missing_charset"}`, func(api *Client) error {
			_, _, err := api.PostMessage("C123", MsgOptionText("hello", false))
			return err
		}},
		{"apps.icon.set", `{"ok":true,"warning":"missing_charset"}`, func(api *Client) error {
			return api.SetAppIcon("config-token", "A123", SetAppIconParameters{File: strings.NewReader("icon")})
		}},
		{"users.setPhoto", `{"ok":true,"warning":"missing_charset"}`, func(api *Client) error {
			return api.SetUserPhoto(photo, UserSetPhotoParams{})
		}},
		{"audit/v1/logs", `{"entries":[],"warning":"missing_charset"}`, func(api *Client) error {
			_, _, err := api.GetAuditLogs(AuditLogParameters{})
			return err
		}},
	}

	for _, test := range tests {
		t.Run(test.method, func(t *testing.T) {
			client := httpClientFunc(func(req *http.Request) (*http.Response, error) {
				if req.Body != nil {
					defer req.Body.Close()
					if _, err := io.Copy(io.Discard, req.Body); err != nil {
						return nil, err
					}
				}
				response := httptest.NewRecorder()
				response.Header().Set("X-Slack-Req-Id", "request-id")
				_, _ = io.WriteString(response, test.response)
				return response.Result(), nil
			})

			var headerMethod, warningMethod, requestID string
			var warning *Warning
			api := New("test-token", OptionHTTPClient(client),
				OptionOnResponseHeaders(func(path string, headers http.Header) {
					headerMethod = path
					requestID = headers.Get("X-Slack-Req-Id")
				}),
				OptionOnWarning(func(path string, request any, w *Warning) {
					warningMethod = path
					warning = w
					values, ok := request.(url.Values)
					require.True(t, ok)
					assert.NotContains(t, values, "token")
				}),
			)

			require.NoError(t, test.call(api))
			assert.Equal(t, test.method, headerMethod)
			assert.Equal(t, "request-id", requestID)
			assert.Equal(t, test.method, warningMethod)
			require.NotNil(t, warning)
			assert.Equal(t, []string{"missing_charset"}, warning.Codes)
		})
	}
}

func TestOptionOnResponseHeaders(t *testing.T) {
	t.Run("callback fires", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.Header().Set("X-Oauth-Scopes", "users:read")
			w.Write([]byte(`{"ok":true,"url":"https://example.slack.com","team":"T","user":"U","team_id":"T1","user_id":"U1"}`))
		}))
		defer ts.Close()

		var gotPath string
		var gotHeaders http.Header
		api := New("test-token",
			OptionAPIURL(ts.URL+"/"),
			OptionOnResponseHeaders(func(path string, headers http.Header) {
				gotPath = path
				gotHeaders = headers
			}),
		)

		_, err := api.AuthTestContext(t.Context())
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if gotPath != "auth.test" {
			t.Fatalf("expected path auth.test, got %q", gotPath)
		}
		if gotHeaders == nil {
			t.Fatal("expected headers, got nil")
		}
		if got := gotHeaders.Get("X-Oauth-Scopes"); got != "users:read" {
			t.Fatalf("expected X-OAuth-Scopes=users:read, got %q", got)
		}
	})

	t.Run("no callback is safe", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.Header().Set("X-Oauth-Scopes", "users:read")
			w.Write([]byte(`{"ok":true,"url":"https://example.slack.com","team":"T","user":"U","team_id":"T1","user_id":"U1"}`))
		}))
		defer ts.Close()

		api := New("test-token", OptionAPIURL(ts.URL+"/"))
		_, err := api.AuthTestContext(t.Context())
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("callback fires on error", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.Header().Set("X-Oauth-Scopes", "users:read")
			w.Write([]byte(`{"ok":false,"error":"invalid_auth"}`))
		}))
		defer ts.Close()

		var gotHeaders http.Header
		api := New("test-token",
			OptionAPIURL(ts.URL+"/"),
			OptionOnResponseHeaders(func(path string, headers http.Header) {
				gotHeaders = headers
			}),
		)

		_, err := api.AuthTestContext(t.Context())
		if err == nil {
			t.Fatal("expected error")
		}
		if gotHeaders == nil {
			t.Fatal("expected headers even on error response, got nil")
		}
		if got := gotHeaders.Get("X-Oauth-Scopes"); got != "users:read" {
			t.Fatalf("expected X-OAuth-Scopes=users:read, got %q", got)
		}
	})
}
