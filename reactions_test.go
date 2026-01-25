package slack

import (
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

type reactionsHandler struct {
	gotParams map[string]string
	response  string
}

func newReactionsHandler() *reactionsHandler {
	return &reactionsHandler{
		gotParams: make(map[string]string),
		response:  `{ "ok": true }`,
	}
}

func (rh *reactionsHandler) accumulateFormValue(k string, r *http.Request) {
	if v := r.FormValue(k); v != "" {
		rh.gotParams[k] = v
	}
}

func (rh *reactionsHandler) handler(w http.ResponseWriter, r *http.Request) {
	rh.accumulateFormValue("channel", r)
	rh.accumulateFormValue("count", r)
	rh.accumulateFormValue("file", r)
	rh.accumulateFormValue("file_comment", r)
	rh.accumulateFormValue("full", r)
	rh.accumulateFormValue("name", r)
	rh.accumulateFormValue("page", r)
	rh.accumulateFormValue("timestamp", r)
	rh.accumulateFormValue("user", r)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(rh.response))
}

func TestSlack_AddReaction(t *testing.T) {
	once.Do(startServer)
	api := New("testing-token", OptionAPIURL("http://"+serverAddr+"/"))
	tests := []struct {
		name       string
		ref        ItemRef
		wantParams map[string]string
	}{
		{
			"thumbsup",
			NewRefToMessage("ChannelID", "123"),
			map[string]string{
				"name":      "thumbsup",
				"channel":   "ChannelID",
				"timestamp": "123",
			},
		},
		{
			"thumbsup",
			NewRefToFile("FileID"),
			map[string]string{
				"name": "thumbsup",
				"file": "FileID",
			},
		},
		{
			"thumbsup",
			NewRefToComment("FileCommentID"),
			map[string]string{
				"name":         "thumbsup",
				"file_comment": "FileCommentID",
			},
		},
	}
	var rh *reactionsHandler
	http.HandleFunc("/reactions.add", func(w http.ResponseWriter, r *http.Request) { rh.handler(w, r) })
	for i, test := range tests {
		rh = newReactionsHandler()
		err := api.AddReaction(test.name, test.ref)
		if err != nil {
			t.Fatalf("%d: Unexpected error: %s", i, err)
		}
		if !reflect.DeepEqual(rh.gotParams, test.wantParams) {
			t.Errorf("%d: Got params %#v, want %#v", i, rh.gotParams, test.wantParams)
		}
	}
}

func TestSlack_RemoveReaction(t *testing.T) {
	once.Do(startServer)
	api := New("testing-token", OptionAPIURL("http://"+serverAddr+"/"))
	tests := []struct {
		name       string
		ref        ItemRef
		wantParams map[string]string
	}{
		{
			"thumbsup",
			NewRefToMessage("ChannelID", "123"),
			map[string]string{
				"name":      "thumbsup",
				"channel":   "ChannelID",
				"timestamp": "123",
			},
		},
		{
			"thumbsup",
			NewRefToFile("FileID"),
			map[string]string{
				"name": "thumbsup",
				"file": "FileID",
			},
		},
		{
			"thumbsup",
			NewRefToComment("FileCommentID"),
			map[string]string{
				"name":         "thumbsup",
				"file_comment": "FileCommentID",
			},
		},
	}
	var rh *reactionsHandler
	http.HandleFunc("/reactions.remove", func(w http.ResponseWriter, r *http.Request) { rh.handler(w, r) })
	for i, test := range tests {
		rh = newReactionsHandler()
		err := api.RemoveReaction(test.name, test.ref)
		if err != nil {
			t.Fatalf("%d: Unexpected error: %s", i, err)
		}
		if !reflect.DeepEqual(rh.gotParams, test.wantParams) {
			t.Errorf("%d: Got params %#v, want %#v", i, rh.gotParams, test.wantParams)
		}
	}
}

func TestSlack_GetReactions(t *testing.T) {
	once.Do(startServer)
	api := New("testing-token", OptionAPIURL("http://"+serverAddr+"/"))
	tests := []struct {
		ref             ItemRef
		params          GetReactionsParameters
		wantParams      map[string]string
		json            string
		wantReactedItem ReactedItem
	}{
		{
			NewRefToMessage("ChannelID", "123"),
			GetReactionsParameters{},
			map[string]string{
				"channel":   "ChannelID",
				"timestamp": "123",
			},
			`{"ok": true,
		 "type": "message",
		 "channel": "ChannelID",
		 "message": {
			"text": "lorem ipsum dolor sit amet",
			"ts": "123",
			"user": "U2147483828",
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
		 }}`,
			ReactedItem{
				Item: Item{
					Type:    "message",
					Channel: "ChannelID",
					Message: &Message{
						Msg: Msg{
							Text:      "lorem ipsum dolor sit amet",
							User:      "U2147483828",
							Timestamp: "123",
						},
					},
				},
				Reactions: []ItemReaction{
					{Name: "astonished", Count: 3, Users: []string{"U1", "U2", "U3"}},
					{Name: "clock1", Count: 3, Users: []string{"U1", "U2"}},
				},
			},
		},
		{
			NewRefToFile("FileID"),
			GetReactionsParameters{Full: true},
			map[string]string{
				"file": "FileID",
				"full": "true",
			},
			`{"ok": true,
    "type": "file",
    "file": {
		  "id": "F0A12BCDE",
		  "created": 1531763342,
		  "timestamp": 1531763342,
		  "name": "tedair.gif",
		  "title": "tedair.gif",
		  "mimetype": "image/gif",
		  "filetype": "gif",
		  "pretty_type": "GIF",
		  "user": "U012A3BCD",
		  "editable": false,
		  "size": 137531,
		  "mode": "hosted",
    	  "is_external": false,
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
    }}`,
			ReactedItem{
				Item: Item{
					Type: "file", File: &File{
						Name:      "tedair.gif",
						ID:        "F0A12BCDE",
						Created:   1531763342,
						Timestamp: 1531763342,
						User:      "U012A3BCD",
						Editable:  false,
						Size:      137531,
					},
				},
				Reactions: []ItemReaction{
					{Name: "astonished", Count: 3, Users: []string{"U1", "U2", "U3"}},
					{Name: "clock1", Count: 3, Users: []string{"U1", "U2"}},
				},
			},
		},
		{
			NewRefToComment("FileCommentID"),
			GetReactionsParameters{},
			map[string]string{
				"file_comment": "FileCommentID",
			},
			`{"ok": true,
    "type": "file_comment",
    "file": {
	 	  "id": "F0A12BCDE",
		  "created": 1531763342,
		  "timestamp": 1531763342,
		  "name": "tedair.gif",
		  "title": "tedair.gif",
		  "mimetype": "image/gif",
		  "filetype": "gif",
		  "pretty_type": "GIF",
		  "user": "U012A3BCD",
		  "editable": false,
		  "size": 137531,
		  "is_external": false
	 },
    "comment": {
		  "comment": "lorem ipsum dolor sit amet comment",
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
    }}`,
			ReactedItem{
				Item: Item{
					Type: "file_comment", File: &File{
						Name: "tedair.gif",
					},
					Comment: &Comment{
						Comment: "lorem ipsum dolor sit amet comment",
					},
				},
				Reactions: []ItemReaction{
					{Name: "astonished", Count: 3, Users: []string{"U1", "U2", "U3"}},
					{Name: "clock1", Count: 3, Users: []string{"U1", "U2"}},
				},
			},
		},
	}
	var rh *reactionsHandler
	http.HandleFunc("/reactions.get", func(w http.ResponseWriter, r *http.Request) { rh.handler(w, r) })
	for i, test := range tests {
		rh = newReactionsHandler()
		rh.response = test.json
		got, err := api.GetReactions(test.ref, test.params)
		if err != nil {
			t.Fatalf("%d: Unexpected error: %s", i, err)
		}
		if !reflect.DeepEqual(got.Reactions, test.wantReactedItem.Reactions) {
			t.Errorf("%d: Got reaction %#v, want %#v", i, got.Reactions, test.wantReactedItem.Reactions)
		}
		if !reflect.DeepEqual(rh.gotParams, test.wantParams) {
			t.Errorf("%d: Got params %#v, want %#v", i, rh.gotParams, test.wantParams)
		}

		switch got.Type {
		case "message":
			if got.Message == nil {
				t.Fatalf("%d: Got message %#v, want %#v", i, got.Message, test.wantReactedItem.Message)
			}

			if got.Message.Text != test.wantReactedItem.Message.Text {
				t.Errorf("%d: Got message text %#v, want %#v", i, got.Message.Text, test.wantReactedItem.Message.Text)
			}
			if got.Channel != test.wantReactedItem.Channel {
				t.Errorf("%d: Got channel %#v, want %#v", i, got.Channel, test.wantReactedItem.Channel)
			}
			if got.Message.User != test.wantReactedItem.Message.User {
				t.Errorf("%d: Got message user %#v, want %#v", i, got.Message.User, test.wantReactedItem.Message.User)
			}
			if got.Message.Timestamp != test.wantReactedItem.Message.Timestamp {
				t.Errorf("%d: Got message timestamp %#v, want %#v", i, got.Message.Timestamp, test.wantReactedItem.Message.Timestamp)
			}
		case "file":
			if got.File == nil {
				t.Fatalf("%d: Got file %#v, want %#v", i, got.File, test.wantReactedItem.File)
			}
			if got.File.Name != test.wantReactedItem.File.Name {
				t.Errorf("%d: Got file name %#v, want %#v", i, got.File.Name, test.wantReactedItem.File.Name)
			}
			if got.File.ID != test.wantReactedItem.File.ID {
				t.Errorf("%d: Got file ID %#v, want %#v", i, got.File.ID, test.wantReactedItem.File.ID)
			}
			if got.File.Created != test.wantReactedItem.File.Created {
				t.Errorf("%d: Got file created %#v, want %#v", i, got.File.Created, test.wantReactedItem.File.Created)
			}
			if got.File.Timestamp != test.wantReactedItem.File.Timestamp {
				t.Errorf("%d: Got file timestamp %#v, want %#v", i, got.File.Timestamp, test.wantReactedItem.File.Timestamp)
			}
			if got.File.User != test.wantReactedItem.File.User {
				t.Errorf("%d: Got file user %#v, want %#v", i, got.File.User, test.wantReactedItem.File.User)
			}
			if got.File.Editable != test.wantReactedItem.File.Editable {
				t.Errorf("%d: Got file editable %#v, want %#v", i, got.File.Editable, test.wantReactedItem.File.Editable)
			}
			if got.File.Size != test.wantReactedItem.File.Size {
				t.Errorf("%d: Got file size %#v, want %#v", i, got.File.Size, test.wantReactedItem.File.Size)
			}
		case "file_comment":
			if got.Comment == nil {
				t.Fatalf("%d: Got comment %#v, want %#v", i, got.Comment, test.wantReactedItem.Comment)
			}
			if got.File == nil {
				t.Fatalf("%d: Got file %#v, want %#v", i, got.File, test.wantReactedItem.File)
			}
			if got.File.Name != test.wantReactedItem.File.Name {
				t.Errorf("%d: Got file name %#v, want %#v", i, got.File.Name, test.wantReactedItem.File.Name)
			}
			if got.Comment.Comment != test.wantReactedItem.Comment.Comment {
				t.Errorf("%d: Got comment comment %#v, want %#v", i, got.Comment.Comment, test.wantReactedItem.Comment.Comment)
			}
		}
	}
}

func TestSlack_ListReactions(t *testing.T) {
	once.Do(startServer)
	api := New("testing-token", OptionAPIURL("http://"+serverAddr+"/"))
	rh := newReactionsHandler()
	http.HandleFunc("/reactions.list", func(w http.ResponseWriter, r *http.Request) { rh.handler(w, r) })
	rh.response = `{"ok": true,
    "items": [
        {
            "type": "message",
            "channel": "C1",
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
            "file": {
                "name": "toy"
            },
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
    ],
    "paging": {
        "count": 100,
        "total": 4,
        "page": 1,
        "pages": 1
    }}`
	want := []ReactedItem{
		{
			Item: NewMessageItem("C1", &Message{Msg: Msg{
				Text: "hello",
				Reactions: []ItemReaction{
					{Name: "astonished", Count: 3, Users: []string{"U1", "U2", "U3"}},
					{Name: "clock1", Count: 3, Users: []string{"U1", "U2"}},
				},
			}}),
			Reactions: []ItemReaction{
				{Name: "astonished", Count: 3, Users: []string{"U1", "U2", "U3"}},
				{Name: "clock1", Count: 3, Users: []string{"U1", "U2"}},
			},
		},
		{
			Item: NewFileItem(&File{Name: "toy"}),
			Reactions: []ItemReaction{
				{Name: "clock1", Count: 3, Users: []string{"U1", "U2"}},
			},
		},
		{
			Item: NewFileCommentItem(&File{Name: "toy"}, &Comment{Comment: "cool toy"}),
			Reactions: []ItemReaction{
				{Name: "astonished", Count: 3, Users: []string{"U1", "U2", "U3"}},
			},
		},
	}
	wantParams := map[string]string{
		"user":  "User",
		"count": "200",
		"page":  "2",
		"full":  "true",
	}
	params := NewListReactionsParameters()
	params.User = "User"
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
	if !reflect.DeepEqual(rh.gotParams, wantParams) {
		t.Errorf("Got params %#v, want %#v", rh.gotParams, wantParams)
	}
	if reflect.DeepEqual(paging, Paging{}) {
		t.Errorf("Want paging data, got empty struct")
	}
}
