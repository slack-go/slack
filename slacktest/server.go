package slacktest

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/slack-go/slack"
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

// Customize the server's responses.
type Customize interface {
	Handle(pattern string, handler http.HandlerFunc)
}

type binder func(Customize)

// NewTestServer returns a slacktest.Server ready to be started
func NewTestServer(custom ...binder) *Server {
	serverChans := newMessageChannels()

	channels := &serverChannels{}
	groups := &serverGroups{}
	s := &Server{
		registered:           map[string]struct{}{},
		mux:                  http.NewServeMux(),
		seenInboundMessages:  &messageCollection{},
		seenOutboundMessages: &messageCollection{},
	}

	for _, c := range custom {
		c(s)
	}

	s.Handle("/conversations.info", s.conversationsInfoHandler)
	s.Handle("/ws", s.wsHandler)
	s.Handle("/rtm.start", rtmStartHandler)
	s.Handle("/rtm.connect", RTMConnectHandler)
	s.Handle("/chat.postMessage", s.postMessageHandler)
	s.Handle("/conversations.create", createConversationHandler)
	s.Handle("/conversations.setTopic", setConversationTopicHandler)
	s.Handle("/conversations.setPurpose", setConversationPurposeHandler)
	s.Handle("/conversations.rename", renameConversationHandler)
	s.Handle("/conversations.invite", inviteConversationHandler)
	s.Handle("/users.info", usersInfoHandler)
	s.Handle("/users.lookupByEmail", usersInfoHandler)
	s.Handle("/bots.info", botsInfoHandler)
	s.Handle("/auth.test", authTestHandler)
	s.Handle("/reactions.add", reactionAddHandler)

	httpserver := httptest.NewUnstartedServer(s.mux)
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

// Handle allow for customizing endpoints
func (sts *Server) Handle(pattern string, handler http.HandlerFunc) {
	if _, found := sts.registered[pattern]; found {
		// log.Printf("route already registered: %s\n", pattern)
		return
	}

	sts.registered[pattern] = struct{}{}
	sts.mux.Handle(pattern, contextHandler(sts, handler))
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

// GetSeenInboundMessages returns all messages seen via websocket excluding pings
func (sts *Server) GetSeenInboundMessages() []string {
	sts.seenInboundMessages.RLock()
	m := sts.seenInboundMessages.messages
	sts.seenInboundMessages.RUnlock()
	return m
}

// GetSeenOutboundMessages returns all messages seen via websocket excluding pings
func (sts *Server) GetSeenOutboundMessages() []string {
	sts.seenOutboundMessages.RLock()
	m := sts.seenOutboundMessages.messages
	sts.seenOutboundMessages.RUnlock()
	return m
}

// SawOutgoingMessage checks if a message was sent to connected websocket clients
func (sts *Server) SawOutgoingMessage(msg string) bool {
	sts.seenOutboundMessages.RLock()
	defer sts.seenOutboundMessages.RUnlock()
	for _, m := range sts.seenOutboundMessages.messages {
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
	sts.seenInboundMessages.RLock()
	defer sts.seenInboundMessages.RUnlock()
	for _, m := range sts.seenInboundMessages.messages {
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

// GetAPIURL returns the api url you can pass to slack.SLACK_API
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
	go sts.queueForWebsocket(string(j), sts.ServerAddr)
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
	go sts.queueForWebsocket(string(j), sts.ServerAddr)
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
	go sts.queueForWebsocket(stringMsg, sts.ServerAddr)
}

// SendToWebsocket send `s` as is to connected clients.
// This is useful for sending your own custom json to the websocket
func (sts *Server) SendToWebsocket(s string) {
	go sts.queueForWebsocket(s, sts.ServerAddr)
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
func (sts *Server) GetTestRTMInstance() *slack.RTM {
	api := slack.New("ABCEFG", slack.OptionAPIURL(sts.GetAPIURL()))
	rtm := api.NewRTM()
	return rtm
}
