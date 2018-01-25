package slack

import (
	"context"
	"errors"
	"net/url"
)

// AddReminder adds a reminder for a channel or a user.
//
// Only one of ChannelID or UserID should be provided. Both can be omitted, in
// which case the reminder is set for the authenticated user.
//
// See https://api.slack.com/methods/reminders.add (NOTE: the ability to set
// reminders on a channel is currently undocumented but has been tested to
// work)
func (api *Client) AddReminder(channelID, userID, text, time string) error {
	return api.AddReminderContext(context.Background(), channelID, userID, text, time)
}

// AddReminderContext adds a reminder for a channel or a user with a custom context.
//
// See AddReminder for full details.
func (api *Client) AddReminderContext(ctx context.Context, channelID, userID, text, time string) error {
	values := url.Values{
		"token":   {api.token},
		"text":    {text},
		"time":    {time},
		"channel": {channelID},
	}
	if channelID != "" {
		values.Set("channel", channelID)
	} else if userID != "" {
		values.Set("user", userID)
	}

	response := &SlackResponse{}
	if err := post(ctx, api.httpclient, "reminders.add", values, response, api.debug); err != nil {
		return err
	}
	if !response.Ok {
		return errors.New(response.Error)
	}
	return nil
}
