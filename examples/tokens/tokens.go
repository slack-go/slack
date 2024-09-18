package tokens

import (
	"fmt"
	"github.com/slack-go/slack"
)

func main() {
	api := slack.New(
		"YOUR_TOKEN_HERE",
		// You may choose to provide your config tokens when creating your Slack client
		// or when invoking the method calls
		slack.OptionConfigToken("YOUR_CONFIG_ACCESS_TOKEN_HERE"),
		slack.OptionConfigRefreshToken("YOUR_REFRESH_TOKEN_HERE"),
	)

	// Obtain a fresh set of tokens
	// You may pass your tokens as a parameter here as well, if you didn't do it above
	freshTokens, err := api.RotateTokens("", "")
	if err != nil {
		fmt.Printf("error rotating tokens: %v\n", err)
		return
	}

	fmt.Printf("new access token: %s\n", freshTokens.Token)
	fmt.Printf("new refresh token: %s\n", freshTokens.RefreshToken)
	fmt.Printf("new tokenset expires at: %d\n", freshTokens.ExpiresAt)

	// Optionally: update the tokens inside the running Slack client
	// This isn't necessary if you restart the application after storing the tokens elsewhere,
	// or pass them as parameters to RotateTokens() explicitly
	api.UpdateConfigTokens(freshTokens)
}
