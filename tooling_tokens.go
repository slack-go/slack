package slack

import (
	"context"
	"net/url"
)

// ToolingTokensRotate contains our Auth response from the tooling.tokens.rotate endpoint
type ToolingTokensRotate struct {
	SlackResponse        // Contains the "ok", and "Error", if any
	Token         string `json:"token"`
	RefreshToken  string `json:"refresh_token"`
	TeamID        string `json:"team_id"`
	UserID        string `json:"user_id"`
	Iat           int64  `json:"iat"`
	Exp           int64  `json:"exp"`
}

// ToolingTokensRotate will send a refresh for our token
func (api *Client) ToolingTokensRotate(refresh_token string) (*ToolingTokensRotate, error) {
	return api.ToolingTokensRotateContext(context.Background(), refresh_token)
}

// ToolingTokensRotateContext will send a refresh request for our token
func (api *Client) ToolingTokensRotateContext(ctx context.Context, refresh_token string) (*ToolingTokensRotate, error) {
	response := &ToolingTokensRotate{}
	err := api.postMethod(ctx, "tooling.tokens.rotate", url.Values{"refresh_token": {refresh_token}}, response)
	if err != nil {
		return nil, err
	}

	return response, response.Err()
}
