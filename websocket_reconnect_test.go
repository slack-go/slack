package slack_test

import (
	"bufio"
	"net"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"testing/synctest"
	"time"

	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/require"

	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slacktest"
)

func TestRTMReconnectKeepsReplacementConnection(t *testing.T) {
	for _, test := range []struct {
		name              string
		failedConnections int32
		goodbye           bool
		wantCause         error
	}{
		{name: "deadman", failedConnections: 1, wantCause: slack.ErrRTMDeadman},
		{name: "repeated deadman", failedConnections: 2, wantCause: slack.ErrRTMDeadman},
		{name: "goodbye", failedConnections: 1, goodbye: true, wantCause: slack.ErrRTMGoodbye},
	} {
		t.Run(test.name, func(t *testing.T) {
			synctest.Test(t, func(t *testing.T) {
				var connections atomic.Int32
				dialer := pipeWebsocketDialer(t, slacktest.Websocket(func(conn *websocket.Conn) {
					failed := connections.Add(1) <= test.failedConnections
					if failed && test.goodbye {
						if err := slacktest.RTMServerSendGoodbye(conn); err != nil {
							return
						}
					}
					for {
						var ping slack.Ping
						if err := conn.ReadJSON(&ping); err != nil {
							return
						}
						// The failed connections stay open but never answer pings.
						if !failed && ping.Type == "ping" {
							if err := conn.WriteJSON(slack.Pong{Type: "pong", ReplyTo: ping.ID, Timestamp: ping.Timestamp}); err != nil {
								return
							}
						}
					}
				}))
				client := rtmConnectTestClient(func(w http.ResponseWriter, r *http.Request) {
					if r.URL.Path != "/api/rtm.connect" {
						t.Errorf("unexpected RTM endpoint: %s", r.URL.Path)
					}
					w.Header().Set("Content-Type", "application/json")
					_, _ = w.Write([]byte(`{"ok":true,"url":"ws://slack.test/ws"}`))
				})
				api := slack.New(testToken, slack.OptionHTTPClient(client))
				rtm := api.NewRTM(slack.RTMOptionDialer(dialer), slack.RTMOptionPingInterval(100*time.Millisecond))
				startRTMTestConnection(t, rtm)

				deadline := time.NewTimer(10 * time.Second)
				defer deadline.Stop()
				var connected, disconnected int32
				pongs := 0
				for pongs < 6 { // Keep the replacement alive beyond a full deadman interval.
					select {
					case event := <-rtm.IncomingEvents:
						switch ev := event.Data.(type) {
						case *slack.ConnectedEvent:
							connected++
							require.LessOrEqual(t, connected, test.failedConnections+1, "healthy replacement was reconnected")
						case *slack.DisconnectedEvent:
							require.False(t, ev.Intentional)
							disconnected++
							require.LessOrEqual(t, disconnected, test.failedConnections, "unexpected extra disconnect: %v", ev.Cause)
							require.ErrorIs(t, ev.Cause, test.wantCause)
						case *slack.LatencyReport:
							pongs++
						}
					case <-deadline.C:
						t.Fatal("timed out waiting for a stable replacement connection")
					}
				}
				require.Equal(t, test.failedConnections+1, connected)
				require.Equal(t, test.failedConnections, disconnected)
				require.Equal(t, connected, connections.Load())
			})
		})
	}
}

type rtmConnectTestClient http.HandlerFunc

func (h rtmConnectTestClient) Do(r *http.Request) (*http.Response, error) {
	if r.Body != nil {
		defer r.Body.Close()
	}
	recorder := httptest.NewRecorder()
	h(recorder, r)
	return recorder.Result(), nil
}

func startRTMTestConnection(t *testing.T, rtm *slack.RTM) {
	t.Helper()
	done := make(chan struct{})
	go func() {
		defer close(done)
		rtm.ManageConnection()
	}()
	t.Cleanup(func() {
		go func() { _ = rtm.Disconnect() }()
		timeout := time.NewTimer(5 * time.Second)
		defer timeout.Stop()
		for {
			select {
			case <-rtm.IncomingEvents:
			case <-done:
				return
			case <-timeout.C:
				t.Error("RTM manager did not stop after Disconnect")
				return
			}
		}
	})
}

// pipeWebsocketDialer runs a real WebSocket handshake over in-memory I/O so
// the virtual clock can advance while both ends wait for frames.
func pipeWebsocketDialer(t *testing.T, serve http.HandlerFunc) *websocket.Dialer {
	t.Helper()
	return &websocket.Dialer{
		HandshakeTimeout: time.Second,
		NetDial: func(_, _ string) (net.Conn, error) {
			client, server := net.Pipe()
			go func() {
				defer server.Close()
				rw := bufio.NewReadWriter(bufio.NewReader(server), bufio.NewWriter(server))
				req, err := http.ReadRequest(rw.Reader)
				if err != nil {
					t.Errorf("read WebSocket upgrade request: %v", err)
					return
				}
				serve(&pipeResponseWriter{
					ResponseRecorder: httptest.NewRecorder(),
					conn:             server,
					rw:               rw,
				}, req)
			}()
			return client, nil
		},
	}
}

type pipeResponseWriter struct {
	*httptest.ResponseRecorder
	conn net.Conn
	rw   *bufio.ReadWriter
}

func (w *pipeResponseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	return w.conn, w.rw, nil
}
