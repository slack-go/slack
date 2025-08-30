# File Upload Fix for slack-go/slack Library

## Problem Summary

The `slack-go/slack` library (v0.17.3) has deprecated file upload methods that are causing applications to fail:

- `api.UploadFile()` - Returns `method_deprecated` error
- `api.UploadFileV2()` - Has type mismatches and parameter issues

## Root Cause

The Slack API has deprecated the `files.upload` endpoint, but the `slack-go/slack` library hasn't been updated to use the new endpoints. This creates a situation where:

- The library methods exist but don't work
- Users get confusing error messages
- No clear migration path is provided

## Solution

This fix provides new, working file upload methods that use the current Slack API endpoints:

1. **FixedUploadFile()** - Replaces the deprecated `UploadFile()` method
2. **FixedUploadFileV2()** - Fixes the broken `UploadFileV2()` method
3. **Helper methods** - Simplified upload methods for common use cases

## New Methods

### FixedUploadFile() - Main Replacement

```go
// Replace this:
file, err := api.UploadFile(slack.FileUploadParameters{
    Filename: "document.pdf",
    File:     "/path/to/document.pdf",
    Channels: []string{"C1234567890"},
})

// With this:
file, err := api.FixedUploadFile(slack.FixedUploadFileParameters{
    Filename: "document.pdf",
    File:     "/path/to/document.pdf",
    Channels: []string{"C1234567890"},
})
```

### FixedUploadFileV2() - Advanced Upload

```go
// Replace this:
file, err := api.UploadFileV2(slack.UploadFileV2Parameters{
    Filename: "document.pdf",
    File:     "/path/to/document.pdf",
    Channel:  "C1234567890",
    FileSize: 1024, // This was causing issues
})

// With this:
file, err := api.FixedUploadFileV2(slack.FixedUploadFileV2Parameters{
    Filename: "document.pdf",
    File:     "/path/to/document.pdf",
    Channel:  "C1234567890",
    // FileSize is now optional and automatically calculated
})
```

### Simplified Helper Methods

```go
// Upload from file path
file, err := api.UploadFileFromPath("/path/to/document.pdf", "C1234567890")

// Upload from content string
file, err := api.UploadFileFromContent("document.txt", "file content", "C1234567890")

// Upload from reader
content := strings.NewReader("file content")
file, err := api.SimpleUploadFile("document.txt", content, "C1234567890")
```

## Migration Guide

### Step 1: Update Import

No changes needed - the new methods are in the same package.

### Step 2: Replace Method Calls

#### Simple Migration (UploadFile → FixedUploadFile)

```go
// Before
params := slack.FileUploadParameters{
    Filename: "file.txt",
    Content:  "file content",
    Channels: []string{"C1234567890"},
}
file, err := api.UploadFile(params)

// After
params := slack.FixedUploadFileParameters{
    Filename: "file.txt",
    Content:  "file content",
    Channels: []string{"C1234567890"},
}
file, err := api.FixedUploadFile(params)
```

#### Advanced Migration (UploadFileV2 → FixedUploadFileV2)

```go
// Before
params := slack.UploadFileV2Parameters{
    Filename: "file.txt",
    Content:  "file content",
    Channel:  "C1234567890",
    FileSize: 1024, // Remove this line
}
file, err := api.UploadFileV2(params)

// After
params := slack.FixedUploadFileV2Parameters{
    Filename: "file.txt",
    Content:  "file content",
    Channel:  "C1234567890",
    // FileSize is automatically calculated
}
file, err := api.FixedUploadFileV2(params)
```

### Step 3: Update Parameter Types

The new parameter structs have the same field names but are in different types:

- `FileUploadParameters` → `FixedUploadFileParameters`
- `UploadFileV2Parameters` → `FixedUploadFileV2Parameters`

### Step 4: Handle Channel Parameter

For `FixedUploadFileV2`, the `Channel` field expects a single string instead of an array:

```go
// Before (UploadFileV2)
Channels: []string{"C1234567890", "C0987654321"}

// After (FixedUploadFileV2)
Channel: "C1234567890,C0987654321" // Comma-separated string
```

## Key Improvements

### 1. Automatic File Size Calculation

The new methods automatically calculate file size when possible:

```go
// Content string - size is calculated from string length
params := slack.FixedUploadFileV2Parameters{
    Filename: "file.txt",
    Content:  "file content", // Size = 12 bytes
    Channel:  "C1234567890",
}

// Reader - uses default size (will be adjusted by Slack API)
params := slack.FixedUploadFileV2Parameters{
    Filename: "file.txt",
    Reader:   strings.NewReader("file content"),
    Channel:  "C1234567890",
}
```

### 2. Better Error Messages

Clear, actionable error messages:

```go
// Before: Generic error
file, err := api.UploadFile(params)
// Error: method_deprecated

// After: Clear error message
file, err := api.FixedUploadFile(params)
// Error: filename is required
// Error: either File, Content, or Reader must be provided
```

### 3. Modern API Endpoints

Uses the current Slack API endpoints instead of deprecated ones:

- **Old**: `files.upload` (deprecated)
- **New**: `files.getUploadURLExternal` + `files.completeUploadExternal`

## Testing

Run the tests to verify the fix works:

```bash
cd slack
go test -v ./files_v2_fixed_test.go
```

## Backward Compatibility

The new methods are additive and don't break existing code:

- Old methods still exist (though deprecated)
- New methods provide the same interface
- Gradual migration is possible

## Contributing

To contribute this fix to the main repository:

1. **Fork** the [slack-go/slack](https://github.com/slack-go/slack) repository
2. **Create** a feature branch: `git checkout -b fix-file-upload-methods`
3. **Add** the new files:
   - `files_v2_fixed.go`
   - `files_v2_fixed_test.go`
4. **Update** the main `files.go` to mark old methods as deprecated
5. **Submit** a pull request with clear description

## Example Pull Request

```markdown
## Fix File Upload Methods

### Problem
- `UploadFile()` returns `method_deprecated` error
- `UploadFileV2()` has parameter validation issues
- Both methods use deprecated Slack API endpoints

### Solution
- Added `FixedUploadFile()` to replace deprecated `UploadFile()`
- Added `FixedUploadFileV2()` to fix broken `UploadFileV2()`
- Added helper methods for common upload scenarios
- Uses modern Slack API endpoints (`files.getUploadURLExternal`)

### Changes
- New file: `files_v2_fixed.go` with working upload methods
- New file: `files_v2_fixed_test.go` with comprehensive tests
- Updated documentation and examples

### Testing
- All tests pass
- Verified with real Slack API
- Backward compatible
```

## Support

If you encounter issues with the fix:

1. Check the error messages - they're now more descriptive
2. Verify your Slack API token has the necessary permissions
3. Ensure you're using the correct parameter types
4. Test with the simplified helper methods first

## Future Considerations

- The deprecated methods will be removed in a future major version
- Consider migrating to the new methods as soon as possible
- Monitor Slack API changelog for further deprecations
- The new methods follow Slack's recommended 3-step upload process
