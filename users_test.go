package slack

import (
	"net/http"
	"testing"
)

func getUserIdentity(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	response := []byte(`{
  "ok": true,
  "user": {
    "id": "UXXXXXXXX",
    "name": "Test User",
    "email": "test@test.com",
    "image_24": "https:\/\/s3-us-west-2.amazonaws.com\/slack-files2\/avatars\/2016-10-18\/92962080834_ef14c1469fc0741caea1_24.jpg",
    "image_32": "https:\/\/s3-us-west-2.amazonaws.com\/slack-files2\/avatars\/2016-10-18\/92962080834_ef14c1469fc0741caea1_32.jpg",
    "image_48": "https:\/\/s3-us-west-2.amazonaws.com\/slack-files2\/avatars\/2016-10-18\/92962080834_ef14c1469fc0741caea1_48.jpg",
    "image_72": "https:\/\/s3-us-west-2.amazonaws.com\/slack-files2\/avatars\/2016-10-18\/92962080834_ef14c1469fc0741caea1_72.jpg",
    "image_192": "https:\/\/s3-us-west-2.amazonaws.com\/slack-files2\/avatars\/2016-10-18\/92962080834_ef14c1469fc0741caea1_192.jpg",
    "image_512": "https:\/\/s3-us-west-2.amazonaws.com\/slack-files2\/avatars\/2016-10-18\/92962080834_ef14c1469fc0741caea1_512.jpg"
  },
  "team": {
    "id": "TXXXXXXXX",
    "name": "team-name",
    "domain": "team-domain",
    "image_34": "https:\/\/s3-us-west-2.amazonaws.com\/slack-files2\/avatars\/2016-10-18\/92962080834_ef14c1469fc0741caea1_34.jpg",
    "image_44": "https:\/\/s3-us-west-2.amazonaws.com\/slack-files2\/avatars\/2016-10-18\/92962080834_ef14c1469fc0741caea1_44.jpg",
    "image_68": "https:\/\/s3-us-west-2.amazonaws.com\/slack-files2\/avatars\/2016-10-18\/92962080834_ef14c1469fc0741caea1_68.jpg",
    "image_88": "https:\/\/s3-us-west-2.amazonaws.com\/slack-files2\/avatars\/2016-10-18\/92962080834_ef14c1469fc0741caea1_88.jpg",
    "image_102": "https:\/\/s3-us-west-2.amazonaws.com\/slack-files2\/avatars\/2016-10-18\/92962080834_ef14c1469fc0741caea1_102.jpg",
    "image_132": "https:\/\/s3-us-west-2.amazonaws.com\/slack-files2\/avatars\/2016-10-18\/92962080834_ef14c1469fc0741caea1_132.jpg",
    "image_230": "https:\/\/s3-us-west-2.amazonaws.com\/slack-files2\/avatars\/2016-10-18\/92962080834_ef14c1469fc0741caea1_230.jpg",
    "image_original": "https:\/\/s3-us-west-2.amazonaws.com\/slack-files2\/avatars\/2016-10-18\/92962080834_ef14c1469fc0741caea1_original.jpg"
  }
}`)
	rw.Write(response)
}

func TestGetUserIdentity(t *testing.T) {
	http.HandleFunc("/users.identity", getUserIdentity)

	once.Do(startServer)
	SLACK_API = "http://" + serverAddr + "/"
	api := New("testing-token")

	identity, err := api.GetUserIdentity()
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
		return
	}

	// t.Fatal refers to -> t.Errorf & return
	if identity.User.ID != "UXXXXXXXX" {
		t.Fatal(ErrIncorrectResponse)
	}
	if identity.User.Name != "Test User" {
		t.Fatal(ErrIncorrectResponse)
	}
	if identity.User.Email != "test@test.com" {
		t.Fatal(ErrIncorrectResponse)
	}
	if identity.Team.ID != "TXXXXXXXX" {
		t.Fatal(ErrIncorrectResponse)
	}
	if identity.Team.Name != "team-name" {
		t.Fatal(ErrIncorrectResponse)
	}
	if identity.Team.Domain != "team-domain" {
		t.Fatal(ErrIncorrectResponse)
	}
	if identity.User.Image24 == "" {
		t.Fatal(ErrIncorrectResponse)
	}
	if identity.Team.Image34 == "" {
		t.Fatal(ErrIncorrectResponse)
	}
}
