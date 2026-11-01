package slack

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetConversationInfoAgentProperties(t *testing.T) {
	// Session and origin fields follow the Node response type and reference.
	// Context-bar ownership and the managed views also reflect live responses.
	const body = `{"ok":true,"channel":{"id":"C123ABC456","name":"SDK review","properties":{
		"record_channel":{"record_type":"agent_channel"},
		"agent_session":{
			"status":"active","title":"SDK review","agent_bot_user_ids":[],
			"encoded_agent_bot_user_ids":["U123ABC456"],"encoded_initiator_user_id":"U234ABC567",
			"origin_link":{"channel_id":"C234ABC567","ts":"1717171717.123456"}
		},
		"code_channel":{
			"context_bar_items":[
				{"key":"repo","label":"Repository","icon":"code","url":"https://github.com/slack-go/slack","bot_user_id":"U123ABC456"},
				{"key":"capture-test","label":"Test action","item_type":"action","bot_user_id":"U123ABC456"}
			],
			"summary_message":{"message_ts":"1717182000.456789","thread_ts":"1717171717.123456"},
			"agent_session_views":[
				{"view_id":"Ct1234567890","type":"block_kit","view_key":"review/status","name":"Status","content_version":1,"date_added":1793197616,"date_updated":1793197616,"file_id":"F1234567890","owner_bot_user_id":"U123ABC456"},
				{"view_id":"code","type":"diff","content_version":1,"date_added":1793197642,"date_updated":1793197642,"file_id":"F2345678901"},
				{"view_id":"Ct2345678901","type":"canvas","name":"Plan","date_added":1793199493,"date_updated":1793199493,"canvas_id":"F3456789012","owner_bot_user_id":"U123ABC456"}
			]
		}
	}}}`
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/conversations.info", r.URL.Path)
		if !assert.NoError(t, r.ParseForm()) {
			return
		}
		assert.Equal(t, "C123ABC456", r.PostForm.Get("channel"))
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(body))
	}))
	defer server.Close()
	api := New("testing-token", OptionAPIURL(server.URL+"/"))
	channel, err := api.GetConversationInfo(&GetConversationInfoInput{ChannelID: "C123ABC456"})
	require.NoError(t, err)
	require.NotNil(t, channel.Properties)
	assert.Equal(t, "agent_channel", channel.Properties.RecordChannel.RecordType)
	assert.Equal(t, &AgentSessionProperties{
		Status: "active", Title: "SDK review", AgentBotUserIDs: []string{},
		EncodedAgentBotUserIDs: []string{"U123ABC456"}, EncodedInitiatorUserID: "U234ABC567",
		OriginLink: &AgentSessionOriginLink{ChannelID: "C234ABC567", Timestamp: "1717171717.123456"},
	}, channel.Properties.AgentSession)
	require.NotNil(t, channel.Properties.CodeChannel)
	code := channel.Properties.CodeChannel
	assert.Equal(t, []AgentContextBarItemInfo{
		{AgentContextBarItem: AgentContextBarItem{Key: "repo", Label: "Repository", Icon: "code", URL: "https://github.com/slack-go/slack"}, BotUserID: "U123ABC456"},
		{AgentContextBarItem: AgentContextBarItem{Key: "capture-test", Label: "Test action", ItemType: "action"}, BotUserID: "U123ABC456"},
	}, code.ContextBarItems)
	assert.Equal(t, &AgentSummaryMessage{MessageTS: "1717182000.456789", ThreadTS: "1717171717.123456"}, code.SummaryMessage)
	assert.Equal(t, []AgentCodeChannelView{
		{AgentConversationSessionView: AgentConversationSessionView{ViewID: "Ct1234567890", Type: AgentViewTypeBlockKit, ViewKey: "review/status", Name: "Status", ContentVersion: 1, DateAdded: 1793197616, FileID: "F1234567890", OwnerBotUserID: "U123ABC456"}, DateUpdated: 1793197616},
		{AgentConversationSessionView: AgentConversationSessionView{ViewID: "code", Type: AgentViewTypeDiff, ContentVersion: 1, DateAdded: 1793197642, FileID: "F2345678901"}, DateUpdated: 1793197642},
		{AgentConversationSessionView: AgentConversationSessionView{ViewID: "Ct2345678901", Type: AgentViewTypeCanvas, Name: "Plan", DateAdded: 1793199493, OwnerBotUserID: "U123ABC456"}, DateUpdated: 1793199493, CanvasID: "F3456789012"},
	}, code.AgentSessionViews)
}

func TestChannelAgentPropertiesOptional(t *testing.T) {
	for _, body := range []string{
		`{"id":"C123ABC456"}`,
		`{"id":"C123ABC456","properties":{"record_channel":{"record_type":"Account"}}}`,
	} {
		var channel Channel
		require.NoError(t, json.Unmarshal([]byte(body), &channel))
		if channel.Properties != nil {
			assert.Nil(t, channel.Properties.AgentSession)
			assert.Nil(t, channel.Properties.CodeChannel)
		}
		encoded, err := json.Marshal(channel)
		require.NoError(t, err)
		assert.NotContains(t, string(encoded), `"agent_session"`)
		assert.NotContains(t, string(encoded), `"code_channel"`)
	}

	var channel Channel
	require.NoError(t, json.Unmarshal([]byte(`{"properties":{"agent_session":{},"code_channel":{}}}`), &channel))
	require.NotNil(t, channel.Properties.AgentSession)
	require.NotNil(t, channel.Properties.CodeChannel)
	assert.Empty(t, channel.Properties.CodeChannel.ContextBarItems)
	assert.Empty(t, channel.Properties.CodeChannel.AgentSessionViews)
}
