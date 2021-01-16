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
	testParsing(t, EventHello, &ClientEvent{Type: "hello", Data: &slack.HelloEvent{}})
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

func parse(raw string) (*ClientEvent, error) {
	c := &Client{
		IncomingEvents: make(chan ClientEvent, 1),
	}

	tpe := c.handleWebSocketMessage(json.RawMessage([]byte(raw)))

	if tpe == "" {
		return nil, errors.New("failed handling raw event: type was empty")
	}

	select {
	case evt := <-c.IncomingEvents:
		return &evt, nil
	default:
		return nil, errors.New("no expected event emitted")
	}
}
