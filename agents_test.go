package slack

import (
	"encoding/json"
	"net/http"
	"testing"
)

func agentSessionSetStatusHandler(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Content-Type", "application/json")

	if r.FormValue("status") == "" || r.FormValue("channel_id") == "" {
		rw.Write([]byte(`{ "ok": false, "error": "invalid_arguments" }`))
		return
	}

	resp, _ := json.Marshal(&AgentSessionSetStatusResponse{
		SlackResponse: SlackResponse{Ok: true},
		Status:        r.FormValue("status"),
		AgentStatus:   r.FormValue("status"),
		Title:         r.FormValue("title"),
	})
	rw.Write(resp)
}

func TestSetAgentSessionStatus(t *testing.T) {
	http.HandleFunc("/agents.sessions.setStatus", agentSessionSetStatusHandler)
	once.Do(startServer)
	api := New("testing-token", OptionAPIURL("http://"+serverAddr+"/"))

	params := AgentSessionSetStatusParameters{
		Status:          AgentSessionStatusProcessing,
		ChannelID:       "CXXXXXXXX",
		ThreadTS:        "1234567890.123456",
		Title:           "Scuba diving research",
		InitiatorUserID: "UXXXXXXXX",
		IconEmoji:       ":thinking_face:",
		IconURL:         "https://example.com/icon.png",
		Username:        "custom name",
	}

	response, err := api.SetAgentSessionStatus(params)
	if err != nil {
		t.Fatalf("Unexpected error: %s", err)
	}

	if !response.Ok {
		t.Fatalf("Expected Ok to be true")
	}

	if response.Status != "processing" {
		t.Fatalf("Expected Status to be 'processing', got %q", response.Status)
	}

	if response.AgentStatus != "processing" {
		t.Fatalf("Expected AgentStatus to be 'processing', got %q", response.AgentStatus)
	}

	if response.Title != "Scuba diving research" {
		t.Fatalf("Expected Title to be 'Scuba diving research', got %q", response.Title)
	}

	_, err = api.SetAgentSessionStatus(AgentSessionSetStatusParameters{
		Status: AgentSessionStatusActive,
	})
	if err == nil {
		t.Fatalf("Expected error, got nil")
	}

	if err.Error() != "invalid_arguments" {
		t.Fatalf("Expected error 'invalid_arguments', got %q", err.Error())
	}
}

func agentSessionRenameHandler(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Content-Type", "application/json")

	if r.FormValue("title") == "" || r.FormValue("channel_id") == "" {
		rw.Write([]byte(`{ "ok": false, "error": "invalid_arguments" }`))
		return
	}

	resp, _ := json.Marshal(&AgentSessionRenameResponse{
		SlackResponse: SlackResponse{Ok: true},
		Title:         r.FormValue("title"),
	})
	rw.Write(resp)
}

func TestRenameAgentSession(t *testing.T) {
	http.HandleFunc("/agents.sessions.rename", agentSessionRenameHandler)
	once.Do(startServer)
	api := New("testing-token", OptionAPIURL("http://"+serverAddr+"/"))

	params := AgentSessionRenameParameters{
		Title:     "Bora Bora trip prep",
		ChannelID: "CXXXXXXXX",
		ThreadTS:  "1234567890.123456",
	}

	response, err := api.RenameAgentSession(params)
	if err != nil {
		t.Fatalf("Unexpected error: %s", err)
	}

	if !response.Ok {
		t.Fatalf("Expected Ok to be true")
	}

	if response.Title != "Bora Bora trip prep" {
		t.Fatalf("Expected Title to be 'Bora Bora trip prep', got %q", response.Title)
	}

	_, err = api.RenameAgentSession(AgentSessionRenameParameters{})
	if err == nil {
		t.Fatalf("Expected error, got nil")
	}

	if err.Error() != "invalid_arguments" {
		t.Fatalf("Expected error 'invalid_arguments', got %q", err.Error())
	}
}
