package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/slack-go/slack"
)

func main() {
	// Get signing secret from environment variable
	signingSecret := os.Getenv("SLACK_SIGNING_SECRET")
	if signingSecret == "" {
		fmt.Println("SLACK_SIGNING_SECRET environment variable is required")
		os.Exit(1)
	}

	http.HandleFunc("/slash", func(w http.ResponseWriter, r *http.Request) {
		verifier, err := slack.NewSecretsVerifier(r.Header, signingSecret)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		r.Body = io.NopCloser(io.TeeReader(r.Body, &verifier))
		s, err := slack.SlashCommandParse(r)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if err = verifier.Ensure(); err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		switch s.Command {
		case "/echo":
			params := &slack.Msg{Text: s.Text}
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
