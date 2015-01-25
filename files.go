package slack

import (
	"errors"
	"log"
	"net/url"
	"strconv"
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

type Comment struct {
	Id        string   `json:"id"`
	Timestamp JSONTime `json:"timestamp"`
	UserId    string   `json:"user"`
	Comment   string   `json:"comment"`
}

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

	Url                string `json:"url"`
	UrlDownload        string `json:"url_download"`
	UrlPrivate         string `json:"url_private"`
	UrlPrivateDownload string `json:"url_private_download"`

	Thumb64     string `json:"thumb_64"`
	Thumb80     string `json:"thumb_80"`
	Thumb360    string `json:"thumb_360"`
	Thumb360Gif string `json:"thumb_360_gif"`
	Thumb360W   string `json:"thumb_360_w"`
	Thumb360H   string `json:"thumb_360_h"`

	Permalink        string `json:"permalink"`
	EditLink         string `json:"edit_link"`
	Preview          string `json:"preview"`
	PreviewHighlight string `json:"preview_highlight"`
	Lines            int    `json:"lines"`
	LinesMore        int    `json:"lines_more"`

	IsPublic        bool     `json:"is_public"`
	PublicUrlShared bool     `json:"public_url_shared"`
	Channels        []string `json:"channels"`
	Groups          []string `json:"groups"`
	InitialComment  Comment  `json:"initial_comment"`
	NumStars        int      `json:"num_stars"`
	IsStarred       bool     `json:"is_starred"`
}

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
	err := ParseResponse(path, values, response, debug)
	if err != nil {
		return nil, err
	}
	if !response.Ok {
		return nil, errors.New(response.Error)
	}
	return response, nil
}

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

func (api *Slack) UploadFile() {
	log.Fatal("Not implemented yet")
}
