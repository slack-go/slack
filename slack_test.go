package slack

import (
	"log"
	"net/http"
	"net/http/httptest"
	"sync"
)

type httpClientFunc func(*http.Request) (*http.Response, error)

func (f httpClientFunc) Do(r *http.Request) (*http.Response, error) { return f(r) }

const (
	validToken = "testing-token"
)

var (
	serverAddr string
	once       sync.Once
)

func startServer() {
	server := httptest.NewServer(nil)
	serverAddr = server.Listener.Addr().String()
	log.Print("Test WebSocket server listening on ", serverAddr)
}
