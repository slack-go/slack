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
