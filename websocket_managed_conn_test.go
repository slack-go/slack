package slack_test

import (
	"testing"

	slacktest "github.com/lusis/slack-test"
	"github.com/nlopes/slack"
	"github.com/stretchr/testify/assert"
)

const (
	testMessage = "test message"
	testToken   = "TEST_TOKEN"
)

func TestRTMSingleConnect(t *testing.T) {
	// Set up the test server.
	testServer := slacktest.NewTestServer()
	go testServer.Start()

	// Setup and start the RTM.
	slack.SLACK_API = testServer.GetAPIURL()
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
					done <- struct{}{}
					return
				}
				t.Logf("Discarding message with content %+v", ev)
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
