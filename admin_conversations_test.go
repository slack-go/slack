package slack

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
)

func TestAdminConversationsInvite(t *testing.T) {
	http.DefaultServeMux = new(http.ServeMux)
	http.HandleFunc("/admin.conversations.invite", mockAdminInviteHandler(t))
	once.Do(startServer)
	api := New("testing-token", OptionAPIURL("http://"+serverAddr+"/"))

	err := api.AdminConversationsInvite(context.Background(), AdminConversationsInviteParams{
		ChannelID: "C1234567890",
		UserIDs:   []string{"U123", "U456"},
	})
	if err != nil {
		t.Errorf("unexpected error: %s", err)
		return
	}
}

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

func TestAdminConversationsArchive(t *testing.T) {
	http.DefaultServeMux = new(http.ServeMux)
	http.HandleFunc("/admin.conversations.archive", mockAdminChannelIDHandler(t))
	once.Do(startServer)
	api := New("testing-token", OptionAPIURL("http://"+serverAddr+"/"))

	err := api.AdminConversationsArchive(context.Background(), "C1234567890")
	if err != nil {
		t.Errorf("unexpected error: %s", err)
		return
	}
}

func TestAdminConversationsUnarchive(t *testing.T) {
	http.DefaultServeMux = new(http.ServeMux)
	http.HandleFunc("/admin.conversations.unarchive", mockAdminChannelIDHandler(t))
	once.Do(startServer)
	api := New("testing-token", OptionAPIURL("http://"+serverAddr+"/"))

	err := api.AdminConversationsUnarchive(context.Background(), "C1234567890")
	if err != nil {
		t.Errorf("unexpected error: %s", err)
		return
	}
}

func TestAdminConversationsRename(t *testing.T) {
	http.DefaultServeMux = new(http.ServeMux)
	http.HandleFunc("/admin.conversations.rename", mockAdminRenameHandler(t))
	once.Do(startServer)
	api := New("testing-token", OptionAPIURL("http://"+serverAddr+"/"))

	err := api.AdminConversationsRename(context.Background(), "C1234567890", "new-channel-name")
	if err != nil {
		t.Errorf("unexpected error: %s", err)
		return
	}
}

func TestAdminConversationsDelete(t *testing.T) {
	http.DefaultServeMux = new(http.ServeMux)
	http.HandleFunc("/admin.conversations.delete", mockAdminChannelIDHandler(t))
	once.Do(startServer)
	api := New("testing-token", OptionAPIURL("http://"+serverAddr+"/"))

	err := api.AdminConversationsDelete(context.Background(), "C1234567890")
	if err != nil {
		t.Errorf("unexpected error: %s", err)
		return
	}
}

func TestAdminConversationsCreate(t *testing.T) {
	http.DefaultServeMux = new(http.ServeMux)
	http.HandleFunc("/admin.conversations.create", mockAdminCreateHandler(t))
	once.Do(startServer)
	api := New("testing-token", OptionAPIURL("http://"+serverAddr+"/"))

	channelID, err := api.AdminConversationsCreate(context.Background(), "test-channel", true,
		AdminConversationsCreateOptionDescription("A test channel"),
		AdminConversationsCreateOptionTeamID("T123"),
	)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
		return
	}

	if channelID != "C1234567890" {
		t.Errorf("unexpected channel_id: %s", channelID)
		return
	}
}

func TestAdminConversationsGetTeams(t *testing.T) {
	http.DefaultServeMux = new(http.ServeMux)
	http.HandleFunc("/admin.conversations.getTeams", mockAdminGetTeamsHandler(t))
	once.Do(startServer)
	api := New("testing-token", OptionAPIURL("http://"+serverAddr+"/"))

	teamIDs, cursor, err := api.AdminConversationsGetTeams(context.Background(), AdminConversationsGetTeamsParams{
		ChannelID: "C1234567890",
		Limit:     100,
	})
	if err != nil {
		t.Errorf("unexpected error: %s", err)
		return
	}

	if len(teamIDs) != 2 {
		t.Errorf("unexpected team count: %d", len(teamIDs))
		return
	}

	if cursor != "next_cursor_value" {
		t.Errorf("unexpected cursor: %s", cursor)
		return
	}
}

func TestAdminConversationsSearch(t *testing.T) {
	http.DefaultServeMux = new(http.ServeMux)
	http.HandleFunc("/admin.conversations.search", mockAdminSearchHandler(t))
	once.Do(startServer)
	api := New("testing-token", OptionAPIURL("http://"+serverAddr+"/"))

	response, err := api.AdminConversationsSearch(context.Background(),
		AdminConversationsSearchOptionQuery("test"),
		AdminConversationsSearchOptionLimit(100),
	)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
		return
	}

	if len(response.Conversations) != 1 {
		t.Errorf("unexpected conversation count: %d", len(response.Conversations))
		return
	}

	if response.Conversations[0].ID != "C1234567890" {
		t.Errorf("unexpected conversation ID: %s", response.Conversations[0].ID)
		return
	}
}

func TestAdminConversationsBulkArchive(t *testing.T) {
	http.DefaultServeMux = new(http.ServeMux)
	http.HandleFunc("/admin.conversations.bulkArchive", mockAdminChannelIDsHandler(t))
	once.Do(startServer)
	api := New("testing-token", OptionAPIURL("http://"+serverAddr+"/"))

	err := api.AdminConversationsBulkArchive(context.Background(), []string{"C123", "C456"})
	if err != nil {
		t.Errorf("unexpected error: %s", err)
		return
	}
}

func TestAdminConversationsBulkDelete(t *testing.T) {
	http.DefaultServeMux = new(http.ServeMux)
	http.HandleFunc("/admin.conversations.bulkDelete", mockAdminChannelIDsHandler(t))
	once.Do(startServer)
	api := New("testing-token", OptionAPIURL("http://"+serverAddr+"/"))

	err := api.AdminConversationsBulkDelete(context.Background(), []string{"C123", "C456"})
	if err != nil {
		t.Errorf("unexpected error: %s", err)
		return
	}
}

func TestAdminConversationsBulkMove(t *testing.T) {
	http.DefaultServeMux = new(http.ServeMux)
	http.HandleFunc("/admin.conversations.bulkMove", mockAdminBulkMoveHandler(t))
	once.Do(startServer)
	api := New("testing-token", OptionAPIURL("http://"+serverAddr+"/"))

	err := api.AdminConversationsBulkMove(context.Background(), AdminConversationsBulkMoveParams{
		ChannelIDs:   []string{"C123", "C456"},
		TargetTeamID: "T789",
	})
	if err != nil {
		t.Errorf("unexpected error: %s", err)
		return
	}
}

func TestAdminConversationsGetConversationPrefs(t *testing.T) {
	http.DefaultServeMux = new(http.ServeMux)
	http.HandleFunc("/admin.conversations.getConversationPrefs", mockAdminGetConversationPrefsHandler(t))
	once.Do(startServer)
	api := New("testing-token", OptionAPIURL("http://"+serverAddr+"/"))

	prefs, err := api.AdminConversationsGetConversationPrefs(context.Background(), "C1234567890")
	if err != nil {
		t.Errorf("unexpected error: %s", err)
		return
	}

	if prefs.WhoCanPost == nil || len(prefs.WhoCanPost.Type) != 1 || prefs.WhoCanPost.Type[0] != "admin" {
		t.Errorf("unexpected prefs: %+v", prefs)
		return
	}
}

func TestAdminConversationsSetCustomRetention(t *testing.T) {
	http.DefaultServeMux = new(http.ServeMux)
	http.HandleFunc("/admin.conversations.setCustomRetention", mockAdminRetentionHandler(t))
	once.Do(startServer)
	api := New("testing-token", OptionAPIURL("http://"+serverAddr+"/"))

	err := api.AdminConversationsSetCustomRetention(context.Background(), "C1234567890", 90)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
		return
	}
}

func TestAdminConversationsRemoveCustomRetention(t *testing.T) {
	http.DefaultServeMux = new(http.ServeMux)
	http.HandleFunc("/admin.conversations.removeCustomRetention", mockAdminChannelIDHandler(t))
	once.Do(startServer)
	api := New("testing-token", OptionAPIURL("http://"+serverAddr+"/"))

	err := api.AdminConversationsRemoveCustomRetention(context.Background(), "C1234567890")
	if err != nil {
		t.Errorf("unexpected error: %s", err)
		return
	}
}

func TestAdminConversationsDisconnectShared(t *testing.T) {
	http.DefaultServeMux = new(http.ServeMux)
	http.HandleFunc("/admin.conversations.disconnectShared", mockAdminChannelIDHandler(t))
	once.Do(startServer)
	api := New("testing-token", OptionAPIURL("http://"+serverAddr+"/"))

	err := api.AdminConversationsDisconnectShared(context.Background(), "C1234567890",
		AdminConversationsDisconnectSharedOptionLeavingTeamIDs([]string{"T123", "T456"}),
	)
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

		_, _ = rw.Write(response)
	}
}

func mockAdminInviteHandler(t *testing.T) func(rw http.ResponseWriter, r *http.Request) {
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

		if len(r.Form["user_ids"]) == 0 {
			t.Error("missing user_ids in request")
			return
		}

		rw.Header().Set("Content-Type", "application/json")
		response, _ := json.Marshal(SlackResponse{
			Ok: true,
		})

		_, _ = rw.Write(response)
	}
}

func mockAdminRenameHandler(t *testing.T) func(rw http.ResponseWriter, r *http.Request) {
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

		if len(r.Form["name"]) == 0 {
			t.Error("missing name in request")
			return
		}

		rw.Header().Set("Content-Type", "application/json")
		response, _ := json.Marshal(SlackResponse{
			Ok: true,
		})

		_, _ = rw.Write(response)
	}
}

func mockAdminCreateHandler(t *testing.T) func(rw http.ResponseWriter, r *http.Request) {
	return func(rw http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST request, got %s", r.Method)
		}

		if err := r.ParseForm(); err != nil {
			t.Errorf("unexpected error: %s", err)
			return
		}

		if len(r.Form["name"]) == 0 {
			t.Error("missing name in request")
			return
		}

		if len(r.Form["is_private"]) == 0 {
			t.Error("missing is_private in request")
			return
		}

		rw.Header().Set("Content-Type", "application/json")
		response, _ := json.Marshal(AdminConversationsCreateResponse{
			SlackResponse: SlackResponse{Ok: true},
			ChannelID:     "C1234567890",
		})

		_, _ = rw.Write(response)
	}
}

func mockAdminGetTeamsHandler(t *testing.T) func(rw http.ResponseWriter, r *http.Request) {
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
		response, _ := json.Marshal(AdminConversationsGetTeamsResponse{
			SlackResponse: SlackResponse{
				Ok:               true,
				ResponseMetadata: ResponseMetadata{Cursor: "next_cursor_value"},
			},
			TeamIDs: []string{"T123", "T456"},
		})

		_, _ = rw.Write(response)
	}
}

func mockAdminSearchHandler(t *testing.T) func(rw http.ResponseWriter, r *http.Request) {
	return func(rw http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST request, got %s", r.Method)
		}

		rw.Header().Set("Content-Type", "application/json")
		response, _ := json.Marshal(AdminConversationsSearchResponse{
			SlackResponse: SlackResponse{Ok: true},
			Conversations: []AdminConversation{
				{
					ID:          "C1234567890",
					Name:        "test-channel",
					IsPrivate:   false,
					MemberCount: 10,
				},
			},
			TotalCount: 1,
		})

		_, _ = rw.Write(response)
	}
}

func mockAdminChannelIDsHandler(t *testing.T) func(rw http.ResponseWriter, r *http.Request) {
	return func(rw http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST request, got %s", r.Method)
		}

		if err := r.ParseForm(); err != nil {
			t.Errorf("unexpected error: %s", err)
			return
		}

		if len(r.Form["channel_ids"]) == 0 {
			t.Error("missing channel_ids in request")
			return
		}

		rw.Header().Set("Content-Type", "application/json")
		response, _ := json.Marshal(SlackResponse{
			Ok: true,
		})

		_, _ = rw.Write(response)
	}
}

func mockAdminBulkMoveHandler(t *testing.T) func(rw http.ResponseWriter, r *http.Request) {
	return func(rw http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST request, got %s", r.Method)
		}

		if err := r.ParseForm(); err != nil {
			t.Errorf("unexpected error: %s", err)
			return
		}

		if len(r.Form["channel_ids"]) == 0 {
			t.Error("missing channel_ids in request")
			return
		}

		if len(r.Form["target_team_id"]) == 0 {
			t.Error("missing target_team_id in request")
			return
		}

		rw.Header().Set("Content-Type", "application/json")
		response, _ := json.Marshal(SlackResponse{
			Ok: true,
		})

		_, _ = rw.Write(response)
	}
}

func mockAdminGetConversationPrefsHandler(t *testing.T) func(rw http.ResponseWriter, r *http.Request) {
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
		response, _ := json.Marshal(AdminConversationsGetConversationPrefsResponse{
			SlackResponse: SlackResponse{Ok: true},
			Prefs: AdminConversationPrefs{
				WhoCanPost: &AdminConversationPref{Type: []string{"admin"}},
			},
		})

		_, _ = rw.Write(response)
	}
}

func mockAdminRetentionHandler(t *testing.T) func(rw http.ResponseWriter, r *http.Request) {
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

		if len(r.Form["duration_days"]) == 0 {
			t.Error("missing duration_days in request")
			return
		}

		rw.Header().Set("Content-Type", "application/json")
		response, _ := json.Marshal(SlackResponse{
			Ok: true,
		})

		_, _ = rw.Write(response)
	}
}
