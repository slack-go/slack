package slack

import (
	"errors"
	"net/http"
	"testing"
)

var (
	ErrIncorrectResponse = errors.New("Response is incorrect")
)

func getTeamInfo(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	response := []byte(`{"ok": true, "team": {
			"id": "F0UWHUX",
			"name": "notalar",
			"domain": "notalar",
			"icon": {
              "image_34": "https://slack.global.ssl.fastly.net/66f9/img/avatars-teams/ava_0002-34.png",
              "image_44": "https://slack.global.ssl.fastly.net/66f9/img/avatars-teams/ava_0002-44.png",
              "image_55": "https://slack.global.ssl.fastly.net/66f9/img/avatars-teams/ava_0002-55.png",
              "image_default": true
          }
		}}`)
	rw.Write(response)
}

func TestGetTeamInfo(t *testing.T) {
	http.HandleFunc("/team.info", getTeamInfo)

	once.Do(startServer)
	SLACK_API = "http://" + serverAddr + "/"
	api := New("testing-token")

	teamInfo, err := api.GetTeamInfo()
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
		return
	}

	// t.Fatal refers to -> t.Errorf & return
	if teamInfo.ID != "F0UWHUX" {
		t.Fatal(ErrIncorrectResponse)
	}
	if teamInfo.Domain != "notalar" {
		t.Fatal(ErrIncorrectResponse)
	}
	if teamInfo.Name != "notalar" {
		t.Fatal(ErrIncorrectResponse)
	}
	if teamInfo.Icon == nil {
		t.Fatal(ErrIncorrectResponse)
	}
}
