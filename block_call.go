package slack

// CallBlock defines information
type CallBlock struct {
	Type                   MessageBlockType `json:"type"`
	CallID                 string           `json:"call_id"`
	BlockID                string           `json:"block_id"`
	APIDecorationAvailable bool             `json:"api_decoration_available"`
	Call                   Call             `json:"call"`
}

type AppIconUrls struct {
	Image32       string `json:"image_32"`
	Image36       string `json:"image_36"`
	Image48       string `json:"image_48"`
	Image64       string `json:"image_64"`
	Image72       string `json:"image_72"`
	Image96       string `json:"image_96"`
	Image128      string `json:"image_128"`
	Image192      string `json:"image_192"`
	Image512      string `json:"image_512"`
	Image1024     string `json:"image_1024"`
	ImageOriginal string `json:"image_original"`
}
type CallInfo struct {
	ID                 string        `json:"id"`
	AppID              string        `json:"app_id"`
	AppIconUrls        AppIconUrls   `json:"app_icon_urls"`
	DateStart          int           `json:"date_start"`
	ActiveParticipants []interface{} `json:"active_participants"`
	AllParticipants    []interface{} `json:"all_participants"`
	DisplayID          string        `json:"display_id"`
	JoinURL            string        `json:"join_url"`
	DesktopAppJoinURL  string        `json:"desktop_app_join_url"`
	Name               string        `json:"name"`
	CreatedBy          string        `json:"created_by"`
	DateEnd            int           `json:"date_end"`
	Channels           []string      `json:"channels"`
	IsDmCall           bool          `json:"is_dm_call"`
	WasRejected        bool          `json:"was_rejected"`
	WasMissed          bool          `json:"was_missed"`
	WasAccepted        bool          `json:"was_accepted"`
	HasEnded           bool          `json:"has_ended"`
}
type Call struct {
	CallInfo         CallInfo `json:"v1"`
	MediaBackendType string   `json:"media_backend_type"`
}

// BlockType returns the type of the block
func (s CallBlock) BlockType() MessageBlockType {
	return s.Type
}
