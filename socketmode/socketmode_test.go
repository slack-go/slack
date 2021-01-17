package socketmode

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/pkg/errors"
	"github.com/slack-go/slack"
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
)

func TestEventParsing(t *testing.T) {
	testParsing(t, EventHello, &Event{Type: "hello", Data: &slack.HelloEvent{}})
}

func testParsing(t *testing.T, raw string, want interface{}) {
	t.Helper()

	got, err := parse(raw)
	if err != nil {
		t.Fatalf("unexpected error parsing event %q: %v", raw, err)
	}

	if !reflect.DeepEqual(want, got) {
		t.Fatalf("unexpected parse result: want %v, got %v", want, got)
	}
}

func parse(raw string) (*Event, error) {
	c := &Client{}

	evt, err := c.parseEvent(json.RawMessage([]byte(raw)))

	if evt == nil {
		return nil, errors.New("failed handling raw event: event was empty")
	}

	return evt, err
}
