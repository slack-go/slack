package errorsx

import (
	"errors"
	"testing"
)

func TestIs(t *testing.T) {
	test := func(receiver String, target error, expected bool) {
		if result := receiver.Is(target); result != expected {
			t.Errorf(`expected "%s Is %s" should be %t , got: %t`, receiver, target, expected, result)
		}
	}

	test(String("test error"), errors.New("test error"), true)
	test(String("test error"), errors.New("different error"), false)
	test(String("test error"), nil, false)
}
