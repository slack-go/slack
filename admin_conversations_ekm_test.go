package slack

import (
	"context"
	"io"
	"net/http"
	"slices"
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
	if got := response.Channels[0].InternalTeamIDs; !slices.Equal(got, []string{"T001", "T002"}) {
		t.Errorf("unexpected internal team IDs: %v", got)
	}
}

func mockEKMListHandler(t *testing.T) func(rw http.ResponseWriter, r *http.Request) {
	return func(rw http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST request, got %s", r.Method)
		}

		rw.Header().Set("Content-Type", "application/json")
		_, _ = io.WriteString(rw, `{"ok":true,"channels":[{
			"id":"C1234567890",
			"original_connected_host_id":"T123",
			"original_connected_channel_id":"C9876543210",
			"internal_team_ids":["T001","T002"]
		}]}`)
	}
}
