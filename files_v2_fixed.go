package slack

import (
	"context"
	"fmt"
	"io"
	"strings"
)

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
	AltTxt          string
	SnippetType     string
}

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
}

func (api *Client) FixedUploadFile(params FixedUploadFileParameters) (*File, error) {
	return api.FixedUploadFileContext(context.Background(), params)
}

func (api *Client) FixedUploadFileContext(ctx context.Context, params FixedUploadFileParameters) (*File, error) {
	if params.Filename == "" {
		return nil, fmt.Errorf("filename is required")
	}
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
	file, _, _, err := api.GetFileInfoContext(ctx, fileSummary.ID, 0, 0)
	if err != nil {
		return nil, fmt.Errorf("failed to get file info after upload: %w", err)
	}
	return file, nil
}

func (api *Client) FixedUploadFileV2(params FixedUploadFileV2Parameters) (*FileSummary, error) {
	return api.FixedUploadFileV2Context(context.Background(), params)
}

func (api *Client) FixedUploadFileV2Context(ctx context.Context, params FixedUploadFileV2Parameters) (*FileSummary, error) {
	if params.Filename == "" {
		return nil, fmt.Errorf("filename is required")
	}
	var fileSize int
	if params.Reader != nil {
		fileSize = 1024
	} else if params.File != "" {
		fileSize = 1024
	} else if params.Content != "" {
		fileSize = len(params.Content)
	} else {
		return nil, fmt.Errorf("either File, Content, or Reader must be provided")
	}

	uploadResponse, err := api.GetUploadURLExternalContext(ctx, GetUploadURLExternalParameters{
		AltTxt:      params.AltTxt,
		FileName:    params.Filename,
		FileSize:    fileSize,
		SnippetType: params.SnippetType,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get upload URL: %w", err)
	}

	err = api.UploadToURL(ctx, UploadToURLParameters{
		UploadURL: uploadResponse.UploadURL,
		Reader:    params.Reader,
		File:      params.File,
		Content:   params.Content,
		Filename:  params.Filename,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to upload file to URL: %w", err)
	}

	completeResponse, err := api.CompleteUploadExternalContext(ctx, CompleteUploadExternalParameters{
		Files: []FileSummary{{
			ID:    uploadResponse.FileID,
			Title: params.Title,
		}},
		Channel:         params.Channel,
		InitialComment:  params.InitialComment,
		ThreadTimestamp: params.ThreadTimestamp,
		Blocks:          params.Blocks,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to complete upload: %w", err)
	}
	if len(completeResponse.Files) != 1 {
		return nil, fmt.Errorf("expected 1 file, got %d", len(completeResponse.Files))
	}
	return &completeResponse.Files[0], nil
}

// Helper methods for common upload scenarios

func (api *Client) SimpleUploadFile(filename string, content io.Reader, channel string) (*File, error) {
	return api.SimpleUploadFileContext(context.Background(), filename, content, channel)
}

func (api *Client) SimpleUploadFileContext(ctx context.Context, filename string, content io.Reader, channel string) (*File, error) {
	return api.FixedUploadFileContext(ctx, FixedUploadFileParameters{
		Filename: filename,
		Reader:   content,
		Channels: []string{channel},
	})
}

func (api *Client) UploadFileFromPath(filePath, channel string) (*File, error) {
	return api.UploadFileFromPathContext(context.Background(), filePath, channel)
}

func (api *Client) UploadFileFromPathContext(ctx context.Context, filePath, channel string) (*File, error) {
	return api.FixedUploadFileContext(ctx, FixedUploadFileParameters{
		File:     filePath,
		Filename: filePath,
		Channels: []string{channel},
	})
}

func (api *Client) UploadFileFromContent(filename, content, channel string) (*File, error) {
	return api.UploadFileFromContentContext(context.Background(), filename, content, channel)
}

func (api *Client) UploadFileFromContentContext(ctx context.Context, filename, content, channel string) (*File, error) {
	return api.FixedUploadFileContext(ctx, FixedUploadFileParameters{
		Content:  content,
		Filename: filename,
		Channels: []string{channel},
	})
}
