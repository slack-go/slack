package main

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/slack-go/slack"
)

// runWorkflow exercises all seven user endpoints on a newly created skill only.
// Cleanup uses a separate deadline so Ctrl-C still attempts to remove the fixture.
func runWorkflow(ctx context.Context, api *slack.Client) (result error) {
	const initial = "Slackbot Skills API verification fixture"
	const updated = "Updated Slackbot Skills API verification fixture"
	created, err := api.CreateSlackbotSkillContext(ctx,
		fmt.Sprintf("Skills API verification %d", time.Now().UnixNano()),
		slack.CreateSlackbotSkillOptionInstructions("# Verification\n\n"+initial),
		slack.CreateSlackbotSkillOptionCategory("engineering"))
	id := ""
	if created != nil {
		id = created.Skill.ID
	}
	deleted := false
	defer func() {
		if deleted || id == "" {
			return
		}
		cleanupCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		if err := api.DeleteSlackbotSkillContext(cleanupCtx, id); err != nil {
			result = errors.Join(result, fmt.Errorf("cleanup skill %s: %w", id, err))
		}
	}()
	if err := output(created, err); err != nil {
		return fmt.Errorf("create: %w", err)
	}
	if id == "" {
		return errors.New("create returned an empty skill ID")
	}
	info, err := api.GetSlackbotSkillInfoContext(ctx, id)
	if err != nil {
		return fmt.Errorf("info: %w", err)
	}
	if info.Skill.Instructions == nil || !strings.Contains(*info.Skill.Instructions, initial) {
		return errors.New("info did not return the created instructions")
	}
	changed, err := api.UpdateSlackbotSkillContext(ctx, id,
		slack.UpdateSlackbotSkillOptionInstructions("# Verification\n\n"+updated),
		slack.UpdateSlackbotSkillOptionDescription("Temporary verification skill"))
	if err := output(changed, err); err != nil {
		return fmt.Errorf("update: %w", err)
	}
	info, err = api.GetSlackbotSkillInfoContext(ctx, id)
	if err != nil {
		return fmt.Errorf("info after update: %w", err)
	}
	if info.Skill.Instructions == nil || !strings.Contains(*info.Skill.Instructions, updated) || info.Skill.Description != "Temporary verification skill" {
		return errors.New("updated instructions or description did not round-trip")
	}
	if err := waitForMembership(ctx, api, "owned", id, true); err != nil {
		return err
	}
	if err := api.RemoveSlackbotSkillFromLibraryContext(ctx, id); err != nil {
		return fmt.Errorf("library.remove: %w", err)
	}
	if err := waitForMembership(ctx, api, "library", id, false); err != nil {
		return err
	}
	if err := api.AddSlackbotSkillToLibraryContext(ctx, id); err != nil {
		return fmt.Errorf("library.add: %w", err)
	}
	if err := waitForMembership(ctx, api, "library", id, true); err != nil {
		return err
	}
	if err := api.DeleteSlackbotSkillContext(ctx, id); err != nil {
		return fmt.Errorf("delete: %w", err)
	}
	deleted = true
	_, err = api.GetSlackbotSkillInfoContext(ctx, id)
	var slackErr slack.SlackErrorResponse
	if !errors.As(err, &slackErr) || slackErr.Err != "skill_not_found" {
		return fmt.Errorf("expected skill_not_found after deletion; got %v", err)
	}
	fmt.Println("PASS: all seven slackbot.skills methods verified; temporary skill deleted")
	return nil
}

func waitForMembership(ctx context.Context, api *slack.Client, filter, id string, want bool) error {
	ctx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()
	for {
		found := false
		cursor := ""
		for {
			page, err := api.ListSlackbotSkillsContext(ctx,
				slack.ListSlackbotSkillsOptionFilter(filter), slack.ListSlackbotSkillsOptionCursor(cursor))
			if err != nil {
				return fmt.Errorf("list %s: %w", filter, err)
			}
			for _, skill := range page.Skills {
				if skill.ID == id {
					found = true
				}
			}
			cursor = page.ResponseMetadata.Cursor
			if found || cursor == "" {
				break
			}
		}
		if found == want {
			return nil
		}
		select {
		case <-ctx.Done():
			return fmt.Errorf("skill %s in %s: wanted membership %t: %w", id, filter, want, ctx.Err())
		case <-time.After(1500 * time.Millisecond):
		}
	}
}
