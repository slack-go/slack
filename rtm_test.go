package slack

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"slices"
	"testing"

	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/require"
)

func TestRTMOptionDialHeaders(t *testing.T) {
	for _, tt := range []struct {
		name     string
		noOption bool
		headers  http.Header
		want     http.Header
	}{
		{
			name: "default", noOption: true,
			want: http.Header{"Origin": {"https://api.slack.com"}, "X-Custom-Header": nil, "X-Test": nil},
		},
		{
			name: "nil headers",
			want: http.Header{"Origin": {"https://api.slack.com"}, "X-Custom-Header": nil, "X-Test": nil},
		},
		{
			name: "custom headers",
			headers: http.Header{
				"X-Custom-Header": {"example"}, "X-Test": {"first", "second"},
			},
			want: http.Header{
				"Origin":          {"https://api.slack.com"},
				"X-Custom-Header": {"example"}, "X-Test": {"sent", "second"},
			},
		},
		{
			name:    "empty origin uses default",
			headers: http.Header{"Origin": {""}},
			want:    http.Header{"Origin": {"https://api.slack.com"}, "X-Custom-Header": nil, "X-Test": nil},
		},
		{
			name:    "custom origin",
			headers: http.Header{"Origin": {"https://example.test"}},
			want:    http.Header{"Origin": {"https://example.test"}, "X-Custom-Header": nil, "X-Test": nil},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			upgrader := websocket.Upgrader{CheckOrigin: func(*http.Request) bool { return true }}
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path == "/rtm.connect" {
					w.Header().Set("Content-Type", "application/json")
					if err := json.NewEncoder(w).Encode(map[string]any{"ok": true, "url": "ws://" + r.Host + "/ws"}); err != nil {
						t.Error(err)
					}
					return
				}
				for name, want := range tt.want {
					if got := r.Header.Values(name); !slices.Equal(got, want) {
						t.Errorf("%s header = %q, want %q", name, got, want)
					}
				}
				conn, err := upgrader.Upgrade(w, r, nil)
				if err != nil {
					t.Error(err)
					return
				}
				_ = conn.Close()
			}))
			t.Cleanup(server.Close)
			dialer := &websocket.Dialer{Proxy: func(r *http.Request) (*url.URL, error) {
				// Mutating one handshake must not affect the next connection.
				if values := r.Header.Values("X-Test"); len(values) > 0 {
					if !slices.Equal(values, []string{"first", "second"}) {
						t.Errorf("headers changed between dials: %q", values)
					}
					values[0] = "sent"
				}
				return nil, nil
			}}
			headers := tt.headers.Clone()
			options := []RTMOption{RTMOptionDialer(dialer)}
			if !tt.noOption {
				options = append(options, RTMOptionDialHeaders(headers))
			}
			rtm := New("test-token", OptionAPIURL(server.URL+"/")).NewRTM(options...)
			t.Cleanup(func() { rtm.pingDeadman.Stop() })
			if headers != nil {
				headers.Set("X-Custom-Header", "caller mutation")
				if values := headers.Values("X-Test"); len(values) > 0 {
					values[0] = "caller mutation"
				}
			}
			for range 2 {
				_, conn, err := rtm.startRTMAndDial(false)
				require.NoError(t, err)
				_ = conn.Close()
			}
		})
	}
}
