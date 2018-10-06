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
	APIURL = "http://" + serverAddr + "/"
	api := New("testing-token")
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
	APIURL = "http://" + serverAddr + "/"
	api := New("testing-token")
	err := api.ArchiveConversation("CXXXXXXXX")
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
		return
	}
}

func TestUnArchiveConversation(t *testing.T) {
	http.HandleFunc("/conversations.unarchive", okJSONHandler)
	once.Do(startServer)
	APIURL = "http://" + serverAddr + "/"
	api := New("testing-token")
	err := api.UnArchiveConversation("CXXXXXXXX")
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
		return
	}
}

func getTestChannel() *Channel {
	return &Channel{
		groupConversation: groupConversation{
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

func TestSetTopicOfConversation(t *testing.T) {
	http.HandleFunc("/conversations.setTopic", okChannelJsonHandler)
	once.Do(startServer)
	APIURL = "http://" + serverAddr + "/"
	api := New("testing-token")
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
	APIURL = "http://" + serverAddr + "/"
	api := New("testing-token")
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
	APIURL = "http://" + serverAddr + "/"
	api := New("testing-token")
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
	APIURL = "http://" + serverAddr + "/"
	api := New("testing-token")
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

func TestKickUserFromConversation(t *testing.T) {
	http.HandleFunc("/conversations.kick", okJSONHandler)
	once.Do(startServer)
	APIURL = "http://" + serverAddr + "/"
	api := New("testing-token")
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
	APIURL = "http://" + serverAddr + "/"
	api := New("testing-token")
	_, _, err := api.CloseConversation("CXXXXXXXX")
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
		return
	}
}

func TestCreateConversation(t *testing.T) {
	http.HandleFunc("/conversations.create", okChannelJsonHandler)
	once.Do(startServer)
	APIURL = "http://" + serverAddr + "/"
	api := New("testing-token")
	channel, err := api.CreateConversation("CXXXXXXXX", false)
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
	APIURL = "http://" + serverAddr + "/"
	api := New("testing-token")
	channel, err := api.GetConversationInfo("CXXXXXXXX", false)
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
		return
	}
	if channel == nil {
		t.Error("channel should not be nil")
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
	APIURL = "http://" + serverAddr + "/"
	api := New("testing-token")
	_, err := api.LeaveConversation("CXXXXXXXX")
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
		return
	}
}

func getConversationRepliesHander(rw http.ResponseWriter, r *http.Request) {
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
	http.HandleFunc("/conversations.replies", getConversationRepliesHander)
	once.Do(startServer)
	APIURL = "http://" + serverAddr + "/"
	api := New("testing-token")
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

func getConversationsHander(rw http.ResponseWriter, r *http.Request) {
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
	http.HandleFunc("/conversations.list", getConversationsHander)
	once.Do(startServer)
	APIURL = "http://" + serverAddr + "/"
	api := New("testing-token")
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
	APIURL = "http://" + serverAddr + "/"
	api := New("testing-token")
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
	APIURL = "http://" + serverAddr + "/"
	api := New("testing-token")
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
	APIURL = "http://" + serverAddr + "/"
	api := New("testing-token")
	params := GetConversationHistoryParameters{ChannelID: "CXXXXXXXX"}
	_, err := api.GetConversationHistory(&params)
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
		return
	}
}
