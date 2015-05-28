Slack API in Go
===============

This library supports most if not all of the `api.slack.com` REST
calls, as well as the Real-Time Messaging protocol over websocket, in
a fully managed way.

This fork breaks many things `github.com/nlopes/slack`, and improves
the RTM a lot.


## Installing

### *go get*

    $ go get github.com/abourget/slack

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

## Minimal RTM usage:




## Contributing

You are more than welcome to contribute to this project.  Fork and
make a Pull Request, or create an Issue if you see any problem.

## License

BSD 2 Clause license
