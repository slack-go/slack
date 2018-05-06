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
