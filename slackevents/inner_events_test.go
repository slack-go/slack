package slackevents

import (
	"encoding/json"
	"testing"
)

func TestAppMention(t *testing.T) {
	rawE := []byte(`
			{
				"type": "app_mention",
				"user": "U061F7AUR",
				"text": "<@U0LAN0Z89> is it everything a river should be?",
				"ts": "1515449522.000016",
				"thread_ts": "1515449522.000016",
				"channel": "C0LAN2Q65",
				"event_ts": "1515449522000016"
		}
	`)
	err := json.Unmarshal(rawE, &AppMentionEvent{})
	if err != nil {
		t.Error(err)
	}
}

func TestAppUninstalled(t *testing.T) {
	rawE := []byte(`
		{
			"type": "app_uninstalled"
		}
	`)
	err := json.Unmarshal(rawE, &AppUninstalledEvent{})
	if err != nil {
		t.Error(err)
	}
}

func TestGridMigrationFinishedEvent(t *testing.T) {
	rawE := []byte(`
			{
				"type": "grid_migration_finished",
				"enterprise_id": "EXXXXXXXX"
			}
	`)
	err := json.Unmarshal(rawE, &GridMigrationFinishedEvent{})
	if err != nil {
		t.Error(err)
	}
}

func TestGridMigrationStartedEvent(t *testing.T) {
	rawE := []byte(`
			{
				"token": "XXYYZZ",
				"team_id": "TXXXXXXXX",
				"api_app_id": "AXXXXXXXXX",
				"event": {
						"type": "grid_migration_started",
						"enterprise_id": "EXXXXXXXX"
				},
				"type": "event_callback",
				"event_id": "EvXXXXXXXX",
				"event_time": 1234567890
		}
	`)
	err := json.Unmarshal(rawE, &GridMigrationStartedEvent{})
	if err != nil {
		t.Error(err)
	}
}

func TestLinkSharedEvent(t *testing.T) {
	rawE := []byte(`
			{
				"type": "link_shared",
				"channel": "Cxxxxxx",
				"user": "Uxxxxxxx",
				"message_ts": "123456789.9875",
				"links":
						[
								{
										"domain": "example.com",
										"url": "https://example.com/12345"
								},
								{
										"domain": "example.com",
										"url": "https://example.com/67890"
								},
								{
										"domain": "another-example.com",
										"url": "https://yet.another-example.com/v/abcde"
								}
						]
		}
	`)
	err := json.Unmarshal(rawE, &LinkSharedEvent{})
	if err != nil {
		t.Error(err)
	}
}

func TestMessageEvent(t *testing.T) {
	rawE := []byte(`
			{
				"type": "message",
				"channel": "G024BE91L",
				"user": "U2147483697",
				"text": "Live long and prospect.",
				"ts": "1355517523.000005",
				"event_ts": "1355517523.000005",
				"channel_type": "channel",
				"message": {
					"text": "To infinity and beyond.",
					"edited": {
						"user": "U2147483697",
						"ts": "1355517524.000000"
					}
				},
				"previous_message": {
					"text": "Live long and prospect."
				}
		}
	`)
	err := json.Unmarshal(rawE, &MessageEvent{})
	if err != nil {
		t.Error(err)
	}
}

func TestBotMessageEvent(t *testing.T) {
	rawE := []byte(`
			{
				"type": "message",
				"subtype": "bot_message",
				"ts": "1358877455.000010",
				"text": "Pushing is the answer",
				"bot_id": "BB12033",
				"username": "github",
				"icons": {}
		}
	`)
	err := json.Unmarshal(rawE, &MessageEvent{})
	if err != nil {
		t.Error(err)
	}
}

func TestPinAdded(t *testing.T) {
	rawE := []byte(`
			{
				"type": "pin_added",
				"user": "U061F7AUR",
				"item": {
					"type": "message",
					"channel":"C0LAN2Q65",
					"message":{
						"type":"message",
						"user":"U061F7AUR",
						"text": "<@U0LAN0Z89> is it everything a river should be?",
						"ts":"1539904112.000100",
						"pinned_to":["C0LAN2Q65"],
						"replace_original":false,
						"delete_original":false
					}
				},
				"channel_id":"C0LAN2Q65",
				"event_ts": "1515449522000016"
		}
	`)
	err := json.Unmarshal(rawE, &PinAddedEvent{})
	if err != nil {
		t.Error(err)
	}
}

func TestPinRemoved(t *testing.T) {
	rawE := []byte(`
			{
				"type": "pin_removed",
				"user": "U061F7AUR",
				"item": {
					"type": "message",
					"channel":"C0LAN2Q65",
					"message":{
						"type":"message",
						"user":"U061F7AUR",
						"text": "<@U0LAN0Z89> is it everything a river should be?",
						"ts":"1539904112.000100",
						"pinned_to":["C0LAN2Q65"],
						"replace_original":false,
						"delete_original":false
					}
				},
				"channel_id":"C0LAN2Q65",
				"event_ts": "1515449522000016"
		}
	`)
	err := json.Unmarshal(rawE, &PinRemovedEvent{})
	if err != nil {
		t.Error(err)
	}
}

func TestTokensRevoked(t *testing.T) {
	rawE := []byte(`
	{
		"type": "tokens_revoked",
		"tokens": {
				"oauth": [
						"OUXXXXXXXX"
				],
				"bot": [
						"BUXXXXXXXX"
				]
		}
	}
`)
	tre := TokensRevokedEvent{}
	err := json.Unmarshal(rawE, &tre)
	if err != nil {
		t.Error(err)
	}

	if tre.Type != "tokens_revoked" {
		t.Fail()
	}

	if len(tre.Tokens.Bot) != 1 || tre.Tokens.Bot[0] != "BUXXXXXXXX" {
		t.Fail()
	}

	if len(tre.Tokens.Oauth) != 1 || tre.Tokens.Oauth[0] != "OUXXXXXXXX" {
		t.Fail()
	}
}
