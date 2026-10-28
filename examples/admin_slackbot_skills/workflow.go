package main

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"time"

	"github.com/slack-go/slack"
)

// runWorkflow uses a fresh skill and only the groups explicitly selected by the
// user. The same org-admin token needs skills:write for setup and cleanup.
func runWorkflow(ctx context.Context, api *slack.Client, groups []string) (result error) {
	created, err := api.CreateSlackbotSkillContext(ctx,
		fmt.Sprintf("Admin Skills API verification %d", time.Now().UnixNano()),
		slack.CreateSlackbotSkillOptionInstructions("# API verification fixture\n\nWhen invoked, reply that this is a temporary API verification skill."),
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
		_, disableErr := api.AdminSlackbotSkillsSetAutoEnabledContext(cleanupCtx, id, "no_one")
		_, unpublishErr := api.AdminSlackbotSkillsUnpublishContext(cleanupCtx, []string{id})
		deleteErr := api.DeleteSlackbotSkillContext(cleanupCtx, id)
		if err := errors.Join(disableErr, unpublishErr, deleteErr); err != nil {
			result = errors.Join(result, fmt.Errorf("cleanup skill %s: %w", id, err))
		}
	}()
	if err := output(created, err); err != nil {
		return fmt.Errorf("create fixture: %w", err)
	}
	if id == "" {
		return errors.New("create returned an empty skill ID")
	}
	if err := waitForInventory(ctx, api, id, false, "no_one", nil); err != nil {
		return err
	}
	published, err := api.AdminSlackbotSkillsPublishContext(ctx, []string{id})
	if err := output(published, err); err != nil {
		return fmt.Errorf("publish: %w", err)
	}
	if err := waitForInventory(ctx, api, id, true, "no_one", nil); err != nil {
		return err
	}
	deployment, err := api.AdminSlackbotSkillsSetAutoEnabledContext(ctx, id, "named_entities",
		slack.AdminSlackbotSkillsSetAutoEnabledOptionUsergroupIDs(groups...))
	if err := output(deployment, err); err != nil {
		return fmt.Errorf("setAutoEnabled: %w", err)
	}
	wantGroups := slices.Clone(groups)
	slices.Sort(wantGroups)
	wantGroups = slices.Compact(wantGroups)
	if deployment.SkillID != id || deployment.PermissionType != "named_entities" || !slices.Equal(deployment.UsergroupIDs, wantGroups) {
		return errors.New("deployment response did not match requested audience")
	}
	if err := waitForInventory(ctx, api, id, true, "named_entities", wantGroups); err != nil {
		return err
	}
	unpublished, err := api.AdminSlackbotSkillsUnpublishContext(ctx, []string{id})
	if err := output(unpublished, err); err != nil {
		return fmt.Errorf("unpublish: %w", err)
	}
	// Unpublishing must preserve deployment; inventory proves that distinction.
	if err := waitForInventory(ctx, api, id, false, "named_entities", wantGroups); err != nil {
		return err
	}
	disabled, err := api.AdminSlackbotSkillsSetAutoEnabledContext(ctx, id, "no_one")
	if err := output(disabled, err); err != nil {
		return fmt.Errorf("disable: %w", err)
	}
	if err := waitForInventory(ctx, api, id, false, "no_one", nil); err != nil {
		return err
	}
	if err := api.DeleteSlackbotSkillContext(ctx, id); err != nil {
		return fmt.Errorf("delete fixture: %w", err)
	}
	deleted = true
	fmt.Println("PASS: all four admin.slackbot.skills methods verified; temporary skill deleted")
	return nil
}

func waitForInventory(ctx context.Context, api *slack.Client, id string, published bool, permission string, groups []string) error {
	ctx, cancel := context.WithTimeout(ctx, time.Minute)
	defer cancel()
	for {
		cursor := ""
		for {
			page, err := api.AdminSlackbotSkillsListContext(ctx, slack.AdminSlackbotSkillsListOptionCursor(cursor))
			if err != nil {
				return fmt.Errorf("admin inventory: %w", err)
			}
			for _, skill := range page.Skills {
				// Inventory does not promise the group's ordering, unlike the
				// setAutoEnabled response. Compare the effective audience as a set.
				actualGroups := slices.Clone(skill.AutoEnabled.UsergroupIDs)
				slices.Sort(actualGroups)
				actualGroups = slices.Compact(actualGroups)
				if skill.ID == id && skill.IsOfficial == published && skill.AutoEnabled.PermissionType == permission && slices.Equal(actualGroups, groups) {
					return nil
				}
			}
			cursor = page.ResponseMetadata.Cursor
			if cursor == "" {
				break
			}
		}
		select {
		case <-ctx.Done():
			return fmt.Errorf("inventory did not show skill %s published=%t audience=%s: %w", id, published, permission, ctx.Err())
		case <-time.After(3 * time.Second):
		}
	}
}
