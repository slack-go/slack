package slack

import "time"

/**
 * Internal events, created by this lib and not mapped to Slack APIs.
 */
type ConnectedEvent struct {
	ConnectionCount int // 1 = first time, 2 = second time
	Info            *Info
}

type ConnectingEvent struct {
	Attempt         int // 1 = first attempt, 2 = second attempt
	ConnectionCount int
}

type DisconnectedEvent struct{}

type LatencyReport struct {
	Value time.Duration
}

type InvalidAuthEvent struct{}
