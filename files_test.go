package slack

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"reflect"
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

func TestSlack_DeleteFileComment(t *testing.T) {
	once.Do(startServer)
	SLACK_API = "http://" + serverAddr + "/"
	api := New("testing-token")
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
	SLACK_API = "http://" + serverAddr + "/"
	api := New("testing-token")
	params := FileUploadParameters{
		Filename: "test.txt", Content: "test content",
		Channels: []string{"CXXXXXXXX"}}
	if _, err := api.UploadFile(params); err != nil {
		t.Errorf("Unexpected error: %s", err)
	}

	reader := bytes.NewBufferString("test reader")
	params = FileUploadParameters{
		Filename: "test.txt", Reader: reader,
		Channels: []string{"CXXXXXXXX"}}
	if _, err := api.UploadFile(params); err != nil {
		t.Errorf("Unexpected error: %s", err)
	}
}
