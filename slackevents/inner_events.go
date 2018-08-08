// inner_events.go provides EventsAPI particular inner events

package slackevents

import "encoding/json"

// EventsAPIInnerEvent the inner event of a EventsAPI event_callback Event.
type EventsAPIInnerEvent struct {
	Type string `json:"type"`
	Data interface{}
}

// AppMentionEvent is an (inner) EventsAPI subscribable event.
type AppMentionEvent struct {
	Type            string      `json:"type"`
	User            string      `json:"user"`
	Text            string      `json:"text"`
	TimeStamp       string      `json:"ts"`
	ThreadTimeStamp string      `json:"thread_ts"`
	Channel         string      `json:"channel"`
	EventTimeStamp  json.Number `json:"event_ts"`
}

// AppUninstalledEvent Your Slack app was uninstalled.
type AppUninstalledEvent struct {
	Type string `json:"type"`
}

// GridMigrationFinishedEvent An enterprise grid migration has finished on this workspace.
type GridMigrationFinishedEvent struct {
	Type         string `json:"type"`
	EnterpriseID string `json:"enterprise_id"`
}

// GridMigrationStartedEvent An enterprise grid migration has started on this workspace.
type GridMigrationStartedEvent struct {
	Type         string `json:"type"`
	EnterpriseID string `json:"enterprise_id"`
}

// LinkSharedEvent A message was posted containing one or more links relevant to your application
type LinkSharedEvent struct {
	Type             string        `json:"type"`
	User             string        `json:"user"`
	TimeStamp        string        `json:"ts"`
	Channel          string        `json:"channel"`
	MessageTimeStamp json.Number   `json:"message_ts"`
	Links            []sharedLinks `json:"links"`
}

type sharedLinks struct {
	Domain string `json:"domain"`
	URL    string `json:"url"`
}

// MessageEvent occurs when a variety of types of messages has been posted.
// Parse ChannelType to see which
// if ChannelType = "group", this is a private channel message
// if ChannelType = "channel", this message was sent to a channel
// if ChannelType = "im", this is a private message
// if ChannelType = "mim", A message was posted in a multiparty direct message channel
// TODO: Improve this so that it is not required to manually parse ChannelType
type MessageEvent struct {
	// Basic Message Event - https://api.slack.com/events/message
	Type            string      `json:"type"`
	User            string      `json:"user"`
	Text            string      `json:"text"`
	ThreadTimeStamp string      `json:"thread_ts"`
	TimeStamp       string      `json:"ts"`
	Channel         string      `json:"channel"`
	ChannelType     string      `json:"channel_type"`
	EventTimeStamp  json.Number `json:"event_ts"`

	// Edited Message
	Message         *MessageEvent `json:"message,omitempty"`
	PreviousMessage *MessageEvent `json:"previous_message,omitempty"`
	Edited          *Edited       `json:"edited,omitempty"`

	// Message Subtypes
	SubType string `json:"subtype,omitempty"`

	// bot_message (https://api.slack.com/events/message/bot_message)
	BotID    string `json:"bot_id,omitempty"`
	Username string `json:"username,omitempty"`
	Icons    *Icon  `json:"icons,omitempty"`

	Upload bool   `json:"upload"`
	Files  []File `json:"files"`
}

// File is a file upload
type File struct {
	ID                 string `json:"id"`
	Created            int    `json:"created"`
	Timestamp          int    `json:"timestamp"`
	Name               string `json:"name"`
	Title              string `json:"title"`
	Mimetype           string `json:"mimetype"`
	Filetype           string `json:"filetype"`
	PrettyType         string `json:"pretty_type"`
	User               string `json:"user"`
	Editable           bool   `json:"editable"`
	Size               int    `json:"size"`
	Mode               string `json:"mode"`
	IsExternal         bool   `json:"is_external"`
	ExternalType       string `json:"external_type"`
	IsPublic           bool   `json:"is_public"`
	PublicURLShared    bool   `json:"public_url_shared"`
	DisplayAsBot       bool   `json:"display_as_bot"`
	Username           string `json:"username"`
	URLPrivate         string `json:"url_private"`
	URLPrivateDownload string `json:"url_private_download"`
	Thumb64            string `json:"thumb_64"`
	Thumb80            string `json:"thumb_80"`
	Thumb360           string `json:"thumb_360"`
	Thumb360W          int    `json:"thumb_360_w"`
	Thumb360H          int    `json:"thumb_360_h"`
	Thumb480           string `json:"thumb_480"`
	Thumb480W          int    `json:"thumb_480_w"`
	Thumb480H          int    `json:"thumb_480_h"`
	Thumb160           string `json:"thumb_160"`
	Thumb720           string `json:"thumb_720"`
	Thumb720W          int    `json:"thumb_720_w"`
	Thumb720H          int    `json:"thumb_720_h"`
	Thumb800           string `json:"thumb_800"`
	Thumb800W          int    `json:"thumb_800_w"`
	Thumb800H          int    `json:"thumb_800_h"`
	Thumb960           string `json:"thumb_960"`
	Thumb960W          int    `json:"thumb_960_w"`
	Thumb960H          int    `json:"thumb_960_h"`
	Thumb1024          string `json:"thumb_1024"`
	Thumb1024W         int    `json:"thumb_1024_w"`
	Thumb1024H         int    `json:"thumb_1024_h"`
	ImageExifRotation  int    `json:"image_exif_rotation"`
	OriginalW          int    `json:"original_w"`
	OriginalH          int    `json:"original_h"`
	Permalink          string `json:"permalink"`
	PermalinkPublic    string `json:"permalink_public"`
}

// Edited is included when a Message is edited
type Edited struct {
	User      string `json:"user"`
	TimeStamp string `json:"ts"`
}

// Icon is used for bot messages
type Icon struct {
	IconURL   string `json:"icon_url,omitempty"`
	IconEmoji string `json:"icon_emoji,omitempty"`
}

// IsEdited checks if the MessageEvent is caused by an edit
func (e MessageEvent) IsEdited() bool {
	return e.Message != nil &&
		e.Message.Edited != nil
}

const (
	// AppMention is an Events API subscribable event
	AppMention = "app_mention"
	// AppUninstalled Your Slack app was uninstalled.
	AppUninstalled = "app_uninstalled"
	// GridMigrationFinished An enterprise grid migration has finished on this workspace.
	GridMigrationFinished = "grid_migration_finished"
	// GridMigrationStarted An enterprise grid migration has started on this workspace.
	GridMigrationStarted = "grid_migration_started"
	// LinkShared A message was posted containing one or more links relevant to your application
	LinkShared = "link_shared"
	// Message A message was posted to a channel, private channel (group), im, or mim
	Message = "message"
)

// EventsAPIInnerEventMapping maps INNER Event API events to their corresponding struct
// implementations. The structs should be instances of the unmarshalling
// target for the matching event type.
var EventsAPIInnerEventMapping = map[string]interface{}{
	AppMention:            AppMentionEvent{},
	AppUninstalled:        AppUninstalledEvent{},
	GridMigrationFinished: GridMigrationFinishedEvent{},
	GridMigrationStarted:  GridMigrationStartedEvent{},
	LinkShared:            LinkSharedEvent{},
	Message:               MessageEvent{},
}
