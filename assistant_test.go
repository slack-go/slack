package slack

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
		Username:        "custom name",
		IconURL:         "https://example.com/icon.png",
		IconEmoji:       ":thinking_face:",
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
			SlackResponse: SlackResponse{
				Ok:               true,
				ResponseMetadata: ResponseMetadata{Cursor: "next_cursor_value"},
			},
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

	if response.ResponseMetadata.Cursor != "next_cursor_value" {
		t.Fatalf("Expected Cursor to be 'next_cursor_value', got %s", response.ResponseMetadata.Cursor)
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
								"author_name": "Before User",
								"text": "Before message",
								"ts": "1234567889.000000",
								"user_id": "U222"
							}
						],
						"after": [
							{
								"author_name": "After User",
								"text": "After message",
								"ts": "1234567891.000000",
								"user_id": "U333"
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
	if msg.ContextMessages.Before[0].Text != "Before message" {
		t.Errorf("Expected before text 'Before message', got %q", msg.ContextMessages.Before[0].Text)
	}
	if msg.ContextMessages.Before[0].UserID != "U222" {
		t.Errorf("Expected before UserID 'U222', got %q", msg.ContextMessages.Before[0].UserID)
	}
	if len(msg.ContextMessages.After) != 1 {
		t.Fatalf("Expected 1 after context message, got %d", len(msg.ContextMessages.After))
	}
	if msg.ContextMessages.After[0].TS != "1234567891.000000" {
		t.Errorf("Expected after TS '1234567891.000000', got %q", msg.ContextMessages.After[0].TS)
	}
	if msg.ContextMessages.After[0].AuthorName != "After User" {
		t.Errorf("Expected after AuthorName 'After User', got %q", msg.ContextMessages.After[0].AuthorName)
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
	if response.ResponseMetadata.Cursor != "cursor123" {
		t.Errorf("Expected Cursor 'cursor123', got %q", response.ResponseMetadata.Cursor)
	}
}

// TestSearchAssistantContextKeepsContextMessages checks that context messages
// survive decoding. They use text, user_id and ts, unlike top-level results.
func TestSearchAssistantContextKeepsContextMessages(t *testing.T) {
	// Keys from a live response and
	// https://docs.slack.dev/reference/methods/assistant.search.context
	const contextMessages = `{
		"before": [
			{"author_name": "Ada", "text": "before the match", "ts": "1763900000.000100", "user_id": "U012AB3CD"}
		],
		"after": [
			{"author_name": "Grace", "text": "after the match", "ts": "1763900000.000300", "user_id": "U045EF6GH"}
		]
	}`
	ts := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		rw.Header().Set("Content-Type", "application/json")
		_, _ = rw.Write([]byte(`{
			"ok": true,
			"results": {
				"messages": [
					{
						"author_name": "Ada",
						"author_user_id": "U012AB3CD",
						"channel_id": "C012AB3CD",
						"channel_name": "general",
						"content": "the match",
						"context_messages": ` + contextMessages + `,
						"is_author_bot": false,
						"message_ts": "1763900000.000200",
						"permalink": "https://example.slack.com/archives/C012AB3CD/p1763900000000200",
						"team_id": "T012AB3CD"
					}
				]
			}
		}`))
	}))
	t.Cleanup(ts.Close)
	api := New("testing-token", OptionAPIURL(ts.URL+"/"))

	response, err := api.SearchAssistantContext(AssistantSearchContextParameters{
		Query:                  "match",
		IncludeContextMessages: true,
	})
	require.NoError(t, err)
	require.Len(t, response.Results.Messages, 1)

	got, err := json.Marshal(response.Results.Messages[0].ContextMessages)
	require.NoError(t, err)
	assert.JSONEq(t, contextMessages, string(got))

	surrounding := response.Results.Messages[0].ContextMessages
	require.Len(t, surrounding.Before, 1)
	require.Len(t, surrounding.After, 1)
	assert.Equal(t, "U012AB3CD", surrounding.Before[0].UserID)
	assert.Equal(t, "Ada", surrounding.Before[0].AuthorName)
	assert.Equal(t, "before the match", surrounding.Before[0].Text)
	assert.Equal(t, "1763900000.000100", surrounding.Before[0].TS)
	assert.Equal(t, "U045EF6GH", surrounding.After[0].UserID)
	assert.Equal(t, "after the match", surrounding.After[0].Text)
	assert.Equal(t, "1763900000.000300", surrounding.After[0].TS)
}

func newAssistantSearchInfoTestClient(t *testing.T, body string) *Client {
	t.Helper()
	// The handler runs on the server goroutine, so it uses assert, not require.
	ts := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/assistant.search.info", r.URL.Path)
		assert.NoError(t, r.ParseForm())
		assert.Equal(t, url.Values{"token": {"testing-token"}}, r.PostForm)

		rw.Header().Set("Content-Type", "application/json")
		_, _ = rw.Write([]byte(body))
	}))
	t.Cleanup(ts.Close)
	return New("testing-token", OptionAPIURL(ts.URL+"/"))
}

func TestGetAssistantSearchInfo(t *testing.T) {
	ctx := context.Background()
	tests := []struct {
		name string
		call func(*Client) (*AssistantSearchInfoResponse, error)
	}{
		{"GetAssistantSearchInfo", func(api *Client) (*AssistantSearchInfoResponse, error) {
			return api.GetAssistantSearchInfo()
		}},
		{"GetAssistantSearchInfoContext", func(api *Client) (*AssistantSearchInfoResponse, error) {
			return api.GetAssistantSearchInfoContext(ctx)
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			api := newAssistantSearchInfoTestClient(t, `{"ok": true, "is_ai_search_enabled": true}`)

			response, err := tt.call(api)
			require.NoError(t, err)
			assert.True(t, response.IsAISearchEnabled)
		})
	}
}

func TestGetAssistantSearchInfoError(t *testing.T) {
	api := newAssistantSearchInfoTestClient(t, `{"ok": false, "error": "missing_scope"}`)

	response, err := api.GetAssistantSearchInfo()
	var slackErr SlackErrorResponse
	require.ErrorAs(t, err, &slackErr)
	assert.Equal(t, "missing_scope", slackErr.Err)
	require.NotNil(t, response)
	assert.False(t, response.IsAISearchEnabled)
}

func newAssistantSearchContextTestClient(t *testing.T, body string) *Client {
	t.Helper()
	ts := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/assistant.search.context", r.URL.Path)
		rw.Header().Set("Content-Type", "application/json")
		_, _ = rw.Write([]byte(body))
	}))
	t.Cleanup(ts.Close)
	return New("testing-token", OptionAPIURL(ts.URL+"/"))
}

func TestSearchAssistantContextErrorMetadata(t *testing.T) {
	api := newAssistantSearchContextTestClient(t, `{
		"ok": false,
		"error": "invalid_arguments",
		"response_metadata": {
			"messages": ["[ERROR] must provide a valid sort value [json-pointer:/sort]"]
		}
	}`)

	_, err := api.SearchAssistantContext(AssistantSearchContextParameters{Query: "test", Sort: "oldest"})
	var slackErr SlackErrorResponse
	require.ErrorAs(t, err, &slackErr)
	assert.Equal(t, "invalid_arguments", slackErr.Err)
	assert.Equal(t, []string{"[ERROR] must provide a valid sort value [json-pointer:/sort]"}, slackErr.ResponseMetadata.Messages)
}

func TestSearchAssistantContextResponseMetadata(t *testing.T) {
	api := newAssistantSearchContextTestClient(t, `{
		"ok": true,
		"warning": "superfluous_charset",
		"results": {"messages": []},
		"response_metadata": {
			"next_cursor": "Q1VSUkVOVF9QQUdFOjI=",
			"warnings": ["superfluous_charset"]
		}
	}`)

	response, err := api.SearchAssistantContext(AssistantSearchContextParameters{Query: "test"})
	require.NoError(t, err)
	assert.Equal(t, ResponseMetadata{
		Cursor:   "Q1VSUkVOVF9QQUdFOjI=",
		Warnings: []string{"superfluous_charset"},
	}, response.ResponseMetadata)
}

func TestSearchAssistantContextUsersAndChannels(t *testing.T) {
	users := `[
		{
			"user_id": "U012AB3CD",
			"full_name": "Jane Doe",
			"title": "Staff Engineer",
			"timezone": "Europe/London",
			"email": "jane.doe@example.com",
			"permalink": "https://example.slack.com/team/U012AB3CD",
			"profile_pic_permalink": "https://avatars.slack-edge.com/2024-01-15/1234567890123_0123456789abcdef0123_192.png"
		}
	]`
	channels := `[
		{
			"team_id": "T012AB3CD",
			"creator_user_id": "U012AB3CD",
			"creator_name": "Jane Doe",
			"date_created": 1700000000,
			"date_updated": 1710000000,
			"name": "project-gizmo",
			"topic": "Launch date: Q4 2025",
			"purpose": "Coordinate the Gizmo launch",
			"permalink": "https://example.slack.com/archives/C012AB3CD",
			"is_archived": false,
			"channel_type": "public_channel"
		}
	]`
	// The handler runs on the server goroutine, so it uses assert, not require.
	ts := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/assistant.search.context", r.URL.Path)
		assert.NoError(t, r.ParseForm())
		assert.Equal(t, "users,channels", r.PostForm.Get("content_types"))

		rw.Header().Set("Content-Type", "application/json")
		_, _ = rw.Write([]byte(`{"ok": true, "results": {"users": ` + users + `, "channels": ` + channels + `}}`))
	}))
	t.Cleanup(ts.Close)
	api := New("testing-token", OptionAPIURL(ts.URL+"/"))

	response, err := api.SearchAssistantContext(AssistantSearchContextParameters{
		Query:        "gizmo",
		ContentTypes: []string{"users", "channels"},
	})
	require.NoError(t, err)

	require.Len(t, response.Results.Users, 1)
	assert.Equal(t, "jane.doe@example.com", response.Results.Users[0].Email)
	require.Len(t, response.Results.Channels, 1)
	assert.Equal(t, "public_channel", response.Results.Channels[0].ChannelType)

	gotUsers, err := json.Marshal(response.Results.Users)
	require.NoError(t, err)
	assert.JSONEq(t, users, string(gotUsers))

	gotChannels, err := json.Marshal(response.Results.Channels)
	require.NoError(t, err)
	assert.JSONEq(t, channels, string(gotChannels))
}
