package slackevents

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/slack-go/slack"
)

func TestParserOuterCallBackEvent(t *testing.T) {
	eventsAPIRawCallbackEvent := `
			{
				"token": "XXYYZZ",
				"team_id": "TXXXXXXXX",
				"api_app_id": "AXXXXXXXXX",
				"event": {
								"type": "app_mention",
								"event_ts": "1234567890.123456",
								"user": "UXXXXXXX1"
				},
				"type": "event_callback",
				"authed_users": [ "UXXXXXXX1" ],
				"event_id": "Ev08MFMKH6",
				"event_time": 1234567890
		}
	`
	msg, e := ParseEvent(json.RawMessage(eventsAPIRawCallbackEvent), OptionVerifyToken(&TokenComparator{"XXYYZZ"}))
	if e != nil {
		fmt.Println(e)
		t.Fail()
	}
	switch ev := msg.Data.(type) {
	case *EventsAPICallbackEvent:
		{
		}
	case *slack.UnmarshallingErrorEvent:
		{
			fmt.Println("Unmarshalling Error!")
			fmt.Println(ev)
			t.Fail()
		}
	default:
		{
			fmt.Println(ev)
			t.Fail()
		}
	}
}

func TestParseURLVerificationEvent(t *testing.T) {
	urlVerificationEvent := `
		{
			"token": "fake-token",
			"challenge": "aljdsflaji3jj",
			"type": "url_verification"
		}
	`
	msg, e := ParseEvent(json.RawMessage(urlVerificationEvent), OptionVerifyToken(&TokenComparator{"fake-token"}))
	if e != nil {
		fmt.Println(e)
		t.Fail()
	}
	switch ev := msg.Data.(type) {
	case *EventsAPIURLVerificationEvent:
		{
		}
	default:
		{
			fmt.Println(ev)
			t.Fail()
		}
	}
}

func TestParseAppRateLimitedEvent(t *testing.T) {
	event := `
		{
			"token": "fake-token",
			"team_id": "T123ABC456",
			"minute_rate_limited": 1518467820,
			"api_app_id": "A123ABC456",
			"type": "app_rate_limited"
		}
	`
	msg, e := ParseEvent(json.RawMessage(event), OptionVerifyToken(&TokenComparator{"fake-token"}))
	if e != nil {
		fmt.Println(e)
		t.Fail()
	}
	switch ev := msg.Data.(type) {
	case *EventsAPIAppRateLimited:
		{
		}
	default:
		{
			fmt.Println(ev)
			t.Fail()
		}
	}
}

func TestThatOuterCallbackEventHasInnerEvent(t *testing.T) {
	eventsAPIRawCallbackEvent := `
			{
				"token": "XXYYZZ",
				"team_id": "TXXXXXXXX",
				"api_app_id": "AXXXXXXXXX",
				"event": {
								"type": "app_mention",
								"event_ts": "1234567890.123456",
								"user": "UXXXXXXX1"
				},
				"type": "event_callback",
				"authed_users": [ "UXXXXXXX1" ],
				"event_id": "Ev08MFMKH6",
				"event_time": 1234567890
		}
	`
	msg, e := ParseEvent(json.RawMessage(eventsAPIRawCallbackEvent), OptionVerifyToken(&TokenComparator{"XXYYZZ"}))
	if e != nil {
		fmt.Println(e)
		t.Fail()
	}
	switch outerEvent := msg.Data.(type) {
	case *EventsAPICallbackEvent:
		{
			switch innerEvent := msg.InnerEvent.Data.(type) {
			case *AppMentionEvent:
				{
				}
			default:
				fmt.Println(innerEvent)
				t.Fail()
			}
		}
	default:
		{
			fmt.Println(outerEvent)
			t.Fail()
		}
	}
}

func TestBadTokenVerification(t *testing.T) {
	urlVerificationEvent := `
		{
			"token": "fake-token",
			"challenge": "aljdsflaji3jj",
			"type": "url_verification"
		}
	`
	_, e := ParseEvent(json.RawMessage(urlVerificationEvent), OptionVerifyToken(TokenComparator{"real-token"}))
	if e == nil {
		t.Fail()
	}
}

func TestNoTokenVerification(t *testing.T) {
	urlVerificationEvent := `
		{
			"token": "fake-token",
			"challenge": "aljdsflaji3jj",
			"type": "url_verification"
		}
	`
	_, e := ParseEvent(json.RawMessage(urlVerificationEvent), OptionNoVerifyToken())
	if e != nil {
		fmt.Println(e)
		t.Fail()
	}
}

func TestParseEventAPIAppMentionWithAssistantThread(t *testing.T) {
	eventsAPIRawCallbackEvent := `
		{
			"token": "XXYYZZ",
			"team_id": "TXXXXXXXX",
			"api_app_id": "AXXXXXXXXX",
			"event": {
				"type": "app_mention",
				"event_ts": "1234567890.123456",
				"user": "UXXXXXXX1",
				"text": "<@U0LAN0Z89> help me with something",
				"ts": "1515449522.000016",
				"channel": "C0LAN2Q65",
				"assistant_thread": {
					"action_token": "1234567.abcdefg"
				}
			},
			"type": "event_callback",
			"authed_users": [ "UXXXXXXX1" ],
			"event_id": "Ev08MFMKH6",
			"event_time": 1234567890
		}
	`
	msg, e := ParseEvent(json.RawMessage(eventsAPIRawCallbackEvent), OptionVerifyToken(&TokenComparator{"XXYYZZ"}))
	if e != nil {
		fmt.Println(e)
		t.Fail()
	}

	switch outerEvent := msg.Data.(type) {
	case *EventsAPICallbackEvent:
		{
			switch innerEvent := msg.InnerEvent.Data.(type) {
			case *AppMentionEvent:
				{
					if innerEvent.AssistantThread == nil {
						t.Error("Expected AssistantThread to be non-nil")
					}
					if innerEvent.AssistantThread.ActionToken != "1234567.abcdefg" {
						t.Errorf("Expected ActionToken to be '1234567.abcdefg', got %s", innerEvent.AssistantThread.ActionToken)
					}
				}
			default:
				fmt.Println(innerEvent)
				t.Fail()
			}
		}
	default:
		{
			fmt.Println(outerEvent)
			t.Fail()
		}
	}
}

func TestParseEventAPIMessageIMWithAssistantThread(t *testing.T) {
	eventsAPIRawCallbackEvent := `
		{
			"token": "XXYYZZ",
			"team_id": "TXXXXXXXX",
			"api_app_id": "AXXXXXXXXX",
			"event": {
				"type": "message",
				"channel": "D024BE91L",
				"user": "U2147483697",
				"text": "Hello, I need help with something.",
				"ts": "1355517523.000005",
				"event_ts": "1355517523.000005",
				"channel_type": "im",
				"assistant_thread": {
					"action_token": "9876543.hijklmnop"
				}
			},
			"type": "event_callback",
			"authed_users": [ "U2147483697" ],
			"event_id": "Ev08MFMKH7",
			"event_time": 1234567890
		}
	`
	msg, e := ParseEvent(json.RawMessage(eventsAPIRawCallbackEvent), OptionVerifyToken(&TokenComparator{"XXYYZZ"}))
	if e != nil {
		fmt.Println(e)
		t.Fail()
	}

	switch outerEvent := msg.Data.(type) {
	case *EventsAPICallbackEvent:
		{
			switch innerEvent := msg.InnerEvent.Data.(type) {
			case *MessageEvent:
				{
					if innerEvent.AssistantThread == nil {
						t.Error("Expected AssistantThread to be non-nil")
					}
					if innerEvent.AssistantThread.ActionToken != "9876543.hijklmnop" {
						t.Errorf("Expected ActionToken to be '9876543.hijklmnop', got %s", innerEvent.AssistantThread.ActionToken)
					}
					if innerEvent.ChannelType != "im" {
						t.Errorf("Expected ChannelType to be 'im', got %s", innerEvent.ChannelType)
					}
				}
			default:
				fmt.Println(innerEvent)
				t.Fail()
			}
		}
	default:
		{
			fmt.Println(outerEvent)
			t.Fail()
		}
	}
}

func TestParseEventAPIMessageChannelWithAssistantThread(t *testing.T) {
	eventsAPIRawCallbackEvent := `
		{
			"token": "XXYYZZ",
			"team_id": "TXXXXXXXX",
			"api_app_id": "AXXXXXXXXX",
			"event": {
				"type": "message",
				"channel": "C024BE91L",
				"user": "U2147483697",
				"text": "Hello everyone, I need help with something.",
				"ts": "1355517523.000005",
				"event_ts": "1355517523.000005",
				"channel_type": "channel",
				"assistant_thread": {
					"action_token": "abcd1234.qwerty"
				}
			},
			"type": "event_callback",
			"authed_users": [ "U2147483697" ],
			"event_id": "Ev08MFMKH8",
			"event_time": 1234567890
		}
	`
	msg, e := ParseEvent(json.RawMessage(eventsAPIRawCallbackEvent), OptionVerifyToken(&TokenComparator{"XXYYZZ"}))
	if e != nil {
		fmt.Println(e)
		t.Fail()
	}

	switch outerEvent := msg.Data.(type) {
	case *EventsAPICallbackEvent:
		{
			switch innerEvent := msg.InnerEvent.Data.(type) {
			case *MessageEvent:
				{
					if innerEvent.AssistantThread == nil {
						t.Error("Expected AssistantThread to be non-nil")
					}
					if innerEvent.AssistantThread.ActionToken != "abcd1234.qwerty" {
						t.Errorf("Expected ActionToken to be 'abcd1234.qwerty', got %s", innerEvent.AssistantThread.ActionToken)
					}
					if innerEvent.ChannelType != "channel" {
						t.Errorf("Expected ChannelType to be 'channel', got %s", innerEvent.ChannelType)
					}
				}
			default:
				fmt.Println(innerEvent)
				t.Fail()
			}
		}
	default:
		{
			fmt.Println(outerEvent)
			t.Fail()
		}
	}
}

func TestParseEventAPIMessageMPIMWithAssistantThread(t *testing.T) {
	eventsAPIRawCallbackEvent := `
		{
			"token": "XXYYZZ",
			"team_id": "TXXXXXXXX",
			"api_app_id": "AXXXXXXXXX",
			"event": {
				"type": "message",
				"channel": "G024BE91L",
				"user": "U2147483697",
				"text": "Hey team, I need some assistance.",
				"ts": "1355517523.000005",
				"event_ts": "1355517523.000005",
				"channel_type": "mpim",
				"assistant_thread": {
					"action_token": "xyz789.multiparty"
				}
			},
			"type": "event_callback",
			"authed_users": [ "U2147483697" ],
			"event_id": "Ev08MFMKH9",
			"event_time": 1234567890
		}
	`
	msg, e := ParseEvent(json.RawMessage(eventsAPIRawCallbackEvent), OptionVerifyToken(&TokenComparator{"XXYYZZ"}))
	if e != nil {
		fmt.Println(e)
		t.Fail()
	}

	switch outerEvent := msg.Data.(type) {
	case *EventsAPICallbackEvent:
		{
			switch innerEvent := msg.InnerEvent.Data.(type) {
			case *MessageEvent:
				{
					if innerEvent.AssistantThread == nil {
						t.Error("Expected AssistantThread to be non-nil")
					}
					if innerEvent.AssistantThread.ActionToken != "xyz789.multiparty" {
						t.Errorf("Expected ActionToken to be 'xyz789.multiparty', got %s", innerEvent.AssistantThread.ActionToken)
					}
					if innerEvent.ChannelType != "mpim" {
						t.Errorf("Expected ChannelType to be 'mpim', got %s", innerEvent.ChannelType)
					}
				}
			default:
				fmt.Println(innerEvent)
				t.Fail()
			}
		}
	default:
		{
			fmt.Println(outerEvent)
			t.Fail()
		}
	}
}

func TestParseEventAPIMessageGroupWithAssistantThread(t *testing.T) {
	eventsAPIRawCallbackEvent := `
		{
			"token": "XXYYZZ",
			"team_id": "TXXXXXXXX",
			"api_app_id": "AXXXXXXXXX",
			"event": {
				"type": "message",
				"channel": "G124BE91L",
				"user": "U2147483697",
				"text": "Private group message with assistant request.",
				"ts": "1355517523.000005",
				"event_ts": "1355517523.000005",
				"channel_type": "group",
				"assistant_thread": {
					"action_token": "group123.private"
				}
			},
			"type": "event_callback",
			"authed_users": [ "U2147483697" ],
			"event_id": "Ev08MFMK10",
			"event_time": 1234567890
		}
	`
	msg, e := ParseEvent(json.RawMessage(eventsAPIRawCallbackEvent), OptionVerifyToken(&TokenComparator{"XXYYZZ"}))
	if e != nil {
		fmt.Println(e)
		t.Fail()
	}

	switch outerEvent := msg.Data.(type) {
	case *EventsAPICallbackEvent:
		{
			switch innerEvent := msg.InnerEvent.Data.(type) {
			case *MessageEvent:
				{
					if innerEvent.AssistantThread == nil {
						t.Error("Expected AssistantThread to be non-nil")
					}
					if innerEvent.AssistantThread.ActionToken != "group123.private" {
						t.Errorf("Expected ActionToken to be 'group123.private', got %s", innerEvent.AssistantThread.ActionToken)
					}
					if innerEvent.ChannelType != "group" {
						t.Errorf("Expected ChannelType to be 'group', got %s", innerEvent.ChannelType)
					}
				}
			default:
				fmt.Println(innerEvent)
				t.Fail()
			}
		}
	default:
		{
			fmt.Println(outerEvent)
			t.Fail()
		}
	}
}
