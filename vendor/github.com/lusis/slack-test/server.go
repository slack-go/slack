package slacktest

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"time"

	slack "github.com/nlopes/slack"
)

func newMessageChannels() *messageChannels {
	sent := make(chan (string))
	seen := make(chan (string))
	mc := messageChannels{
		seen: seen,
		sent: sent,
	}
	return &mc
}

// NewTestServer returns a slacktest.Server ready to be started
func NewTestServer() *Server {
	serverChans := newMessageChannels()
	seenInboundMessages = &messageCollection{}
	seenOutboundMessages = &messageCollection{}
	channels := &serverChannels{}
	groups := &serverGroups{}
	s := &Server{}
	mux := http.NewServeMux()
	mux.Handle("/ws", contextHandler(s, wsHandler))
	mux.Handle("/rtm.start", contextHandler(s, rtmStartHandler))
	mux.Handle("/chat.postMessage", contextHandler(s, postMessageHandler))
	mux.Handle("/channels.list", contextHandler(s, listChannelsHandler))
	mux.Handle("/groups.list", contextHandler(s, listGroupsHandler))
	mux.Handle("/users.info", contextHandler(s, usersInfoHandler))
	mux.Handle("/bots.info", contextHandler(s, botsInfoHandler))
	httpserver := httptest.NewUnstartedServer(mux)
	addr := httpserver.Listener.Addr().String()

	s.ServerAddr = addr
	s.server = httpserver
	s.BotName = defaultBotName
	s.BotID = defaultBotID
	s.SeenFeed = serverChans.seen
	s.channels = channels
	s.groups = groups
	addErr := addServerToHub(s, serverChans)
	if addErr != nil {
		log.Printf("Unable to add server to hub: %s", addErr.Error())
	}
	return s
}

// GetChannels returns all the fake channels registered
func (sts *Server) GetChannels() []slack.Channel {
	sts.channels.RLock()
	defer sts.channels.RUnlock()
	return sts.channels.channels
}

// GetGroups returns all the fake groups registered
func (sts *Server) GetGroups() []slack.Group {
	return sts.groups.channels
}

/*
// These are placeholders for now
// AddChannel adds a new fake channel
func (sts *Server) AddChannel(c slack.Channel) {
	sts.channels.Lock()
	sts.channels.channels = append(sts.channels.channels, c)
	sts.channels.Unlock()
}

// AddGroup adds a new fake group
func (sts *Server) AddGroup(c slack.Group) {
	sts.groups.Lock()
	sts.groups.channels = append(sts.groups.channels, c)
	sts.groups.Unlock()
}
*/

// GetSeenInboundMessages returns all messages seen via websocket excluding pings
func (sts *Server) GetSeenInboundMessages() []string {
	seenInboundMessages.RLock()
	m := seenInboundMessages.messages
	seenInboundMessages.RUnlock()
	return m
}

// GetSeenOutboundMessages returns all messages seen via websocket excluding pings
func (sts *Server) GetSeenOutboundMessages() []string {
	seenOutboundMessages.RLock()
	m := seenOutboundMessages.messages
	seenOutboundMessages.RUnlock()
	return m
}

// SawOutgoingMessage checks if a message was sent to connected websocket clients
func (sts *Server) SawOutgoingMessage(msg string) bool {
	seenOutboundMessages.RLock()
	defer seenOutboundMessages.RUnlock()
	for _, m := range seenOutboundMessages.messages {
		evt := &slack.MessageEvent{}
		jErr := json.Unmarshal([]byte(m), evt)
		if jErr != nil {
			continue
		}
		if evt.Text == msg {
			return true
		}
	}
	return false
}

// SawMessage checks if an incoming message was seen
func (sts *Server) SawMessage(msg string) bool {
	seenInboundMessages.RLock()
	defer seenInboundMessages.RUnlock()
	for _, m := range seenInboundMessages.messages {
		evt := &slack.MessageEvent{}
		jErr := json.Unmarshal([]byte(m), evt)
		if jErr != nil {
			// This event isn't a message event so we'll skip it
			continue
		}
		if evt.Text == msg {
			return true
		}
	}
	return false
}

// GetAPIURL returns the api url you can pass to slack.APIURL
func (sts *Server) GetAPIURL() string {
	return "http://" + sts.ServerAddr + "/"
}

// GetWSURL returns the websocket url
func (sts *Server) GetWSURL() string {
	return "ws://" + sts.ServerAddr + "/ws"
}

// Stop stops the test server
func (sts *Server) Stop() {
	sts.server.Close()
}

// Start starts the test server
func (sts *Server) Start() {
	log.Print("starting server")
	sts.server.Start()
}

// SendMessageToBot sends a message addressed to the Bot
func (sts *Server) SendMessageToBot(channel, msg string) {
	m := slack.Message{}
	m.Type = slack.TYPE_MESSAGE
	m.Channel = channel
	m.User = defaultNonBotUserID
	m.Text = fmt.Sprintf("<@%s> %s", sts.BotID, msg)
	m.Timestamp = fmt.Sprintf("%d", time.Now().Unix())
	j, jErr := json.Marshal(m)
	if jErr != nil {
		log.Printf("Unable to marshal message for bot: %s", jErr.Error())
		return
	}
	go queueForWebsocket(string(j), sts.ServerAddr)
}

// SendDirectMessageToBot sends a direct message to the bot
func (sts *Server) SendDirectMessageToBot(msg string) {
	m := slack.Message{}
	m.Type = slack.TYPE_MESSAGE
	m.Channel = "D024BE91L"
	m.User = defaultNonBotUserID
	m.Text = msg
	m.Timestamp = fmt.Sprintf("%d", time.Now().Unix())
	j, jErr := json.Marshal(m)
	if jErr != nil {
		log.Printf("Unable to marshal private message for bot: %s", jErr.Error())
		return
	}
	go queueForWebsocket(string(j), sts.ServerAddr)
}

// SendMessageToChannel sends a message to a channel
func (sts *Server) SendMessageToChannel(channel, msg string) {
	m := slack.Message{}
	m.Type = slack.TYPE_MESSAGE
	m.Channel = channel
	m.Text = msg
	m.User = defaultNonBotUserID
	m.Timestamp = fmt.Sprintf("%d", time.Now().Unix())
	j, jErr := json.Marshal(m)
	if jErr != nil {
		log.Printf("Unable to marshal message for channel: %s", jErr.Error())
		return
	}
	stringMsg := string(j)
	go queueForWebsocket(stringMsg, sts.ServerAddr)
}

// SendToWebsocket send `s` as is to connected clients.
// This is useful for sending your own custom json to the websocket
func (sts *Server) SendToWebsocket(s string) {
	go queueForWebsocket(s, sts.ServerAddr)
}

// SetBotName sets a custom botname
func (sts *Server) SetBotName(b string) {
	sts.BotName = b
}

// SendBotChannelInvite invites the bot to a channel
func (sts *Server) SendBotChannelInvite() {
	joinMsg := `
	{
			"type":"channel_joined",
			"channel":
					{
							"id": "C024BE92L",
							"name": "bot-playground",
							"is_channel": true,
							"created": 1360782804,
							"creator": "W012A3CDE",
							"is_archived": false,
							"is_general": true,
							"members": [
									"W012A3CDE"
							],
							"topic": {
									"value": "Fun times",
									"creator": "W012A3CDE",
									"last_set": 1360782804
							},
							"purpose": {
									"value": "This channel is for fun",
									"creator": "W012A3CDE",
									"last_set": 1360782804
							},
							"is_member": true
					}
	}`
	sts.SendToWebsocket(joinMsg)
}

// SendBotGroupInvite invites the bot to a channel
func (sts *Server) SendBotGroupInvite() {
	joinMsg := `
	{
			"type":"group_joined",
			"channel":
			{
				"id": "G024BE91L",
				"name": "secretplans",
				"is_group": true,
				"created": 1360782804,
				"creator": "W012A3CDE",
				"is_archived": false,
				"members": [
					"W012A3CDE"
				],
				"topic": {
					"value": "Secret plans on hold",
					"creator": "W012A3CDE",
					"last_set": 1360782804
				},
				"purpose": {
					"value": "Discuss secret plans that no-one else should know",
					"creator": "W012A3CDE",
					"last_set": 1360782804
				}
			}
	}`
	sts.SendToWebsocket(joinMsg)
}

// GetTestRTMInstance will give you an RTM instance in the context of the current fake server
func (sts *Server) GetTestRTMInstance() (*slack.Client, *slack.RTM) {
	slack.APIURL = sts.GetAPIURL()
	api := slack.New("ABCEFG")
	rtm := api.NewRTM()
	return api, rtm
}
