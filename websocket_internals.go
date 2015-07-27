package slack

import (
	"fmt"
	"time"
)

/**
 * Internal events, created by this lib and not mapped to Slack APIs.
 */
type ConnectedEvent struct {
	ConnectionCount int // 1 = first time, 2 = second time
	Info            *Info
}

type ConnectionErrorEvent struct {
	Attempt  int
	ErrorObj error
}

func (c *ConnectionErrorEvent) Error() string {
	return c.ErrorObj.Error()
}

type ConnectingEvent struct {
	Attempt         int // 1 = first attempt, 2 = second attempt
	ConnectionCount int
}

type DisconnectedEvent struct {
	Intentional bool
}

type LatencyReport struct {
	Value time.Duration
}

type InvalidAuthEvent struct{}

type UnmarshallingErrorEvent struct {
	ErrorObj error
}

func (u UnmarshallingErrorEvent) Error() string {
	return u.ErrorObj.Error()
}

type MessageTooLongEvent struct {
	Message   OutgoingMessage
	MaxLength int
}

func (m *MessageTooLongEvent) Error() string {
	return fmt.Sprintf("Message too long (max %d characters)", m.MaxLength)
}

type OutgoingErrorEvent struct {
	Message  OutgoingMessage
	ErrorObj error
}

func (o OutgoingErrorEvent) Error() string {
	return o.ErrorObj.Error()
}

type IncomingEventError struct {
	ErrorObj error
}

func (i *IncomingEventError) Error() string {
	return i.ErrorObj.Error()
}

type AckErrorEvent struct {
	ErrorObj error
}

func (a *AckErrorEvent) Error() string {
	return a.ErrorObj.Error()
}

type SlackErrorEvent struct {
	ErrorObj error
}

func (s SlackErrorEvent) Error() string {
	return s.ErrorObj.Error()
}
