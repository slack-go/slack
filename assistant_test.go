package slack

import (
	"encoding/json"
	"net/http"
	"testing"
)

func TestAssistantThreadsSuggestedPrompts(t *testing.T) {
	s := startServer()
	s.RegisterHandler("/assistant.threads.setSuggestedPrompts", okJSONHandler)
	defer s.Close()

	api := New("testing-token", OptionAPIURL("http://"+serverAddr+"/"))

	params := AssistantThreadsSetSuggestedPromptsParameters{
		ChannelID: "CXXXXXXXX",
		ThreadTS:  "1234567890.123456",
	}

	params.AddPrompt("title1", "message1")
	params.AddPrompt("title2", "message2")

	err := api.SetAssistantThreadsSuggestedPrompts(params)
	if err != nil {
		t.Fatalf("Unexpected error: %s", err)
	}

}

func TestSetAssistantThreadsStatus(t *testing.T) {
	s := startServer()
	s.RegisterHandler("/assistant.threads.setStatus", okJSONHandler)
	defer s.Close()

	api := New("testing-token", OptionAPIURL("http://"+serverAddr+"/"))

	params := AssistantThreadsSetStatusParameters{
		ChannelID: "CXXXXXXXX",
		ThreadTS:  "1234567890.123456",
		Status:    "updated status",
	}

	err := api.SetAssistantThreadsStatus(params)
	if err != nil {
		t.Fatalf("Unexpected error: %s", err)
	}

}

func assistantThreadsTitleHandler(rw http.ResponseWriter, r *http.Request) {
	channelID := r.FormValue("channel_id")
	threadTS := r.FormValue("thread_ts")
	title := r.FormValue("title")

	rw.Header().Set("Content-Type", "application/json")

	if channelID != "" && threadTS != "" && title != "" {
		resp, _ := json.Marshal(&addBookmarkResponse{
			SlackResponse: SlackResponse{Ok: true},
		})
		rw.Write(resp)
	} else {
		rw.Write([]byte(`{ "ok": false, "error": "errored" }`))
	}
}

func TestSetAssistantThreadsTitle(t *testing.T) {
	s := startServer()
	s.RegisterHandler("/assistant.threads.setTitle", assistantThreadsTitleHandler)
	defer s.Close()

	api := New("testing-token", OptionAPIURL("http://"+serverAddr+"/"))

	params := AssistantThreadsSetTitleParameters{
		ChannelID: "CXXXXXXXX",
		ThreadTS:  "1234567890.123456",
		Title:     "updated title",
	}

	err := api.SetAssistantThreadsTitle(params)
	if err != nil {
		t.Fatalf("Unexpected error: %s", err)
	}
}
