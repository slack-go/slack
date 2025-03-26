package slack

import (
	"encoding/json"
	"net/http"
	"reflect"
	"testing"
)

func createCanvasHandler(rw http.ResponseWriter, r *http.Request) {
	title := r.FormValue("title")
	documentContent := r.FormValue("document_content")

	rw.Header().Set("Content-Type", "application/json")

	if title != "" && documentContent != "" {
		resp, _ := json.Marshal(&struct {
			SlackResponse
			CanvasID string `json:"canvas_id"`
		}{
			SlackResponse: SlackResponse{Ok: true},
			CanvasID:      "F1234ABCD",
		})
		rw.Write(resp)
	} else {
		rw.Write([]byte(`{ "ok": false, "error": "errored" }`))
	}
}

func TestCreateCanvas(t *testing.T) {
	s := startServer()
	defer s.Close()

	s.RegisterHandler("/canvases.create", createCanvasHandler)
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
	s := startServer()
	defer s.Close()

	s.RegisterHandler("/canvases.delete", deleteCanvasHandler)
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
	s := startServer()
	defer s.Close()

	s.RegisterHandler("/canvases.edit", editCanvasHandler)
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
	s := startServer()
	defer s.Close()

	s.RegisterHandler("/canvases.access.set", setCanvasAccessHandler)
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
	s := startServer()
	defer s.Close()

	s.RegisterHandler("/canvases.access.delete", deleteCanvasAccessHandler)

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
	s := startServer()
	defer s.Close()

	s.RegisterHandler("/canvases.sections.lookup", lookupCanvasSectionsHandler)

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
