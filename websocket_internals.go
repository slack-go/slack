package slack

import "time"

/**
 * Internal events, created by this lib and not mapped to Slack APIs.
 */
type Disconnected struct{}
type Reconnecting struct {
	Attempt int
}
type Connected struct{}

type LatencyReport struct {
	Value time.Duration
}
