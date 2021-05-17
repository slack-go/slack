package slack

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
)

const (
	DEFAULT_LOGINS_COUNT = 100
	DEFAULT_LOGINS_PAGE  = 1
	DEFAULT_LOGS_COUNT   = 100
	DEFAULT_LOGS_PAGE    = 1
)

type TeamResponse struct {
	Team TeamInfo `json:"team"`
	SlackResponse
}

type TeamInfo struct {
	ID          string                 `json:"id"`
	Name        string                 `json:"name"`
	Domain      string                 `json:"domain"`
	EmailDomain string                 `json:"email_domain"`
	Icon        map[string]interface{} `json:"icon"`
}

type LoginResponse struct {
	Logins []Login `json:"logins"`
	Paging `json:"paging"`
	SlackResponse
}

type Login struct {
	UserID    string `json:"user_id"`
	Username  string `json:"username"`
	DateFirst int    `json:"date_first"`
	DateLast  int    `json:"date_last"`
	Count     int    `json:"count"`
	IP        string `json:"ip"`
	UserAgent string `json:"user_agent"`
	ISP       string `json:"isp"`
	Country   string `json:"country"`
	Region    string `json:"region"`
}

type LogResponse struct {
	Logs   []Log `json:"logs"`
	Paging `json:"paging"`
	SlackResponse
}

// IntAsString exists so we can unmarshal both string and int as string due to apperant bug in slack
// https://github.com/slack-go/slack/pull/920#issuecomment-823655954
type IntAsString string

// UnmarshalJSON will unmarshal both string and int JSON values
func (i *IntAsString) UnmarshalJSON(buf []byte) error {
	var v interface{}
	if err := json.Unmarshal(buf, &v); err != nil {
		return err
	}
	switch v := v.(type) {
	case string:
		*i = IntAsString(v)
	case float64:
		*i = IntAsString(strconv.FormatInt(int64(v), 10))
	default:
		return fmt.Errorf("slack: unknown IntAsString type: %+v", v)
	}
	return nil
}

type Log struct {
	// according to example in https://api.slack.com/methods/team.integrationLogs `service_id` can be both int or string
	ServiceID   IntAsString `json:"service_id"`
	ServiceType string      `json:"service_type"`
	AppID       string      `json:"app_id"`
	AppType     string      `json:"app_type"`
	UserID      string      `json:"user_id"`
	UserName    string      `json:"user_name"`
	Channel     string      `json:"channel"`
	Date        string      `json:"date"`
	ChangeType  string      `json:"change_type"`
	Reason      string      `json:"reason"`
	Scope       string      `json:"scope"`
}

type BillableInfoResponse struct {
	BillableInfo map[string]BillingActive `json:"billable_info"`
	SlackResponse
}

type BillingActive struct {
	BillingActive bool `json:"billing_active"`
}

// AccessLogParameters contains all the parameters necessary (including the optional ones) for a GetAccessLogs() request
type AccessLogParameters struct {
	Count int
	Page  int
}

// IntegrationLogParameters contains all the parameters necessary (including the optional ones) for a GetIntegrationLogs() request
type IntegrationLogParameters struct {
	AppID      string
	ChangeType string
	ServiceID  string
	TeamID     string
	User       string
	Count      int
	Page       int
}

// NewAccessLogParameters provides an instance of AccessLogParameters with all the sane default values set
func NewAccessLogParameters() AccessLogParameters {
	return AccessLogParameters{
		Count: DEFAULT_LOGINS_COUNT,
		Page:  DEFAULT_LOGINS_PAGE,
	}
}

// NewIntegrationLogParameters provides an instance of IntegrationLogParameters with all the sane default values set
func NewIntegrationLogParameters() IntegrationLogParameters {
	return IntegrationLogParameters{
		Count: DEFAULT_LOGS_COUNT,
		Page:  DEFAULT_LOGS_PAGE,
	}
}

func (api *Client) teamRequest(ctx context.Context, path string, values url.Values) (*TeamResponse, error) {
	response := &TeamResponse{}
	err := api.postMethod(ctx, path, values, response)
	if err != nil {
		return nil, err
	}

	return response, response.Err()
}

func (api *Client) billableInfoRequest(ctx context.Context, path string, values url.Values) (map[string]BillingActive, error) {
	response := &BillableInfoResponse{}
	err := api.postMethod(ctx, path, values, response)
	if err != nil {
		return nil, err
	}

	return response.BillableInfo, response.Err()
}

func (api *Client) accessLogsRequest(ctx context.Context, path string, values url.Values) (*LoginResponse, error) {
	response := &LoginResponse{}
	err := api.postMethod(ctx, path, values, response)
	if err != nil {
		return nil, err
	}
	return response, response.Err()
}

func (api *Client) integrationLogsRequest(ctx context.Context, path string, values url.Values) (*LogResponse, error) {
	response := &LogResponse{}
	err := api.postMethod(ctx, path, values, response)
	if err != nil {
		return nil, err
	}
	return response, response.Err()
}

// GetTeamInfo gets the Team Information of the user
func (api *Client) GetTeamInfo() (*TeamInfo, error) {
	return api.GetTeamInfoContext(context.Background())
}

// GetTeamInfoContext gets the Team Information of the user with a custom context
func (api *Client) GetTeamInfoContext(ctx context.Context) (*TeamInfo, error) {
	values := url.Values{
		"token": {api.token},
	}

	response, err := api.teamRequest(ctx, "team.info", values)
	if err != nil {
		return nil, err
	}
	return &response.Team, nil
}

// GetAccessLogs retrieves a page of logins according to the parameters given
func (api *Client) GetAccessLogs(params AccessLogParameters) ([]Login, *Paging, error) {
	return api.GetAccessLogsContext(context.Background(), params)
}

// GetAccessLogsContext retrieves a page of logins according to the parameters given with a custom context
func (api *Client) GetAccessLogsContext(ctx context.Context, params AccessLogParameters) ([]Login, *Paging, error) {
	values := url.Values{
		"token": {api.token},
	}
	if params.Count != DEFAULT_LOGINS_COUNT {
		values.Add("count", strconv.Itoa(params.Count))
	}
	if params.Page != DEFAULT_LOGINS_PAGE {
		values.Add("page", strconv.Itoa(params.Page))
	}

	response, err := api.accessLogsRequest(ctx, "team.accessLogs", values)
	if err != nil {
		return nil, nil, err
	}
	return response.Logins, &response.Paging, nil
}

// GetIntegrationLogs retrieves a page of activity logs according to the parameters given
func (api *Client) GetIntegrationLogs(params IntegrationLogParameters) ([]Log, *Paging, error) {
	return api.GetIntegrationLogsContext(context.Background(), params)
}

// GetIntegrationLogsContext retrieves a page of activity logs according to the parameters given with a custom context
func (api *Client) GetIntegrationLogsContext(ctx context.Context, params IntegrationLogParameters) ([]Log, *Paging, error) {
	values := url.Values{
		"token": {api.token},
	}
	if params.AppID != "" {
		values.Add("app_id", params.AppID)
	}
	if params.ChangeType != "" {
		values.Add("change_type", params.ChangeType)
	}
	if params.ServiceID != "" {
		values.Add("service_id", params.ServiceID)
	}
	if params.TeamID != "" {
		values.Add("team_id", params.TeamID)
	}
	if params.User != "" {
		values.Add("user", params.User)
	}
	if params.Count != DEFAULT_LOGINS_COUNT {
		values.Add("count", strconv.Itoa(params.Count))
	}
	if params.Page != DEFAULT_LOGINS_PAGE {
		values.Add("page", strconv.Itoa(params.Page))
	}

	response, err := api.integrationLogsRequest(ctx, "team.integrationLogs", values)
	if err != nil {
		return nil, nil, err
	}
	return response.Logs, &response.Paging, nil
}

// GetBillableInfo ...
func (api *Client) GetBillableInfo(user string) (map[string]BillingActive, error) {
	return api.GetBillableInfoContext(context.Background(), user)
}

// GetBillableInfoContext ...
func (api *Client) GetBillableInfoContext(ctx context.Context, user string) (map[string]BillingActive, error) {
	values := url.Values{
		"token": {api.token},
		"user":  {user},
	}

	return api.billableInfoRequest(ctx, "team.billableInfo", values)
}

// GetBillableInfoForTeam returns the billing_active status of all users on the team.
func (api *Client) GetBillableInfoForTeam() (map[string]BillingActive, error) {
	return api.GetBillableInfoForTeamContext(context.Background())
}

// GetBillableInfoForTeamContext returns the billing_active status of all users on the team with a custom context
func (api *Client) GetBillableInfoForTeamContext(ctx context.Context) (map[string]BillingActive, error) {
	values := url.Values{
		"token": {api.token},
	}

	return api.billableInfoRequest(ctx, "team.billableInfo", values)
}
