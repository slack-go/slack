package slack

import (
	"context"
	"errors"
	"net/url"
)

func (api *Client) doReminder(ctx context.Context, path string, values url.Values) error {
	response := &SlackResponse{}
	if err := postSlackMethod(ctx, api.httpclient, path, values, response, api); err != nil {
		return err
	}
	if !response.Ok {
		return errors.New(response.Error)
	}
	return nil
}

// AddChannelReminder adds a reminder for a channel.
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
	return api.doReminder(context.Background(), "reminders.add", values)
}

// AddUserReminder adds a reminder for a user.
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
	return api.doReminder(context.Background(), "reminders.add", values)
}

// DeleteReminder deletes an existing reminder.
//
// See https://api.slack.com/methods/reminders.delete
func (api *Client) DeleteReminder(name string) error {
	values := url.Values{
		"token":    {api.token},
		"reminder": {name},
	}
	return api.doReminder(context.Background(), "reminders.delete", values)
}
