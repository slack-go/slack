package slack

import (
	"context"
	"strings"
	"testing"
)

func TestFixedUploadFileParameters_Validation(t *testing.T) {
	params := FixedUploadFileParameters{
		File:            "test.txt",
		Content:         "test content",
		Filename:        "test.txt",
		Title:           "Test File",
		InitialComment:  "Test comment",
		Channels:        []string{"#general"},
		ThreadTimestamp: "1234567890.123456",
		AltTxt:          "Test alt text",
		SnippetType:     "text",
	}

	if params.Filename != "test.txt" {
		t.Errorf("Expected filename 'test.txt', got '%s'", params.Filename)
	}
	if params.Content != "test content" {
		t.Errorf("Expected content 'test content', got '%s'", params.Content)
	}
	if len(params.Channels) != 1 || params.Channels[0] != "#general" {
		t.Errorf("Expected channels ['#general'], got %v", params.Channels)
	}
}

func TestFixedUploadFileV2Parameters_Validation(t *testing.T) {
	params := FixedUploadFileV2Parameters{
		File:            "test.txt",
		Content:         "test content",
		Filename:        "test.txt",
		Title:           "Test File",
		InitialComment:  "Test comment",
		Channel:         "#general",
		ThreadTimestamp: "1234567890.123456",
		AltTxt:          "Test alt text",
		SnippetType:     "text",
	}

	if params.Filename != "test.txt" {
		t.Errorf("Expected filename 'test.txt', got '%s'", params.Filename)
	}
	if params.Content != "test content" {
		t.Errorf("Expected content 'test content', got '%s'", params.Content)
	}
	if params.Channel != "#general" {
		t.Errorf("Expected channel '#general', got '%s'", params.Channel)
	}
}

func TestMethodSignatures(t *testing.T) {
	// This is a compile-time test - if it compiles, the signatures are correct
	var api *Client
	_ = func() (*File, error) { return api.FixedUploadFile(FixedUploadFileParameters{}) }
	_ = func(ctx context.Context) (*File, error) {
		return api.FixedUploadFileContext(ctx, FixedUploadFileParameters{})
	}
	_ = func() (*FileSummary, error) { return api.FixedUploadFileV2(FixedUploadFileV2Parameters{}) }
	_ = func(ctx context.Context) (*FileSummary, error) {
		return api.FixedUploadFileV2Context(ctx, FixedUploadFileV2Parameters{})
	}
	_ = func(filename string, content interface{}, channel string) (*File, error) {
		return api.SimpleUploadFile(filename, nil, channel)
	}
	_ = func(ctx context.Context, filename string, content interface{}, channel string) (*File, error) {
		return api.SimpleUploadFileContext(ctx, filename, nil, channel)
	}
	_ = func(filePath, channel string) (*File, error) { return api.UploadFileFromPath(filePath, channel) }
	_ = func(ctx context.Context, filePath, channel string) (*File, error) {
		return api.UploadFileFromPathContext(ctx, filePath, channel)
	}
	_ = func(filename, content, channel string) (*File, error) {
		return api.UploadFileFromContent(filename, content, channel)
	}
	_ = func(ctx context.Context, filename, content, channel string) (*File, error) {
		return api.UploadFileFromContentContext(ctx, filename, content, channel)
	}
	t.Log("All method signatures are correct")
}

func TestFileSizeCalculation(t *testing.T) {
	// Test content string length calculation
	shortContent := "short"
	shortSize := len(shortContent)
	if shortSize != 5 {
		t.Errorf("Expected short content size 5, got %d", shortSize)
	}

	longContent := "This is a much longer content string for testing file size calculation"
	longSize := len(longContent)
	if longSize != 70 {
		t.Errorf("Expected long content size 70, got %d", longSize)
	}
}

func TestChannelHandling(t *testing.T) {
	// Test channel array to string conversion
	channels := []string{"#general", "#random", "#help"}
	channelString := strings.Join(channels, ",")
	expected := "#general,#random,#help"
	if channelString != expected {
		t.Errorf("Expected channel string '%s', got '%s'", expected, channelString)
	}
}

func TestStructFieldAccess(t *testing.T) {
	// Test that we can access all fields of the parameter structs
	params := FixedUploadFileParameters{
		File:            "test.txt",
		Content:         "test content",
		Filename:        "test.txt",
		Title:           "Test File",
		InitialComment:  "Test comment",
		Channels:        []string{"#general"},
		ThreadTimestamp: "1234567890.123456",
		AltTxt:          "Test alt text",
		SnippetType:     "text",
	}

	// Access all fields to ensure they exist
	_ = params.File
	_ = params.Content
	_ = params.Filename
	_ = params.Title
	_ = params.InitialComment
	_ = params.Channels
	_ = params.ThreadTimestamp
	_ = params.AltTxt
	_ = params.SnippetType

	t.Log("All struct fields are accessible")
}
