package slack

import (
	"net/http"
	"reflect"
	"testing"
)

var (
	addedReaction Reaction
)

func addReactionHandler(w http.ResponseWriter, r *http.Request) {
	addedReaction.Name = r.FormValue("name")
	addedReaction.File = r.FormValue("file")
	addedReaction.FileComment = r.FormValue("file_comment")
	addedReaction.Channel = r.FormValue("channel")
	addedReaction.Timestamp = r.FormValue("timestamp")
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{ "ok": true }`))
}

func TestSlack_AddReaction(t *testing.T) {
	http.HandleFunc("/reactions.add", addReactionHandler)
	once.Do(startServer)
	SLACK_API = "http://" + serverAddr + "/"
	api := New("testing-token")
	r := Reaction{}
	r.File = "FileID"
	r.FileComment = "FileCommentID"
	r.Channel = "ChannelID"
	r.Timestamp = "123"
	addedReaction = Reaction{}
	err := api.AddReaction(r)
	if err != nil {
		t.Fatalf("Unexpected error: %s", err)
	}
	if got := addedReaction; !reflect.DeepEqual(got, r) {
		t.Errorf("Got reaction %#v, want %#v", got, r)
	}
}
