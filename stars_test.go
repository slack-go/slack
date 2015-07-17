package slack

import (
	"net/http"
	"reflect"
	"testing"
)

type starsHandler struct {
	response string
}

func (rh *starsHandler) handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(rh.response))
}

func TestSlack_GetStarred(t *testing.T) {
	once.Do(startServer)
	SLACK_API = "http://" + serverAddr + "/"
	api := New("testing-token")
	tests := []struct {
		json         string
		starredItems []StarredItem
		paging       *Paging
	}{
		{
			`{"ok": true,
    "items": [
        {
            "type": "message",
            "channel": "C2147483705",
            "message": {
                "text": "hello"
            }
        },
        {
            "type": "file",
            "file": {
                "name": "toy"
            }
        },
        {
            "type": "file_comment",
            "file": {
                "name": "toy"
            },
            "comment": {
                "comment": "nice"
            }
        },
        {
            "type": "channel",
            "channel": "C2147483705"
        },
        {
            "type": "im",
            "channel": "D1"
        },
        {
            "type": "group",
            "channel": "G1"
        }
    ],
    "paging": {
        "count": 100,
        "total": 4,
        "page": 1,
        "pages": 1
    }}`,
			[]StarredItem{
				{Item: NewMessageItem("C2147483705", &Message{Msg: Msg{Text: "hello"}})},
				{Item: NewFileItem(&File{Name: "toy"})},
				{Item: NewFileCommentItem(&File{Name: "toy"}, &Comment{Comment: "nice"})},
				{Item: NewChannelItem("C2147483705")},
				{Item: NewIMItem("D1")},
				{Item: NewGroupItem("G1")},
			},
			&Paging{
				Count: 100,
				Total: 4,
				Page:  1,
				Pages: 1,
			},
		},
	}
	var sh *starsHandler
	http.HandleFunc("/stars.list", func(w http.ResponseWriter, r *http.Request) { sh.handler(w, r) })
	for i, test := range tests {
		sh = &starsHandler{}
		sh.response = test.json
		response_items, response_paging, err := api.GetStarred(NewStarsParameters())
		if err != nil {
			t.Fatalf("%d Unexpected error: %s", i, err)
		}
		if !reflect.DeepEqual(response_items, test.starredItems) {
			t.Errorf("%d got %v; want %v", i, response_items, test.starredItems)
		}
		if !reflect.DeepEqual(response_paging, test.paging) {
			t.Errorf("%d got %v; want %v", i, response_paging, test.paging)
		}
	}
}
