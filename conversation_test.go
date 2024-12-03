package slack

import (
	"encoding/json"
	"net/http"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Channel
var simpleChannel = `{
    "id": "C024BE91L",
    "name": "fun",
    "is_channel": true,
    "created": 1360782804,
    "creator": "U024BE7LH",
    "is_archived": false,
    "is_general": false,
    "members": [
        "U024BE7LH"
    ],
    "topic": {
        "value": "Fun times",
        "creator": "U024BE7LV",
        "last_set": 1369677212
    },
    "purpose": {
        "value": "This channel is for fun",
        "creator": "U024BE7LH",
        "last_set": 1360782804
    },
    "is_member": true,
    "last_read": "1401383885.000061",
    "unread_count": 0,
    "unread_count_display": 0
}`

func unmarshalChannel(j string) (*Channel, error) {
	channel := &Channel{}
	if err := json.Unmarshal([]byte(j), &channel); err != nil {
		return nil, err
	}
	return channel, nil
}

func TestSimpleChannel(t *testing.T) {
	channel, err := unmarshalChannel(simpleChannel)
	assert.Nil(t, err)
	assertSimpleChannel(t, channel)
}

func assertSimpleChannel(t *testing.T, channel *Channel) {
	assert.NotNil(t, channel)
	assert.Equal(t, "C024BE91L", channel.ID)
	assert.Equal(t, "fun", channel.Name)
	assert.Equal(t, true, channel.IsChannel)
	assert.Equal(t, JSONTime(1360782804), channel.Created)
	assert.Equal(t, "U024BE7LH", channel.Creator)
	assert.Equal(t, false, channel.IsArchived)
	assert.Equal(t, false, channel.IsGeneral)
	assert.Equal(t, true, channel.IsMember)
	assert.Equal(t, "1401383885.000061", channel.LastRead)
	assert.Equal(t, 0, channel.UnreadCount)
	assert.Equal(t, 0, channel.UnreadCountDisplay)
}

func TestCreateSimpleChannel(t *testing.T) {
	channel := &Channel{}
	channel.ID = "C024BE91L"
	channel.Name = "fun"
	channel.IsChannel = true
	channel.Created = JSONTime(1360782804)
	channel.Creator = "U024BE7LH"
	channel.IsArchived = false
	channel.IsGeneral = false
	channel.IsMember = true
	channel.LastRead = "1401383885.000061"
	channel.UnreadCount = 0
	channel.UnreadCountDisplay = 0
	assertSimpleChannel(t, channel)
}

// Shared Channel
var sharedChannel = `{
	"id": "C024BE91L",
	"name": "fun",
	"is_channel": true,
	"created": 1360782804,
	"creator": "U024BE7LH",
	"is_archived": false,
	"is_general": false,
	"members": [
			"U024BE7LH"
	],
	"is_shared": true,
	"context_team_id": "T1ABCD2E12",
	"is_ext_shared": true,
	"shared_team_ids": [
			"T07XY8FPJ5C"
	],
	"internal_team_ids": [],
	"connected_team_ids": [
			"T07XY8FPJ5C",
			"T1ABCD2E12"
	],
	"connected_limited_team_ids": [],
	"pending_connected_team_ids": [],
	"conversation_host_id": "T07XY8FPJ5C",
	"topic": {
			"value": "Fun times",
			"creator": "U024BE7LV",
			"last_set": 1369677212
	},
	"purpose": {
			"value": "This channel is for fun",
			"creator": "U024BE7LH",
			"last_set": 1360782804
	},
	"is_member": true,
	"last_read": "1401383885.000061",
	"unread_count": 0,
	"unread_count_display": 0
}`

func unmarshalSharedChannel(j string) (*Channel, error) {
	channel := &Channel{}
	if err := json.Unmarshal([]byte(j), &channel); err != nil {
		return nil, err
	}
	return channel, nil
}

func TestSharedChannel(t *testing.T) {
	channel, err := unmarshalSharedChannel(sharedChannel)
	assert.Nil(t, err)
	assertSharedChannel(t, channel)
}

func assertSharedChannel(t *testing.T, channel *Channel) {
	assertSimpleChannel(t, channel)
	assert.Equal(t, true, channel.IsShared)
	assert.Equal(t, true, channel.IsExtShared)
	assert.Equal(t, "T1ABCD2E12", channel.ContextTeamID)
	assert.Equal(t, "T07XY8FPJ5C", channel.ConversationHostID)
	if !reflect.DeepEqual([]string{"T07XY8FPJ5C"}, channel.SharedTeamIDs) {
		t.Fatal(ErrIncorrectResponse)
	}
	if !reflect.DeepEqual([]string{"T07XY8FPJ5C", "T1ABCD2E12"}, channel.ConnectedTeamIDs) {
		t.Fatal(ErrIncorrectResponse)
	}
}

func TestCreateSharedChannel(t *testing.T) {
	channel := &Channel{}
	channel.ID = "C024BE91L"
	channel.Name = "fun"
	channel.IsChannel = true
	channel.Created = JSONTime(1360782804)
	channel.Creator = "U024BE7LH"
	channel.IsArchived = false
	channel.IsGeneral = false
	channel.IsMember = true
	channel.LastRead = "1401383885.000061"
	channel.UnreadCount = 0
	channel.UnreadCountDisplay = 0
	channel.IsShared = true
	channel.IsExtShared = true
	channel.ContextTeamID = "T1ABCD2E12"
	channel.ConversationHostID = "T07XY8FPJ5C"
	channel.SharedTeamIDs = []string{"T07XY8FPJ5C"}
	channel.ConnectedTeamIDs = []string{"T07XY8FPJ5C", "T1ABCD2E12"}
	assertSharedChannel(t, channel)
}

// Group
var simpleGroup = `{
    "id": "G024BE91L",
    "name": "secretplans",
    "is_group": true,
    "created": 1360782804,
    "creator": "U024BE7LH",
    "is_archived": false,
    "members": [
        "U024BE7LH"
    ],
    "topic": {
        "value": "Secret plans on hold",
        "creator": "U024BE7LV",
        "last_set": 1369677212
    },
    "purpose": {
        "value": "Discuss secret plans that no-one else should know",
        "creator": "U024BE7LH",
        "last_set": 1360782804
    },
    "last_read": "1401383885.000061",
    "unread_count": 0,
    "unread_count_display": 0
}`

func unmarshalGroup(j string) (*Group, error) {
	group := &Group{}
	if err := json.Unmarshal([]byte(j), &group); err != nil {
		return nil, err
	}
	return group, nil
}

func TestSimpleGroup(t *testing.T) {
	group, err := unmarshalGroup(simpleGroup)
	assert.Nil(t, err)
	assertSimpleGroup(t, group)
}

func assertSimpleGroup(t *testing.T, group *Group) {
	assert.NotNil(t, group)
	assert.Equal(t, "G024BE91L", group.ID)
	assert.Equal(t, "secretplans", group.Name)
	assert.Equal(t, true, group.IsGroup)
	assert.Equal(t, JSONTime(1360782804), group.Created)
	assert.Equal(t, "U024BE7LH", group.Creator)
	assert.Equal(t, false, group.IsArchived)
	assert.Equal(t, "1401383885.000061", group.LastRead)
	assert.Equal(t, 0, group.UnreadCount)
	assert.Equal(t, 0, group.UnreadCountDisplay)
}

func TestCreateSimpleGroup(t *testing.T) {
	group := &Group{}
	group.ID = "G024BE91L"
	group.Name = "secretplans"
	group.IsGroup = true
	group.Created = JSONTime(1360782804)
	group.Creator = "U024BE7LH"
	group.IsArchived = false
	group.LastRead = "1401383885.000061"
	group.UnreadCount = 0
	group.UnreadCountDisplay = 0
	assertSimpleGroup(t, group)
}

// Channel with Canvas
var channelWithCanvas = `{
    "id": "C024BE91L",
    "name": "fun",
    "is_channel": true,
    "created": 1360782804,
    "creator": "U024BE7LH",
    "is_archived": false,
    "is_general": false,
    "members": [
        "U024BE7LH"
    ],
    "topic": {
        "value": "Fun times",
        "creator": "U024BE7LV",
        "last_set": 1369677212
    },
    "purpose": {
        "value": "This channel is for fun",
        "creator": "U024BE7LH",
        "last_set": 1360782804
    },
    "is_member": true,
    "last_read": "1401383885.000061",
    "unread_count": 0,
    "unread_count_display": 0,
	"properties": {
        "canvas": {
            "file_id": "F05RQ01LJU0",
            "is_empty": true,
            "quip_thread_id": "XFB9AAlvIyJ"
        }
    }
}`

func unmarshalChannelWithCanvas(j string) (*Channel, error) {
	channel := &Channel{}
	if err := json.Unmarshal([]byte(j), &channel); err != nil {
		return nil, err
	}
	return channel, nil
}

func TestChannelWithCanvas(t *testing.T) {
	channel, err := unmarshalChannelWithCanvas(channelWithCanvas)
	assert.Nil(t, err)
	assertChannelWithCanvas(t, channel)
}

func assertChannelWithCanvas(t *testing.T, channel *Channel) {
	assertSimpleChannel(t, channel)
	assert.Equal(t, "F05RQ01LJU0", channel.Properties.Canvas.FileId)
	assert.Equal(t, true, channel.Properties.Canvas.IsEmpty)
	assert.Equal(t, "XFB9AAlvIyJ", channel.Properties.Canvas.QuipThreadId)
}

func TestCreateChannelWithCanvas(t *testing.T) {
	channel := &Channel{}
	channel.ID = "C024BE91L"
	channel.Name = "fun"
	channel.IsChannel = true
	channel.Created = JSONTime(1360782804)
	channel.Creator = "U024BE7LH"
	channel.IsArchived = false
	channel.IsGeneral = false
	channel.IsMember = true
	channel.LastRead = "1401383885.000061"
	channel.UnreadCount = 0
	channel.UnreadCountDisplay = 0
	channel.Properties = &Properties{
		Canvas: Canvas{
			FileId:       "F05RQ01LJU0",
			IsEmpty:      true,
			QuipThreadId: "XFB9AAlvIyJ",
		},
	}
	assertChannelWithCanvas(t, channel)
}

// IM
var simpleIM = `{
    "id": "D024BFF1M",
    "is_im": true,
    "user": "U024BE7LH",
    "created": 1360782804,
    "is_user_deleted": false,
    "is_open": true,
    "last_read": "1401383885.000061",
    "unread_count": 0,
    "unread_count_display": 0
}`

func unmarshalIM(j string) (*IM, error) {
	im := &IM{}
	if err := json.Unmarshal([]byte(j), &im); err != nil {
		return nil, err
	}
	return im, nil
}

func TestSimpleIM(t *testing.T) {
	im, err := unmarshalIM(simpleIM)
	assert.Nil(t, err)
	assertSimpleIM(t, im)
}

func assertSimpleIM(t *testing.T, im *IM) {
	assert.NotNil(t, im)
	assert.Equal(t, "D024BFF1M", im.ID)
	assert.Equal(t, true, im.IsIM)
	assert.Equal(t, "U024BE7LH", im.User)
	assert.Equal(t, JSONTime(1360782804), im.Created)
	assert.Equal(t, false, im.IsUserDeleted)
	assert.Equal(t, true, im.IsOpen)
	assert.Equal(t, "1401383885.000061", im.LastRead)
	assert.Equal(t, 0, im.UnreadCount)
	assert.Equal(t, 0, im.UnreadCountDisplay)
}

func TestCreateSimpleIM(t *testing.T) {
	im := &IM{}
	im.ID = "D024BFF1M"
	im.IsIM = true
	im.User = "U024BE7LH"
	im.Created = JSONTime(1360782804)
	im.IsUserDeleted = false
	im.IsOpen = true
	im.LastRead = "1401383885.000061"
	im.UnreadCount = 0
	im.UnreadCountDisplay = 0
	assertSimpleIM(t, im)
}

func getTestMembers() []string {
	return []string{"test"}
}

func getUsersInConversation(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	response, _ := json.Marshal(struct {
		SlackResponse
		Members          []string         `json:"members"`
		ResponseMetaData responseMetaData `json:"response_metadata"`
	}{
		SlackResponse:    SlackResponse{Ok: true},
		Members:          getTestMembers(),
		ResponseMetaData: responseMetaData{NextCursor: ""},
	})
	rw.Write(response)
}

func TestGetUsersInConversation(t *testing.T) {
	http.HandleFunc("/conversations.members", getUsersInConversation)
	once.Do(startServer)
	api := New("testing-token", OptionAPIURL("http://"+serverAddr+"/"))
	params := GetUsersInConversationParameters{
		ChannelID: "CXXXXXXXX",
	}

	expectedMembers := getTestMembers()

	members, _, err := api.GetUsersInConversation(&params)
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
		return
	}
	if !reflect.DeepEqual(expectedMembers, members) {
		t.Fatal(ErrIncorrectResponse)
	}
}

func TestArchiveConversation(t *testing.T) {
	http.HandleFunc("/conversations.archive", okJSONHandler)
	once.Do(startServer)
	api := New("testing-token", OptionAPIURL("http://"+serverAddr+"/"))
	err := api.ArchiveConversation("CXXXXXXXX")
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
		return
	}
}

func TestUnArchiveConversation(t *testing.T) {
	http.HandleFunc("/conversations.unarchive", okJSONHandler)
	once.Do(startServer)
	api := New("testing-token", OptionAPIURL("http://"+serverAddr+"/"))
	err := api.UnArchiveConversation("CXXXXXXXX")
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
		return
	}
}

func getTestChannel() *Channel {
	return &Channel{
		GroupConversation: GroupConversation{
			Topic: Topic{
				Value: "response topic",
			},
			Purpose: Purpose{
				Value: "response purpose",
			},
		}}
}

func okChannelJsonHandler(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	response, _ := json.Marshal(struct {
		SlackResponse
		Channel *Channel `json:"channel"`
	}{
		SlackResponse: SlackResponse{Ok: true},
		Channel:       getTestChannel(),
	})
	rw.Write(response)
}

func okInviteSharedJsonHandler(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	response, _ := json.Marshal(struct {
		SlackResponse
		InviteID              string `json:"invite_id"`
		IsLegacySharedChannel bool   `json:"is_legacy_shared_channel"`
	}{
		SlackResponse:         SlackResponse{Ok: true},
		InviteID:              "I01234567",
		IsLegacySharedChannel: false,
	})
	rw.Write(response)
}

func TestSetTopicOfConversation(t *testing.T) {
	http.HandleFunc("/conversations.setTopic", okChannelJsonHandler)
	once.Do(startServer)
	api := New("testing-token", OptionAPIURL("http://"+serverAddr+"/"))
	inputChannel := getTestChannel()
	channel, err := api.SetTopicOfConversation("CXXXXXXXX", inputChannel.Topic.Value)
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
		return
	}
	if channel.Topic.Value != inputChannel.Topic.Value {
		t.Fatalf(`topic = '%s', want '%s'`, channel.Topic.Value, inputChannel.Topic.Value)
	}
}

func TestSetPurposeOfConversation(t *testing.T) {
	http.HandleFunc("/conversations.setPurpose", okChannelJsonHandler)
	once.Do(startServer)
	api := New("testing-token", OptionAPIURL("http://"+serverAddr+"/"))
	inputChannel := getTestChannel()
	channel, err := api.SetPurposeOfConversation("CXXXXXXXX", inputChannel.Purpose.Value)
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
		return
	}
	if channel.Purpose.Value != inputChannel.Purpose.Value {
		t.Fatalf(`purpose = '%s', want '%s'`, channel.Purpose.Value, inputChannel.Purpose.Value)
	}
}

func TestRenameConversation(t *testing.T) {
	http.HandleFunc("/conversations.rename", okChannelJsonHandler)
	once.Do(startServer)
	api := New("testing-token", OptionAPIURL("http://"+serverAddr+"/"))
	inputChannel := getTestChannel()
	channel, err := api.RenameConversation("CXXXXXXXX", inputChannel.Name)
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
		return
	}
	if channel.Name != inputChannel.Name {
		t.Fatalf(`channelName = '%s', want '%s'`, channel.Name, inputChannel.Name)
	}
}

func TestInviteUsersToConversation(t *testing.T) {
	http.HandleFunc("/conversations.invite", okChannelJsonHandler)
	once.Do(startServer)
	api := New("testing-token", OptionAPIURL("http://"+serverAddr+"/"))
	users := []string{"UXXXXXXX1", "UXXXXXXX2"}
	channel, err := api.InviteUsersToConversation("CXXXXXXXX", users...)
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
		return
	}
	if channel == nil {
		t.Error("channel should not be nil")
		return
	}
}

func TestInviteSharedToConversation(t *testing.T) {
	http.HandleFunc("/conversations.inviteShared", okInviteSharedJsonHandler)
	once.Do(startServer)
	api := New("testing-token", OptionAPIURL("http://"+serverAddr+"/"))

	t.Run("user_ids", func(t *testing.T) {
		userIDs := []string{"UXXXXXXX1", "UXXXXXXX2"}
		inviteID, isLegacySharedChannel, err := api.InviteSharedUserIDsToConversation("CXXXXXXXX", userIDs...)
		if err != nil {
			t.Errorf("Unexpected error: %s", err)
			return
		}
		if inviteID == "" {
			t.Error("invite id should have a value")
			return
		}
		if isLegacySharedChannel {
			t.Error("is legacy shared channel should be false")
		}
	})

	t.Run("emails", func(t *testing.T) {
		emails := []string{"nopcoder@slack.com", "nopcoder@example.com"}
		inviteID, isLegacySharedChannel, err := api.InviteSharedEmailsToConversation("CXXXXXXXX", emails...)
		if err != nil {
			t.Errorf("Unexpected error: %s", err)
			return
		}
		if inviteID == "" {
			t.Error("invite id should have a value")
			return
		}
		if isLegacySharedChannel {
			t.Error("is legacy shared channel should be false")
		}
	})

	t.Run("external_limited", func(t *testing.T) {
		userIDs := []string{"UXXXXXXX1", "UXXXXXXX2"}
		externalLimited := true
		inviteID, isLegacySharedChannel, err := api.InviteSharedToConversation(InviteSharedToConversationParams{
			ChannelID:       "CXXXXXXXX",
			UserIDs:         userIDs,
			ExternalLimited: &externalLimited,
		})
		if err != nil {
			t.Errorf("Unexpected error: %s", err)
			return
		}
		if inviteID == "" {
			t.Error("invite id should have a value")
			return
		}
		if isLegacySharedChannel {
			t.Error("is legacy shared channel should be false")
		}
	})
}

func TestKickUserFromConversation(t *testing.T) {
	http.HandleFunc("/conversations.kick", okJSONHandler)
	once.Do(startServer)
	api := New("testing-token", OptionAPIURL("http://"+serverAddr+"/"))
	err := api.KickUserFromConversation("CXXXXXXXX", "UXXXXXXXX")
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
		return
	}
}

func closeConversationHandler(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	response, _ := json.Marshal(struct {
		SlackResponse
		NoOp          bool `json:"no_op"`
		AlreadyClosed bool `json:"already_closed"`
	}{
		SlackResponse: SlackResponse{Ok: true}})
	rw.Write(response)
}

func TestCloseConversation(t *testing.T) {
	http.HandleFunc("/conversations.close", closeConversationHandler)
	once.Do(startServer)
	api := New("testing-token", OptionAPIURL("http://"+serverAddr+"/"))
	_, _, err := api.CloseConversation("CXXXXXXXX")
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
		return
	}
}

func TestCreateConversation(t *testing.T) {
	http.HandleFunc("/conversations.create", okChannelJsonHandler)
	once.Do(startServer)
	api := New("testing-token", OptionAPIURL("http://"+serverAddr+"/"))
	channel, err := api.CreateConversation(CreateConversationParams{ChannelName: "CXXXXXXXX"})
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
		return
	}
	if channel == nil {
		t.Error("channel should not be nil")
		return
	}
}

func TestGetConversationInfo(t *testing.T) {
	http.HandleFunc("/conversations.info", okChannelJsonHandler)
	once.Do(startServer)
	api := New("testing-token", OptionAPIURL("http://"+serverAddr+"/"))
	channel, err := api.GetConversationInfo(&GetConversationInfoInput{
		ChannelID: "CXXXXXXXX",
	})
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
		return
	}
	if channel == nil {
		t.Error("channel should not be nil")
		return
	}

	// Nil Input Error
	api = New("testing-token", OptionAPIURL("http://"+serverAddr+"/"))
	_, err = api.GetConversationInfo(nil)
	if err == nil {
		t.Errorf("Unexpected pass where there should have been nil input error")
		return
	}

	// No Channel Error
	api = New("testing-token", OptionAPIURL("http://"+serverAddr+"/"))
	_, err = api.GetConversationInfo(&GetConversationInfoInput{})
	if err == nil {
		t.Errorf("Unexpected pass where there should have been missing channel error")
		return
	}
}

func leaveConversationHandler(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	response, _ := json.Marshal(struct {
		SlackResponse
		NotInChannel bool `json:"not_in_channel"`
	}{
		SlackResponse: SlackResponse{Ok: true}})
	rw.Write(response)
}

func TestLeaveConversation(t *testing.T) {
	http.HandleFunc("/conversations.leave", leaveConversationHandler)
	once.Do(startServer)
	api := New("testing-token", OptionAPIURL("http://"+serverAddr+"/"))
	_, err := api.LeaveConversation("CXXXXXXXX")
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
		return
	}
}

func getConversationRepliesHandler(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	response, _ := json.Marshal(struct {
		SlackResponse
		HasMore          bool `json:"has_more"`
		ResponseMetaData struct {
			NextCursor string `json:"next_cursor"`
		} `json:"response_metadata"`
		Messages []Message `json:"messages"`
	}{
		SlackResponse: SlackResponse{Ok: true},
		Messages:      []Message{}})
	rw.Write(response)
}

func TestGetConversationReplies(t *testing.T) {
	http.HandleFunc("/conversations.replies", getConversationRepliesHandler)
	once.Do(startServer)
	api := New("testing-token", OptionAPIURL("http://"+serverAddr+"/"))
	params := GetConversationRepliesParameters{
		ChannelID: "CXXXXXXXX",
		Timestamp: "1234567890.123456",
	}
	_, _, _, err := api.GetConversationReplies(&params)
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
		return
	}
}

func getConversationsHandler(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	response, _ := json.Marshal(struct {
		SlackResponse
		ResponseMetaData struct {
			NextCursor string `json:"next_cursor"`
		} `json:"response_metadata"`
		Channels []Channel `json:"channels"`
	}{
		SlackResponse: SlackResponse{Ok: true},
		Channels:      []Channel{}})
	rw.Write(response)
}

func TestGetConversations(t *testing.T) {
	http.HandleFunc("/conversations.list", getConversationsHandler)
	once.Do(startServer)
	api := New("testing-token", OptionAPIURL("http://"+serverAddr+"/"))
	params := GetConversationsParameters{}
	_, _, err := api.GetConversations(&params)
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
		return
	}
}

func openConversationHandler(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	response, _ := json.Marshal(struct {
		SlackResponse
		NoOp        bool     `json:"no_op"`
		AlreadyOpen bool     `json:"already_open"`
		Channel     *Channel `json:"channel"`
	}{
		SlackResponse: SlackResponse{Ok: true}})
	rw.Write(response)
}

func TestOpenConversation(t *testing.T) {
	http.HandleFunc("/conversations.open", openConversationHandler)
	once.Do(startServer)
	api := New("testing-token", OptionAPIURL("http://"+serverAddr+"/"))
	params := OpenConversationParameters{ChannelID: "CXXXXXXXX"}
	_, _, _, err := api.OpenConversation(&params)
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
		return
	}
}

func joinConversationHandler(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	response, _ := json.Marshal(struct {
		Channel          *Channel `json:"channel"`
		Warning          string   `json:"warning"`
		ResponseMetaData *struct {
			Warnings []string `json:"warnings"`
		} `json:"response_metadata"`
		SlackResponse
	}{
		SlackResponse: SlackResponse{Ok: true}})
	rw.Write(response)
}

func TestJoinConversation(t *testing.T) {
	http.HandleFunc("/conversations.join", joinConversationHandler)
	once.Do(startServer)
	api := New("testing-token", OptionAPIURL("http://"+serverAddr+"/"))
	_, _, _, err := api.JoinConversation("CXXXXXXXX")
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
		return
	}
}

func getConversationHistoryHandler(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	response, _ := json.Marshal(GetConversationHistoryResponse{
		SlackResponse: SlackResponse{Ok: true}})
	rw.Write(response)
}

func TestGetConversationHistory(t *testing.T) {
	http.HandleFunc("/conversations.history", getConversationHistoryHandler)
	once.Do(startServer)
	api := New("testing-token", OptionAPIURL("http://"+serverAddr+"/"))
	params := GetConversationHistoryParameters{ChannelID: "CXXXXXXXX"}
	_, err := api.GetConversationHistory(&params)
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
		return
	}
}

func markConversationHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	response, _ := json.Marshal(GetConversationHistoryResponse{
		SlackResponse: SlackResponse{Ok: true}})
	w.Write(response)
}

func TestMarkConversation(t *testing.T) {
	http.HandleFunc("/conversations.mark", markConversationHandler)
	once.Do(startServer)
	api := New("testing-token", OptionAPIURL("http://"+serverAddr+"/"))
	err := api.MarkConversation("CXXXXXXXX", "1401383885.000061")
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
		return
	}
}

func createChannelCanvasHandler(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	response, _ := json.Marshal(struct {
		SlackResponse
		CanvasID string `json:"canvas_id"`
	}{
		SlackResponse: SlackResponse{Ok: true},
		CanvasID:      "F05RQ01LJU0",
	})
	rw.Write(response)
}

func TestCreateChannelCanvas(t *testing.T) {
	http.HandleFunc("/conversations.canvases.create", createChannelCanvasHandler)
	once.Do(startServer)
	api := New("testing-token", OptionAPIURL("http://"+serverAddr+"/"))

	documentContent := DocumentContent{
		Type:     "markdown",
		Markdown: "> channel canvas!",
	}

	canvasID, err := api.CreateChannelCanvas("C1234567890", documentContent)
	if err != nil {
		t.Errorf("Failed to create channel canvas: %v", err)
		return
	}

	assert.Equal(t, "F05RQ01LJU0", canvasID)
}
