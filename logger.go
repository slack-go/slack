package slack

import (
	"fmt"
)

// logProvider is a logger interface compatible with both stdlib and some
// 3rd party loggers such as logrus.
type logProvider interface {
	Output(int, string) error
}

// logInternal represents the internal logging api we use.
type logInternal interface {
	Print(...interface{})
	Printf(string, ...interface{})
	Println(...interface{})
	Output(int, string) error
}

type debug interface {
	Debug() bool

	// Debugf print a formatted debug line.
	Debugf(format string, v ...interface{})
	// Debugln print a debug line.
	Debugln(v ...interface{})
}

// ilogger implements the additional methods used by our internal logging.
type ilogger struct {
	logProvider
}

// Println replicates the behaviour of the standard logger.
func (t ilogger) Println(v ...interface{}) {
	t.Output(2, fmt.Sprintln(v...))
}

// Printf replicates the behaviour of the standard logger.
func (t ilogger) Printf(format string, v ...interface{}) {
	t.Output(2, fmt.Sprintf(format, v...))
}

// Print replicates the behaviour of the standard logger.
func (t ilogger) Print(v ...interface{}) {
	t.Output(2, fmt.Sprint(v...))
}

type discard struct{}

func (t discard) Debug() bool {
	return false
}

// Debugf print a formatted debug line.
func (t discard) Debugf(format string, v ...interface{}) {}

// Debugln print a debug line.
func (t discard) Debugln(v ...interface{}) {}
