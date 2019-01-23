package slack

import (
	"net/http"
	"testing"
)

func handleAppsUninstall(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Content-Type", "application/json")

	clientID := r.FormValue("client_id")
	clientSecret := r.FormValue("client_secret")
	token := r.FormValue("token")

	// First, nothing should be empty
	if clientID == "" {
		rw.Write([]byte(`{"ok":false,"error":"client_id is empty"}`))
		return
	}
	if clientSecret == "" {
		rw.Write([]byte(`{"ok":false,"error":"client_secret is empty"}`))
		return
	}
	if token == "" {
		rw.Write([]byte(`{"ok":false,"error":"token is empty"}`))
		return
	}
	if token != validToken {
		rw.Write([]byte(`{"ok":false,"error":"invalid_token"}`))
		return
	}
	response := []byte(`{"ok": true}`)
	rw.Write(response)
}

func TestAppsUninstall(t *testing.T) {
	http.HandleFunc("/apps.uninstall", handleAppsUninstall)

	once.Do(startServer)
	APIURL = "http://" + serverAddr + "/"
	api := New("testing-token")

	resp, err := api.SendAppsUninstall("client ID!", "Client Secret!", "testing-token")
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}

	if resp.Ok != true {
		t.Errorf("Got not ok: %v (%v)", resp, resp.Ok)
	}

	// Test token substitution from api
	resp, err = api.SendAppsUninstall("client ID!", "Client Secret!", "")
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}

	if resp.Ok != true {
		t.Errorf("Got not ok: %v (%v)", resp, resp.Ok)
	}

	// Now some negative tests

	// Test bad token (Not a very good test, but, it does test the negative side of handleAppsUninstall)
	// But, that's ok since our code does not throw errors internally, and just hands them back from
	// slack
	resp, err = api.SendAppsUninstall("client ID!", "Client Secret!", "ng-token")
	if err == nil {
		t.Errorf("Should have gotten an error!")
	}

	if resp.Ok != false {
		t.Errorf("Got ok!: %v (%v)", resp, resp.Ok)
	}

}
