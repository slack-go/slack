package slack

import (
	"context"
)

// SocketModeConnection contains various details about the SocketMode connection.
// It is returned by an "apps.connections.open" API call.
type SocketModeConnection struct {
	URL  string                 `json:"url,omitempty"`
	Data map[string]interface{} `json:"-"`
}

type openResponseFull struct {
	SlackResponse
	SocketModeConnection
}

// StartSocketModeContext calls the "apps.connections.open" endpoint and returns the provided URL and the full Info block with a custom context.
//
// To have a fully managed Websocket connection, use `NewRTM`, and call `Run()` on it.
func (api *Client) StartSocketModeContext(ctx context.Context) (info *SocketModeConnection, websocketURL string, err error) {
	response := &openResponseFull{}
	err = api.postMethod(ctx, "apps.connections.open", url.Values{"token": {api.token}}, response)
	if err != nil {
		return nil, "", err
	}

	api.Debugln("Using URL:", response.SocketModeConnection.URL)
	return &response.SocketModeConnection, response.SocketModeConnection.URL, response.Err()
}
