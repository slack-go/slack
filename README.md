Slack API in Go [![GoDoc](https://godoc.org/github.com/essentialkaos/slack?status.svg)](https://godoc.org/github.com/essentialkaos/slack) [![Go Report Card](https://goreportcard.com/badge/github.com/essentialkaos/slack)](https://goreportcard.com/report/github.com/essentialkaos/slack) [![Build Status](https://travis-ci.org/essentialkaos/slack.svg)](https://travis-ci.org/essentialkaos/slack)
===============

This library supports most if not all of the `api.slack.com` REST calls, as well as the Real-Time Messaging protocol over websocket, in a fully managed way.

### Installing

Make sure you have a working Go 1.7+ workspace ([instructions](https://golang.org/doc/install)), then:

```
go get github.com/essentialkaos/slack
```

For update to latest stable release, do:

```
go get -u github.com/essentialkaos/slack
```

### Examples

#### Getting all groups

```golang
import (
  "fmt"

  "github.com/essentialkaos/slack"
)

func main() {
  api := slack.New("YOUR_TOKEN_HERE")
  // If you set debugging, it will log all requests to the console
  // Useful when encountering issues
  // api.SetDebug(true)
  groups, err := api.GetGroups(false)

  if err != nil {
    fmt.Printf("%s\n", err)
    return
  }

  for _, group := range groups {
    fmt.Printf("ID: %s, Name: %s\n", group.ID, group.Name)
  }
}
```

#### Getting user information

```golang
import (
  "fmt"

  "github.com/essentialkaos/slack"
)

func main() {
  api := slack.New("YOUR_TOKEN_HERE")
  user, err := api.GetUserInfo("U023BECGF")

  if err != nil {
    fmt.Printf("%s\n", err)
    return
  }

  fmt.Printf("ID: %s, Fullname: %s, Email: %s\n", user.ID, user.Profile.RealName, user.Profile.Email)
}
```

#### Minimal RTM usage:

See [example](examples/websocket/websocket.go).

### Build Status

| Branch | Status |
|------------|--------|
| `master` (_Stable_) | [![Build Status](https://travis-ci.org/essentialkaos/slack.svg?branch=master)](https://travis-ci.org/essentialkaos/slack) |
| `develop` (_Unstable_) | [![Build Status](https://travis-ci.org/essentialkaos/slack.svg?branch=develop)](https://travis-ci.org/essentialkaos/slack) |

### License

BSD 2 Clause license
