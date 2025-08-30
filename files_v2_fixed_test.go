package slack

import (
	"context"
	"strings"
	"testing"
)

func TestFixedUploadFileParameters_Validation(t *testing.T) {
	// Test parameter struct creation and field access
	params := FixedUploadFileParameters{
		Filename: "test.txt",
		Content:  "test content",
		Channels: []string{"C1234567890"},
		Title:    "Test File",
		AltTxt:   "Alternative text",
	}

	// Verify fields are set correctly
	if params.Filename != "test.txt" {
		t.Errorf("Expected filename 'test.txt', got '%s'", params.Filename)
	}

	if params.Content != "test content" {
		t.Errorf("Expected content 'test content', got '%s'", params.Content)
	}

	if len(params.Channels) != 1 || params.Channels[0] != "C1234567890" {
		t.Errorf("Expected channel 'C1234567890', got '%v'", params.Channels)
	}

	if params.Title != "Test File" {
		t.Errorf("Expected title 'Test File', got '%s'", params.Title)
	}

	if params.AltTxt != "Alternative text" {
		t.Errorf("Expected alt text 'Alternative text', got '%s'", params.AltTxt)
	}
}

func TestFixedUploadFileV2Parameters_Validation(t *testing.T) {
	// Test parameter struct creation and field access
	params := FixedUploadFileV2Parameters{
		Filename: "test.txt",
		Content:  "test content",
		Channel:  "C1234567890",
		Title:    "Test File V2",
		AltTxt:   "Alternative text V2",
	}

	// Verify fields are set correctly
	if params.Filename != "test.txt" {
		t.Errorf("Expected filename 'test.txt', got '%s'", params.Filename)
	}

	if params.Content != "test content" {
		t.Errorf("Expected content 'test content', got '%s'", params.Content)
	}

	if params.Channel != "C1234567890" {
		t.Errorf("Expected channel 'C1234567890', got '%s'", params.Channel)
	}

	if params.Title != "Test File V2" {
		t.Errorf("Expected title 'Test File V2', got '%s'", params.Title)
	}

	if params.AltTxt != "Alternative text V2" {
		t.Errorf("Expected alt text 'Alternative text V2', got '%s'", params.AltTxt)
	}
}

func TestMethodSignatures(t *testing.T) {
	// Test that the methods exist and have correct signatures
	// This is a compile-time test - if it compiles, the signatures are correct

	// Create a mock client for testing method signatures
	var api *Client

	// Test FixedUploadFile method signature
	_ = func() (*File, error) {
		params := FixedUploadFileParameters{
			Filename: "test.txt",
			Content:  "test content",
			Channels: []string{"C1234567890"},
		}
		return api.FixedUploadFile(params)
	}

	// Test FixedUploadFileContext method signature
	_ = func(ctx context.Context) (*File, error) {
		params := FixedUploadFileParameters{
			Filename: "test.txt",
			Content:  "test content",
			Channels: []string{"C1234567890"},
		}
		return api.FixedUploadFileContext(ctx, params)
	}

	// Test FixedUploadFileV2 method signature
	_ = func() (*FileSummary, error) {
		params := FixedUploadFileV2Parameters{
			Filename: "test.txt",
			Content:  "test content",
			Channel:  "C1234567890",
		}
		return api.FixedUploadFileV2(params)
	}

	// Test FixedUploadFileV2Context method signature
	_ = func(ctx context.Context) (*FileSummary, error) {
		params := FixedUploadFileV2Parameters{
			Filename: "test.txt",
			Content:  "test content",
			Channel:  "C1234567890",
		}
		return api.FixedUploadFileV2Context(ctx, params)
	}

	// Test helper method signatures
	_ = func() (*File, error) {
		content := strings.NewReader("test content")
		return api.SimpleUploadFile("test.txt", content, "C1234567890")
	}

	_ = func() (*File, error) {
		return api.UploadFileFromContent("test.txt", "test content", "C1234567890")
	}

	_ = func() (*File, error) {
		return api.UploadFileFromPath("/path/to/test.txt", "C1234567890")
	}

	t.Log("All method signatures are correct")
}

func TestParameterConversion(t *testing.T) {
	// Test that parameters can be converted between types correctly
	params := FixedUploadFileParameters{
		Filename: "test.txt",
		Content:  "test content",
		Channels: []string{"C1234567890", "C0987654321"},
		Title:    "Test File",
		AltTxt:   "Alternative text",
	}

	// Convert to V2 parameters
	v2Params := FixedUploadFileV2Parameters{
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
	}

	// Verify conversion
	if v2Params.Filename != params.Filename {
		t.Errorf("Expected filename %s, got %s", params.Filename, v2Params.Filename)
	}

	if v2Params.Channel != "C1234567890,C0987654321" {
		t.Errorf("Expected channel %s, got %s", "C1234567890,C0987654321", v2Params.Channel)
	}

	if v2Params.Title != params.Title {
		t.Errorf("Expected title %s, got %s", params.Title, v2Params.Title)
	}
}

func TestChannelHandling(t *testing.T) {
	// Test channel array to string conversion
	channels := []string{"C1234567890", "C0987654321"}
	channelString := strings.Join(channels, ",")

	if channelString != "C1234567890,C0987654321" {
		t.Errorf("Expected channel string 'C1234567890,C0987654321', got '%s'", channelString)
	}

	// Test single channel
	singleChannel := []string{"C1234567890"}
	singleChannelString := strings.Join(singleChannel, ",")

	if singleChannelString != "C1234567890" {
		t.Errorf("Expected single channel string 'C1234567890', got '%s'", singleChannelString)
	}

	// Test empty channels
	emptyChannels := []string{}
	emptyChannelString := strings.Join(emptyChannels, ",")

	if emptyChannelString != "" {
		t.Errorf("Expected empty channel string '', got '%s'", emptyChannelString)
	}
}

func TestFileSizeCalculation(t *testing.T) {
	// Test file size calculation logic
	content := "test content"
	expectedSize := len(content)

	if expectedSize != 12 {
		t.Errorf("Expected content size 12, got %d", expectedSize)
	}

	// Test with empty content
	emptyContent := ""
	emptySize := len(emptyContent)
	if emptySize != 0 {
		t.Errorf("Expected empty content size 0, got %d", emptySize)
	}

	// Test with longer content
	longContent := "This is a much longer content string for testing file size calculation"
	longSize := len(longContent)
	if longSize != 70 {
		t.Errorf("Expected long content size 70, got %d", longSize)
	}
}

func TestStructFieldAccess(t *testing.T) {
	// Test that all struct fields can be accessed and modified
	params := FixedUploadFileParameters{}

	// Test setting fields
	params.Filename = "test.txt"
	params.Content = "test content"
	params.Channels = []string{"C1234567890"}
	params.Title = "Test Title"
	params.AltTxt = "Alt Text"
	params.SnippetType = "text"

	// Test getting fields
	if params.Filename != "test.txt" {
		t.Errorf("Filename not set correctly")
	}

	if params.Content != "test content" {
		t.Errorf("Content not set correctly")
	}

	if len(params.Channels) != 1 || params.Channels[0] != "C1234567890" {
		t.Errorf("Channels not set correctly")
	}

	if params.Title != "Test Title" {
		t.Errorf("Title not set correctly")
	}

	if params.AltTxt != "Alt Text" {
		t.Errorf("AltTxt not set correctly")
	}

	if params.SnippetType != "text" {
		t.Errorf("SnippetType not set correctly")
	}
}
