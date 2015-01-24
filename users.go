package slack

import (
	"errors"
	"net/url"
)

type UserProfile struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	RealName  string `json:"real_name"`
	Email     string `json:"email"`
	Skype     string `json:"skype"`
	Phone     string `json:"phone"`
	Image24   string `json:"image_24"`
	Image32   string `json:"image_32"`
	Image48   string `json:"image_48"`
	Image72   string `json:"image_72"`
	Image192  string `json:"image_192"`
}

type User struct {
	Id                string      `json:"id"`
	Name              string      `json:"name"`
	Deleted           bool        `json:"deleted"`
	Color             string      `json:"color"`
	Profile           UserProfile `json:"profile"`
	IsAdmin           bool        `json:"is_admin"`
	IsOwner           bool        `json:"is_owner"`
	IsPrimaryOwner    bool        `json:"is_primary_owner"`
	IsRestricted      bool        `json:"is_restricted"`
	IsUltraRestricted bool        `json:"is_ultra_restricted"`
	HasFiles          bool        `json:"has_files"`
}

type UserPresence struct {
	Presence        string   `json:"presence,omitempty"`
	Online          bool     `json:"online,omitempty"`
	AutoAway        bool     `json:"auto_away,omitempty"`
	ManualAway      bool     `json:"manual_away,omitempty"`
	ConnectionCount int      `json:"connection_count,omitempty"`
	LastActivity    JSONTime `json:"last_activity,omitempty"`
}

type userResponseFull struct {
	Members      []User                  `json:"members,omitempty"` // ListUsers
	User         `json:"user,omitempty"` // GetUserInfo
	UserPresence                         // GetUserPresence
	SlackResponse
}

func userRequest(path string, values url.Values, debug bool) (*userResponseFull, error) {
	response := &userResponseFull{}
	err := ParseResponse(path, values, response, debug)
	if err != nil {
		return nil, err
	}
	if !response.Ok {
		return nil, errors.New(response.Error)
	}
	return response, nil
}

func (api *Slack) GetUserPresence(userId string) (*UserPresence, error) {
	values := url.Values{
		"token": {api.config.token},
		"user":  {userId},
	}
	response, err := userRequest("users.getPresence", values, api.debug)
	if err != nil {
		return nil, err
	}
	return &response.UserPresence, nil
}

func (api *Slack) GetUserInfo(userId string) (*User, error) {
	values := url.Values{
		"token": {api.config.token},
		"user":  {userId},
	}
	response, err := userRequest("users.info", values, api.debug)
	if err != nil {
		return nil, err
	}
	return &response.User, nil
}

func (api *Slack) GetUsers() ([]User, error) {
	values := url.Values{
		"token": {api.config.token},
	}
	response, err := userRequest("users.list", values, api.debug)
	if err != nil {
		return nil, err
	}
	return response.Members, nil
}

func (api *Slack) SetUserAsActive() error {
	values := url.Values{
		"token": {api.config.token},
	}
	_, err := userRequest("users.setActive", values, api.debug)
	if err != nil {
		return err
	}
	return nil
}

func (api *Slack) SetUserPresence(presence string) error {
	values := url.Values{
		"token":    {api.config.token},
		"presence": {presence},
	}
	_, err := userRequest("users.setPresence", values, api.debug)
	if err != nil {
		return err
	}
	return nil

}
