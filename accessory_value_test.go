package slack

import "testing"

func TestNewAccessoryWithUnknownValueDoesNotPanic(t *testing.T) {
	a := NewAccessory(UnknownBlockElement{Type: "future_element"})
	if a == nil || a.UnknownElement == nil {
		t.Fatal("NewAccessory returned no unknown element")
	}
	if a.UnknownElement.Type != MessageElementType("future_element") {
		t.Errorf("unknown element type = %q, want future_element", a.UnknownElement.Type)
	}
}
