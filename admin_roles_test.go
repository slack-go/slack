package slack

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
)

func TestAdminRolesAddAssignments(t *testing.T) {
	http.DefaultServeMux = new(http.ServeMux)
	http.HandleFunc("/admin.roles.addAssignments", mockAdminRolesAddAssignmentsHandler(t))
	once.Do(startServer)
	api := New("testing-token", OptionAPIURL("http://"+serverAddr+"/"))

	resp, err := api.AdminRolesAddAssignments(context.Background(), AdminRolesAddAssignmentsParams{
		RoleID:    "Rl0L",
		UserIDs:   []string{"U123", "U456"},
		EntityIDs: []string{"E123"},
	})
	if err != nil {
		t.Errorf("unexpected error: %s", err)
		return
	}

	if !resp.Ok {
		t.Errorf("expected Ok to be true")
	}
}

func mockAdminRolesAddAssignmentsHandler(t *testing.T) func(rw http.ResponseWriter, r *http.Request) {
	return func(rw http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST request, got %s", r.Method)
		}

		if err := r.ParseForm(); err != nil {
			t.Errorf("unexpected error: %s", err)
			return
		}

		if len(r.Form["role_id"]) == 0 {
			t.Error("missing role_id in request")
			return
		}

		if len(r.Form["user_ids"]) == 0 && len(r.Form["entity_ids"]) == 0 {
			t.Error("missing user_ids or entity_ids in request")
			return
		}

		rw.Header().Set("Content-Type", "application/json")
		response, _ := json.Marshal(AdminRolesAddAssignmentsResponse{
			SlackResponse: SlackResponse{Ok: true},
		})

		_, _ = rw.Write(response)
	}
}

func TestAdminRolesListAssignments(t *testing.T) {
	http.DefaultServeMux = new(http.ServeMux)
	http.HandleFunc("/admin.roles.listAssignments", mockAdminRolesListAssignmentsHandler(t))
	once.Do(startServer)
	api := New("testing-token", OptionAPIURL("http://"+serverAddr+"/"))

	resp, err := api.AdminRolesListAssignments(context.Background(), AdminRolesListAssignmentsParams{
		RoleIDs: []string{"Rl0L"},
		Limit:   10,
	})
	if err != nil {
		t.Errorf("unexpected error: %s", err)
		return
	}

	if !resp.Ok {
		t.Errorf("expected Ok to be true")
	}

	if len(resp.RoleAssignments) != 1 {
		t.Errorf("expected 1 role assignment, got %d", len(resp.RoleAssignments))
	}

	if resp.RoleAssignments[0].RoleID != "Rl0L" {
		t.Errorf("expected role ID Rl0L, got %s", resp.RoleAssignments[0].RoleID)
	}
}

func mockAdminRolesListAssignmentsHandler(t *testing.T) func(rw http.ResponseWriter, r *http.Request) {
	return func(rw http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST request, got %s", r.Method)
		}

		if err := r.ParseForm(); err != nil {
			t.Errorf("unexpected error: %s", err)
			return
		}

		rw.Header().Set("Content-Type", "application/json")
		response, _ := json.Marshal(AdminRolesListAssignmentsResponse{
			SlackResponse: SlackResponse{Ok: true},
			RoleAssignments: []RoleAssignment{
				{
					RoleID: "Rl0L",
					UserID: "U123",
				},
			},
		})

		_, _ = rw.Write(response)
	}
}

func TestAdminRolesRemoveAssignments(t *testing.T) {
	http.DefaultServeMux = new(http.ServeMux)
	http.HandleFunc("/admin.roles.removeAssignments", mockAdminRolesRemoveAssignmentsHandler(t))
	once.Do(startServer)
	api := New("testing-token", OptionAPIURL("http://"+serverAddr+"/"))

	resp, err := api.AdminRolesRemoveAssignments(context.Background(), AdminRolesRemoveAssignmentsParams{
		RoleID:    "Rl0L",
		UserIDs:   []string{"U123", "U456"},
		EntityIDs: []string{"E123"},
	})
	if err != nil {
		t.Errorf("unexpected error: %s", err)
		return
	}

	if !resp.Ok {
		t.Errorf("expected Ok to be true")
	}
}

func mockAdminRolesRemoveAssignmentsHandler(t *testing.T) func(rw http.ResponseWriter, r *http.Request) {
	return func(rw http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST request, got %s", r.Method)
		}

		if err := r.ParseForm(); err != nil {
			t.Errorf("unexpected error: %s", err)
			return
		}

		if len(r.Form["role_id"]) == 0 {
			t.Error("missing role_id in request")
			return
		}

		if len(r.Form["user_ids"]) == 0 && len(r.Form["entity_ids"]) == 0 {
			t.Error("missing user_ids or entity_ids in request")
			return
		}

		rw.Header().Set("Content-Type", "application/json")
		response, _ := json.Marshal(AdminRolesRemoveAssignmentsResponse{
			SlackResponse: SlackResponse{Ok: true},
		})

		_, _ = rw.Write(response)
	}
}
