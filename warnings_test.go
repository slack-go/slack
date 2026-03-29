package slack

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestWarn(t *testing.T) {
	t.Run("top-level only", func(t *testing.T) {
		resp := SlackResponse{Warning: "missing_charset"}
		got := resp.Warn()
		if got == nil {
			t.Fatal("expected warning, got nil")
		}
		if len(got.Codes) != 1 || got.Codes[0] != "missing_charset" {
			t.Fatalf("expected codes [missing_charset], got %v", got.Codes)
		}
	})

	t.Run("metadata only", func(t *testing.T) {
		resp := SlackResponse{
			ResponseMetadata: ResponseMetadata{
				Warnings: []string{"superfluous_charset"},
			},
		}
		got := resp.Warn()
		if got == nil {
			t.Fatal("expected warning, got nil")
		}
		if len(got.Warnings) != 1 || got.Warnings[0] != "superfluous_charset" {
			t.Fatalf("expected warnings [superfluous_charset], got %v", got.Warnings)
		}
	})

	t.Run("both sources", func(t *testing.T) {
		resp := SlackResponse{
			Warning: "missing_charset,deprecated",
			ResponseMetadata: ResponseMetadata{
				Warnings: []string{"missing_charset", "other"},
			},
		}
		got := resp.Warn()
		if got == nil {
			t.Fatal("expected warning, got nil")
		}
		if len(got.Codes) != 2 {
			t.Fatalf("expected 2 codes, got %v", got.Codes)
		}
		if len(got.Warnings) != 2 {
			t.Fatalf("expected 2 warnings, got %v", got.Warnings)
		}
	})

	t.Run("no warnings", func(t *testing.T) {
		resp := SlackResponse{Ok: true}
		got := resp.Warn()
		if got != nil {
			t.Fatalf("expected nil, got %+v", got)
		}
	})
}

func TestOptionOnWarning(t *testing.T) {
	t.Run("callback fires on warning", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(`{
				"ok": true,
				"warning": "missing_charset",
				"response_metadata": {
					"warnings": ["missing_charset"]
				}
			}`))
		}))
		defer ts.Close()

		var gotWarning *Warning
		var gotPath string
		var gotRequest any
		api := New("test-token",
			OptionAPIURL(ts.URL+"/"),
			OptionOnWarning(func(path string, request any, w *Warning) {
				gotWarning = w
				gotPath = path
				gotRequest = request
			}),
		)

		_, err := api.AuthTestContext(t.Context())
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if gotWarning == nil {
			t.Fatal("expected warning, got nil")
		}
		if len(gotWarning.Codes) != 1 || gotWarning.Codes[0] != "missing_charset" {
			t.Fatalf("expected codes [missing_charset], got %v", gotWarning.Codes)
		}
		if gotPath != "auth.test" {
			t.Fatalf("expected path auth.test, got %s", gotPath)
		}
		if _, ok := gotRequest.(url.Values); !ok {
			t.Fatalf("expected url.Values request, got %T", gotRequest)
		}
	})

	t.Run("callback not fired when no warnings", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(`{"ok": true}`))
		}))
		defer ts.Close()

		called := false
		api := New("test-token",
			OptionAPIURL(ts.URL+"/"),
			OptionOnWarning(func(path string, request any, w *Warning) {
				called = true
			}),
		)

		_, err := api.AuthTestContext(t.Context())
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if called {
			t.Fatal("callback should not have been called")
		}
	})

	t.Run("callback fires on error response with warnings", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(`{
				"ok": false,
				"error": "invalid_auth",
				"warning": "missing_charset",
				"response_metadata": {
					"warnings": ["missing_charset"]
				}
			}`))
		}))
		defer ts.Close()

		var gotWarning *Warning
		api := New("test-token",
			OptionAPIURL(ts.URL+"/"),
			OptionOnWarning(func(path string, request any, w *Warning) {
				gotWarning = w
			}),
		)

		_, err := api.AuthTestContext(t.Context())
		if err == nil {
			t.Fatal("expected error")
		}
		if gotWarning == nil {
			t.Fatal("expected warning, got nil")
		}
		if len(gotWarning.Codes) != 1 || gotWarning.Codes[0] != "missing_charset" {
			t.Fatalf("expected codes [missing_charset], got %v", gotWarning.Codes)
		}
	})

	t.Run("no callback registered is safe", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(`{
				"ok": true,
				"warning": "missing_charset"
			}`))
		}))
		defer ts.Close()

		api := New("test-token", OptionAPIURL(ts.URL+"/"))
		_, err := api.AuthTestContext(t.Context())
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})
}
