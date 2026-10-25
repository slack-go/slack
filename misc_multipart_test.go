package slack

import (
	"context"
	"errors"
	"io"
	"net"
	"net/http"
	"net/url"
	"strings"
	"testing"
	"testing/iotest"
	"testing/synctest"
)

func TestMultipartUploadRequestFailure(t *testing.T) {
	dialErr := errors.New("connection failed")
	tests := []struct {
		name     string
		path     string
		canceled bool
		wantErr  error
	}{
		{name: "invalid request URL", path: "://invalid"},
		{name: "connection failure", path: "https://example.test/upload", wantErr: dialErr},
		{name: "canceled request", path: "https://example.test/upload", canceled: true, wantErr: context.Canceled},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// synctest fails if the upload returns with a writer still blocked.
			synctest.Test(t, func(t *testing.T) {
				transport := &http.Transport{
					DialContext: func(context.Context, string, string) (net.Conn, error) {
						return nil, dialErr
					},
				}
				defer transport.CloseIdleConnections()
				ctx, cancel := context.WithCancel(t.Context())
				defer cancel()
				if tt.canceled {
					cancel()
				}

				err := postWithMultipartResponse(ctx, &http.Client{Transport: transport}, tt.path,
					"file.txt", "file", "test-token", url.Values{"title": {"test"}}, strings.NewReader("contents"), nil, discard{})
				if err == nil {
					t.Fatal("expected the upload to fail")
				}
				if tt.wantErr != nil && !errors.Is(err, tt.wantErr) {
					t.Fatalf("got %v, want %v", err, tt.wantErr)
				}
			})
		})
	}
}

type multipartTestClient func(*http.Request) (*http.Response, error)

func (f multipartTestClient) Do(req *http.Request) (*http.Response, error) {
	return f(req)
}

func TestMultipartUploadReaderError(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		readErr := errors.New("file read failed")
		client := multipartTestClient(func(req *http.Request) (*http.Response, error) {
			defer req.Body.Close()
			if _, err := io.Copy(io.Discard, req.Body); err != nil {
				return nil, err
			}
			return &http.Response{StatusCode: http.StatusOK, Body: io.NopCloser(strings.NewReader("OK"))}, nil
		})

		err := postWithMultipartResponse(t.Context(), client, "https://example.test/upload",
			"file.txt", "file", "test-token", nil, iotest.ErrReader(readErr), nil, discard{})
		if !errors.Is(err, readErr) {
			t.Fatalf("got %v, want the source reader error %v", err, readErr)
		}
	})
}
