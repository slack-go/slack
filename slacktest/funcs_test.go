package slacktest

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWebsocketReconnectReceivesNextMessage(t *testing.T) {
	var s *Server
	disconnected := make(chan struct{}, 2)
	s = NewTestServer(func(c Customize) {
		c.Handle("/ws", func(w http.ResponseWriter, r *http.Request) {
			s.wsHandler(w, r)
			disconnected <- struct{}{}
		})
	})
	t.Cleanup(s.Stop)
	s.Start()

	first := dialTestWebsocket(t, s)
	s.SendToWebsocket(`{"type":"hello"}`)
	requireWebsocketMessage(t, first, `{"type":"hello"}`)
	require.NoError(t, first.WriteMessage(websocket.CloseMessage,
		websocket.FormatCloseMessage(websocket.CloseNormalClosure, "")))
	select {
	case <-disconnected:
	case <-time.After(5 * time.Second):
		t.Fatal("server did not close the first connection")
	}

	second := dialTestWebsocket(t, s)
	s.SendToWebsocket(`{"type":"goodbye"}`)
	requireWebsocketMessage(t, second, `{"type":"goodbye"}`)
}

func dialTestWebsocket(t *testing.T, s *Server) *websocket.Conn {
	t.Helper()
	c, _, err := websocket.DefaultDialer.Dial(s.GetWSURL(), nil)
	require.NoError(t, err)
	t.Cleanup(func() { _ = c.Close() })
	return c
}

func requireWebsocketMessage(t *testing.T, c *websocket.Conn, want string) {
	t.Helper()
	require.NoError(t, c.SetReadDeadline(time.Now().Add(5*time.Second)))
	_, got, err := c.ReadMessage()
	require.NoError(t, err)
	assert.Equal(t, want, string(got))
}

func TestGenerateDefaultRTMInfo(t *testing.T) {
	wsurl := "ws://127.0.0.1:5555/ws"
	ctx := context.TODO()
	info := generateRTMInfo(ctx, wsurl)
	assert.Equal(t, wsurl, info.URL)
	assert.True(t, info.Ok)
	assert.Equal(t, defaultBotID, info.User.ID)
	assert.Equal(t, defaultBotName, info.User.Name)
	assert.Equal(t, defaultTeamID, info.Team.ID)
	assert.Equal(t, defaultTeamName, info.Team.Name)
	assert.Equal(t, defaultTeamDomain, info.Team.Domain)
}

func TestCustomDefaultRTMInfo(t *testing.T) {
	wsurl := "ws://127.0.0.1:5555/ws"
	defaultInfo := generateRTMInfo(context.Background(), wsurl)
	ctx := context.TODO()
	ctx = context.WithValue(ctx, ServerBotIDContextKey, "U1234567890")
	ctx = context.WithValue(ctx, ServerBotNameContextKey, "SomeTestBotThing")
	info := generateRTMInfo(ctx, wsurl)
	assert.Equal(t, wsurl, info.URL)
	assert.True(t, info.Ok)
	assert.Equal(t, "U1234567890", info.User.ID)
	assert.Equal(t, "SomeTestBotThing", info.User.Name)
	assert.Equal(t, defaultBotID, defaultInfo.User.ID)
	assert.Equal(t, defaultBotName, defaultInfo.User.Name)
	assert.Equal(t, defaultTeamID, info.Team.ID)
	assert.Equal(t, defaultTeamName, info.Team.Name)
	assert.Equal(t, defaultTeamDomain, info.Team.Domain)
}

func TestGetHubMissingServerAddr(t *testing.T) {
	mc, err := getHubForServer("")
	assert.Nil(t, mc.seen, "seen should be nil")
	assert.Nil(t, mc.sent, "sent should be nil")
	assert.Nil(t, mc.posted, "posted should be nil")
	require.Error(t, err, "should return an error")
	assert.EqualError(t, err, ErrPassedEmptyServerAddr.Error())
}

func TestGetHubNoQueuesForServer(t *testing.T) {
	mc, err := getHubForServer("foo")
	assert.Nil(t, mc.seen, "seen should be nil")
	assert.Nil(t, mc.sent, "sent should be nil")
	assert.Nil(t, mc.posted, "posted should be nil")
	require.Error(t, err, "should return an error")
	assert.EqualError(t, err, ErrNoQueuesRegisteredForServer.Error())
}

func TestUnableToAddToHub(t *testing.T) {
	err := addServerToHub(&Server{}, &messageChannels{})
	require.Error(t, err, "should return and error")
	assert.EqualError(t, err, ErrEmptyServerToHub.Error())
}
