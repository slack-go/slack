package slack

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"reflect"
	"strings"
	"testing"
)

type fileCommentHandler struct {
	gotParams map[string]string
}

func newFileCommentHandler() *fileCommentHandler {
	return &fileCommentHandler{
		gotParams: make(map[string]string),
	}
}

func (h *fileCommentHandler) accumulateFormValue(k string, r *http.Request) {
	if v := r.FormValue(k); v != "" {
		h.gotParams[k] = v
	}
}

func (h *fileCommentHandler) handler(w http.ResponseWriter, r *http.Request) {
	h.accumulateFormValue("token", r)
	h.accumulateFormValue("file", r)
	h.accumulateFormValue("id", r)

	w.Header().Set("Content-Type", "application/json")
	if h.gotParams["id"] == "trigger-error" {
		w.Write([]byte(`{ "ok": false, "error": "errored" }`))
	} else {
		w.Write([]byte(`{ "ok": true }`))
	}
}

type mockHTTPClient struct{}

func (m *mockHTTPClient) Do(*http.Request) (*http.Response, error) {
	return &http.Response{StatusCode: 200, Body: ioutil.NopCloser(bytes.NewBufferString(`OK`))}, nil
}

func TestSlack_GetFile(t *testing.T) {
	api := &Client{
		endpoint:   "http://" + serverAddr + "/",
		token:      "testing-token",
		httpclient: &mockHTTPClient{},
	}

	tests := []struct {
		title       string
		downloadURL string
		expectError bool
	}{
		{
			title:       "Testing with valid file",
			downloadURL: "https://files.slack.com/files-pri/T99999999-FGGGGGGGG/download/test.csv",
			expectError: false,
		},
		{
			title:       "Testing with invalid file (empty URL)",
			downloadURL: "",
			expectError: true,
		},
	}

	for _, test := range tests {
		err := api.GetFile(test.downloadURL, &bytes.Buffer{})

		if !test.expectError && err != nil {
			log.Fatalf("%s: Unexpected error: %s in test", test.title, err)
		} else if test.expectError == true && err == nil {
			log.Fatalf("Expected error but got none")
		}
	}
}

func TestSlack_DeleteFileComment(t *testing.T) {
	once.Do(startServer)
	api := New("testing-token", OptionAPIURL("http://"+serverAddr+"/"))
	tests := []struct {
		title       string
		body        url.Values
		wantParams  map[string]string
		expectError bool
	}{
		{
			title: "Testing with proper body",
			body: url.Values{
				"file": {"file12345"},
				"id":   {"id12345"},
			},
			wantParams: map[string]string{
				"token": "testing-token",
				"file":  "file12345",
				"id":    "id12345",
			},
			expectError: false,
		},
		{
			title: "Testing with false body",
			body: url.Values{
				"file": {""},
				"id":   {""},
			},
			wantParams:  map[string]string{},
			expectError: true,
		},
		{
			title: "Testing with error",
			body: url.Values{
				"file": {"file12345"},
				"id":   {"trigger-error"},
			},
			wantParams: map[string]string{
				"token": "testing-token",
				"file":  "file12345",
				"id":    "trigger-error",
			},
			expectError: true,
		},
	}

	var fch *fileCommentHandler
	http.HandleFunc("/files.comments.delete", func(w http.ResponseWriter, r *http.Request) {
		fch.handler(w, r)
	})

	for _, test := range tests {
		fch = newFileCommentHandler()
		err := api.DeleteFileComment(test.body["id"][0], test.body["file"][0])

		if !test.expectError && err != nil {
			log.Fatalf("%s: Unexpected error: %s in test", test.title, err)
		} else if test.expectError == true && err == nil {
			log.Fatalf("Expected error but got none")
		}

		if !reflect.DeepEqual(fch.gotParams, test.wantParams) {
			log.Fatalf("%s: Got params [%#v]\nBut received [%#v]\n", test.title, fch.gotParams, test.wantParams)
		}
	}
}

func authTestHandler(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	response, _ := json.Marshal(authTestResponseFull{
		SlackResponse: SlackResponse{Ok: true}})
	rw.Write(response)
}

func uploadFileHandler(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	response, _ := json.Marshal(fileResponseFull{
		SlackResponse: SlackResponse{Ok: true}})
	rw.Write(response)
}

func TestUploadFile(t *testing.T) {
	http.HandleFunc("/auth.test", authTestHandler)
	http.HandleFunc("/files.upload", uploadFileHandler)
	once.Do(startServer)
	api := New("testing-token", OptionAPIURL("http://"+serverAddr+"/"))
	params := FileUploadParameters{
		Filename: "test.txt", Content: "test content",
		Channels: []string{"CXXXXXXXX"}}
	if _, err := api.UploadFile(params); err != nil {
		t.Errorf("Unexpected error: %s", err)
	}

	reader := bytes.NewBufferString("test reader")
	params = FileUploadParameters{
		Filename: "test.txt",
		Reader:   reader,
		Channels: []string{"CXXXXXXXX"}}
	if _, err := api.UploadFile(params); err != nil {
		t.Errorf("Unexpected error: %s", err)
	}

	largeByt := make([]byte, 107374200)
	reader = bytes.NewBuffer(largeByt)
	params = FileUploadParameters{
		Filename: "test.txt", Reader: reader,
		Channels: []string{"CXXXXXXXX"}}
	if _, err := api.UploadFile(params); err != nil {
		t.Errorf("Unexpected error: %s", err)
	}
}

func TestUploadFileWithoutFilename(t *testing.T) {
	once.Do(startServer)
	api := New("testing-token", OptionAPIURL("http://"+serverAddr+"/"))

	reader := bytes.NewBufferString("test reader")
	params := FileUploadParameters{
		Reader:   reader,
		Channels: []string{"CXXXXXXXX"}}
	_, err := api.UploadFile(params)
	if err == nil {
		t.Fatal("Expected error when omitting filename, instead got nil")
	}

	if !strings.Contains(err.Error(), ".Filename is mandatory") {
		t.Errorf("Error message should mention empty FileUploadParameters.Filename")
	}
}

func uploadURLHandler(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	response, _ := json.Marshal(getUploadURLExternalResponse{
		FileID:        "RandomID",
		UploadURL:     "http://" + serverAddr + "/abc",
		SlackResponse: SlackResponse{Ok: true}})
	rw.Write(response)
}

func urlFileUploadHandler(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Content-Type", "text")
	rw.Write([]byte("Ok: 200, file uploaded"))
}

func completeURLUpload(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	response, _ := json.Marshal(completeUploadExternalResponse{
		Files: []FileSummary{
			{
				ID:    "RandomID",
				Title: "",
			},
		},
		SlackResponse: SlackResponse{Ok: true}})
	rw.Write(response)
}

func TestUploadFileV2(t *testing.T) {
	http.HandleFunc("/files.getUploadURLExternal", uploadURLHandler)
	http.HandleFunc("/abc", urlFileUploadHandler)
	http.HandleFunc("/files.completeUploadExternal", completeURLUpload)
	once.Do(startServer)
	api := New("testing-token", OptionAPIURL("http://"+serverAddr+"/"))

	params := UploadFileV2Parameters{
		Filename: "test.txt", Content: "test content", FileSize: 10,
		Channel: "CXXXXXXXX",
	}
	if _, err := api.UploadFileV2(params); err != nil {
		t.Errorf("Unexpected error: %s", err)
	}

	reader := bytes.NewBufferString("test reader")
	params = UploadFileV2Parameters{
		Filename: "test.txt",
		Reader:   reader,
		FileSize: 10,
		Channel:  "CXXXXXXXX"}
	if _, err := api.UploadFileV2(params); err != nil {
		t.Errorf("Unexpected error: %s", err)
	}

	largeByt := make([]byte, 107374200)
	reader = bytes.NewBuffer(largeByt)
	params = UploadFileV2Parameters{
		Filename: "test.txt", Reader: reader, FileSize: len(largeByt),
		Channel: "CXXXXXXXX"}
	if _, err := api.UploadFileV2(params); err != nil {
		t.Errorf("Unexpected error: %s", err)
	}
}
