package slack

import (
	"encoding/json"
	"fmt"
)

// ConditionType is the kind of check a Condition makes.
type ConditionType string

const (
	// ConditionTypeClient matches the Slack client the message is viewed on.
	ConditionTypeClient ConditionType = "client"
)

// ConditionOperator compares a Condition's value with the viewer's.
type ConditionOperator string

const (
	ConditionOperatorEqual ConditionOperator = "="
)

// ConditionClient is a Slack client a ConditionTypeClient condition can match.
type ConditionClient string

const (
	ConditionClientDesktop ConditionClient = "desktop"
	ConditionClientIOS     ConditionClient = "ios"
	ConditionClientAndroid ConditionClient = "android"
)

// Condition is one check in a ConditionCase.
type Condition struct {
	Type  ConditionType     `json:"type"`
	Op    ConditionOperator `json:"op"`
	Value string            `json:"value"`
}

// NewClientCondition returns a condition that matches the given Slack client.
func NewClientCondition(client ConditionClient) Condition {
	return Condition{Type: ConditionTypeClient, Op: ConditionOperatorEqual, Value: string(client)}
}

// ConditionCase renders Blocks when all of its Conditions match.
type ConditionCase struct {
	Conditions []Condition `json:"conditions"`
	Blocks     Blocks      `json:"blocks"`
}

// NewConditionCase returns a case that renders blocks when every condition matches.
func NewConditionCase(conditions []Condition, blocks ...Block) ConditionCase {
	return ConditionCase{Conditions: conditions, Blocks: Blocks{BlockSet: blocks}}
}

// ConditionBlock renders the blocks of the first case whose conditions match the
// viewer, or Default when none does. An empty Default renders nothing.
//
// Slack does not document this block yet. Its shape matches what Block Kit Builder
// produces.
type ConditionBlock struct {
	Type    MessageBlockType `json:"type"`
	BlockID string           `json:"block_id,omitempty"`
	Cases   []ConditionCase  `json:"cases"`
	Default Blocks           `json:"default"`
}

// BlockType returns the type of the block.
func (s ConditionBlock) BlockType() MessageBlockType {
	return s.Type
}

// ID returns the ID of the block.
func (s ConditionBlock) ID() string {
	return s.BlockID
}

// MarshalJSON writes empty block lists as [] rather than null.
func (s ConditionBlock) MarshalJSON() ([]byte, error) {
	type alias ConditionBlock
	a := alias(s)
	a.Cases = make([]ConditionCase, len(s.Cases))
	for i, c := range s.Cases {
		if c.Blocks.BlockSet == nil {
			c.Blocks.BlockSet = []Block{}
		}
		a.Cases[i] = c
	}
	if a.Default.BlockSet == nil {
		a.Default.BlockSet = []Block{}
	}
	return json.Marshal(a)
}

// Validate checks that the block has at least one case, and that every case has at
// least one condition.
func (s ConditionBlock) Validate() error {
	if s.Type != MBTCondition {
		return fmt.Errorf("type must be %q", MBTCondition)
	}
	if len(s.Cases) == 0 {
		return fmt.Errorf("cases must have at least 1 case")
	}
	for i, c := range s.Cases {
		if len(c.Conditions) == 0 {
			return fmt.Errorf("case %d must have at least 1 condition", i)
		}
	}
	return nil
}

// NewConditionBlock returns a new condition block with the given cases and an empty
// default.
func NewConditionBlock(cases ...ConditionCase) *ConditionBlock {
	return &ConditionBlock{
		Type:  MBTCondition,
		Cases: cases,
	}
}

// WithBlockID sets the block ID for the ConditionBlock.
func (s *ConditionBlock) WithBlockID(blockID string) *ConditionBlock {
	s.BlockID = blockID
	return s
}

// WithDefault sets the blocks rendered when no case matches.
func (s *ConditionBlock) WithDefault(blocks ...Block) *ConditionBlock {
	s.Default = Blocks{BlockSet: blocks}
	return s
}

// AddCase appends a case. Cases are checked in order.
func (s *ConditionBlock) AddCase(c ConditionCase) *ConditionBlock {
	s.Cases = append(s.Cases, c)
	return s
}
