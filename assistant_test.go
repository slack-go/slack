package slack

import (
	"encoding/json"
	"net/http"
	"testing"
)

func assistantThreadsSuggestedPromptsHandler(rw http.ResponseWriter, r *http.Request) {

	channelID := r.FormValue("channel_id")
	threadTS := r.FormValue("thread_ts")
	promptStr := r.FormValue("prompts")

	var prompts []AssistantThreadsPrompt
	err := json.Unmarshal([]byte(promptStr), &prompts)
	if err != nil {
		rw.Write([]byte(`{ "ok": false, "error": "errored" }`))
		return
	}

	rw.Header().Set("Content-Type", "application/json")

	if channelID == "" {
		rw.Write([]byte(`{ "ok": false, "error": "channel_id missing" }`))
		return
	}

	if threadTS == "" {
		rw.Write([]byte(`{ "ok": false, "error": "thread_ts missing" }`))
		return
	}

	if len(prompts) != 2 {
		rw.Write([]byte(`{ "ok": false, "error": "incorrect prompt count" }`))
		return
	}

	resp, _ := json.Marshal(&addBookmarkResponse{
		SlackResponse: SlackResponse{Ok: true},
	})
	rw.Write(resp)

}

func TestAssistantThreadsSuggestedPrompts(t *testing.T) {

	http.HandleFunc("/assistant.threads.setSuggestedPrompts", assistantThreadsSuggestedPromptsHandler)
	once.Do(startServer)
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

func setAssistantThreadsStatusHandler(rw http.ResponseWriter, r *http.Request) {

	channelID := r.FormValue("channel_id")
	threadTS := r.FormValue("thread_ts")
	status := r.FormValue("status")

	rw.Header().Set("Content-Type", "application/json")

	if channelID == "" {
		rw.Write([]byte(`{ "ok": false, "error": "channel_id missing" }`))
		return
	}

	if threadTS == "" {
		rw.Write([]byte(`{ "ok": false, "error": "thread_ts missing" }`))
		return
	}

	if status == "" {
		rw.Write([]byte(`{ "ok": false, "error": "status missing" }`))
		return
	}

	resp, _ := json.Marshal(&addBookmarkResponse{
		SlackResponse: SlackResponse{Ok: true},
	})
	rw.Write(resp)

}

func TestSetAssistantThreadsStatus(t *testing.T) {

	http.HandleFunc("/assistant.threads.setStatus", setAssistantThreadsStatusHandler)
	once.Do(startServer)
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

	http.HandleFunc("/assistant.threads.setTitle", assistantThreadsTitleHandler)
	once.Do(startServer)
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
