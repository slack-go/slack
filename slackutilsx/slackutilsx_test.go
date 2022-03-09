package slackutilsx

import (
	"testing"
)

func TestDetectChannelType(t *testing.T) {
	test := func(channelID string, expected ChannelType) {
		if computed := DetectChannelType(channelID); computed != expected {
			t.Errorf("expected channelID %s to have type %s, got: %s", channelID, expected, computed)
		}
	}

	test("G11111111", CTypeGroup)
	test("D11111111", CTypeDM)
	test("C11111111", CTypeChannel)
	test("", CTypeUnknown)
	test("X11111111", CTypeUnknown)
}

func TestEscapeMessage(t *testing.T) {
	test := func(message string, expected string) {
		if computed := EscapeMessage(message); computed != expected {
			t.Errorf("expected message %s to be converted to %s, got: %s", message, expected, computed)
		}
	}
	test("A & B", "A &amp; B")
	test("A < B", "A &lt; B")
	test("A > B", "A &gt; B")
}

func BenchmarkEscapeMessage(b *testing.B) {
	for i := 0; i < b.N; i++ {
		EscapeMessage("A & B")
	}
}
