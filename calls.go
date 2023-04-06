package slack

// "external_unique_id": "025169F6-E37A-4E62-BB54-7F93A0FC4C1F",
// "join_url": "https://callmebeepme.com/calls/1234567890",
// "desktop_app_join_url": "callapp://join/1234567890",
// "external_display_id": "705-292-868",
// "title": "Kimpossible sync up",
// "users": [
// 	{
// 		"slack_id": "U0MQG83FD"
// 	},
// 	{
// 		"external_id": "54321678",
// 		"display_name": "Kim Possible",
// 		"avatar_url": "https://callmebeepme.com/users/avatar1234.jpg"
// 	}
// ]

type Call struct {
	ID               string   `json:"id"`
	DateStart        JSONTime `json:"date_start"`
	ExternalUniqueID string   `json:"external_unique_id"`
	JoinUrl          string   `json:"join_url"`

	DesktopAppJoinUrl string `json:"desktop_app_join_url"`
	ExternalDisplayID string `json:"external_display_id"`

	Users []CallUser `json:"users"`
}

type CallUser struct{}
