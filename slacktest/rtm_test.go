package slacktest

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/slack-go/slack"
)

func startTestRTM(t *testing.T, s *Server, options ...slack.RTMOption) *slack.RTM {
	t.Helper()
	api := slack.New("ABCDEFG", slack.OptionAPIURL(s.GetAPIURL()))
	rtm := api.NewRTM(options...)
	done := make(chan struct{})
	go func() {
		defer close(done)
		rtm.ManageConnection()
	}()
	t.Cleanup(func() {
		stopped := make(chan struct{})
		go func() {
			_ = rtm.Disconnect()
			<-done
			close(stopped)
		}()
		select {
		case <-stopped:
		case <-time.After(5 * time.Second):
			t.Error("connection manager did not stop")
		}
	})
	return rtm
}

func waitForRTMEvent[E any](t *testing.T, rtm *slack.RTM) E {
	t.Helper()
	timer := time.NewTimer(5 * time.Second)
	defer timer.Stop()
	for {
		select {
		case event := <-rtm.IncomingEvents:
			if data, ok := event.Data.(E); ok {
				return data
			}
		case <-timer.C:
			var want E
			t.Fatalf("timed out waiting for %T", want)
		}
	}
}

func TestRTMInfo(t *testing.T) {
	s := NewTestServer()
	t.Cleanup(s.Stop)
	s.Start()
	rtm := startTestRTM(t, s)
	m := waitForRTMEvent[*slack.ConnectedEvent](t, rtm)
	assert.Equal(t, s.BotID, m.Info.User.ID, "bot id did not match")
	assert.Equal(t, s.BotName, m.Info.User.Name, "bot name did not match")
}

func TestRTMPing(t *testing.T) {
	s := NewTestServer()
	t.Cleanup(s.Stop)
	s.Start()
	rtm := startTestRTM(t, s, slack.RTMOptionPingInterval(50*time.Millisecond))
	m := waitForRTMEvent[*slack.LatencyReport](t, rtm)
	assert.Positive(t, m.Value, "latency report should have a positive value")
}

func TestRTMDirectMessage(t *testing.T) {
	s := NewTestServer()
	t.Cleanup(s.Stop)
	s.Start()
	rtm := startTestRTM(t, s)
	s.SendDirectMessageToBot("some text")
	m := waitForRTMEvent[*slack.MessageEvent](t, rtm)
	assert.Equal(t, defaultNonBotUserID, m.User)
	assert.Equal(t, "D024BE91L", m.Channel)
	assert.Equal(t, "some text", m.Text)
}

func TestRTMChannelMessage(t *testing.T) {
	s := NewTestServer()
	t.Cleanup(s.Stop)
	s.Start()
	rtm := startTestRTM(t, s)
	s.SendMessageToChannel("#foochan", "some text")
	m := waitForRTMEvent[*slack.MessageEvent](t, rtm)
	assert.Equal(t, "#foochan", m.Channel)
	assert.Equal(t, "some text", m.Text)
}
