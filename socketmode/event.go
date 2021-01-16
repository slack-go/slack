package socketmode

type ErrorWriteFailed struct {
	Cause    error
	Response *Response
}
