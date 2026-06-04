package slack

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestBlockRealPayloadRoundTrip is a guardrail against silent data loss when decoding the
// JSON that Slack actually sends.
//
// Each case is a single block payload written as raw JSON in the shape Slack delivers
// (i.e. authored from the JSON in, not built via the SDK constructors — that distinction
// matters, see below). Every case is decoded through the same path inbound events take —
// Blocks.UnmarshalJSON, which dispatches to the concrete block type — and the test asserts:
//
//  1. it decodes into a recognised concrete block (never UnknownBlock), and
//  2. re-marshalling the decoded block reproduces the input exactly (semantic JSON equality).
//
// (2) is the important part: a field the concrete type fails to model — like the per-cell
// text that went missing in issue #1558 — disappears on the way back out and fails the
// comparison. Constructor-based tests can't catch this, because they only ever serialise
// what the SDK already knows how to represent.
//
// Note this deliberately does NOT use BlockFromJSON: that helper preserves the raw bytes and
// echoes them back unchanged, so it would round-trip any payload perfectly while exercising
// none of the concrete UnmarshalJSON logic this guardrail exists to protect.
//
// To cover a new block, add a case below with a payload in the shape Slack sends. If it
// doesn't round-trip, the SDK is dropping something Slack sent — fix the type, don't trim
// the payload.
func TestBlockRealPayloadRoundTrip(t *testing.T) {
	tests := []struct {
		name    string
		payload string
	}{
		{
			// A table pasted from a spreadsheet (issue #1558): a bold rich_text header
			// row, raw_text/raw_number data cells, and a null cell for an empty space.
			name: "table with mixed cell types and an empty cell",
			payload: `{
				"type": "table",
				"block_id": "tbl_pasted",
				"rows": [
					[
						{"type": "rich_text", "elements": [{"type": "rich_text_section", "elements": [{"type": "text", "text": "Name", "style": {"bold": true}}]}]},
						{"type": "rich_text", "elements": [{"type": "rich_text_section", "elements": [{"type": "text", "text": "Score", "style": {"bold": true}}]}]}
					],
					[
						{"type": "raw_text", "text": "Alice"},
						{"type": "raw_number", "value": 42}
					],
					[
						{"type": "raw_text", "text": "Bob"},
						null
					]
				]
			}`,
		},
		{
			name: "data_table with raw_text, raw_number and rich_text cells",
			payload: `{
				"type": "data_table",
				"block_id": "dt-1",
				"caption": "A Fabulous Table",
				"page_size": 5,
				"rows": [
					[
						{"type": "raw_text", "text": "Name"},
						{"type": "raw_text", "text": "Department"},
						{"type": "raw_text", "text": "Badge"}
					],
					[
						{"type": "raw_text", "text": "Helly"},
						{"type": "raw_text", "text": "MDR"},
						{"type": "rich_text", "elements": [{"type": "rich_text_section", "elements": [{"type": "text", "text": "Blue", "style": {"bold": true}}]}]}
					],
					[
						{"type": "raw_text", "text": "Score"},
						{"type": "raw_text", "text": "Wellness"},
						{"type": "raw_number", "value": 97}
					]
				]
			}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Slack delivers blocks as a JSON array; Blocks.UnmarshalJSON is what
			// dispatches each element to its concrete type.
			var blocks Blocks
			require.NoError(t, json.Unmarshal([]byte("["+tt.payload+"]"), &blocks), "payload failed to decode")
			require.Len(t, blocks.BlockSet, 1)

			block := blocks.BlockSet[0]
			_, isUnknown := block.(*UnknownBlock)
			assert.False(t, isUnknown,
				"payload decoded to an UnknownBlock; the block type is not modelled")

			marshalled, err := json.Marshal(block)
			require.NoError(t, err)

			var want, got any
			require.NoError(t, json.Unmarshal([]byte(tt.payload), &want))
			require.NoError(t, json.Unmarshal(marshalled, &got))
			assert.Equal(t, want, got,
				"round-trip lost or altered data: the concrete block type is not preserving everything Slack sent")
		})
	}
}
