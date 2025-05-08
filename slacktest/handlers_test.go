package slacktest

import (
	"testing"

	"github.com/stretchr/testify/assert"

	slack "github.com/slack-go/slack"
)

func TestAuthTestHandler(t *testing.T) {
	s := NewTestServer()
	go s.Start()

	client := slack.New("ABCDEFG", slack.OptionAPIURL(s.GetAPIURL()))
	user, err := client.AuthTest()
	assert.NoError(t, err, "should not error out")
	assert.Equal(t, defaultTeamName, user.Team, "user ID should be correct")
	assert.Equal(t, defaultTeamID, user.TeamID, "user ID should be correct")
	assert.Equal(t, defaultNonBotUserID, user.UserID, "user ID should be correct")
	assert.Equal(t, defaultNonBotUserName, user.User, "user ID should be correct")
}

func TestPostMessageHandler(t *testing.T) {
	s := NewTestServer()
	go s.Start()

	client := slack.New("ABCDEFG", slack.OptionAPIURL(s.GetAPIURL()))
	channel, tstamp, err := client.PostMessage("foo", slack.MsgOptionText("some text", false), slack.MsgOptionPostMessageParameters(slack.PostMessageParameters{}))
	assert.NoError(t, err, "should not error out")
	assert.Equal(t, "foo", channel, "channel should be correct")
	assert.NotEmpty(t, tstamp, "timestamp should not be empty")
}

func TestServerCreateConversationHandler(t *testing.T) {
	s := NewTestServer()
	go s.Start()

	client := slack.New("ABCDEFG", slack.OptionAPIURL(s.GetAPIURL()))
	conversation, err := client.CreateConversation(slack.CreateConversationParams{ChannelName: "test"})
	assert.NoError(t, err)
	assert.Equal(t, "C0EAQDV4Z", conversation.ID)
	assert.Equal(t, "U023BECGF", conversation.Creator)
	assert.Equal(t, "endeavor", conversation.Name)
}

func TestServerSetConversationTopicHandler(t *testing.T) {
	s := NewTestServer()
	go s.Start()

	client := slack.New("ABCDEFG", slack.OptionAPIURL(s.GetAPIURL()))
	response, err := client.SetTopicOfConversation("test", "Apply topically for best effects")
	assert.NoError(t, err)
	assert.Equal(t, "Apply topically for best effects", response.Topic.Value)
	assert.Equal(t, "U023BECGF", response.Topic.Creator)
	assert.NotEmpty(t, response.Topic.LastSet)
}

func TestServerSetConversationPurposeHandler(t *testing.T) {
	s := NewTestServer()
	go s.Start()

	client := slack.New("ABCDEFG", slack.OptionAPIURL(s.GetAPIURL()))
	response, err := client.SetPurposeOfConversation("test", "I didn't set this purpose on purpose!")
	assert.NoError(t, err)
	assert.Equal(t, "I didn't set this purpose on purpose!", response.Purpose.Value)
	assert.Equal(t, "U023BECGF", response.Purpose.Creator)
	assert.NotEmpty(t, response.Purpose.LastSet)
}

func TestServerRenameConversationHandler(t *testing.T) {
	s := NewTestServer()
	go s.Start()

	client := slack.New("ABCDEFG", slack.OptionAPIURL(s.GetAPIURL()))
	response, err := client.RenameConversation("ID", "newName")
	assert.NoError(t, err)
	assert.Equal(t, "newName", response.Name)
}

func TestServerInviteConversationHandler(t *testing.T) {
	s := NewTestServer()
	go s.Start()

	client := slack.New("ABCDEFG", slack.OptionAPIURL(s.GetAPIURL()))
	response, err := client.InviteUsersToConversation("conversationID", "username")
	assert.NoError(t, err)
	assert.Equal(t, 1, response.NumMembers)
}

func TestUserInfoHandler(t *testing.T) {
	s := NewTestServer()
	go s.Start()

	client := slack.New("ABCDEFG", slack.OptionAPIURL(s.GetAPIURL()))
	user, err := client.GetUserInfo("123456")
	assert.NoError(t, err)
	assert.Equal(t, "W012A3CDE", user.ID)
	assert.Equal(t, "spengler", user.Name)
	assert.True(t, user.IsAdmin)
}

func TestUserLookUpByEmailHandler(t *testing.T) {
	s := NewTestServer()
	go s.Start()

	client := slack.New("ABCDEFG", slack.OptionAPIURL(s.GetAPIURL()))
	user, err := client.GetUserByEmail("user@email.com")
	assert.NoError(t, err)
	assert.Equal(t, "W012A3CDE", user.ID)
	assert.Equal(t, "spengler", user.Name)
	assert.True(t, user.IsAdmin)
}

func TestBotInfoHandler(t *testing.T) {
	s := NewTestServer()
	go s.Start()

	client := slack.New("ABCDEFG", slack.OptionAPIURL(s.GetAPIURL()))
	bot, err := client.GetBotInfo(slack.GetBotInfoParameters{Bot: s.BotID})
	assert.NoError(t, err)
	assert.Equal(t, s.BotID, bot.ID)
	assert.Equal(t, s.BotName, bot.Name)
	assert.False(t, bot.Deleted)
}
