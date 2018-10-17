package slack

import (
	"context"
	"errors"
	"net/url"
)

// APIResponseFull contains our API response from the api.test endpoint
type APIResponseFull struct {
	SlackResponse
	Args map[string]string `json:"args,omitempty"`
}

func apiRequest(ctx context.Context, client httpClient, path string, values url.Values, debug debug) (*APIResponseFull, error) {
	response := &APIResponseFull{}
	err := postSlackMethod(ctx, client, path, values, response, debug)
	if err != nil {
		return nil, err
	}
	if !response.Ok {
		return nil, errors.New(response.Error)
	}
	return response, nil
}

// GetApiTest will test the API by hitting the api.test endpoint, and
// returning that response. It is a good way to make sure you're authenticated,
// and talking to slack.
func (api *Client) GetApiTest() (*APIResponseFull, error) {
	return api.GetApiTestContext(context.Background())
}

// GetApiTestContext will retrieve the satus from api.test
func (api *Client) GetApiTestContext(ctx context.Context) (*APIResponseFull, error) {
	values := url.Values{
		"token": {api.token},
	}

	response, err := apiRequest(ctx, api.httpclient, "api.test", values, api)
	if err != nil {
		return nil, err
	}
	return response, nil
}
