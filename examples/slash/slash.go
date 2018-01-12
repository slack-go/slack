package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"net/http"

	"github.com/nlopes/slack"
)

func main() {
	var (
		verificationToken string
	)

	flag.StringVar(&verificationToken, "token", "YOUR_VERIFICATION_TOKEN_HERE", "Your Slash Verification Token")
	flag.Parse()

	// Example 1 (very simple)
	http.HandleFunc("/slash", slack.SlashHandler(verificationToken, slash))

	// Example 2 (customize)
	http.HandleFunc("/slash2", func(w http.ResponseWriter, r *http.Request) {
		s := &slack.Slash{}
		s.Parse(r)

		if s.Token != verificationToken {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		switch s.Command {
		case "/echo":
			params := &slack.PostMessageParameters{Text: s.Text}
			b, err := json.Marshal(params)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.Write(b)
		default:
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	})
	fmt.Println("[INFO] Server listening")
	http.ListenAndServe(":3000", nil)
}

func slash(s *slack.Slash) (params *slack.PostMessageParameters, err error) {
	switch s.Command {
	case "/echo":
		params = &slack.PostMessageParameters{Text: s.Text}
	default:
		return nil, errors.New("Invalid command")
	}
	return params, nil
}
