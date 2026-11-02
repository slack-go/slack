package slack

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func createCanvasHandler(rw http.ResponseWriter, r *http.Request) {
	title := r.FormValue("title")
	documentContent := r.FormValue("document_content")
	channelID := r.FormValue("channel_id")

	rw.Header().Set("Content-Type", "application/json")

	if title != "" && documentContent != "" {
		// Channel canvases get their own ID, so the tests can tell whether
		// the expected channel_id made it into the request.
		canvasID := "F1234ABCD"
		if channelID == "C1234ABCD" {
			canvasID = "F5678EFGH"
		}

		resp, _ := json.Marshal(&struct {
			SlackResponse
			CanvasID string `json:"canvas_id"`
		}{
			SlackResponse: SlackResponse{Ok: true},
			CanvasID:      canvasID,
		})
		rw.Write(resp)
	} else {
		rw.Write([]byte(`{ "ok": false, "error": "errored" }`))
	}
}

func TestCreateCanvas(t *testing.T) {
	http.HandleFunc("/canvases.create", createCanvasHandler)
	once.Do(startServer)
	api := New("testing-token", OptionAPIURL("http://"+serverAddr+"/"))

	documentContent := DocumentContent{
		Type:     "markdown",
		Markdown: "Test Content",
	}

	canvasID, err := api.CreateCanvas("Test Canvas", documentContent)
	if err != nil {
		t.Fatalf("Unexpected error: %s", err)
	}

	if canvasID != "F1234ABCD" {
		t.Fatalf("Expected canvas ID to be F1234ABCD, got %s", canvasID)
	}
}

func TestCreateCanvasWithChannelID(t *testing.T) {
	http.DefaultServeMux = new(http.ServeMux)
	http.HandleFunc("/canvases.create", createCanvasHandler)
	once.Do(startServer)
	api := New("testing-token", OptionAPIURL("http://"+serverAddr+"/"))

	documentContent := DocumentContent{
		Type:     "markdown",
		Markdown: "Test Content",
	}

	canvasID, err := api.CreateCanvas("Test Canvas", documentContent, CreateCanvasOptionChannelID("C1234ABCD"))
	if err != nil {
		t.Fatalf("Unexpected error: %s", err)
	}

	if canvasID != "F5678EFGH" {
		t.Fatalf("Expected canvas ID to be F5678EFGH, got %s", canvasID)
	}
}

func deleteCanvasHandler(rw http.ResponseWriter, r *http.Request) {
	canvasID := r.FormValue("canvas_id")

	rw.Header().Set("Content-Type", "application/json")

	if canvasID == "F1234ABCD" {
		rw.Write([]byte(`{ "ok": true }`))
	} else {
		rw.Write([]byte(`{ "ok": false, "error": "errored" }`))
	}
}

func TestDeleteCanvas(t *testing.T) {
	http.HandleFunc("/canvases.delete", deleteCanvasHandler)
	once.Do(startServer)
	api := New("testing-token", OptionAPIURL("http://"+serverAddr+"/"))

	err := api.DeleteCanvas("F1234ABCD")
	if err != nil {
		t.Fatalf("Unexpected error: %s", err)
	}
}

func editCanvasHandler(rw http.ResponseWriter, r *http.Request) {
	canvasID := r.FormValue("canvas_id")

	rw.Header().Set("Content-Type", "application/json")

	if canvasID == "F1234ABCD" {
		rw.Write([]byte(`{ "ok": true }`))
	} else {
		rw.Write([]byte(`{ "ok": false, "error": "errored" }`))
	}
}

func TestEditCanvas(t *testing.T) {
	http.HandleFunc("/canvases.edit", editCanvasHandler)
	once.Do(startServer)
	api := New("testing-token", OptionAPIURL("http://"+serverAddr+"/"))

	params := EditCanvasParams{
		CanvasID: "F1234ABCD",
		Changes: []CanvasChange{
			{
				Operation: "update",
				SectionID: "S1234",
				DocumentContent: DocumentContent{
					Type:     "markdown",
					Markdown: "Updated Content",
				},
			},
		},
	}

	err := api.EditCanvas(params)
	if err != nil {
		t.Fatalf("Unexpected error: %s", err)
	}
}

func setCanvasAccessHandler(rw http.ResponseWriter, r *http.Request) {
	canvasID := r.FormValue("canvas_id")

	rw.Header().Set("Content-Type", "application/json")

	if canvasID == "F1234ABCD" {
		rw.Write([]byte(`{ "ok": true }`))
	} else {
		rw.Write([]byte(`{ "ok": false, "error": "errored" }`))
	}
}

func TestSetCanvasAccess(t *testing.T) {
	http.HandleFunc("/canvases.access.set", setCanvasAccessHandler)
	once.Do(startServer)
	api := New("testing-token", OptionAPIURL("http://"+serverAddr+"/"))

	params := SetCanvasAccessParams{
		CanvasID:    "F1234ABCD",
		AccessLevel: "read",
		ChannelIDs:  []string{"C1234ABCD"},
		UserIDs:     []string{"U1234ABCD"},
	}

	err := api.SetCanvasAccess(params)
	if err != nil {
		t.Fatalf("Unexpected error: %s", err)
	}
}

func deleteCanvasAccessHandler(rw http.ResponseWriter, r *http.Request) {
	canvasID := r.FormValue("canvas_id")

	rw.Header().Set("Content-Type", "application/json")

	if canvasID == "F1234ABCD" {
		rw.Write([]byte(`{ "ok": true }`))
	} else {
		rw.Write([]byte(`{ "ok": false, "error": "errored" }`))
	}
}

func TestDeleteCanvasAccess(t *testing.T) {
	http.HandleFunc("/canvases.access.delete", deleteCanvasAccessHandler)
	once.Do(startServer)
	api := New("testing-token", OptionAPIURL("http://"+serverAddr+"/"))

	params := DeleteCanvasAccessParams{
		CanvasID:   "F1234ABCD",
		ChannelIDs: []string{"C1234ABCD"},
		UserIDs:    []string{"U1234ABCD"},
	}

	err := api.DeleteCanvasAccess(params)
	if err != nil {
		t.Fatalf("Unexpected error: %s", err)
	}
}

func lookupCanvasSectionsHandler(rw http.ResponseWriter, r *http.Request) {
	canvasID := r.FormValue("canvas_id")

	rw.Header().Set("Content-Type", "application/json")

	if canvasID == "F1234ABCD" {
		sections := []CanvasSection{
			{ID: "S1234"},
			{ID: "S5678"},
		}

		resp, _ := json.Marshal(&LookupCanvasSectionsResponse{
			SlackResponse: SlackResponse{Ok: true},
			Sections:      sections,
		})
		rw.Write(resp)
	} else {
		rw.Write([]byte(`{ "ok": false, "error": "errored" }`))
	}
}

func TestLookupCanvasSections(t *testing.T) {
	http.HandleFunc("/canvases.sections.lookup", lookupCanvasSectionsHandler)
	once.Do(startServer)
	api := New("testing-token", OptionAPIURL("http://"+serverAddr+"/"))

	params := LookupCanvasSectionsParams{
		CanvasID: "F1234ABCD",
		Criteria: LookupCanvasSectionsCriteria{
			SectionTypes: []string{"h1", "h2"},
			ContainsText: "Test",
		},
	}

	sections, err := api.LookupCanvasSections(params)
	if err != nil {
		t.Fatalf("Unexpected error: %s", err)
	}

	expectedSections := []CanvasSection{
		{ID: "S1234"},
		{ID: "S5678"},
	}

	if !reflect.DeepEqual(expectedSections, sections) {
		t.Fatalf("Expected sections %v, got %v", expectedSections, sections)
	}
}

func TestGetCanvasContent(t *testing.T) {
	for _, tt := range []struct {
		name    string
		options []GetCanvasContentOption
		want    url.Values
		content string
	}{
		{"default", nil, url.Values{}, "# Plan\n\n- Item"},
		{"markdown", []GetCanvasContentOption{GetCanvasContentOptionContentType("markdown")}, url.Values{"content_type": {"markdown"}}, "# Plan"},
		{"html", []GetCanvasContentOption{GetCanvasContentOptionContentType("html")}, url.Values{"content_type": {"html"}}, "<h1>Plan</h1>"},
		{"nil and empty", []GetCanvasContentOption{nil, GetCanvasContentOptionContentType("")}, url.Values{}, ""},
		{"last option wins", []GetCanvasContentOption{GetCanvasContentOptionContentType("markdown"), GetCanvasContentOptionContentType("html")}, url.Values{"content_type": {"html"}}, "<p>Plan</p>"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, http.MethodPost, r.Method)
				assert.Equal(t, "/canvases.getContent", r.URL.Path)
				if !assert.NoError(t, r.ParseForm()) {
					return
				}
				assert.Equal(t, "testing-token", r.PostForm.Get("token"))
				assert.Equal(t, "F123", r.PostForm.Get("canvas_id"))
				r.PostForm.Del("token")
				r.PostForm.Del("canvas_id")
				assert.Equal(t, tt.want, r.PostForm)
				w.Header().Set("Content-Type", "application/json")
				_ = json.NewEncoder(w).Encode(map[string]any{"ok": true, "content": tt.content})
			}))
			defer server.Close()
			api := New("testing-token", OptionAPIURL(server.URL+"/"))
			content, err := api.GetCanvasContent("F123", tt.options...)
			require.NoError(t, err)
			assert.Equal(t, tt.content, content)
		})
	}
}

func TestGetCanvasContentErrors(t *testing.T) {
	for _, method := range []struct {
		name string
		call func(*Client, context.Context) error
	}{
		{"canvas content", func(api *Client, ctx context.Context) error {
			_, err := api.GetCanvasContentContext(ctx, "F123", nil)
			return err
		}},
	} {
		t.Run(method.name, func(t *testing.T) {
			for _, tt := range []struct {
				name, body string
				status     int
			}{
				{"Slack error", `{"ok":false,"error":"access_denied","response_metadata":{"messages":["Request rejected"]}}`, http.StatusOK},
				{"invalid JSON", `{"ok":`, http.StatusOK},
				{"HTTP error", `server error`, http.StatusInternalServerError},
			} {
				t.Run(tt.name, func(t *testing.T) {
					server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
						w.WriteHeader(tt.status)
						_, _ = w.Write([]byte(tt.body))
					}))
					defer server.Close()
					api := New("testing-token", OptionAPIURL(server.URL+"/"))
					err := method.call(api, context.Background())
					require.Error(t, err)
					if tt.name == "Slack error" {
						var slackErr SlackErrorResponse
						require.ErrorAs(t, err, &slackErr)
						assert.Equal(t, "access_denied", slackErr.Err)
						assert.Equal(t, []string{"Request rejected"}, slackErr.ResponseMetadata.Messages)
					}
				})
			}
			t.Run("canceled context", func(t *testing.T) {
				server := httptest.NewServer(http.HandlerFunc(func(http.ResponseWriter, *http.Request) { t.Error("canceled request reached server") }))
				defer server.Close()
				api := New("testing-token", OptionAPIURL(server.URL+"/"))
				ctx, cancel := context.WithCancel(context.Background())
				cancel()
				require.ErrorIs(t, method.call(api, ctx), context.Canceled)
			})
		})
	}
}
