package slack

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
)

func TestAdminConversationsEKMListOriginalConnectedChannelInfo(t *testing.T) {
	http.DefaultServeMux = new(http.ServeMux)
	http.HandleFunc("/admin.conversations.ekm.listOriginalConnectedChannelInfo", mockEKMListHandler(t))
	once.Do(startServer)
	api := New("testing-token", OptionAPIURL("http://"+serverAddr+"/"))

	response, err := api.AdminConversationsEKMListOriginalConnectedChannelInfo(context.Background(),
		AdminConversationsEKMListOriginalConnectedChannelInfoOptionChannelIDs([]string{"C123", "C456"}),
		AdminConversationsEKMListOriginalConnectedChannelInfoOptionTeamIDs([]string{"T789"}),
		AdminConversationsEKMListOriginalConnectedChannelInfoOptionLimit(100),
	)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
		return
	}

	if len(response.Channels) != 1 {
		t.Errorf("unexpected channel count: %d", len(response.Channels))
		return
	}

	if response.Channels[0].ID != "C1234567890" {
		t.Errorf("unexpected channel ID: %s", response.Channels[0].ID)
		return
	}
}

func mockEKMListHandler(t *testing.T) func(rw http.ResponseWriter, r *http.Request) {
	return func(rw http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST request, got %s", r.Method)
		}

		rw.Header().Set("Content-Type", "application/json")
		response, _ := json.Marshal(AdminConversationsEKMListOriginalConnectedChannelInfoResponse{
			SlackResponse: SlackResponse{Ok: true},
			Channels: []AdminConversationsEKMOriginalConnectedChannelInfo{
				{
					ID:                         "C1234567890",
					OriginalConnectedHostID:    "T123",
					OriginalConnectedChannelID: "C9876543210",
					InternalTeamIDs:            []string{"T001", "T002"},
				},
			},
		})

		_, _ = rw.Write(response)
	}
}
