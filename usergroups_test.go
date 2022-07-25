package slack

import (
	"net/http"
	"reflect"
	"testing"
)

type userGroupsHandler struct {
	gotParams map[string]string
	response  string
}

func newUserGroupsHandler() *userGroupsHandler {
	return &userGroupsHandler{
		gotParams: make(map[string]string),
		response: `{
    "ok": true,
    "usergroup": {
        "id": "S0615G0KT",
        "team_id": "T060RNRCH",
        "is_usergroup": true,
        "name": "Marketing Team",
        "description": "Marketing gurus, PR experts and product advocates.",
        "handle": "marketing-team",
        "is_external": false,
        "date_create": 1446746793,
        "date_update": 1446746793,
        "date_delete": 0,
        "auto_type": null,
        "created_by": "U060RNRCZ",
        "updated_by": "U060RNRCZ",
        "deleted_by": null,
        "prefs": {
            "channels": [

            ],
            "groups": [

            ]
        },
        "user_count": 0
    }
}`,
	}
}

func (ugh *userGroupsHandler) handler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	for k, v := range r.Form {
		ugh.gotParams[k] = v[0]
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(ugh.response))
}

func TestCreateUserGroup(t *testing.T) {
	once.Do(startServer)
	api := New("testing-token", OptionAPIURL("http://"+serverAddr+"/"))

	tests := []struct {
		userGroup  UserGroup
		wantParams map[string]string
	}{
		{
			UserGroup{
				Name:        "Marketing Team",
				Description: "Marketing gurus, PR experts and product advocates.",
				Handle:      "marketing-team"},
			map[string]string{
				"token":       "testing-token",
				"name":        "Marketing Team",
				"description": "Marketing gurus, PR experts and product advocates.",
				"handle":      "marketing-team",
			},
		},
	}

	var rh *userGroupsHandler
	http.HandleFunc("/usergroups.create", func(w http.ResponseWriter, r *http.Request) { rh.handler(w, r) })

	for i, test := range tests {
		rh = newUserGroupsHandler()
		_, err := api.CreateUserGroup(test.userGroup)
		if err != nil {
			t.Fatalf("%d: Unexpected error: %s", i, err)
		}
		if !reflect.DeepEqual(rh.gotParams, test.wantParams) {
			t.Errorf("%d: Got params %#v, want %#v", i, rh.gotParams, test.wantParams)
		}
	}
}

func getUserGroups(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	response := []byte(`{
    "ok": true,
    "usergroups": [
        {
            "id": "S0614TZR7",
            "team_id": "T060RNRCH",
            "is_usergroup": true,
            "name": "Team Admins",
            "description": "A group of all Administrators on your team.",
            "handle": "admins",
            "is_external": false,
            "date_create": 1446598059,
            "date_update": 1446670362,
            "date_delete": 0,
            "auto_type": "admin",
            "created_by": "USLACKBOT",
            "updated_by": "U060RNRCZ",
            "deleted_by": null,
            "prefs": {
                "channels": [
                  "channel1",
                  "channel2"
                ],
                "groups": [
                  "group1",
                  "group2",
                  "group3"
                ]
			},
            "users": [
                "user1",
                "user2"
            ],
            "user_count": 2
        }
    ]
}`)
	rw.Write(response)
}

func TestGetUserGroups(t *testing.T) {
	http.HandleFunc("/usergroups.list", getUserGroups)

	once.Do(startServer)
	api := New("testing-token", OptionAPIURL("http://"+serverAddr+"/"))

	userGroups, err := api.GetUserGroups(GetUserGroupsOptionIncludeUsers(true))
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
		return
	}

	// t.Fatal refers to -> t.Errorf & return
	if len(userGroups) != 1 {
		t.Fatal(ErrIncorrectResponse)
	}

	S0614TZR7 := UserGroup{
		ID:          "S0614TZR7",
		TeamID:      "T060RNRCH",
		IsUserGroup: true,
		Name:        "Team Admins",
		Description: "A group of all Administrators on your team.",
		Handle:      "admins",
		IsExternal:  false,
		DateCreate:  1446598059,
		DateUpdate:  1446670362,
		DateDelete:  0,
		AutoType:    "admin",
		CreatedBy:   "USLACKBOT",
		UpdatedBy:   "U060RNRCZ",
		DeletedBy:   "",
		Prefs: UserGroupPrefs{
			Channels: []string{"channel1", "channel2"},
			Groups:   []string{"group1", "group2", "group3"},
		},
		Users: []string{
			"user1",
			"user2",
		},
		UserCount: 2,
	}

	if !reflect.DeepEqual(userGroups[0], S0614TZR7) {
		t.Errorf("Got %#v, want %#v", userGroups[0], S0614TZR7)
	}
}

func updateUserGroupsHandler() *userGroupsHandler {
	return &userGroupsHandler{
		gotParams: make(map[string]string),
		response: `{
    "ok": true,
    "usergroup": {
        "id": "S0615G0KT",
        "team_id": "T060RNRCH",
        "is_usergroup": true,
        "name": "Marketing Team",
        "description": "Marketing gurus, PR experts and product advocates.",
        "handle": "marketing-team",
        "is_external": false,
        "date_create": 1446746793,
        "date_update": 1446746793,
        "date_delete": 0,
        "auto_type": null,
        "created_by": "U060RNRCZ",
        "updated_by": "U060RNRCZ",
        "deleted_by": null,
        "prefs": {
            "channels": [
				"channel1",
				"channel2"
            ],
            "groups": [

            ]
        },
        "user_count": 0
    }
}`,
	}
}
func TestUpdateUserGroup(t *testing.T) {
	once.Do(startServer)
	api := New("testing-token", OptionAPIURL("http://"+serverAddr+"/"))

	emptyDescription := ""
	presenceDescription := "Marketing gurus, PR experts and product advocates."

	tests := []struct {
		options    []UpdateUserGroupsOption
		wantParams map[string]string
	}{
		{
			[]UpdateUserGroupsOption{
				UpdateUserGroupsOptionName("Marketing Team"),
				UpdateUserGroupsOptionHandle("marketing-team"),
			},
			map[string]string{
				"token":     "testing-token",
				"usergroup": "S0615G0KT",
				"name":      "Marketing Team",
				"handle":    "marketing-team",
			},
		},
		{
			[]UpdateUserGroupsOption{
				UpdateUserGroupsOptionDescription(&presenceDescription),
				UpdateUserGroupsOptionChannels([]string{"channel1", "channel2"}),
			},
			map[string]string{
				"token":       "testing-token",
				"usergroup":   "S0615G0KT",
				"description": "Marketing gurus, PR experts and product advocates.",
				"channels":    "channel1,channel2",
			},
		},
		{
			[]UpdateUserGroupsOption{
				UpdateUserGroupsOptionDescription(&emptyDescription),
				UpdateUserGroupsOptionChannels([]string{}),
			},
			map[string]string{
				"token":       "testing-token",
				"usergroup":   "S0615G0KT",
				"description": "",
				"channels":    "",
			},
		},
	}

	var rh *userGroupsHandler
	http.HandleFunc("/usergroups.update", func(w http.ResponseWriter, r *http.Request) { rh.handler(w, r) })

	for i, test := range tests {
		rh = updateUserGroupsHandler()
		_, err := api.UpdateUserGroup("S0615G0KT", test.options...)
		if err != nil {
			t.Fatalf("%d: Unexpected error: %s", i, err)
		}
		if !reflect.DeepEqual(rh.gotParams, test.wantParams) {
			t.Errorf("%d: Got params %#v, want %#v", i, rh.gotParams, test.wantParams)
		}
	}
}
