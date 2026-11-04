package slack_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slacktest"
)

const (
	testMessage = "test message"
	testToken   = "TEST_TOKEN"
)

func TestRTMBeforeEvents(t *testing.T) {
	// Set up the test server.
	testServer := slacktest.NewTestServer()
	go testServer.Start()

	// Setup and start the RTM.
	api := slack.New(testToken, slack.OptionAPIURL(testServer.GetAPIURL()))
	rtm := api.NewRTM()

	done := make(chan struct{})
	go func() {
		for msg := range rtm.IncomingEvents {
			switch ev := msg.Data.(type) {
			case *slack.DisconnectedEvent:
				if ev.Intentional {
					close(done)
					return
				}
			default:
				// t.Logf("Discarded event of type '%s' with content '%#v'", msg.Type, ev)
			}
		}
	}()
	go rtm.Disconnect()
	go rtm.ManageConnection()
	select {
	case <-done:
	case <-time.After(5 * time.Second):
		t.Error("timed out waiting for disconnect")
		t.Fail()
	}
}

func TestRTMGoodbye(t *testing.T) {
	// Set up the test server.
	testServer := slacktest.NewTestServer(
		func(c slacktest.Customize) {
			c.Handle("/ws", slacktest.Websocket(func(conn *websocket.Conn) {
				if err := slacktest.RTMServerSendGoodbye(conn); err != nil {
					log.Println("failed to send goodbye", err)
				}
			}))
		},
	)
	go testServer.Start()

	// Setup and start the RTM.
	api := slack.New(
		testToken,
		slack.OptionAPIURL(testServer.GetAPIURL()),
	)

	rtm := api.NewRTM(
		slack.RTMOptionPingInterval(100 * time.Millisecond),
	)

	done := make(chan struct{})
	go rtm.ManageConnection()
	connected := 0
	disconnected := 0
	func() {
		for msg := range rtm.IncomingEvents {
			switch ev := msg.Data.(type) {
			case *slack.ConnectedEvent:
				connected += 1
				if connected > 5 {
					rtm.Disconnect()
				}
			case *slack.DisconnectedEvent:
				// t.Log("disconnect event received", ev.Intentional, ev.Cause)
				if ev.Intentional {
					close(done)
					return
				}
				disconnected += 1
			default:
				// t.Logf("Discarded event of type '%s' with content '%#v'", msg.Type, ev)
			}
		}
	}()

	select {
	case <-done:
		// magic numbers from empirical testing.
		assert.LessOrEqual(t, connected, 7)
		assert.LessOrEqual(t, disconnected, 12)
	case <-time.After(5 * time.Second):
		t.Error("timed out waiting for disconnect")
		t.Fail()
	}
}

func TestRTMDeadConnection(t *testing.T) {
	// Set up the test server.
	testServer := slacktest.NewTestServer(
		func(c slacktest.Customize) {
			c.Handle("/ws", slacktest.Websocket(func(conn *websocket.Conn) {
				// closes immediately
			}))
		},
	)
	go testServer.Start()

	// Setup and start the RTM.
	api := slack.New(
		testToken,
		slack.OptionAPIURL(testServer.GetAPIURL()),
	)

	rtm := api.NewRTM(
		slack.RTMOptionPingInterval(100 * time.Millisecond),
	)

	go rtm.ManageConnection()
	done := make(chan struct{})
	connected := 0
	disconnected := 0
	func() {
		for msg := range rtm.IncomingEvents {
			switch ev := msg.Data.(type) {
			case *slack.ConnectedEvent:
				connected += 1
				if connected > 5 {
					rtm.Disconnect()
				}
			case *slack.DisconnectedEvent:
				// t.Log("disconnect event received", ev.Intentional, ev.Cause)
				if ev.Intentional {
					close(done)
					return
				}
				disconnected += 1
			default:
				// t.Logf("Discarded event of type '%s' with content '%#v'", msg.Type, ev)
			}
		}
	}()

	select {
	case <-done:
		// magic numbers from empirical testing.
		assert.LessOrEqual(t, connected, 7)
		assert.LessOrEqual(t, disconnected, 7)
	case <-time.After(5 * time.Second):
		t.Error("timed out waiting for disconnect")
		t.Fail()
	}
}

func TestRTMDisconnect(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/rtm.connect" {
			t.Errorf("unexpected endpoint: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = fmt.Fprintln(w, `{"ok":false,"error":"invalid_auth"}`)
	}))
	t.Cleanup(server.Close)

	api := slack.New(testToken, slack.OptionAPIURL(server.URL+"/"))
	rtm := api.NewRTM()
	done := make(chan struct{})
	go func() {
		defer close(done)
		rtm.ManageConnection()
	}()
	t.Cleanup(func() {
		_ = rtm.Disconnect()
		select {
		case <-done:
		case <-time.After(5 * time.Second):
			t.Error("connection manager did not stop during cleanup")
		}
	})

	timer := time.NewTimer(5 * time.Second)
	defer timer.Stop()
	for _, want := range []any{&slack.ConnectingEvent{}, &slack.InvalidAuthEvent{}} {
		select {
		case event := <-rtm.IncomingEvents:
			require.IsType(t, want, event.Data)
		case <-timer.C:
			t.Fatalf("timed out waiting for %T", want)
		}
	}

	select {
	case <-done:
	case <-timer.C:
		t.Fatal("connection manager did not stop after invalid auth")
	}
}

func TestRTMConnectRateLimit(t *testing.T) {
	// Set up the test server.
	testServer := slacktest.NewTestServer(
		func(c slacktest.Customize) {
			c.Handle("/rtm.connect", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Add("Retry-After", "1")
				w.WriteHeader(http.StatusTooManyRequests)
			}))
		},
	)
	go testServer.Start()

	// Setup and start the RTM.
	api := slack.New(testToken, slack.OptionAPIURL(testServer.GetAPIURL()))
	rtm := api.NewRTM()
	go rtm.ManageConnection()

	// Observe incoming failures
	connectionFailure := make(chan *slack.ConnectionErrorEvent)
	go func() {
		for msg := range rtm.IncomingEvents {
			switch ev := msg.Data.(type) {
			case *slack.ConnectingEvent:
			case *slack.ConnectionErrorEvent:
				connectionFailure <- ev
				if ev.Attempt > 5 {
					rtm.Disconnect()
				}
			case *slack.DisconnectedEvent:
				if ev.Intentional {
					close(connectionFailure)
					return
				}
			default:
				t.Logf("Discarded event of type '%s' with content '%#v'", msg.Type, ev)
			}
		}
	}()

	previous := time.Duration(0)
	for ev := range connectionFailure {
		assert.LessOrEqual(t, previous, ev.Backoff, "backoff should increase during rate limits: %v <= %v", previous, ev.Backoff)
		previous = ev.Backoff
	}
	testServer.Stop()
}

func TestRTMSingleConnect(t *testing.T) {
	// Set up the test server.
	testServer := slacktest.NewTestServer()
	go testServer.Start()

	// Setup and start the RTM.
	api := slack.New(testToken, slack.OptionAPIURL(testServer.GetAPIURL()))
	rtm := api.NewRTM()
	go rtm.ManageConnection()

	// Observe incoming messages.
	done := make(chan struct{})
	connectingReceived := false
	connectedReceived := false
	testMessageReceived := false
	go func() {
		for msg := range rtm.IncomingEvents {
			switch ev := msg.Data.(type) {
			case *slack.ConnectingEvent:
				if connectingReceived {
					t.Error("Received multiple connecting events.")
					t.Fail()
				}
				connectingReceived = true
			case *slack.ConnectedEvent:
				if connectedReceived {
					t.Error("Received multiple connected events.")
					t.Fail()
				}
				connectedReceived = true
			case *slack.MessageEvent:
				if ev.Text == testMessage {
					testMessageReceived = true
					rtm.Disconnect()
				}
				t.Logf("Discarding message with content %+v", ev)
			case *slack.DisconnectedEvent:
				if ev.Intentional {
					done <- struct{}{}
					return
				}
			default:
				t.Logf("Discarded event of type '%s' with content '%#v'", msg.Type, ev)
			}
		}
	}()

	// Send a message and sleep for some time to make sure the message can be processed client-side.
	testServer.SendDirectMessageToBot(testMessage)
	<-done
	testServer.Stop()

	// Verify that all expected events have been received by the RTM client.
	assert.True(t, connectingReceived, "Should have received a connecting event from the RTM instance.")
	assert.True(t, connectedReceived, "Should have received a connected event from the RTM instance.")
	assert.True(t, testMessageReceived, "Should have received a test message from the server.")
}

func TestRTMUnmappedError(t *testing.T) {
	const unmappedEventName = "some_unknown_event"
	// Set up the test server.
	testServer := slacktest.NewTestServer()
	go testServer.Start()

	// Setup and start the RTM.
	api := slack.New(testToken, slack.OptionAPIURL(testServer.GetAPIURL()))
	rtm := api.NewRTM()
	go rtm.ManageConnection()

	// Observe incoming messages.
	done := make(chan struct{})
	var gotUnmarshallingError *slack.UnmarshallingErrorEvent
	go func() {
		for msg := range rtm.IncomingEvents {
			switch ev := msg.Data.(type) {
			case *slack.UnmarshallingErrorEvent:
				gotUnmarshallingError = ev
				rtm.Disconnect()
			case *slack.DisconnectedEvent:
				if ev.Intentional {
					done <- struct{}{}
					return
				}
			default:
				t.Logf("Discarded event of type '%s' with content '%#v'", msg.Type, ev)
			}
		}
	}()

	// Send a message and sleep for some time to make sure the message can be processed client-side.
	testServer.SendToWebsocket(fixSlackMessage(t, unmappedEventName))
	<-done
	testServer.Stop()

	// Verify that we got the expected error with details
	unmappedErr, ok := errors.AsType[*slack.UnmappedError](gotUnmarshallingError.ErrorObj)
	require.True(t, ok)
	assert.Equal(t, unmappedEventName, unmappedErr.EventType)
}

func fixSlackMessage(t *testing.T, eType string) string {
	t.Helper()

	m := slack.Message{
		Msg: slack.Msg{
			Type:      eType,
			Text:      "Fixture Slack message",
			Timestamp: fmt.Sprintf("%d", time.Now().Unix()),
		},
	}
	msg, err := json.Marshal(m)
	require.NoError(t, err)

	return string(msg)
}
