Slack API in Go

## Installing

### *go get*

    $ go get github.com/nlopes/slack

## Example

### Getting all groups

    import (
		"fmt"

		"github.com/nlopes/slack"
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
			fmt.Printf("Id: %s, Name: %s\n", group.Id, group.Name)
		}
	}

### Getting User Information

    import (
	    "fmt"

	    "github.com/nlopes/slack"
    )

    func main() {
	    api := slack.New("YOUR_TOKEN_HERE")
	    user, err := api.GetUserInfo("U023BECGF")
	    if err != nil {
		    fmt.Printf("%s\n", err)
		    return
	    }
	    fmt.Printf("Id: %s, Fullname: %s, Email: %s\n", user.Id, user.Profile.RealName, user.Profile.Email)
    }

## Why?
I am currently learning Go and this seemed like a good idea.

## Stability
As with any other piece of software expect bugs. Also, the design isn't finalized yet because I am not happy with how I laid out some things. Especially the websocket stuff. It is functional but very incomplete and buggy.

## Help
Anyone is welcome to contribute. Either open a PR or create an issue.

## License
BSD 2 Clause license