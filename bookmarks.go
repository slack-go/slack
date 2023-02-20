package slack

import (
	"context"
	"errors"
	"net/url"
)

type Bookmark struct {
	ID                  string  `json:"id"`
	ChannelID           string  `json:"channel_id"`
	Title               string  `json:"title"`
	Link                string  `json:"link"`
	Emoji               string  `json:"emoji,omitempty"`
	IconURL             *string `json:"icon_url"`
	Type                string  `json:"type"`
	DateCreated         uint64  `json:"date_created"`
	DateUpdated         uint64  `json:"date_updated"`
	Rank                string  `json:"rank"`
	LastUpdatedByUserID *string `json:"last_updated_by_user_id"`
	LastUpdatedByTeamID *string `json:"last_updated_by_team_id"`
	ShortcutID          *string `json:"shortcut_id"`
	EntityID            *string `json:"entity_id"`
	AppID               *string `json:"app_id"`
}

// ListBookmarks returns all the bookmarks in the given channel
func (api *Client) ListBookmarks(channelID string) ([]Bookmark, error) {
	return api.ListBookmarksContext(context.Background(), channelID)
}

// ListBookmarksContext returns all the bookmarks in the given channel
func (api *Client) ListBookmarksContext(ctx context.Context, channelID string) ([]Bookmark, error) {
	values := url.Values{
		"token":      {api.token},
		"channel_id": {channelID},
	}

	response := &listBookmarksResponseFull{}
	err := api.postMethod(ctx, "bookmarks.list", values, response)
	if err != nil {
		return nil, err
	}
	if !response.Ok {
		return nil, errors.New(response.Error)
	}
	return response.Bookmarks, nil
}

type AddBookmarkParameters struct {
	Title     string `json:"title"`
	Type      string `json:"type"`
	Link      string `json:"link,omitempty"`
	Emoji     string `json:"emoji,omitempty"`
	EntityID  string `json:"entity_id,omitempty"`
	ParentID  string `json:"parent_id,omitempty"`
	ChannelID string `json:"channel_id"`
}

// AddBookmark creates a new bookmark. ChannelID, Title, and Type are required
// (`Type=link` is the sensible default!). The other params are all optional.
func (api *Client) AddBookmark(params AddBookmarkParameters) (*Bookmark, error) {
	return api.AddBookmarkContext(context.Background(), params)
}

// AddBookmarkContext creates a new bookmark. ChannelID, Title, and Type are required
// (`Type: "link"` is the sensible default!). The other params are all optional.
func (api *Client) AddBookmarkContext(ctx context.Context, params AddBookmarkParameters) (*Bookmark, error) {
	response := &singleBookmarkResponse{}
	values := url.Values{
		"token":      {api.token},
		"channel_id": {params.ChannelID},
		"title":      {params.Title},
		"type":       {params.Type},
	}

	if params.Emoji != "" {
		values["emoji"] = []string{params.Emoji}
	}

	if params.EntityID != "" {
		values["entity_id"] = []string{params.EntityID}
	}

	if params.Link != "" {
		values["link"] = []string{params.Link}
	}

	if params.ParentID != "" {
		values["parent_id"] = []string{params.ParentID}
	}

	err := api.postMethod(ctx, "bookmarks.add", values, response)
	if err != nil {
		return nil, err
	}
	if err := response.Err(); err != nil {
		return nil, err
	}

	return &response.Bookmark, nil
}

type EditBookmarkParameters struct {
	Title      string `json:"title,omitempty"`
	Emoji      string `json:"emoji,omitempty"`
	Link       string `json:"link,omitempty"`
	ChannelID  string `json:"channel_id"`
	BookmarkID string `json:"bookmark_id"`
	Type       string `json:"type,omitempty"`
}

// EditBookmark updates an existing bookmark. ChannelID and BookmarkID are
// required, other params are optional.
func (api *Client) EditBookmark(params EditBookmarkParameters) (*Bookmark, error) {
	return api.EditBookmarkContext(context.Background(), params)
}

// EditBookmarkContext updates an existing bookmark. ChannelID and BookmarkID
// are required, other params are optional.
func (api *Client) EditBookmarkContext(ctx context.Context, params EditBookmarkParameters) (*Bookmark, error) {
	response := &singleBookmarkResponse{}
	values := url.Values{
		"token":       {api.token},
		"channel_id":  {params.ChannelID},
		"bookmark_id": {params.BookmarkID},
	}

	if params.Type != "" {
		values["type"] = []string{params.Type}
	}

	if params.Emoji != "" {
		values["emoji"] = []string{params.Emoji}
	}

	if params.Link != "" {
		values["link"] = []string{params.Link}
	}

	if params.Title != "" {
		values["title"] = []string{params.Title}
	}

	err := api.postMethod(ctx, "bookmarks.edit", values, &response)

	if err != nil {
		return nil, err
	}
	if err := response.Err(); err != nil {
		return nil, err
	}

	return &response.Bookmark, nil
}

// RemoveBookmark deletes a bookmark from the given channel
func (api *Client) RemoveBookmark(channelID, bookmarkID string) error {
	return api.RemoveBookmarkContext(context.Background(), channelID, bookmarkID)
}

// RemoveBookmarkContext deletes a bookmark from the given channel
func (api *Client) RemoveBookmarkContext(ctx context.Context, channelID, bookmarkID string) error {
	response := &SlackResponse{}
	values := url.Values{
		"token":       {api.token},
		"channel_id":  {channelID},
		"bookmark_id": {bookmarkID},
	}

	err := api.postMethod(ctx, "bookmarks.remove", values, response)
	if err != nil {
		return err
	}
	return response.Err()
}

type listBookmarksResponseFull struct {
	Bookmarks []Bookmark
	SlackResponse
}

type singleBookmarkResponse struct {
	Bookmark Bookmark `json:"bookmark"`
	SlackResponse
}
