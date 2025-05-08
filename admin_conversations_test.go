package slack

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
)

func TestAdminConversationsSetTeams(t *testing.T) {
	http.DefaultServeMux = new(http.ServeMux)
	http.HandleFunc("/admin.conversations.setTeams", mockAdminChannelIDHandler(t))
	once.Do(startServer)
	api := New("testing-token", OptionAPIURL("http://"+serverAddr+"/"))

	orgChannel := true
	teamID := "T789"

	err := api.AdminConversationsSetTeams(context.Background(), AdminConversationsSetTeamsParams{
		ChannelID:     "C1234567890",
		OrgChannel:    &orgChannel,
		TargetTeamIDs: []string{"T123", "T456"},
		TeamID:        &teamID,
	})
	if err != nil {
		t.Errorf("unexpected error: %s", err)
		return
	}
}

func TestAdminConversationsConvertToPrivate(t *testing.T) {
	http.DefaultServeMux = new(http.ServeMux)
	http.HandleFunc("/admin.conversations.convertToPrivate", mockAdminChannelIDHandler(t))
	once.Do(startServer)
	api := New("testing-token", OptionAPIURL("http://"+serverAddr+"/"))

	err := api.AdminConversationsConvertToPrivate(context.Background(), "C1234567890")
	if err != nil {
		t.Errorf("unexpected error: %s", err)
		return
	}
}

func TestAdminConversationsConvertToPublic(t *testing.T) {
	http.DefaultServeMux = new(http.ServeMux)
	http.HandleFunc("/admin.conversations.convertToPublic", mockAdminChannelIDHandler(t))
	once.Do(startServer)
	api := New("testing-token", OptionAPIURL("http://"+serverAddr+"/"))

	err := api.AdminConversationsConvertToPublic(context.Background(), "C1234567890")
	if err != nil {
		t.Errorf("unexpected error: %s", err)
		return
	}
}

// mockAdminChannelIDHandler returns a handler which expects a channel_id to be present
// in the request, and will fail the test if it isn't present.
func mockAdminChannelIDHandler(t *testing.T) func(rw http.ResponseWriter, r *http.Request) {
	return func(rw http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST request, got %s", r.Method)
		}

		if err := r.ParseForm(); err != nil {
			t.Errorf("unexpected error: %s", err)
			return
		}

		if len(r.Form["channel_id"]) == 0 {
			t.Error("missing channel_id in request")
			return
		}

		rw.Header().Set("Content-Type", "application/json")
		response, _ := json.Marshal(SlackResponse{
			Ok: true,
		})

		rw.Write(response)
	}
}
