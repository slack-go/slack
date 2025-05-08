package slack

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func getTestCall(callID string) Call {
	return Call{
		ID:               callID,
		Title:            "test call",
		JoinURL:          "https://example.com/example",
		ExternalUniqueID: "123",
	}
}

func testClient(api string, f http.HandlerFunc) *Client {
	http.HandleFunc(api, f)
	once.Do(startServer)
	return New("testing-token", OptionAPIURL("http://"+serverAddr+"/"))
}

var callTestId = 999

func addCallHandler(t *testing.T) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		rw.Header().Set("Content-Type", "application/json")
		if err := r.ParseForm(); err != nil {
			httpTestErrReply(rw, true, fmt.Sprintf("err parsing form: %s", err.Error()))
			return
		}
		call := Call{
			ID:                fmt.Sprintf("R%d", callTestId),
			Title:             r.FormValue("title"),
			JoinURL:           r.FormValue("join_url"),
			ExternalUniqueID:  r.FormValue("external_unique_id"),
			ExternalDisplayID: r.FormValue("external_display_id"),
			DesktopAppJoinURL: r.FormValue("desktop_app_join_url"),
		}
		callTestId += 1
		json.Unmarshal([]byte(r.FormValue("users")), &call.Participants)
		if start := r.FormValue("date_start"); start != "" {
			dateStart, err := strconv.ParseInt(start, 10, 64)
			require.NoError(t, err)
			call.DateStart = JSONTime(dateStart)
		}
		resp, _ := json.Marshal(callResponse{Call: call, SlackResponse: SlackResponse{Ok: true}})
		rw.Write(resp)
	}
}

func TestAddCall(t *testing.T) {
	api := testClient("/calls.add", addCallHandler(t))
	params := AddCallParameters{
		Title:            "test call",
		JoinURL:          "https://example.com/example",
		ExternalUniqueID: "123",
	}
	call, err := api.AddCall(params)
	require.NoError(t, err)
	assert.Equal(t, params.Title, call.Title)
	assert.Equal(t, params.JoinURL, call.JoinURL)
	assert.Equal(t, params.ExternalUniqueID, call.ExternalUniqueID)
}

func getCallHandler(calls []Call) func(rw http.ResponseWriter, r *http.Request) {
	return func(rw http.ResponseWriter, r *http.Request) {
		callID := r.FormValue("id")

		rw.Header().Set("Content-Type", "application/json")
		for _, call := range calls {
			if call.ID == callID {
				resp, _ := json.Marshal(callResponse{Call: call, SlackResponse: SlackResponse{Ok: true}})
				rw.Write(resp)
				return
			}
		}
		// Fail if the call doesn't exist
		rw.Write([]byte(`{ "ok": false, "error": "not_found" }`))
	}
}

func TestGetCall(t *testing.T) {
	calls := []Call{
		getTestCall("R1234567890"),
		getTestCall("R1234567891"),
	}
	http.HandleFunc("/calls.info", getCallHandler(calls))
	once.Do(startServer)
	api := New("testing-token", OptionAPIURL("http://"+serverAddr+"/"))

	for _, call := range calls {
		resp, err := api.GetCall(call.ID)
		require.NoError(t, err)
		assert.Equal(t, call, resp)
	}
	// Test a call that doesn't exist
	_, err := api.GetCall("R1234567892")
	require.Error(t, err)
}

func updateCallHandler(calls []Call) func(rw http.ResponseWriter, r *http.Request) {
	return func(rw http.ResponseWriter, r *http.Request) {
		callID := r.FormValue("id")

		rw.Header().Set("Content-Type", "application/json")
		if err := r.ParseForm(); err != nil {
			httpTestErrReply(rw, true, fmt.Sprintf("err parsing form: %s", err.Error()))
			return
		}

		for _, call := range calls {
			if call.ID == callID {
				if title := r.FormValue("title"); title != "" {
					call.Title = title
				}
				if joinURL := r.FormValue("join_url"); joinURL != "" {
					call.JoinURL = joinURL
				}
				if desktopAppJoinURL := r.FormValue("desktop_app_join_url"); desktopAppJoinURL != "" {
					call.DesktopAppJoinURL = desktopAppJoinURL
				}
				resp, _ := json.Marshal(callResponse{Call: call, SlackResponse: SlackResponse{Ok: true}})
				rw.Write(resp)
				return
			}
		}
		// Fail if the call doesn't exist
		rw.Write([]byte(`{ "ok": false, "error": "not_found" }`))
	}
}

func TestUpdateCall(t *testing.T) {
	calls := []Call{
		getTestCall("R1234567890"),
		getTestCall("R1234567891"),
		getTestCall("R1234567892"),
		getTestCall("R1234567893"),
		getTestCall("R1234567894"),
	}
	http.HandleFunc("/calls.update", updateCallHandler(calls))
	once.Do(startServer)
	api := New("testing-token", OptionAPIURL("http://"+serverAddr+"/"))

	changes := []struct {
		callID string
		params UpdateCallParameters
	}{
		{
			callID: "R1234567890",
			params: UpdateCallParameters{Title: "test"},
		},
		{
			callID: "R1234567891",
			params: UpdateCallParameters{JoinURL: "https://example.com/join"},
		},
		{
			callID: "R1234567892",
			params: UpdateCallParameters{DesktopAppJoinURL: "https://example.com/join"},
		},
		{ // Change multiple fields at once
			callID: "R1234567893",
			params: UpdateCallParameters{
				Title:   "test",
				JoinURL: "https://example.com/join",
			},
		},
	}

	for _, change := range changes {
		call, err := api.UpdateCall(change.callID, change.params)
		require.NoError(t, err)
		if change.params.Title != "" && call.Title != change.params.Title {
			t.Fatalf("Expected title to be %s, got %s", change.params.Title, call.Title)
		}
		if change.params.JoinURL != "" && call.JoinURL != change.params.JoinURL {
			t.Fatalf("Expected join_url to be %s, got %s", change.params.JoinURL, call.JoinURL)
		}
		if change.params.DesktopAppJoinURL != "" && call.DesktopAppJoinURL != change.params.DesktopAppJoinURL {
			t.Fatalf("Expected desktop_app_join_url to be %s, got %s", change.params.DesktopAppJoinURL, call.DesktopAppJoinURL)
		}
	}
}
