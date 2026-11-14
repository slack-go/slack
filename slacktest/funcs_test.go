package slacktest

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHandlePendingMessagesStopsOnWriteError(t *testing.T) {
	closedConnections := make(chan *websocket.Conn, 1)
	server := httptest.NewServer(http.HandlerFunc(Websocket(func(conn *websocket.Conn) {
		_ = conn.Close()
		closedConnections <- conn
	})))
	defer server.Close()
	client, _, err := websocket.DefaultDialer.Dial("ws"+strings.TrimPrefix(server.URL, "http"), nil)
	require.NoError(t, err)
	defer client.Close()

	channels := newMessageChannels()
	require.NoError(t, addServerToHub(&Server{ServerAddr: t.Name()}, channels))
	defer func() {
		masterHub.Lock()
		delete(masterHub.serverChannels, t.Name())
		masterHub.Unlock()
	}()
	conn := <-closedConnections
	done := make(chan struct{})
	go func() {
		defer close(done)
		handlePendingMessages(conn, t.Name())
	}()
	defer func() {
		close(channels.sent)
		<-done
	}()

	channels.sent <- "first message"
	select {
	case channels.sent <- "second message":
		t.Fatal("failed writer consumed another message after its write error")
	case <-done:
	}
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
