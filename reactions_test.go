package slack

import (
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

func init() {
	http.HandleFunc("/reactions.add", addReactionHandler)
	http.HandleFunc("/reactions.get", getReactionHandler)
	http.HandleFunc("/reactions.list", listReactionHandler)
}

var (
	gotParams       map[string]string
	getReactionRes  string
	listReactionRes string
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

func listReactionHandler(w http.ResponseWriter, r *http.Request) {
	accumulateFormValue("user", r)
	accumulateFormValue("count", r)
	accumulateFormValue("full", r)
	accumulateFormValue("page", r)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(listReactionRes))
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

func TestSlack_ListReactions(t *testing.T) {
	once.Do(startServer)
	SLACK_API = "http://" + serverAddr + "/"
	api := New("testing-token")
	listReactionRes = `{"ok": true,
    "items": [
        {
            "type": "message",
            "message": {
                "text": "hello",
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
        },
        {
            "type": "file",
            "file": {
                "name": "toy",
                "reactions": [
                    {
                        "name": "clock1",
                        "count": 3,
                        "users": [ "U1", "U2" ]
                    }
                ]
            }
        },
        {
            "type": "file_comment",
            "file_comment": {
                "file": {},
                "comment": {
                    "comment": "cool toy",
                    "reactions": [
                        {
                            "name": "astonished",
                            "count": 3,
                            "users": [ "U1", "U2", "U3" ]
                        }
                    ]
                }
            }
        }
    ],
    "paging": {
        "count": 100,
        "total": 4,
        "page": 1,
        "pages": 1
    }}`
	want := []ReactedItem{
		ReactedItem{
			Type:    "message",
			Message: &Message{Msg: Msg{Text: "hello"}},
			Reactions: []ItemReaction{
				ItemReaction{Name: "astonished", Count: 3, Users: []string{"U1", "U2", "U3"}},
				ItemReaction{Name: "clock1", Count: 3, Users: []string{"U1", "U2"}},
			},
		},
		ReactedItem{
			Type: "file",
			File: &File{Name: "toy"},
			Reactions: []ItemReaction{
				ItemReaction{Name: "clock1", Count: 3, Users: []string{"U1", "U2"}},
			},
		},
		ReactedItem{
			Type:    "file_comment",
			Comment: &Comment{Comment: "cool toy"},
			Reactions: []ItemReaction{
				ItemReaction{Name: "astonished", Count: 3, Users: []string{"U1", "U2", "U3"}},
			},
		},
	}
	wantParams := map[string]string{
		"user":  "UserID",
		"count": "200",
		"page":  "2",
		"full":  "true",
	}
	gotParams = map[string]string{}
	params := NewListReactionsParameters("UserID")
	params.Count = 200
	params.Page = 2
	params.Full = true
	got, paging, err := api.ListReactions(params)
	if err != nil {
		t.Fatalf("Unexpected error: %s", err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Got reaction %#v, want %#v", got, want)
		for i, item := range got {
			fmt.Printf("Item %d, Type: %s\n", i, item.Type)
			fmt.Printf("Message  %#v\n", item.Message)
			fmt.Printf("File     %#v\n", item.File)
			fmt.Printf("Comment  %#v\n", item.Comment)
			fmt.Printf("Reactions %#v\n", item.Reactions)
		}
	}
	if !reflect.DeepEqual(gotParams, wantParams) {
		t.Errorf("Got params %#v, want %#v", gotParams, wantParams)
	}
	if reflect.DeepEqual(paging, Paging{}) {
		t.Errorf("Want paging data, got empty struct")
	}
}
