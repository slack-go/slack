# Contributing to the File Upload Fix

## Overview

This guide helps developers contribute to fixing the file upload methods in the `slack-go/slack` library. The fix addresses broken file upload functionality that's affecting users.

## Getting Started

### Prerequisites
- Go 1.22 or later
- Git
- A GitHub account
- Basic understanding of Go and the Slack API

### Setup
1. **Fork** the [slack-go/slack](https://github.com/slack-go/slack) repository
2. **Clone** your fork locally:
   ```bash
   git clone https://github.com/YOUR_USERNAME/slack.git
   cd slack
   ```
3. **Add** the upstream remote:
   ```bash
   git remote add upstream https://github.com/slack-go/slack.git
   ```

## Understanding the Problem

### Current Issues
- `UploadFile()` returns `method_deprecated` error
- `UploadFileV2()` has parameter validation issues
- Both methods use deprecated Slack API endpoints

### Root Cause
The Slack API has deprecated the `files.upload` endpoint, but the library hasn't been updated to use the new endpoints.

## Development Workflow

### 1. Create Feature Branch
```bash
git checkout -b fix-file-upload-methods
```

### 2. Add New Files
The fix includes several new files:
- `files_v2_fixed.go` - Core implementation
- `files_v2_fixed_test.go` - Tests
- `examples/file_upload_fix_example.go` - Examples

### 3. Update Existing Files
- `files.go` - Add deprecation warnings
- Documentation files

### 4. Test Your Changes
```bash
# Run all tests
go test -v ./...

# Run specific test files
go test -v ./files_v2_fixed_test.go

# Build examples
go build ./examples/file_upload_fix_example.go
```

### 5. Commit Your Changes
```bash
git add .
git commit -m "Fix file upload methods with modern Slack API endpoints

- Add FixedUploadFile() to replace deprecated UploadFile()
- Add FixedUploadFileV2() to fix broken UploadFileV2()
- Add helper methods for common upload scenarios
- Use modern Slack API endpoints (files.getUploadURLExternal)
- Maintain backward compatibility
- Add comprehensive tests and examples"
```

### 6. Push and Create Pull Request
```bash
git push origin fix-file-upload-methods
```

Then create a pull request on GitHub using the template in `PULL_REQUEST.md`.

## Code Standards

### Go Conventions
- Follow Go formatting standards (`gofmt`)
- Use meaningful variable and function names
- Add comments for exported functions
- Keep functions focused and concise

### Testing Requirements
- All new methods must have tests
- Tests should cover success and error cases
- Use descriptive test names
- Mock external dependencies when possible

### Documentation
- Update README files
- Add inline comments for complex logic
- Provide usage examples
- Document breaking changes

## Testing Strategy

### Unit Tests
- Parameter validation
- Error conditions
- Method signatures
- Edge cases

### Integration Tests
- Real Slack API calls (when possible)
- File upload scenarios
- Error handling
- Performance testing

### Test Commands
```bash
# Run all tests
go test -v ./...

# Run specific test files
go test -v ./files_v2_fixed_test.go

# Run tests with coverage
go test -v -cover ./...

# Run tests in parallel
go test -v -parallel 4 ./...
```

## Code Review Process

### Before Submitting
- [ ] All tests pass
- [ ] Code follows Go conventions
- [ ] Documentation is updated
- [ ] Examples work correctly
- [ ] No breaking changes (unless intentional)

### Review Checklist
- [ ] Code is readable and well-structured
- [ ] Error handling is robust
- [ ] Performance considerations addressed
- [ ] Security considerations addressed
- [ ] Backward compatibility maintained

## Common Issues and Solutions

### Import Errors
If you get import errors, ensure:
- All required packages are imported
- No unused imports
- Correct package paths

### Test Failures
- Check that all dependencies are available
- Verify test environment setup
- Look for timing issues in tests

### Compilation Errors
- Ensure Go version compatibility
- Check for syntax errors
- Verify all required types are defined

## Slack API Integration

### Required Permissions
The new methods require these Slack app permissions:
- `files:write` - Upload files
- `channels:read` - Read channel information
- `groups:read` - Read group information

### API Endpoints Used
- `files.getUploadURLExternal` - Get upload URL
- `files.completeUploadExternal` - Complete upload
- `files.info` - Get file information

### Rate Limiting
- Respect Slack API rate limits
- Implement proper error handling
- Use exponential backoff for retries

## Performance Considerations

### File Size Handling
- Small files (< 1MB): Direct content upload
- Large files (> 1MB): Streamed upload with reader
- Automatic size detection and optimization

### Memory Usage
- Streaming uploads for large files
- Minimal memory footprint
- Efficient parameter handling

### Network Efficiency
- Single API call for small files
- Optimized multipart uploads
- Proper timeout handling

## Security Considerations

### Token Handling
- Secure token storage
- Environment variable usage
- No hardcoded credentials

### File Validation
- File type checking
- Size limits
- Content validation

### API Permissions
- Required Slack app permissions
- Token scope validation
- Error handling for permission issues

## Troubleshooting

### Common Problems
1. **Tests failing**: Check Go version and dependencies
2. **Import errors**: Verify package paths and imports
3. **Compilation errors**: Check syntax and type definitions
4. **API errors**: Verify Slack token and permissions

### Getting Help
- Check existing issues on GitHub
- Review the documentation
- Ask questions in the PR discussion
- Join the Slack community

## Future Enhancements

### Potential Improvements
1. **Async Upload Support**: Background upload with progress callbacks
2. **Batch Upload**: Multiple file uploads in single operation
3. **Resumable Uploads**: Resume interrupted uploads
4. **Upload Progress**: Real-time upload progress tracking
5. **File Compression**: Automatic file compression for large files

### API Evolution
1. **Webhook Support**: File upload via webhooks
2. **OAuth Integration**: Enhanced OAuth flow for file uploads
3. **Enterprise Features**: Support for enterprise file policies
4. **Audit Logging**: File upload audit trails

## Contributing Guidelines

### Code of Conduct
- Be respectful and inclusive
- Help others learn and grow
- Focus on the code and technical issues
- Provide constructive feedback

### Communication
- Use clear, descriptive language
- Provide context for your changes
- Respond to feedback promptly
- Ask questions when unsure

### Recognition
- Contributors will be credited in the project
- Significant contributions may be highlighted
- All contributors are valued and appreciated

## Resources

### Documentation
- [Go Documentation](https://golang.org/doc/)
- [Slack API Documentation](https://api.slack.com/)
- [GitHub Pull Request Guide](https://docs.github.com/en/pull-requests)

### Community
- [Go Community](https://golang.org/community/)
- [Slack Developer Community](https://slack.dev/)
- [GitHub Discussions](https://github.com/slack-go/slack/discussions)

### Tools
- [Go Playground](https://play.golang.org/)
- [Go Modules](https://golang.org/ref/mod)
- [Go Testing](https://golang.org/pkg/testing/)

---

Thank you for contributing to the file upload fix! Your help makes the library better for everyone.
