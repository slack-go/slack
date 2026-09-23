package slack

import (
	"context"
	"net/url"
)

// Lifecycle statuses accepted by agents.sessions.setStatus.
const (
	AgentSessionStatusActive     = "active"
	AgentSessionStatusProcessing = "processing"
	AgentSessionStatusSuspended  = "suspended"
	AgentSessionStatusClosed     = "closed"
)

// AgentSessionSetStatusParameters are the parameters for SetAgentSessionStatus.
// Title and InitiatorUserID only apply when the call creates the session.
//
// IconEmoji, IconURL and Username form one set of overrides that persists across
// calls. A call that leaves all three empty keeps the current overrides. A call
// that sets any of them replaces the whole set, so a field left empty is cleared.
type AgentSessionSetStatusParameters struct {
	Status          string `json:"status"`
	ChannelID       string `json:"channel_id,omitempty"`
	ThreadTS        string `json:"thread_ts,omitempty"`
	Title           string `json:"title,omitempty"`
	InitiatorUserID string `json:"initiator_user_id,omitempty"`
	IconEmoji       string `json:"icon_emoji,omitempty"`
	IconURL         string `json:"icon_url,omitempty"`
	Username        string `json:"username,omitempty"`
}

// AgentSessionSetStatusResponse is the response from agents.sessions.setStatus.
// Status covers every agent in the session; AgentStatus is the caller's own.
type AgentSessionSetStatusResponse struct {
	SlackResponse
	Status      string `json:"status"`
	AgentStatus string `json:"agent_status"`
	Title       string `json:"title,omitempty"`
}

// AgentSessionRenameParameters are the parameters for RenameAgentSession.
type AgentSessionRenameParameters struct {
	Title     string `json:"title"`
	ChannelID string `json:"channel_id,omitempty"`
	ThreadTS  string `json:"thread_ts,omitempty"`
}

// AgentSessionRenameResponse is the response from agents.sessions.rename.
type AgentSessionRenameResponse struct {
	SlackResponse
	Title string `json:"title"`
}

// SetAgentSessionStatus sets an agent session's lifecycle status, creating the session if needed.
//
// Slack API docs: https://docs.slack.dev/reference/methods/agents.sessions.setStatus
func (api *Client) SetAgentSessionStatus(params AgentSessionSetStatusParameters) (*AgentSessionSetStatusResponse, error) {
	return api.SetAgentSessionStatusContext(context.Background(), params)
}

// SetAgentSessionStatusContext sets an agent session's lifecycle status with a custom context.
//
// Slack API docs: https://docs.slack.dev/reference/methods/agents.sessions.setStatus
func (api *Client) SetAgentSessionStatusContext(ctx context.Context, params AgentSessionSetStatusParameters) (*AgentSessionSetStatusResponse, error) {
	values := url.Values{
		"token":  {api.token},
		"status": {params.Status},
	}

	if params.ChannelID != "" {
		values.Add("channel_id", params.ChannelID)
	}

	if params.ThreadTS != "" {
		values.Add("thread_ts", params.ThreadTS)
	}

	if params.Title != "" {
		values.Add("title", params.Title)
	}

	if params.InitiatorUserID != "" {
		values.Add("initiator_user_id", params.InitiatorUserID)
	}

	if params.IconEmoji != "" {
		values.Add("icon_emoji", params.IconEmoji)
	}

	if params.IconURL != "" {
		values.Add("icon_url", params.IconURL)
	}

	if params.Username != "" {
		values.Add("username", params.Username)
	}

	response := &AgentSessionSetStatusResponse{}

	err := api.postMethod(ctx, "agents.sessions.setStatus", values, response)
	if err != nil {
		return nil, err
	}

	return response, response.Err()
}

// RenameAgentSession sets the title of an agent session.
//
// Slack API docs: https://docs.slack.dev/reference/methods/agents.sessions.rename
func (api *Client) RenameAgentSession(params AgentSessionRenameParameters) (*AgentSessionRenameResponse, error) {
	return api.RenameAgentSessionContext(context.Background(), params)
}

// RenameAgentSessionContext sets the title of an agent session with a custom context.
//
// Slack API docs: https://docs.slack.dev/reference/methods/agents.sessions.rename
func (api *Client) RenameAgentSessionContext(ctx context.Context, params AgentSessionRenameParameters) (*AgentSessionRenameResponse, error) {
	values := url.Values{
		"token": {api.token},
		"title": {params.Title},
	}

	if params.ChannelID != "" {
		values.Add("channel_id", params.ChannelID)
	}

	if params.ThreadTS != "" {
		values.Add("thread_ts", params.ThreadTS)
	}

	response := &AgentSessionRenameResponse{}

	err := api.postMethod(ctx, "agents.sessions.rename", values, response)
	if err != nil {
		return nil, err
	}

	return response, response.Err()
}
