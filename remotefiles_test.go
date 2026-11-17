package slack

import (
	"encoding/json"
	"errors"
	"io"
	"maps"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
	if _, err := api.AddRemoteFile(params); !errors.Is(err, ErrParametersMissing) {
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
		call func(*Client, ListRemoteFilesParameters) ([]RemoteFile, string, error)
	}{
		{"ListRemoteFiles", (*Client).ListRemoteFiles},
		{"ListRemoteFilesContext", func(api *Client, params ListRemoteFilesParameters) ([]RemoteFile, string, error) {
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
					_, nextCursor, err := method.call(api, tt.params)
					require.NoError(t, err)
					assert.Empty(t, nextCursor)
				})
			}
		})
	}
}

func TestListRemoteFilesPagination(t *testing.T) {
	methods := []struct {
		name string
		call func(*Client, ListRemoteFilesParameters) ([]RemoteFile, string, error)
	}{
		{"ListRemoteFiles", (*Client).ListRemoteFiles},
		{"ListRemoteFilesContext", func(api *Client, params ListRemoteFilesParameters) ([]RemoteFile, string, error) {
			return api.ListRemoteFilesContext(t.Context(), params)
		}},
	}
	for _, method := range methods {
		t.Run(method.name, func(t *testing.T) {
			for _, initialCursor := range []string{"", "start+/="} {
				t.Run("cursor="+initialCursor, func(t *testing.T) {
					// Request i sends cursors[i] and returns cursors[i+1].
					cursors := []string{initialCursor, "next+/=", "last+/=", ""}
					responses := []string{
						`{"ok":true,"files":[{"id":"F1"}],"response_metadata":{"next_cursor":"next+/="}}`,
						`{"ok":true,"files":[],"response_metadata":{"next_cursor":"last+/="}}`,
						`{"ok":true,"files":[{"id":"F2"},{"id":"F3"}],"response_metadata":{"next_cursor":""}}`,
					}
					var requests int
					api := New("testing-token", OptionHTTPClient(httpClientFunc(func(r *http.Request) (*http.Response, error) {
						if requests >= len(responses) {
							return nil, errors.New("unexpected extra remote-file list request")
						}
						if err := r.ParseForm(); err != nil {
							return nil, err
						}
						want := url.Values{
							"token": {"testing-token"}, "channel": {"C123"}, "limit": {"2"}, "ts_from": {"12"}, "ts_to": {"34"},
						}
						if cursors[requests] != "" {
							want.Set("cursor", cursors[requests])
						}
						assert.Equal(t, "/api/files.remote.list", r.URL.Path)
						assert.Equal(t, want, r.PostForm)
						body := responses[requests]
						requests++
						return &http.Response{StatusCode: http.StatusOK, Body: io.NopCloser(strings.NewReader(body))}, nil
					})))
					params := ListRemoteFilesParameters{
						Channel: "C123", Cursor: initialCursor, Limit: 2, TimestampFrom: 12, TimestampTo: 34,
					}
					var ids []string
					for page := range responses {
						files, nextCursor, err := method.call(api, params)
						require.NoError(t, err)
						assert.Equal(t, page+1, requests, "each call makes one request")
						assert.Equal(t, cursors[page+1], nextCursor)
						for _, file := range files {
							ids = append(ids, file.ID)
						}
						params.Cursor = nextCursor
					}
					assert.Equal(t, []string{"F1", "F2", "F3"}, ids)
				})
			}
		})
	}
}

func TestListRemoteFilesError(t *testing.T) {
	api := New("testing-token", OptionHTTPClient(httpClientFunc(func(*http.Request) (*http.Response, error) {
		return &http.Response{StatusCode: http.StatusOK, Body: io.NopCloser(strings.NewReader(
			`{"ok":false,"error":"invalid_cursor","response_metadata":{"next_cursor":"next"}}`,
		))}, nil
	})))
	files, nextCursor, err := api.ListRemoteFilesContext(t.Context(), ListRemoteFilesParameters{Cursor: "stale"})
	require.ErrorContains(t, err, "invalid_cursor")
	assert.Nil(t, files)
	assert.Empty(t, nextCursor)
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
	if _, err := api.ShareRemoteFile([]string{}, "ExternalID", ""); !errors.Is(err, ErrParametersMissing) {
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
