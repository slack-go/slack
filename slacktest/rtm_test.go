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
	if testing.Short() {
		t.Skip("skipping timered test")
	}
	maxWait := 45 * time.Second
	s := NewTestServer()
	go s.Start()

	api := slack.New("ABCDEFG", slack.OptionAPIURL(s.GetAPIURL()))
	rtm := api.NewRTM()
	go rtm.ManageConnection()
	messageChan := make(chan (*slack.LatencyReport), 1)
	go func() {
		for msg := range rtm.IncomingEvents {
			if ev, ok := msg.Data.(*slack.LatencyReport); ok {
				messageChan <- ev
			}
		}
	}()
	select {
	case m := <-messageChan:
		assert.NotEmpty(t, m.Value, "latency report should value a value")
		assert.Positive(t, m.Value, "latency report should be greater than 0")
		break
	case <-time.After(maxWait):
		assert.FailNow(t, "did not get latency report in time")

	}
}

func TestRTMDirectMessage(t *testing.T) {
	maxWait := 5 * time.Second
	s := NewTestServer()
	go s.Start()

	api := slack.New("ABCDEFG", slack.OptionAPIURL(s.GetAPIURL()))
	rtm := api.NewRTM()
	go rtm.ManageConnection()
	messageChan := make(chan (*slack.MessageEvent), 1)
	go func() {
		for msg := range rtm.IncomingEvents {
			if ev, ok := msg.Data.(*slack.MessageEvent); ok {
				messageChan <- ev
			}
		}
	}()
	s.SendDirectMessageToBot("some text")
	select {
	case m := <-messageChan:
		assert.Equal(t, defaultNonBotUserID, m.User)
		assert.Equal(t, "D024BE91L", m.Channel)
		assert.Equal(t, "some text", m.Text)
		break
	case <-time.After(maxWait):
		assert.FailNow(t, "did not get direct message in time")
	}
}

func TestRTMChannelMessage(t *testing.T) {
	maxWait := 5 * time.Second
	s := NewTestServer()
	go s.Start()

	api := slack.New("ABCDEFG", slack.OptionAPIURL(s.GetAPIURL()))
	rtm := api.NewRTM()
	go rtm.ManageConnection()
	messageChan := make(chan (*slack.MessageEvent), 1)
	go func() {
		for msg := range rtm.IncomingEvents {
			if ev, ok := msg.Data.(*slack.MessageEvent); ok {
				messageChan <- ev
			}
		}
	}()
	s.SendMessageToChannel("#foochan", "some text")
	select {
	case m := <-messageChan:
		assert.Equal(t, "#foochan", m.Channel)
		assert.Equal(t, "some text", m.Text)
		break
	case <-time.After(maxWait):
		assert.FailNow(t, "did not get channel message in time")
	}

}
