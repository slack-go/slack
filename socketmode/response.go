package socketmode

type Response struct {
	EnvelopeID string      `json:"envelope_id"`
	Payload    interface{} `json:"payload,omitempty"`

	// rawJSON holds the pre-marshaled JSON bytes when set by SendCtx.
	// This avoids double-marshaling: once for the size check and once for
	// the WebSocket write.
	rawJSON []byte `json:"-"`
}
