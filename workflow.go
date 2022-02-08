package slack

import (
	"context"
	"encoding/json"
)

type WorkflowUpdateRequest struct {
	EditID  string                       `json:"workflow_step_edit_id"`
	Inputs  map[string]WorkflowStepInput `json:"inputs,omitempty"`
	Outputs []WorkflowStepOutput         `json:"outputs,omitempty"`
}

type WorkflowStepInput struct {
	Value                   string                 `json:"value,omitempty"`
	SkipVariableReplacement bool                   `json:"skip_variable_replacement,omitempty"`
	Variables               map[string]interface{} `json:"variables,omitempty"`
}

type WorkflowStepOutput struct {
	Name  string `json:"name,omitempty"`
	Type  string `json:"type,omitempty"`
	Label string `json:"label,omitempty"`
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
