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
	Type           string      `json:"type"`
	User           string      `json:"user"`
	Text           string      `json:"text"`
	TimeStamp      string      `json:"ts"`
	Channel        string      `json:"channel"`
	EventTimeStamp json.Number `json:"event_ts"`
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
	Type            string      `json:"type"`
	Subtype         string      `json:"subtype"`
	User            string      `json:"user"`
	Username        string      `json:"username"`
	Text            string      `json:"text"`
	TimeStamp       string      `json:"ts"`
	ThreadTimeStamp string      `json:"thread_ts"`
	Channel         string      `json:"channel"`
	ChannelType     string      `json:"channel_type"`
	EventTimeStamp  json.Number `json:"event_ts"`
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
