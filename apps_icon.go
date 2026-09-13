package slack

import (
	"context"
	"errors"
	"io"
	"net/url"
)

// SetAppIconParameters contains arguments for SetAppIcon.
type SetAppIconParameters struct {
	AppID    string
	URL      string
	File     io.Reader
	Filename string
}

// SetAppIcon sets an app icon using a caller-supplied token.
// For more details, see SetAppIconContext documentation.
func (api *Client) SetAppIcon(token string, params SetAppIconParameters) error {
	return api.SetAppIconContext(context.Background(), token, params)
}

// SetAppIconContext sets an app icon using a caller-supplied token and a custom context.
// The token must have the app_configurations:write scope.
//
// Slack API docs: https://docs.slack.dev/reference/methods/apps.icon.set
func (api *Client) SetAppIconContext(ctx context.Context, token string, params SetAppIconParameters) error {
	if token == "" {
		return errors.New("apps.icon.set: token cannot be empty")
	}
	if params.AppID == "" {
		return errors.New("apps.icon.set: app ID cannot be empty")
	}
	if (params.URL == "") == (params.File == nil) {
		return errors.New("apps.icon.set: exactly one of URL or file must be provided")
	}

	values := url.Values{
		"app_id": {params.AppID},
	}
	response := &SlackResponse{}

	if params.URL != "" {
		values.Add("token", token)
		values.Add("url", params.URL)
		if err := api.postMethod(ctx, "apps.icon.set", values, response); err != nil {
			return err
		}
	} else {
		filename := params.Filename
		if filename == "" {
			filename = "icon.png"
		}
		if err := postWithMultipartResponse(ctx, api.httpclient, api.endpoint+"apps.icon.set", filename, "file", token, values, params.File, response, api); err != nil {
			return err
		}
	}

	return response.Err()
}
