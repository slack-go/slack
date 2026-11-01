package socketmode

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
)

func TestCodeChannelLifecycleDispatch(t *testing.T) {
	for _, tt := range []struct {
		name  string
		kind  slackevents.EventsAPIType
		inner string
		want  any
	}{
		{
			name: "stop whole channel without thread timestamp",
			kind: slackevents.AgentSessionStopped,
			inner: `{"type":"agent_session_stopped","channel":"C9876543210","user":"U2345678901",
				"streaming_message_ts":["1784744784.686639"],"event_ts":"1784744788.887216"}`,
			want: &slackevents.AgentSessionStoppedEvent{
				Type: "agent_session_stopped", Channel: "C9876543210", User: "U2345678901",
				StreamingMessageTimestamps: []string{"1784744784.686639"}, EventTimestamp: "1784744788.887216",
			},
		},
		{
			name: "stop specific thread",
			kind: slackevents.AgentSessionStopped,
			inner: `{"type":"agent_session_stopped","channel":"C9876543210","user":"U2345678901",
				"thread_ts":"1784744767.192879","streaming_message_ts":[],"event_ts":"1784744788.887216"}`,
			want: &slackevents.AgentSessionStoppedEvent{
				Type: "agent_session_stopped", Channel: "C9876543210", User: "U2345678901",
				ThreadTimestamp: "1784744767.192879", StreamingMessageTimestamps: []string{}, EventTimestamp: "1784744788.887216",
			},
		},
		{
			name: "rename whole channel in enterprise install",
			kind: slackevents.AgentSessionTitleChanged,
			inner: `{"type":"agent_session_title_changed","channel":"C9876543210","user":"U2345678901",
				"team_id":"T123ABC456","enterprise_id":"E123ABC456","title":"Fix login",
				"previous_title":"Investigate login","event_ts":"1784744788.887216"}`,
			want: &slackevents.AgentSessionTitleChangedEvent{
				Type: "agent_session_title_changed", Channel: "C9876543210", User: "U2345678901",
				TeamID: "T123ABC456", EnterpriseID: "E123ABC456", Title: "Fix login",
				PreviousTitle: "Investigate login", EventTimestamp: "1784744788.887216",
			},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			wire, err := json.Marshal(map[string]any{
				"type": "events_api", "envelope_id": "code-channel-envelope", "accepts_response_payload": false,
				"payload": map[string]any{
					"type": "event_callback", "team_id": "T123ABC456", "api_app_id": "A123ABC456",
					"event_id": "Ev123ABC456", "event": json.RawMessage(tt.inner),
				},
			})
			require.NoError(t, err)
			client := New(slack.New("testing-token"))
			event, err := client.parseEvent(wire)
			require.NoError(t, err)
			require.NotNil(t, event)
			assert.Equal(t, EventTypeEventsAPI, event.Type)
			assert.Equal(t, "code-channel-envelope", event.Request.EnvelopeID)
			parsed, ok := event.Data.(slackevents.EventsAPIEvent)
			require.True(t, ok)
			assert.Equal(t, string(tt.kind), parsed.InnerEvent.Type)
			assert.Equal(t, tt.want, parsed.InnerEvent.Data)

			handler := NewSocketmodeHandler(client)
			received := make(chan *Event, 1)
			handler.HandleEvents(tt.kind, func(event *Event, _ *Client) { received <- event })
			handler.HandleDefault(func(_ *Event, _ *Client) { received <- nil })
			handler.DispatchEvent(*event)
			select {
			case dispatched := <-received:
				require.NotNil(t, dispatched, "event reached the fallback instead of its registered handler")
				assert.Equal(t, event.Request, dispatched.Request)
				assert.Equal(t, event.Data, dispatched.Data)
			case <-time.After(time.Second):
				t.Fatal("lifecycle event was not dispatched")
			}
		})
	}
}

func TestAgentViewInteractionDispatch(t *testing.T) {
	for _, tt := range []struct {
		name    string
		payload string
		kind    slack.InteractionType
	}{
		{
			name: "button handler by action ID",
			kind: slack.InteractionTypeBlockActions,
			payload: `{"type":"block_actions","channel":{"id":"C9876543210"},
				"container":{"type":"view","view_id":"Ct1234567890","channel_id":"C9876543210"},
				"actions":[{"type":"button","action_id":"approve","block_id":"decision","value":"approve"}]}`,
		},
		{
			name: "external select handler by interaction type",
			kind: slack.InteractionTypeBlockSuggestion,
			payload: `{"type":"block_suggestion","channel":{"id":"C9876543210"},
				"container":{"type":"view","view_id":"Ct1234567890","channel_id":"C9876543210"},
				"action_id":"repo_lookup","block_id":"repository","value":"billing"}`,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			wire, err := json.Marshal(map[string]any{
				"type": "interactive", "envelope_id": "view-envelope", "accepts_response_payload": true,
				"payload": json.RawMessage(tt.payload),
			})
			require.NoError(t, err)
			client := New(slack.New("testing-token"))
			event, err := client.parseEvent(wire)
			require.NoError(t, err)
			require.NotNil(t, event)
			assert.Equal(t, EventTypeInteractive, event.Type)
			assert.True(t, event.Request.AcceptsResponsePayload)
			handler := NewSocketmodeHandler(client)
			received := make(chan *Event, 1)
			callback := func(event *Event, _ *Client) { received <- event }
			if tt.kind == slack.InteractionTypeBlockActions {
				handler.HandleInteractionBlockAction("approve", callback)
			} else {
				handler.HandleInteraction(tt.kind, callback)
			}
			handler.HandleDefault(func(_ *Event, _ *Client) { received <- nil })
			handler.DispatchEvent(*event)
			select {
			case dispatched := <-received:
				require.NotNil(t, dispatched, "interaction reached the fallback instead of its registered handler")
				parsed, ok := dispatched.Data.(slack.InteractionCallback)
				require.True(t, ok)
				assert.Equal(t, tt.kind, parsed.Type)
				assert.Equal(t, "Ct1234567890", parsed.Container.ViewID)
				assert.Equal(t, "C9876543210", parsed.Container.ChannelID)
				assert.Equal(t, "view-envelope", dispatched.Request.EnvelopeID)
				if tt.kind == slack.InteractionTypeBlockSuggestion {
					assert.Equal(t, "billing", parsed.Value)
					assert.Equal(t, "repo_lookup", parsed.ActionID)
				}
			case <-time.After(time.Second):
				t.Fatal("view interaction was not dispatched")
			}
		})
	}
}

func TestAgentSlashCommandDispatch(t *testing.T) {
	wire := json.RawMessage(`{
		"type":"slash_commands","envelope_id":"command-envelope","accepts_response_payload":true,
		"payload":{"token":"XXYYZZ","team_id":"T123ABC456","channel_id":"C9876543210",
			"user_id":"U2345678901","command":"/create-pr","text":"Fix login",
			"api_app_id":"A123ABC456","is_enterprise_install":false,
			"response_url":"https://hooks.slack.com/commands/T123ABC456/1234567890/example",
			"trigger_id":"13345224609.738474920"}
	}`)
	client := New(slack.New("testing-token"))
	event, err := client.parseEvent(wire)
	require.NoError(t, err)
	require.NotNil(t, event)
	assert.Equal(t, EventTypeSlashCommand, event.Type)
	handler := NewSocketmodeHandler(client)
	received := make(chan *Event, 1)
	handler.HandleSlashCommand("/create-pr", func(event *Event, _ *Client) { received <- event })
	handler.HandleDefault(func(_ *Event, _ *Client) { received <- nil })
	handler.DispatchEvent(*event)
	select {
	case dispatched := <-received:
		require.NotNil(t, dispatched, "command reached the fallback instead of its registered handler")
		command, ok := dispatched.Data.(slack.SlashCommand)
		require.True(t, ok)
		assert.Equal(t, "/create-pr", command.Command)
		assert.Equal(t, "Fix login", command.Text)
		assert.Equal(t, "C9876543210", command.ChannelID)
		assert.Equal(t, "A123ABC456", command.APIAppID)
		assert.Equal(t, "13345224609.738474920", command.TriggerID)
		assert.Equal(t, "command-envelope", dispatched.Request.EnvelopeID)
	case <-time.After(time.Second):
		t.Fatal("agent slash command was not dispatched")
	}
}

func TestCodeChannelContextBarActionDispatch(t *testing.T) {
	// Based on a live context-bar click captured on 2026-10-28, with IDs and
	// credentials replaced. Slack delivered a slash command without the
	// is_enterprise_install field, rather than a code_channel_action event.
	wire := json.RawMessage(`{
		"type":"slash_commands","envelope_id":"context-bar-envelope","accepts_response_payload":true,
		"payload":{"api_app_id":"A123ABC456","channel_id":"C9876543210",
			"command":"/capture-test","response_url":"https://hooks.slack.com/commands/example",
			"team_id":"E123ABC456","text":"","trigger_id":"13345224609.738474920",
			"user_id":"U2345678901"}
	}`)
	client := New(slack.New("testing-token"))
	event, err := client.parseEvent(wire)
	require.NoError(t, err)
	require.NotNil(t, event)
	assert.Equal(t, EventTypeSlashCommand, event.Type)
	require.NotNil(t, event.Request)
	assert.Equal(t, "context-bar-envelope", event.Request.EnvelopeID)
	assert.True(t, event.Request.AcceptsResponsePayload)

	handler := NewSocketmodeHandler(client)
	received := make(chan *Event, 1)
	handler.HandleSlashCommand("/capture-test", func(event *Event, _ *Client) { received <- event })
	handler.HandleDefault(func(_ *Event, _ *Client) { received <- nil })
	handler.DispatchEvent(*event)
	select {
	case dispatched := <-received:
		require.NotNil(t, dispatched, "context-bar action reached the fallback instead of its slash-command handler")
		command, ok := dispatched.Data.(slack.SlashCommand)
		require.True(t, ok)
		assert.Equal(t, slack.SlashCommand{
			APIAppID: "A123ABC456", ChannelID: "C9876543210", Command: "/capture-test",
			ResponseURL: "https://hooks.slack.com/commands/example", TeamID: "E123ABC456",
			TriggerID: "13345224609.738474920", UserID: "U2345678901",
		}, command)
		assert.Equal(t, event.Request, dispatched.Request)
	case <-time.After(time.Second):
		t.Fatal("context-bar action was not dispatched")
	}
}
