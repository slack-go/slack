package slack

import (
	"fmt"
	"net/http"
	"net/url"
	"reflect"
	"strings"
	"testing"
)

func TestSlash_ServeHTTP(t *testing.T) {
	once.Do(startServer)
	serverURL := fmt.Sprintf("http://%s/slash", serverAddr)

	tests := []struct {
		body           url.Values
		wantParams     SlashCommand
		wantStatusCode int
	}{
		{
			body: url.Values{
				"command":         []string{"/command"},
				"team_domain":     []string{"team"},
				"enterprise_id":   []string{"E0001"},
				"enterprise_name": []string{"Globular%20Construct%20Inc"},
				"channel_id":      []string{"C1234ABCD"},
				"text":            []string{"text"},
				"team_id":         []string{"T1234ABCD"},
				"user_id":         []string{"U1234ABCD"},
				"user_name":       []string{"username"},
				"response_url":    []string{"https://hooks.slack.com/commands/XXXXXXXX/00000000000/YYYYYYYYYYYYYY"},
				"token":           []string{"valid"},
				"channel_name":    []string{"channel"},
				"trigger_id":      []string{"0000000000.1111111111.222222222222aaaaaaaaaaaaaa"},
			},
			wantParams: SlashCommand{
				Command:        "/command",
				TeamDomain:     "team",
				EnterpriseID:   "E0001",
				EnterpriseName: "Globular%20Construct%20Inc",
				ChannelID:      "C1234ABCD",
				Text:           "text",
				TeamID:         "T1234ABCD",
				UserID:         "U1234ABCD",
				UserName:       "username",
				ResponseURL:    "https://hooks.slack.com/commands/XXXXXXXX/00000000000/YYYYYYYYYYYYYY",
				Token:          "valid",
				ChannelName:    "channel",
				TriggerID:      "0000000000.1111111111.222222222222aaaaaaaaaaaaaa",
			},
			wantStatusCode: http.StatusOK,
		},
		{
			body: url.Values{
				"token": []string{"invalid"},
			},
			wantParams: SlashCommand{
				Token: "invalid",
			},
			wantStatusCode: http.StatusUnauthorized,
		},
	}

	var slashCommand SlashCommand
	client := &http.Client{}
	http.HandleFunc("/slash", func(w http.ResponseWriter, r *http.Request) {
		var err error
		slashCommand, err = SlashCommandParse(r)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
		acceptableTokens := []string{"valid", "valid2"}
		if !slashCommand.ValidateToken(acceptableTokens...) {
			w.WriteHeader(http.StatusUnauthorized)
		}
	})

	for i, test := range tests {
		req, err := http.NewRequest(http.MethodPost, serverURL, strings.NewReader(test.body.Encode()))
		if err != nil {
			t.Fatalf("%d: Unexpected error: %s", i, err)
		}
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("%d: Unexpected error: %s", i, err)
		}

		if resp.StatusCode != test.wantStatusCode {
			t.Errorf("%d: Got status code %d, want %d", i, resp.StatusCode, test.wantStatusCode)
		}
		if !reflect.DeepEqual(slashCommand, test.wantParams) {
			t.Errorf("%d: Got params %#v, want %#v", i, slashCommand, test.wantParams)
		}
		resp.Body.Close()
	}
}
