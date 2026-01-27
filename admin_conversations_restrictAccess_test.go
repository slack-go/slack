package slack

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
)

func TestAdminConversationsRestrictAccessAddGroup(t *testing.T) {
	http.DefaultServeMux = new(http.ServeMux)
	http.HandleFunc("/admin.conversations.restrictAccess.addGroup", mockRestrictAccessHandler(t, "group_id"))
	once.Do(startServer)
	api := New("testing-token", OptionAPIURL("http://"+serverAddr+"/"))

	err := api.AdminConversationsRestrictAccessAddGroup(context.Background(),
		"C1234567890",
		"G123",
		AdminConversationsRestrictAccessAddGroupOptionTeamID("T789"),
	)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
		return
	}
}

func TestAdminConversationsRestrictAccessListGroups(t *testing.T) {
	http.DefaultServeMux = new(http.ServeMux)
	http.HandleFunc("/admin.conversations.restrictAccess.listGroups", mockRestrictAccessListGroupsHandler(t))
	once.Do(startServer)
	api := New("testing-token", OptionAPIURL("http://"+serverAddr+"/"))

	groupIDs, err := api.AdminConversationsRestrictAccessListGroups(context.Background(),
		"C1234567890",
		AdminConversationsRestrictAccessListGroupsOptionTeamID("T789"),
	)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
		return
	}

	if len(groupIDs) != 2 {
		t.Errorf("unexpected group count: %d", len(groupIDs))
		return
	}
}

func TestAdminConversationsRestrictAccessRemoveGroup(t *testing.T) {
	http.DefaultServeMux = new(http.ServeMux)
	http.HandleFunc("/admin.conversations.restrictAccess.removeGroup", mockRestrictAccessHandler(t, "group_id"))
	once.Do(startServer)
	api := New("testing-token", OptionAPIURL("http://"+serverAddr+"/"))

	err := api.AdminConversationsRestrictAccessRemoveGroup(context.Background(),
		"C1234567890",
		"G123",
		AdminConversationsRestrictAccessRemoveGroupOptionTeamID("T789"),
	)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
		return
	}
}

func mockRestrictAccessHandler(t *testing.T, requiredField string) func(rw http.ResponseWriter, r *http.Request) {
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

		if len(r.Form[requiredField]) == 0 {
			t.Errorf("missing %s in request", requiredField)
			return
		}

		rw.Header().Set("Content-Type", "application/json")
		response, _ := json.Marshal(SlackResponse{
			Ok: true,
		})

		_, _ = rw.Write(response)
	}
}

func mockRestrictAccessListGroupsHandler(t *testing.T) func(rw http.ResponseWriter, r *http.Request) {
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
		response, _ := json.Marshal(AdminConversationsRestrictAccessListGroupsResponse{
			SlackResponse: SlackResponse{Ok: true},
			GroupIDs:      []string{"G123", "G456"},
		})

		_, _ = rw.Write(response)
	}
}
