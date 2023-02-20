package socketmode

import (
	"bytes"
	"encoding/json"
	"errors"
	"reflect"
	"testing"

	"github.com/slack-go/slack/slackevents"
)

const (
	EventDisconnect = `{
  "type": "disconnect",
  "reason": "warning",
  "debug_info": {
    "host": "applink-7fc4fdbb64-4x5xq"
  }
}
`
	EventHello = `{
  "type": "hello",
  "num_connections": 4,
  "debug_info": {
    "host": "applink-7fc4fdbb64-4x5xq",
    "build_number": 10,
    "approximate_connection_time": 18060
  },
  "connection_info": {
    "app_id": "A01K58AR4RF"
  }
}
`

	EventAppMention = `{
  "envelope_id": "c67a03d0-4094-4744-90ca-d286e00a3ab1",
  "payload": {
    "token": "redacted",
    "team_id": "redacted",
    "api_app_id": "redacted",
    "event": {
      "client_msg_id": "c714568f-67df-42d7-a343-0a8e4d9c6030",
      "type": "app_mention",
      "text": "\u003c@U01JKSB8T7Y\u003e test",
      "user": "redacted",
      "ts": "1610927831.000200",
      "team": "redacted",
      "blocks": [
        {
          "type": "rich_text",
          "block_id": "2Le",
          "elements": [
            {
              "type": "rich_text_section",
              "elements": [
                {
                  "type": "user",
                  "user_id": "redacted"
                },
                {
                  "type": "text",
                  "text": " test39"
                }
              ]
            }
          ]
        }
      ],
      "channel": "redacted",
      "event_ts": "1610927831.000200"
    },
    "type": "event_callback",
    "event_id": "Ev01JZ2T7S3U",
    "event_time": 1610927831,
    "authorizations": [
      {
        "enterprise_id": null,
        "team_id": "redacted",
        "user_id": "redacted",
        "is_bot": true,
        "is_enterprise_install": false
      }
    ],
    "is_ext_shared_channel": false,
    "event_context": "1-app_mention-redacted-redacted"
  },
  "type": "events_api",
  "accepts_response_payload": false,
  "retry_attempt": 0,
  "retry_reason": ""
}`
)

func TestEventParsing(t *testing.T) {
	testParsing(t,
		EventHello,
		&Event{
			Type: EventTypeHello,
			Request: &Request{
				Type:           RequestTypeHello,
				NumConnections: 4,
				DebugInfo: DebugInfo{
					Host:                      "applink-7fc4fdbb64-4x5xq",
					BuildNumber:               10,
					ApproximateConnectionTime: 18060,
				},
				ConnectionInfo: ConnectionInfo{
					AppID: "A01K58AR4RF",
				},
			},
		})

	testParsing(t,
		EventDisconnect,
		&Event{
			Type: EventTypeDisconnect,
			Request: &Request{
				Type:   RequestTypeDisconnect,
				Reason: "warning",
				DebugInfo: DebugInfo{
					Host: "applink-7fc4fdbb64-4x5xq",
				},
			},
		})

	rawAppMention := json.RawMessage(`{
      "client_msg_id": "c714568f-67df-42d7-a343-0a8e4d9c6030",
      "type": "app_mention",
      "text": "\u003c@U01JKSB8T7Y\u003e test",
      "user": "redacted",
      "ts": "1610927831.000200",
      "team": "redacted",
      "blocks": [
        {
          "type": "rich_text",
          "block_id": "2Le",
          "elements": [
            {
              "type": "rich_text_section",
              "elements": [
                {
                  "type": "user",
                  "user_id": "redacted"
                },
                {
                  "type": "text",
                  "text": " test39"
                }
              ]
            }
          ]
        }
      ],
      "channel": "redacted",
      "event_ts": "1610927831.000200"
    }`)

	rawAppMentionReqPayload := json.RawMessage(`{
    "token": "redacted",
    "team_id": "redacted",
    "api_app_id": "redacted",
    "event": {
      "client_msg_id": "c714568f-67df-42d7-a343-0a8e4d9c6030",
      "type": "app_mention",
      "text": "\u003c@U01JKSB8T7Y\u003e test",
      "user": "redacted",
      "ts": "1610927831.000200",
      "team": "redacted",
      "blocks": [
        {
          "type": "rich_text",
          "block_id": "2Le",
          "elements": [
            {
              "type": "rich_text_section",
              "elements": [
                {
                  "type": "user",
                  "user_id": "redacted"
                },
                {
                  "type": "text",
                  "text": " test39"
                }
              ]
            }
          ]
        }
      ],
      "channel": "redacted",
      "event_ts": "1610927831.000200"
    },
    "type": "event_callback",
    "event_id": "Ev01JZ2T7S3U",
    "event_time": 1610927831,
    "authorizations": [
      {
        "enterprise_id": null,
        "team_id": "redacted",
        "user_id": "redacted",
        "is_bot": true,
        "is_enterprise_install": false
      }
    ],
    "is_ext_shared_channel": false,
    "event_context": "1-app_mention-redacted-redacted"
  }`)
	testParsing(t,
		EventAppMention,
		&Event{
			Type: EventTypeEventsAPI,
			Data: slackevents.EventsAPIEvent{
				Token:    "redacted",
				TeamID:   "redacted",
				Type:     "event_callback",
				APIAppID: "redacted",
				Data: &slackevents.EventsAPICallbackEvent{
					Type:         "event_callback",
					Token:        "redacted",
					TeamID:       "redacted",
					APIAppID:     "redacted",
					InnerEvent:   &rawAppMention,
					AuthedUsers:  nil,
					AuthedTeams:  nil,
					EventID:      "Ev01JZ2T7S3U",
					EventTime:    1610927831,
					EventContext: "1-app_mention-redacted-redacted",
				},
				InnerEvent: slackevents.EventsAPIInnerEvent{
					Type: string(slackevents.AppMention),
					Data: &slackevents.AppMentionEvent{
						Type:            string(slackevents.AppMention),
						User:            "redacted",
						Text:            "<@U01JKSB8T7Y> test",
						TimeStamp:       "1610927831.000200",
						ThreadTimeStamp: "",
						Channel:         "redacted",
						EventTimeStamp:  "1610927831.000200",
					},
				},
			},
			Request: &Request{
				Type:                   RequestTypeEventsAPI,
				EnvelopeID:             "c67a03d0-4094-4744-90ca-d286e00a3ab1",
				Payload:                rawAppMentionReqPayload,
				AcceptsResponsePayload: false,
				RetryAttempt:           0,
				RetryReason:            "",
			},
		})
}

func testParsing(t *testing.T, raw string, want interface{}) {
	t.Helper()

	got, err := parse(raw)
	if err != nil {
		t.Fatalf("unexpected error parsing event %q: %v", raw, err)
	}

	if !reflect.DeepEqual(want, got) {
		t.Fatalf("unexpected parse result: want %s, got %s", dump(t, want), dump(t, got))
	}
}

func dump(t *testing.T, data interface{}) string {
	t.Helper()

	var buf bytes.Buffer

	e := json.NewEncoder(&buf)
	e.SetIndent("", "  ")

	if err := e.Encode(data); err != nil {
		t.Fatalf("encoding data to json: %v", err)
	}

	return buf.String()
}

func parse(raw string) (*Event, error) {
	c := &Client{}

	evt, err := c.parseEvent(json.RawMessage([]byte(raw)))

	if evt == nil {
		return nil, errors.New("failed handling raw event: event was empty")
	}

	return evt, err
}
