package slack

import (
	"net/http"
	"reflect"
	"testing"
)

var (
	addedReaction     Reaction
	gottenReactionRef ItemRef
)

func addReactionHandler(w http.ResponseWriter, r *http.Request) {
	addedReaction.Name = r.FormValue("name")
	addedReaction.FileId = r.FormValue("file")
	addedReaction.FileCommentId = r.FormValue("file_comment")
	addedReaction.ChannelId = r.FormValue("channel")
	addedReaction.Timestamp = r.FormValue("timestamp")
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{ "ok": true }`))
}

func getReactionHandler(w http.ResponseWriter, r *http.Request) {
	gottenReactionRef.FileId = r.FormValue("file")
	gottenReactionRef.FileCommentId = r.FormValue("file_comment")
	gottenReactionRef.ChannelId = r.FormValue("channel")
	gottenReactionRef.Timestamp = r.FormValue("timestamp")
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"ok": true,
    "message": {
        "type": "message",
        "channel": "C2147483705",
        "message": {
            "type": "message",
            "ts": "111",
            "reactions": [
                {
                    "name": "astonished",
                    "count": 3,
                    "users": [ "U1", "U2", "U3" ]
                },
                {
                    "name": "clock1",
                    "count": 3,
                    "users": [ "U1", "U2" ]
                }
            ]
        }
    }}`))
}

func TestSlack_AddReaction(t *testing.T) {
	http.HandleFunc("/reactions.add", addReactionHandler)
	once.Do(startServer)
	SLACK_API = "http://" + serverAddr + "/"
	api := New("testing-token")
	want := Reaction{}
	want.Name = "thumbsup"
	want.FileId = "FileID"
	want.FileCommentId = "FileCommentID"
	want.ChannelId = "ChannelID"
	want.Timestamp = "123"
	addedReaction = Reaction{}
	err := api.AddReaction("thumbsup", want.ItemRef)
	if err != nil {
		t.Fatalf("Unexpected error: %s", err)
	}
	if got := addedReaction; !reflect.DeepEqual(got, want) {
		t.Errorf("Got reaction %#v, want %#v", got, want)
	}
}

func TestSlack_GetReaction(t *testing.T) {
	http.HandleFunc("/reactions.get", getReactionHandler)
	once.Do(startServer)
	SLACK_API = "http://" + serverAddr + "/"
	api := New("testing-token")
	ref := ItemRef{}
	ref.FileId = "FileID"
	ref.FileCommentId = "FileCommentID"
	ref.ChannelId = "ChannelID"
	ref.Timestamp = "123"
	gottenReactionRef = ItemRef{}
	got, err := api.GetReactions(ref)
	if err != nil {
		t.Fatalf("Unexpected error: %s", err)
	}
	if got, want := gottenReactionRef, ref; !reflect.DeepEqual(got, want) {
		t.Errorf("Got reaction ref %#v, want %#v", got, want)
	}
	want := []ItemReaction{
		ItemReaction{Name: "astonished", Count: 3, Users: []string{"U1", "U2", "U3"}},
		ItemReaction{Name: "clock1", Count: 3, Users: []string{"U1", "U2"}},
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Got reaction %#v, want %#v", got, want)
	}
}
