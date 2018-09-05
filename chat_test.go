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
	APIURL = "http://" + serverAddr + "/"
	api := New("testing-token")
	_, _, err := api.PostMessage("CXXXXXXXX", MsgOptionText("hello", false))
	if err == nil {
		t.Errorf("Expected error: channel_not_found; instead succeeded")
		return
	}

	if err.Error() != "channel_not_found" {
		t.Errorf("Expected error: channel_not_found; received: %s", err)
		return
	}
}
