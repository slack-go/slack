package slack

import (
	"encoding/json"
	"net/http"
	"testing"
)

func workflowStepHandler(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	response, _ := json.Marshal(SlackResponse{
		Ok: true,
	})
	rw.Write(response)
}

func TestWorkflowStepCompleted(t *testing.T) {
	http.HandleFunc("/workflows.stepCompleted", workflowStepHandler)
	once.Do(startServer)
	api := New("testing-token", OptionAPIURL("http://"+serverAddr+"/"))

	if err := api.WorkflowStepCompleted("executeID"); err != nil {
		t.Errorf("Unexpected error: %s", err)
	}
}

func TestWorkflowStepFailed(t *testing.T) {
	http.HandleFunc("/workflows.stepFailed", workflowStepHandler)
	once.Do(startServer)
	api := New("testing-token", OptionAPIURL("http://"+serverAddr+"/"))

	if err := api.WorkflowStepFailed("executeID", "error message"); err != nil {
		t.Errorf("Unexpected error: %s", err)
	}
}
