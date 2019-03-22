package slack_test

import (
	"testing"
	"time"

	slacktest "github.com/lusis/slack-test"
	"github.com/nlopes/slack"
	"github.com/stretchr/testify/assert"
)

const (
	testMessage = "test message"
	testToken   = "TEST_TOKEN"
)

func TestRTMDisconnect(t *testing.T) {
	// actually connect to slack here w/ an invalid token
	api := slack.New(testToken)
	rtm := api.NewRTM()
	go rtm.ManageConnection()

	// Observe incoming messages.
	done := make(chan struct{})
	connectingReceived := false
	disconnectedReceived := false

	go func() {
		for msg := range rtm.IncomingEvents {
			switch ev := msg.Data.(type) {
			case *slack.InvalidAuthEvent:
				t.Log("invalid auth event received")
				disconnectedReceived = true
				close(done)
			case *slack.ConnectingEvent:
				connectingReceived = true
			case *slack.ConnectedEvent:
				t.Error("received connected events on an invalid connection")
				t.Fail()
			default:
				t.Logf("discarded event of type '%s' with content '%#v'", msg.Type, ev)
			}
		}
	}()

	select {
	case <-done:
	case <-time.After(5 * time.Second):
		t.Error("timed out waiting for disconnect")
		t.Fail()
	}

	// Verify that all expected events have been received by the RTM client.
	assert.True(t, connectingReceived, "Should have received a connecting event from the RTM instance.")
	assert.True(t, disconnectedReceived, "Should have received a disconnected event from the RTM instance.")
}

func TestRTMSingleConnect(t *testing.T) {
	// Set up the test server.
	testServer := slacktest.NewTestServer()
	go testServer.Start()

	// Setup and start the RTM.
	slack.APIURL = testServer.GetAPIURL()
	api := slack.New(testToken)
	rtm := api.NewRTM(slack.RTMOptionUseStart(true))
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
