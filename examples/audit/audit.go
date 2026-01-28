// This example demonstrates the Audit Logs API.
// The Audit Logs API requires an Enterprise Grid organization with an app
// that has the auditlogs:read scope.
//
// Usage:
//
//	export SLACK_USER_TOKEN="xoxp-..."  # User token with auditlogs:read scope
//	go run audit.go
//
// Note: The Audit Logs API uses a different endpoint (api.slack.com) than
// the regular Slack API (slack.com/api). This is handled automatically by
// the library.
package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/slack-go/slack"
)

func main() {
	token := os.Getenv("SLACK_USER_TOKEN")
	if token == "" {
		fmt.Fprintln(os.Stderr, "SLACK_USER_TOKEN environment variable is required")
		fmt.Fprintln(os.Stderr, "This must be a user token with auditlogs:read scope")
		os.Exit(1)
	}

	api := slack.New(token)

	// Fetch the last 24 hours of audit logs
	now := time.Now()
	yesterday := now.Add(-24 * time.Hour)

	fmt.Println("Fetching audit logs from the last 24 hours...")
	fmt.Printf("  From: %s\n", yesterday.Format(time.RFC3339))
	fmt.Printf("  To:   %s\n\n", now.Format(time.RFC3339))

	params := slack.AuditLogParameters{
		Limit:  10,
		Oldest: int(yesterday.Unix()),
		Latest: int(now.Unix()),
	}

	entries, nextCursor, err := api.GetAuditLogsContext(context.Background(), params)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error fetching audit logs: %s\n", err)
		os.Exit(1)
	}

	fmt.Printf("Found %d audit log entries (limited to 10)\n", len(entries))
	if nextCursor != "" {
		fmt.Printf("More entries available (cursor: %s)\n", nextCursor)
	}
	fmt.Println()

	for i, entry := range entries {
		fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
		fmt.Printf("Entry %d: %s\n", i+1, entry.ID)
		fmt.Printf("  Action: %s\n", entry.Action)
		fmt.Printf("  Date:   %s\n", time.Unix(int64(entry.DateCreate), 0).Format(time.RFC3339))

		// Actor info
		fmt.Printf("  Actor:  %s (%s)\n", entry.Actor.User.Name, entry.Actor.User.Email)

		// Entity info
		fmt.Printf("  Entity Type: %s\n", entry.Entity.Type)
		switch entry.Entity.Type {
		case "user":
			fmt.Printf("  Entity: %s (%s)\n", entry.Entity.User.Name, entry.Entity.User.Email)
		case "channel":
			fmt.Printf("  Entity: #%s (%s)\n", entry.Entity.Channel.Name, entry.Entity.Channel.ID)
		case "file":
			fmt.Printf("  Entity: %s (%s)\n", entry.Entity.File.Name, entry.Entity.File.ID)
		case "app":
			fmt.Printf("  Entity: %s (%s)\n", entry.Entity.App.Name, entry.Entity.App.ID)
		case "workspace":
			fmt.Printf("  Entity: %s (%s)\n", entry.Entity.Workspace.Name, entry.Entity.Workspace.Domain)
		case "enterprise":
			fmt.Printf("  Entity: %s (%s)\n", entry.Entity.Enterprise.Name, entry.Entity.Enterprise.Domain)
		}

		// Context
		if entry.Context.Location.Name != "" {
			fmt.Printf("  Location: %s (%s)\n", entry.Context.Location.Name, entry.Context.Location.Type)
		}
		if entry.Context.IPAddress != "" {
			fmt.Printf("  IP Address: %s\n", entry.Context.IPAddress)
		}

		fmt.Println()
	}

	if len(entries) == 0 {
		fmt.Println("No audit log entries found in the specified time range.")
		fmt.Println("This could mean:")
		fmt.Println("  - No actions were taken in the last 24 hours")
		fmt.Println("  - The token doesn't have the auditlogs:read scope")
		fmt.Println("  - The workspace is not part of an Enterprise Grid")
	}
}
