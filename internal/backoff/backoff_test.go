package backoff

import (
	"testing"
	"time"
)

func TestBackoffDurationDoublesThenCapsAtMax(t *testing.T) {
	b := &Backoff{Initial: 100 * time.Millisecond, Max: time.Second}

	want := []time.Duration{
		100 * time.Millisecond,
		200 * time.Millisecond,
		400 * time.Millisecond,
		800 * time.Millisecond,
		time.Second, // 1.6s capped
		time.Second,
		time.Second,
	}
	for i, w := range want {
		if got := b.Duration(); got != w {
			t.Errorf("attempt %d: got %v, want %v", i, got, w)
		}
	}
}

// Duration must stay at Max no matter how many attempts have been made. Without
// a cap the delay grows as 2^attempts*Initial, which reaches hours after a few
// dozen attempts, and eventually overflows to a non-positive value.
func TestBackoffDurationStaysCappedForLargeAttemptCounts(t *testing.T) {
	b := &Backoff{Initial: 100 * time.Millisecond, Max: 5 * time.Minute}

	for i := range 200 {
		got := b.Duration()
		if got <= 0 {
			t.Fatalf("attempt %d: got non-positive duration %v", i, got)
		}
		if got > b.Max {
			t.Fatalf("attempt %d: got %v, want at most %v", i, got, b.Max)
		}
	}
}

func TestBackoffDurationAppliesDefaults(t *testing.T) {
	b := &Backoff{}

	if got := b.Duration(); got != 100*time.Millisecond {
		t.Errorf("first duration: got %v, want the default Initial of 100ms", got)
	}
	if b.Max != 10*time.Second {
		t.Errorf("Max: got %v, want the default of 10s", b.Max)
	}
}

func TestBackoffDurationWithJitter(t *testing.T) {
	b := &Backoff{Initial: 100 * time.Millisecond, Max: 200 * time.Millisecond, Jitter: 50 * time.Millisecond}

	// Jitter is added after the cap, so the upper bound is Max+Jitter.
	for i := range 50 {
		got := b.Duration()
		if got < 100*time.Millisecond || got >= b.Max+b.Jitter {
			t.Fatalf("attempt %d: got %v, want within [100ms, %v)", i, got, b.Max+b.Jitter)
		}
	}
}

func TestBackoffReset(t *testing.T) {
	b := &Backoff{Initial: 100 * time.Millisecond, Max: time.Second}

	for range 5 {
		_ = b.Duration()
	}
	if b.Attempts() != 5 {
		t.Errorf("Attempts: got %d, want 5", b.Attempts())
	}

	b.Reset()
	if b.Attempts() != 0 {
		t.Errorf("Attempts after Reset: got %d, want 0", b.Attempts())
	}
	if got := b.Duration(); got != 100*time.Millisecond {
		t.Errorf("first duration after Reset: got %v, want 100ms", got)
	}
}
