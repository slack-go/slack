package slack

import (
	"context"
	"net/http"
	"reflect"
	"testing"
)

type workflowsFeaturedHandler struct {
	rawResponse string
}

func (h *workflowsFeaturedHandler) handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(h.rawResponse))
}

func TestSlack_WorkflowsFeaturedAdd(t *testing.T) {
	once.Do(startServer)
	api := New("testing-token", OptionAPIURL("http://"+serverAddr+"/"))

	cases := []struct {
		caseName    string
		input       *WorkflowsFeaturedAddInput
		rawResp     string
		expectedErr error
	}{
		{
			caseName: "success",
			input: &WorkflowsFeaturedAddInput{
				ChannelID:  "C012345678",
				TriggerIDs: []string{"Ft1234", "Ft5678"},
			},
			rawResp:     `{"ok": true}`,
			expectedErr: nil,
		},
	}

	h := &workflowsFeaturedHandler{}
	http.HandleFunc("/workflows.featured.add", h.handler)
	for _, c := range cases {
		t.Run(c.caseName, func(t *testing.T) {
			h.rawResponse = c.rawResp

			err := api.WorkflowsFeaturedAdd(context.Background(), c.input)
			if c.expectedErr == nil && err != nil {
				t.Fatalf("unexpected error: %s\n", err)
			}
			if c.expectedErr != nil && err == nil {
				t.Fatalf("expected %s, but did not raise an error", c.expectedErr)
			}
			if c.expectedErr != nil && err != nil && c.expectedErr.Error() != err.Error() {
				t.Fatalf("expected %s as error but got %s\n", c.expectedErr, err)
			}
		})
	}
}

func TestSlack_WorkflowsFeaturedList(t *testing.T) {
	once.Do(startServer)
	api := New("testing-token", OptionAPIURL("http://"+serverAddr+"/"))

	cases := []struct {
		caseName     string
		input        *WorkflowsFeaturedListInput
		rawResp      string
		expectedResp *WorkflowsFeaturedListOutput
		expectedErr  error
	}{
		{
			caseName: "success",
			input: &WorkflowsFeaturedListInput{
				ChannelIDs: []string{"C012345678", "C987654321"},
			},
			rawResp: `{
				"ok": true,
				"featured_workflows": [
					{
						"channel_id": "C012345678",
						"triggers": [
							{"id": "Ft1234", "title": "Tabby workflow"},
							{"id": "Ft5678", "title": "Tortoise workflow"}
						]
					},
					{
						"channel_id": "C987654321",
						"triggers": [
							{"id": "Ft1234", "title": "Ragdoll workflow"}
						]
					}
				]
			}`,
			expectedResp: &WorkflowsFeaturedListOutput{
				FeaturedWorkflows: []FeaturedWorkflow{
					{
						ChannelID: "C012345678",
						Triggers: []FeaturedWorkflowTrigger{
							{ID: "Ft1234", Title: "Tabby workflow"},
							{ID: "Ft5678", Title: "Tortoise workflow"},
						},
					},
					{
						ChannelID: "C987654321",
						Triggers: []FeaturedWorkflowTrigger{
							{ID: "Ft1234", Title: "Ragdoll workflow"},
						},
					},
				},
			},
			expectedErr: nil,
		},
	}

	h := &workflowsFeaturedHandler{}
	http.HandleFunc("/workflows.featured.list", h.handler)
	for _, c := range cases {
		t.Run(c.caseName, func(t *testing.T) {
			h.rawResponse = c.rawResp

			resp, err := api.WorkflowsFeaturedList(context.Background(), c.input)
			if c.expectedErr == nil && err != nil {
				t.Fatalf("unexpected error: %s\n", err)
			}
			if c.expectedErr != nil && err == nil {
				t.Fatalf("expected %s, but did not raise an error", c.expectedErr)
			}
			if c.expectedErr != nil && err != nil && c.expectedErr.Error() != err.Error() {
				t.Fatalf("expected %s as error but got %s\n", c.expectedErr, err)
			}
			if resp == nil || c.expectedResp == nil {
				return
			}
			if !reflect.DeepEqual(resp, c.expectedResp) {
				t.Fatalf("expected:\n\t%v\n but got:\n\t%v\n", c.expectedResp, resp)
			}
		})
	}
}

func TestSlack_WorkflowsFeaturedRemove(t *testing.T) {
	once.Do(startServer)
	api := New("testing-token", OptionAPIURL("http://"+serverAddr+"/"))

	cases := []struct {
		caseName    string
		input       *WorkflowsFeaturedRemoveInput
		rawResp     string
		expectedErr error
	}{
		{
			caseName: "success",
			input: &WorkflowsFeaturedRemoveInput{
				ChannelID:  "C012345678",
				TriggerIDs: []string{"Ft1234"},
			},
			rawResp:     `{"ok": true}`,
			expectedErr: nil,
		},
	}

	h := &workflowsFeaturedHandler{}
	http.HandleFunc("/workflows.featured.remove", h.handler)
	for _, c := range cases {
		t.Run(c.caseName, func(t *testing.T) {
			h.rawResponse = c.rawResp

			err := api.WorkflowsFeaturedRemove(context.Background(), c.input)
			if c.expectedErr == nil && err != nil {
				t.Fatalf("unexpected error: %s\n", err)
			}
			if c.expectedErr != nil && err == nil {
				t.Fatalf("expected %s, but did not raise an error", c.expectedErr)
			}
			if c.expectedErr != nil && err != nil && c.expectedErr.Error() != err.Error() {
				t.Fatalf("expected %s as error but got %s\n", c.expectedErr, err)
			}
		})
	}
}

func TestSlack_WorkflowsFeaturedSet(t *testing.T) {
	once.Do(startServer)
	api := New("testing-token", OptionAPIURL("http://"+serverAddr+"/"))

	cases := []struct {
		caseName    string
		input       *WorkflowsFeaturedSetInput
		rawResp     string
		expectedErr error
	}{
		{
			caseName: "success",
			input: &WorkflowsFeaturedSetInput{
				ChannelID:  "C012345678",
				TriggerIDs: []string{"Ft1234", "Ft5678"},
			},
			rawResp:     `{"ok": true}`,
			expectedErr: nil,
		},
	}

	h := &workflowsFeaturedHandler{}
	http.HandleFunc("/workflows.featured.set", h.handler)
	for _, c := range cases {
		t.Run(c.caseName, func(t *testing.T) {
			h.rawResponse = c.rawResp

			err := api.WorkflowsFeaturedSet(context.Background(), c.input)
			if c.expectedErr == nil && err != nil {
				t.Fatalf("unexpected error: %s\n", err)
			}
			if c.expectedErr != nil && err == nil {
				t.Fatalf("expected %s, but did not raise an error", c.expectedErr)
			}
			if c.expectedErr != nil && err != nil && c.expectedErr.Error() != err.Error() {
				t.Fatalf("expected %s as error but got %s\n", c.expectedErr, err)
			}
		})
	}
}
