package socketmode

import "encoding/json"

type Response struct {
	EnvelopeID string          `json:"envelope_id"`
	Payload    json.RawMessage `json:"payload,omitempty"`
}
