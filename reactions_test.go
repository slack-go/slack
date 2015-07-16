package slack

import (
	"net/http"
	"reflect"
	"testing"
)

func init() {
	http.HandleFunc("/reactions.add", addReactionHandler)
	http.HandleFunc("/reactions.get", getReactionHandler)
}

var (
	gotParams         map[string]string
	addedReaction     Reaction
	getReactionRes    string
	gottenReactionRef ItemRef
)

func accumulateFormValue(k string, r *http.Request) {
	if v := r.FormValue(k); v != "" {
		gotParams[k] = v
	}
}

func addReactionHandler(w http.ResponseWriter, r *http.Request) {
	accumulateFormValue("name", r)
	accumulateFormValue("file", r)
	accumulateFormValue("file_comment", r)
	accumulateFormValue("channel", r)
	accumulateFormValue("timestamp", r)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{ "ok": true }`))
}

func getReactionHandler(w http.ResponseWriter, r *http.Request) {
	accumulateFormValue("file", r)
	accumulateFormValue("file_comment", r)
	accumulateFormValue("channel", r)
	accumulateFormValue("timestamp", r)
	accumulateFormValue("full", r)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(getReactionRes))
}

func TestSlack_AddReaction_ToMessage(t *testing.T) {
	once.Do(startServer)
	SLACK_API = "http://" + serverAddr + "/"
	api := New("testing-token")
	wantParams := map[string]string{
		"name":      "thumbsup",
		"channel":   "ChannelID",
		"timestamp": "123",
	}
	gotParams = map[string]string{}
	params := NewAddReactionParameters("thumbsup", NewRefToMessage("ChannelID", "123"))
	err := api.AddReaction(params)
	if err != nil {
		t.Fatalf("Unexpected error: %s", err)
	}
	if !reflect.DeepEqual(gotParams, wantParams) {
		t.Errorf("Got params %#v, want %#v", gotParams, wantParams)
	}
}

func TestSlack_AddReaction_ToFile(t *testing.T) {
	once.Do(startServer)
	SLACK_API = "http://" + serverAddr + "/"
	api := New("testing-token")
	wantParams := map[string]string{
		"name": "thumbsup",
		"file": "FileID",
	}
	gotParams = map[string]string{}
	params := NewAddReactionParameters("thumbsup", NewRefToFile("FileID"))
	err := api.AddReaction(params)
	if err != nil {
		t.Fatalf("Unexpected error: %s", err)
	}
	if !reflect.DeepEqual(gotParams, wantParams) {
		t.Errorf("Got params %#v, want %#v", gotParams, wantParams)
	}
}

func TestSlack_AddReaction_ToFileComment(t *testing.T) {
	once.Do(startServer)
	SLACK_API = "http://" + serverAddr + "/"
	api := New("testing-token")
	wantParams := map[string]string{
		"name":         "thumbsup",
		"file_comment": "FileCommentID",
	}
	gotParams = map[string]string{}
	params := NewAddReactionParameters("thumbsup", NewRefToFileComment("FileCommentID"))
	err := api.AddReaction(params)
	if err != nil {
		t.Fatalf("Unexpected error: %s", err)
	}
	if !reflect.DeepEqual(gotParams, wantParams) {
		t.Errorf("Got params %#v, want %#v", gotParams, wantParams)
	}
}

func TestSlack_GetReaction_ToMessage(t *testing.T) {
	once.Do(startServer)
	SLACK_API = "http://" + serverAddr + "/"
	api := New("testing-token")
	getReactionRes = `{"ok": true,
    "message": {
        "type": "message",
        "message": {
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
    }}`
	want := []ItemReaction{
		ItemReaction{Name: "astonished", Count: 3, Users: []string{"U1", "U2", "U3"}},
		ItemReaction{Name: "clock1", Count: 3, Users: []string{"U1", "U2"}},
	}
	wantParams := map[string]string{
		"channel":   "ChannelID",
		"timestamp": "123",
	}
	gotParams = map[string]string{}
	params := NewGetReactionParameters(NewRefToMessage("ChannelID", "123"))
	got, err := api.GetReactions(params)
	if err != nil {
		t.Fatalf("Unexpected error: %s", err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Got reaction %#v, want %#v", got, want)
	}
	if !reflect.DeepEqual(gotParams, wantParams) {
		t.Errorf("Got params %#v, want %#v", gotParams, wantParams)
	}
}

func TestSlack_GetReaction_ToFile(t *testing.T) {
	once.Do(startServer)
	SLACK_API = "http://" + serverAddr + "/"
	api := New("testing-token")
	getReactionRes = `{"ok": true,
    "message": {
        "type": "file",
        "file": {
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
    }}`
	want := []ItemReaction{
		ItemReaction{Name: "astonished", Count: 3, Users: []string{"U1", "U2", "U3"}},
		ItemReaction{Name: "clock1", Count: 3, Users: []string{"U1", "U2"}},
	}
	wantParams := map[string]string{
		"file": "FileID",
		"full": "true",
	}
	gotParams = map[string]string{}
	params := NewGetReactionParameters(NewRefToFile("FileID"))
	params.Full = true
	got, err := api.GetReactions(params)
	if err != nil {
		t.Fatalf("Unexpected error: %s", err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Got reaction %#v, want %#v", got, want)
	}
	if !reflect.DeepEqual(gotParams, wantParams) {
		t.Errorf("Got params %#v, want %#v", gotParams, wantParams)
	}
}

func TestSlack_GetReaction_ToFileComment(t *testing.T) {
	once.Do(startServer)
	SLACK_API = "http://" + serverAddr + "/"
	api := New("testing-token")
	getReactionRes = `{"ok": true,
    "message": {
        "type": "file_comment",
        "file_comment": {
	    "comment": {
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
        }
    }}`
	want := []ItemReaction{
		ItemReaction{Name: "astonished", Count: 3, Users: []string{"U1", "U2", "U3"}},
		ItemReaction{Name: "clock1", Count: 3, Users: []string{"U1", "U2"}},
	}
	wantParams := map[string]string{
		"file_comment": "FileCommentID",
	}
	gotParams = map[string]string{}
	params := NewGetReactionParameters(NewRefToFileComment("FileCommentID"))
	got, err := api.GetReactions(params)
	if err != nil {
		t.Fatalf("Unexpected error: %s", err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Got reaction %#v, want %#v", got, want)
	}
	if !reflect.DeepEqual(gotParams, wantParams) {
		t.Errorf("Got params %#v, want %#v", gotParams, wantParams)
	}
}
