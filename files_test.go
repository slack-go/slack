package slack

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
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

// TestGetFileSignInRedirect replays files.slack.com (live check, 2026-10-04): with a
// token that cannot read the file it answers 302 to the workspace sign-in page with
// redir set to the file path; with access it serves the file, HTML files included.
func TestGetFileSignInRedirect(t *testing.T) {
	const (
		downloadPath = "/files-pri/T1-F1/download/todo.html"
		inlinePath   = "/files-pri/T1-F1/todo.html"
		page         = "<!DOCTYPE html><html><body>a real HTML file</body></html>"
	)
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authorized := r.Header.Get("Authorization") == "Bearer xoxb-valid"
		switch {
		case r.URL.Path == downloadPath && authorized:
			w.Header().Set("Content-Type", "application/force-download")
			_, _ = w.Write([]byte(page))
		case r.URL.Path == inlinePath && authorized:
			w.Header().Set("Content-Type", "text/html")
			_, _ = w.Write([]byte(page))
		case r.URL.Path == downloadPath || r.URL.Path == inlinePath:
			http.Redirect(w, r, "/?redir="+url.QueryEscape(r.URL.Path), http.StatusFound)
		case r.URL.Path == "/" && r.URL.Query().Get("redir") != "":
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			_, _ = w.Write([]byte("<!DOCTYPE html><html><body>Sign in to Slack</body></html>"))
		default:
			t.Errorf("unexpected request %s", r.URL)
		}
	}))
	defer ts.Close()

	for _, path := range []string{downloadPath, inlinePath} {
		t.Run("token can read "+path, func(t *testing.T) {
			var buf bytes.Buffer
			err := New("xoxb-valid", OptionAPIURL(ts.URL+"/")).GetFile(ts.URL+path, &buf)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if buf.String() != page {
				t.Fatalf("wrote %q, want the file", buf.String())
			}
		})
	}

	t.Run("token cannot read the file", func(t *testing.T) {
		var buf bytes.Buffer
		err := New("xoxb-invalid", OptionAPIURL(ts.URL+"/")).GetFile(ts.URL+downloadPath, &buf)
		if err == nil || !strings.Contains(err.Error(), "sign-in page") {
			t.Fatalf("expected a sign-in redirect error, got %v", err)
		}
		if buf.Len() != 0 {
			t.Fatalf("wrote %d bytes of the sign-in page", buf.Len())
		}
	})
}

// recordingHTTPClient answers every request with 200 and keeps the requests it got.
type recordingHTTPClient struct {
	requests []*http.Request
}

func (c *recordingHTTPClient) Do(req *http.Request) (*http.Response, error) {
	c.requests = append(c.requests, req)
	if req.Body != nil {
		// Let the multipart writer goroutine of UploadToURL finish.
		_, _ = io.Copy(io.Discard, req.Body)
	}
	return &http.Response{StatusCode: http.StatusOK, Body: io.NopCloser(strings.NewReader("OK")), Request: req}, nil
}

// TestTokenOnlySentToSlack covers GHSA-3q3v-34v2-g88f: GetFile and UploadToURL attach
// the token to the URL they are given, and a remote file's url_private is the
// external_url its creator chose. Only https Slack URLs and the API endpoint may get it.
func TestTokenOnlySentToSlack(t *testing.T) {
	const endpoint = "http://127.0.0.1:8080/"
	cases := []struct {
		url     string
		allowed bool
	}{
		{"https://files.slack.com/files-pri/T1-F1/a.txt", true},
		{"https://FILES.Slack.com/files-pri/T1-F1/a.txt", true},
		{"https://files.slack.com:443/files-pri/T1-F1/a.txt", true},
		{"https://acme.enterprise.slack.com/files-pri/T1-F1/a.txt", true},
		{"https://files.slack-gov.com/files-pri/T1-F1/a.txt", true},
		{endpoint + "files-pri/T1-F1/a.txt", true},
		{"https://evil.example/a.txt", false},
		{"https://files.slack.com@evil.example/a.txt", false},
		{"https://files.slack.com.evil.example/a.txt", false},
		{"https://evilslack.com/a.txt", false},
		{"http://files.slack.com/files-pri/T1-F1/a.txt", false},
		{"http://127.0.0.1:8081/files-pri/T1-F1/a.txt", false},
		{"files.slack.com/files-pri/T1-F1/a.txt", false},
	}
	calls := map[string]func(api *Client, u string) error{
		"GetFile": func(api *Client, u string) error {
			return api.GetFile(u, io.Discard)
		},
		"UploadToURL": func(api *Client, u string) error {
			return api.UploadToURL(context.Background(), UploadToURLParameters{UploadURL: u, Filename: "a.txt", Content: "a"})
		},
	}

	for name, call := range calls {
		for _, tc := range cases {
			t.Run(name+" "+tc.url, func(t *testing.T) {
				rc := &recordingHTTPClient{}
				api := New("xoxb-secret", OptionHTTPClient(rc), OptionAPIURL(endpoint))
				err := call(api, tc.url)
				if !tc.allowed {
					if err == nil {
						t.Error("expected an error")
					}
					for _, req := range rc.requests {
						t.Errorf("sent the token to %s", req.URL)
					}
					return
				}
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				if len(rc.requests) != 1 {
					t.Fatalf("sent %d requests, want 1", len(rc.requests))
				}
				if got := rc.requests[0].Header.Get("Authorization"); got != "Bearer xoxb-secret" {
					t.Errorf("Authorization = %q, want the token", got)
				}
			})
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
		{
			// A comment ID (Fc...) in the file position is the old argument order:
			// fail without calling Slack.
			title: "Testing with swapped IDs",
			body: url.Values{
				"file": {"Fc1234567890"},
				"id":   {"F1234567890"},
			},
			wantParams:  map[string]string{},
			expectError: true,
		},
	}

	var fch *fileCommentHandler
	http.HandleFunc("/files.comments.delete", func(w http.ResponseWriter, r *http.Request) {
		fch.handler(w, r)
	})

	for _, test := range tests {
		// Both functions take the file ID first (#1591).
		calls := []struct {
			name string
			call func() error
		}{
			{"DeleteFileComment", func() error {
				return api.DeleteFileComment(test.body["file"][0], test.body["id"][0])
			}},
			{"DeleteFileCommentContext", func() error {
				return api.DeleteFileCommentContext(context.Background(), test.body["file"][0], test.body["id"][0])
			}},
		}
		for _, c := range calls {
			fch = newFileCommentHandler()
			err := c.call()

			if !test.expectError && err != nil {
				log.Fatalf("%s: %s: Unexpected error: %s in test", c.name, test.title, err)
			} else if test.expectError == true && err == nil {
				log.Fatalf("%s: Expected error but got none", c.name)
			}

			if !reflect.DeepEqual(fch.gotParams, test.wantParams) {
				log.Fatalf("%s: %s: Got params [%#v]\nBut received [%#v]\n", c.name, test.title, fch.gotParams, test.wantParams)
			}
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
			title: "Testing with multiple channels",
			params: CompleteUploadExternalParameters{
				Files: []FileSummary{
					{
						ID: "ID1",
					},
				},
				Channels:       []string{"test-channel-1", "test-channel-2"},
				InitialComment: "test-comment",
			},
			wantResponse: CompleteUploadExternalResponse{
				Files: []FileSummary{
					{
						ID: "ID1",
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
			var resBody map[string]any
			if !tc.wantErr {
				resBody = map[string]any{
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
				resBody = map[string]any{
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

func TestDeleteFileCommentSwappedIDsError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error("files.comments.delete must not be called with swapped IDs")
	}))
	defer ts.Close()
	api := New("testing-token", OptionAPIURL(ts.URL+"/"))

	err := api.DeleteFileComment("Fc1234567890", "F1234567890")
	if err == nil || !strings.Contains(err.Error(), "fileID, commentID") {
		t.Fatalf("expected an error that names the argument order, got %v", err)
	}
}
