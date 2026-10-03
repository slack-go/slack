package slack

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newSetAppIconTestClient(t *testing.T, handler http.HandlerFunc, options ...Option) *Client {
	t.Helper()
	ts := httptest.NewServer(handler)
	t.Cleanup(ts.Close)
	return New("client-token", append(options, OptionAPIURL(ts.URL+"/"))...)
}

// The handlers run on the server goroutine, so they use assert, not require.

func TestSetAppIconURL(t *testing.T) {
	api := newSetAppIconTestClient(t, func(rw http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/apps.icon.set", r.URL.Path)
		assert.Equal(t, "application/x-www-form-urlencoded", r.Header.Get("Content-Type"))
		assert.NoError(t, r.ParseForm())
		assert.Equal(t, "caller-token", r.FormValue("token"))
		assert.Equal(t, "A123", r.FormValue("app_id"))
		assert.Equal(t, "https://example.com/icon.png", r.FormValue("url"))

		rw.Header().Set("Content-Type", "application/json")
		_, _ = rw.Write([]byte(`{"ok": true}`))
	})

	err := api.SetAppIcon("caller-token", "A123", SetAppIconParameters{
		URL: "https://example.com/icon.png",
	})
	require.NoError(t, err)
}

func TestSetAppIconFile(t *testing.T) {
	api := newSetAppIconTestClient(t, func(rw http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "Bearer caller-token", r.Header.Get("Authorization"))
		if !assert.NoError(t, r.ParseMultipartForm(1024)) {
			return
		}
		assert.Equal(t, "A123", r.FormValue("app_id"))
		assert.Empty(t, r.FormValue("url"))

		file, header, err := r.FormFile("file")
		if !assert.NoError(t, err) {
			return
		}
		defer file.Close()
		contents, err := io.ReadAll(file)
		assert.NoError(t, err)
		assert.Equal(t, "app-icon.png", header.Filename)
		assert.Equal(t, "image contents", string(contents))

		rw.Header().Set("Content-Type", "application/json")
		_, _ = rw.Write([]byte(`{"ok": true}`))
	})

	err := api.SetAppIconContext(context.Background(), "caller-token", "A123", SetAppIconParameters{
		File:     strings.NewReader("image contents"),
		Filename: "app-icon.png",
	})
	require.NoError(t, err)
}

func TestSetAppIconDefaultFilename(t *testing.T) {
	api := newSetAppIconTestClient(t, func(rw http.ResponseWriter, r *http.Request) {
		if !assert.NoError(t, r.ParseMultipartForm(1024)) {
			return
		}
		file, header, err := r.FormFile("file")
		if !assert.NoError(t, err) {
			return
		}
		defer file.Close()
		assert.Equal(t, "icon.png", header.Filename)

		rw.Header().Set("Content-Type", "application/json")
		_, _ = rw.Write([]byte(`{"ok": true}`))
	})

	err := api.SetAppIconContext(context.Background(), "caller-token", "A123", SetAppIconParameters{
		File: strings.NewReader("image contents"),
	})
	require.NoError(t, err)
}

func TestSetAppIconConfigTokenFallback(t *testing.T) {
	api := newSetAppIconTestClient(t, func(rw http.ResponseWriter, r *http.Request) {
		assert.NoError(t, r.ParseForm())
		assert.Equal(t, "config-token", r.FormValue("token"))

		rw.Header().Set("Content-Type", "application/json")
		_, _ = rw.Write([]byte(`{"ok": true}`))
	}, OptionConfigToken("config-token"))

	err := api.SetAppIcon("", "A123", SetAppIconParameters{
		URL: "https://example.com/icon.png",
	})
	require.NoError(t, err)
}

func TestSetAppIconValidation(t *testing.T) {
	tests := []struct {
		name  string
		token string
		appID string
		icon  SetAppIconParameters
		error string
	}{
		{
			name:  "missing token and no configuration token",
			appID: "A123",
			icon:  SetAppIconParameters{URL: "https://example.com/icon.png"},
			error: "apps.icon.set: token cannot be empty",
		},
		{
			name:  "missing app ID",
			token: "caller-token",
			icon:  SetAppIconParameters{URL: "https://example.com/icon.png"},
			error: "apps.icon.set: app ID cannot be empty",
		},
		{
			name:  "missing image source",
			token: "caller-token",
			appID: "A123",
			error: "apps.icon.set: exactly one of URL or file must be provided",
		},
		{
			name:  "multiple image sources",
			token: "caller-token",
			appID: "A123",
			icon: SetAppIconParameters{
				URL:  "https://example.com/icon.png",
				File: strings.NewReader("image contents"),
			},
			error: "apps.icon.set: exactly one of URL or file must be provided",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			api := New("client-token")
			err := api.SetAppIconContext(context.Background(), test.token, test.appID, test.icon)
			assert.EqualError(t, err, test.error)
		})
	}
}

func TestSetAppIconErrors(t *testing.T) {
	t.Run("URL", func(t *testing.T) {
		api := newSetAppIconTestClient(t, func(rw http.ResponseWriter, r *http.Request) {
			rw.Header().Set("Content-Type", "application/json")
			_, _ = rw.Write([]byte(`{"ok": false, "error": "invalid_icon_size"}`))
		})

		err := api.SetAppIcon("caller-token", "A123", SetAppIconParameters{
			URL: "https://example.com/small.png",
		})
		assert.EqualError(t, err, "invalid_icon_size")
	})

	t.Run("file", func(t *testing.T) {
		api := newSetAppIconTestClient(t, func(rw http.ResponseWriter, r *http.Request) {
			assert.NoError(t, r.ParseMultipartForm(1024))
			rw.Header().Set("Content-Type", "application/json")
			_, _ = rw.Write([]byte(`{"ok": false, "error": "error_bad_format"}`))
		})

		err := api.SetAppIcon("caller-token", "A123", SetAppIconParameters{
			File: strings.NewReader("not an image"),
		})
		assert.EqualError(t, err, "error_bad_format")
	})
}
