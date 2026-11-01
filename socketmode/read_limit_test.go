package socketmode

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/require"

	"github.com/slack-go/slack"
)

func TestSocketModeReadLimit(t *testing.T) {
	const limit = 128
	for _, tt := range []struct {
		name       string
		options    []Option
		size       int
		fragmented bool
		wantLimit  bool
	}{
		{name: "unlimited by default", size: 2 * limit},
		{name: "zero disables limit", options: []Option{OptionReadLimit(0)}, size: 2 * limit},
		{name: "negative disables limit", options: []Option{OptionReadLimit(-1)}, size: 2 * limit},
		{name: "below limit", options: []Option{OptionReadLimit(limit)}, size: limit - 1},
		{name: "at limit", options: []Option{OptionReadLimit(limit)}, size: limit},
		{name: "over limit", options: []Option{OptionReadLimit(limit)}, size: limit + 1, wantLimit: true},
		{name: "fragmented over limit", options: []Option{OptionReadLimit(limit)}, size: limit + 1, fragmented: true, wantLimit: true},
		{name: "custom dialer", options: []Option{OptionReadLimit(limit), OptionDialer(&websocket.Dialer{})}, size: limit + 1, wantLimit: true},
	} {
		t.Run(tt.name, func(t *testing.T) {
			const prefix = `{"type":"hello","padding":"`
			message := prefix + strings.Repeat("x", tt.size-len(prefix)-2) + `"}`
			upgrader := websocket.Upgrader{CheckOrigin: func(*http.Request) bool { return true }}
			if tt.fragmented {
				upgrader.WriteBufferSize = 64
			}
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path == "/apps.connections.open" {
					w.Header().Set("Content-Type", "application/json")
					json.NewEncoder(w).Encode(map[string]any{"ok": true, "url": "ws://" + r.Host + "/ws"})
					return
				}
				conn, err := upgrader.Upgrade(w, r, nil)
				if err != nil {
					t.Error(err)
					return
				}
				defer conn.Close()
				// NextWriter uses the small write buffer to split the fragmented case.
				writer, err := conn.NextWriter(websocket.TextMessage)
				if err != nil {
					t.Error(err)
					return
				}
				_, _ = writer.Write([]byte(message))
				_ = writer.Close()
			}))
			defer server.Close()

			client := New(slack.New("test-token", slack.OptionAPIURL(server.URL+"/")), tt.options...)
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			for range 2 { // The same options must apply to replacement connections.
				_, conn, err := client.openAndDial(ctx, nil)
				require.NoError(t, err)
				require.NoError(t, conn.SetReadDeadline(time.Now().Add(5*time.Second)))
				events := make(chan json.RawMessage, 1)
				err = client.receiveMessagesInto(ctx, conn, events)
				conn.Close()
				if tt.wantLimit {
					require.ErrorIs(t, err, websocket.ErrReadLimit)
					require.Empty(t, events, "oversized messages must not reach handlers")
				} else {
					require.NoError(t, err)
					require.Len(t, events, 1)
					require.Equal(t, message, string(<-events))
				}
			}
		})
	}
}
