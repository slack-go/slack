package socketmode

import "time"

type deadmanTimer struct {
	timeout time.Duration
	timer   *time.Timer
}

func newDeadmanTimer(timeout time.Duration) *deadmanTimer {
	return &deadmanTimer{
		timeout: timeout,
		timer:   time.NewTimer(timeout),
	}
}

func (smc *deadmanTimer) Elapsed() <-chan time.Time {
	return smc.timer.C
}

func (smc *deadmanTimer) Reset() {
	// FIXME: Race on "deadmanTimer", timer channel cannot be read concurrently while resetting.
	// "This should not be done concurrent to other receives from the Timer's channel."
	// https://pkg.go.dev/time#Timer.Reset
	// See socket_mode_managed_conn.go lines ~59 & ~151.
	if !smc.timer.Stop() {
		select {
		case <-smc.timer.C:
		default:
		}
	}

	smc.timer.Reset(smc.timeout)
}
