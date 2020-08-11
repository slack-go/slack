package slack

import (
	"net/http"
	"testing"
)

func getBotInfo(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	response := []byte(`{"ok": true, "bot": {
			"id":"B02875YLA",
			"deleted":false,
			"name":"github",
			"updated": 1449272004,
			"app_id":"A161CLERW",
			"user_id": "U012ABCDEF",
			"icons": {
              "image_36":"https:\/\/a.slack-edge.com\/2fac\/plugins\/github\/assets\/service_36.png",
              "image_48":"https:\/\/a.slack-edge.com\/2fac\/plugins\/github\/assets\/service_48.png",
              "image_72":"https:\/\/a.slack-edge.com\/2fac\/plugins\/github\/assets\/service_72.png"
            }
        }}`)
	rw.Write(response)
}

func TestGetBotInfo(t *testing.T) {
	http.HandleFunc("/bots.info", getBotInfo)

	once.Do(startServer)
	api := New("testing-token", OptionAPIURL("http://"+serverAddr+"/"))

	bot, err := api.GetBotInfo("B02875YLA")
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
		return
	}

	if bot.ID != "B02875YLA" {
		t.Fatal("Incorrect ID")
	}
	if bot.Name != "github" {
		t.Fatal("Incorrect Name")
	}
	if bot.AppID != "A161CLERW" {
		t.Fatal("Incorrect App ID")
	}
	if bot.UserID != "U012ABCDEF" {
		t.Fatal("Incorrect User ID")
	}
	if bot.Updated != 1449272004 {
		t.Fatal("Incorrect Updated")
	}
	if len(bot.Icons.Image36) == 0 {
		t.Fatal("Missing Image36")
	}
	if len(bot.Icons.Image48) == 0 {
		t.Fatal("Missing Image38")
	}
	if len(bot.Icons.Image72) == 0 {
		t.Fatal("Missing Image72")
	}
}

func getBotInfoSelf(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	response := []byte(`{"ok":true,
		"url":"cerberus-xxxx.lb.slack-msgs.com",
		"team":{
			"id":"T1234567","name":"Github","domain":"github"
		},"self":{
			"id":"U012ABCDEF","name":"github"}}`)
	rw.Write(response)
}

func TestGetBotInfoSelf(t *testing.T) {
	http.DefaultServeMux = new(http.ServeMux)
	http.HandleFunc("/bots.info", getBotInfoSelf)

	once.Do(startServer)
	api := New("testing-token", OptionAPIURL("http://"+serverAddr+"/"))

	bot, err := api.GetBotInfoSelf()
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
		return
	}

	if bot.URL != "cerberus-xxxx.lb.slack-msgs.com" {
		t.Fatal("Missing URL")
	}

	if bot.Team == nil {
		t.Fatal("Missing Team")
	}

	if bot.Team.ID != "T1234567" {
		t.Fatal("Missing Team ID")
	}

	if bot.Team.Name != "Github" {
		t.Fatal("Missing Team Name")
	}

	if bot.Team.Domain != "github" {
		t.Fatal("Missing Team Domain")
	}

	if bot.User == nil {
		t.Fatal("Missing User")
	}

	if bot.User.ID != "U012ABCDEF" {
		t.Fatal("Missing User ID")
	}

	if bot.User.Name != "github" {
		t.Fatal("Missing User Name")
	}
}
