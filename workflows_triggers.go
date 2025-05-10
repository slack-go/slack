package slack

import (
	"context"
	"errors"
	"fmt"
)

type (
	WorkflowsTriggersPermissionsAddInput struct {
		TriggerId  string   `json:"trigger_id"`
		ChannelIds []string `json:"channel_ids,omitempty"`
		OrgIds     []string `json:"org_ids,omitempty"`
		TeamIds    []string `json:"team_ids,omitempty"`
		UserIds    []string `json:"user_ids,omitempty"`
	}

	WorkflowsTriggersPermissionsAddOutput struct {
		PermissionType string   `json:"permission_type"`
		ChannelIds     []string `json:"channel_ids,omitempty"`
		OrgIds         []string `json:"org_ids,omitempty"`
		TeamIds        []string `json:"team_ids,omitempty"`
		UserIds        []string `json:"user_ids,omitempty"`
	}

	WorkflowsTriggersPermissionsListInput struct {
		TriggerId string `json:"trigger_id"`
	}

	WorkflowsTriggersPermissionsListOutput struct {
		PermissionType string   `json:"permission_type"`
		ChannelIds     []string `json:"channel_ids,omitempty"`
		OrgIds         []string `json:"org_ids,omitempty"`
		TeamIds        []string `json:"team_ids,omitempty"`
		UserIds        []string `json:"user_ids,omitempty"`
	}

	WorkflowsTriggersPermissionsRemoveInput struct {
		TriggerId  string   `json:"trigger_id"`
		ChannelIds []string `json:"channel_ids,omitempty"`
		OrgIds     []string `json:"org_ids,omitempty"`
		TeamIds    []string `json:"team_ids,omitempty"`
		UserIds    []string `json:"user_ids,omitempty"`
	}

	WorkflowsTriggersPermissionsRemoveOutput struct {
		PermissionType string   `json:"permission_type"`
		ChannelIds     []string `json:"channel_ids,omitempty"`
		OrgIds         []string `json:"org_ids,omitempty"`
		TeamIds        []string `json:"team_ids,omitempty"`
		UserIds        []string `json:"user_ids,omitempty"`
	}

	WorkflowsTriggersPermissionsSetInput struct {
		PermissionType string   `json:"permission_type"`
		TriggerId      string   `json:"trigger_id"`
		ChannelIds     []string `json:"channel_ids,omitempty"`
		OrgIds         []string `json:"org_ids,omitempty"`
		TeamIds        []string `json:"team_ids,omitempty"`
		UserIds        []string `json:"user_ids,omitempty"`
	}

	WorkflowsTriggersPermissionsSetOutput struct {
		PermissionType string   `json:"permission_type"`
		ChannelIds     []string `json:"channel_ids,omitempty"`
		OrgIds         []string `json:"org_ids,omitempty"`
		TeamIds        []string `json:"team_ids,omitempty"`
		UserIds        []string `json:"user_ids,omitempty"`
	}
)

// WorkflowsTriggersPermissionsAdd allows users to run a trigger that has its permission type set to named_entities.
//
// Slack API Docs:https://api.dev.slack.com/methods/workflows.triggers.permissions.add
func (api *Client) WorkflowsTriggersPermissionsAdd(ctx context.Context, input *WorkflowsTriggersPermissionsAddInput) (*WorkflowsTriggersPermissionsAddOutput, error) {
	response := struct {
		*ResponsePointer
		*WorkflowsTriggersPermissionsAddOutput
	}{}

	err := api.postJSON(ctx, "workflows.triggers.permissions.add", input, &response)
	if err != nil {
		return nil, err
	}

	if response.Error != nil {
		return nil, errors.New(fmt.Sprintf("error: %s", *response.Error))
	}

	return response.WorkflowsTriggersPermissionsAddOutput, nil
}

// WorkflowsTriggersPermissionsList returns the permission type of a trigger and if applicable, includes the entities that have been granted access.
//
// Slack API Docs:https://api.dev.slack.com/methods/workflows.triggers.permissions.list
func (api *Client) WorkflowsTriggersPermissionsList(ctx context.Context, input *WorkflowsTriggersPermissionsListInput) (*WorkflowsTriggersPermissionsListOutput, error) {
	response := struct {
		*ResponsePointer
		*WorkflowsTriggersPermissionsListOutput
	}{}

	err := api.postJSON(ctx, "workflows.triggers.permissions.list", input, &response)
	if err != nil {
		return nil, err
	}

	if response.Error != nil {
		return nil, errors.New(fmt.Sprintf("error: %s", *response.Error))
	}

	return response.WorkflowsTriggersPermissionsListOutput, nil
}

// WorkflowsTriggersPermissionsRemove revoke an entity's access to a trigger that has its permission type set to named_entities.
//
// Slack API Docs:https://api.dev.slack.com/methods/workflows.triggers.permissions.remove
func (api *Client) WorkflowsTriggersPermissionsRemove(ctx context.Context, input *WorkflowsTriggersPermissionsRemoveInput) (*WorkflowsTriggersPermissionsRemoveOutput, error) {
	response := struct {
		*ResponsePointer
		*WorkflowsTriggersPermissionsRemoveOutput
	}{}

	err := api.postJSON(ctx, "workflows.triggers.permissions.remove", input, &response)
	if err != nil {
		return nil, err
	}

	if response.Error != nil {
		return nil, errors.New(fmt.Sprintf("error: %s", *response.Error))
	}

	return response.WorkflowsTriggersPermissionsRemoveOutput, nil
}

// WorkflowsTriggersPermissionsSet sets the permission type for who can run a trigger.
//
// Slack API Docs:https://api.dev.slack.com/methods/workflows.triggers.permissions.set
func (api *Client) WorkflowsTriggersPermissionsSet(ctx context.Context, input *WorkflowsTriggersPermissionsSetInput) (*WorkflowsTriggersPermissionsSetOutput, error) {
	response := struct {
		*ResponsePointer
		*WorkflowsTriggersPermissionsSetOutput
	}{}

	err := api.postJSON(ctx, "workflows.triggers.permissions.set", input, &response)
	if err != nil {
		return nil, err
	}

	if response.Error != nil {
		return nil, errors.New(fmt.Sprintf("error: %s", *response.Error))
	}

	return response.WorkflowsTriggersPermissionsSetOutput, nil
}
