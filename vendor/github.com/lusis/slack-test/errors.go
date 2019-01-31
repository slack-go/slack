package slacktest

import (
	"fmt"
)

// ErrEmptyServerToHub is the error when attempting an empty server address to the hub
var ErrEmptyServerToHub = fmt.Errorf("Unable to add an empty server address to hub")

// ErrPassedEmptyServerAddr is the error when being passed an empty server address
var ErrPassedEmptyServerAddr = fmt.Errorf("Passed an empty server address")

// ErrNoQueuesRegisteredForServer is the error when there are no queues for a server in the hub
var ErrNoQueuesRegisteredForServer = fmt.Errorf("No queues registered for server")
