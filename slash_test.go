package slack

import (
	"fmt"
	"net/http"
	"net/url"
	"reflect"
	"strings"
	"testing"
)

type slashHandler struct {
	gotParams *Slash
}

func newSlashHandler() *slashHandler {
	return &slashHandler{
		gotParams: &Slash{},
	}
}
func (sh *slashHandler) handler(s *Slash) (*PostMessageParameters, error) {
	sh.gotParams = s
	response := &PostMessageParameters{Text: "success"}
	return response, nil
}

func TestSlash_ServeHTTP(t *testing.T) {
	once.Do(startServer)
	serverURL := fmt.Sprintf("http://%s/slash", serverAddr)

	tests := []struct {
		body           url.Values
		wantParams     *Slash
		wantStatusCode int
	}{
		{
			body: url.Values{
				"command":      []string{"/command"},
				"team_domain":  []string{"team"},
				"channel_id":   []string{"C1234ABCD"},
				"text":         []string{"text"},
				"team_id":      []string{"T1234ABCD"},
				"user_id":      []string{"U1234ABCD"},
				"user_name":    []string{"username"},
				"response_url": []string{"https://hooks.slack.com/commands/XXXXXXXX/00000000000/YYYYYYYYYYYYYY"},
				"token":        []string{"valid"},
				"channel_name": []string{"channel"},
				"trigger_id":   []string{"0000000000.1111111111.222222222222aaaaaaaaaaaaaa"},
			},
			wantParams: &Slash{
				Command:     "/command",
				TeamDomain:  "team",
				ChannelID:   "C1234ABCD",
				Text:        "text",
				TeamID:      "T1234ABCD",
				UserID:      "U1234ABCD",
				UserName:    "username",
				ResponseURL: "https://hooks.slack.com/commands/XXXXXXXX/00000000000/YYYYYYYYYYYYYY",
				Token:       "valid",
				ChannelName: "channel",
				TriggerID:   "0000000000.1111111111.222222222222aaaaaaaaaaaaaa",
			},
			wantStatusCode: http.StatusOK,
		},
		{
			body: url.Values{
				"token": []string{"invalid"},
			},
			wantParams:     &Slash{},
			wantStatusCode: http.StatusUnauthorized,
		},
	}

	client := &http.Client{}
	h := newSlashHandler()
	verificationToken := "valid"
	http.HandleFunc("/slash", SlashHandler(verificationToken, h.handler))

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
		if !reflect.DeepEqual(h.gotParams, test.wantParams) {
			t.Errorf("%d: Got params %#v, want %#v", i, h.gotParams, test.wantParams)
		}
		resp.Body.Close()
		h = newSlashHandler()
	}
}
