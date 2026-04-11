package slack

import (
	"encoding/json"
	"net/http"
	"reflect"
	"strings"
	"testing"
)

func TestCreateManifest(t *testing.T) {
	http.HandleFunc("/apps.manifest.create", handleCreateManifest)
	once.Do(startServer)

	api := New("testing-token", OptionAPIURL("http://"+serverAddr+"/"))

	manif := getTestManifest()
	resp, err := api.CreateManifest(&manif, "token")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
		return
	}

	if !reflect.DeepEqual(resp, getTestManifestResponse()) {
		t.Fatal(ErrIncorrectResponse)
	}
}

func handleCreateManifest(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Content-Type", "application/json")

	response, _ := json.Marshal(getTestManifestResponse())
	rw.Write(response)
}

func TestDeleteManifest(t *testing.T) {
	http.HandleFunc("/apps.manifest.delete", handleDeleteManifest)
	expectedResponse := SlackResponse{Ok: true}

	once.Do(startServer)
	api := New("testing-token", OptionAPIURL("http://"+serverAddr+"/"))

	resp, err := api.DeleteManifest("token", "app id")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
		return
	}

	if !reflect.DeepEqual(expectedResponse, *resp) {
		t.Fatal(ErrIncorrectResponse)
	}
}

func handleDeleteManifest(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Content-Type", "application/json")

	response, _ := json.Marshal(SlackResponse{Ok: true})
	rw.Write(response)
}

func TestExportManifest(t *testing.T) {
	http.HandleFunc("/apps.manifest.export", handleExportManifest)
	expectedResponse := getTestManifest()

	api := New("testing-token", OptionAPIURL("http://"+serverAddr+"/"))

	resp, err := api.ExportManifest("token", "app id")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
		return
	}

	if !reflect.DeepEqual(expectedResponse, *resp) {
		t.Fatal(ErrIncorrectResponse)
	}
}

func handleExportManifest(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Content-Type", "application/json")

	response, _ := json.Marshal(ExportManifestResponse{Manifest: getTestManifest()})
	rw.Write(response)
}

func TestUpdateManifest(t *testing.T) {
	http.HandleFunc("/apps.manifest.update", handleUpdateManifest)
	expectedResponse := UpdateManifestResponse{AppId: "app id"}

	api := New("testing-token", OptionAPIURL("http://"+serverAddr+"/"))

	manif := getTestManifest()
	resp, err := api.UpdateManifest(&manif, "token", "app id")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
		return
	}

	if !reflect.DeepEqual(expectedResponse, *resp) {
		t.Fatal(ErrIncorrectResponse)
	}
}

func handleUpdateManifest(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Content-Type", "application/json")

	response, _ := json.Marshal(UpdateManifestResponse{AppId: "app id"})
	rw.Write(response)
}

func TestValidateManifest(t *testing.T) {
	http.HandleFunc("/apps.manifest.validate", handleValidateManifest)
	expectedResponse := ManifestResponse{SlackResponse: SlackResponse{Ok: true}}

	api := New("testing-token", OptionAPIURL("http://"+serverAddr+"/"))

	manif := getTestManifest()
	resp, err := api.ValidateManifest(&manif, "token", "app id")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
		return
	}

	if !reflect.DeepEqual(expectedResponse, *resp) {
		t.Fatal(ErrIncorrectResponse)
	}
}

func handleValidateManifest(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Content-Type", "application/json")

	response, _ := json.Marshal(ManifestResponse{SlackResponse: SlackResponse{Ok: true}})
	rw.Write(response)
}

func getTestManifest() Manifest {
	return Manifest{
		Display: Display{
			Name:        "test",
			Description: "this is a test",
		},
	}
}

func TestOAuthScopesOptionalFields(t *testing.T) {
	scopes := OAuthScopes{
		Bot:          []string{"chat:write", "commands"},
		User:         []string{"users:read"},
		BotOptional:  []string{"files:read", "reactions:read"},
		UserOptional: []string{"channels:read"},
	}

	data, err := json.Marshal(scopes)
	if err != nil {
		t.Fatalf("Marshal error: %s", err)
	}

	var roundtrip OAuthScopes
	if err := json.Unmarshal(data, &roundtrip); err != nil {
		t.Fatalf("Unmarshal error: %s", err)
	}

	if !reflect.DeepEqual(scopes, roundtrip) {
		t.Errorf("Round-trip mismatch: got %+v, want %+v", roundtrip, scopes)
	}

	// Verify omitempty: empty optional fields should not appear
	minimal := OAuthScopes{Bot: []string{"chat:write"}}
	data, err = json.Marshal(minimal)
	if err != nil {
		t.Fatalf("Marshal error: %s", err)
	}
	s := string(data)
	if strings.Contains(s, "bot_optional") {
		t.Errorf("Expected bot_optional to be omitted from JSON: %s", s)
	}
	if strings.Contains(s, "user_optional") {
		t.Errorf("Expected user_optional to be omitted from JSON: %s", s)
	}
}

func getTestManifestResponse() *ManifestResponse {
	return &ManifestResponse{
		SlackResponse: SlackResponse{
			Ok: true,
		},
	}
}
