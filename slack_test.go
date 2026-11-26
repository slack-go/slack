package slack

import (
	"log"
	"net/http"
	"net/http/httptest"
	"sync"
)

type httpClientFunc func(*http.Request) (*http.Response, error)

func (f httpClientFunc) Do(r *http.Request) (*http.Response, error) { return f(r) }

// handlerTestClient lets tests use HTTP handlers inside a synctest bubble.
type handlerTestClient http.HandlerFunc

func (h handlerTestClient) Do(r *http.Request) (*http.Response, error) {
	if r.Body != nil {
		defer r.Body.Close()
	}
	recorder := httptest.NewRecorder()
	h(recorder, r)
	return recorder.Result(), nil
}

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
