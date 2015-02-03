package slack

import (
	"errors"
	"net/url"
	"strconv"
)

const (
	DEFAULT_STARS_USERID = ""
	DEFAULT_STARS_COUNT  = 100
	DEFAULT_STARS_PAGE   = 1
)

type StarsParameters struct {
	User  string
	Count int
	Page  int
}

// TODO: Verify this. The whole thing is complicated. I don't like the way they mixed things
// It also appears to be a bug in parsing the message
type StarredItem struct {
	Type      string `json:"type"`
	ChannelId string `json:"channel"`
	Message   `json:"message,omitempty"`
	File      `json:"file,omitempty"`
	Comment   `json:"comment,omitempty"`
}

type starsResponseFull struct {
	Items  []StarredItem `json:"items"`
	Paging `json:"paging"`
	SlackResponse
}

func NewStarsParameters() StarsParameters {
	return StarsParameters{
		User:  DEFAULT_STARS_USERID,
		Count: DEFAULT_STARS_COUNT,
		Page:  DEFAULT_STARS_PAGE,
	}
}

// GetStarred returns a list of StarredItem items. The user then has to iterate over them and figure out what they should
// be looking at according to what is in the Type.
//    for _, item := range items {
//        switch c.Type {
//        case "file_comment":
//            log.Println(c.Comment)
//        case "file":
//             ...
//        }
//    }
func (api *Slack) GetStarred(params StarsParameters) ([]StarredItem, *Paging, error) {
	response := &starsResponseFull{}
	values := url.Values{
		"token": {api.config.token},
	}
	if params.User != DEFAULT_STARS_USERID {
		values.Add("user", params.User)
	}
	if params.Count != DEFAULT_STARS_COUNT {
		values.Add("count", strconv.Itoa(params.Count))
	}
	if params.Page != DEFAULT_STARS_PAGE {
		values.Add("page", strconv.Itoa(params.Page))
	}
	err := parseResponse("stars.list", values, response, api.debug)
	if err != nil {
		return nil, nil, err
	}
	if !response.Ok {
		return nil, nil, errors.New(response.Error)
	}
	return response.Items, &response.Paging, nil
}
