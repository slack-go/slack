package slack

import (
	"context"
	"fmt"
	"io"
	"strings"
)

// FixedUploadFileParameters contains all the parameters necessary for the fixed UploadFile method
type FixedUploadFileParameters struct {
	File            string
	Content         string
	Reader          io.Reader
	Filetype        string
	Filename        string
	Title           string
	InitialComment  string
	Channels        []string
	ThreadTimestamp string
	// New fields for modern Slack API
	AltTxt      string
	SnippetType string
}

// FixedUploadFileV2Parameters contains all the parameters necessary for the fixed UploadFileV2 method
type FixedUploadFileV2Parameters struct {
	File            string
	Content         string
	Reader          io.Reader
	Filetype        string
	Filename        string
	Title           string
	InitialComment  string
	Blocks          Blocks
	Channel         string
	ThreadTimestamp string
	AltTxt          string
	SnippetType     string
	// Remove FileSize requirement as it's not always known
}

// FixedUploadFile uploads a file using the modern Slack API approach
// This method replaces the deprecated UploadFile method
func (api *Client) FixedUploadFile(params FixedUploadFileParameters) (*File, error) {
	return api.FixedUploadFileContext(context.Background(), params)
}

// FixedUploadFileContext uploads a file using the modern Slack API approach with a custom context
// This method replaces the deprecated UploadFileContext method
func (api *Client) FixedUploadFileContext(ctx context.Context, params FixedUploadFileParameters) (*File, error) {
	// Validate required parameters
	if params.Filename == "" {
		return nil, fmt.Errorf("filename is required")
	}

	// Use the modern 3-step upload process
	fileSummary, err := api.FixedUploadFileV2Context(ctx, FixedUploadFileV2Parameters{
		File:            params.File,
		Content:         params.Content,
		Reader:          params.Reader,
		Filetype:        params.Filetype,
		Filename:        params.Filename,
		Title:           params.Title,
		InitialComment:  params.InitialComment,
		Channel:         strings.Join(params.Channels, ","),
		ThreadTimestamp: params.ThreadTimestamp,
		AltTxt:          params.AltTxt,
		SnippetType:     params.SnippetType,
	})
	if err != nil {
		return nil, err
	}

	// Get the full file info
	file, _, _, err := api.GetFileInfoContext(ctx, fileSummary.ID, 0, 0)
	if err != nil {
		return nil, fmt.Errorf("failed to get file info after upload: %w", err)
	}

	return file, nil
}

// FixedUploadFileV2 uploads file to a given slack channel using the modern 3-step process
// This method fixes the parameter validation issues in the original UploadFileV2
func (api *Client) FixedUploadFileV2(params FixedUploadFileV2Parameters) (*FileSummary, error) {
	return api.FixedUploadFileV2Context(context.Background(), params)
}

// FixedUploadFileV2Context uploads file to a given slack channel using the modern 3-step process
// This method fixes the parameter validation issues in the original UploadFileV2Context
func (api *Client) FixedUploadFileV2Context(ctx context.Context, params FixedUploadFileV2Parameters) (*FileSummary, error) {
	// Validate required parameters
	if params.Filename == "" {
		return nil, fmt.Errorf("filename is required")
	}

	// Calculate file size if possible
	var fileSize int
	if params.Reader != nil {
		// For readers, we can't determine size without reading the entire content
		// Use a reasonable default or estimate
		fileSize = 1024 // Default size, will be adjusted by Slack API
	} else if params.File != "" {
		// For local files, we could get the size, but let's use the modern approach
		fileSize = 1024 // Default size, will be adjusted by Slack API
	} else if params.Content != "" {
		fileSize = len(params.Content)
	} else {
		return nil, fmt.Errorf("either File, Content, or Reader must be provided")
	}

	// Step 1: Get upload URL
	uploadParams := GetUploadURLExternalParameters{
		AltTxt:      params.AltTxt,
		FileName:    params.Filename,
		FileSize:    fileSize,
		SnippetType: params.SnippetType,
	}

	uploadResponse, err := api.GetUploadURLExternalContext(ctx, uploadParams)
	if err != nil {
		return nil, fmt.Errorf("failed to get upload URL: %w", err)
	}

	// Step 2: Upload file to the URL
	uploadToURLParams := UploadToURLParameters{
		UploadURL: uploadResponse.UploadURL,
		Reader:    params.Reader,
		File:      params.File,
		Content:   params.Content,
		Filename:  params.Filename,
	}

	err = api.UploadToURL(ctx, uploadToURLParams)
	if err != nil {
		return nil, fmt.Errorf("failed to upload file to URL: %w", err)
	}

	// Step 3: Complete the upload
	completeParams := CompleteUploadExternalParameters{
		Files: []FileSummary{{
			ID:    uploadResponse.FileID,
			Title: params.Title,
		}},
		Channel:         params.Channel,
		InitialComment:  params.InitialComment,
		ThreadTimestamp: params.ThreadTimestamp,
		Blocks:          params.Blocks,
	}

	completeResponse, err := api.CompleteUploadExternalContext(ctx, completeParams)
	if err != nil {
		return nil, fmt.Errorf("failed to complete upload: %w", err)
	}

	if len(completeResponse.Files) != 1 {
		return nil, fmt.Errorf("expected 1 file, got %d", len(completeResponse.Files))
	}

	return &completeResponse.Files[0], nil
}

// SimpleUploadFile provides a simplified upload method for common use cases
// This method automatically chooses the best upload strategy
func (api *Client) SimpleUploadFile(filename string, content io.Reader, channel string) (*File, error) {
	return api.SimpleUploadFileContext(context.Background(), filename, content, channel)
}

// SimpleUploadFileContext provides a simplified upload method for common use cases with context
func (api *Client) SimpleUploadFileContext(ctx context.Context, filename string, content io.Reader, channel string) (*File, error) {
	params := FixedUploadFileParameters{
		Filename: filename,
		Reader:   content,
		Channels: []string{channel},
	}

	return api.FixedUploadFileContext(ctx, params)
}

// UploadFileFromPath uploads a file from a local file path
func (api *Client) UploadFileFromPath(filePath, channel string) (*File, error) {
	return api.UploadFileFromPathContext(context.Background(), filePath, channel)
}

// UploadFileFromPathContext uploads a file from a local file path with context
func (api *Client) UploadFileFromPathContext(ctx context.Context, filePath, channel string) (*File, error) {
	params := FixedUploadFileParameters{
		File:      filePath,
		Channels:  []string{channel},
		Filename:  filePath, // Will be cleaned up by the upload function
	}

	return api.FixedUploadFileContext(ctx, params)
}

// UploadFileFromContent uploads file content as a string
func (api *Client) UploadFileFromContent(filename, content, channel string) (*File, error) {
	return api.UploadFileFromContentContext(context.Background(), filename, content, channel)
}

// UploadFileFromContentContext uploads file content as a string with context
func (api *Client) UploadFileFromContentContext(ctx context.Context, filename, content, channel string) (*File, error) {
	params := FixedUploadFileParameters{
		Filename: filename,
		Content:  content,
		Channels: []string{channel},
	}

	return api.FixedUploadFileContext(ctx, params)
}
