package slack

import (
	"bytes"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLogging(t *testing.T) {
	buf := bytes.NewBufferString("")
	logger := ilogger{logProvider: log.New(buf, "", 0|log.Lshortfile)}
	logger.Println("test line 123")
	assert.Equal(t, buf.String(), "logger_test.go:14: test line 123\n")
	buf.Truncate(0)
	logger.Print("test line 123")
	assert.Equal(t, buf.String(), "logger_test.go:17: test line 123\n")
	buf.Truncate(0)
	logger.Printf("test line 123\n")
	assert.Equal(t, buf.String(), "logger_test.go:20: test line 123\n")
	buf.Truncate(0)
	logger.Output(1, "test line 123\n")
	assert.Equal(t, buf.String(), "logger_test.go:23: test line 123\n")
	buf.Truncate(0)
}
