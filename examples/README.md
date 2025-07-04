# Slack Examples

This directory contains examples demonstrating how to use the slack-go library for various Slack API operations.

## Development Guidelines

### Environment Variables vs Command Line Arguments

When developing examples, follow these patterns for handling different types of data:

#### Environment Variables (Sensitive Data)

Use environment variables for **sensitive information** that should not be exposed in command history or process lists:

- **Tokens**: `SLACK_BOT_TOKEN`, `SLACK_APP_TOKEN`, `SLACK_USER_TOKEN`
- **Secrets**: `SLACK_SIGNING_SECRET`
- **URLs with credentials**: webhook URLs, etc.

**Pattern:**
```go
token := os.Getenv("SLACK_BOT_TOKEN")
if token == "" {
    fmt.Fprintf(os.Stderr, "SLACK_BOT_TOKEN environment variable is required\n")
    os.Exit(1)
}

// Optional: Validate token format
if !strings.HasPrefix(token, "xoxb-") {
    fmt.Fprintf(os.Stderr, "SLACK_BOT_TOKEN must be a bot token (xoxb-)\n")
    os.Exit(1)
}
```

**Common environment variables:**
- `SLACK_BOT_TOKEN` - Bot user OAuth token (starts with `xoxb-`)
- `SLACK_APP_TOKEN` - App-level token (starts with `xapp-`)
- `SLACK_USER_TOKEN` - User OAuth token (starts with `xoxp-`)
- `SLACK_SIGNING_SECRET` - For webhook signature verification

#### Command Line Arguments (Non-Sensitive Data)

Use command line flags for **operational parameters** that are safe to expose:

- **IDs**: channel IDs, user IDs, team IDs
- **Configuration**: timeouts, limits, modes
- **Options**: boolean flags, enum values

**Pattern:**
```go
import "flag"

var (
    channelID = flag.String("channel", "", "Channel ID (required)")
    userID    = flag.String("user", "", "User ID (required)")
    verbose   = flag.Bool("verbose", false, "Enable verbose logging")
)

func main() {
    flag.Parse()

    if *channelID == "" {
        fmt.Fprintf(os.Stderr, "Error: -channel flag is required\n")
        os.Exit(1)
    }

    if *userID == "" {
        fmt.Fprintf(os.Stderr, "Error: -user flag is required\n")
        os.Exit(1)
    }

    // Use the flags
    fmt.Printf("Operating on channel %s with user %s\n", *channelID, *userID)
}
```

### Error Handling Standards

**Environment variable errors:**
- Use `fmt.Fprintf(os.Stderr, ...)` for error output
- Include clear variable name in error message
- Exit with `os.Exit(1)` for missing required variables

**Command line argument errors:**
- Use `fmt.Fprintf(os.Stderr, ...)` for error output
- Include flag name in error message
- Exit with `os.Exit(1)` for missing required flags

### Security Considerations

1. **Never hardcode sensitive values** in example code
2. **Think of validating required environment variables** before using them
3. **Use token format validation** when applicable (e.g., `xoxb-` prefix for bot tokens)
4. **Keep sensitive data out of command line arguments** to prevent exposure in process lists

This pattern ensures consistent security practices across all examples and makes them easier to understand and use safely.
