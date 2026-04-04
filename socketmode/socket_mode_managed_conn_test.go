//go:build go1.13
// +build go1.13

package socketmode

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slacktest"

	"github.com/stretchr/testify/assert"
)

func Test_passContext(t *testing.T) {
	s := slacktest.NewTestServer()
	go s.Start()

	api := slack.New("ABCDEFG", slack.OptionAPIURL(s.GetAPIURL()))
	cli := New(api)

	ctx, cancel := context.WithTimeout(context.TODO(), time.Nanosecond)
	defer cancel()

	t.Run("RunWithContext", func(t *testing.T) {
		// should fail imidiatly.
		assert.EqualError(t, cli.RunContext(ctx), context.DeadlineExceeded.Error())
	})

	t.Run("openAndDial", func(t *testing.T) {
		_, _, err := cli.openAndDial(ctx, func(_ string) error { return nil })

		// should fail imidiatly.
		assert.EqualError(t, errors.Unwrap(err), context.DeadlineExceeded.Error())
	})

	t.Run("OpenWithContext", func(t *testing.T) {
		_, _, err := cli.OpenContext(ctx)

		// should fail imidiatly.
		assert.EqualError(t, errors.Unwrap(err), context.DeadlineExceeded.Error())
	})
}

func TestSendCtx_PayloadQueued(t *testing.T) {
	api := slack.New("ABCDEFG")
	cli := New(api)

	err := cli.SendCtx(context.Background(), Response{
		EnvelopeID: "test-envelope",
		Payload:    "small payload",
	})

	assert.NoError(t, err)

	select {
	case res := <-cli.socketModeResponses:
		assert.Equal(t, "test-envelope", res.EnvelopeID)
	default:
		t.Fatal("expected response to be queued on socketModeResponses channel")
	}
}

func TestSendCtx_PayloadAtLimit(t *testing.T) {
	api := slack.New("ABCDEFG")
	cli := New(api)

	// Build a response that serializes to exactly maxResponseSize bytes.
	// The envelope and JSON overhead consume some bytes, so we pad the payload.
	envelope := "test-envelope"
	overhead, _ := json.Marshal(Response{EnvelopeID: envelope, Payload: ""})
	padding := strings.Repeat("x", maxResponseSize-len(overhead))

	// Verify we hit exactly the limit.
	exact, _ := json.Marshal(Response{EnvelopeID: envelope, Payload: padding})
	assert.Equal(t, maxResponseSize, len(exact), "test setup: payload should serialize to exactly %d bytes", maxResponseSize)

	err := cli.SendCtx(context.Background(), Response{
		EnvelopeID: envelope,
		Payload:    padding,
	})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "silently dropped")
}

func TestSendCtx_PayloadJustUnderLimit(t *testing.T) {
	api := slack.New("ABCDEFG")
	cli := New(api)

	envelope := "test-envelope"
	overhead, _ := json.Marshal(Response{EnvelopeID: envelope, Payload: ""})
	padding := strings.Repeat("x", maxResponseSize-len(overhead)-1)

	// Verify we're one byte under.
	under, _ := json.Marshal(Response{EnvelopeID: envelope, Payload: padding})
	assert.Equal(t, maxResponseSize-1, len(under), "test setup: payload should serialize to exactly %d bytes", maxResponseSize-1)

	err := cli.SendCtx(context.Background(), Response{
		EnvelopeID: envelope,
		Payload:    padding,
	})
	assert.NoError(t, err)

	select {
	case res := <-cli.socketModeResponses:
		assert.Equal(t, envelope, res.EnvelopeID)
	default:
		t.Fatal("expected response to be queued on socketModeResponses channel")
	}
}

func TestSendCtx_MarshalError(t *testing.T) {
	api := slack.New("ABCDEFG")
	cli := New(api)

	err := cli.SendCtx(context.Background(), Response{
		EnvelopeID: "test-envelope",
		Payload:    func() {},
	})

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "marshalling socket mode response")
}

func TestAck_ReturnsNoError(t *testing.T) {
	api := slack.New("ABCDEFG")
	cli := New(api)

	err := cli.Ack(Request{EnvelopeID: "test-envelope"}, "payload")
	assert.NoError(t, err)

	select {
	case res := <-cli.socketModeResponses:
		assert.Equal(t, "test-envelope", res.EnvelopeID)
	default:
		t.Fatal("expected response to be queued on socketModeResponses channel")
	}
}

func TestAck_ReturnsErrorWhenPayloadTooLarge(t *testing.T) {
	api := slack.New("ABCDEFG")
	cli := New(api)

	largePayload := strings.Repeat("x", maxResponseSize)
	err := cli.Ack(Request{EnvelopeID: "test-envelope"}, largePayload)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "silently dropped")
}
