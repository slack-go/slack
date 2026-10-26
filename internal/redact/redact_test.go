package redact

import "testing"

func TestURLOriginOmitsCredentials(t *testing.T) {
	for _, tt := range []struct {
		url  string
		want string
	}{
		{"wss://wss.slack.com/link/?ticket=secret", "wss://wss.slack.com"},
		{"wss://user:secret@wss.slack.com/private-secret?ticket=secret#secret", "wss://wss.slack.com"},
		{"ws://[::1]:8080/link/?ticket=secret", "ws://[::1]:8080"},
		{"wss://wss.slack.com/%secret", "[invalid URL]"},
		{"/relative?ticket=secret", "[invalid URL]"},
		{"secret", "[invalid URL]"},
	} {
		if got := URLOrigin(tt.url); got != tt.want {
			t.Errorf("URLOrigin(%q) = %q, want %q", tt.url, got, tt.want)
		}
	}
}
