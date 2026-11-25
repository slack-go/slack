package slack

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewConditionBlock(t *testing.T) {
	desktop := NewSectionBlock(NewTextBlockObject(MarkdownType, "Desktop", false, false), nil, nil)
	ios := NewSectionBlock(NewTextBlockObject(MarkdownType, "iOS", false, false), nil, nil)
	fallback := NewDividerBlock()

	block := NewConditionBlock(
		NewConditionCase([]Condition{NewClientCondition(ConditionClientDesktop)}, desktop),
	).
		AddCase(NewConditionCase([]Condition{NewClientCondition(ConditionClientIOS)}, ios)).
		WithDefault(fallback).
		WithBlockID("condition-1")

	assert.Equal(t, MBTCondition, block.BlockType())
	assert.Equal(t, "condition-1", block.ID())
	require.Len(t, block.Cases, 2)
	assert.Equal(t, Condition{Type: ConditionTypeClient, Op: ConditionOperatorEqual, Value: "ios"}, block.Cases[1].Conditions[0])
	assert.Equal(t, []Block{fallback}, block.Default.BlockSet)
	assert.NoError(t, block.Validate())
}

func TestConditionBlockMarshalJSON(t *testing.T) {
	block := NewConditionBlock(NewConditionCase([]Condition{NewClientCondition(ConditionClientDesktop)}))

	got, err := json.Marshal(block)
	require.NoError(t, err)
	assert.JSONEq(t, `{
		"type": "condition",
		"cases": [{"conditions": [{"type": "client", "op": "=", "value": "desktop"}], "blocks": []}],
		"default": []
	}`, string(got))
}

func TestConditionBlockRoundTrip(t *testing.T) {
	in := `[{
		"type": "condition",
		"cases": [{
			"conditions": [{"type": "client", "op": "=", "value": "android"}],
			"blocks": [{"type": "divider"}]
		}],
		"default": [{"type": "section", "text": {"type": "mrkdwn", "text": "Elsewhere"}}]
	}]`

	var blocks Blocks
	require.NoError(t, json.Unmarshal([]byte(in), &blocks))
	require.Len(t, blocks.BlockSet, 1)

	block, ok := blocks.BlockSet[0].(*ConditionBlock)
	require.True(t, ok, "expected *ConditionBlock, got %T", blocks.BlockSet[0])
	require.Len(t, block.Cases, 1)
	assert.Equal(t, "android", block.Cases[0].Conditions[0].Value)
	assert.IsType(t, &DividerBlock{}, block.Cases[0].Blocks.BlockSet[0])
	assert.IsType(t, &SectionBlock{}, block.Default.BlockSet[0])

	out, err := json.Marshal(blocks)
	require.NoError(t, err)
	assert.JSONEq(t, in, string(out))
}

func TestConditionBlockValidate(t *testing.T) {
	tests := []struct {
		name    string
		block   *ConditionBlock
		wantErr bool
	}{
		{
			name:  "valid",
			block: NewConditionBlock(NewConditionCase([]Condition{NewClientCondition(ConditionClientIOS)})),
		},
		{
			name:    "no cases",
			block:   NewConditionBlock(),
			wantErr: true,
		},
		{
			name:    "case without conditions",
			block:   NewConditionBlock(NewConditionCase(nil, NewDividerBlock())),
			wantErr: true,
		},
		{
			name:    "wrong type",
			block:   &ConditionBlock{Type: MBTContainer, Cases: []ConditionCase{NewConditionCase([]Condition{NewClientCondition(ConditionClientIOS)})}},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.block.Validate()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
