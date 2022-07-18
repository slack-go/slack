package slack

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
)

/*----------------------------------------------------------------------------------------------------------------
The original package doesn't have support for most Admin functions, and I just need a few, but I don't want to spend the
effort fully defining everything, so I'm taking a different approach here and going generic.

The idea is that you construct a request, add parameters to it, then execute it (telling it which response field you're
interested in tracking)
It will then page through and gather those items up as an array of interface{} objects.
You can then extract those into either a strongly typed object, or easily unmarshall them into whatever other custom
types you like without having to modify this package.

The original package is restricted to golang 1.18 so no generics, which makes things a little uglier, but we have what
we have
----------------------------------------------------------------------------------------------------------------*/

/*----------------------------------------------------------------------------------------------------------------
Sample usage for https://api.slack.com/methods/admin.conversations.search

items, err := adminClient.
	GenericAdminRequest("admin.conversations.search"). //
	UrlParamString("query", "something").              //
	UrlParamString("limit", "10").                     //
	Query("conversations")
arr, _ := items.ExtractAdminConversations()
----------------------------------------------------------------------------------------------------------------*/

// This represents our request
type genericAdminRequest struct {
	api       *Client
	cmd       string
	urlValues url.Values
}

// This is how you start off a request... pass in the urlFragment from the API command, it will be appended to the general Slack WebAPI URL
func (api *Client) GenericAdminRequest(cmd string) *genericAdminRequest {
	return &genericAdminRequest{
		api: api,
		cmd: cmd,
		urlValues: url.Values{
			"token": {api.token}, // Could use api.appLevelToken instead, but .token works with more stuff without me having to change anything
		},
	}
}

// These are how you add Url Parameters to the given request
func (req *genericAdminRequest) UrlParamString(name, value string) *genericAdminRequest {
	return req.UrlParamStringArr(name, []string{value})
}

func (req *genericAdminRequest) UrlParamStringArr(name string, values []string) *genericAdminRequest {
	req.urlValues[name] = values
	return req
}

//----------------------------------------------------------------------------------------------------------------
// While we store the map results so we get all fields, we do unmarshall into this type so we get
// structured access to the common Slack status and token fields
// Note: We could union the above types in here and have just a single unmarshall, but without a way to dynamically
//   choose the field we want to extract items from, it doesn't really add much value
type SlackCursor struct {
	SlackResponse
	NextCursor string `json:"next_cursor"`
}

//----------------------------------------------------------------------------------------------------------------
// This is returned by Execute and is used to store the collection of response (as interface{}s)
type GenericExecutionResponse struct {
	Responses  []interface{}
	OutputHint string // Records the fieldName the responses list was extracted from
}

//----------------------------------------------------------------------------------------------------------------
//-- Execute a Slack WebAPI Command, returning only whether there was an error or not
func (req *genericAdminRequest) Execute() error {
	_, err := req.QueryContext(context.Background(), "")
	return err
}

//----------------------------------------------------------------------------------------------------------------
//-- Execute a Slack WebAPI Command, returning a collated collection of expected objects
func (req *genericAdminRequest) Query(returnFieldName string) (GenericExecutionResponse, error) {
	return req.QueryContext(context.Background(), returnFieldName)
}

//----------------------------------------------------------------------------------------------------------------
//-- This will execute the request and return any errors
//-- If a non-empty returnFieldName is provided, that named field will be collated from the SlackResponse objects and returned
func (req *genericAdminRequest) QueryContext(ctx context.Context, returnFieldName string) (GenericExecutionResponse, error) {
	cursor := ""
	result := GenericExecutionResponse{
		OutputHint: returnFieldName,
	}

	for {
		// Make our own copy of the map so we can add/tweak necessary options without affecting the original
		curValues := make(url.Values, len(req.urlValues))
		for k, v := range req.urlValues {
			curValues[k] = make([]string, len(v))
			copy(curValues[k], v)
		}

		// Add cursor if necessary
		if cursor != "" {
			curValues["cursor"] = []string{cursor}
		}

		// Do the request
		resp := make(map[string]interface{})
		endpoint := fmt.Sprintf("%s%s", APIURL, req.cmd)
		if err := postForm(ctx, req.api.httpclient, endpoint, curValues, &resp, req.api); err != nil {
			return result, err
		}

		// Coerce our map into the common Slack Response fields so we can check status
		sr := SlackCursor{}
		err := castFieldToData(resp, &sr)
		if err != nil {
			return result, err
		}

		// We errored out
		if sr.Err() != nil {
			return result, err
		}

		// Squirrel our field info away
		if returnFieldName != "" {
			items, ok := resp[returnFieldName].([]interface{})
			if !ok {
				// If we're not an array, then just treat us as a single and add us anyway
				item := resp[returnFieldName].(interface{})
				result.Responses = append(result.Responses, item)
			} else {
				result.Responses = append(result.Responses, items...)
			}
		}

		// Next Cursor could be in one of a couple of places
		nextCursor := sr.NextCursor
		if nextCursor == "" {
			nextCursor = sr.ResponseMetadata.Cursor
		}

		// If we ran out of results or we don't care about results, get out, otherwise keep looping around
		if nextCursor == "" || returnFieldName == "" {
			return result, err
		} else {
			cursor = nextCursor
		}

	}
}

//----------------------------------------------------------------------------------------------------------------
// Utility function for above...
//  src is the source data we intend to remarshall (convert from interface{} to jason then to a struct)
//  dst is a pointer to the object we want to unmarshall into and governs the unmarshalling
func castFieldToData(src interface{}, dst interface{}) error {
	jsonbody, err := json.Marshal(src)
	if err != nil {
		return err
	}

	err = json.Unmarshal(jsonbody, dst)
	if err != nil {
		return err
	}

	return nil
}

//----------------------------------------------------------------------------------------------------------------
// Given a pointer to a typed array, this will unmarkshall into that array, it's the simplest generic output handler
func (r GenericExecutionResponse) Extract(arr interface{}) error {
	return castFieldToData(r.Responses, arr)
}
