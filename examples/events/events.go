package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ACollectionOfAtoms/slack"
)

var api = slack.New("TOKEN")

func main() {
	http.HandleFunc("/events-endpoint", func(w http.ResponseWriter, r *http.Request) {
		buf := new(bytes.Buffer)
		buf.ReadFrom(r.Body)
		body := buf.String()
		eventsAPIEvent, e := api.ParseEventsAPIEvent(json.RawMessage(body))
		if e != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}

		if eventsAPIEvent.Type == slack.URLVerification {
			var r *slack.ChallengeResponse
			err := json.Unmarshal([]byte(body), &r)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
			}
			w.Header().Set("Content-Type", "text")
			w.Write([]byte(r.Challenge))
		}
		if eventsAPIEvent.Type == slack.CallbackEvent {
			postParams := slack.PostMessageParameters{}
			innerEvent := eventsAPIEvent.InnerEvent
			switch ev := innerEvent.Data.(type) {
			case *slack.AppMentionEvent:
				api.PostMessage(ev.Channel, "Yes, hello.", postParams)
			}
		}
	})
	fmt.Println("[INFO] Server listening")
	http.ListenAndServe(":3000", nil)
}
