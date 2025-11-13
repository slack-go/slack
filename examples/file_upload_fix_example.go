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
	token := os.Getenv("SLACK_TOKEN")
	if token == "" {
		log.Fatal("SLACK_TOKEN environment variable is required")
	}

	// Create Slack client
	api := slack.New(token)

	// Example 1: Upload file using FixedUploadFile
	fmt.Println("=== Example 1: FixedUploadFile ===")
	err := uploadFileExample(api)
	if err != nil {
		log.Printf("Error in uploadFileExample: %v", err)
	}

	// Example 2: Upload file using FixedUploadFileV2
	fmt.Println("\n=== Example 2: FixedUploadFileV2 ===")
	err = uploadFileV2Example(api)
	if err != nil {
		log.Printf("Error in uploadFileV2Example: %v", err)
	}

	// Example 3: Simple helper methods
	fmt.Println("\n=== Example 3: Helper Methods ===")
	err = helperMethodsExample(api)
	if err != nil {
		log.Printf("Error in helperMethodsExample: %v", err)
	}

	// Example 4: Context usage
	fmt.Println("\n=== Example 4: Context Usage ===")
	err = contextExample(api)
	if err != nil {
		log.Printf("Error in contextExample: %v", err)
	}

	fmt.Println("\n=== All examples completed ===")
}

func uploadFileExample(api *slack.Client) error {
	// Create parameters for FixedUploadFile
	params := slack.FixedUploadFileParameters{
		Filename: "example.txt",
		Content:  "This is an example file uploaded using FixedUploadFile",
		Title:    "Example File",
		Channels: []string{"#general"},
		AltTxt:   "Example text file",
	}

	// Upload the file
	file, err := api.FixedUploadFile(params)
	if err != nil {
		return fmt.Errorf("failed to upload file: %w", err)
	}

	fmt.Printf("File uploaded successfully!\n")
	fmt.Printf("  ID: %s\n", file.ID)
	fmt.Printf("  Name: %s\n", file.Name)
	fmt.Printf("  Title: %s\n", file.Title)
	fmt.Printf("  Size: %d bytes\n", file.Size)

	return nil
}

func uploadFileV2Example(api *slack.Client) error {
	// Create parameters for FixedUploadFileV2
	params := slack.FixedUploadFileV2Parameters{
		Filename: "example_v2.txt",
		Content:  "This is an example file uploaded using FixedUploadFileV2",
		Title:    "Example File V2",
		Channel:  "#general",
		AltTxt:   "Example text file V2",
	}

	// Upload the file
	fileSummary, err := api.FixedUploadFileV2(params)
	if err != nil {
		return fmt.Errorf("failed to upload file V2: %w", err)
	}

	fmt.Printf("File V2 uploaded successfully!\n")
	fmt.Printf("  ID: %s\n", fileSummary.ID)
	fmt.Printf("  Title: %s\n", fileSummary.Title)

	return nil
}

func helperMethodsExample(api *slack.Client) error {
	// Example 1: Upload file from content
	content := "This is content uploaded using UploadFileFromContent helper"
	file1, err := api.UploadFileFromContent("helper_content.txt", content, "#general")
	if err != nil {
		return fmt.Errorf("failed to upload file from content: %w", err)
	}
	fmt.Printf("Helper method 1 - Content upload: %s\n", file1.Name)

	// Example 2: Upload file from path (if file exists)
	filePath := "example_file.txt"
	if _, err := os.Stat(filePath); err == nil {
		file2, err := api.UploadFileFromPath(filePath, "#general")
		if err != nil {
			return fmt.Errorf("failed to upload file from path: %w", err)
		}
		fmt.Printf("Helper method 2 - Path upload: %s\n", file2.Name)
	} else {
		fmt.Printf("Helper method 2 - Path upload: File %s not found, skipping\n", filePath)
	}

	// Example 3: Simple upload with reader
	reader := strings.NewReader("This is content from a reader")
	file3, err := api.SimpleUploadFile("reader_example.txt", reader, "#general")
	if err != nil {
		return fmt.Errorf("failed to upload file with reader: %w", err)
	}
	fmt.Printf("Helper method 3 - Reader upload: %s\n", file3.Name)

	return nil
}

func contextExample(api *slack.Client) error {
	// Create a context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Upload file with context
	params := slack.FixedUploadFileParameters{
		Filename: "context_example.txt",
		Content:  "This file was uploaded with a context timeout",
		Title:    "Context Example",
		Channels: []string{"#general"},
	}

	file, err := api.FixedUploadFileContext(ctx, params)
	if err != nil {
		return fmt.Errorf("failed to upload file with context: %w", err)
	}

	fmt.Printf("Context upload successful: %s\n", file.Name)
	return nil
}
