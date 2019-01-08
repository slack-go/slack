package slack

import (
	"context"
	"errors"
	"net/url"
)

func (api *Client) addReminder(ctx context.Context, values url.Values) error {
	response := &SlackResponse{}
	if err := postSlackMethod(ctx, api.httpclient, "reminders.add", values, response, api); err != nil {
		return err
	}
	if !response.Ok {
		return errors.New(response.Error)
	}
	return nil
}

// AddChannelReminder adds a reminder for a channel with a custom context.
//
// See https://api.slack.com/methods/reminders.add (NOTE: the ability to set
// reminders on a channel is currently undocumented but has been tested to
// work)
func (api *Client) AddChannelReminder(channelID, text, time string) error {
	values := url.Values{
		"token":   {api.token},
		"text":    {text},
		"time":    {time},
		"channel": {channelID},
	}
	return api.addReminder(context.Background(), values)
}

// AddUserReminder adds a reminder for a user with a custom context.
//
// See https://api.slack.com/methods/reminders.add (NOTE: the ability to set
// reminders on a channel is currently undocumented but has been tested to
// work)
func (api *Client) AddUserReminder(userID, text, time string) error {
	values := url.Values{
		"token": {api.token},
		"text":  {text},
		"time":  {time},
		"user":  {userID},
	}
	return api.addReminder(context.Background(), values)
}
