package slack

type TeamJoinEvent struct {
	Type string `json:"type"`
	User User   `json:"user,omitempty"`
}

type TeamRenameEvent struct {
	Type           string         `json:"type"`
	Name           string         `json:"name,omitempty"`
	EventTimestamp JSONTimeString `json:"event_ts,omitempty"`
}

type TeamPrefChangeEvent struct {
	Type  string   `json:"type"`
	Name  string   `json:"name,omitempty"`
	Value []string `json:"value,omitempty"`
}

type TeamDomainChangeEvent struct {
	Type   string `json:"type"`
	Url    string `json:"url"`
	Domain string `json:"domain"`
}

type TeamMigrationStartedEvent struct {
	Type string `json:"type"`
}
