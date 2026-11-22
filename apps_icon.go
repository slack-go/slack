package slack

import (
	"cmp"
	"context"
	"errors"
	"io"
	"net/url"
)

// SetAppIconParameters contains the image for SetAppIcon. Set exactly one of
// URL or File.
type SetAppIconParameters struct {
	URL      string
	File     io.Reader
	Filename string
}

// SetAppIcon sets the icon of an app.
// For more details, see SetAppIconContext documentation.
func (api *Client) SetAppIcon(token, appID string, icon SetAppIconParameters) error {
	return api.SetAppIconContext(context.Background(), token, appID, icon)
}

// SetAppIconContext sets the icon of an app with a custom context. An empty token
// falls back to the configuration token set with OptionConfigToken.
//
// Slack API docs: https://docs.slack.dev/reference/methods/apps.icon.set
func (api *Client) SetAppIconContext(ctx context.Context, token, appID string, icon SetAppIconParameters) error {
	token = cmp.Or(token, api.configToken)
	switch {
	case token == "":
		return errors.New("apps.icon.set: token cannot be empty")
	case appID == "":
		return errors.New("apps.icon.set: app ID cannot be empty")
	case (icon.URL == "") == (icon.File == nil):
		return errors.New("apps.icon.set: exactly one of URL or file must be provided")
	}

	response := &SlackResponse{}
	var err error
	if icon.URL != "" {
		values := url.Values{"token": {token}, "app_id": {appID}, "url": {icon.URL}}
		err = api.postMethod(ctx, "apps.icon.set", values, response)
	} else {
		filename := cmp.Or(icon.Filename, "icon.png")
		values := url.Values{"app_id": {appID}}
		err = api.postMultipartMethod(ctx, "apps.icon.set", filename, "file", token, values, icon.File, response)
	}
	if err != nil {
		return err
	}
	return response.Err()
}
