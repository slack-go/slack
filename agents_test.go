package slack

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateAgentConversation(t *testing.T) {
	for _, tt := range []struct {
		name        string
		channelName string
		options     []CreateAgentConversationOption
		want        url.Values
	}{
		{
			name: "defaults", channelName: "SDK probe",
			want: url.Values{"name": {"SDK probe"}},
		},
		{
			name: "all options", channelName: "SDK probe",
			options: []CreateAgentConversationOption{
				CreateAgentConversationOptionTeamID("T1234567890"),
				CreateAgentConversationOptionSessionID("ses_abc123"),
				CreateAgentConversationOptionIsPrivate(true),
				CreateAgentConversationOptionOriginChannelID("C1234567890"),
				CreateAgentConversationOptionOriginMessageTimestamp("1717171717.123456"),
			},
			want: url.Values{
				"name": {"SDK probe"}, "team_id": {"T1234567890"}, "session_id": {"ses_abc123"},
				"is_private": {"true"}, "origin_channel_id": {"C1234567890"}, "origin_message_ts": {"1717171717.123456"},
			},
		},
		{
			name: "origin name and privacy",
			options: []CreateAgentConversationOption{
				CreateAgentConversationOptionOriginChannelID("C1234567890"),
				CreateAgentConversationOptionOriginMessageTimestamp("1717171717.123456"),
			},
			want: url.Values{"origin_channel_id": {"C1234567890"}, "origin_message_ts": {"1717171717.123456"}},
		},
		{
			name: "explicit public channel", channelName: "SDK probe",
			options: []CreateAgentConversationOption{CreateAgentConversationOptionIsPrivate(false)},
			want:    url.Values{"name": {"SDK probe"}, "is_private": {"false"}},
		},
		{
			name: "nil and empty options", channelName: "SDK probe",
			options: []CreateAgentConversationOption{
				nil,
				CreateAgentConversationOptionTeamID(""),
				CreateAgentConversationOptionSessionID(""),
				CreateAgentConversationOptionOriginChannelID(""),
				CreateAgentConversationOptionOriginMessageTimestamp(""),
			},
			want: url.Values{"name": {"SDK probe"}},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, http.MethodPost, r.Method)
				assert.Equal(t, "/agents.conversations.create", r.URL.Path)
				assert.NoError(t, r.ParseForm())
				assert.Equal(t, "testing-token", r.PostForm.Get("token"))
				r.PostForm.Del("token")
				assert.Equal(t, tt.want, r.PostForm)
				w.Header().Set("Content-Type", "application/json")
				_, _ = w.Write([]byte(`{"ok":true,"channel_id":"C9876543210"}`))
			}))
			defer server.Close()
			api := New("testing-token", OptionAPIURL(server.URL+"/"))
			response, err := api.CreateAgentConversation(tt.channelName, tt.options...)
			require.NoError(t, err)
			require.NotNil(t, response)
			assert.True(t, response.Ok)
			assert.Equal(t, "C9876543210", response.ChannelID)
		})
	}
}

func TestArchiveAgentConversation(t *testing.T) {
	for _, tt := range []struct {
		name    string
		options []ArchiveAgentConversationOption
		want    url.Values
	}{
		{
			name: "defaults",
			want: url.Values{"channel_id": {"C9876543210"}},
		},
		{
			name:    "summary message",
			options: []ArchiveAgentConversationOption{ArchiveAgentConversationOptionSummaryMessageTimestamp("1717182000.456789")},
			want:    url.Values{"channel_id": {"C9876543210"}, "summary_message_ts": {"1717182000.456789"}},
		},
		{
			name:    "nil and empty options",
			options: []ArchiveAgentConversationOption{nil, ArchiveAgentConversationOptionSummaryMessageTimestamp("")},
			want:    url.Values{"channel_id": {"C9876543210"}},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, http.MethodPost, r.Method)
				assert.Equal(t, "/agents.conversations.archive", r.URL.Path)
				assert.NoError(t, r.ParseForm())
				assert.Equal(t, "testing-token", r.PostForm.Get("token"))
				r.PostForm.Del("token")
				assert.Equal(t, tt.want, r.PostForm)
				w.Header().Set("Content-Type", "application/json")
				_, _ = w.Write([]byte(`{"ok":true}`))
			}))
			defer server.Close()
			api := New("testing-token", OptionAPIURL(server.URL+"/"))
			require.NoError(t, api.ArchiveAgentConversation("C9876543210", tt.options...))
		})
	}
}

func TestAgentConversationLifecycleErrors(t *testing.T) {
	for _, method := range []struct {
		name       string
		slackError string
		call       func(*Client, context.Context) error
	}{
		{
			name: "create", slackError: "feature_disabled",
			call: func(api *Client, ctx context.Context) error {
				_, err := api.CreateAgentConversationContext(ctx, "SDK probe")
				return err
			},
		},
		{
			name: "archive", slackError: "already_archived",
			call: func(api *Client, ctx context.Context) error {
				return api.ArchiveAgentConversationContext(ctx, "C9876543210")
			},
		},
	} {
		t.Run(method.name, func(t *testing.T) {
			for _, tt := range []struct {
				name   string
				status int
				body   string
			}{
				{"Slack error", http.StatusOK, `{"ok":false,"error":"` + method.slackError + `","response_metadata":{"messages":["Request rejected"]}}`},
				{"HTTP error", http.StatusInternalServerError, `server error`},
				{"invalid JSON", http.StatusOK, `{"ok":`},
			} {
				t.Run(tt.name, func(t *testing.T) {
					server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
						w.Header().Set("Content-Type", "application/json")
						w.WriteHeader(tt.status)
						_, _ = w.Write([]byte(tt.body))
					}))
					defer server.Close()
					api := New("testing-token", OptionAPIURL(server.URL+"/"))
					err := method.call(api, context.Background())
					require.Error(t, err)
					if tt.name == "Slack error" {
						var slackErr SlackErrorResponse
						require.ErrorAs(t, err, &slackErr)
						assert.Equal(t, method.slackError, slackErr.Err)
						assert.Equal(t, []string{"Request rejected"}, slackErr.ResponseMetadata.Messages)
					}
				})
			}
			t.Run("canceled context", func(t *testing.T) {
				server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					t.Error("canceled request reached server")
				}))
				defer server.Close()
				api := New("testing-token", OptionAPIURL(server.URL+"/"))
				ctx, cancel := context.WithCancel(context.Background())
				cancel()
				require.ErrorIs(t, method.call(api, ctx), context.Canceled)
			})
		})
	}
}

func TestSetAgentConversationCommands(t *testing.T) {
	for _, tt := range []struct {
		name     string
		commands []AgentCommand
		wantJSON string
	}{
		{
			name: "commands",
			commands: []AgentCommand{
				{Name: "create-pr", Description: "Open a pull request for the current branch", ArgumentHint: "[title]", ShouldEscape: true},
				{Name: "run-tests", Description: "Run the test suite and report back"},
			},
			wantJSON: `[
				{"name":"create-pr","description":"Open a pull request for the current branch","argument_hint":"[title]","should_escape":true},
				{"name":"run-tests","description":"Run the test suite and report back"}
			]`,
		},
		{
			name:     "optional fields omitted",
			commands: []AgentCommand{{Name: "summarize"}},
			wantJSON: `[{"name":"summarize"}]`,
		},
		{name: "empty clears commands", commands: []AgentCommand{}, wantJSON: `[]`},
		{name: "nil clears commands", wantJSON: `[]`},
	} {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, http.MethodPost, r.Method)
				assert.Equal(t, "/agents.conversations.setCommands", r.URL.Path)
				assert.NoError(t, r.ParseForm())
				assert.Equal(t, "testing-token", r.PostForm.Get("token"))
				assert.Equal(t, "C9876543210", r.PostForm.Get("channel_id"))
				assert.JSONEq(t, tt.wantJSON, r.PostForm.Get("commands"))
				w.Header().Set("Content-Type", "application/json")
				_, _ = fmt.Fprintf(w, `{"ok":true,"channel_id":"C9876543210","command_count":%d}`, len(tt.commands))
			}))
			defer server.Close()

			api := New("testing-token", OptionAPIURL(server.URL+"/"))
			response, err := api.SetAgentConversationCommands("C9876543210", tt.commands)
			require.NoError(t, err)
			require.NotNil(t, response)
			assert.True(t, response.Ok)
			assert.Equal(t, "C9876543210", response.ChannelID)
			assert.Equal(t, len(tt.commands), response.CommandCount)
		})
	}
}

func TestSetAgentConversationCommandsErrors(t *testing.T) {
	for _, tt := range []struct {
		name   string
		status int
		body   string
	}{
		{"Slack error", http.StatusOK, `{"ok":false,"error":"colliding_with_builtin","response_metadata":{"messages":["Choose another command name"]}}`},
		{"HTTP error", http.StatusInternalServerError, `server error`},
		{"invalid JSON", http.StatusOK, `{"ok":`},
	} {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tt.status)
				_, _ = w.Write([]byte(tt.body))
			}))
			defer server.Close()

			api := New("testing-token", OptionAPIURL(server.URL+"/"))
			response, err := api.SetAgentConversationCommandsContext(context.Background(), "C9876543210", []AgentCommand{{Name: "help"}})
			require.Error(t, err)
			if tt.name == "Slack error" {
				var slackErr SlackErrorResponse
				require.ErrorAs(t, err, &slackErr)
				assert.Equal(t, "colliding_with_builtin", slackErr.Err)
				assert.Equal(t, []string{"Choose another command name"}, slackErr.ResponseMetadata.Messages)
				require.NotNil(t, response)
				assert.False(t, response.Ok)
			} else {
				assert.Nil(t, response)
			}
		})
	}

	t.Run("canceled context", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			t.Error("canceled request reached server")
		}))
		defer server.Close()
		api := New("testing-token", OptionAPIURL(server.URL+"/"))
		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		response, err := api.SetAgentConversationCommandsContext(ctx, "C9876543210", nil)
		require.ErrorIs(t, err, context.Canceled)
		assert.Nil(t, response)
	})
}

func agentSessionSetStatusHandler(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Content-Type", "application/json")

	if r.FormValue("status") == "" || r.FormValue("channel_id") == "" {
		rw.Write([]byte(`{ "ok": false, "error": "invalid_arguments" }`))
		return
	}

	resp, _ := json.Marshal(&AgentSessionSetStatusResponse{
		SlackResponse: SlackResponse{Ok: true},
		Status:        r.FormValue("status"),
		AgentStatus:   r.FormValue("status"),
		Title:         r.FormValue("title"),
	})
	rw.Write(resp)
}

func TestSetAgentSessionStatus(t *testing.T) {
	http.HandleFunc("/agents.sessions.setStatus", agentSessionSetStatusHandler)
	once.Do(startServer)
	api := New("testing-token", OptionAPIURL("http://"+serverAddr+"/"))

	params := AgentSessionSetStatusParameters{
		Status:          AgentSessionStatusProcessing,
		ChannelID:       "CXXXXXXXX",
		ThreadTS:        "1234567890.123456",
		Title:           "Scuba diving research",
		InitiatorUserID: "UXXXXXXXX",
		IconEmoji:       ":thinking_face:",
		IconURL:         "https://example.com/icon.png",
		Username:        "custom name",
	}

	response, err := api.SetAgentSessionStatus(params)
	if err != nil {
		t.Fatalf("Unexpected error: %s", err)
	}

	if !response.Ok {
		t.Fatalf("Expected Ok to be true")
	}

	if response.Status != "processing" {
		t.Fatalf("Expected Status to be 'processing', got %q", response.Status)
	}

	if response.AgentStatus != "processing" {
		t.Fatalf("Expected AgentStatus to be 'processing', got %q", response.AgentStatus)
	}

	if response.Title != "Scuba diving research" {
		t.Fatalf("Expected Title to be 'Scuba diving research', got %q", response.Title)
	}

	_, err = api.SetAgentSessionStatus(AgentSessionSetStatusParameters{
		Status: AgentSessionStatusActive,
	})
	if err == nil {
		t.Fatalf("Expected error, got nil")
	}

	if err.Error() != "invalid_arguments" {
		t.Fatalf("Expected error 'invalid_arguments', got %q", err.Error())
	}
}

func agentSessionRenameHandler(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Content-Type", "application/json")

	if r.FormValue("title") == "" || r.FormValue("channel_id") == "" {
		rw.Write([]byte(`{ "ok": false, "error": "invalid_arguments" }`))
		return
	}

	resp, _ := json.Marshal(&AgentSessionRenameResponse{
		SlackResponse: SlackResponse{Ok: true},
		Title:         r.FormValue("title"),
	})
	rw.Write(resp)
}

func TestRenameAgentSession(t *testing.T) {
	http.HandleFunc("/agents.sessions.rename", agentSessionRenameHandler)
	once.Do(startServer)
	api := New("testing-token", OptionAPIURL("http://"+serverAddr+"/"))

	params := AgentSessionRenameParameters{
		Title:     "Bora Bora trip prep",
		ChannelID: "CXXXXXXXX",
		ThreadTS:  "1234567890.123456",
	}

	response, err := api.RenameAgentSession(params)
	if err != nil {
		t.Fatalf("Unexpected error: %s", err)
	}

	if !response.Ok {
		t.Fatalf("Expected Ok to be true")
	}

	if response.Title != "Bora Bora trip prep" {
		t.Fatalf("Expected Title to be 'Bora Bora trip prep', got %q", response.Title)
	}

	_, err = api.RenameAgentSession(AgentSessionRenameParameters{})
	if err == nil {
		t.Fatalf("Expected error, got nil")
	}

	if err.Error() != "invalid_arguments" {
		t.Fatalf("Expected error 'invalid_arguments', got %q", err.Error())
	}
}
