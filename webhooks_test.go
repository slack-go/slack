package slack

import (
	"encoding/json"
	"net/http"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPostWebhook_OK(t *testing.T) {
	once.Do(startServer)

	var receivedPayload WebhookMessage

	http.HandleFunc("/webhook", func(rw http.ResponseWriter, r *http.Request) {
		rw.Header().Set("Content-Type", "application/json")

		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&receivedPayload)
		if err != nil {
			t.Errorf("Request contained invalid JSON, %s", err)
		}

		response := []byte(`{}`)
		rw.Write(response)
	})

	url := "http://" + serverAddr + "/webhook"

	payload := &WebhookMessage{
		Text: "Test Text",
		Attachments: []Attachment{
			{
				Text: "Foo",
			},
		},
	}

	err := PostWebhook(url, payload)

	if err != nil {
		t.Errorf("Expected not to receive error: %s", err)
	}

	if !reflect.DeepEqual(payload, &receivedPayload) {
		t.Errorf("Payload did not match\nwant: %#v\n got: %#v", payload, receivedPayload)
	}
}

func TestPostWebhook_NotOK(t *testing.T) {
	once.Do(startServer)

	http.HandleFunc("/webhook2", func(rw http.ResponseWriter, r *http.Request) {
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte("500 - Something bad happened!"))
	})

	url := "http://" + serverAddr + "/webhook2"

	err := PostWebhook(url, &WebhookMessage{})

	if err == nil {
		t.Errorf("Expected to receive error")
	}
}

func TestPostWebhook_MessageLimitExceeded(t *testing.T) {
	once.Do(startServer)

	http.HandleFunc("/message_limit_exceeded", func(rw http.ResponseWriter, r *http.Request) {
		// When a workspace's message limit is exceeded we get a 429 without a Retry-After header
		rw.WriteHeader(http.StatusTooManyRequests)
		rw.Write([]byte("message_limit_exceeded"))
	})

	url := "http://" + serverAddr + "/message_limit_exceeded"

	err := PostWebhook(url, &WebhookMessage{})

	if err == nil {
		t.Errorf("Expected to receive error")
	}
	var statusErr StatusCodeError
	assert.ErrorAs(t, err, &statusErr)
}

func TestWebhookMessage_UnfurlFields(t *testing.T) {
	t.Run("nil omits fields", func(t *testing.T) {
		msg := WebhookMessage{Text: "hello"}
		raw, err := json.Marshal(msg)
		require.NoError(t, err)
		assert.NotContains(t, string(raw), "unfurl_links")
		assert.NotContains(t, string(raw), "unfurl_media")
	})

	t.Run("false is preserved", func(t *testing.T) {
		msg := WebhookMessage{
			Text:        "hello",
			UnfurlLinks: new(false),
			UnfurlMedia: new(false),
		}
		raw, err := json.Marshal(msg)
		require.NoError(t, err)
		assert.Contains(t, string(raw), `"unfurl_links":false`)
		assert.Contains(t, string(raw), `"unfurl_media":false`)
	})

	t.Run("true is preserved", func(t *testing.T) {
		msg := WebhookMessage{
			Text:        "hello",
			UnfurlLinks: new(true),
			UnfurlMedia: new(true),
		}
		raw, err := json.Marshal(msg)
		require.NoError(t, err)
		assert.Contains(t, string(raw), `"unfurl_links":true`)
		assert.Contains(t, string(raw), `"unfurl_media":true`)
	})

	t.Run("round-trip preserves values", func(t *testing.T) {
		original := WebhookMessage{
			Text:        "hello",
			UnfurlLinks: new(false),
			UnfurlMedia: new(true),
		}
		raw, err := json.Marshal(original)
		require.NoError(t, err)

		var decoded WebhookMessage
		err = json.Unmarshal(raw, &decoded)
		require.NoError(t, err)
		assert.Equal(t, original.Text, decoded.Text)
		require.NotNil(t, decoded.UnfurlLinks)
		assert.False(t, *decoded.UnfurlLinks)
		require.NotNil(t, decoded.UnfurlMedia)
		assert.True(t, *decoded.UnfurlMedia)
	})
}

func TestWebhookMessage_WithBlocks(t *testing.T) {
	textBlockObject := NewTextBlockObject("plain_text", "text", false, false)
	sectionBlock := NewSectionBlock(textBlockObject, nil, nil)

	singleBlock := &Blocks{BlockSet: []Block{sectionBlock}}
	twoBlocks := &Blocks{BlockSet: []Block{sectionBlock, sectionBlock}}

	msgSingleBlock := WebhookMessage{Blocks: singleBlock}
	assert.Len(t, msgSingleBlock.Blocks.BlockSet, 1)

	msgJsonSingleBlock, _ := json.Marshal(msgSingleBlock)
	assert.JSONEq(t, `{"blocks":[{"type":"section","text":{"type":"plain_text","text":"text","emoji":false}}],"replace_original":false,"delete_original":false}`, string(msgJsonSingleBlock))

	msgTwoBlocks := WebhookMessage{Blocks: twoBlocks}
	assert.Len(t, msgTwoBlocks.Blocks.BlockSet, 2)

	msgJsonTwoBlocks, _ := json.Marshal(msgTwoBlocks)
	assert.JSONEq(t, `{"blocks":[{"type":"section","text":{"type":"plain_text","text":"text","emoji":false}},{"type":"section","text":{"type":"plain_text","text":"text","emoji":false}}],"replace_original":false,"delete_original":false}`, string(msgJsonTwoBlocks))

	msgNoBlocks := WebhookMessage{Text: "foo"}
	msgJsonNoBlocks, _ := json.Marshal(msgNoBlocks)
	assert.JSONEq(t, `{"text":"foo","replace_original":false,"delete_original":false}`, string(msgJsonNoBlocks))
}
