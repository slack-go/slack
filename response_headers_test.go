package slack

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestResponseHeaders(t *testing.T) {
	t.Run("AuthTest captures headers", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.Header().Set("X-OAuth-Scopes", "users:read,channels:read")
			w.Header().Set("X-Accepted-OAuth-Scopes", "users:read")
			w.Write([]byte(`{"ok":true,"url":"https://example.slack.com","team":"T","user":"U","team_id":"T1","user_id":"U1"}`))
		}))
		defer ts.Close()

		api := New("test-token", OptionAPIURL(ts.URL+"/"))
		resp, err := api.AuthTestContext(t.Context())
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got := resp.Header.Get("X-OAuth-Scopes"); got != "users:read,channels:read" {
			t.Fatalf("expected X-OAuth-Scopes=users:read,channels:read, got %q", got)
		}
		if got := resp.Header.Get("X-Accepted-OAuth-Scopes"); got != "users:read" {
			t.Fatalf("expected X-Accepted-OAuth-Scopes=users:read, got %q", got)
		}
	})

	t.Run("SlackResponse embeds headers via callback", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.Header().Set("X-OAuth-Scopes", "admin")
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
		if got := gotHeaders.Get("X-OAuth-Scopes"); got != "admin" {
			t.Fatalf("expected X-OAuth-Scopes=admin, got %q", got)
		}
	})
}

func TestOptionOnResponseHeaders(t *testing.T) {
	t.Run("callback fires", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.Header().Set("X-OAuth-Scopes", "users:read")
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
		if got := gotHeaders.Get("X-OAuth-Scopes"); got != "users:read" {
			t.Fatalf("expected X-OAuth-Scopes=users:read, got %q", got)
		}
	})

	t.Run("no callback is safe", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.Header().Set("X-OAuth-Scopes", "users:read")
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
			w.Header().Set("X-OAuth-Scopes", "users:read")
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
		if got := gotHeaders.Get("X-OAuth-Scopes"); got != "users:read" {
			t.Fatalf("expected X-OAuth-Scopes=users:read, got %q", got)
		}
	})
}
