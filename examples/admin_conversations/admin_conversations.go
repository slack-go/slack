// This example demonstrates the admin.conversations.* API methods.
// These methods require an Enterprise Grid organization and an app installed
// at the org level with admin.* scopes.
//
// Usage:
//
//	export SLACK_USER_TOKEN="xoxp-..."
//	export SLACK_TEAM_ID="T..."  # Optional: workspace ID for scoping operations
//	go run admin_conversations.go
//
// The example provides a menu to test different operations. Read-only operations
// are safe to run. Destructive operations are clearly marked.
package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/slack-go/slack"
)

var (
	api    *slack.Client
	teamID string
)

func main() {
	token := os.Getenv("SLACK_USER_TOKEN")
	if token == "" {
		fmt.Fprintln(os.Stderr, "SLACK_USER_TOKEN environment variable is required")
		fmt.Fprintln(os.Stderr, "This must be an org-level token with admin.conversations:* scopes")
		os.Exit(1)
	}

	teamID = os.Getenv("SLACK_TEAM_ID")

	api = slack.New(token)

	reader := bufio.NewReader(os.Stdin)

	for {
		printMenu()
		fmt.Print("\nChoice: ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		switch input {
		case "1":
			testSearch()
		case "2":
			testGetTeams(reader)
		case "3":
			testGetConversationPrefs(reader)
		case "4":
			testGetCustomRetention(reader)
		case "5":
			testLookup()
		case "6":
			testRestrictAccessListGroups(reader)
		case "7":
			testEKMListOriginalConnectedChannelInfo()
		case "8":
			testCreate(reader)
		case "9":
			testInvite(reader)
		case "10":
			testRename(reader)
		case "11":
			testSetConversationPrefs(reader)
		case "12":
			testSetCustomRetention(reader)
		case "13":
			testArchive(reader)
		case "14":
			testUnarchive(reader)
		case "15":
			testDelete(reader)
		case "q", "Q":
			fmt.Println("Goodbye!")
			return
		default:
			fmt.Println("Invalid choice")
		}
		fmt.Println()
	}
}

func printMenu() {
	fmt.Println("\n" + strings.Repeat("=", 60))
	fmt.Println("Admin Conversations API Demo")
	fmt.Println(strings.Repeat("=", 60))
	fmt.Println("\nREAD-ONLY (safe):")
	fmt.Println("  1. Search conversations")
	fmt.Println("  2. Get teams for a channel")
	fmt.Println("  3. Get conversation preferences")
	fmt.Println("  4. Get custom retention policy")
	fmt.Println("  5. Lookup channels by activity")
	fmt.Println("  6. List restrict access groups")
	fmt.Println("  7. List EKM original connected channel info")
	fmt.Println("\nCREATE/MODIFY:")
	fmt.Println("  8. Create a test channel")
	fmt.Println("  9. Invite users to a channel")
	fmt.Println(" 10. Rename a channel")
	fmt.Println(" 11. Set conversation preferences")
	fmt.Println(" 12. Set custom retention policy")
	fmt.Println("\nARCHIVE/DELETE (use with caution):")
	fmt.Println(" 13. Archive a channel")
	fmt.Println(" 14. Unarchive a channel")
	fmt.Println(" 15. Delete a channel [DESTRUCTIVE]")
	fmt.Println("\n  q. Quit")
}

func testSearch() {
	fmt.Println("\n--- Searching conversations ---")

	options := []slack.AdminConversationsSearchOption{
		slack.AdminConversationsSearchOptionLimit(10),
		slack.AdminConversationsSearchOptionSort("member_count"),
	}

	if teamID != "" {
		options = append(options, slack.AdminConversationsSearchOptionTeamIDs([]string{teamID}))
	}

	response, err := api.AdminConversationsSearch(context.Background(), options...)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return
	}

	fmt.Printf("Found %d conversations (showing up to 10)\n", response.TotalCount)
	fmt.Printf("Next cursor: %q\n\n", response.NextCursor)

	for _, conv := range response.Conversations {
		fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
		fmt.Printf("Channel: #%s (%s)\n", conv.Name, conv.ID)
		fmt.Printf("  Purpose: %s\n", conv.Purpose)
		fmt.Printf("  Created: %d | Creator: %s\n", conv.Created, conv.CreatorID)
		fmt.Printf("  Members: %d | External users: %d | Channel managers: %d\n",
			conv.MemberCount, conv.ExternalUserCount, conv.ChannelManagerCount)
		fmt.Printf("  Last activity: %d\n", conv.LastActivityTimestamp)

		// Visibility flags
		fmt.Printf("  Flags: ")
		flags := []string{}
		if conv.IsPrivate {
			flags = append(flags, "private")
		} else {
			flags = append(flags, "public")
		}
		if conv.IsArchived {
			flags = append(flags, "archived")
		}
		if conv.IsGeneral {
			flags = append(flags, "general")
		}
		if conv.IsFrozen {
			flags = append(flags, "frozen")
		}
		if conv.IsOrgShared {
			flags = append(flags, "org-shared")
		}
		if conv.IsExtShared {
			flags = append(flags, "ext-shared")
		}
		if conv.IsGlobalShared {
			flags = append(flags, "global-shared")
		}
		if conv.IsPendingExtShared {
			flags = append(flags, "pending-ext-shared")
		}
		if conv.IsDisconnectInProgress {
			flags = append(flags, "disconnect-in-progress")
		}
		if conv.IsOrgDefault {
			flags = append(flags, "org-default")
		}
		if conv.IsOrgMandatory {
			flags = append(flags, "org-mandatory")
		}
		fmt.Printf("%s\n", strings.Join(flags, ", "))

		// Team IDs
		if len(conv.ConnectedTeamIDs) > 0 {
			fmt.Printf("  Connected teams: %v\n", conv.ConnectedTeamIDs)
		}
		if len(conv.ConnectedLimitedTeamIDs) > 0 {
			fmt.Printf("  Connected limited teams: %v\n", conv.ConnectedLimitedTeamIDs)
		}
		if len(conv.PendingConnectedTeamIDs) > 0 {
			fmt.Printf("  Pending connected teams: %v\n", conv.PendingConnectedTeamIDs)
		}
		if len(conv.InternalTeamIDs) > 0 {
			fmt.Printf("  Internal teams: %v\n", conv.InternalTeamIDs)
		}
		if conv.InternalTeamIDsCount > 0 {
			fmt.Printf("  Internal teams count: %d (sample: %s)\n",
				conv.InternalTeamIDsCount, conv.InternalTeamIDsSampleTeam)
		}
		if conv.ContextTeamID != "" {
			fmt.Printf("  Context team: %s\n", conv.ContextTeamID)
		}
		if conv.ConversationHostID != "" {
			fmt.Printf("  Conversation host: %s\n", conv.ConversationHostID)
		}

		// Email addresses
		if len(conv.ChannelEmailAddresses) > 0 {
			fmt.Printf("  Email addresses:\n")
			for _, email := range conv.ChannelEmailAddresses {
				fmt.Printf("    - %s (team: %s, creator: %s)\n",
					email.Address, email.TeamID, email.CreatorID)
			}
		}

		// Canvas/Lists
		if conv.Canvas != nil {
			fmt.Printf("  Canvas: total_count=%d\n", conv.Canvas.TotalCount)
			for _, od := range conv.Canvas.OwnershipDetails {
				fmt.Printf("    - team %s: %d\n", od.TeamID, od.Count)
			}
		}
		if conv.Lists != nil {
			fmt.Printf("  Lists: total_count=%d\n", conv.Lists.TotalCount)
			for _, od := range conv.Lists.OwnershipDetails {
				fmt.Printf("    - team %s: %d\n", od.TeamID, od.Count)
			}
		}

		// Properties
		if conv.Properties != nil {
			fmt.Printf("  Properties: present\n")
		}

		fmt.Println()
	}
}

func testGetTeams(reader *bufio.Reader) {
	fmt.Println("\n--- Get teams for channel ---")
	channelID := prompt(reader, "Channel ID (e.g., C1234567890): ")

	teamIDs, cursor, err := api.AdminConversationsGetTeams(context.Background(), slack.AdminConversationsGetTeamsParams{
		ChannelID: channelID,
		Limit:     100,
	})
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return
	}

	fmt.Printf("Teams connected to %s:\n", channelID)
	for _, tid := range teamIDs {
		fmt.Printf("  - %s\n", tid)
	}
	if cursor != "" {
		fmt.Printf("(more results available, cursor: %s)\n", cursor)
	}
}

func testGetConversationPrefs(reader *bufio.Reader) {
	fmt.Println("\n--- Get conversation preferences ---")
	channelID := prompt(reader, "Channel ID: ")

	prefs, err := api.AdminConversationsGetConversationPrefs(context.Background(), channelID)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return
	}

	fmt.Printf("Preferences for %s:\n", channelID)
	if prefs.WhoCanPost != nil {
		fmt.Printf("  Who can post: types=%v, users=%v\n", prefs.WhoCanPost.Type, prefs.WhoCanPost.User)
	}
	if prefs.CanThread != nil {
		fmt.Printf("  Can thread: types=%v\n", prefs.CanThread.Type)
	}
	if prefs.CanHuddle != nil {
		fmt.Printf("  Can huddle: types=%v\n", prefs.CanHuddle.Type)
	}
}

func testGetCustomRetention(reader *bufio.Reader) {
	fmt.Println("\n--- Get custom retention policy ---")
	channelID := prompt(reader, "Channel ID: ")

	resp, err := api.AdminConversationsGetCustomRetention(context.Background(), channelID)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return
	}

	if !resp.IsPolicyEnabled {
		fmt.Printf("Channel %s does not have a custom retention policy enabled\n", channelID)
	} else {
		fmt.Printf("Channel %s has custom retention: %d days\n", channelID, resp.DurationDays)
	}
}

func testLookup() {
	fmt.Println("\n--- Lookup channels by activity ---")

	if teamID == "" {
		fmt.Println("Error: SLACK_TEAM_ID environment variable is required for lookup")
		return
	}

	// Find channels with no activity in the last 90 days
	cutoff := time.Now().AddDate(0, 0, -90).Unix()

	channels, cursor, err := api.AdminConversationsLookup(context.Background(),
		[]string{teamID}, cutoff,
		slack.AdminConversationsLookupOptionLimit(10),
	)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return
	}

	fmt.Printf("Channels with no activity in last 90 days:\n")
	if len(channels) == 0 {
		fmt.Println("  (none found)")
	}
	for _, ch := range channels {
		fmt.Printf("  - %s\n", ch)
	}
	if cursor != "" {
		fmt.Printf("(more results available)\n")
	}
}

func testRestrictAccessListGroups(reader *bufio.Reader) {
	fmt.Println("\n--- List restrict access groups ---")
	channelID := prompt(reader, "Channel ID: ")

	var options []slack.AdminConversationsRestrictAccessListGroupsOption
	if teamID != "" {
		options = append(options, slack.AdminConversationsRestrictAccessListGroupsOptionTeamID(teamID))
	}

	groupIDs, err := api.AdminConversationsRestrictAccessListGroups(context.Background(), channelID, options...)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return
	}

	fmt.Printf("IDP groups with access to %s:\n", channelID)
	if len(groupIDs) == 0 {
		fmt.Println("  (no IDP group restrictions)")
	}
	for _, gid := range groupIDs {
		fmt.Printf("  - %s\n", gid)
	}
}

func testEKMListOriginalConnectedChannelInfo() {
	fmt.Println("\n--- List EKM original connected channel info ---")

	var options []slack.AdminConversationsEKMListOriginalConnectedChannelInfoOption
	options = append(options, slack.AdminConversationsEKMListOriginalConnectedChannelInfoOptionLimit(10))

	if teamID != "" {
		options = append(options, slack.AdminConversationsEKMListOriginalConnectedChannelInfoOptionTeamIDs([]string{teamID}))
	}

	response, err := api.AdminConversationsEKMListOriginalConnectedChannelInfo(context.Background(), options...)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return
	}

	fmt.Printf("Found %d channels\n\n", len(response.Channels))

	for _, ch := range response.Channels {
		fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
		fmt.Printf("Channel: %s\n", ch.ID)
		fmt.Printf("  Original connected host ID: %s\n", ch.OriginalConnectedHostID)
		fmt.Printf("  Original connected channel ID: %s\n", ch.OriginalConnectedChannelID)
		fmt.Printf("  Internal team IDs: %v\n", ch.InternalTeamIDs)
		fmt.Println()
	}

	if len(response.Channels) == 0 {
		fmt.Println("  (no EKM connected channel info found)")
	}
}

func testCreate(reader *bufio.Reader) {
	fmt.Println("\n--- Create test channel ---")
	name := prompt(reader, "Channel name (e.g., test-admin-api): ")
	isPrivate := strings.ToLower(prompt(reader, "Private? (y/n): ")) == "y"
	description := prompt(reader, "Description (optional): ")

	var options []slack.AdminConversationsCreateOption
	if description != "" {
		options = append(options, slack.AdminConversationsCreateOptionDescription(description))
	}
	if teamID != "" {
		options = append(options, slack.AdminConversationsCreateOptionTeamID(teamID))
	}

	channelID, err := api.AdminConversationsCreate(context.Background(), name, isPrivate, options...)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return
	}

	fmt.Printf("Created channel: %s\n", channelID)
	fmt.Println("Save this ID to test other operations!")
}

func testInvite(reader *bufio.Reader) {
	fmt.Println("\n--- Invite users to channel ---")
	channelID := prompt(reader, "Channel ID: ")
	userIDsStr := prompt(reader, "User IDs (comma-separated, e.g., U123,U456): ")

	userIDs := strings.Split(userIDsStr, ",")
	for i := range userIDs {
		userIDs[i] = strings.TrimSpace(userIDs[i])
	}

	err := api.AdminConversationsInvite(context.Background(), slack.AdminConversationsInviteParams{
		ChannelID: channelID,
		UserIDs:   userIDs,
	})
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return
	}

	fmt.Printf("Invited %d user(s) to %s\n", len(userIDs), channelID)
}

func testRename(reader *bufio.Reader) {
	fmt.Println("\n--- Rename channel ---")
	channelID := prompt(reader, "Channel ID: ")
	newName := prompt(reader, "New name: ")

	err := api.AdminConversationsRename(context.Background(), channelID, newName)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return
	}

	fmt.Printf("Renamed %s to #%s\n", channelID, newName)
}

func testSetConversationPrefs(reader *bufio.Reader) {
	fmt.Println("\n--- Set conversation preferences ---")
	channelID := prompt(reader, "Channel ID: ")

	fmt.Println("Who can post?")
	fmt.Println("  1. Everyone")
	fmt.Println("  2. Admins only")
	fmt.Println("  3. Org admins only")
	choice := prompt(reader, "Choice: ")

	var whoCanPost []string
	switch choice {
	case "1":
		whoCanPost = []string{"owner", "admin", "org_admin", "member", "ra_member", "guest"}
	case "2":
		whoCanPost = []string{"owner", "admin", "org_admin"}
	case "3":
		whoCanPost = []string{"owner", "org_admin"}
	default:
		fmt.Println("Invalid choice")
		return
	}

	err := api.AdminConversationsSetConversationPrefs(context.Background(), slack.AdminConversationsSetConversationPrefsParams{
		ChannelID: channelID,
		Prefs: slack.AdminConversationPrefs{
			WhoCanPost: &slack.AdminConversationPref{Type: whoCanPost},
		},
	})
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return
	}

	fmt.Printf("Updated posting preferences for %s\n", channelID)
}

func testSetCustomRetention(reader *bufio.Reader) {
	fmt.Println("\n--- Set custom retention policy ---")
	channelID := prompt(reader, "Channel ID: ")
	daysStr := prompt(reader, "Retention days (e.g., 90): ")

	var days int
	_, err := fmt.Sscanf(daysStr, "%d", &days)
	if err != nil || days <= 0 {
		fmt.Println("Invalid number of days")
		return
	}

	err = api.AdminConversationsSetCustomRetention(context.Background(), channelID, days)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return
	}

	fmt.Printf("Set %d-day retention policy for %s\n", days, channelID)
}

func testArchive(reader *bufio.Reader) {
	fmt.Println("\n--- Archive channel ---")
	fmt.Println("WARNING: This will archive the channel!")
	channelID := prompt(reader, "Channel ID: ")

	confirm := prompt(reader, "Type 'archive' to confirm: ")
	if confirm != "archive" {
		fmt.Println("Cancelled")
		return
	}

	err := api.AdminConversationsArchive(context.Background(), channelID)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return
	}

	fmt.Printf("Archived %s\n", channelID)
}

func testUnarchive(reader *bufio.Reader) {
	fmt.Println("\n--- Unarchive channel ---")
	channelID := prompt(reader, "Channel ID: ")

	err := api.AdminConversationsUnarchive(context.Background(), channelID)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return
	}

	fmt.Printf("Unarchived %s\n", channelID)
}

func testDelete(reader *bufio.Reader) {
	fmt.Println("\n--- Delete channel ---")
	fmt.Println("!!! WARNING: THIS IS PERMANENT AND CANNOT BE UNDONE !!!")
	channelID := prompt(reader, "Channel ID: ")

	confirm := prompt(reader, "Type 'DELETE' (all caps) to confirm: ")
	if confirm != "DELETE" {
		fmt.Println("Cancelled")
		return
	}

	err := api.AdminConversationsDelete(context.Background(), channelID)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return
	}

	fmt.Printf("Deleted %s\n", channelID)
}

func prompt(reader *bufio.Reader, message string) string {
	fmt.Print(message)
	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(input)
}
