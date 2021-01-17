package socketmode

// ClientEvent is the event sent to the consumer of Client
type ClientEvent struct {
	Type EventType
	Data interface{}

	// Request is the json-decoded raw WebSocket message that is received via the Slack Socket Mode
	// WebSocket connection.
	Request *Request
}

type ErrorWriteFailed struct {
	Cause    error
	Response *Response
}
