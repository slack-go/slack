package slack

//----------------------------------------------------------------------------------------------------------------
// Admin Teams functions
//----------------------------------------------------------------------------------------------------------------
//-- Support types for AdminTeams list
type AdminTeamPrimaryOwner struct {
	PrimaryOwner struct {
		UserID string `json:"user_id"`
		Email  string `json:"email"`
	}
}
type AdminTeam struct {
	ID              string                `json:"id"`
	Name            string                `json:"name"`
	Discoverability string                `json:"discoverability"`
	PrimaryOwner    AdminTeamPrimaryOwner `json:"primary_owner"`
	TeamURL         string                `json:"team_url"`
}

//----------------------------------------------------------------------------------------------------------------
func (admin *Client) AdminTeamsList() ([]AdminTeam, error) {
	// Make the API Call
	items, err := admin.
		GenericAdminRequest("admin.teams.list").
		UrlParamString("limit", "100"). // This is the default limit anyway
		Query("teams")

	// Return an empty list on error
	if err != nil {
		return []AdminTeam{}, err
	}

	// Extract our map[string]interface{} into an actual typed array
	arr := make([]AdminTeam, len(items.Responses))
	err = items.Extract(&arr)

	// Return an empty list on error
	if err != nil {
		return []AdminTeam{}, err
	}

	// Return the actual results
	return arr, err
}

//--

//----------------------------------------------------------------------------------------------------------------
// Admin Conversations functions
//----------------------------------------------------------------------------------------------------------------
//-- Support types for Admin conversation queries
type AdminConversation struct {
	ID                        string   `json:"id"`
	Name                      string   `json:"name"`
	Purpose                   string   `json:"purpose"`
	MemberCount               int      `json:"member_count,omitempty"`
	Created                   JSONTime `json:"created"`
	CreatorID                 string   `json:"creator_id"`
	IsPrivate                 bool     `json:"is_private"`
	IsArchived                bool     `json:"is_archived"`
	IsGeneral                 bool     `json:"is_general"`
	LastActivityTimestamp     JSONTime `json:"last_activity_ts"`
	IsExtShared               bool     `json:"is_ext_shared"`
	IsGlobalShared            bool     `json:"is_global_shared"`
	IsOrgDefault              bool     `json:"is_org_default"`
	IsOrgMandatory            bool     `json:"is_org_mandatory"`
	IsOrgShared               bool     `json:"is_org_shared"`
	IsFrozen                  bool     `json:"is_frozen"`
	ConnectedTeamIDs          []string `json:"connected_team_ids"`
	InternalTeamIDsCount      int      `json:"internal_team_ids_count,omitempty"`
	InternalTeamIDsSampleTeam string   `json:"internal_team_ids_sample_team,omitempty"`
	PendingConnectedTeamIDs   []string `json:"pending_connected_team_ids"`
	IsPendingExtShared        bool     `json:"is_pending_ext_shared"`
}

//----------------------------------------------------------------------------------------------------------------
func (admin *Client) AdminConversationsSearch(query string) ([]AdminConversation, error) {
	// Make the API Call
	items, err := admin.
		GenericAdminRequest("admin.conversations.search").
		UrlParamString("query", "nowak").
		UrlParamString("limit", "25"). // Enough to give a reasonable paging, but not so much that th equery times out
		Query("conversations")

	// Return an empty list on error
	if err != nil {
		return []AdminConversation{}, err
	}

	// Extract our map[string]interface{} into an actual typed array
	arr := make([]AdminConversation, len(items.Responses))
	err = items.Extract(&arr)

	// Return an empty list on error
	if err != nil {
		return []AdminConversation{}, err
	}

	// Return the actual results
	return arr, err
}
