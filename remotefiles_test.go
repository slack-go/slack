package slack

import (
	"encoding/json"
	"net/http"
	"strings"
	"testing"
)

func addRemoteFileHandler(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	response, _ := json.Marshal(remoteFileResponseFull{
		SlackResponse: SlackResponse{Ok: true}})
	rw.Write(response)
}

func TestAddRemoteFile(t *testing.T) {
	s := startServer()
	defer s.Close()

	s.RegisterHandler("/files.remote.add", addRemoteFileHandler)

	api := New("testing-token", OptionAPIURL("http://"+serverAddr+"/"))
	params := RemoteFileParameters{
		ExternalID:  "externalID",
		ExternalURL: "http://example.com/",
		Title:       "example",
	}
	if _, err := api.AddRemoteFile(params); err != nil {
		t.Errorf("Unexpected error: %s", err)
	}
}

func TestAddRemoteFileWithoutTitle(t *testing.T) {
	api := New("testing-token", OptionAPIURL("http://"+serverAddr+"/"))
	params := RemoteFileParameters{
		ExternalID:  "externalID",
		ExternalURL: "http://example.com/",
	}
	if _, err := api.AddRemoteFile(params); err != ErrParametersMissing {
		t.Errorf("Expected ErrParametersMissing. got %s", err)
	}
}

func listRemoteFileHandler(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	response, _ := json.Marshal(remoteFileResponseFull{
		SlackResponse: SlackResponse{Ok: true}})
	rw.Write(response)
}

func TestListRemoteFile(t *testing.T) {
	s := startServer()
	defer s.Close()

	s.RegisterHandler("/files.remote.list", listRemoteFileHandler)

	api := New("testing-token", OptionAPIURL("http://"+serverAddr+"/"))
	params := ListRemoteFilesParameters{}
	if _, err := api.ListRemoteFiles(params); err != nil {
		t.Errorf("Unexpected error: %s", err)
	}
}

func getRemoteFileInfoHandler(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	response, _ := json.Marshal(remoteFileResponseFull{
		SlackResponse: SlackResponse{Ok: true}})
	rw.Write(response)
}

func TestGetRemoteFileInfo(t *testing.T) {
	s := startServer()
	defer s.Close()

	s.RegisterHandler("/files.remote.info", getRemoteFileInfoHandler)

	api := New("testing-token", OptionAPIURL("http://"+serverAddr+"/"))
	if _, err := api.GetRemoteFileInfo("ExternalID", ""); err != nil {
		t.Errorf("Unexpected error: %s", err)
	}
}

func TestGetRemoteFileInfoWithoutID(t *testing.T) {
	api := New("testing-token", OptionAPIURL("http://"+serverAddr+"/"))
	_, err := api.GetRemoteFileInfo("", "")
	if err == nil {
		t.Fatal("Expected error when both externalID and fileID is not provided, instead got nil")
	}
	if !strings.Contains(err.Error(), "either externalID or fileID is required") {
		t.Errorf("Error message should mention a required field")
	}
}

func TestGetRemoteFileInfoWithFileIDAndExternalID(t *testing.T) {
	api := New("testing-token", OptionAPIURL("http://"+serverAddr+"/"))
	_, err := api.GetRemoteFileInfo("ExternalID", "FileID")
	if err == nil {
		t.Fatal("Expected error when both externalID and fileID are both provided, instead got nil")
	}
	if !strings.Contains(err.Error(), "don't provide both externalID and fileID") {
		t.Errorf("Error message should mention don't providing both externalID and fileID")
	}
}

func shareRemoteFileHandler(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	response, _ := json.Marshal(remoteFileResponseFull{
		SlackResponse: SlackResponse{Ok: true}})
	rw.Write(response)
}

func TestShareRemoteFile(t *testing.T) {
	s := startServer()
	defer s.Close()

	s.RegisterHandler("/files.remote.share", shareRemoteFileHandler)

	api := New("testing-token", OptionAPIURL("http://"+serverAddr+"/"))
	if _, err := api.ShareRemoteFile([]string{"channel"}, "ExternalID", ""); err != nil {
		t.Errorf("Unexpected error: %s", err)
	}
}

func TestShareRemoteFileWithoutChannels(t *testing.T) {
	api := New("testing-token", OptionAPIURL("http://"+serverAddr+"/"))
	if _, err := api.ShareRemoteFile([]string{}, "ExternalID", ""); err != ErrParametersMissing {
		t.Errorf("Expected ErrParametersMissing. got %s", err)
	}
}

func TestShareRemoteFileWithoutID(t *testing.T) {
	api := New("testing-token", OptionAPIURL("http://"+serverAddr+"/"))
	_, err := api.ShareRemoteFile([]string{"channel"}, "", "")
	if err == nil {
		t.Fatal("Expected error when both externalID and fileID is not provided, instead got nil")
	}
	if !strings.Contains(err.Error(), "either externalID or fileID is required") {
		t.Errorf("Error message should mention a required field")
	}
}

func updateRemoteFileHandler(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	response, _ := json.Marshal(remoteFileResponseFull{
		SlackResponse: SlackResponse{Ok: true}})
	rw.Write(response)
}

func TestUpdateRemoteFile(t *testing.T) {
	s := startServer()
	defer s.Close()

	s.RegisterHandler("/files.remote.update", updateRemoteFileHandler)

	api := New("testing-token", OptionAPIURL("http://"+serverAddr+"/"))
	params := RemoteFileParameters{
		ExternalURL: "http://example.com/",
		Title:       "example",
	}
	if _, err := api.UpdateRemoteFile("fileID", params); err != nil {
		t.Errorf("Unexpected error: %s", err)
	}
}

func removeRemoteFileHandler(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	response, _ := json.Marshal(remoteFileResponseFull{
		SlackResponse: SlackResponse{Ok: true}})
	rw.Write(response)
}

func TestRemoveRemoteFile(t *testing.T) {
	s := startServer()
	defer s.Close()

	s.RegisterHandler("/files.remote.remove", removeRemoteFileHandler)

	api := New("testing-token", OptionAPIURL("http://"+serverAddr+"/"))
	if err := api.RemoveRemoteFile("ExternalID", ""); err != nil {
		t.Errorf("Unexpected error: %s", err)
	}
}

func TestRemoveRemoteFileWithoutID(t *testing.T) {
	api := New("testing-token", OptionAPIURL("http://"+serverAddr+"/"))
	err := api.RemoveRemoteFile("", "")
	if err == nil {
		t.Fatal("Expected error when both externalID and fileID is not provided, instead got nil")
	}
	if !strings.Contains(err.Error(), "either externalID or fileID is required") {
		t.Errorf("Error message should mention a required field")
	}
}

func TestRemoveRemoteFileWithFileIDAndExternalID(t *testing.T) {
	api := New("testing-token", OptionAPIURL("http://"+serverAddr+"/"))
	err := api.RemoveRemoteFile("ExternalID", "FileID")
	if err == nil {
		t.Fatal("Expected error when both externalID and fileID are both provided, instead got nil")
	}
	if !strings.Contains(err.Error(), "don't provide both externalID and fileID") {
		t.Errorf("Error message should mention don't providing both externalID and fileID")
	}
}
