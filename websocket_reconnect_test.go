package slack_test

import (
	"net/http"
	"sync/atomic"
	"testing"
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
			var connections atomic.Int32
			server := slacktest.NewTestServer(func(c slacktest.Customize) {
				c.Handle("/ws", slacktest.Websocket(func(conn *websocket.Conn) {
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
			})
			server.Start()
			t.Cleanup(server.Stop)

			api := slack.New("test-token", slack.OptionAPIURL(server.GetAPIURL()), slack.OptionHTTPClient(&http.Client{Timeout: time.Second}))
			rtm := api.NewRTM(slack.RTMOptionPingInterval(100 * time.Millisecond))
			done := make(chan struct{})
			go func() {
				rtm.ManageConnection()
				close(done)
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
	}
}
