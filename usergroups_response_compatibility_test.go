package slack

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestUserGroupUserCountJSONForms(t *testing.T) {
	tests := []struct {
		name         string
		payload      string
		initialCount int
		wantCount    int
	}{
		{name: "string", payload: `{"user_count":"4"}`, wantCount: 4},
		{name: "number", payload: `{"user_count":4}`, wantCount: 4},
		{name: "null", payload: `{"user_count":null}`, initialCount: 7, wantCount: 7},
		{name: "omitted", payload: `{}`, initialCount: 7, wantCount: 7},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userGroup := UserGroup{UserCount: tt.initialCount}
			if err := json.Unmarshal([]byte(tt.payload), &userGroup); err != nil {
				t.Fatalf("Unmarshal() error = %v", err)
			}
			if userGroup.UserCount != tt.wantCount {
				t.Fatalf("UserCount = %d, want %d", userGroup.UserCount, tt.wantCount)
			}
		})
	}
}

func TestUserGroupDocumentedFields(t *testing.T) {
	payload := []byte(`{
		"id": "S123",
		"auto_provision": true,
		"channel_count": 2,
		"channels": ["C123", "C456"],
		"enterprise_id": "E123",
		"enterprise_subteam_id": "S456",
		"is_editing_restricted": true,
		"is_idp_group": true,
		"is_membership_locked": true,
		"is_org_level": true,
		"is_section": true,
		"is_subteam": true,
		"is_visible": true,
		"teams": ["T123", "T456"],
		"prefs": {
			"channels": ["C123"],
			"groups": ["G123"],
			"file_id": "F123",
			"additional_channels": ["C456"]
		},
		"user_count": "4"
	}`)

	var userGroup UserGroup
	if err := json.Unmarshal(payload, &userGroup); err != nil {
		t.Fatalf("Unmarshal() error = %v", err)
	}

	if !userGroup.AutoProvision || userGroup.ChannelCount != 2 {
		t.Errorf("AutoProvision = %t, ChannelCount = %d", userGroup.AutoProvision, userGroup.ChannelCount)
	}
	if !reflect.DeepEqual(userGroup.Channels, []string{"C123", "C456"}) {
		t.Errorf("Channels = %#v", userGroup.Channels)
	}
	if userGroup.EnterpriseID != "E123" || userGroup.EnterpriseSubteamID != "S456" {
		t.Errorf("enterprise IDs = %q, %q", userGroup.EnterpriseID, userGroup.EnterpriseSubteamID)
	}
	if !userGroup.IsEditingRestricted || !userGroup.IsIDPGroup || !userGroup.IsMembershipLocked ||
		!userGroup.IsOrgLevel || !userGroup.IsSection || !userGroup.IsSubteam || !userGroup.IsVisible {
		t.Errorf("documented flags were not decoded: %#v", userGroup)
	}
	if !reflect.DeepEqual(userGroup.Teams, []string{"T123", "T456"}) {
		t.Errorf("Teams = %#v", userGroup.Teams)
	}
	if userGroup.Prefs.FileID != "F123" {
		t.Errorf("Prefs.FileID = %q", userGroup.Prefs.FileID)
	}
	if !reflect.DeepEqual(userGroup.Prefs.AdditionalChannels, []string{"C456"}) {
		t.Errorf("Prefs.AdditionalChannels = %#v", userGroup.Prefs.AdditionalChannels)
	}
	if userGroup.UserCount != 4 {
		t.Errorf("UserCount = %d, want 4", userGroup.UserCount)
	}

	var userCount int = userGroup.UserCount
	_ = userCount
}

func TestUserGroupNewZeroFieldsRemainOmitted(t *testing.T) {
	encoded, err := json.Marshal(UserGroup{ID: "S123"})
	if err != nil {
		t.Fatalf("Marshal() error = %v", err)
	}

	var fields map[string]json.RawMessage
	if err := json.Unmarshal(encoded, &fields); err != nil {
		t.Fatalf("Unmarshal() error = %v", err)
	}
	for _, field := range []string{
		"auto_provision",
		"channel_count",
		"channels",
		"enterprise_id",
		"enterprise_subteam_id",
		"is_editing_restricted",
		"is_idp_group",
		"is_membership_locked",
		"is_org_level",
		"is_section",
		"is_subteam",
		"is_visible",
		"teams",
	} {
		if _, ok := fields[field]; ok {
			t.Errorf("zero-value field %q was marshaled", field)
		}
	}

	var prefs map[string]json.RawMessage
	if err := json.Unmarshal(fields["prefs"], &prefs); err != nil {
		t.Fatalf("unmarshal prefs: %v", err)
	}
	for _, field := range []string{"file_id", "additional_channels"} {
		if _, ok := prefs[field]; ok {
			t.Errorf("zero-value prefs field %q was marshaled", field)
		}
	}
}
