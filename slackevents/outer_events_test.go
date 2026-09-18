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
				"event_time": 1234567890,
				"is_ext_shared_channel": true
		}
	`)
	var cb EventsAPICallbackEvent
	err := json.Unmarshal(rawE, &cb)
	if err != nil {
		t.Error(err)
	}
	if !cb.IsExtSharedChannel {
		t.Errorf("expected IsExtSharedChannel to be true, got false")
	}
}

func TestCallBackEventAuthorizations(t *testing.T) {
	rawE := []byte(`
			{
				"token": "XXYYZZ",
				"team_id": "TXXXXXXXX",
				"api_app_id": "AXXXXXXXXX",
				"event": {
								"type": "app_context_changed",
								"context": {}
				},
				"type": "event_callback",
				"authorizations": [
					{
						"enterprise_id": null,
						"team_id": "TXXXXXXXX",
						"user_id": "UXXXXXXX1",
						"is_bot": true,
						"is_enterprise_install": false
					}
				],
				"event_id": "Ev08MFMKH6",
				"event_time": 1234567890
		}
	`)
	var cb EventsAPICallbackEvent
	err := json.Unmarshal(rawE, &cb)
	if err != nil {
		t.Fatal(err)
	}
	if len(cb.Authorizations) != 1 {
		t.Fatalf("expected 1 authorization, got %d", len(cb.Authorizations))
	}
	if cb.Authorizations[0].UserID != "UXXXXXXX1" {
		t.Errorf("expected UserID UXXXXXXX1, got %q", cb.Authorizations[0].UserID)
	}
	if cb.Authorizations[0].TeamID != "TXXXXXXXX" {
		t.Errorf("expected TeamID TXXXXXXXX, got %q", cb.Authorizations[0].TeamID)
	}
	if !cb.Authorizations[0].IsBot {
		t.Errorf("expected IsBot to be true, got false")
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
