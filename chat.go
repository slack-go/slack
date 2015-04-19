package slack

import (
	"encoding/json"
	"errors"
	"net/url"
	"strings"
)

const (
	DEFAULT_MESSAGE_USERNAME     = ""
	DEFAULT_MESSAGE_ASUSER       = false
	DEFAULT_MESSAGE_PARSE        = ""
	DEFAULT_MESSAGE_LINK_NAMES   = 0
	DEFAULT_MESSAGE_UNFURL_LINKS = true
	DEFAULT_MESSAGE_UNFURL_MEDIA = false
	DEFAULT_MESSAGE_ICON_URL     = ""
	DEFAULT_MESSAGE_ICON_EMOJI   = ""
)

type chatResponseFull struct {
	ChannelId string `json:"channel"`
	Timestamp string `json:"ts"`
	Text      string `json:"text"`
	SlackResponse
}

// AttachmentField contains information for an attachment field
// An Attachment can contain multiple of these
type AttachmentField struct {
	Title string `json:"title"`
	Value string `json:"value"`
	Short bool   `json:"short"`
}

// Attachment contains all the information for an attachment
type Attachment struct {
	Fallback string `json:"fallback"`

	Color string `json:"color,omitempty"`

	Pretext string `json:"pretext,omitempty"`

	AuthorName string `json:"author_name,omitempty"`
	AuthorLink string `json:"author_link,omitempty"`
	AuthorIcon string `json:"author_icon,omitempty"`

	Title     string `json:"title,omitempty"`
	TitleLink string `json:"title_link,omitempty"`

	Text string `json:"text"`

	ImageURL string `json:"image_url,omitempty"`

	Fields []AttachmentField `json:"fields,omitempty"`
}

// PostMessageParameters contains all the parameters necessary (including the optional ones) for a PostMessage() request
type PostMessageParameters struct {
	Text        string
	Username    string
	AsUser      bool
	Parse       string
	LinkNames   int
	Attachments []Attachment
	UnfurlLinks bool
	UnfurlMedia bool
	IconURL     string
	IconEmoji   string
}

// NewPostMessageParameters provides an instance of PostMessageParameters with all the sane default values set
func NewPostMessageParameters() PostMessageParameters {
	return PostMessageParameters{
		Username:    DEFAULT_MESSAGE_USERNAME,
		AsUser:      DEFAULT_MESSAGE_ASUSER,
		Parse:       DEFAULT_MESSAGE_PARSE,
		LinkNames:   DEFAULT_MESSAGE_LINK_NAMES,
		Attachments: nil,
		UnfurlLinks: DEFAULT_MESSAGE_UNFURL_LINKS,
		UnfurlMedia: DEFAULT_MESSAGE_UNFURL_MEDIA,
		IconURL:     DEFAULT_MESSAGE_ICON_URL,
		IconEmoji:   DEFAULT_MESSAGE_ICON_EMOJI,
	}
}

func chatRequest(path string, values url.Values, debug bool) (*chatResponseFull, error) {
	response := &chatResponseFull{}
	err := parseResponse(path, values, response, debug)
	if err != nil {
		return nil, err
	}
	if !response.Ok {
		return nil, errors.New(response.Error)
	}
	return response, nil
}

// DeleteMessage deletes a message in a channel
func (api *Slack) DeleteMessage(channelId, messageTimestamp string) (string, string, error) {
	values := url.Values{
		"token":   {api.config.token},
		"channel": {channelId},
		"ts":      {messageTimestamp},
	}
	response, err := chatRequest("chat.delete", values, api.debug)
	if err != nil {
		return "", "", err
	}
	return response.ChannelId, response.Timestamp, nil
}

func escapeMessage(message string) string {
	/*
		& replaced with &amp;
		< replaced with &lt;
		> replaced with &gt;
	*/
	replacer := strings.NewReplacer("&", "&amp;", "<", "&lt;", ">", "&gt;")
	return replacer.Replace(message)
}

// PostMessage sends a message to a channel
// Message is automatically escaped according to https://api.slack.com/docs/formatting
func (api *Slack) PostMessage(channelId string, text string, params PostMessageParameters) (channel string, timestamp string, err error) {
	values := url.Values{
		"token":   {api.config.token},
		"channel": {channelId},
		"text":    {escapeMessage(text)},
	}
	if params.Username != DEFAULT_MESSAGE_USERNAME {
		values.Set("username", string(params.Username))
	}
	if params.AsUser != DEFAULT_MESSAGE_ASUSER {
		values.Set("as_user", string(params.AsUser))
	}
	if params.Parse != DEFAULT_MESSAGE_PARSE {
		values.Set("parse", string(params.Parse))
	}
	if params.LinkNames != DEFAULT_MESSAGE_LINK_NAMES {
		values.Set("link_names", "1")
	}
	if params.Attachments != nil {
		attachments, err := json.Marshal(params.Attachments)
		if err != nil {
			return "", "", err
		}
		values.Set("attachments", string(attachments))
	}
	if params.UnfurlLinks == DEFAULT_MESSAGE_UNFURL_LINKS {
		values.Set("unfurl_links", "false")
	}
	if params.UnfurlMedia != DEFAULT_MESSAGE_UNFURL_MEDIA {
		values.Set("unfurl_media", "true")
	}
	if params.IconURL != DEFAULT_MESSAGE_ICON_URL {
		values.Set("icon_url", params.IconURL)
	}
	if params.IconEmoji != DEFAULT_MESSAGE_ICON_EMOJI {
		values.Set("icon_emoji", params.IconEmoji)
	}

	response, err := chatRequest("chat.postMessage", values, api.debug)
	if err != nil {
		return "", "", err
	}
	return response.ChannelId, response.Timestamp, nil
}

// UpdateMessage updates a message in a channel
func (api *Slack) UpdateMessage(channelId, timestamp, text string) (string, string, string, error) {
	values := url.Values{
		"token":   {api.config.token},
		"channel": {channelId},
		"text":    {escapeMessage(text)},
		"ts":      {timestamp},
	}
	response, err := chatRequest("chat.update", values, api.debug)
	if err != nil {
		return "", "", "", err
	}
	return response.ChannelId, response.Timestamp, response.Text, nil
}
