package slacktest

import (
	"github.com/slack-go/slack/internal/errorsx"
)

const (
	// ErrEmptyServerToHub is the error when attempting an empty server address to the hub
	ErrEmptyServerToHub = errorsx.String("Unable to add an empty server address to hub")
	// ErrPassedEmptyServerAddr is the error when being passed an empty server address
	ErrPassedEmptyServerAddr = errorsx.String("Passed an empty server address")
	// ErrNoQueuesRegisteredForServer is the error when there are no queues for a server in the hub
	ErrNoQueuesRegisteredForServer = errorsx.String("No queues registered for server")
)
