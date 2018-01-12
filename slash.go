package slack

import (
	"encoding/json"
	"net/http"
)

// Slash contains information about a request of the slash command
type Slash struct {
	Token       string `json:"token"`
	TeamID      string `json:"team_id"`
	TeamDomain  string `json:"team_domain"`
	ChannelID   string `json:"channel_id"`
	ChannelName string `json:"channel_name"`
	UserID      string `json:"user_id"`
	UserName    string `json:"user_name"`
	Command     string `json:"command"`
	Text        string `json:"text"`
	ResponseURL string `json:"response_url"`
	TriggerID   string `json:"trigger_id"`
}

// Parse will parse the request of the slash command
func (s *Slash) Parse(r *http.Request) error {
	if err := r.ParseForm(); err != nil {
		return err
	}
	s.Token = r.PostForm.Get("token")
	s.TeamID = r.PostForm.Get("team_id")
	s.TeamDomain = r.PostForm.Get("team_domain")
	s.ChannelID = r.PostForm.Get("channel_id")
	s.ChannelName = r.PostForm.Get("channel_name")
	s.UserID = r.PostForm.Get("user_id")
	s.UserName = r.PostForm.Get("user_name")
	s.Command = r.PostForm.Get("command")
	s.Text = r.PostForm.Get("text")
	s.ResponseURL = r.PostForm.Get("response_url")
	s.TriggerID = r.PostForm.Get("trigger_id")
	return nil
}

// SlashHandler returns func for http.Handler
func SlashHandler(verificationToken string,
	handler func(*Slash) (*PostMessageParameters, error)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
		}

		s := &Slash{}
		if err := s.Parse(r); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if s.Token != verificationToken {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		params, err := handler(s)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		b, err := json.Marshal(params)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(b)
	}
}
