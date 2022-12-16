package slack

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func getTeamList(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	response := []byte(`{
    "ok": true,
    "teams": [
        {
            "name": "Shinichi's workspace",
            "id": "T12345678"
        },
        {
            "name": "Migi's workspace",
            "id": "T12345679"
        }
    ],
    "response_metadata": {
        "next_cursor": "dXNlcl9pZDo5MTQyOTI5Mzkz"
    }
}`)
	rw.Write(response)
}

func TestListTeams(t *testing.T) {
	http.HandleFunc("/auth.teams.list", getTeamList)

	once.Do(startServer)
	api := New("testing-token", OptionAPIURL("http://"+serverAddr+"/"))

	teams, cursor, err := api.ListTeams(ListTeamsParameters{})
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
		return
	}

	assert.Len(t, teams, 2)
	assert.Equal(t, "T12345678", teams[0].ID)
	assert.Equal(t, "Shinichi's workspace", teams[0].Name)

	assert.Equal(t, "T12345679", teams[1].ID)
	assert.Equal(t, "Migi's workspace", teams[1].Name)

	assert.Equal(t, "dXNlcl9pZDo5MTQyOTI5Mzkz", cursor)
}
