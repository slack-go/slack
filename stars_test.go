package slack

import (
	"net/http"
	"reflect"
	"testing"
)

type starsHandler struct {
	gotParams map[string]string
	response  string
}

func newStarsHandler() *starsHandler {
	return &starsHandler{
		gotParams: make(map[string]string),
		response:  `{ "ok": true }`,
	}
}

func (sh *starsHandler) accumulateFormValue(k string, r *http.Request) {
	if v := r.FormValue(k); v != "" {
		sh.gotParams[k] = v
	}
}

func (sh *starsHandler) handler(w http.ResponseWriter, r *http.Request) {
	sh.accumulateFormValue("user", r)
	sh.accumulateFormValue("count", r)
	sh.accumulateFormValue("page", r)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(sh.response))
}

func TestSlack_GetStarred(t *testing.T) {
	once.Do(startServer)
	SLACK_API = "http://" + serverAddr + "/"
	api := New("testing-token")
	tests := []struct {
		params       StarsParameters
		wantParams   map[string]string
		json         string
		starredItems []StarredItem
		paging       *Paging
	}{
		{
			StarsParameters{
				User:  "U1",
				Count: 10,
				Page:  100,
			},
			map[string]string{
				"user":  "U1",
				"count": "10",
				"page":  "100",
			},
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
		sh = newStarsHandler()
		sh.response = test.json
		responseItems, responsePaging, err := api.GetStarred(test.params)
		if err != nil {
			t.Fatalf("%d Unexpected error: %s", i, err)
		}
		if !reflect.DeepEqual(sh.gotParams, test.wantParams) {
			t.Errorf("%d got %v; want %v", i, sh.gotParams, test.wantParams)
		}
		if !reflect.DeepEqual(responseItems, test.starredItems) {
			t.Errorf("%d got %v; want %v", i, responseItems, test.starredItems)
		}
		if !reflect.DeepEqual(responsePaging, test.paging) {
			t.Errorf("%d got %v; want %v", i, responsePaging, test.paging)
		}
	}
}
