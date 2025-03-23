package slack

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestServer provides a wrapper around httptest.Server with added functionality
// for Slack API testing
type TestServer struct {
	Server *httptest.Server
	Mux    *http.ServeMux
	URL    string
	Client *Client
}

// NewTestServer creates a new test server with a fresh mux and a configured Slack client
func NewTestServer(t *testing.T, token string) *TestServer {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)

	client := New(token, OptionAPIURL(server.URL+"/"))

	return &TestServer{
		Server: server,
		Mux:    mux,
		URL:    server.URL,
		Client: client,
	}
}

// Close shuts down the test server
func (ts *TestServer) Close() {
	ts.Server.Close()
}
