package slack

import (
	"encoding/json"
	"net/http"
	"testing"
)

func postMessageInvalidChannelHandler(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	response, _ := json.Marshal(chatResponseFull{
		SlackResponse: SlackResponse{Ok: false, Error: "channel_not_found"},
	})
	rw.Write(response)
}

func TestPostMessageInvalidChannel(t *testing.T) {
	http.HandleFunc("/chat.postMessage", postMessageInvalidChannelHandler)
	once.Do(startServer)
	SLACK_API = "http://" + serverAddr + "/"
	api := New("testing-token")
	_, _, err := api.PostMessage("CXXXXXXXX", "hello", PostMessageParameters{})
	if err == nil {
		t.Errorf("Expected error: %s; instead succeeded", "channel_not_found")
		return
	}

	if err.Error() != "channel_not_found" {
		t.Errorf("Expected error: %s; received: %s", "channel_not_found", err)
		return
	}
}
