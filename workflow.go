package slack

import (
	"context"
	"encoding/json"
)

type WorkflowUpdateRequest struct {
	EditID  string                      `json:"workflow_step_edit_id"`
	Inputs  *WorkflowStepInputRequest   `json:"inputs,omitempty"`
	Outputs []WorkflowStepOutputRequest `json:"outputs,omitempty"`
}

type WorkflowStepInputRequest struct {
}

type WorkflowStepOutputRequest struct {
}

func (api *Client) UpdateWorkflow(ctx context.Context, req WorkflowUpdateRequest) (*SlackResponse, error) {
	if req.EditID == "" {
		return nil, ErrParametersMissing
	}

	encoded, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	endpoint := api.endpoint + "workflows.updateStep"
	resp := &SlackResponse{}
	if err := postJSON(ctx, api.httpclient, endpoint, api.token, encoded, resp, api); err != nil {
		return nil, err
	}
	return resp, resp.Err()
}
