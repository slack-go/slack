package slack

import (
	"context"
	"io"
	"net/url"
	"strconv"
	"strings"
)

const (
	DEFAULT_REMOTE_FILES_CHANNEL = ""
	DEFAULT_REMOTE_FILES_TS_FROM = 0
	DEFAULT_REMOTE_FILES_TS_TO   = -1
	DEFAULT_REMOTE_FILES_COUNT   = 100
)

// RemoteFile contains all the information for a remote file
type RemoteFile struct {
	ID              string   `json:"id"`
	Created         JSONTime `json:"created"`
	Timestamp       JSONTime `json:"timestamp"`
	Name            string   `json:"name"`
	Title           string   `json:"title"`
	Mimetype        string   `json:"mimetype"`
	Filetype        string   `json:"filetype"`
	PrettyType      string   `json:"pretty_type"`
	User            string   `json:"user"`
	Editable        bool     `json:"editable"`
	Size            int      `json:"size"`
	Mode            string   `json:"mode"`
	IsExternal      bool     `json:"is_external"`
	ExternalType    string   `json:"external_type"`
	IsPublic        bool     `json:"is_public"`
	PublicURLShared bool     `json:"public_url_shared"`
	DisplayAsBot    bool     `json:"display_as_bot"`
	Username        string   `json:"username"`
	URLPrivate      string   `json:"url_private"`
	Permalink       string   `json:"permalink"`
	CommentsCount   int      `json:"comments_count"`
	IsStarred       bool     `json:"is_starred"`
	Shares          Share    `json:"shares"`
	Channels        []string `json:"channels"`
	Groups          []string `json:"groups"`
	IMs             []string `json:"ims"`
	ExternalID      string   `json:"external_id"`
	ExternalURL     string   `json:"external_url"`
	HasRichPreview  bool     `json:"has_rich_preview"`
}

type RemoteFileParameters struct {
	ExternalID            string // required
	ExternalURL           string // required
	Title                 string // required
	Filetype              string
	IndexableFileContents string
	PreviewImage          string
	PreviewImageReader    io.Reader
}

type ListRemoteFilesParameters struct {
	Channel       string
	Cursor        string
	Limit         int
	TimestampFrom JSONTime
	TimestampTo   JSONTime
}

type remoteFileResponseFull struct {
	RemoteFile `json:"file"`
	Paging     `json:"paging"`
	Files      []RemoteFile     `json:"files"`
	Metadata   ResponseMetadata `json:"response_metadata"`
	SlackResponse
}

func (api *Client) remoteFileRequest(ctx context.Context, path string, values url.Values) (*remoteFileResponseFull, error) {
	response := &remoteFileResponseFull{}
	err := api.postMethod(ctx, path, values, response)
	if err != nil {
		return nil, err
	}

	return response, response.Err()
}

// AddRemoteFile adds a remote file
func (api *Client) AddRemoteFile(params RemoteFileParameters) (remotefile *RemoteFile, err error) {
	return api.AddRemoteFileContext(context.Background(), params)
}

// AddRemoteFileContext adds a remote file and setting a custom context
func (api *Client) AddRemoteFileContext(ctx context.Context, params RemoteFileParameters) (remotefile *RemoteFile, err error) {
	// Test if user token is valid. This helps because client.Do doesn't like this for some reason. XXX: More
	// investigation needed, but for now this will do.
	_, err = api.AuthTest()
	if err != nil {
		return nil, err
	}
	response := &remoteFileResponseFull{}
	values := url.Values{
		"token":        {api.token},
		"external_id":  {params.ExternalID},
		"external_url": {params.ExternalURL},
		"title":        {params.Title},
	}
	if params.Filetype != "" {
		values.Add("filetype", params.Filetype)
	}
	if params.IndexableFileContents != "" {
		values.Add("indexable_file_contents", params.IndexableFileContents)
	}
	if params.PreviewImage != "" {
		err = postLocalWithMultipartResponse(ctx, api.httpclient, api.endpoint+"files.remote.add", params.PreviewImage, "preview_image", values, response, api)
	} else if params.PreviewImageReader != nil {
		err = postWithMultipartResponse(ctx, api.httpclient, api.endpoint+"files.remote.add", "preview.png", "preview_image", values, params.PreviewImageReader, response, api)
	} else {
		response, err = api.remoteFileRequest(ctx, "files.remote.add", values)
	}

	if err != nil {
		return nil, err
	}

	return &response.RemoteFile, response.Err()
}

// ListRemoteFiles retrieves all remote files according to the parameters given. Uses cursor based pagination.
func (api *Client) ListRemoteFiles(params ListRemoteFilesParameters) ([]RemoteFile, *ListRemoteFilesParameters, error) {
	return api.ListRemoteFilesContext(context.Background(), params)
}

// ListRemoteFilesContext retrieves all remote files according to the parameters given with a custom context. Uses cursor based pagination.
func (api *Client) ListRemoteFilesContext(ctx context.Context, params ListRemoteFilesParameters) ([]RemoteFile, *ListRemoteFilesParameters, error) {
	values := url.Values{
		"token": {api.token},
	}
	if params.Channel != DEFAULT_REMOTE_FILES_CHANNEL {
		values.Add("channel", params.Channel)
	}
	if params.TimestampFrom != DEFAULT_REMOTE_FILES_TS_FROM {
		values.Add("ts_from", strconv.FormatInt(int64(params.TimestampFrom), 10))
	}
	if params.TimestampTo != DEFAULT_REMOTE_FILES_TS_TO {
		values.Add("ts_to", strconv.FormatInt(int64(params.TimestampTo), 10))
	}
	if params.Limit != DEFAULT_REMOTE_FILES_COUNT {
		values.Add("limit", strconv.Itoa(params.Limit))
	}
	if params.Cursor != "" {
		values.Add("cursor", params.Cursor)
	}

	response, err := api.remoteFileRequest(ctx, "files.remote.list", values)
	if err != nil {
		return nil, nil, err
	}

	params.Cursor = response.Metadata.Cursor

	return response.Files, &params, nil
}

// GetRemoteFileInfo retrieves the complete remote file information.
func (api *Client) GetRemoteFileInfo(externalID, fileID string) (remotefile *RemoteFile, err error) {
	return api.GetRemoteFileInfoContext(context.Background(), externalID, fileID)
}

// GetRemoteFileInfoContext retrieves the complete remote file information given with a custom context.
func (api *Client) GetRemoteFileInfoContext(ctx context.Context, externalID, fileID string) (remotefile *RemoteFile, err error) {
	values := url.Values{
		"token": {api.token},
	}
	if fileID != "" {
		values.Add("file", fileID)
	}
	if externalID != "" {
		values.Add("external_id", externalID)
	}
	response, err := api.remoteFileRequest(ctx, "files.remote.info", values)
	if err != nil {
		return nil, err
	}
	return &response.RemoteFile, err
}

// ShareRemoteFile shares a remote file to channels
func (api *Client) ShareRemoteFile(channels []string, externalID, fileID string) (file *RemoteFile, err error) {
	return api.ShareRemoteFileContext(context.Background(), channels, externalID, fileID)
}

// ShareRemoteFileContext shares a remote file to channels with a custom context
func (api *Client) ShareRemoteFileContext(ctx context.Context, channels []string, externalID, fileID string) (file *RemoteFile, err error) {
	values := url.Values{
		"token":    {api.token},
		"channels": {strings.Join(channels, ",")},
	}
	if fileID != "" {
		values.Add("file", fileID)
	}
	if externalID != "" {
		values.Add("external_id", externalID)
	}
	response, err := api.remoteFileRequest(ctx, "files.remote.share", values)
	if err != nil {
		return nil, err
	}
	return &response.RemoteFile, err
}

// UpdateRemoteFile updates a remote file
func (api *Client) UpdateRemoteFile(fileID string, params RemoteFileParameters) (remotefile *RemoteFile, err error) {
	return api.UpdateRemoteFileContext(context.Background(), fileID, params)
}

// UpdateRemoteFileContext updates a remote file with a custom context
func (api *Client) UpdateRemoteFileContext(ctx context.Context, fileID string, params RemoteFileParameters) (remotefile *RemoteFile, err error) {
	// Test if user token is valid. This helps because client.Do doesn't like this for some reason. XXX: More
	// investigation needed, but for now this will do.
	_, err = api.AuthTest()
	if err != nil {
		return nil, err
	}
	response := &remoteFileResponseFull{}
	values := url.Values{
		"token": {api.token},
	}
	if fileID != "" {
		values.Add("file", fileID)
	}
	if params.ExternalID != "" {
		values.Add("external_id", params.ExternalID)
	}
	if params.ExternalURL != "" {
		values.Add("external_url", params.ExternalURL)
	}
	if params.Title != "" {
		values.Add("title", params.Title)
	}
	if params.Filetype != "" {
		values.Add("filetype", params.Filetype)
	}
	if params.IndexableFileContents != "" {
		values.Add("indexable_file_contents", params.IndexableFileContents)
	}
	if params.PreviewImageReader != nil {
		err = postWithMultipartResponse(ctx, api.httpclient, api.endpoint+"files.remote.add", "preview.png", "preview_image", values, params.PreviewImageReader, response, api)
	} else {
		response, err = api.remoteFileRequest(ctx, "files.remote.add", values)
	}

	if err != nil {
		return nil, err
	}

	return &response.RemoteFile, response.Err()
}

// RemoveRemoteFile removes a remote file
func (api *Client) RemoveRemoteFile(externalID, fileID string) (err error) {
	return api.RemoveRemoteFileContext(context.Background(), externalID, fileID)
}

// RemoveRemoteFileContext removes a remote file with a custom context
func (api *Client) RemoveRemoteFileContext(ctx context.Context, externalID, fileID string) (err error) {
	values := url.Values{
		"token": {api.token},
	}
	if fileID != "" {
		values.Add("file", fileID)
	}
	if externalID != "" {
		values.Add("external_id", externalID)
	}
	_, err = api.remoteFileRequest(ctx, "files.remote.remove", values)
	return err
}
