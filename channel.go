package slack

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
)

type ChannelHistory struct {
	Ok       bool      `json:"ok"`
	Latest   string    `json:"latest"`
	Messages []Message `json:"messages"`
	HasMore  bool      `json:"has_more"`

	Error string `json:"error"`
}

type ChannelTopic struct {
	Value   string   `json:"value"`
	Creator string   `json:"creator"`
	LastSet JSONTime `json:"last_set"`
}

type Channel struct {
	Id          string       `json:"id"`
	Name        string       `json:"name"`
	IsChannel   bool         `json:"is_channel"`
	Creator     string       `json:"creator"`
	IsArchived  bool         `json:"is_archived"`
	IsGeneral   bool         `json:"is_general"`
	Members     []string     `json:"members"`
	Topic       ChannelTopic `json:"topic"`
	Created     JSONTime     `json:"created"`
	IsMember    bool         `json:"is_member"`
	LastRead    string       `json:"last_read"`
	Latest      Message      `json:"latest"`
	UnreadCount int          `json:"unread_count"`
}

func (api *SlackAPI) GetChannelHistory(channel_id string, latest string) ChannelHistory {
	channel_history := ChannelHistory{}
	resp, err := http.PostForm(SLACK_API+"channels.history",
		url.Values{
			"token":   {api.config.token},
			"channel": {channel_id},
			"latest":  {latest},
		})
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	decoder := json.NewDecoder(resp.Body)
	if err = decoder.Decode(&channel_history); err != nil {
		log.Fatal(err)
	}
	if !channel_history.Ok {
		log.Fatal(channel_history.Error)
	}
	return channel_history
}
