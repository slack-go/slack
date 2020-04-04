package slack

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestInputBlock_UnmarshalJSON_MultiStaticSelect(t *testing.T) {
	block := NewInputBlock("", nil, NewOptionsMultiSelectBlockElement(MultiOptTypeStatic, nil, ""))

	data, _ := json.Marshal(block)

	var newBlock = InputBlock{}
	err := json.Unmarshal(data, &newBlock)

	assert.Nil(t, err)
}
