package slacktest

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/slack-go/slack"
)

func TestDefaultNewServer(t *testing.T) {
	s := NewTestServer()
	assert.Equal(t, defaultBotID, s.BotID)
	assert.Equal(t, defaultBotName, s.BotName)
	assert.NotEmpty(t, s.ServerAddr)
	s.Stop()
}

func TestCustomNewServer(t *testing.T) {
	s := NewTestServer()
	t.Cleanup(s.Stop)
	s.SetBotName("BobsBot")
	assert.Equal(t, "BobsBot", s.BotName)
}

func TestServerSendMessageToChannel(t *testing.T) {
	s := NewTestServer()
	t.Cleanup(s.Stop)
	s.Start()
	s.SendMessageToChannel("C123456789", "some text")
	assert.True(t, s.SawOutgoingMessage("some text"))
}

func TestServerSendMessageToBot(t *testing.T) {
	s := NewTestServer()
	t.Cleanup(s.Stop)
	s.Start()
	s.SendMessageToBot("C123456789", "some text")
	expectedMsg := fmt.Sprintf("<@%s> %s", s.BotID, "some text")
	assert.True(t, s.SawOutgoingMessage(expectedMsg))
}

func TestBotDirectMessageBotHandler(t *testing.T) {
	s := NewTestServer()
	t.Cleanup(s.Stop)
	s.Start()
	s.SendDirectMessageToBot("some text")
	expectedMsg := "some text"
	assert.True(t, s.SawOutgoingMessage(expectedMsg))
}

func TestGetSeenOutboundMessages(t *testing.T) {
	s := NewTestServer()
	t.Cleanup(s.Stop)
	s.Start()

	s.SendMessageToChannel("foo", "should see this message")

	seenOutbound := s.GetSeenOutboundMessages()
	assert.Len(t, seenOutbound, 1)

	hadMessage := false
	for _, msg := range seenOutbound {
		var m = slack.Message{}
		jerr := json.Unmarshal([]byte(msg), &m)
		require.NoError(t, jerr, "messages should decode as slack.Message")
		if m.Text == "should see this message" {
			hadMessage = true
			break
		}
	}
	assert.True(t, hadMessage, "did not see my sent message")
}

func TestGetSeenInboundMessages(t *testing.T) {
	s := NewTestServer()
	t.Cleanup(s.Stop)
	s.Start()

	rtm := startTestRTM(t, s)
	rtm.SendMessage(&slack.OutgoingMessage{
		Channel: "foo",
		Text:    "should see this inbound message",
	})
	select {
	case <-s.SeenFeed:
	case <-time.After(5 * time.Second):
		t.Fatal("did not observe the inbound message")
	}
	seenInbound := s.GetSeenInboundMessages()
	assert.NotEmpty(t, seenInbound)
	hadMessage := false
	for _, msg := range seenInbound {
		var m = slack.Message{}
		jerr := json.Unmarshal([]byte(msg), &m)
		require.NoError(t, jerr, "messages should decode as slack.Message")
		if m.Text == "should see this inbound message" {
			hadMessage = true
			break
		}
	}
	assert.True(t, hadMessage, "did not see my sent message")
	assert.True(t, s.SawMessage("should see this inbound message"))
}

func TestSendChannelInvite(t *testing.T) {
	s := NewTestServer()
	t.Cleanup(s.Stop)
	s.Start()
	rtm := startTestRTM(t, s)
	s.SendBotChannelInvite()
	m := waitForRTMEvent[*slack.ChannelJoinedEvent](t, rtm)
	assert.Equal(t, "C024BE92L", m.Channel.ID, "channel id should match")
	assert.Equal(t, "Fun times", m.Channel.Topic.Value, "topic should match")
}

func TestSendGroupInvite(t *testing.T) {
	s := NewTestServer()
	t.Cleanup(s.Stop)
	s.Start()
	rtm := startTestRTM(t, s)
	s.SendBotGroupInvite()
	m := waitForRTMEvent[*slack.GroupJoinedEvent](t, rtm)
	assert.Equal(t, "G024BE91L", m.Channel.ID, "channel id should match")
	assert.Equal(t, "Secret plans on hold", m.Channel.Topic.Value, "topic should match")
}

func TestServerSawMessage(t *testing.T) {
	s := NewTestServer()
	t.Cleanup(s.Stop)
	s.Start()
	assert.False(t, s.SawMessage("foo"), "should not have seen any message")
}

func TestServerSawOutgoingMessage(t *testing.T) {
	s := NewTestServer()
	t.Cleanup(s.Stop)
	s.Start()
	assert.False(t, s.SawOutgoingMessage("foo"), "should not have seen any message")
}
