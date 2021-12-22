package slack

import (
	"context"
	"encoding/json"
	"net/url"
)

// SetUserProfile will set a profile for the provided user
func (api *Client) SetUserProfile(user string, profile *UserProfile) error {
	return api.SetUserProfileContext(context.Background(), user, profile)
}

// SetUserProfileContext will set a profile for the provided user with a custom context
func (api *Client) SetUserProfileContext(ctx context.Context, user string, profile *UserProfile) error {
	jsonProfile, err := json.Marshal(profile)

	if err != nil {
		return err
	}

	values := url.Values{
		"token":   {api.token},
		"profile": {string(jsonProfile)},
	}

	// optional field. It should not be set if empty
	if user != "" {
		values["user"] = []string{user}
	}

	response := &userResponseFull{}
	if err = api.postMethod(ctx, "users.profile.set", values, response); err != nil {
		return err
	}

	return response.Err()
}
