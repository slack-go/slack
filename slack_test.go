package slack

import (
	"net/http"
	"net/http/httptest"
)

const validToken = "testing-token"

var serverAddr string

type testServer struct {
	server    *httptest.Server
	mux       *http.ServeMux
	wasCalled bool
}

func (t *testServer) Close() {
	if !t.wasCalled {
		panic("close called on test server, but nothing was registered")
	}
	t.server.Close()
}

func (t *testServer) RegisterHandler(pattern string, handler func(http.ResponseWriter, *http.Request)) {
	t.wasCalled = true
	t.mux.HandleFunc(pattern, handler)
}

func startServer() *testServer {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	serverAddr = server.Listener.Addr().String()

	return &testServer{
		server: server,
		mux:    mux,
	}
}
