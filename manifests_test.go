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

func getTestManifestResponse() *ManifestResponse {
	return &ManifestResponse{
		SlackResponse: SlackResponse{
			Ok: true,
		},
	}
}
