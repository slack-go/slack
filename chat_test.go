package slack

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"
	"testing"
)

func postMessageInvalidChannelHandler(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	response, _ := json.Marshal(chatResponseFull{
		SlackResponse: SlackResponse{Ok: false, Error: "channel_not_found"},
	})
	rw.Write(response)
}

func TestPostMessageInvalidChannel(t *testing.T) {
	http.DefaultServeMux = new(http.ServeMux)
	http.HandleFunc("/chat.postMessage", postMessageInvalidChannelHandler)
	once.Do(startServer)
	api := New("testing-token", OptionAPIURL("http://"+serverAddr+"/"))
	_, _, err := api.PostMessage("CXXXXXXXX", MsgOptionText("hello", false))
	if err == nil {
		t.Errorf("Expected error: channel_not_found; instead succeeded")
		return
	}

	if err.Error() != "channel_not_found" {
		t.Errorf("Expected error: channel_not_found; received: %s", err)
		return
	}
}

func TestGetPermalink(t *testing.T) {
	channel := "C1H9RESGA"
	timeStamp := "p135854651500008"

	http.HandleFunc("/chat.getPermalink", func(rw http.ResponseWriter, r *http.Request) {

		if got, want := r.Header.Get("Content-Type"), "application/x-www-form-urlencoded"; got != want {
			t.Errorf("request uses unexpected content type: got %s, want %s", got, want)
		}

		if got, want := r.URL.Query().Get("channel"), channel; got != want {
			t.Errorf("request contains unexpected channel: got %s, want %s", got, want)
		}

		if got, want := r.URL.Query().Get("message_ts"), timeStamp; got != want {
			t.Errorf("request contains unexpected message timestamp: got %s, want %s", got, want)
		}

		rw.Header().Set("Content-Type", "application/json")
		response := []byte("{\"ok\": true, \"channel\": \"" + channel + "\", \"permalink\": \"https://ghostbusters.slack.com/archives/" + channel + "/" + timeStamp + "\"}")
		rw.Write(response)
	})

	once.Do(startServer)
	api := New("testing-token", OptionAPIURL("http://"+serverAddr+"/"))
	pp := PermalinkParameters{Channel: channel, Ts: timeStamp}
	pl, err := api.GetPermalink(&pp)

	if got, want := pl, "https://ghostbusters.slack.com/archives/C1H9RESGA/p135854651500008"; got != want {
		t.Errorf("unexpected permalink: got %s, want %s", got, want)
	}

	if err != nil {
		t.Errorf("unexpected error returned: %v", err)
	}
}

func TestPostMessage(t *testing.T) {
	type messageTest struct {
		endpoint string
		opt      []MsgOption
		expected url.Values
	}

	blocks := []Block{NewContextBlock("context", NewTextBlockObject(PlainTextType, "hello", false, false))}
	blockStr := `[{"type":"context","block_id":"context","elements":[{"type":"plain_text","text":"hello"}]}]`

	tests := map[string]messageTest{
		"OnlyBasicProperties": {
			endpoint: "/chat.postMessage",
			opt:      []MsgOption{},
			expected: url.Values{
				"channel": []string{"CXXX"},
				"token":   []string{"testing-token"},
			},
		},
		"Blocks": {
			endpoint: "/chat.postMessage",
			opt: []MsgOption{
				MsgOptionBlocks(blocks...),
				MsgOptionText("text", false),
			},
			expected: url.Values{
				"blocks":  []string{blockStr},
				"channel": []string{"CXXX"},
				"text":    []string{"text"},
				"token":   []string{"testing-token"},
			},
		},
		"Attachment": {
			endpoint: "/chat.postMessage",
			opt: []MsgOption{
				MsgOptionAttachments(
					Attachment{
						Blocks: Blocks{BlockSet: blocks},
					}),
			},
			expected: url.Values{
				"attachments": []string{`[{"blocks":` + blockStr + `}]`},
				"channel":     []string{"CXXX"},
				"token":       []string{"testing-token"},
			},
		},
		"Metadata": {
			endpoint: "/chat.postMessage",
			opt: []MsgOption{
				MsgOptionMetadata(
					SlackMetadata{
						EventType: "testing-event",
						EventPayload: map[string]interface{}{
							"id":   13,
							"name": "testing-name",
						},
					}),
			},
			expected: url.Values{
				"metadata": []string{`{"event_type":"testing-event","event_payload":{"id":13,"name":"testing-name"}}`},
				"channel":  []string{"CXXX"},
				"token":    []string{"testing-token"},
			},
		},
		"Unfurl": {
			endpoint: "/chat.unfurl",
			opt: []MsgOption{
				MsgOptionUnfurl("123", map[string]Attachment{"something": {Text: "attachment-test"}}),
			},
			expected: url.Values{
				"channel": []string{"CXXX"},
				"token":   []string{"testing-token"},
				"ts":      []string{"123"},
				"unfurls": []string{`{"something":{"text":"attachment-test","blocks":null}}`},
			},
		},
		"UnfurlAuthURL": {
			endpoint: "/chat.unfurl",
			opt: []MsgOption{
				MsgOptionUnfurlAuthURL("123", "https://auth-url.com"),
			},
			expected: url.Values{
				"channel":       []string{"CXXX"},
				"token":         []string{"testing-token"},
				"ts":            []string{"123"},
				"user_auth_url": []string{"https://auth-url.com"},
			},
		},
		"UnfurlAuthRequired": {
			endpoint: "/chat.unfurl",
			opt: []MsgOption{
				MsgOptionUnfurlAuthRequired("123"),
			},
			expected: url.Values{
				"channel":            []string{"CXXX"},
				"token":              []string{"testing-token"},
				"ts":                 []string{"123"},
				"user_auth_required": []string{"true"},
			},
		},
		"UnfurlAuthMessage": {
			endpoint: "/chat.unfurl",
			opt: []MsgOption{
				MsgOptionUnfurlAuthMessage("123", "Please!"),
			},
			expected: url.Values{
				"channel":           []string{"CXXX"},
				"token":             []string{"testing-token"},
				"ts":                []string{"123"},
				"user_auth_message": []string{"Please!"},
			},
		},
	}

	once.Do(startServer)
	api := New(validToken, OptionAPIURL("http://"+serverAddr+"/"))

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			http.DefaultServeMux = new(http.ServeMux)
			http.HandleFunc(test.endpoint, func(rw http.ResponseWriter, r *http.Request) {
				body, err := ioutil.ReadAll(r.Body)
				if err != nil {
					t.Errorf("unexpected error: %v", err)
					return
				}
				actual, err := url.ParseQuery(string(body))
				if err != nil {
					t.Errorf("unexpected error: %v", err)
					return
				}
				if !reflect.DeepEqual(actual, test.expected) {
					t.Errorf("\nexpected: %s\n  actual: %s", test.expected, actual)
					return
				}
			})

			_, _, _ = api.PostMessage("CXXX", test.opt...)
		})
	}
}

func TestPostMessageWithBlocksWhenMsgOptionResponseURLApplied(t *testing.T) {
	expectedBlocks := []Block{NewContextBlock("context", NewTextBlockObject(PlainTextType, "hello", false, false))}

	http.DefaultServeMux = new(http.ServeMux)
	http.HandleFunc("/response-url", func(rw http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
			return
		}
		var msg Msg
		if err := json.Unmarshal(body, &msg); err != nil {
			t.Errorf("unexpected error: %v", err)
			return
		}
		actualBlocks := msg.Blocks.BlockSet
		if !reflect.DeepEqual(expectedBlocks, actualBlocks) {
			t.Errorf("expected: %#v, got: %#v", expectedBlocks, actualBlocks)
			return
		}
	})

	once.Do(startServer)
	api := New(validToken, OptionAPIURL("http://"+serverAddr+"/"))

	responseURL := api.endpoint + "response-url"

	_, _, _ = api.PostMessage("CXXX", MsgOptionBlocks(expectedBlocks...), MsgOptionText("text", false), MsgOptionResponseURL(responseURL, ResponseTypeInChannel))
}

func TestPostMessageWhenMsgOptionReplaceOriginalApplied(t *testing.T) {
	http.DefaultServeMux = new(http.ServeMux)
	http.HandleFunc("/response-url", func(rw http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
			return
		}
		var msg Msg
		if err := json.Unmarshal(body, &msg); err != nil {
			t.Errorf("unexpected error: %v", err)
			return
		}
		if msg.ReplaceOriginal != true {
			t.Errorf("expected: true, got: %v", msg.ReplaceOriginal)
			return
		}
	})

	once.Do(startServer)
	api := New(validToken, OptionAPIURL("http://"+serverAddr+"/"))

	responseURL := api.endpoint + "response-url"

	_, _, _ = api.PostMessage("CXXX", MsgOptionText("text", false), MsgOptionReplaceOriginal(responseURL))
}

func TestPostMessageWhenMsgOptionDeleteOriginalApplied(t *testing.T) {
	http.DefaultServeMux = new(http.ServeMux)
	http.HandleFunc("/response-url", func(rw http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
			return
		}
		var msg Msg
		if err := json.Unmarshal(body, &msg); err != nil {
			t.Errorf("unexpected error: %v", err)
			return
		}
		if msg.DeleteOriginal != true {
			t.Errorf("expected: true, got: %v", msg.DeleteOriginal)
			return
		}
	})

	once.Do(startServer)
	api := New(validToken, OptionAPIURL("http://"+serverAddr+"/"))

	responseURL := api.endpoint + "response-url"

	_, _, _ = api.PostMessage("CXXX", MsgOptionDeleteOriginal(responseURL))
}
