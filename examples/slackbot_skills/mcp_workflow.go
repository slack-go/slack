package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/slack-go/slack"
)

// runMCPWorkflow changes only newly created skills. It never invokes an MCP tool
// or changes the app's permissions, publication, or other users' libraries.
func runMCPWorkflow(ctx context.Context, api *slack.Client, appID string, supplied []json.RawMessage) (result error) {
	if len(supplied) == 0 {
		return errors.New("mcp-workflow needs non-empty declarations via -mcp-servers or -mcp-servers-file")
	}
	payload, err := json.Marshal(supplied)
	if err != nil {
		return err
	}
	var expected []slack.SlackbotSkillMCPServer
	if err := json.Unmarshal(payload, &expected); err != nil {
		return err
	}
	ids := []string{}
	defer func() {
		for _, id := range ids {
			cleanupCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			err := api.DeleteSlackbotSkillContext(cleanupCtx, id)
			cancel()
			if err != nil {
				result = errors.Join(result, fmt.Errorf("cleanup temporary skill %s: %w", id, err))
			} else {
				fmt.Printf("Deleted temporary skill %s\n", id)
			}
		}
		if result == nil {
			fmt.Println("PASS: non-empty MCP declarations verified on create and update; clearing verified; temporary skills deleted")
		}
	}()
	create := func(options ...slack.CreateSlackbotSkillOption) (*slack.SlackbotSkillResponse, error) {
		options = append([]slack.CreateSlackbotSkillOption{
			slack.CreateSlackbotSkillOptionInstructions("# MCP declaration verification\n\nThis is a temporary API fixture. Do not invoke external tools."),
		}, options...)
		response, err := api.CreateSlackbotSkillContext(ctx, fmt.Sprintf("MCP declaration verification %d", time.Now().UnixNano()), options...)
		if response != nil && response.Skill.ID != "" {
			ids = append(ids, response.Skill.ID)
		}
		if err := output(response, err); err != nil {
			return response, err
		}
		if response.Skill.ID == "" {
			return response, errors.New("create returned no skill ID")
		}
		return response, nil
	}
	baseline, err := create()
	if err != nil {
		return fmt.Errorf("create temporary skill: %w", err)
	}
	id := baseline.Skill.ID
	if len(baseline.Skill.MCPServers) != 0 {
		return errors.New("new skill unexpectedly contains MCP declarations")
	}
	fmt.Printf("Testing mcp_servers=%s\n", payload)
	response, err := api.UpdateSlackbotSkillContext(ctx, id, slack.UpdateSlackbotSkillOptionMCPServers(supplied...))
	if err := output(response, err); err != nil {
		return fmt.Errorf("update declarations: %w", err)
	}
	if err := verifyMCPDeclarations(ctx, api, id, appID, expected); err != nil {
		return fmt.Errorf("update declarations: %w", err)
	}
	_, err = api.UpdateSlackbotSkillContext(ctx, id, slack.UpdateSlackbotSkillOptionMCPServers())
	if err != nil {
		return fmt.Errorf("clear declarations: %w", err)
	}
	cleared, err := api.GetSlackbotSkillInfoContext(ctx, id)
	if err != nil {
		return err
	}
	if len(cleared.Skill.MCPServers) != 0 {
		return errors.New("sending [] did not clear MCP declarations")
	}
	withMCP, err := create(slack.CreateSlackbotSkillOptionMCPServers(supplied...))
	if err != nil {
		return fmt.Errorf("create with declarations: %w", err)
	}
	if err := verifyMCPDeclarations(ctx, api, withMCP.Skill.ID, appID, expected); err != nil {
		return fmt.Errorf("create did not round-trip declarations: %w", err)
	}
	return nil
}

func verifyMCPDeclarations(ctx context.Context, api *slack.Client, skillID, appID string, expected []slack.SlackbotSkillMCPServer) error {
	response, err := api.GetSlackbotSkillInfoContext(ctx, skillID)
	if err := output(response, err); err != nil {
		return err
	}
	actual := response.Skill.MCPServers
	if len(actual) != len(expected) {
		return fmt.Errorf("expected %d MCP declarations, info returned %d", len(expected), len(actual))
	}
	// Match identifiers without depending on Slack's response order or enriched fields.
	remaining := append([]slack.SlackbotSkillMCPServer(nil), actual...)
	for _, want := range expected {
		found := false
		for i, got := range remaining {
			if got.AppID == appID && got.AppID == want.AppID && got.ServerID == want.ServerID && got.ServiceType == want.ServiceType {
				remaining = append(remaining[:i], remaining[i+1:]...)
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("info did not return declaration for app %s, server %s, service type %s", want.AppID, want.ServerID, want.ServiceType)
		}
	}
	return nil
}
