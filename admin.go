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

func (api *Slack) InviteGuest(
	teamName string,
	channelID string,
	firstName string,
	lastName string,
	emailAddress string,
) error {
	values := url.Values{
		"email":            {emailAddress},
		"channels":         {channelID},
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

func (api *Slack) InviteRestricted(
	teamName string,
	channelID string,
	firstName string,
	lastName string,
	emailAddress string,
) error {
	values := url.Values{
		"email":      {emailAddress},
		"channels":   {channelID},
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
