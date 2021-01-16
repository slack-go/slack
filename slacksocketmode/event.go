package slacksocketmode

type ErrorWriteFailed struct {
	Cause    error
	Response *Response
}
