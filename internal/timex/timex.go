package timex

import "time"

// Max returns the maximum duration
func Max(values ...time.Duration) time.Duration {
	var (
		maxDuration time.Duration
	)

	for _, v := range values {
		if v > maxDuration {
			maxDuration = v
		}
	}

	return maxDuration
}
