package slack

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func toolingTokensRotate(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	response := []byte(`{
	"ok": true,
	"token": "xoxe.xoxp-...",
	"refresh_token": "xoxe-...",
	"team_id": "...",
	"user_id": "...",
	"iat": 1633095660,
	"exp": 1633138860
}`)
	rw.Write(response)
}

func TestToolingTokensRotate(t *testing.T) {
	http.HandleFunc("/tooling.tokens.rotate", toolingTokensRotate)

	once.Do(startServer)
	api := New("testing-token", OptionAPIURL("http://"+serverAddr+"/"))

	token, err := api.ToolingTokensRotate("xoxe.xoxp-...")
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
		return
	}

	assert.Equal(t, true, token.Ok)
	assert.Equal(t, "", token.Error)

	assert.Equal(t, "xoxe.xoxp-...", token.Token)
	assert.Equal(t, "xoxe-...", token.RefreshToken)
	assert.Equal(t, "...", token.TeamID)
	assert.Equal(t, "...", token.UserID)
	assert.Equal(t, int64(1633095660), token.Iat)
	assert.Equal(t, int64(1633138860), token.Exp)
}
