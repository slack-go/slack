package slack

import (
	"context"
	"encoding/json"
	"net/url"
	"strconv"
)

type Call struct {
	ID                string     `json:"id"`
	DateStart         JSONTime   `json:"date_start"`
	ExternalUniqueID  string     `json:"external_unique_id"`
	JoinUrl           string     `json:"join_url"`
	DesktopAppJoinUrl string     `json:"desktop_app_join_url"`
	ExternalDisplayID string     `json:"external_display_id"`
	Users             []CallUser `json:"users"`
}

// A thin user representation which has a SlackID, ExternalID, or both.
//
// See: https://api.slack.com/apis/calls#users
type CallUser struct {
	SlackID     string `json:"slack_id,omitempty"`
	ExternalID  string `json:"external_id,omitempty"`
	DisplayName string `json:"display_name,omitempty"`
	AvatarUrl   string `json:"avatar_url,omitempty"`
}

func (u CallUser) Valid() bool {
	return u.SlackID != "" || u.ExternalID != ""
}

type AddCallParameters struct {
	Title             string
	DesktopAppJoinUrl string
	ExternalDisplayID string
	DateStart         JSONTime
	CreatedBy         string
	Users             []CallUser
}

type UpdateCallParameters struct {
	Title             string
	DesktopAppJoinUrl string
	JoinUrl           string
}

type callResponse struct {
	Call Call `json:"call"`
	SlackResponse
}

func (api *Client) AddCall(externalID, joinUrl string, params AddCallParameters) (Call, error) {
	return api.AddCallContext(context.Background(), externalID, joinUrl, params)
}

func (api *Client) AddCallContext(ctx context.Context, externalID, joinUrl string, params AddCallParameters) (Call, error) {
	values := url.Values{
		"token":              {api.token},
		"join_url":           {joinUrl},
		"external_unique_id": {externalID},
	}
	if params.CreatedBy != "" {
		values.Set("created_by", params.CreatedBy)
	}
	if params.DateStart != 0 {
		values.Set("date_start", strconv.FormatInt(int64(params.DateStart), 10))
	}
	if params.DesktopAppJoinUrl != "" {
		values.Set("desktop_app_join_url", params.DesktopAppJoinUrl)
	}
	if params.ExternalDisplayID != "" {
		values.Set("external_display_id", params.ExternalDisplayID)
	}
	if params.Title != "" {
		values.Set("title", params.Title)
	}
	if len(params.Users) > 0 {
		data, err := json.Marshal(params.Users)
		if err != nil {
			return Call{}, err
		}
		values.Set("users", string(data))
	}

	response := &callResponse{}
	if err := api.postMethod(ctx, "calls.add", values, response); err != nil {
		return Call{}, err
	}

	return response.Call, response.Err()
}

func (api *Client) GetCallInfo(callID string) (Call, error) {
	return api.GetCallInfoContext(context.Background(), callID)
}

func (api *Client) GetCallInfoContext(ctx context.Context, callID string) (Call, error) {
	values := url.Values{
		"token": {api.token},
		"id":    {callID},
	}

	response := &callResponse{}
	if err := api.postMethod(ctx, "calls.add", values, response); err != nil {
		return Call{}, err
	}
	return response.Call, response.Err()
}

func (api *Client) UpdateCall(callID string, params UpdateCallParameters) (Call, error) {
	return api.GetCallInfoContext(context.Background(), callID)
}

func (api *Client) UpdateCallContext(ctx context.Context, callID string, params UpdateCallParameters) (Call, error) {
	values := url.Values{
		"token": {api.token},
		"id":    {callID},
	}

	if params.DesktopAppJoinUrl != "" {
		values.Set("desktop_app_join_url", params.DesktopAppJoinUrl)
	}
	if params.JoinUrl != "" {
		values.Set("join_url", params.JoinUrl)
	}
	if params.Title != "" {
		values.Set("title", params.Title)
	}

	response := &callResponse{}
	if err := api.postMethod(ctx, "calls.update", values, response); err != nil {
		return Call{}, err
	}
	return response.Call, response.Err()
}

func (api *Client) EndCall(callID string) error {
	return api.EndCallContext(context.Background(), callID)
}

func (api *Client) EndCallContext(ctx context.Context, callID string) error {
	values := url.Values{
		"token": {api.token},
		"id":    {callID},
	}

	response := &SlackResponse{}
	if err := api.postMethod(ctx, "calls.end", values, response); err != nil {
		return err
	}
	return response.Err()
}

func (api *Client) CallAddUsers(callID string, users []CallUser) error {
	return api.CallAddUsersContext(context.Background(), callID, users)
}

func (api *Client) CallAddUsersContext(ctx context.Context, callID string, users []CallUser) error {
	return api.setCallUsers(ctx, "calls.participants.add", callID, users)
}

func (api *Client) CallRemoveUsers(callID string, users []CallUser) error {
	return api.CallRemoveUsersContext(context.Background(), callID, users)
}

func (api *Client) CallRemoveUsersContext(ctx context.Context, callID string, users []CallUser) error {
	return api.setCallUsers(ctx, "calls.participants.remove", callID, users)
}

func (api *Client) setCallUsers(ctx context.Context, method, callID string, users []CallUser) error {
	values := url.Values{
		"token": {api.token},
		"id":    {callID},
	}

	data, err := json.Marshal(users)
	if err != nil {
		return err
	}
	values.Set("users", string(data))

	response := &SlackResponse{}
	if err := api.postMethod(ctx, method, values, response); err != nil {
		return err
	}
	return response.Err()
}
