package slack

// MessageBlockType defines a named string type to define each block type
// as a constant for use within the package.
type MessageBlockType string

const (
	MBTSection  MessageBlockType = "section"
	MBTDivider  MessageBlockType = "divider"
	MBTImage    MessageBlockType = "image"
	MBTAction   MessageBlockType = "actions"
	MBTContext  MessageBlockType = "context"
	MBTFile     MessageBlockType = "file"
	MBTInput    MessageBlockType = "input"
	MBTHeader   MessageBlockType = "header"
	MBTRichText MessageBlockType = "rich_text"
	MBTCall     MessageBlockType = "call"
	MBTVideo    MessageBlockType = "video"
	MBTMarkdown MessageBlockType = "markdown"
)

// Block defines an interface all block types should implement
// to ensure consistency between blocks.
type Block interface {
	BlockType() MessageBlockType
	ID() string
}

// Blocks is a convenience struct defined to allow dynamic unmarshalling of
// the "blocks" value in Slack's JSON response, which varies depending on block type
type Blocks struct {
	BlockSet []Block `json:"blocks,omitempty" form:"blocks"`
}

// BlockAction is the action callback sent when a block is interacted with
type BlockAction struct {
	ActionID              string              `json:"action_id" form:"action_id"`
	BlockID               string              `json:"block_id" form:"block_id"`
	Type                  ActionType          `json:"type" form:"type"`
	Text                  TextBlockObject     `json:"text" form:"text"`
	Value                 string              `json:"value" form:"value"`
	Files                 []File              `json:"files" form:"files"`
	ActionTs              string              `json:"action_ts" form:"action_ts"`
	SelectedOption        OptionBlockObject   `json:"selected_option" form:"selected_option"`
	SelectedOptions       []OptionBlockObject `json:"selected_options" form:"selected_option"`
	SelectedUser          string              `json:"selected_user" form:"selected_user"`
	SelectedUsers         []string            `json:"selected_users" form:"selected_users"`
	SelectedChannel       string              `json:"selected_channel" form:"selected_channel"`
	SelectedChannels      []string            `json:"selected_channels" form:"selected_channels"`
	SelectedConversation  string              `json:"selected_conversation" form:"selected_conversation"`
	SelectedConversations []string            `json:"selected_conversations" form:"selected_conversations"`
	SelectedDate          string              `json:"selected_date" form:"selected_date"`
	SelectedTime          string              `json:"selected_time" form:"selected_time"`
	SelectedDateTime      int64               `json:"selected_date_time" form:"selected_data_time"`
	Timezone              string              `json:"timezone" form:"timezone"`
	InitialOption         OptionBlockObject   `json:"initial_option" form:"initial_option"`
	InitialUser           string              `json:"initial_user" form:"initial_user"`
	InitialChannel        string              `json:"initial_channel" form:"initial_channel"`
	InitialConversation   string              `json:"initial_conversation" form:"initial_conversation"`
	InitialDate           string              `json:"initial_date" form:"initial_date"`
	InitialTime           string              `json:"initial_time" form:"initial_time"`
	RichTextValue         RichTextBlock       `json:"rich_text_value" form:"rich_text_value"`
}

// actionType returns the type of the action
func (b BlockAction) actionType() ActionType {
	return b.Type
}

// NewBlockMessage creates a new Message that contains one or more blocks to be displayed
func NewBlockMessage(blocks ...Block) Message {
	return Message{
		Msg: Msg{
			Blocks: Blocks{
				BlockSet: blocks,
			},
		},
	}
}

// AddBlockMessage appends a block to the end of the existing list of blocks
func AddBlockMessage(message Message, newBlk Block) Message {
	message.Msg.Blocks.BlockSet = append(message.Msg.Blocks.BlockSet, newBlk)
	return message
}
