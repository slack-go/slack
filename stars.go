package slack

import (
	"errors"
	"net/url"
	"strconv"
)

const (
	DEFAULT_STARS_USER  = ""
	DEFAULT_STARS_COUNT = 100
	DEFAULT_STARS_PAGE  = 1
)

type StarsParameters struct {
	User  string
	Count int
	Page  int
}

// StarredItem is an item that has been starred.
type StarredItem struct {
	Item
}

type starsResponseFull struct {
	Items  []StarredItem `json:"items"`
	Paging `json:"paging"`
	SlackResponse
}

func NewStarsParameters() StarsParameters {
	return StarsParameters{
		User:  DEFAULT_STARS_USER,
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
	values := url.Values{
		"token": {api.config.token},
	}
	if params.User != DEFAULT_STARS_USER {
		values.Add("user", params.User)
	}
	if params.Count != DEFAULT_STARS_COUNT {
		values.Add("count", strconv.Itoa(params.Count))
	}
	if params.Page != DEFAULT_STARS_PAGE {
		values.Add("page", strconv.Itoa(params.Page))
	}
	response := &starsResponseFull{}
	err := parseResponse("stars.list", values, response, api.debug)
	if err != nil {
		return nil, nil, err
	}
	if !response.Ok {
		return nil, nil, errors.New(response.Error)
	}
	return response.Items, &response.Paging, nil
}
