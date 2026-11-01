package slack

import (
	"encoding/json"
	"maps"
	"net/http"
	"net/http/httptest"
	"net/url"
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
	http.HandleFunc("/files.remote.add", addRemoteFileHandler)
	once.Do(startServer)
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
	once.Do(startServer)
	api := New("testing-token", OptionAPIURL("http://"+serverAddr+"/"))
	params := RemoteFileParameters{
		ExternalID:  "externalID",
		ExternalURL: "http://example.com/",
	}
	if _, err := api.AddRemoteFile(params); err != ErrParametersMissing {
		t.Errorf("Expected ErrParametersMissing. got %s", err)
	}
}

func TestListRemoteFile(t *testing.T) {
	tests := []struct {
		name   string
		params ListRemoteFilesParameters
		want   url.Values
	}{
		{name: "empty"},
		{name: "channel only", params: ListRemoteFilesParameters{Channel: "C123"}, want: url.Values{"channel": {"C123"}}},
		{
			name: "sentinel defaults",
			params: ListRemoteFilesParameters{
				TimestampTo: DEFAULT_REMOTE_FILES_TS_TO,
				Limit:       DEFAULT_REMOTE_FILES_COUNT,
			},
		},
		{
			name: "explicit filters",
			params: ListRemoteFilesParameters{
				Channel: "C123", Cursor: "next+/=", Limit: 7, TimestampFrom: 12, TimestampTo: 34,
			},
			want: url.Values{
				"channel": {"C123"}, "cursor": {"next+/="}, "limit": {"7"}, "ts_from": {"12"}, "ts_to": {"34"},
			},
		},
	}
	methods := []struct {
		name string
		call func(*Client, ListRemoteFilesParameters) ([]RemoteFile, error)
	}{
		{"ListRemoteFiles", (*Client).ListRemoteFiles},
		{"ListRemoteFilesContext", func(api *Client, params ListRemoteFilesParameters) ([]RemoteFile, error) {
			return api.ListRemoteFilesContext(t.Context(), params)
		}},
	}
	for _, method := range methods {
		t.Run(method.name, func(t *testing.T) {
			for _, tt := range tests {
				t.Run(tt.name, func(t *testing.T) {
					want := url.Values{"token": {"testing-token"}}
					maps.Copy(want, tt.want)
					server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
						if err := r.ParseForm(); err != nil {
							t.Error(err)
							http.Error(w, err.Error(), http.StatusBadRequest)
							return
						}
						if got := r.PostForm.Encode(); got != want.Encode() {
							t.Errorf("form = %q, want %q", got, want.Encode())
						}
						w.Header().Set("Content-Type", "application/json")
						_, _ = w.Write([]byte(`{"ok":true,"files":[]}`))
					}))
					defer server.Close()

					api := New("testing-token", OptionAPIURL(server.URL+"/"))
					if _, err := method.call(api, tt.params); err != nil {
						t.Fatal(err)
					}
				})
			}
		})
	}
}

func getRemoteFileInfoHandler(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	response, _ := json.Marshal(remoteFileResponseFull{
		SlackResponse: SlackResponse{Ok: true}})
	rw.Write(response)
}

func TestGetRemoteFileInfo(t *testing.T) {
	http.HandleFunc("/files.remote.info", getRemoteFileInfoHandler)
	once.Do(startServer)
	api := New("testing-token", OptionAPIURL("http://"+serverAddr+"/"))
	if _, err := api.GetRemoteFileInfo("ExternalID", ""); err != nil {
		t.Errorf("Unexpected error: %s", err)
	}
}

func TestGetRemoteFileInfoWithoutID(t *testing.T) {
	once.Do(startServer)
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
	once.Do(startServer)
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
	http.HandleFunc("/files.remote.share", shareRemoteFileHandler)
	once.Do(startServer)
	api := New("testing-token", OptionAPIURL("http://"+serverAddr+"/"))
	if _, err := api.ShareRemoteFile([]string{"channel"}, "ExternalID", ""); err != nil {
		t.Errorf("Unexpected error: %s", err)
	}
}

func TestShareRemoteFileWithoutChannels(t *testing.T) {
	once.Do(startServer)
	api := New("testing-token", OptionAPIURL("http://"+serverAddr+"/"))
	if _, err := api.ShareRemoteFile([]string{}, "ExternalID", ""); err != ErrParametersMissing {
		t.Errorf("Expected ErrParametersMissing. got %s", err)
	}
}

func TestShareRemoteFileWithoutID(t *testing.T) {
	once.Do(startServer)
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
	http.HandleFunc("/files.remote.update", updateRemoteFileHandler)
	once.Do(startServer)
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
	http.HandleFunc("/files.remote.remove", removeRemoteFileHandler)
	once.Do(startServer)
	api := New("testing-token", OptionAPIURL("http://"+serverAddr+"/"))
	if err := api.RemoveRemoteFile("ExternalID", ""); err != nil {
		t.Errorf("Unexpected error: %s", err)
	}
}

func TestRemoveRemoteFileWithoutID(t *testing.T) {
	once.Do(startServer)
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
	once.Do(startServer)
	api := New("testing-token", OptionAPIURL("http://"+serverAddr+"/"))
	err := api.RemoveRemoteFile("ExternalID", "FileID")
	if err == nil {
		t.Fatal("Expected error when both externalID and fileID are both provided, instead got nil")
	}
	if !strings.Contains(err.Error(), "don't provide both externalID and fileID") {
		t.Errorf("Error message should mention don't providing both externalID and fileID")
	}
}

func TestAddRemoteFileWithPreviewImageReader(t *testing.T) {
	once.Do(startServer)
	api := New("testing-token", OptionAPIURL("http://"+serverAddr+"/"))
	params := RemoteFileParameters{
		ExternalID:         "externalID",
		ExternalURL:        "http://example.com/",
		Title:              "example",
		PreviewImageReader: strings.NewReader("fake image data"),
		PreviewImageName:   "preview.jpg",
	}
	if _, err := api.AddRemoteFile(params); err != nil {
		t.Errorf("Unexpected error: %s", err)
	}
}

func TestUpdateRemoteFileWithPreviewImageReader(t *testing.T) {
	once.Do(startServer)
	api := New("testing-token", OptionAPIURL("http://"+serverAddr+"/"))
	params := RemoteFileParameters{
		Title:              "updated example",
		PreviewImageReader: strings.NewReader("fake image data"),
		PreviewImageName:   "preview.jpg",
	}
	if _, err := api.UpdateRemoteFile("fileID", params); err != nil {
		t.Errorf("Unexpected error: %s", err)
	}
}
