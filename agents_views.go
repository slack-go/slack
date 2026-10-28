package slack

import (
	"context"
	"encoding/json"
	"net/url"
)

// View types accepted by agents.conversations.setView.
const (
	AgentViewTypeHTML        = "html"
	AgentViewTypeDiff        = "diff"
	AgentViewTypeBlockKit    = "block_kit"
	AgentViewTypeCanvas      = "canvas"
	AgentViewTypePullRequest = "pull_request"
)

// AgentViewCSP declares HTTPS origins for a code channel's HTML view.
// Slack validates the origins. ConnectDomains is stored but is not yet honored
// when rendering the view.
//
// Slack API docs: https://docs.slack.dev/reference/methods/agents.conversations.setView
type AgentViewCSP struct {
	ResourceDomains []string `json:"resource_domains,omitempty"`
	ConnectDomains  []string `json:"connect_domains,omitempty"`
}

type setAgentConversationViewParams struct {
	values url.Values
	blocks []Block
	csp    *AgentViewCSP
}

// SetAgentConversationViewOption configures a code channel view creation or update.
type SetAgentConversationViewOption func(*setAgentConversationViewParams)

// SetAgentConversationViewOptionType sets the view type. Slack defaults to html.
func SetAgentConversationViewOptionType(viewType string) SetAgentConversationViewOption {
	return func(params *setAgentConversationViewParams) {
		params.values.Set("type", viewType)
	}
}

// SetAgentConversationViewOptionViewKey sets the stable identity used to create
// or update the same view. It is required for html, block_kit, and canvas views,
// and ignored for diff views, which are unique within a channel.
func SetAgentConversationViewOptionViewKey(viewKey string) SetAgentConversationViewOption {
	return func(params *setAgentConversationViewParams) {
		params.values.Set("view_key", viewKey)
	}
}

// SetAgentConversationViewOptionContent sets the HTML document or unified diff.
// An empty string is sent explicitly. Slack limits content to 1,000,000 bytes.
func SetAgentConversationViewOptionContent(content string) SetAgentConversationViewOption {
	return func(params *setAgentConversationViewParams) {
		params.values.Set("content", content)
	}
}

// SetAgentConversationViewOptionBlocks sets the blocks for a block_kit view.
// Calling it without blocks sends an empty array.
func SetAgentConversationViewOptionBlocks(blocks ...Block) SetAgentConversationViewOption {
	return func(params *setAgentConversationViewParams) {
		params.blocks = append([]Block{}, blocks...)
	}
}

// SetAgentConversationViewOptionCanvasID sets the canvas to attach as a view.
func SetAgentConversationViewOptionCanvasID(canvasID string) SetAgentConversationViewOption {
	return func(params *setAgentConversationViewParams) {
		params.values.Set("canvas_id", canvasID)
	}
}

// SetAgentConversationViewOptionAccessLevel sets the channel's access to a canvas
// view: read, write, or comment. Slack defaults to write.
func SetAgentConversationViewOptionAccessLevel(accessLevel string) SetAgentConversationViewOption {
	return func(params *setAgentConversationViewParams) {
		params.values.Set("access_level", accessLevel)
	}
}

// SetAgentConversationViewOptionAgentContentHash records an opaque hash of the
// canvas-derived markdown last written by the agent, for detecting human edits.
func SetAgentConversationViewOptionAgentContentHash(hash string) SetAgentConversationViewOption {
	return func(params *setAgentConversationViewParams) {
		params.values.Set("agent_content_hash", hash)
	}
}

// SetAgentConversationViewOptionPRURL sets the URL for a pull_request view.
func SetAgentConversationViewOptionPRURL(prURL string) SetAgentConversationViewOption {
	return func(params *setAgentConversationViewParams) {
		params.values.Set("pr_url", prURL)
	}
}

// SetAgentConversationViewOptionBaseBranch sets the base branch label of a diff view.
func SetAgentConversationViewOptionBaseBranch(branch string) SetAgentConversationViewOption {
	return func(params *setAgentConversationViewParams) {
		params.values.Set("base_branch", branch)
	}
}

// SetAgentConversationViewOptionHeadBranch sets the head branch label of a diff view.
func SetAgentConversationViewOptionHeadBranch(branch string) SetAgentConversationViewOption {
	return func(params *setAgentConversationViewParams) {
		params.values.Set("head_branch", branch)
	}
}

// SetAgentConversationViewOptionName sets the view tab's display name.
func SetAgentConversationViewOptionName(name string) SetAgentConversationViewOption {
	return func(params *setAgentConversationViewParams) {
		params.values.Set("name", name)
	}
}

// SetAgentConversationViewOptionLabel sets the legacy view tab label.
//
// Deprecated: Use SetAgentConversationViewOptionName instead.
func SetAgentConversationViewOptionLabel(label string) SetAgentConversationViewOption {
	return func(params *setAgentConversationViewParams) {
		params.values.Set("label", label)
	}
}

// SetAgentConversationViewOptionCSP sets the HTML view's content security policy
// domain declarations.
func SetAgentConversationViewOptionCSP(csp AgentViewCSP) SetAgentConversationViewOption {
	return func(params *setAgentConversationViewParams) {
		params.csp = &csp
	}
}

// AgentConversationSetViewResponse is the response from agents.conversations.setView.
// ViewID identifies the channel tab; FileID and CanvasID identify backing content.
type AgentConversationSetViewResponse struct {
	SlackResponse
	ChannelID      string `json:"channel_id"`
	ViewID         string `json:"view_id"`
	FileID         string `json:"file_id"`
	CanvasID       string `json:"canvas_id"`
	ContentVersion int    `json:"content_version"`
	Type           string `json:"type"`
}

// SetAgentConversationView creates or updates a view in a code channel.
// HTML and diff views require content; block_kit views require blocks; canvas
// views require a canvas ID; pull_request views require a pull request URL.
//
// Slack API docs: https://docs.slack.dev/reference/methods/agents.conversations.setView
func (api *Client) SetAgentConversationView(channelID string, options ...SetAgentConversationViewOption) (*AgentConversationSetViewResponse, error) {
	return api.SetAgentConversationViewContext(context.Background(), channelID, options...)
}

// SetAgentConversationViewContext creates or updates a code channel view with a
// custom context. Slack validates the view type and its required arguments.
//
// Slack API docs: https://docs.slack.dev/reference/methods/agents.conversations.setView
func (api *Client) SetAgentConversationViewContext(ctx context.Context, channelID string, options ...SetAgentConversationViewOption) (*AgentConversationSetViewResponse, error) {
	params := setAgentConversationViewParams{
		values: url.Values{
			"token":      {api.token},
			"channel_id": {channelID},
		},
	}
	for _, option := range options {
		if option != nil {
			option(&params)
		}
	}
	if params.blocks != nil {
		blocksJSON, err := json.Marshal(params.blocks)
		if err != nil {
			return nil, err
		}
		params.values.Set("blocks", string(blocksJSON))
	}
	if params.csp != nil {
		cspJSON, err := json.Marshal(params.csp)
		if err != nil {
			return nil, err
		}
		params.values.Set("csp", string(cspJSON))
	}

	response := &AgentConversationSetViewResponse{}
	if err := api.postMethod(ctx, "agents.conversations.setView", params.values, response); err != nil {
		return nil, err
	}
	return response, response.Err()
}
