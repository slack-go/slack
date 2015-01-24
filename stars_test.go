package slack

import (
	"net/http"
	"reflect"
	"testing"
)

var starsTests = struct {
	in        []byte
	outItems  []StarredItem
	outPaging *Paging
}{
	[]byte(`{"ok": true,
    "items": [
        {
            "type": "message",
            "channel": "C2147483705"
        },
        {
            "type": "file"
        },
        {
            "type": "file_comment"
        },
        {
            "type": "channel",
            "channel": "C2147483705"
        }

    ],
    "paging": {
        "count": 100,
        "total": 4,
        "page": 1,
        "pages": 1
    }}`),
	[]StarredItem{
		{
			Type:      "message",
			ChannelId: "C2147483705",
		},
		{
			Type: "file",
		},
		{
			Type: "file_comment",
		},
		{
			Type:      "channel",
			ChannelId: "C2147483705",
		},
	},
	&Paging{
		Count: 100,
		Total: 4,
		Page:  1,
		Pages: 1,
	},
}

func getStarredHandler(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	// XXX: I stripped the actual content just to test this test Oo
	// At some point I should add valid content here
	rw.Write(starsTests.in)
}

func TestGetStarred(t *testing.T) {
	http.HandleFunc("/stars.list", getStarredHandler)
	once.Do(startServer)
	SLACK_API = "http://" + serverAddr + "/"
	api := New("testing-token")
	response_items, response_paging, err := api.GetStarred(NewStarsParameters())
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
		return
	}
	eq := reflect.DeepEqual(response_items, starsTests.outItems)
	if !eq {
		t.Errorf("got %v; want %v", response_items, starsTests.outItems)
	}
	eq = reflect.DeepEqual(response_paging, starsTests.outPaging)
	if !eq {
		t.Errorf("got %v; want %v", response_paging, starsTests.outPaging)
	}
}
