package slacksocketmode

import "time"

func deadmanDuration(d time.Duration) time.Duration {
	return d * 4
}
