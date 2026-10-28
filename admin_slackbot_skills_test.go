package slack

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAdminSlackbotSkillsList(t *testing.T) {
	api := skillTestClient(t, "admin.slackbot.skills.list", `{"category":"engineering","cursor":"next","limit":200,"creator_ids":["U1","W2"],"publish_status":"unpublished","auto_enabled_status":"enabled","sort":"date_updated","sort_dir":"desc"}`, `{"ok":true,"skills":[{"id":"F1","display_name":"Review","description":"Review releases","category":"engineering","creator_id":"W2","date_created":1790000000,"date_updated":null,"visibility":"invite_only","is_official":false,"usage_count":8,"library_count":4,"auto_enabled":{"permission_type":"named_entities","usergroup_ids":["S1"]}}],"response_metadata":{"next_cursor":"continue"}}`)
	response, err := api.AdminSlackbotSkillsList(nil,
		AdminSlackbotSkillsListOptionCategory("engineering"), AdminSlackbotSkillsListOptionCursor("next"),
		AdminSlackbotSkillsListOptionLimit(200), AdminSlackbotSkillsListOptionCreatorIDs("U1", "W2"),
		AdminSlackbotSkillsListOptionPublishStatus("unpublished"), AdminSlackbotSkillsListOptionAutoEnabledStatus("enabled"),
		AdminSlackbotSkillsListOptionSort("date_updated"), AdminSlackbotSkillsListOptionSortDir("desc"))
	require.NoError(t, err)
	require.Len(t, response.Skills, 1)
	assert.Nil(t, response.Skills[0].DateUpdated)
	assert.Equal(t, "invite_only", response.Skills[0].Visibility)
	assert.Equal(t, []string{"S1"}, response.Skills[0].AutoEnabled.UsergroupIDs)
	assert.Equal(t, "continue", response.ResponseMetadata.Cursor)

	api = skillTestClient(t, "admin.slackbot.skills.list", `{}`, `{"ok":true,"skills":[],"response_metadata":{"next_cursor":"continue"}}`)
	response, err = api.AdminSlackbotSkillsList()
	require.NoError(t, err)
	assert.Empty(t, response.Skills)
	assert.Equal(t, "continue", response.ResponseMetadata.Cursor)

	api = skillTestClient(t, "admin.slackbot.skills.list", `{"creator_ids":[],"cursor":"","limit":0,"sort_dir":"desc"}`, `{"ok":false,"error":"invalid_arguments"}`)
	response, err = api.AdminSlackbotSkillsList(AdminSlackbotSkillsListOptionCreatorIDs(), AdminSlackbotSkillsListOptionCursor(""), AdminSlackbotSkillsListOptionLimit(0), AdminSlackbotSkillsListOptionSortDir("asc"), AdminSlackbotSkillsListOptionSortDir("desc"))
	require.EqualError(t, err, "invalid_arguments")
	require.NotNil(t, response)
}

func TestAdminSlackbotSkillsPublish(t *testing.T) {
	api := skillTestClient(t, "admin.slackbot.skills.publish", `{"skill_ids":["F1","F1","F2"]}`, `{"ok":true}`)
	response, err := api.AdminSlackbotSkillsPublish([]string{"F1", "F1", "F2"})
	require.NoError(t, err)
	assert.True(t, response.Ok)

	api = skillTestClient(t, "admin.slackbot.skills.publish", `{"skill_ids":["F1","F2"]}`, `{"ok":false,"error":"failed_for_some_skills","failed_skills":[{"skill_id":"F2","error":"restricted_action"}]}`)
	response, err = api.AdminSlackbotSkillsPublish([]string{"F1", "F2"})
	require.EqualError(t, err, "failed_for_some_skills")
	require.NotNil(t, response)
	assert.Equal(t, []SlackbotSkillFailure{{SkillID: "F2", Error: "restricted_action"}}, response.FailedSkills)
}

func TestAdminSlackbotSkillsSetAutoEnabled(t *testing.T) {
	for _, audience := range []string{"everyone", "no_one"} {
		api := skillTestClient(t, "admin.slackbot.skills.setAutoEnabled", `{"skill_id":"F1","permission_type":"`+audience+`"}`, `{"ok":true,"skill_id":"F1","permission_type":"`+audience+`","usergroup_ids":[]}`)
		response, err := api.AdminSlackbotSkillsSetAutoEnabled("F1", audience, nil)
		require.NoError(t, err)
		assert.Equal(t, audience, response.PermissionType)
		assert.NotNil(t, response.UsergroupIDs)
	}
	api := skillTestClient(t, "admin.slackbot.skills.setAutoEnabled", `{"skill_id":"F1","permission_type":"named_entities","usergroup_ids":["S2","S1","S1"]}`, `{"ok":true,"skill_id":"F1","permission_type":"named_entities","usergroup_ids":["S1","S2"]}`)
	response, err := api.AdminSlackbotSkillsSetAutoEnabled("F1", "named_entities",
		AdminSlackbotSkillsSetAutoEnabledOptionUsergroupIDs("ignored"),
		AdminSlackbotSkillsSetAutoEnabledOptionUsergroupIDs("S2", "S1", "S1"))
	require.NoError(t, err)
	assert.Equal(t, []string{"S1", "S2"}, response.UsergroupIDs)

	for _, tt := range []struct{ code, details string }{
		{"skill_not_visible_to_audience", `"usergroup_ids":["S2"]`},
		{"usergroup_not_found", `"invalid_usergroup_ids":["S2"]`},
	} {
		api := skillTestClient(t, "admin.slackbot.skills.setAutoEnabled", `{"skill_id":"F1","permission_type":"named_entities","usergroup_ids":[]}`, `{"ok":false,"error":"`+tt.code+`",`+tt.details+`}`)
		response, err := api.AdminSlackbotSkillsSetAutoEnabled("F1", "named_entities", AdminSlackbotSkillsSetAutoEnabledOptionUsergroupIDs())
		require.EqualError(t, err, tt.code)
		require.NotNil(t, response)
		if tt.code == "usergroup_not_found" {
			assert.Equal(t, []string{"S2"}, response.InvalidUsergroupIDs)
		} else {
			assert.Equal(t, []string{"S2"}, response.UsergroupIDs)
		}
	}
}

func TestAdminSlackbotSkillsUnpublish(t *testing.T) {
	api := skillTestClient(t, "admin.slackbot.skills.unpublish", `{"skill_ids":["F1","F1"]}`, `{"ok":true}`)
	response, err := api.AdminSlackbotSkillsUnpublish([]string{"F1", "F1"})
	require.NoError(t, err)
	assert.True(t, response.Ok)

	api = skillTestClient(t, "admin.slackbot.skills.unpublish", `{"skill_ids":["F1","F2"]}`, `{"ok":false,"error":"failed_for_some_skills","failed_skills":[{"skill_id":"F2","error":"skill_not_found"}]}`)
	response, err = api.AdminSlackbotSkillsUnpublish([]string{"F1", "F2"})
	require.EqualError(t, err, "failed_for_some_skills")
	require.NotNil(t, response)
	assert.Equal(t, []SlackbotSkillFailure{{SkillID: "F2", Error: "skill_not_found"}}, response.FailedSkills)
}
