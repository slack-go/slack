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

func newSetAppIconTestClient(t *testing.T, handler http.HandlerFunc) *Client {
	t.Helper()
	ts := httptest.NewServer(handler)
	t.Cleanup(ts.Close)
	return New("client-token", OptionAPIURL(ts.URL+"/"))
}

func TestSetAppIconURL(t *testing.T) {
	api := newSetAppIconTestClient(t, func(rw http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/apps.icon.set", r.URL.Path)
		assert.Equal(t, "application/x-www-form-urlencoded", r.Header.Get("Content-Type"))
		require.NoError(t, r.ParseForm())
		assert.Equal(t, "caller-token", r.FormValue("token"))
		assert.Equal(t, "A123", r.FormValue("app_id"))
		assert.Equal(t, "https://example.com/icon.png", r.FormValue("url"))

		rw.Header().Set("Content-Type", "application/json")
		_, _ = rw.Write([]byte(`{"ok": true}`))
	})

	err := api.SetAppIcon("caller-token", SetAppIconParameters{
		AppID: "A123",
		URL:   "https://example.com/icon.png",
	})
	require.NoError(t, err)
}

func TestSetAppIconFile(t *testing.T) {
	api := newSetAppIconTestClient(t, func(rw http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "Bearer caller-token", r.Header.Get("Authorization"))
		require.NoError(t, r.ParseMultipartForm(1024))
		assert.Equal(t, "A123", r.FormValue("app_id"))
		assert.Empty(t, r.FormValue("url"))

		file, header, err := r.FormFile("file")
		require.NoError(t, err)
		defer file.Close()
		contents, err := io.ReadAll(file)
		require.NoError(t, err)
		assert.Equal(t, "app-icon.png", header.Filename)
		assert.Equal(t, "image contents", string(contents))

		rw.Header().Set("Content-Type", "application/json")
		_, _ = rw.Write([]byte(`{"ok": true}`))
	})

	err := api.SetAppIconContext(context.Background(), "caller-token", SetAppIconParameters{
		AppID:    "A123",
		File:     strings.NewReader("image contents"),
		Filename: "app-icon.png",
	})
	require.NoError(t, err)
}

func TestSetAppIconDefaultFilename(t *testing.T) {
	api := newSetAppIconTestClient(t, func(rw http.ResponseWriter, r *http.Request) {
		require.NoError(t, r.ParseMultipartForm(1024))
		file, header, err := r.FormFile("file")
		require.NoError(t, err)
		defer file.Close()
		assert.Equal(t, "icon.png", header.Filename)

		rw.Header().Set("Content-Type", "application/json")
		_, _ = rw.Write([]byte(`{"ok": true}`))
	})

	err := api.SetAppIconContext(context.Background(), "caller-token", SetAppIconParameters{
		AppID: "A123",
		File:  strings.NewReader("image contents"),
	})
	require.NoError(t, err)
}

func TestSetAppIconValidation(t *testing.T) {
	tests := []struct {
		name   string
		token  string
		params SetAppIconParameters
		error  string
	}{
		{
			name:   "missing token",
			params: SetAppIconParameters{AppID: "A123", URL: "https://example.com/icon.png"},
			error:  "apps.icon.set: token cannot be empty",
		},
		{
			name:  "missing app ID",
			token: "caller-token",
			params: SetAppIconParameters{
				URL: "https://example.com/icon.png",
			},
			error: "apps.icon.set: app ID cannot be empty",
		},
		{
			name:   "missing image source",
			token:  "caller-token",
			params: SetAppIconParameters{AppID: "A123"},
			error:  "apps.icon.set: exactly one of URL or file must be provided",
		},
		{
			name:  "multiple image sources",
			token: "caller-token",
			params: SetAppIconParameters{
				AppID: "A123",
				URL:   "https://example.com/icon.png",
				File:  strings.NewReader("image contents"),
			},
			error: "apps.icon.set: exactly one of URL or file must be provided",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			api := New("client-token")
			err := api.SetAppIconContext(context.Background(), test.token, test.params)
			assert.EqualError(t, err, test.error)
		})
	}
}

func TestSetAppIconErrors(t *testing.T) {
	t.Run("URL", func(t *testing.T) {
		api := newSetAppIconTestClient(t, func(rw http.ResponseWriter, r *http.Request) {
			rw.Header().Set("Content-Type", "application/json")
			_, _ = rw.Write([]byte(`{"ok": false, "error": "icon_not_accessible"}`))
		})

		err := api.SetAppIcon("caller-token", SetAppIconParameters{
			AppID: "A123",
			URL:   "https://example.com/missing.png",
		})
		assert.EqualError(t, err, "icon_not_accessible")
	})

	t.Run("file", func(t *testing.T) {
		api := newSetAppIconTestClient(t, func(rw http.ResponseWriter, r *http.Request) {
			require.NoError(t, r.ParseMultipartForm(1024))
			rw.Header().Set("Content-Type", "application/json")
			_, _ = rw.Write([]byte(`{"ok": false, "error": "error_bad_format"}`))
		})

		err := api.SetAppIcon("caller-token", SetAppIconParameters{
			AppID: "A123",
			File:  strings.NewReader("not an image"),
		})
		assert.EqualError(t, err, "error_bad_format")
	})
}
