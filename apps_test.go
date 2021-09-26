package slack

import (
	"encoding/json"
	"net/http"
	"testing"
)

func TestListEventAuthorizations(t *testing.T) {
	http.HandleFunc("/apps.event.authorizations.list", testListEventAuthorizationsHandler)
	once.Do(startServer)

	api := New("", OptionAppLevelToken("test-token"), OptionAPIURL("http://"+serverAddr+"/"))

	authorizations, err := api.ListEventAuthorizations("1-message-T012345678-DR12345678")

	if err != nil {
		t.Errorf("Failed, but should have succeeded")
	} else if len(authorizations) != 1 {
		t.Errorf("Didn't get 1 authorization")
	} else if authorizations[0].UserID != "U123456789" {
		t.Errorf("User ID is wrong")
	}
}

func testListEventAuthorizationsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	response, _ := json.Marshal(listEventAuthorizationsResponse{
		SlackResponse: SlackResponse{Ok: true},
		Authorizations: []EventAuthorization{
			{
				UserID: "U123456789",
				TeamID: "T012345678",
			},
		},
	})
	w.Write(response)
}

func TestUninstallApp(t *testing.T) {
	http.HandleFunc("/apps.uninstall", testUninstallAppHandler)
	once.Do(startServer)

	api := New("test-token", OptionAPIURL("http://"+serverAddr+"/"))

	err := api.UninstallApp("", "")

	if err != nil {
		t.Errorf("Failed, but should have succeeded")
	}
}

func testUninstallAppHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	response, _ := json.Marshal(SlackResponse{Ok: true})
	w.Write(response)
}
