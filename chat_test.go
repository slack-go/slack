package slack

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/url"
	"reflect"
	"regexp"
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
	blockStr := `[{"type":"context","block_id":"context","elements":[{"type":"plain_text","text":"hello","emoji":false}]}]`

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
		"UnfurlAuthBlocks": {
			endpoint: "/chat.unfurl",
			opt: []MsgOption{
				MsgOptionUnfurlAuthBlocks("123", NewSectionBlock(NewTextBlockObject(MarkdownType, "*Authenticate* to view", false, false), nil, nil)),
			},
			expected: url.Values{
				"channel":          []string{"CXXX"},
				"token":            []string{"testing-token"},
				"ts":               []string{"123"},
				"user_auth_blocks": []string{`[{"type":"section","text":{"type":"mrkdwn","text":"*Authenticate* to view"}}]`},
			},
		},
		"UnfurlByID": {
			endpoint: "/chat.unfurl",
			opt: []MsgOption{
				MsgOptionUnfurlByID("Uxxxxxxx-909b5454-75f8-4ac4-b325-1b40e230bbd8", "composer", map[string]Attachment{"https://example.com": {Text: "Preview"}}),
			},
			expected: url.Values{
				"token":     []string{"testing-token"},
				"unfurl_id": []string{"Uxxxxxxx-909b5454-75f8-4ac4-b325-1b40e230bbd8"},
				"source":    []string{"composer"},
				"unfurls":   []string{`{"https://example.com":{"text":"Preview","blocks":null}}`},
			},
		},
		"UnfurlByIDWithNilUnfurls": {
			endpoint: "/chat.unfurl",
			opt: []MsgOption{
				MsgOptionUnfurlByID("uf-123", "composer", nil),
			},
			expected: url.Values{
				"token":     []string{"testing-token"},
				"unfurl_id": []string{"uf-123"},
				"source":    []string{"composer"},
				"unfurls":   []string{`{}`},
			},
		},
		"UnfurlWorkObjectMetadataOnly": {
			endpoint: "/chat.unfurl",
			opt: []MsgOption{
				MsgOptionUnfurlWorkObject("123", nil, WorkObjectMetadata{
					Entities: []WorkObjectEntity{{
						URL:           "https://example.com/doc/1",
						ExternalRef:   WorkObjectExternalRef{ID: "1"},
						EntityType:    EntityTypeFile,
						EntityPayload: map[string]interface{}{"title": "Doc"},
					}},
				}),
			},
			expected: url.Values{
				"channel":  []string{"CXXX"},
				"token":    []string{"testing-token"},
				"ts":       []string{"123"},
				"metadata": []string{`{"entities":[{"url":"https://example.com/doc/1","external_ref":{"id":"1"},"entity_type":"slack#/entities/file","entity_payload":{"title":"Doc"}}]}`},
			},
		},
		"LinkNames true": {
			endpoint: "/chat.postMessage",
			opt: []MsgOption{
				MsgOptionLinkNames(true),
			},
			expected: url.Values{
				"channel":    []string{"CXXX"},
				"token":      []string{"testing-token"},
				"link_names": []string{"true"},
			},
		},
		"LinkNames false": {
			endpoint: "/chat.postMessage",
			opt: []MsgOption{
				MsgOptionLinkNames(false),
			},
			expected: url.Values{
				"channel":    []string{"CXXX"},
				"token":      []string{"testing-token"},
				"link_names": []string{"false"},
			},
		},
		"MetadataViaPostMessageParameters": {
			endpoint: "/chat.postMessage",
			opt: []MsgOption{
				MsgOptionPostMessageParameters(PostMessageParameters{
					MetaData: SlackMetadata{
						EventType: "testing-event",
						EventPayload: map[string]interface{}{
							"id":   13,
							"name": "testing-name",
						},
					},
				}),
			},
			expected: url.Values{
				"metadata":     []string{`{"event_type":"testing-event","event_payload":{"id":13,"name":"testing-name"}}`},
				"channel":      []string{"CXXX"},
				"token":        []string{"testing-token"},
				"mrkdwn":       []string{"false"},
				"unfurl_media": []string{"false"},
			},
		},
	}

	once.Do(startServer)
	api := New(validToken, OptionAPIURL("http://"+serverAddr+"/"))

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			http.DefaultServeMux = new(http.ServeMux)
			http.HandleFunc(test.endpoint, func(rw http.ResponseWriter, r *http.Request) {
				body, err := io.ReadAll(r.Body)
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
		body, err := io.ReadAll(r.Body)
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
		body, err := io.ReadAll(r.Body)
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
		body, err := io.ReadAll(r.Body)
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

func TestSendMessageContextRedactsTokenInDebugLog(t *testing.T) {
	tests := []struct {
		name  string
		token string
		want  string
	}{
		{
			name:  "regular token",
			token: "xtest-token-1234-abcd",
			want:  "xtest-REDACTED",
		},
		{
			name:  "refresh token",
			token: "xoxe.xtest-token-1234-abcd",
			want:  "xoxe.xtest-REDACTED",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			once.Do(startServer)
			buf := bytes.NewBufferString("")

			opts := []Option{
				OptionAPIURL("http://" + serverAddr + "/"),
				OptionLog(log.New(buf, "", log.Lshortfile)),
				OptionDebug(true),
			}
			api := New(tt.token, opts...)
			// Why send the token in the message text too? To test that we're not
			// redacting substrings in the request which look like a token but aren't.
			api.SendMessage("CXXX", MsgOptionText(token, false))
			s := buf.String()

			re := regexp.MustCompile(`token=[\w.-]*`)
			want := "token=" + tt.want
			if got := re.FindString(s); got != want {
				t.Errorf("Logged token in SendMessageContext(): got %q, want %q", got, want)
			}
			re = regexp.MustCompile(`text=[\w.-]*`)
			want = "text=" + token
			if got := re.FindString(s); got != want {
				t.Errorf("Logged text in SendMessageContext(): got %q, want %q", got, want)
			}
		})
	}
}

func TestUpdateMessage(t *testing.T) {
	type messageTest struct {
		endpoint string
		opt      []MsgOption
		expected url.Values
	}
	tests := map[string]messageTest{
		"empty file_ids": {
			endpoint: "/chat.update",
			opt:      []MsgOption{},
			expected: url.Values{
				"channel": []string{"CXXX"},
				"token":   []string{"testing-token"},
				"ts":      []string{"1234567890.123456"},
			},
		},
		"with file_ids": {
			endpoint: "/chat.update",
			opt: []MsgOption{
				MsgOptionFileIDs([]string{"F123", "F456"}),
			},
			expected: url.Values{
				"channel":  []string{"CXXX"},
				"token":    []string{"testing-token"},
				"ts":       []string{"1234567890.123456"},
				"file_ids": []string{`["F123","F456"]`},
			},
		},
	}

	once.Do(startServer)
	api := New(validToken, OptionAPIURL("http://"+serverAddr+"/"))

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			http.DefaultServeMux = new(http.ServeMux)
			http.HandleFunc(test.endpoint, func(rw http.ResponseWriter, r *http.Request) {
				body, err := io.ReadAll(r.Body)
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

			_, _, _, _ = api.UpdateMessage("CXXX", "1234567890.123456", test.opt...)
		})
	}
}

func TestStartStream(t *testing.T) {
	type messageTest struct {
		endpoint string
		opt      []MsgOption
		expected url.Values
	}
	tests := map[string]messageTest{
		"basic": {
			endpoint: "/chat.startStream",
			opt: []MsgOption{
				MsgOptionTS("1234567890.123456"),
			},
			expected: url.Values{
				"channel":   []string{"CXXX"},
				"token":     []string{"testing-token"},
				"thread_ts": []string{"1234567890.123456"},
			},
		},
		"with recipients": {
			endpoint: "/chat.startStream",
			opt: []MsgOption{
				MsgOptionTS("1234567890.123456"),
				MsgOptionRecipientTeamID("T12345"),
				MsgOptionRecipientUserID("U12345"),
			},
			expected: url.Values{
				"channel":           []string{"CXXX"},
				"token":             []string{"testing-token"},
				"thread_ts":         []string{"1234567890.123456"},
				"recipient_team_id": []string{"T12345"},
				"recipient_user_id": []string{"U12345"},
			},
		},
	}

	once.Do(startServer)
	api := New(validToken, OptionAPIURL("http://"+serverAddr+"/"))

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			http.DefaultServeMux = new(http.ServeMux)
			http.HandleFunc(test.endpoint, func(rw http.ResponseWriter, r *http.Request) {
				body, err := io.ReadAll(r.Body)
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
				rw.Header().Set("Content-Type", "application/json")
				response, _ := json.Marshal(chatResponseFull{
					Channel:   "CXXX",
					Timestamp: "1234567890.123456",
					SlackResponse: SlackResponse{
						Ok: true,
					},
				})
				rw.Write(response)
			})

			_, _, _ = api.StartStream("CXXX", test.opt...)
		})
	}
}

func TestWorkObjectMetadata(t *testing.T) {
	// Test WorkObjectMetadata marshaling
	metadata := WorkObjectMetadata{
		Entities: []WorkObjectEntity{
			{
				AppUnfurlURL: "https://example.com/document/123?eid=123456&edit=abcxyz",
				URL:          "https://example.com/document/123",
				ExternalRef: WorkObjectExternalRef{
					ID:   "123",
					Type: "document",
				},
				EntityType: "slack#/entities/file",
				EntityPayload: map[string]interface{}{
					"title":       "Test Document",
					"description": "A test document for Work Objects",
				},
			},
		},
	}

	// Test JSON marshaling
	jsonData, err := json.Marshal(metadata)
	if err != nil {
		t.Errorf("Failed to marshal WorkObjectMetadata: %v", err)
	}

	// Test JSON unmarshaling
	var unmarshaled WorkObjectMetadata
	err = json.Unmarshal(jsonData, &unmarshaled)
	if err != nil {
		t.Errorf("Failed to unmarshal WorkObjectMetadata: %v", err)
	}

	// Verify the data
	if len(unmarshaled.Entities) != 1 {
		t.Errorf("Expected 1 entity, got %d", len(unmarshaled.Entities))
	}

	entity := unmarshaled.Entities[0]
	if entity.URL != "https://example.com/document/123" {
		t.Errorf("Expected URL 'https://example.com/document/123', got '%s'", entity.URL)
	}

	if entity.ExternalRef.ID != "123" {
		t.Errorf("Expected external ref ID '123', got '%s'", entity.ExternalRef.ID)
	}

	if entity.EntityType != "slack#/entities/file" {
		t.Errorf("Expected entity type 'slack#/entities/file', got '%s'", entity.EntityType)
	}
}

func TestMsgOptionWorkObjectMetadata(t *testing.T) {
	metadata := WorkObjectMetadata{
		Entities: []WorkObjectEntity{
			{
				URL: "https://example.com/task/456",
				ExternalRef: WorkObjectExternalRef{
					ID: "456",
				},
				EntityType: "slack#/entities/task",
				EntityPayload: map[string]interface{}{
					"title":  "Test Task",
					"status": "in_progress",
				},
			},
		},
	}

	// Create a sendConfig to test the option
	config := &sendConfig{
		values: url.Values{},
	}

	// Apply the option
	opt := MsgOptionWorkObjectMetadata(metadata)
	err := opt(config)
	if err != nil {
		t.Errorf("MsgOptionWorkObjectMetadata returned error: %v", err)
	}

	// Check that metadata was set
	metadataValue := config.values.Get("metadata")
	if metadataValue == "" {
		t.Error("Expected metadata to be set, but it was empty")
	}

	// Verify the JSON structure
	var result WorkObjectMetadata
	err = json.Unmarshal([]byte(metadataValue), &result)
	if err != nil {
		t.Errorf("Failed to unmarshal metadata JSON: %v", err)
	}

	if len(result.Entities) != 1 {
		t.Errorf("Expected 1 entity, got %d", len(result.Entities))
	}
}

func TestMsgOptionWorkObjectMetadataNilEntities(t *testing.T) {
	// When Entities is nil, we should marshal as "entities":[] for API compatibility
	metadata := WorkObjectMetadata{Entities: nil}
	config := &sendConfig{values: url.Values{}}
	opt := MsgOptionWorkObjectMetadata(metadata)
	if err := opt(config); err != nil {
		t.Errorf("MsgOptionWorkObjectMetadata with nil Entities returned error: %v", err)
	}
	metadataValue := config.values.Get("metadata")
	if metadataValue == "" {
		t.Error("Expected metadata to be set")
	}
	if metadataValue != `{"entities":[]}` {
		t.Errorf("Expected metadata with empty entities array, got %q", metadataValue)
	}
}

func TestMsgOptionWorkObjectEntity(t *testing.T) {
	entity := WorkObjectEntity{
		URL: "https://example.com/incident/789",
		ExternalRef: WorkObjectExternalRef{
			ID:   "789",
			Type: "incident",
		},
		EntityType: "slack#/entities/incident",
		EntityPayload: map[string]interface{}{
			"title":    "Production Outage",
			"severity": "high",
		},
	}

	// Create a sendConfig to test the option
	config := &sendConfig{
		values: url.Values{},
	}

	// Apply the option
	opt := MsgOptionWorkObjectEntity(entity)
	err := opt(config)
	if err != nil {
		t.Errorf("MsgOptionWorkObjectEntity returned error: %v", err)
	}

	// Check that metadata was set
	metadataValue := config.values.Get("metadata")
	if metadataValue == "" {
		t.Error("Expected metadata to be set, but it was empty")
	}

	// Verify the JSON structure
	var result WorkObjectMetadata
	err = json.Unmarshal([]byte(metadataValue), &result)
	if err != nil {
		t.Errorf("Failed to unmarshal metadata JSON: %v", err)
	}

	if len(result.Entities) != 1 {
		t.Errorf("Expected 1 entity, got %d", len(result.Entities))
	}

	resultEntity := result.Entities[0]
	if resultEntity.URL != entity.URL {
		t.Errorf("Expected URL '%s', got '%s'", entity.URL, resultEntity.URL)
	}

	if resultEntity.ExternalRef.ID != entity.ExternalRef.ID {
		t.Errorf("Expected external ref ID '%s', got '%s'", entity.ExternalRef.ID, resultEntity.ExternalRef.ID)
	}
}

func TestAppendStream(t *testing.T) {
	type messageTest struct {
		endpoint string
		opt      []MsgOption
		expected url.Values
	}
	tests := map[string]messageTest{
		"basic": {
			endpoint: "/chat.appendStream",
			opt: []MsgOption{
				MsgOptionMarkdownText("Hello, world!"),
			},
			expected: url.Values{
				"channel":       []string{"CXXX"},
				"token":         []string{"testing-token"},
				"ts":            []string{"1234567890.123456"},
				"markdown_text": []string{"Hello, world!"},
			},
		},
	}

	once.Do(startServer)
	api := New(validToken, OptionAPIURL("http://"+serverAddr+"/"))

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			http.DefaultServeMux = new(http.ServeMux)
			http.HandleFunc(test.endpoint, func(rw http.ResponseWriter, r *http.Request) {
				body, err := io.ReadAll(r.Body)
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
				rw.Header().Set("Content-Type", "application/json")
				response, _ := json.Marshal(chatResponseFull{
					Channel:   "CXXX",
					Timestamp: "1234567890.123456",
					SlackResponse: SlackResponse{
						Ok: true,
					},
				})
				rw.Write(response)
			})

			_, _, _ = api.AppendStream("CXXX", "1234567890.123456", test.opt...)
		})
	}
}

func TestStopStream(t *testing.T) {
	type messageTest struct {
		endpoint string
		opt      []MsgOption
		expected url.Values
	}

	blocks := []Block{NewContextBlock("context", NewTextBlockObject(PlainTextType, "feedback", false, false))}
	blockStr := `[{"type":"context","block_id":"context","elements":[{"type":"plain_text","text":"feedback","emoji":false}]}]`

	tests := map[string]messageTest{
		"basic": {
			endpoint: "/chat.stopStream",
			opt:      []MsgOption{},
			expected: url.Values{
				"channel": []string{"CXXX"},
				"token":   []string{"testing-token"},
				"ts":      []string{"1234567890.123456"},
			},
		},
		"with final text and blocks": {
			endpoint: "/chat.stopStream",
			opt: []MsgOption{
				MsgOptionMarkdownText("Final message"),
				MsgOptionBlocks(blocks...),
			},
			expected: url.Values{
				"channel":       []string{"CXXX"},
				"token":         []string{"testing-token"},
				"ts":            []string{"1234567890.123456"},
				"markdown_text": []string{"Final message"},
				"blocks":        []string{blockStr},
			},
		},
	}

	once.Do(startServer)
	api := New(validToken, OptionAPIURL("http://"+serverAddr+"/"))

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			http.DefaultServeMux = new(http.ServeMux)
			http.HandleFunc(test.endpoint, func(rw http.ResponseWriter, r *http.Request) {
				body, err := io.ReadAll(r.Body)
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
				rw.Header().Set("Content-Type", "application/json")
				response, _ := json.Marshal(chatResponseFull{
					Channel:   "CXXX",
					Timestamp: "1234567890.123456",
					SlackResponse: SlackResponse{
						Ok: true,
					},
				})
				rw.Write(response)
			})

			_, _, _ = api.StopStream("CXXX", "1234567890.123456", test.opt...)
		})
	}
}
