package slacktest

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"

	websocket "github.com/gorilla/websocket"

	slack "github.com/slack-go/slack"
)

func contextHandler(server *Server, next http.HandlerFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), ServerURLContextKey, server.GetAPIURL())
		ctx = context.WithValue(ctx, ServerWSContextKey, server.GetWSURL())
		ctx = context.WithValue(ctx, ServerBotNameContextKey, server.BotName)
		ctx = context.WithValue(ctx, ServerBotChannelsContextKey, server.GetChannels())
		ctx = context.WithValue(ctx, ServerBotGroupsContextKey, server.GetGroups())
		ctx = context.WithValue(ctx, ServerBotHubNameContextKey, server.ServerAddr)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// handle auth.test
func authTestHandler(w http.ResponseWriter, _ *http.Request) {
	_, _ = w.Write([]byte(defaultAuthTestJSON))
}

func usersInfoHandler(w http.ResponseWriter, _ *http.Request) {
	_, _ = w.Write([]byte(defaultUsersInfoJSON))
}

func botsInfoHandler(w http.ResponseWriter, r *http.Request) {
	_, _ = w.Write([]byte(defaultBotInfoJSON(r.Context())))
}

type GroupConversationResponse struct {
	Ok      bool                    `json:"ok"`
	Channel slack.GroupConversation `json:"channel"`
}

func (sts *Server) conversationsInfoHandler(w http.ResponseWriter, r *http.Request) {
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		msg := fmt.Sprintf("error reading body: %s", err.Error())
		log.Printf(msg)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}
	values, vErr := url.ParseQuery(string(data))
	if vErr != nil {
		msg := fmt.Sprintf("Unable to decode query params: %s", vErr.Error())
		log.Printf(msg)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}
	ch := values.Get("channel")

	response := GroupConversationResponse{
		Ok: true,
		Channel: slack.GroupConversation{
			Conversation: slack.Conversation{
				ID: ch,
			},
			// Since we don't join channels by name, only ID, let's strip the C prefix and use that as the name.
			Name: ch[1:],
		},
	}
	encoded, err := json.Marshal(&response)
	if err != nil {
		msg := fmt.Sprintf("Unable to encode response: %s", err.Error())
		log.Printf(msg)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}

	_, _ = w.Write(encoded)
}

// handle conversations.create
func createConversationHandler(w http.ResponseWriter, r *http.Request) {
	_, _ = w.Write([]byte(defaultConversationJSON))
}

// handle conversations.setTopic
func setConversationTopicHandler(w http.ResponseWriter, r *http.Request) {
	_, _ = w.Write([]byte(conversionPurposeTopicJSON))
}

// handle conversations.setPurpose
func setConversationPurposeHandler(w http.ResponseWriter, r *http.Request) {
	_, _ = w.Write([]byte(conversionPurposeTopicJSON))
}

// handle conversations.rename
func renameConversationHandler(w http.ResponseWriter, r *http.Request) {
	_, _ = w.Write([]byte(renameConversationJSON))
}

// handle conversations.invite
func inviteConversationHandler(w http.ResponseWriter, r *http.Request) {
	_, _ = w.Write([]byte(inviteConversationJSON))
}

// handle groups.list
func listGroupsHandler(w http.ResponseWriter, _ *http.Request) {
	_, _ = w.Write([]byte(defaultGroupsListJSON))
}

// handle reaction.Add
func reactionAddHandler(w http.ResponseWriter, _ *http.Request) {
	_, _ = w.Write([]byte(defaultOkJSON))
}

// handle chat.postMessage
func (sts *Server) postMessageHandler(w http.ResponseWriter, r *http.Request) {
	serverAddr := r.Context().Value(ServerBotHubNameContextKey).(string)
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		msg := fmt.Sprintf("error reading body: %s", err.Error())
		log.Printf(msg)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}
	values, vErr := url.ParseQuery(string(data))
	if vErr != nil {
		msg := fmt.Sprintf("Unable to decode query params: %s", vErr.Error())
		log.Printf(msg)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}

	ts := time.Now().Unix()
	resp := &struct {
		Ok      bool   `json:"ok"`
		Channel string `json:"channel"`
		Ts      string `json:"ts"`
		Text    string `json:"text"`
	}{
		Ok:      true,
		Channel: values.Get("channel"),
		Ts:      fmt.Sprintf("%d", ts),
		Text:    values.Get("text"),
	}

	m := slack.Message{}
	m.Type = "message"
	m.Channel = values.Get("channel")
	m.Timestamp = fmt.Sprintf("%d", ts)
	m.Text = values.Get("text")
	m.ThreadTimestamp = values.Get("thread_ts")
	if values.Get("as_user") != "true" {
		m.User = defaultNonBotUserID
		m.Username = defaultNonBotUserName
	} else {
		m.User = BotIDFromContext(r.Context())
		m.Username = BotNameFromContext(r.Context())
	}
	attachments := values.Get("attachments")
	if attachments != "" {
		decoded, err := url.QueryUnescape(attachments)
		if err != nil {
			msg := fmt.Sprintf("Unable to decode attachments: %s", err.Error())
			log.Printf(msg)
			http.Error(w, msg, http.StatusInternalServerError)
			return
		}
		var attaches []slack.Attachment
		aJErr := json.Unmarshal([]byte(decoded), &attaches)
		if aJErr != nil {
			msg := fmt.Sprintf("Unable to decode attachments string to json: %s", aJErr.Error())
			log.Printf(msg)
			http.Error(w, msg, http.StatusInternalServerError)
			return
		}
		m.Attachments = attaches
	}
	blocks := values.Get("blocks")
	if blocks != "" {
		decoded, err := url.QueryUnescape(blocks)
		if err != nil {
			msg := fmt.Sprintf("Unable to decode blocks: %s", err.Error())
			log.Printf(msg)
			http.Error(w, msg, http.StatusInternalServerError)
			return
		}
		var decodedBlocks slack.Blocks
		dbJErr := json.Unmarshal([]byte(decoded), &decodedBlocks)
		if dbJErr != nil {
			msg := fmt.Sprintf("Unable to decode blocks string to json: %s", dbJErr.Error())
			log.Printf(msg)
			http.Error(w, msg, http.StatusInternalServerError)
			return
		}
		m.Blocks = decodedBlocks
	}
	jsonMessage, jsonErr := json.Marshal(m)
	if jsonErr != nil {
		msg := fmt.Sprintf("Unable to marshal message: %s", jsonErr.Error())
		log.Printf(msg)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}
	go sts.queueForWebsocket(string(jsonMessage), serverAddr)
	_ = json.NewEncoder(w).Encode(resp)
}

// RTMConnectHandler generates a valid connection
func RTMConnectHandler(w http.ResponseWriter, r *http.Request) {
	_, err := ioutil.ReadAll(r.Body)
	if err != nil {
		msg := fmt.Sprintf("Error reading body: %s", err.Error())
		log.Printf(msg)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}
	wsurl := r.Context().Value(ServerWSContextKey).(string)
	if wsurl == "" {
		msg := "missing webservice url from context"
		log.Printf(msg)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}

	fullresponse := generateRTMInfo(r.Context(), wsurl)
	j, jErr := json.Marshal(fullresponse)
	if jErr != nil {
		msg := fmt.Sprintf("Unable to marshal response: %s", jErr.Error())
		log.Printf("Error: %s", msg)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}
	_, wErr := w.Write(j)
	if wErr != nil {
		log.Printf("Error writing response: %s", wErr.Error())
	}
}

func rtmStartHandler(w http.ResponseWriter, r *http.Request) {
	_, err := ioutil.ReadAll(r.Body)
	if err != nil {
		msg := fmt.Sprintf("Error reading body: %s", err.Error())
		log.Printf(msg)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}
	wsurl := r.Context().Value(ServerWSContextKey).(string)
	if wsurl == "" {
		msg := "missing webservice url from context"
		log.Printf(msg)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}

	fullresponse := generateRTMInfo(r.Context(), wsurl)
	j, jErr := json.Marshal(fullresponse)
	if jErr != nil {
		msg := fmt.Sprintf("Unable to marshal response: %s", jErr.Error())
		log.Printf("Error: %s", msg)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}
	_, wErr := w.Write(j)
	if wErr != nil {
		log.Printf("Error writing response: %s", wErr.Error())
	}
}

func (sts *Server) wsHandler(w http.ResponseWriter, r *http.Request) {
	Websocket(func(c *websocket.Conn) {
		serverAddr := r.Context().Value(ServerBotHubNameContextKey).(string)
		go handlePendingMessages(c, serverAddr)
		for {
			var (
				err   error
				m     json.RawMessage
				mtype string
			)

			if mtype, m, err = RTMRespEventType(c); err != nil {
				if websocket.IsUnexpectedCloseError(err) {
					return
				}

				log.Printf("read error: %s", err.Error())
				continue
			}

			switch mtype {
			case "ping":
				if err = RTMRespPong(c, m); err != nil {
					log.Println("ping error:", err)
				}
			default:
				sts.postProcessMessage(string(m), serverAddr)
			}
		}
	})(w, r)
}

// Websocket handler
func Websocket(delegate func(c *websocket.Conn)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		upgrader := websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		}
		c, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			msg := fmt.Sprintf("Unable to upgrade to ws connection: %s", err.Error())
			log.Print(msg)
			http.Error(w, msg, http.StatusInternalServerError)
			return
		}
		defer func() { _ = c.Close() }()
		delegate(c)
	}
}

// RTMServerSendGoodbye send a goodbye event
func RTMServerSendGoodbye(c *websocket.Conn) error {
	return c.WriteJSON(slack.Event{Type: "goodbye"})
}

// RTMRespEventType retrieve the event type from the next message
func RTMRespEventType(c *websocket.Conn) (t string, m json.RawMessage, err error) {
	var (
		evt slack.Event
	)

	if err = c.ReadJSON(&m); err != nil {
		return "", m, err
	}

	if err = json.Unmarshal(m, &evt); err != nil {
		return "", m, err
	}

	return evt.Type, m, nil
}

// RTMRespPong decode a ping and respond with a pong event.
func RTMRespPong(c *websocket.Conn, m json.RawMessage) (err error) {
	var (
		ping slack.Ping
		pong slack.Pong
	)

	if err = json.Unmarshal(m, &ping); err != nil {
		return err
	}

	pong = slack.Pong{
		Type:    "pong",
		ReplyTo: ping.ID,
	}

	if err = c.WriteJSON(pong); err != nil {
		return err
	}

	return nil
}
