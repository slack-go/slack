package slack

import (
	"fmt"
	"testing"

	"github.com/ACollectionOfAtoms/slack"
)

func TestParseEvent(t *testing.T) {
	body := `
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
	c := slack.New("token")
	e, err := c.ParseEvent(body)
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}
	if e.Type != "event_callback" {
		t.Fail()
	}
	if e.Event.Type != "app_mention" {
		t.Fail()
	}
	fmt.Println(e)
}
