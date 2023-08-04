package slack

import (
	"encoding/json"
	"net/http"
	"reflect"
	"testing"
)

func TestRotateTokens(t *testing.T) {
	http.HandleFunc("/tooling.tokens.rotate", handleRotateToken)
	expected := getTestTokenResponse()

	once.Do(startServer)
	api := New("testing-token", OptionAPIURL("http://"+serverAddr+"/"))

	tok, err := api.RotateTokens("expired-config", "old-refresh")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
		return
	}

	if !reflect.DeepEqual(expected, *tok) {
		t.Fatal(ErrIncorrectResponse)
	}
}

func getTestTokenResponse() TokenResponse {
	return TokenResponse{
		Token:         "token",
		RefreshToken:  "refresh",
		UserId:        "uid",
		TeamId:        "tid",
		IssuedAt:      1,
		ExpiresAt:     1,
		SlackResponse: SlackResponse{Ok: true},
	}
}

func handleRotateToken(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Content-Type", "application/json")

	response, _ := json.Marshal(getTestTokenResponse())
	rw.Write(response)
}
