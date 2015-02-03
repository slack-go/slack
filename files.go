package slack

import (
	"errors"
	"net/url"
	"strconv"
	"strings"
)

const (
	// Add here the defaults in the siten
	DEFAULT_FILES_USERID  = ""
	DEFAULT_FILES_TS_FROM = 0
	DEFAULT_FILES_TS_TO   = -1
	DEFAULT_FILES_TYPES   = "all"
	DEFAULT_FILES_COUNT   = 100
	DEFAULT_FILES_PAGE    = 1
)

// Comment contains all the information relative to a comment
type Comment struct {
	Id        string   `json:"id"`
	Timestamp JSONTime `json:"timestamp"`
	UserId    string   `json:"user"`
	Comment   string   `json:"comment"`
	Created   JSONTime `json:"created,omitempty"`
}

// File contains all the information for a file
type File struct {
	Id        string   `json:"id"`
	Created   JSONTime `json:"created"`
	Timestamp JSONTime `json:"timestamp"`

	Name       string `json:"name"`
	Title      string `json:"title"`
	Mimetype   string `json:"mimetype"`
	Filetype   string `json:"filetype"`
	PrettyType string `json:"pretty_type"`
	UserId     string `json:"user"`

	Mode         string `json:"mode"`
	Editable     bool   `json:"editable"`
	IsExternal   bool   `json:"is_external"`
	ExternalType string `json:"external_type"`

	Size int `json:"size"`

	URL                string `json:"url"`
	URLDownload        string `json:"url_download"`
	URLPrivate         string `json:"url_private"`
	URLPrivateDownload string `json:"url_private_download"`

	Thumb64     string `json:"thumb_64"`
	Thumb80     string `json:"thumb_80"`
	Thumb360    string `json:"thumb_360"`
	Thumb360Gif string `json:"thumb_360_gif"`
	Thumb360W   int    `json:"thumb_360_w"`
	Thumb360H   int    `json:"thumb_360_h"`

	Permalink        string `json:"permalink"`
	EditLink         string `json:"edit_link"`
	Preview          string `json:"preview"`
	PreviewHighlight string `json:"preview_highlight"`
	Lines            int    `json:"lines"`
	LinesMore        int    `json:"lines_more"`

	IsPublic        bool     `json:"is_public"`
	PublicURLShared bool     `json:"public_url_shared"`
	Channels        []string `json:"channels"`
	Groups          []string `json:"groups"`
	InitialComment  Comment  `json:"initial_comment"`
	NumStars        int      `json:"num_stars"`
	IsStarred       bool     `json:"is_starred"`
}

// FileUploadParameters contains all the parameters necessary (including the optional ones) for an UploadFile() request
type FileUploadParameters struct {
	File           string
	Content        string
	Filetype       string
	Filename       string
	Title          string
	InitialComment string
	Channels       []string
}

// GetFilesParameters contains all the parameters necessary (including the optional ones) for a GetFiles() request
type GetFilesParameters struct {
	UserId        string
	TimestampFrom JSONTime
	TimestampTo   JSONTime
	Types         string
	Count         int
	Page          int
}

type fileResponseFull struct {
	File     `json:"file"`
	Paging   `json:"paging"`
	Comments []Comment `json:"comments"`
	Files    []File    `json:"files"`

	SlackResponse
}

// NewGetFilesParameters provides an instance of GetFilesParameters with all the sane default values set
func NewGetFilesParameters() GetFilesParameters {
	return GetFilesParameters{
		UserId:        DEFAULT_FILES_USERID,
		TimestampFrom: DEFAULT_FILES_TS_FROM,
		TimestampTo:   DEFAULT_FILES_TS_TO,
		Types:         DEFAULT_FILES_TYPES,
		Count:         DEFAULT_FILES_COUNT,
		Page:          DEFAULT_FILES_PAGE,
	}
}

func fileRequest(path string, values url.Values, debug bool) (*fileResponseFull, error) {
	response := &fileResponseFull{}
	err := parseResponse(path, values, response, debug)
	if err != nil {
		return nil, err
	}
	if !response.Ok {
		return nil, errors.New(response.Error)
	}
	return response, nil
}

// GetFileInfo retrieves a file and related comments
func (api *Slack) GetFileInfo(fileId string, count, page int) (*File, []Comment, *Paging, error) {
	values := url.Values{
		"token": {api.config.token},
		"file":  {fileId},
		"count": {strconv.Itoa(count)},
		"page":  {strconv.Itoa(page)},
	}
	response, err := fileRequest("files.info", values, api.debug)
	if err != nil {
		return nil, nil, nil, err
	}
	return &response.File, response.Comments, &response.Paging, nil
}

// GetFiles retrieves all files according to the parameters given
func (api *Slack) GetFiles(params GetFilesParameters) ([]File, *Paging, error) {
	values := url.Values{
		"token": {api.config.token},
	}
	if params.UserId != DEFAULT_FILES_USERID {
		values.Add("user", params.UserId)
	}
	// XXX: this is broken. fix it with a proper unix timestamp
	if params.TimestampFrom != DEFAULT_FILES_TS_FROM {
		values.Add("ts_from", params.TimestampFrom.String())
	}
	if params.TimestampTo != DEFAULT_FILES_TS_TO {
		values.Add("ts_to", params.TimestampTo.String())
	}
	if params.Types != DEFAULT_FILES_TYPES {
		values.Add("types", params.Types)
	}
	if params.Count != DEFAULT_FILES_COUNT {
		values.Add("count", strconv.Itoa(params.Count))
	}
	if params.Page != DEFAULT_FILES_PAGE {
		values.Add("page", strconv.Itoa(params.Page))
	}
	response, err := fileRequest("files.list", values, api.debug)
	if err != nil {
		return nil, nil, err
	}
	return response.Files, &response.Paging, nil
}

// UploadFile uploads a file
func (api *Slack) UploadFile(params FileUploadParameters) (file *File, err error) {
	// Test if user token is valid. This helps because client.Do doesn't like this for some reason. XXX: More
	// investigation needed, but for now this will do.
	_, err = api.AuthTest()
	if err != nil {
		return nil, err
	}
	response := &fileResponseFull{}
	values := url.Values{
		"token": {api.config.token},
	}
	if params.Filetype != "" {
		values.Add("filetype", params.Filetype)
	}
	if params.Filename != "" {
		values.Add("filename", params.Filename)
	}
	if params.Title != "" {
		values.Add("title", params.Title)
	}
	if params.InitialComment != "" {
		values.Add("initial_comment", params.InitialComment)
	}
	if len(params.Channels) != 0 {
		values.Add("channels", strings.Join(params.Channels, ","))
	}
	if params.Content != "" {
		values.Add("content", params.Content)
		err = parseResponse("files.upload", values, response, api.debug)
	} else if params.File != "" {
		err = parseResponseMultipart("files.upload", params.File, values, response, api.debug)
	}
	if err != nil {
		return nil, err
	}
	if !response.Ok {
		return nil, errors.New(response.Error)
	}
	return &response.File, nil
}

// DeleteFile deletes a file
func (api *Slack) DeleteFile(fileId string) error {
	values := url.Values{
		"token": {api.config.token},
		"file":  {fileId},
	}
	_, err := fileRequest("files.delete", values, api.debug)
	if err != nil {
		return err
	}
	return nil

}
