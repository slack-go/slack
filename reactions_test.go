package slack

import (
	"net/http"
	"reflect"
	"testing"
)

var (
	addedReaction ReactionParameters
)

func addReactionHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	addedReaction.Name = r.FormValue("name")
	addedReaction.File = r.FormValue("file")
	addedReaction.FileComment = r.FormValue("file_comment")
	addedReaction.Channel = r.FormValue("channel")
	addedReaction.Timestamp = r.FormValue("timestamp")
	w.Write([]byte(`{ "ok": true }`))
}

func TestSlack_AddReaction(t *testing.T) {
	http.HandleFunc("/reactions.add", addReactionHandler)
	once.Do(startServer)
	SLACK_API = "http://" + serverAddr + "/"
	api := New("testing-token")
	params := NewReactionParameters()
	params.Name = "thumbsup"
	params.File = "FileID"
	params.FileComment = "FileCommentID"
	params.Channel = "ChannelID"
	params.Timestamp = "123"
	addedReaction = ReactionParameters{}
	err := api.AddReaction(params)
	if err != nil {
		t.Fatalf("Unexpected error: %s", err)
	}
	if !reflect.DeepEqual(params, addedReaction) {
		t.Fatalf("Got reaction %#v, want %#v", addedReaction, params)
	}
}
