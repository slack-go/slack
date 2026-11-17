package slack

import (
	"context"
	"encoding/json"
)

type (
	FunctionCompleteSuccessRequest struct {
		FunctionExecutionID string            `json:"function_execution_id"`
		Outputs             map[string]string `json:"outputs"`
	}

	FunctionCompleteErrorRequest struct {
		FunctionExecutionID string `json:"function_execution_id"`
		Error               string `json:"error"`
	}
)

type FunctionCompleteSuccessRequestOption func(opt *FunctionCompleteSuccessRequest) error

func FunctionCompleteSuccessRequestOptionOutput(outputs map[string]string) FunctionCompleteSuccessRequestOption {
	return func(opt *FunctionCompleteSuccessRequest) error {
		if len(outputs) > 0 {
			opt.Outputs = outputs
		}
		return nil
	}
}

// FunctionCompleteSuccess indicates function is completed
func (api *Client) FunctionCompleteSuccess(functionExecutionId string, options ...FunctionCompleteSuccessRequestOption) error {
	return api.FunctionCompleteSuccessContext(context.Background(), functionExecutionId, options...)
}

// FunctionCompleteSuccessContext indicates successful function completion with a custom context.
//
// Slack API docs: https://docs.slack.dev/reference/methods/functions.completeSuccess
func (api *Client) FunctionCompleteSuccessContext(ctx context.Context, functionExecutionId string, options ...FunctionCompleteSuccessRequestOption) error {
	r := &FunctionCompleteSuccessRequest{
		FunctionExecutionID: functionExecutionId,
	}
	for _, option := range options {
		if err := option(r); err != nil {
			return err
		}
	}
	if r.Outputs == nil {
		r.Outputs = map[string]string{}
	}

	jsonData, err := json.Marshal(r)
	if err != nil {
		return err
	}

	response := &SlackResponse{}
	if err := api.postJSONMethod(ctx, "functions.completeSuccess", api.token, jsonData, response); err != nil {
		return err
	}

	if !response.Ok {
		return response.Err()
	}

	return nil
}

// FunctionCompleteError indicates function is completed with error
func (api *Client) FunctionCompleteError(functionExecutionID string, errorMessage string) error {
	return api.FunctionCompleteErrorContext(context.Background(), functionExecutionID, errorMessage)
}

// FunctionCompleteErrorContext indicates function is completed with error
func (api *Client) FunctionCompleteErrorContext(ctx context.Context, functionExecutionID string, errorMessage string) error {
	// More information: https://api.slack.com/methods/functions.completeError
	r := FunctionCompleteErrorRequest{
		FunctionExecutionID: functionExecutionID,
	}
	r.Error = errorMessage

	jsonData, err := json.Marshal(r)
	if err != nil {
		return err
	}

	response := &SlackResponse{}
	if err := api.postJSONMethod(ctx, "functions.completeError", api.token, jsonData, response); err != nil {
		return err
	}

	if !response.Ok {
		return response.Err()
	}

	return nil
}
