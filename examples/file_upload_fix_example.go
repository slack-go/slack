package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/slack-go/slack"
)

func main() {
	// Get Slack token from environment variable
	token := os.Getenv("SLACK_BOT_TOKEN")
	if token == "" {
		log.Fatal("SLACK_BOT_TOKEN environment variable is required")
	}

	// Create Slack client
	api := slack.New(token)

	// Example channel ID (replace with your actual channel ID)
	channelID := "C1234567890" // Replace with your channel ID

	fmt.Println("🚀 Testing Fixed File Upload Methods")
	fmt.Println("=====================================")

	// Test 1: Upload file from content string
	fmt.Println("\n1. Testing FixedUploadFile with content string...")
	err := testUploadFromContent(api, channelID)
	if err != nil {
		fmt.Printf("❌ Failed: %v\n", err)
	} else {
		fmt.Println("✅ Success!")
	}

	// Test 2: Upload file from reader
	fmt.Println("\n2. Testing FixedUploadFile with reader...")
	err = testUploadFromReader(api, channelID)
	if err != nil {
		fmt.Printf("❌ Failed: %v\n", err)
	} else {
		fmt.Println("✅ Success!")
	}

	// Test 3: Upload file from path (if file exists)
	fmt.Println("\n3. Testing FixedUploadFile with file path...")
	err = testUploadFromPath(api, channelID)
	if err != nil {
		fmt.Printf("❌ Failed: %v\n", err)
	} else {
		fmt.Println("✅ Success!")
	}

	// Test 4: Test FixedUploadFileV2
	fmt.Println("\n4. Testing FixedUploadFileV2...")
	err = testFixedUploadFileV2(api, channelID)
	if err != nil {
		fmt.Printf("❌ Failed: %v\n", err)
	} else {
		fmt.Println("✅ Success!")
	}

	// Test 5: Test helper methods
	fmt.Println("\n5. Testing helper methods...")
	err = testHelperMethods(api, channelID)
	if err != nil {
		fmt.Printf("❌ Failed: %v\n", err)
	} else {
		fmt.Println("✅ Success!")
	}

	fmt.Println("\n🎉 All tests completed!")
}

func testUploadFromContent(api *slack.Client, channelID string) error {
	content := fmt.Sprintf("Test file content generated at %s\nThis is a test file to verify the fixed upload methods work correctly.", time.Now().Format(time.RFC3339))

	params := slack.FixedUploadFileParameters{
		Filename: "test_content.txt",
		Content:  content,
		Title:    "Test File from Content",
		Channels: []string{channelID},
		AltTxt:   "A test text file with timestamp",
	}

	file, err := api.FixedUploadFile(params)
	if err != nil {
		return fmt.Errorf("FixedUploadFile failed: %w", err)
	}

	fmt.Printf("   📁 File uploaded successfully: %s (ID: %s)\n", file.Name, file.ID)
	return nil
}

func testUploadFromReader(api *slack.Client, channelID string) error {
	content := fmt.Sprintf("Test file content from reader generated at %s\nThis demonstrates uploading from an io.Reader interface.", time.Now().Format(time.RFC3339))
	reader := strings.NewReader(content)

	params := slack.FixedUploadFileParameters{
		Filename: "test_reader.txt",
		Reader:   reader,
		Title:    "Test File from Reader",
		Channels: []string{channelID},
		AltTxt:   "A test text file uploaded from reader",
	}

	file, err := api.FixedUploadFile(params)
	if err != nil {
		return fmt.Errorf("FixedUploadFile with reader failed: %w", err)
	}

	fmt.Printf("   📁 File uploaded successfully: %s (ID: %s)\n", file.Name, file.ID)
	return nil
}

func testUploadFromPath(api *slack.Client, channelID string) error {
	// Create a temporary test file
	tempFile := "temp_test_file.txt"
	content := fmt.Sprintf("Temporary test file created at %s\nThis file will be uploaded and then cleaned up.", time.Now().Format(time.RFC3339))

	err := os.WriteFile(tempFile, []byte(content), 0644)
	if err != nil {
		return fmt.Errorf("failed to create temp file: %w", err)
	}

	// Clean up the temp file after the function returns
	defer os.Remove(tempFile)

	params := slack.FixedUploadFileParameters{
		Filename: "test_path.txt",
		File:     tempFile,
		Title:    "Test File from Path",
		Channels: []string{channelID},
		AltTxt:   "A test text file uploaded from local path",
	}

	file, err := api.FixedUploadFile(params)
	if err != nil {
		return fmt.Errorf("FixedUploadFile with path failed: %w", err)
	}

	fmt.Printf("   📁 File uploaded successfully: %s (ID: %s)\n", file.Name, file.ID)
	return nil
}

func testFixedUploadFileV2(api *slack.Client, channelID string) error {
	content := fmt.Sprintf("Test file content for V2 method generated at %s\nThis tests the FixedUploadFileV2 method specifically.", time.Now().Format(time.RFC3339))

	params := slack.FixedUploadFileV2Parameters{
		Filename: "test_v2.txt",
		Content:  content,
		Title:    "Test File V2",
		Channel:  channelID,
		AltTxt:   "A test text file using V2 method",
	}

	file, err := api.FixedUploadFileV2(params)
	if err != nil {
		return fmt.Errorf("FixedUploadFileV2 failed: %w", err)
	}

	fmt.Printf("   📁 File uploaded successfully: %s (ID: %s)\n", file.Title, file.ID)
	return nil
}

func testHelperMethods(api *slack.Client, channelID string) error {
	// Test SimpleUploadFile
	content := strings.NewReader("Simple upload test content")
	filename := "simple_test.txt"

	file, err := api.SimpleUploadFile(filename, content, channelID)
	if err != nil {
		return fmt.Errorf("SimpleUploadFile failed: %w", err)
	}

	fmt.Printf("   📁 SimpleUploadFile: %s (ID: %s)\n", file.Name, file.ID)

	// Test UploadFileFromContent
	contentStr := "Content string test"
	filename2 := "content_test.txt"

	file2, err := api.UploadFileFromContent(filename2, contentStr, channelID)
	if err != nil {
		return fmt.Errorf("UploadFileFromContent failed: %w", err)
	}

	fmt.Printf("   📁 UploadFileFromContent: %s (ID: %s)\n", file2.Name, file2.ID)

	return nil
}

// Example of how to migrate from old methods to new methods
func migrationExample() {
	fmt.Println("\n🔄 Migration Example")
	fmt.Println("====================")

	// OLD WAY (deprecated and broken)
	fmt.Println("❌ OLD WAY (deprecated):")
	fmt.Println("   api.UploadFile(slack.FileUploadParameters{...})")
	fmt.Println("   api.UploadFileV2(slack.UploadFileV2Parameters{...})")

	// NEW WAY (fixed and working)
	fmt.Println("\n✅ NEW WAY (fixed):")
	fmt.Println("   api.FixedUploadFile(slack.FixedUploadFileParameters{...})")
	fmt.Println("   api.FixedUploadFileV2(slack.FixedUploadFileV2Parameters{...})")

	// SIMPLIFIED WAY (helper methods)
	fmt.Println("\n🚀 SIMPLIFIED WAY (helper methods):")
	fmt.Println("   api.SimpleUploadFile(filename, reader, channel)")
	fmt.Println("   api.UploadFileFromContent(filename, content, channel)")
	fmt.Println("   api.UploadFileFromPath(filepath, channel)")
}

// Example of error handling with the new methods
func errorHandlingExample(api *slack.Client, channelID string) {
	fmt.Println("\n⚠️  Error Handling Examples")
	fmt.Println("============================")

	// Test missing filename
	fmt.Println("Testing missing filename...")
	params := slack.FixedUploadFileParameters{
		Content:  "test content",
		Channels: []string{channelID},
		// Missing Filename - should return error
	}

	_, err := api.FixedUploadFile(params)
	if err != nil {
		fmt.Printf("   ✅ Expected error caught: %v\n", err)
	} else {
		fmt.Println("   ❌ Expected error but got none")
	}

	// Test missing content source
	fmt.Println("Testing missing content source...")
	params2 := slack.FixedUploadFileParameters{
		Filename: "test.txt",
		Channels: []string{channelID},
		// Missing Content, File, and Reader - should return error
	}

	_, err = api.FixedUploadFile(params2)
	if err != nil {
		fmt.Printf("   ✅ Expected error caught: %v\n", err)
	} else {
		fmt.Println("   ❌ Expected error but got none")
	}
}

// Example of using context for timeout and cancellation
func contextExample(api *slack.Client, channelID string) {
	fmt.Println("\n⏱️  Context Example")
	fmt.Println("====================")

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	content := "File uploaded with context and timeout"
	params := slack.FixedUploadFileParameters{
		Filename: "context_test.txt",
		Content:  content,
		Channels: []string{channelID},
	}

	file, err := api.FixedUploadFileContext(ctx, params)
	if err != nil {
		fmt.Printf("   ❌ Context upload failed: %v\n", err)
		return
	}

	fmt.Printf("   📁 Context upload successful: %s (ID: %s)\n", file.Name, file.ID)
}
