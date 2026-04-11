package slack

import (
	"encoding/json"
	"net/http"
	"testing"
)

func TestAssistantThreadsSuggestedPrompts(t *testing.T) {

	http.HandleFunc("/assistant.threads.setSuggestedPrompts", okJSONHandler)
	once.Do(startServer)
	api := New("testing-token", OptionAPIURL("http://"+serverAddr+"/"))

	params := AssistantThreadsSetSuggestedPromptsParameters{
		ChannelID: "CXXXXXXXX",
		ThreadTS:  "1234567890.123456",
	}

	params.AddPrompt("title1", "message1")
	params.AddPrompt("title2", "message2")

	err := api.SetAssistantThreadsSuggestedPrompts(params)
	if err != nil {
		t.Fatalf("Unexpected error: %s", err)
	}

}

func TestSetAssistantThreadsStatus(t *testing.T) {

	http.HandleFunc("/assistant.threads.setStatus", okJSONHandler)
	once.Do(startServer)
	api := New("testing-token", OptionAPIURL("http://"+serverAddr+"/"))

	params := AssistantThreadsSetStatusParameters{
		ChannelID:       "CXXXXXXXX",
		ThreadTS:        "1234567890.123456",
		Status:          "updated status",
		LoadingMessages: []string{"updating status..."},
	}

	err := api.SetAssistantThreadsStatus(params)
	if err != nil {
		t.Fatalf("Unexpected error: %s", err)
	}

}

func assistantThreadsTitleHandler(rw http.ResponseWriter, r *http.Request) {

	channelID := r.FormValue("channel_id")
	threadTS := r.FormValue("thread_ts")
	title := r.FormValue("title")

	rw.Header().Set("Content-Type", "application/json")

	if channelID != "" && threadTS != "" && title != "" {

		resp, _ := json.Marshal(&addBookmarkResponse{
			SlackResponse: SlackResponse{Ok: true},
		})
		rw.Write(resp)
	} else {
		rw.Write([]byte(`{ "ok": false, "error": "errored" }`))
	}

}

func TestSetAssistantThreadsTitle(t *testing.T) {

	http.HandleFunc("/assistant.threads.setTitle", assistantThreadsTitleHandler)
	once.Do(startServer)
	api := New("testing-token", OptionAPIURL("http://"+serverAddr+"/"))

	params := AssistantThreadsSetTitleParameters{
		ChannelID: "CXXXXXXXX",
		ThreadTS:  "1234567890.123456",
		Title:     "updated title",
	}

	err := api.SetAssistantThreadsTitle(params)
	if err != nil {
		t.Fatalf("Unexpected error: %s", err)
	}

}

func assistantSearchContextHandler(rw http.ResponseWriter, r *http.Request) {
	query := r.FormValue("query")
	rw.Header().Set("Content-Type", "application/json")

	if query != "" {
		resp, _ := json.Marshal(&AssistantSearchContextResponse{
			SlackResponse: SlackResponse{Ok: true},
			Results: AssistantSearchContextResults{
				Messages: []AssistantSearchContextMessage{
					{
						AuthorUserID: "U1234567890",
						TeamID:       "T1234567890",
						ChannelID:    "C1234567890",
						MessageTS:    "1234567890.123456",
						Content:      "This is a test message",
						IsAuthorBot:  false,
						Permalink:    "https://example.slack.com/archives/C1234567890/p1234567890123456",
					},
				},
			},
			ResponseMetadata: struct {
				NextCursor string `json:"next_cursor"`
			}{
				NextCursor: "next_cursor_value",
			},
		})
		rw.Write(resp)
	} else {
		rw.Write([]byte(`{ "ok": false, "error": "invalid_arguments" }`))
	}
}

func TestSearchAssistantContext(t *testing.T) {
	http.HandleFunc("/assistant.search.context", assistantSearchContextHandler)
	once.Do(startServer)
	api := New("testing-token", OptionAPIURL("http://"+serverAddr+"/"))

	params := AssistantSearchContextParameters{
		Query:            "test query",
		ActionToken:      "test_action_token",
		ChannelTypes:     []string{"public_channel", "private_channel"},
		ContentTypes:     []string{"messages"},
		ContextChannelID: "C1234567890",
		Cursor:           "cursor_value",
		IncludeBots:      true,
		Limit:            10,
	}

	response, err := api.SearchAssistantContext(params)
	if err != nil {
		t.Fatalf("Unexpected error: %s", err)
	}

	if !response.Ok {
		t.Fatalf("Expected Ok to be true")
	}

	if len(response.Results.Messages) != 1 {
		t.Fatalf("Expected 1 message, got %d", len(response.Results.Messages))
	}

	message := response.Results.Messages[0]
	if message.AuthorUserID != "U1234567890" {
		t.Fatalf("Expected AuthorUserID to be U1234567890, got %s", message.AuthorUserID)
	}

	if message.TeamID != "T1234567890" {
		t.Fatalf("Expected TeamID to be T1234567890, got %s", message.TeamID)
	}

	if message.ChannelID != "C1234567890" {
		t.Fatalf("Expected ChannelID to be C1234567890, got %s", message.ChannelID)
	}

	if message.MessageTS != "1234567890.123456" {
		t.Fatalf("Expected MessageTS to be '1234567890.123456', got %s", message.MessageTS)
	}

	if message.Content != "This is a test message" {
		t.Fatalf("Expected Content to be 'This is a test message', got %s", message.Content)
	}

	if message.IsAuthorBot != false {
		t.Fatalf("Expected IsAuthorBot to be false, got %v", message.IsAuthorBot)
	}

	if response.ResponseMetadata.NextCursor != "next_cursor_value" {
		t.Fatalf("Expected NextCursor to be 'next_cursor_value', got %s", response.ResponseMetadata.NextCursor)
	}
}

func TestSearchAssistantContextExpandedResponse(t *testing.T) {
	raw := `{
		"ok": true,
		"results": {
			"messages": [
				{
					"author_user_id": "U111",
					"author_name": "Test User",
					"team_id": "T111",
					"channel_id": "C111",
					"channel_name": "general",
					"message_ts": "1234567890.123456",
					"content": "Hello world",
					"is_author_bot": false,
					"permalink": "https://example.slack.com/archives/C111/p1234567890123456",
					"context_messages": {
						"before": [
							{
								"author_user_id": "U222",
								"team_id": "T111",
								"channel_id": "C111",
								"message_ts": "1234567889.000000",
								"content": "Before message",
								"is_author_bot": false,
								"permalink": "https://example.slack.com/archives/C111/p1234567889000000"
							}
						],
						"after": [
							{
								"author_user_id": "U333",
								"team_id": "T111",
								"channel_id": "C111",
								"message_ts": "1234567891.000000",
								"content": "After message",
								"is_author_bot": true,
								"permalink": "https://example.slack.com/archives/C111/p1234567891000000"
							}
						]
					}
				}
			],
			"files": [
				{
					"uploader_user_id": "U111",
					"author_user_id": "U111",
					"author_name": "Test User",
					"team_id": "T111",
					"file_id": "F111",
					"date_created": 1700000000,
					"date_updated": 1700001000,
					"title": "test.pdf",
					"file_type": "pdf",
					"permalink": "https://example.slack.com/files/U111/F111/test.pdf",
					"content": "File content excerpt"
				}
			],
			"channels": [
				{
					"team_id": "T111",
					"creator_user_id": "U111",
					"creator_name": "Test User",
					"date_created": 1600000000,
					"date_updated": 1700000000,
					"name": "general",
					"topic": "General discussion",
					"purpose": "Company-wide announcements",
					"permalink": "https://example.slack.com/archives/C111"
				}
			]
		},
		"response_metadata": {
			"next_cursor": "cursor123"
		}
	}`

	var response AssistantSearchContextResponse
	if err := json.Unmarshal([]byte(raw), &response); err != nil {
		t.Fatalf("Unmarshal error: %s", err)
	}

	if !response.Ok {
		t.Fatalf("Expected Ok to be true")
	}

	// Verify messages
	if len(response.Results.Messages) != 1 {
		t.Fatalf("Expected 1 message, got %d", len(response.Results.Messages))
	}

	msg := response.Results.Messages[0]
	if msg.AuthorName != "Test User" {
		t.Errorf("Expected AuthorName 'Test User', got %q", msg.AuthorName)
	}
	if msg.ChannelName != "general" {
		t.Errorf("Expected ChannelName 'general', got %q", msg.ChannelName)
	}

	// Verify context messages
	if msg.ContextMessages == nil {
		t.Fatal("Expected ContextMessages to be non-nil")
	}
	if len(msg.ContextMessages.Before) != 1 {
		t.Fatalf("Expected 1 before context message, got %d", len(msg.ContextMessages.Before))
	}
	if msg.ContextMessages.Before[0].Content != "Before message" {
		t.Errorf("Expected before content 'Before message', got %q", msg.ContextMessages.Before[0].Content)
	}
	if len(msg.ContextMessages.After) != 1 {
		t.Fatalf("Expected 1 after context message, got %d", len(msg.ContextMessages.After))
	}
	if !msg.ContextMessages.After[0].IsAuthorBot {
		t.Errorf("Expected after message IsAuthorBot true")
	}

	// Verify files
	if len(response.Results.Files) != 1 {
		t.Fatalf("Expected 1 file, got %d", len(response.Results.Files))
	}
	file := response.Results.Files[0]
	if file.FileID != "F111" {
		t.Errorf("Expected FileID 'F111', got %q", file.FileID)
	}
	if file.Title != "test.pdf" {
		t.Errorf("Expected Title 'test.pdf', got %q", file.Title)
	}
	if file.FileType != "pdf" {
		t.Errorf("Expected FileType 'pdf', got %q", file.FileType)
	}
	if file.DateCreated != 1700000000 {
		t.Errorf("Expected DateCreated 1700000000, got %d", file.DateCreated)
	}

	// Verify channels
	if len(response.Results.Channels) != 1 {
		t.Fatalf("Expected 1 channel, got %d", len(response.Results.Channels))
	}
	ch := response.Results.Channels[0]
	if ch.Name != "general" {
		t.Errorf("Expected Name 'general', got %q", ch.Name)
	}
	if ch.Topic != "General discussion" {
		t.Errorf("Expected Topic 'General discussion', got %q", ch.Topic)
	}
	if ch.Purpose != "Company-wide announcements" {
		t.Errorf("Expected Purpose 'Company-wide announcements', got %q", ch.Purpose)
	}

	// Verify cursor
	if response.ResponseMetadata.NextCursor != "cursor123" {
		t.Errorf("Expected NextCursor 'cursor123', got %q", response.ResponseMetadata.NextCursor)
	}
}
