//go:build go1.13
// +build go1.13

package socketmode

import (
	"context"
	"errors"
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
