package slack

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBlockElementFieldsRoundTrip(t *testing.T) {
	tests := []struct {
		name       string
		payload    string
		element    BlockElement
		containers []string
	}{
		{
			name: "button",
			payload: `{
				"type": "button",
				"text": {"type": "plain_text", "text": "Summarize"},
				"action_id": "summarize",
				"accessibility_label": "Summarize this conversation",
				"agent_prompt": "Summarize the conversation so far",
				"agent_prompt_display": "Summarize conversation",
				"visible_to_user_ids": ["U0123456789", "U9876543210"]
			}`,
			element: &ButtonBlockElement{
				Type:               METButton,
				Text:               &TextBlockObject{Type: "plain_text", Text: "Summarize"},
				ActionID:           "summarize",
				AccessibilityLabel: "Summarize this conversation",
				AgentPrompt:        "Summarize the conversation so far",
				AgentPromptDisplay: "Summarize conversation",
				VisibleToUserIDs:   []string{"U0123456789", "U9876543210"},
			},
			containers: []string{"actions", "section"},
		},
		{
			name: "checkboxes",
			payload: `{
				"type": "checkboxes",
				"action_id": "notify",
				"options": [{"text": {"type": "plain_text", "text": "Notify me"}, "value": "notify"}],
				"visible_to_user_ids": ["U0123456789", "U9876543210"]
			}`,
			element: &CheckboxGroupsBlockElement{
				Type:     METCheckboxGroups,
				ActionID: "notify",
				Options: []*OptionBlockObject{{
					Text:  &TextBlockObject{Type: "plain_text", Text: "Notify me"},
					Value: "notify",
				}},
				VisibleToUserIDs: []string{"U0123456789", "U9876543210"},
			},
			containers: []string{"actions", "section", "input"},
		},
		{
			name: "conversation multi-select",
			payload: `{
				"type": "multi_conversations_select",
				"action_id": "conversations",
				"default_to_current_conversation": true
			}`,
			element: &MultiSelectBlockElement{
				Type:                         MultiOptTypeConversations,
				ActionID:                     "conversations",
				DefaultToCurrentConversation: true,
			},
			containers: []string{"actions", "section", "input"},
		},
		{
			name: "rich text input",
			payload: `{
				"type": "rich_text_input",
				"action_id": "description",
				"min_lines": 1,
				"max_lines": 100
			}`,
			element: &RichTextInputBlockElement{
				Type:     METRichTextInput,
				ActionID: "description",
				MinLines: 1,
				MaxLines: 100,
			},
			containers: []string{"input"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name+"/marshal", func(t *testing.T) {
			encoded, err := json.Marshal(tt.element)
			require.NoError(t, err)
			assert.JSONEq(t, tt.payload, string(encoded))
		})
		for _, container := range tt.containers {
			t.Run(tt.name+"/"+container, func(t *testing.T) {
				var payload string
				switch container {
				case "actions":
					payload = fmt.Sprintf(`[{"type":"actions","elements":[%s]}]`, tt.payload)
				case "section":
					payload = fmt.Sprintf(`[{"type":"section","text":{"type":"mrkdwn","text":"Choose an option"},"accessory":%s}]`, tt.payload)
				case "input":
					payload = fmt.Sprintf(`[{"type":"input","label":{"type":"plain_text","text":"Choose an option"},"element":%s}]`, tt.payload)
				}

				var blocks Blocks
				require.NoError(t, json.Unmarshal([]byte(payload), &blocks))
				require.Len(t, blocks.BlockSet, 1)

				var element BlockElement
				switch block := blocks.BlockSet[0].(type) {
				case *ActionBlock:
					require.NotNil(t, block.Elements)
					require.Len(t, block.Elements.ElementSet, 1)
					element = block.Elements.ElementSet[0]
				case *SectionBlock:
					require.NotNil(t, block.Accessory)
					element = toBlockElement(block.Accessory)
				case *InputBlock:
					element = block.Element
				default:
					t.Fatalf("unexpected block type %T", block)
				}
				require.IsType(t, tt.element, element)
				assert.Equal(t, tt.element, element)

				encoded, err := json.Marshal(blocks)
				require.NoError(t, err)
				assert.JSONEq(t, payload, string(encoded))
			})
		}
	}
}

func TestBlockElementFieldsOmitZeroValues(t *testing.T) {
	text := &TextBlockObject{Type: "plain_text", Text: "Choose"}
	option := NewOptionBlockObject("choose", text, nil)
	tests := []struct {
		name    string
		element BlockElement
		payload string
	}{
		{
			name:    "button",
			element: NewButtonBlockElement("choose", "", text),
			payload: `{"type":"button","text":{"type":"plain_text","text":"Choose"},"action_id":"choose"}`,
		},
		{
			name:    "checkboxes",
			element: NewCheckboxGroupsBlockElement("choose", option),
			payload: `{"type":"checkboxes","action_id":"choose","options":[{"text":{"type":"plain_text","text":"Choose"},"value":"choose"}]}`,
		},
		{
			name:    "conversation multi-select",
			element: NewOptionsMultiSelectBlockElement(MultiOptTypeConversations, nil, "choose"),
			payload: `{"type":"multi_conversations_select","action_id":"choose"}`,
		},
		{
			name:    "rich text input",
			element: NewRichTextInputBlockElement(nil, "choose"),
			payload: `{"type":"rich_text_input","action_id":"choose"}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			encoded, err := json.Marshal(tt.element)
			require.NoError(t, err)
			assert.JSONEq(t, tt.payload, string(encoded))
		})
	}
}
