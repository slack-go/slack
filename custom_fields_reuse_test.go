package slack

import (
	"encoding/json"
	"testing"
)

func TestUserProfileCustomFieldsUnmarshalReusedValue(t *testing.T) {
	var fields UserProfileCustomFields
	if err := json.Unmarshal([]byte(`{"Xxxxxx":{"value":"old","alt":""}}`), &fields); err != nil {
		t.Fatal(err)
	}
	if err := json.Unmarshal([]byte(`{"Yyyyyy":{"value":"new","alt":""}}`), &fields); err != nil {
		t.Fatal(err)
	}
	if got := fields.Len(); got != 1 {
		t.Errorf("Len() = %d after reusing the decoder target, want 1: %v", got, fields.ToMap())
	}
	if _, ok := fields.ToMap()["Xxxxxx"]; ok {
		t.Error("stale field from the first decode survived the second decode")
	}
}

func TestUserProfileCustomFieldsEmptyArrayClearsValue(t *testing.T) {
	var fields UserProfileCustomFields
	if err := json.Unmarshal([]byte(`{"Xxxxxx":{"value":"old","alt":""}}`), &fields); err != nil {
		t.Fatal(err)
	}
	if err := json.Unmarshal([]byte(`[]`), &fields); err != nil {
		t.Fatal(err)
	}
	if got := fields.Len(); got != 0 {
		t.Errorf("Len() = %d after decoding [], want 0: %v", got, fields.ToMap())
	}
}
