package slack

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/go-test/deep"
)

func TestAttachment_UnmarshalMarshalJSON_WithBlocks(t *testing.T) {

	originalAttachmentJson := `{
    "id": 1,
    "blocks": [
      {
        "type": "section",
        "block_id": "xxxx",
        "text": {
          "type": "mrkdwn",
          "text": "Pick something:",
          "verbatim": true
        },
        "accessory": {
          "type": "static_select",
          "action_id": "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx",
          "placeholder": {
            "type": "plain_text",
            "text": "Select one item",
            "emoji": true
          },
          "options": [
            {
              "text": {
                "type": "plain_text",
                "text": "ghi",
                "emoji": true
              },
              "value": "ghi"
            }
          ]
        }
      }
    ],
    "color": "#13A554",
    "hide_color": true,
    "fallback": "[no preview available]"
  }`

	attachment := new(Attachment)
	err := json.Unmarshal([]byte(originalAttachmentJson), attachment)
	if err != nil {
		t.Fatalf("expected no error unmarshaling attachment with blocks, got: %v", err)
	}

	actualAttachmentJson, err := json.Marshal(attachment)
	if err != nil {
		t.Fatal(err)
	}

	var (
		actual   any
		expected any
	)
	if err = json.Unmarshal([]byte(originalAttachmentJson), &expected); err != nil {
		t.Fatal(err)
	}
	if err = json.Unmarshal(actualAttachmentJson, &actual); err != nil {
		t.Fatal(err)
	}
	if diff := deep.Equal(actual, expected); diff != nil {
		t.Fatal("actual does not match expected\n", strings.Join(diff, "\n"))
	}
}

func TestAttachment_UnmarshalJSON_ForwardedMessage(t *testing.T) {
	forwardedAttachmentJson := `{
    "is_share": true,
    "is_msg_unfurl": true,
    "from_url": "https://example.slack.com/archives/C123/p1700000000000100",
    "files": [
      {
        "id": "F123",
        "mimetype": "image/png",
        "url_private_download": "https://files.slack.com/files-pri/T123-F123/download/a.png"
      }
    ]
  }`

	attachment := new(Attachment)
	if err := json.Unmarshal([]byte(forwardedAttachmentJson), attachment); err != nil {
		t.Fatalf("expected no error unmarshaling forwarded attachment, got: %v", err)
	}

	if !attachment.IsShare || !attachment.IsMsgUnfurl {
		t.Fatalf("expected IsShare and IsMsgUnfurl, got %v and %v", attachment.IsShare, attachment.IsMsgUnfurl)
	}
	if len(attachment.Files) != 1 {
		t.Fatalf("expected 1 file, got %d", len(attachment.Files))
	}
	if diff := deep.Equal(attachment.Files[0], File{
		ID:                 "F123",
		Mimetype:           "image/png",
		URLPrivateDownload: "https://files.slack.com/files-pri/T123-F123/download/a.png",
	}); diff != nil {
		t.Fatal("file does not match expected\n", strings.Join(diff, "\n"))
	}
}
