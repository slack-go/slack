package slack

import (
	"context"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestListAllRetriesRateLimitedFirstPage(t *testing.T) {
	for _, method := range []string{"conversations.list", "stars.list"} {
		t.Run(method, func(t *testing.T) {
			var calls atomic.Int32
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, "/"+method, r.URL.Path)
				if calls.Add(1) == 1 {
					w.Header().Set("Retry-After", "0")
					w.WriteHeader(http.StatusTooManyRequests)
					return
				}
				assert.Empty(t, r.FormValue("cursor"))
				w.Header().Set("Content-Type", "application/json")
				_, _ = w.Write([]byte(`{"ok":true,"channels":[{"id":"C1"}],"items":[{"type":"message","channel":"C1","ts":"1"}],"response_metadata":{"next_cursor":""}}`))
			}))
			t.Cleanup(server.Close)
			api := New("test-token", OptionAPIURL(server.URL+"/"))
			if method == "conversations.list" {
				channels, err := api.GetAllConversationsContext(context.Background())
				require.NoError(t, err)
				require.Len(t, channels, 1)
				require.Equal(t, "C1", channels[0].ID)
			} else {
				items, err := api.ListAllStarsContext(context.Background())
				require.NoError(t, err)
				require.Len(t, items, 1)
				require.Equal(t, "C1", items[0].Channel)
			}
			require.EqualValues(t, 2, calls.Load())
		})
	}
}
