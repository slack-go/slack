package slack

import (
	"errors"
	"net/url"
)


// User contains all the information of a user group
type UserGroup struct {
	ID                string      `json:"id"`
  Name              string      `json:"name"`
  TeamId            string      `json:"team_id"`
  IsUserGroup       bool        `json:"is_usergroup"`
  Description       string      `json:"description"`
  Handle            string      `json:"handle"`
  IsExternal        bool        `json:"is_external"`
  AutoType          string      `json:"auto_type"`
  CreatedBy         string      `json:"created_by"`
  UpdatedBy         string      `json:"updated_by"`
  DeletedBy         string      `json:"deleted_by"`
  Prefs             struct{
    Channels        []string    `json:"channels"`
    Groups          []string    `json:"groups"`
  }                             `json:"prefs"`
  Users             []string
  UserCount         int         `json:"user_count"`
}

type userGroupResponseFull struct {
  UserGroups   []UserGroup             `json:"usergroups"`
  UserGroup    UserGroup               `json:"usergroup"`
	SlackResponse
}

func userGroupRequest(path string, values url.Values, debug bool) (*userGroupResponseFull, error) {
	response := &userGroupResponseFull{}
	err := post(path, values, response, debug)
	if err != nil {
		return nil, err
	}
	if !response.Ok {
		return nil, errors.New(response.Error)
	}
	return response, nil
}

// GetUserGroups returns a list of user groups for the team
func (api *Client) GetUserGroups() ([]UserGroup, error){
  values := url.Values{
		"token": {api.config.token},
	}
	response, err := userGroupRequest("usergroups.list", values, api.debug)
	if err != nil {
		return nil, err
	}
	return response.UserGroups, nil
}

func (api *Client) UpdateUserGroup(userGroup UserGroup) (UserGroup, error){
  values := url.Values{
		"token": {api.config.token},
    "usergroup": {userGroup.ID},
	}

  if userGroup.Name != ""{
    values["name"] = []string{userGroup.Name}
  }

  if userGroup.Handle != ""{
    values["handle"] = []string{userGroup.Handle}
  }

  if userGroup.Description != ""{
    values["description"] = []string{userGroup.Description}
  }

	response, err := userGroupRequest("usergroups.update", values, api.debug)
	if err != nil {
		return UserGroup{}, err
	}
	return response.UserGroup, nil
}

func (api *Client) UpdateUserGroupMembers(userGroup string, members string) (UserGroup, error){
  values := url.Values{
		"token": {api.config.token},
    "usergroup": {userGroup.ID},
    "users": {members}
	}

	response, err := userGroupRequest("usergroups.users.update", values, api.debug)
	if err != nil {
		return UserGroup{}, err
	}
	return response.UserGroup, nil
}
