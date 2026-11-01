package slack

import (
	"context"
	"errors"
	"net/url"
	"strconv"
	"time"
)

const (
	DEFAULT_STARS_USER = ""
)

type StarsParameters struct {
	User   string
	Cursor string
	Limit  int
	TeamID string
}

type StarredItem Item

type listResponseFull struct {
	Items []Item `json:"items"`
	SlackResponse
	ResponseMetadata `json:"response_metadata"`
}

// NewStarsParameters initialises StarsParameters with default values
func NewStarsParameters() StarsParameters {
	return StarsParameters{
		User: DEFAULT_STARS_USER,
	}
}

// AddStar stars an item in a channel.
// For more information see the AddStarContext documentation.
//
// Deprecated: Slack's stars APIs do not manage Later items. Stars are no longer
// visible to end-users, and Slack recommends retiring features that use them.
func (api *Client) AddStar(channel string, item ItemRef) error {
	return api.AddStarContext(context.Background(), channel, item)
}

// AddStarContext stars an item in a channel with a custom context.
//
// Slack API docs: https://docs.slack.dev/reference/methods/stars.add
//
// Deprecated: Slack's stars APIs do not manage Later items. Stars are no longer
// visible to end-users, and Slack recommends retiring features that use them.
func (api *Client) AddStarContext(ctx context.Context, channel string, item ItemRef) error {
	values := url.Values{
		"channel": {channel},
		"token":   {api.token},
	}
	if item.Timestamp != "" {
		values.Set("timestamp", item.Timestamp)
	}
	if item.File != "" {
		values.Set("file", item.File)
	}
	if item.Comment != "" {
		values.Set("file_comment", item.Comment)
	}

	response := &SlackResponse{}
	if err := api.postMethod(ctx, "stars.add", values, response); err != nil {
		return err
	}

	return response.Err()
}

// RemoveStar removes a starred item from a channel.
// For more information see the RemoveStarContext documentation.
//
// Deprecated: Slack's stars APIs do not manage Later items. Stars are no longer
// visible to end-users, and Slack recommends retiring features that use them.
func (api *Client) RemoveStar(channel string, item ItemRef) error {
	return api.RemoveStarContext(context.Background(), channel, item)
}

// RemoveStarContext removes a starred item from a channel with a custom context.
//
// Slack API docs: https://docs.slack.dev/reference/methods/stars.remove
//
// Deprecated: Slack's stars APIs do not manage Later items. Stars are no longer
// visible to end-users, and Slack recommends retiring features that use them.
func (api *Client) RemoveStarContext(ctx context.Context, channel string, item ItemRef) error {
	values := url.Values{
		"channel": {channel},
		"token":   {api.token},
	}
	if item.Timestamp != "" {
		values.Set("timestamp", item.Timestamp)
	}
	if item.File != "" {
		values.Set("file", item.File)
	}
	if item.Comment != "" {
		values.Set("file_comment", item.Comment)
	}

	response := &SlackResponse{}
	if err := api.postMethod(ctx, "stars.remove", values, response); err != nil {
		return err
	}

	return response.Err()
}

// ListStars returns information about the stars a user added.
// For more information see the ListStarsContext documentation.
//
// Deprecated: Slack's stars APIs do not manage Later items. Stars are no longer
// visible to end-users, and Slack recommends retiring features that use them.
func (api *Client) ListStars(params StarsParameters) ([]Item, string, error) {
	return api.ListStarsContext(context.Background(), params)
}

// ListStarsContext returns information about the stars a user added with a custom context.
//
// Slack API docs: https://docs.slack.dev/reference/methods/stars.list
//
// Deprecated: Slack's stars APIs do not manage Later items. Stars are no longer
// visible to end-users, and Slack recommends retiring features that use them.
func (api *Client) ListStarsContext(ctx context.Context, params StarsParameters) ([]Item, string, error) {
	values := url.Values{
		"token": {api.token},
	}
	if params.User != DEFAULT_STARS_USER {
		values.Add("user", params.User)
	}
	if params.Cursor != "" {
		values.Add("cursor", params.Cursor)
	}
	if params.Limit != 0 {
		values.Add("limit", strconv.Itoa(params.Limit))
	}
	if params.TeamID != "" {
		values.Add("team_id", params.TeamID)
	}

	response := &listResponseFull{}
	err := api.postMethod(ctx, "stars.list", values, response)
	if err != nil {
		return nil, "", err
	}

	if err := response.Err(); err != nil {
		return nil, "", err
	}

	return response.Items, response.ResponseMetadata.Cursor, nil
}

// GetStarred returns a list of StarredItem items.
//
// The user then has to iterate over them and figure out what they should
// be looking at according to what is in the Type:
//
//	for _, item := range items {
//		switch c.Type {
//		case "file_comment":
//			log.Println(c.Comment)
//		case "file":
//			...
//	}
//
// This function still exists to maintain backwards compatibility.
// I exposed it as returning []StarredItem, so it shall stay as StarredItem.
//
// Deprecated: Slack's stars APIs do not manage Later items. Stars are no longer
// visible to end-users, and Slack recommends retiring features that use them.
func (api *Client) GetStarred(params StarsParameters) ([]StarredItem, string, error) {
	return api.GetStarredContext(context.Background(), params)
}

// GetStarredContext returns a list of StarredItem items with a custom context
// For more details see GetStarred
//
// Deprecated: Slack's stars APIs do not manage Later items. Stars are no longer
// visible to end-users, and Slack recommends retiring features that use them.
func (api *Client) GetStarredContext(ctx context.Context, params StarsParameters) ([]StarredItem, string, error) {
	items, nextCursor, err := api.ListStarsContext(ctx, params)
	if err != nil {
		return nil, "", err
	}
	starredItems := make([]StarredItem, len(items))
	for i, item := range items {
		starredItems[i] = StarredItem(item)
	}
	return starredItems, nextCursor, nil
}

type listResponsePaginated struct {
	Items []Item `json:"items"`
	SlackResponse
	Metadata ResponseMetadata `json:"response_metadata"`
}

// StarredItemPagination allows for paginating over legacy starred items.
//
// Deprecated: Slack's stars APIs do not manage Later items. Stars are no longer
// visible to end-users, and Slack recommends retiring features that use them.
type StarredItemPagination struct {
	Items        []Item
	limit        int
	previousResp *ResponseMetadata
	c            *Client
}

// ListStarsOption configures [Client.ListStarsPaginated].
type ListStarsOption func(*StarredItemPagination)

// ListAllStars returns the complete list of legacy starred items.
//
// Deprecated: Slack's stars APIs do not manage Later items. Stars are no longer
// visible to end-users, and Slack recommends retiring features that use them.
func (api *Client) ListAllStars() ([]Item, error) {
	return api.ListAllStarsContext(context.Background())
}

// ListAllStarsContext returns all legacy starred items with a custom context.
//
// Deprecated: Slack's stars APIs do not manage Later items. Stars are no longer
// visible to end-users, and Slack recommends retiring features that use them.
func (api *Client) ListAllStarsContext(ctx context.Context) (results []Item, err error) {
	p := api.ListStarsPaginated()
	for err == nil {
		p, err = p.next(ctx)
		if err == nil {
			results = append(results, p.Items...)
		} else if rateLimitedError, ok := errors.AsType[*RateLimitedError](err); ok {
			select {
			case <-ctx.Done():
				err = ctx.Err()
			case <-time.After(rateLimitedError.RetryAfter):
				err = nil
			}
		}
	}

	return results, p.failure(err)
}

// ListStarsPaginated initializes pagination over legacy starred items.
//
// Deprecated: Slack's stars APIs do not manage Later items. Stars are no longer
// visible to end-users, and Slack recommends retiring features that use them.
func (api *Client) ListStarsPaginated(options ...ListStarsOption) StarredItemPagination {
	return newStarPagination(api, options...)
}

func newStarPagination(c *Client, options ...ListStarsOption) (sip StarredItemPagination) {
	sip = StarredItemPagination{
		c:     c,
		limit: 200, // per slack api documentation.
	}

	for _, opt := range options {
		opt(&sip)
	}

	return sip
}

// done checks if the pagination has completed
func (StarredItemPagination) done(err error) bool {
	return errors.Is(err, errPaginationComplete)
}

// done checks if pagination failed.
func (t StarredItemPagination) failure(err error) error {
	if t.done(err) {
		return nil
	}

	return err
}

// next gets the next list of starred items based on the cursor value
func (t StarredItemPagination) next(ctx context.Context) (_ StarredItemPagination, err error) {
	var resp listResponsePaginated

	if t.c == nil || (t.previousResp != nil && t.previousResp.Cursor == "") {
		return t, errPaginationComplete
	}

	previousResp := t.previousResp.initialize()

	values := url.Values{
		"limit":  {strconv.Itoa(t.limit)},
		"token":  {t.c.token},
		"cursor": {previousResp.Cursor},
	}

	if err = t.c.postMethod(ctx, "stars.list", values, &resp); err != nil {
		return t, err
	}
	if err = resp.Err(); err != nil {
		return t, err
	}

	t.previousResp = &resp.Metadata
	t.Items = resp.Items

	return t, nil
}
