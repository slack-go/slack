package slack

import (
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStartSocketModeContextDebugMalformedURL(t *testing.T) {
	client := httpClientFunc(func(req *http.Request) (*http.Response, error) {
		defer req.Body.Close()
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(strings.NewReader(`{"ok":true,"url":"://bad"}`)),
		}, nil
	})
	api := New("test-token", OptionHTTPClient(client), OptionDebug(true))

	connection, websocketURL, err := api.StartSocketModeContext(t.Context())
	require.NoError(t, err)
	require.NotNil(t, connection)
	assert.Equal(t, "://bad", websocketURL)
	assert.Equal(t, websocketURL, connection.URL)
}
