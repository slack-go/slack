package slack

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

var simpleMessage = `{
    "type": "message",
    "channel": "C2147483705",
    "user": "U2147483697",
    "text": "Hello world",
    "ts": "1355517523.000005"
}`

func unmarshalMessage(j string) (*Message, error) {
	message := &Message{}
	if err := json.Unmarshal([]byte(j), &message); err != nil {
		return nil, err
	}
	return message, nil
}

func TestSimpleMessage(t *testing.T) {
	message, err := unmarshalMessage(simpleMessage)
	assert.Nil(t, err)
	assert.NotNil(t, message)
	assert.Equal(t, "message", message.Type)
	assert.Equal(t, "C2147483705", message.Channel)
	assert.Equal(t, "U2147483697", message.User)
	assert.Equal(t, "Hello world", message.Text)
	assert.Equal(t, "1355517523.000005", message.Timestamp)
}

var starredMessage = `{
    "text": "is testing",
    "type": "message",
    "subtype": "me_message",
    "user": "U2147483697",
    "ts": "1433314126.000003",
    "is_starred": true
}`

func TestStarredMessage(t *testing.T) {
	message, err := unmarshalMessage(starredMessage)
	assert.Nil(t, err)
	assert.NotNil(t, message)
	assert.Equal(t, "is testing", message.Text)
	assert.Equal(t, "message", message.Type)
	assert.Equal(t, "me_message", message.SubType)
	assert.Equal(t, "U2147483697", message.User)
	assert.Equal(t, "1433314126.000003", message.Timestamp)
	assert.Equal(t, true, message.IsStarred)
}

var editedMessage = `{
    "type": "message",
    "user": "U2147483697",
    "text": "hello edited",
    "edited": {
        "user": "U2147483697",
        "ts": "1433314416.000000"
    },
    "ts": "1433314408.000004"
}`

func TestEditedMessage(t *testing.T) {
	message, err := unmarshalMessage(editedMessage)
	assert.Nil(t, err)
	assert.NotNil(t, message)
	assert.Equal(t, "message", message.Type)
	assert.Equal(t, "U2147483697", message.User)
	assert.Equal(t, "hello edited", message.Text)
	assert.NotNil(t, message.Edited)
	assert.Equal(t, "U2147483697", message.Edited.User)
	assert.Equal(t, "1433314416.000000", message.Edited.Timestamp)
	assert.Equal(t, "1433314408.000004", message.Timestamp)
}

var uploadedFile = `{
    "type": "message",
    "subtype": "file_share",
    "text": "<@U2147483697|tester> uploaded a file: <https:\/\/test.slack.com\/files\/tester\/abc\/test.txt|test.txt> and commented: test comment here",
    "file": {
        "id": "abc",
        "created": 1433314757,
        "timestamp": 1433314757,
        "name": "test.txt",
        "title": "test.txt",
        "mimetype": "text\/plain",
        "filetype": "text",
        "pretty_type": "Plain Text",
        "user": "U2147483697",
        "editable": true,
        "size": 5,
        "mode": "snippet",
        "is_external": false,
        "external_type": "",
        "is_public": true,
        "public_url_shared": false,
        "url": "https:\/\/slack-files.com\/files-pub\/abc-def-ghi\/test.txt",
        "url_download": "https:\/\/slack-files.com\/files-pub\/abc-def-ghi\/download\/test.txt",
        "url_private": "https:\/\/files.slack.com\/files-pri\/abc-def\/test.txt",
        "url_private_download": "https:\/\/files.slack.com\/files-pri\/abc-def\/download\/test.txt",
        "permalink": "https:\/\/test.slack.com\/files\/tester\/abc\/test.txt",
        "permalink_public": "https:\/\/slack-files.com\/abc-def-ghi",
        "edit_link": "https:\/\/test.slack.com\/files\/tester\/abc\/test.txt\/edit",
        "preview": "test\n",
        "preview_highlight": "<div class=\"sssh-code\"><div class=\"sssh-line\"><pre>test<\/pre><\/div>\n<div class=\"sssh-line\"><pre><\/pre><\/div>\n<\/div>",
        "lines": 2,
        "lines_more": 0,
        "channels": [
            "C2147483705"
        ],
        "groups": [],
        "ims": [],
        "comments_count": 1,
        "initial_comment": {
            "id": "Fc066YLGKH",
            "created": 1433314757,
            "timestamp": 1433314757,
            "user": "U2147483697",
            "comment": "test comment here"
        }
    },
    "user": "U2147483697",
    "upload": true,
    "ts": "1433314757.000006"
}`

func TestUploadedFile(t *testing.T) {
	message, err := unmarshalMessage(uploadedFile)
	assert.Nil(t, err)
	assert.NotNil(t, message)
	assert.Equal(t, "message", message.Type)
	assert.Equal(t, "file_share", message.SubType)
	assert.Equal(t, "<@U2147483697|tester> uploaded a file: <https://test.slack.com/files/tester/abc/test.txt|test.txt> and commented: test comment here", message.Text)
	// TODO: Assert File
	assert.Equal(t, "U2147483697", message.User)
	assert.True(t, message.Upload)
	assert.Equal(t, "1433314757.000006", message.Timestamp)
}

var testPost = `{
    "type": "message",
    "subtype": "file_share",
    "text": "<@U2147483697|tester> shared a file: <https:\/\/test.slack.com\/files\/tester\/abc\/test_post|test post>",
    "file": {
        "id": "abc",
        "created": 1433315398,
        "timestamp": 1433315398,
        "name": "test_post",
        "title": "test post",
        "mimetype": "text\/plain",
        "filetype": "post",
        "pretty_type": "Post",
        "user": "U2147483697",
        "editable": true,
        "size": 14,
        "mode": "post",
        "is_external": false,
        "external_type": "",
        "is_public": true,
        "public_url_shared": false,
        "url": "https:\/\/slack-files.com\/files-pub\/abc-def-ghi\/test_post",
        "url_download": "https:\/\/slack-files.com\/files-pub\/abc-def-ghi\/download\/test_post",
        "url_private": "https:\/\/files.slack.com\/files-pri\/abc-def\/test_post",
        "url_private_download": "https:\/\/files.slack.com\/files-pri\/abc-def\/download\/test_post",
        "permalink": "https:\/\/test.slack.com\/files\/tester\/abc\/test_post",
        "permalink_public": "https:\/\/slack-files.com\/abc-def-ghi",
        "edit_link": "https:\/\/test.slack.com\/files\/tester\/abc\/test_post\/edit",
        "preview": "test post body",
        "channels": [
            "C2147483705"
        ],
        "groups": [],
        "ims": [],
        "comments_count": 1
    },
    "user": "U2147483697",
    "upload": false,
    "ts": "1433315416.000008"
}`

func TestPost(t *testing.T) {
	message, err := unmarshalMessage(testPost)
	assert.Nil(t, err)
	assert.NotNil(t, message)
	assert.Equal(t, "message", message.Type)
	assert.Equal(t, "file_share", message.SubType)
	assert.Equal(t, "<@U2147483697|tester> shared a file: <https://test.slack.com/files/tester/abc/test_post|test post>", message.Text)
	// TODO: Assert File
	assert.Equal(t, "U2147483697", message.User)
	assert.False(t, message.Upload)
	assert.Equal(t, "1433315416.000008", message.Timestamp)
}

var testComment = `{
    "type": "message",
    "subtype": "file_comment",
    "text": "<@U2147483697|tester> commented on <@U2147483697|tester>'s file <https:\/\/test.slack.com\/files\/tester\/abc\/test_post|test post>: another comment",
    "file": {
        "id": "abc",
        "created": 1433315398,
        "timestamp": 1433315398,
        "name": "test_post",
        "title": "test post",
        "mimetype": "text\/plain",
        "filetype": "post",
        "pretty_type": "Post",
        "user": "U2147483697",
        "editable": true,
        "size": 14,
        "mode": "post",
        "is_external": false,
        "external_type": "",
        "is_public": true,
        "public_url_shared": false,
        "url": "https:\/\/slack-files.com\/files-pub\/abc-def-ghi\/test_post",
        "url_download": "https:\/\/slack-files.com\/files-pub\/abc-def-ghi\/download\/test_post",
        "url_private": "https:\/\/files.slack.com\/files-pri\/abc-def\/test_post",
        "url_private_download": "https:\/\/files.slack.com\/files-pri\/abc-def\/download\/test_post",
        "permalink": "https:\/\/test.slack.com\/files\/tester\/abc\/test_post",
        "permalink_public": "https:\/\/slack-files.com\/abc-def-ghi",
        "edit_link": "https:\/\/test.slack.com\/files\/tester\/abc\/test_post\/edit",
        "preview": "test post body",
        "channels": [
            "C2147483705"
        ],
        "groups": [],
        "ims": [],
        "comments_count": 2
    },
    "comment": {
        "id": "xyz",
        "created": 1433316360,
        "timestamp": 1433316360,
        "user": "U2147483697",
        "comment": "another comment"
    },
    "ts": "1433316360.000009"
}`

func TestComment(t *testing.T) {
	message, err := unmarshalMessage(testComment)
	fmt.Println(err)
	assert.Nil(t, err)
	assert.NotNil(t, message)
	assert.Equal(t, "message", message.Type)
	assert.Equal(t, "file_comment", message.SubType)
	assert.Equal(t, "<@U2147483697|tester> commented on <@U2147483697|tester>'s file <https://test.slack.com/files/tester/abc/test_post|test post>: another comment", message.Text)
	// TODO: Assert File
	// TODO: Assert Comment
	assert.Equal(t, "1433316360.000009", message.Timestamp)
}

var botMessage = `{
    "type": "message",
    "subtype": "bot_message",
    "ts": "1358877455.000010",
    "text": "Pushing is the answer",
    "bot_id": "BB12033",
    "username": "github",
    "icons": {}
}`

func TestBotMessage(t *testing.T) {
	message, err := unmarshalMessage(botMessage)
	assert.Nil(t, err)
	assert.NotNil(t, message)
	assert.Equal(t, "message", message.Type)
	assert.Equal(t, "bot_message", message.SubType)
	assert.Equal(t, "1358877455.000010", message.Timestamp)
	assert.Equal(t, "Pushing is the answer", message.Text)
	assert.Equal(t, "BB12033", message.BotID)
	assert.Equal(t, "github", message.Username)
	assert.NotNil(t, message.Icons)
	assert.Empty(t, message.Icons.IconURL)
	assert.Empty(t, message.Icons.IconEmoji)
}

var meMessage = `{
    "type": "message",
    "subtype": "me_message",
    "channel": "C2147483705",
    "user": "U2147483697",
    "text": "is doing that thing",
    "ts": "1355517523.000005"
}`

func TestMeMessage(t *testing.T) {
	message, err := unmarshalMessage(meMessage)
	assert.Nil(t, err)
	assert.NotNil(t, message)
	assert.Equal(t, "message", message.Type)
	assert.Equal(t, "me_message", message.SubType)
	assert.Equal(t, "C2147483705", message.Channel)
	assert.Equal(t, "U2147483697", message.User)
	assert.Equal(t, "is doing that thing", message.Text)
	assert.Equal(t, "1355517523.000005", message.Timestamp)
}
