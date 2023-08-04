package slack

import (
	"encoding/json"
	"net/http"
	"reflect"
	"testing"
)

func TestCreateManifest(t *testing.T) {
	http.HandleFunc("/apps.manifest.create", handleCreateManifest)
	once.Do(startServer)

	api := New("testing-token", OptionAPIURL("http://"+serverAddr+"/"))

	resp, err := api.CreateManifest(getTestManifest(), "token")
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

func getTestManifest() *Manifest {
	return &Manifest{
		Display: Display{
			Name:        "test",
			Description: "this is a test",
		},
	}
}

func getTestManifestResponse() *ManifestResponse {
	return &ManifestResponse{
		SlackResponse: SlackResponse{
			Ok: true,
		},
	}
}
