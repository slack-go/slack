package slack

import (
	"encoding/json"
	"testing"
)

func TestAccessoryUnknownElementMarshalPreservesType(t *testing.T) {
	var a Accessory
	if err := json.Unmarshal([]byte(`{"type":"future_element"}`), &a); err != nil {
		t.Fatal(err)
	}
	if a.UnknownElement == nil {
		t.Fatal("expected an unknown element after unmarshal")
	}

	out, err := json.Marshal(&a)
	if err != nil {
		t.Fatal(err)
	}
	if string(out) == "null" {
		t.Fatalf("marshal of an unknown accessory produced null, losing the element: %s", out)
	}
	var got map[string]any
	if err := json.Unmarshal(out, &got); err != nil {
		t.Fatal(err)
	}
	if got["type"] != "future_element" {
		t.Errorf("marshaled type = %v, want future_element (out=%s)", got["type"], out)
	}
}
