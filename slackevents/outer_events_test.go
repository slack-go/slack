package slackevents

import (
	"encoding/json"
	"testing"
)

func TestURLVerificationEvent(t *testing.T) {
	rawE := []byte(`
			{
				"token": "Jhj5dZrVaK7ZwHHjRyZWjbDl",
				"challenge": "3eZbrw1aBm2rZgRNFdxV2595E9CY3gmdALWMmHkvFXO7tYXAYM8P",
				"type": "url_verification"
		}
	`)
	err := json.Unmarshal(rawE, &EventsAPIURLVerificationEvent{})
	if err != nil {
		t.Error(err)
	}
}

func TestCallBackEvent(t *testing.T) {
	rawE := []byte(`
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
	`)
	err := json.Unmarshal(rawE, &EventsAPICallbackEvent{})
	if err != nil {
		t.Error(err)
	}
}

func TestAppRateLimitedEvent(t *testing.T) {
	rawE := []byte(`
			{
				"token": "Jhj5dZrVaK7ZwHHjRyZWjbDl",
				"type": "app_rate_limited",
				"team_id": "T123456",
				"minute_rate_limited": 1518467820,
				"api_app_id": "A123456"
		}
	`)
	err := json.Unmarshal(rawE, &EventsAPIAppRateLimited{})
	if err != nil {
		t.Error(err)
	}
}
