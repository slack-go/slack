package slack

import (
	"errors"
	"fmt"
	"net/url"
)

type adminResponse struct {
	OK    bool   `json:"ok"`
	Error string `json:"error"`
}

func adminRequest(method string, teamName string, values url.Values, debug bool) (*adminResponse, error) {
	adminResponse := &adminResponse{}
	err := parseAdminResponse(method, teamName, values, adminResponse, debug)
	if err != nil {
		return nil, err
	}

	if !adminResponse.OK {
		return nil, errors.New(adminResponse.Error)
	}

	return adminResponse, nil
}

// DisableUser disables a user account, given a user ID
func (api *Slack) DisableUser(teamName string, uid string) error {
	values := url.Values{
		"user":       {uid},
		"token":      {api.config.token},
		"set_active": {"true"},
		"_attempts":  {"1"},
	}

	_, err := adminRequest("setInactive", teamName, values, api.debug)
	if err != nil {
		return fmt.Errorf("Failed to disable user with id '%s': %s", uid, err)
	}

	return nil
}

// InviteGuest invites a user to Slack as a single-channel guest
func (api *Slack) InviteGuest(
	teamName string,
	channel string,
	firstName string,
	lastName string,
	emailAddress string,
) error {
	values := url.Values{
		"email":            {emailAddress},
		"channels":         {channel},
		"first_name":       {firstName},
		"last_name":        {lastName},
		"ultra_restricted": {"1"},
		"token":            {api.config.token},
		"set_active":       {"true"},
		"_attempts":        {"1"},
	}

	_, err := adminRequest("invite", teamName, values, api.debug)
	if err != nil {
		return fmt.Errorf("Failed to invite single-channel guest: %s", err)
	}

	return nil
}

// InviteRestricted invites a user to Slack as a restricted account
func (api *Slack) InviteRestricted(
	teamName string,
	channel string,
	firstName string,
	lastName string,
	emailAddress string,
) error {
	values := url.Values{
		"email":      {emailAddress},
		"channels":   {channel},
		"first_name": {firstName},
		"last_name":  {lastName},
		"restricted": {"1"},
		"token":      {api.config.token},
		"set_active": {"true"},
		"_attempts":  {"1"},
	}

	_, err := adminRequest("invite", teamName, values, api.debug)
	if err != nil {
		return fmt.Errorf("Failed to restricted account: %s", err)
	}

	return nil
}

// SetRegular enables the specified user
func (api *Slack) SetRegular(teamName string, user string) error {
	values := url.Values{
		"user":       {user},
		"token":      {api.config.token},
		"set_active": {"true"},
		"_attempts":  {"1"},
	}

	_, err := adminRequest("setRegular", teamName, values, api.debug)
	if err != nil {
		return fmt.Errorf("Failed to change the user (%s) to a regular user: %s", user, err)
	}

	return nil
}

// SendSSOBindingEmail sends an SSO binding email to the specified user
func (api *Slack) SendSSOBindingEmail(teamName string, user string) error {
	values := url.Values{
		"user":       {user},
		"token":      {api.config.token},
		"set_active": {"true"},
		"_attempts":  {"1"},
	}

	_, err := adminRequest("sendSSOBind", teamName, values, api.debug)
	if err != nil {
		return fmt.Errorf("Failed to send SSO binding email for user (%s): %s", user, err)
	}

	return nil
}

// SetUltraRestricted converts a user into a single-channel guest
func (api *Slack) SetUltraRestricted(
	teamName string,
	uid string,
	channel string,
) error {
	values := url.Values{
		"user":       {uid},
		"channel":    {channel},
		"token":      {api.config.token},
		"set_active": {"true"},
		"_attempts":  {"1"},
	}

	_, err := adminRequest("setUltraRestricted", teamName, values, api.debug)
	if err != nil {
		return fmt.Errorf("Failed to ultra-restrict account: %s", err)
	}

	return nil
}

// SetRestricted converts a user into a restricted account
func (api *Slack) SetRestricted(
	teamName string,
	uid string,
) error {
	values := url.Values{
		"user":       {uid},
		"token":      {api.config.token},
		"set_active": {"true"},
		"_attempts":  {"1"},
	}

	_, err := adminRequest("setRestricted", teamName, values, api.debug)
	if err != nil {
		return fmt.Errorf("Failed to restrict account: %s", err)
	}

	return nil
}
