package slack

// CallBlock defines data that is used to display call fields.
//
// More Information: https://api.slack.com/reference/block-kit/blocks#call
type CallBlock struct {
	Type                   MessageBlockType `json:"type"`
	BlockID                string           `json:"block_id,omitempty"`
	CallID                 string           `json:"call_id,omitempty"`
	Call                   CallObject       `json:"call,omitempty"`
	APIDecorationAvailable bool             `json:"api_decoration_available,omitempty"`
}

// CallObject declares the interface for call blocks.
type CallObject interface{}

// MessageBlockType defines a named string type to define each media backend type
// as a constant for use within the package.
type MediaBackendType string

const (
	MBETPlatformCall MediaBackendType = "platform_call"
)

type ZoomCall struct {
	CallObject
	MediaBackendType MediaBackendType `json:"media_backend_type"`
	Info             ZoomCallInfo     `json:"v1,omitempty"`
}

type ZoomCallInfo struct {
	AppId              string             `json:"app_id"`
	ActiveParticipants []*CallParticipant `json:"active_participants"`
	AllParticipants    []*CallParticipant `json:"all_participants"`
	Channels           []string           `json:"channels"`
	CreatedBy          string             `json:"created_by"`
	DateEnd            uint64             `json:"date_end"`
	DateStart          uint64             `json:"date_start"`
	DesktopAppJoinUrl  string             `json:"desktop_app_join_url"`
	DisplayId          string             `json:"display_id"`
	HasEnded           bool               `json:"has_ended"`
	Id                 string             `json:"id"`
	IsDmCall           bool               `json:"is_dm_call"`
	JoinUrl            string             `json:"join_url"`
	Name               string             `json:"name"`
	WasAccepted        bool               `json:"was_accepted"`
	WasMissed          bool               `json:"was_missed"`
	WasRejected        bool               `json:"was_rejected"`
}

type CallParticipant struct {
	ID          string `json:"external_id"`
	AvatarUrl   string `json:"avatar_url,omitempty"`
	DisplayName string `json:"display_name,omitempty"`
}

// BlockType returns the type of the block
func (s CallBlock) BlockType() MessageBlockType {
	return s.Type
}

// NewCallBlock returns a new instance of an input block
func NewCallBlock(blockID, callID string, call CallObject) *CallBlock {
	return &CallBlock{
		Type:    MBTCall,
		BlockID: blockID,
		CallID:  callID,
		Call:    call,
	}
}
