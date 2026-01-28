package slack

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
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

type mockHTTPClient struct{}

func (m *mockHTTPClient) Do(*http.Request) (*http.Response, error) {
	return &http.Response{StatusCode: 200, Body: io.NopCloser(bytes.NewBufferString(`OK`))}, nil
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

func uploadURLHandler(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	response, _ := json.Marshal(GetUploadURLExternalResponse{
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
	response, _ := json.Marshal(CompleteUploadExternalResponse{
		Files: []FileSummary{
			{
				ID:    "RandomID",
				Title: "",
			},
		},
		SlackResponse: SlackResponse{Ok: true}})
	rw.Write(response)
}

func TestUploadFile(t *testing.T) {
	http.HandleFunc("/files.getUploadURLExternal", uploadURLHandler)
	http.HandleFunc("/abc", urlFileUploadHandler)
	http.HandleFunc("/files.completeUploadExternal", completeURLUpload)
	once.Do(startServer)
	api := New("testing-token", OptionAPIURL("http://"+serverAddr+"/"))

	params := UploadFileParameters{
		Filename: "test.txt", Content: "test content", FileSize: 10,
		Channel: "CXXXXXXXX",
	}
	if _, err := api.UploadFile(params); err != nil {
		t.Errorf("Unexpected error: %s", err)
	}

	reader := bytes.NewBufferString("test reader")
	params = UploadFileParameters{
		Filename: "test.txt",
		Reader:   reader,
		FileSize: 10,
		Channel:  "CXXXXXXXX"}
	if _, err := api.UploadFile(params); err != nil {
		t.Errorf("Unexpected error: %s", err)
	}

	largeByt := make([]byte, 107374200)
	reader = bytes.NewBuffer(largeByt)
	params = UploadFileParameters{
		Filename: "test.txt", Reader: reader, FileSize: len(largeByt),
		Channel: "CXXXXXXXX"}
	if _, err := api.UploadFile(params); err != nil {
		t.Errorf("Unexpected error: %s", err)
	}

	reader = bytes.NewBufferString("test no channel")
	params = UploadFileParameters{
		Filename: "test.txt",
		Reader:   reader,
		FileSize: 15}
	if _, err := api.UploadFile(params); err != nil {
		t.Errorf("Unexpected error: %s", err)
	}
}

type mockGetUploadURLExternalHttpClient struct {
	ResponseStatus int
	ResponseBody   []byte
}

func (m *mockGetUploadURLExternalHttpClient) Do(req *http.Request) (*http.Response, error) {
	if req.URL.Path != "files.getUploadURLExternal" {
		return nil, fmt.Errorf("invalid path: %s", req.URL.Path)
	}

	return &http.Response{
		StatusCode: m.ResponseStatus,
		Body:       io.NopCloser(bytes.NewBuffer(m.ResponseBody)),
	}, nil
}

func TestGetUploadURLExternalContext(t *testing.T) {
	type testCase struct {
		title             string
		params            GetUploadURLExternalParameters
		wantSlackResponse []byte
		wantResponse      GetUploadURLExternalResponse
		wantErr           error
	}
	testCases := []testCase{
		{
			title: "Testing with required parameters",
			params: GetUploadURLExternalParameters{
				FileName: "test.txt",
				FileSize: 10,
			},
			wantSlackResponse: []byte(`{"ok":true,"file_id":"RandomID","upload_url":"http://test-server/abc"}`),
			wantResponse: GetUploadURLExternalResponse{
				FileID:    "RandomID",
				UploadURL: "http://test-server/abc",
				SlackResponse: SlackResponse{
					Ok: true,
				},
			},
		},
		{
			title: "Testing with optional parameters",
			params: GetUploadURLExternalParameters{
				FileSize:    10,
				FileName:    "test.txt",
				AltTxt:      "test-alt-text",
				SnippetType: "test-snippet-type",
			},
			wantSlackResponse: []byte(`{"ok":true,"file_id":"RandomID","upload_url":"http://test-server/abc"}`),
			wantResponse: GetUploadURLExternalResponse{
				FileID:    "RandomID",
				UploadURL: "http://test-server/abc",
				SlackResponse: SlackResponse{
					Ok: true,
				},
			},
		},
		{
			title: "Testing with request error",
			params: GetUploadURLExternalParameters{
				FileName: "test.txt",
				FileSize: 10,
			},
			wantSlackResponse: []byte(`{"ok":false,"error":"errored"}`),
			wantErr:           fmt.Errorf("errored"),
		},
		{
			title: "Testing with invalid parameters: empty file name",
			params: GetUploadURLExternalParameters{
				FileName: "",
				FileSize: 10,
			},
			wantErr: fmt.Errorf("FileName cannot be empty"),
		},
		{
			title: "Testing with invalid parameters: file size 0",
			params: GetUploadURLExternalParameters{
				FileName: "test.txt",
				FileSize: 0,
			},
			wantErr: fmt.Errorf("FileSize cannot be 0"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			api := &Client{
				token: validToken,
				httpclient: &mockGetUploadURLExternalHttpClient{
					ResponseStatus: 200,
					ResponseBody:   tc.wantSlackResponse,
				},
			}

			gotResponse, err := api.GetUploadURLExternalContext(context.Background(), tc.params)

			if err != nil {
				if tc.wantErr == nil {
					t.Fatalf("GetUploadURLExternalContext() error = %v, want nil", err)
				}
				if err.Error() != tc.wantErr.Error() {
					t.Errorf("GetUploadURLExternalContext() error = %v, want %v", err, tc.wantErr)
				}
			} else {
				if tc.wantErr != nil {
					t.Fatalf("GetUploadURLExternalContext() error = nil, want %v", tc.wantErr)
				}
				if !reflect.DeepEqual(gotResponse, &tc.wantResponse) {
					t.Errorf("GetUploadURLExternalContext() = %v, want %v", gotResponse, tc.wantResponse)
				}
			}
		})
	}
}

type mockCompleteUploadExternalHttpClient struct {
	ResponseStatus int
	ResponseBody   []byte
}

func (m *mockCompleteUploadExternalHttpClient) Do(req *http.Request) (*http.Response, error) {
	if req.URL.Path != "files.completeUploadExternal" {
		return nil, fmt.Errorf("invalid path: %s", req.URL.Path)
	}

	return &http.Response{
		StatusCode: m.ResponseStatus,
		Body:       io.NopCloser(bytes.NewBuffer(m.ResponseBody)),
	}, nil
}

func TestCompleteUploadExternalContext(t *testing.T) {
	type testCase struct {
		title        string
		params       CompleteUploadExternalParameters
		wantResponse CompleteUploadExternalResponse
		wantErr      bool
	}
	testCases := []testCase{
		{
			title: "Testing with required parameters",
			params: CompleteUploadExternalParameters{
				Files: []FileSummary{
					{
						ID: "ID1",
					},
					{
						ID: "ID2",
					},
				},
			},
			wantResponse: CompleteUploadExternalResponse{
				Files: []FileSummary{
					{
						ID: "ID1",
					},
					{
						ID: "ID2",
					},
				},
				SlackResponse: SlackResponse{Ok: true},
			},
		},
		{
			title: "Testing with optional parameters",
			params: CompleteUploadExternalParameters{
				Files: []FileSummary{
					{
						ID: "ID1",
					},
					{
						ID:    "ID2",
						Title: "Title2",
					},
				},
				Channel:         "test-channel",
				InitialComment:  "test-comment",
				ThreadTimestamp: "1234567890.123456",
			},
			wantResponse: CompleteUploadExternalResponse{
				Files: []FileSummary{
					{
						ID: "ID1",
					},
					{
						ID:    "ID2",
						Title: "Title2",
					},
				},
				SlackResponse: SlackResponse{Ok: true},
			},
		},
		{
			title: "Testing with blocks",
			params: CompleteUploadExternalParameters{
				Files: []FileSummary{
					{
						ID: "ID1",
					},
					{
						ID:    "ID2",
						Title: "Title2",
					},
				},
				Channel:         "test-channel",
				ThreadTimestamp: "1234567890.123456",
				Blocks: Blocks{BlockSet: []Block{
					NewSectionBlock(
						NewTextBlockObject("plain_text", "This is a section block", false, false), nil, nil),
				},
				},
			},
			wantResponse: CompleteUploadExternalResponse{
				Files: []FileSummary{
					{
						ID: "ID1",
					},
					{
						ID:    "ID2",
						Title: "Title2",
					},
				},
				SlackResponse: SlackResponse{Ok: true},
			},
		},
		{
			title: "Testing with error",
			params: CompleteUploadExternalParameters{
				Files: []FileSummary{
					{
						ID: "ID1",
					},
				},
			},
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			var resBody map[string]interface{}
			if !tc.wantErr {
				resBody = map[string]interface{}{
					"ok": true,
				}
				files := make([]map[string]string, 0)
				for _, file := range tc.params.Files {
					m := map[string]string{
						"id": file.ID,
					}
					if file.Title != "" {
						m["title"] = file.Title
					}
					files = append(files, m)
				}
				resBody["files"] = files
			} else {
				resBody = map[string]interface{}{
					"ok":    false,
					"error": "errored",
				}
			}

			resBodyBytes, err := json.Marshal(resBody)
			if err != nil {
				t.Fatalf("failed to marshal response body: %v", err)
			}

			api := &Client{
				token: validToken,
				httpclient: &mockCompleteUploadExternalHttpClient{
					ResponseStatus: 200,
					ResponseBody:   resBodyBytes,
				},
			}

			gotResponse, err := api.CompleteUploadExternalContext(context.Background(), tc.params)

			if err != nil {
				if !tc.wantErr {
					t.Errorf("CompleteUploadExternalContext() error = %v, want nil", err)
				}
			} else {
				if tc.wantErr {
					t.Fatalf("CompleteUploadExternalContext() error = nil, want %v", tc.wantErr)
				}
				if !reflect.DeepEqual(gotResponse, &tc.wantResponse) {
					t.Errorf("CompleteUploadExternalContext() = %v, want %v", gotResponse, tc.wantResponse)
				}
			}
		})
	}
}
