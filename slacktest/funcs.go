package slacktest

import (
	"context"
	"fmt"
	"log"
	"time"

	websocket "github.com/gorilla/websocket"

	slack "github.com/slack-go/slack"
)

func (sts *Server) queueForWebsocket(s, hubname string) {
	channel, err := getHubForServer(hubname)
	if err != nil {
		log.Printf("Unable to get server's channels: %s", err.Error())
	}
	sts.seenOutboundMessages.Lock()
	sts.seenOutboundMessages.messages = append(sts.seenOutboundMessages.messages, s)
	sts.seenOutboundMessages.Unlock()
	channel.sent <- s
}

func handlePendingMessages(c *websocket.Conn, hubname string) {
	channel, err := getHubForServer(hubname)
	if err != nil {
		log.Printf("Unable to get server's channels: %s", err.Error())
		return
	}
	for m := range channel.sent {
		err := c.WriteMessage(websocket.TextMessage, []byte(m))
		if err != nil {
			log.Printf("error writing message to websocket: %s", err.Error())
			continue
		}
	}
}

func (sts *Server) postProcessMessage(m, hubname string) {
	channel, err := getHubForServer(hubname)
	if err != nil {
		log.Printf("Unable to get server's channels: %s", err.Error())
		return
	}
	sts.seenInboundMessages.Lock()
	sts.seenInboundMessages.messages = append(sts.seenInboundMessages.messages, m)
	sts.seenInboundMessages.Unlock()
	// send to firehose
	channel.seen <- m
}

func newHub() *hub {
	h := &hub{}
	c := make(map[string]*messageChannels)
	h.serverChannels = c
	return h
}

func addServerToHub(s *Server, channels *messageChannels) error {
	if s.ServerAddr == "" {
		return ErrEmptyServerToHub
	}
	masterHub.Lock()
	masterHub.serverChannels[s.ServerAddr] = channels
	masterHub.Unlock()
	return nil
}

func getHubForServer(serverAddr string) (*messageChannels, error) {
	if serverAddr == "" {
		return &messageChannels{}, ErrPassedEmptyServerAddr
	}
	masterHub.RLock()
	defer masterHub.RUnlock()
	channels, ok := masterHub.serverChannels[serverAddr]
	if !ok {
		return &messageChannels{}, ErrNoQueuesRegisteredForServer
	}
	return channels, nil
}

// BotNameFromContext returns the botname from a provided context
func BotNameFromContext(ctx context.Context) string {
	botname, ok := ctx.Value(ServerBotNameContextKey).(string)
	if !ok {
		return defaultBotName
	}
	return botname
}

// BotIDFromContext returns the bot userid from a provided context
func BotIDFromContext(ctx context.Context) string {
	botname, ok := ctx.Value(ServerBotIDContextKey).(string)
	if !ok {
		return defaultBotID
	}
	return botname
}

// generate a full rtminfo response for initial rtm connections
func generateRTMInfo(ctx context.Context, wsurl string) *fullInfoSlackResponse {
	rtmInfo := slack.Info{
		URL:  wsurl,
		Team: defaultTeam,
		User: defaultBotInfo,
	}
	rtmInfo.User.ID = BotIDFromContext(ctx)
	rtmInfo.User.Name = BotNameFromContext(ctx)
	return &fullInfoSlackResponse{
		rtmInfo,
		okWebResponse,
	}
}

func nowAsJSONTime() slack.JSONTime {
	return slack.JSONTime(time.Now().Unix())
}

func defaultBotInfoJSON(ctx context.Context) string {
	botid := BotIDFromContext(ctx)
	botname := BotNameFromContext(ctx)
	return fmt.Sprintf(`
		{
			"ok":true,
			"bot":{
					"id": "%s",
					"app_id": "A4H1JB4AZ",
					"deleted": false,
					"name": "%s",
					"icons": {
						"image_36": "https://localhost.localdomain/img36.png",
						"image_48": "https://localhost.localdomain/img48.png",
						"image_72": "https://localhost.localdomain/img72.png"
					}
				}
		}
		`, botid, botname)
}
