package slack

import (
	"net/http"
	"reflect"
	"testing"
)

func getAdminEmojiHandler(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	response := []byte(`{"ok": true, "emoji": {
			"bowtie": {"url":"https://my.slack.com/emoji/bowtie/46ec6f2bb0.png","date_created":1644445892,"uploaded_by":"Churchill"},
			"squirrel": {"url":"https://my.slack.com/emoji/squirrel/f35f40c0e0.png","date_created":1644445893,"uploaded_by":"TheSquirrelHimself"},
			"shipit": {"url":"alias:squirrel","date_created":1644445894,"uploaded_by":"BoredGuy"}
		}}`)
	rw.Write(response)
}

func TestGetAdminEmoji(t *testing.T) {
	http.HandleFunc("/admin.emoji.list", getAdminEmojiHandler)

	once.Do(startServer)
	api := New("testing-token", OptionAPIURL("http://"+serverAddr+"/"))
	emojisResponse := map[string]*Emoji{
		"bowtie": {
			URL:         "https://my.slack.com/emoji/bowtie/46ec6f2bb0.png",
			DateCreated: 1644445892,
			UploadedBy:  "Churchill",
		},
		"squirrel": {
			URL:         "https://my.slack.com/emoji/squirrel/f35f40c0e0.png",
			DateCreated: 1644445893,
			UploadedBy:  "TheSquirrelHimself",
		},
		"shipit": {
			URL:         "alias:squirrel",
			DateCreated: 1644445894,
			UploadedBy:  "BoredGuy",
		},
	}

	emojis, err := api.GetAdminEmoji()
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
		return
	}
	eq := reflect.DeepEqual(emojis, emojisResponse)
	if !eq {
		t.Errorf("got %v; want %v", emojis, emojisResponse)
	}
}
