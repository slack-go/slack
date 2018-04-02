package slack

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/ACollectionOfAtoms/slack"
)

var c = slack.New("my-token")

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
	msg := c.ParseEventsAPIEvent(json.RawMessage(eventsAPIRawCallbackEvent))
	switch ev := msg.Data.(type) {
	case *slack.EventsAPICallbackEvent:
		{
			fmt.Println(ev)
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
	msg := c.ParseEventsAPIEvent(json.RawMessage(urlVerificationEvent))
	switch ev := msg.Data.(type) {
	case *slack.EventsAPIURLVerificationEvent:
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
	msg := c.ParseEventsAPIEvent(json.RawMessage(eventsAPIRawCallbackEvent))
	switch outterEvent := msg.Data.(type) {
	case *slack.EventsAPICallbackEvent:
		{
			switch innerEvent := msg.InnerEvent.Data.(type) {
			case *slack.AppMentionEvent:
				{
				}
			default:
				fmt.Println(innerEvent)
				t.Fail()
			}
		}
	default:
		{
			fmt.Println(outterEvent)
			t.Fail()
		}
	}
}
