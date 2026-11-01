package slack

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestListAllStarsReturnsAPIErrors(t *testing.T) {
	for _, failOnPage := range []int32{1, 2} {
		t.Run(strconv.Itoa(int(failOnPage)), func(t *testing.T) {
			var calls atomic.Int32
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				if calls.Add(1) == failOnPage {
					_, _ = w.Write([]byte(`{"ok":false,"error":"missing_scope","needed":"stars:read","provided":"chat:write"}`))
					return
				}
				_ = json.NewEncoder(w).Encode(map[string]any{
					"ok":                true,
					"items":             []Item{{Type: "message", Channel: "C1"}},
					"response_metadata": ResponseMetadata{Cursor: "page2"},
				})
			}))
			t.Cleanup(server.Close)
			api := New("test-token", OptionAPIURL(server.URL+"/"))
			items, err := api.ListAllStarsContext(context.Background())
			require.EqualError(t, err, "missing_scope")
			var responseErr SlackErrorResponse
			require.ErrorAs(t, err, &responseErr)
			require.Equal(t, "stars:read", responseErr.Needed)
			require.Equal(t, "chat:write", responseErr.Provided)
			require.Len(t, items, int(failOnPage-1))
			require.Equal(t, failOnPage, calls.Load())
		})
	}
}
