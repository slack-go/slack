package slack

const (
	DEFAULT_HISTORY_LATEST = ""
	DEFAULT_HISTORY_OLDEST = "0"
	DEFAULT_HISTORY_COUNT  = 100
)

type HistoryParameters struct {
	Latest string
	Oldest string
	Count  int
}

type History struct {
	Latest   string    `json:"latest"`
	Messages []Message `json:"messages"`
	HasMore  bool      `json:"has_more"`
}

func NewHistoryParameters() HistoryParameters {
	return HistoryParameters{
		Latest: DEFAULT_HISTORY_LATEST,
		Oldest: DEFAULT_HISTORY_OLDEST,
		Count:  DEFAULT_HISTORY_COUNT,
	}
}
