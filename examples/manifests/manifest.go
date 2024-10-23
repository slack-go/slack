package manifests

import (
	"fmt"
	"github.com/slack-go/slack"
)

// createManifest programmatically creates a Slack app manifest
func createManifest() *slack.Manifest {
	return &slack.Manifest{
		Display: slack.Display{
			Name: "Your Application",
		},
		// ... other configuration here
	}
}

func main() {
	api := slack.New(
		"YOUR_TOKEN_HERE",
		// You may choose to provide your access token when creating your Slack client
		// or when invoking the method calls
		slack.OptionConfigToken("YOUR_CONFIG_ACCESS_TOKEN_HERE"),
	)

	// Create a new Manifest object
	manifest := createManifest()

	// Update your application using the new manifest
	// You may pass your token as a parameter here as well, if you didn't do it above
	response, err := api.UpdateManifest(manifest, "", "YOUR_APP_ID_HERE")
	if err != nil {
		fmt.Printf("error updating Slack application: %v\n", err)
		return
	}

	if !response.Ok {
		fmt.Printf("unable to update Slack application: %v\n", response.Errors)
	}

	fmt.Println("successfully updated Slack application")

	// The access token is now invalid, so it should be rotated for future use
	// Refer to the examples about tokens for more details
}
