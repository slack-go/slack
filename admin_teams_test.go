package slack

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAdminTeamsList(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/admin.teams.list", r.URL.Path)
		if !assert.NoError(t, r.ParseForm()) {
			rw.WriteHeader(http.StatusBadRequest)
			return
		}
		assert.Equal(t, "testing-token", r.FormValue("token"))
		assert.Equal(t, "250", r.FormValue("limit"))
		assert.Equal(t, "cursor-1", r.FormValue("cursor"))

		rw.Header().Set("Content-Type", "application/json")
		_, _ = rw.Write([]byte(`{
			"ok": true,
			"teams": [{
				"id": "T1234",
				"name": "My Team",
				"discoverability": "hidden",
				"primary_owner": {
					"user_id": "W1234",
					"email": "owner@example.com"
				},
				"team_url": "https://example.slack.com/"
			}],
			"response_metadata": {
				"next_cursor": "cursor-2"
			}
		}`))
	}))
	defer ts.Close()

	api := New("testing-token", OptionAPIURL(ts.URL+"/"))
	response, err := api.AdminTeamsList(context.Background(), AdminTeamsListParams{
		Limit:  250,
		Cursor: "cursor-1",
	})
	require.NoError(t, err)
	require.Len(t, response.Teams, 1)

	team := response.Teams[0]
	assert.True(t, response.Ok)
	assert.Equal(t, "T1234", team.ID)
	assert.Equal(t, "My Team", team.Name)
	assert.Equal(t, TeamDiscoverability("hidden"), team.Discoverability)
	assert.Equal(t, "W1234", team.PrimaryOwner.UserID)
	assert.Equal(t, "owner@example.com", team.PrimaryOwner.Email)
	assert.Equal(t, "https://example.slack.com/", team.TeamURL)
	assert.Equal(t, "cursor-2", response.ResponseMetadata.Cursor)
}

func TestAdminTeamsListError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		if !assert.NoError(t, r.ParseForm()) {
			rw.WriteHeader(http.StatusBadRequest)
			return
		}
		assert.Empty(t, r.FormValue("limit"))
		assert.Empty(t, r.FormValue("cursor"))

		rw.Header().Set("Content-Type", "application/json")
		_, _ = rw.Write([]byte(`{"ok": false, "error": "invalid_cursor"}`))
	}))
	defer ts.Close()

	api := New("testing-token", OptionAPIURL(ts.URL+"/"))
	response, err := api.AdminTeamsList(context.Background(), AdminTeamsListParams{})
	assert.EqualError(t, err, "invalid_cursor")
	assert.NotNil(t, response)
}

func getAdminTeamsSettingsInfo(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	rw.Write([]byte(`{
		"ok": true,
		"team": {
			"id": "T12345",
			"name": "Test Workspace",
			"url": "https://test-workspace.slack.com/",
			"domain": "test-workspace",
			"email_domain": "example.com",
			"avatar_base_url": "https://ca.slack-edge.com/",
			"is_verified": false,
			"icon": {
				"image_default": true,
				"image_34": "https://example.com/icon_34.png",
				"image_44": "https://example.com/icon_44.png",
				"image_68": "https://example.com/icon_68.png",
				"image_88": "https://example.com/icon_88.png",
				"image_102": "https://example.com/icon_102.png",
				"image_132": "https://example.com/icon_132.png",
				"image_230": "https://example.com/icon_230.png"
			},
			"enterprise_id": "E12345",
			"enterprise_name": "Test Enterprise",
			"enterprise_domain": "test-enterprise",
			"default_channels": ["C12345", "C67890"]
		}
	}`))
}

func TestAdminTeamsSettingsInfo(t *testing.T) {
	http.HandleFunc("/admin.teams.settings.info", getAdminTeamsSettingsInfo)

	once.Do(startServer)
	api := New("testing-token", OptionAPIURL("http://"+serverAddr+"/"))

	settings, err := api.AdminTeamsSettingsInfo(context.Background(), "T12345")
	require.NoError(t, err)

	assert.Equal(t, "T12345", settings.ID)
	assert.Equal(t, "Test Workspace", settings.Name)
	assert.Equal(t, "https://test-workspace.slack.com/", settings.URL)
	assert.Equal(t, "test-workspace", settings.Domain)
	assert.Equal(t, "example.com", settings.EmailDomain)
	assert.Equal(t, "https://ca.slack-edge.com/", settings.AvatarBaseURL)
	assert.False(t, settings.IsVerified)
	assert.Equal(t, "E12345", settings.EnterpriseID)
	assert.Equal(t, "Test Enterprise", settings.EnterpriseName)
	assert.Equal(t, "test-enterprise", settings.EnterpriseDomain)
	assert.Equal(t, []string{"C12345", "C67890"}, settings.DefaultChannels)
	assert.True(t, settings.Icon.ImageDefault)
	assert.Equal(t, "https://example.com/icon_34.png", settings.Icon.Image34)
}

func okHandler(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	rw.Write([]byte(`{"ok": true}`))
}

func TestAdminTeamsSettingsSetDefaultChannels(t *testing.T) {
	http.HandleFunc("/admin.teams.settings.setDefaultChannels", okHandler)

	once.Do(startServer)
	api := New("testing-token", OptionAPIURL("http://"+serverAddr+"/"))

	err := api.AdminTeamsSettingsSetDefaultChannels(context.Background(), "T12345", "C111", "C222")
	require.NoError(t, err)
}

func TestAdminTeamsSettingsSetDescription(t *testing.T) {
	http.HandleFunc("/admin.teams.settings.setDescription", okHandler)

	once.Do(startServer)
	api := New("testing-token", OptionAPIURL("http://"+serverAddr+"/"))

	err := api.AdminTeamsSettingsSetDescription(context.Background(), "T12345", "A test workspace")
	require.NoError(t, err)
}

func TestAdminTeamsSettingsSetDiscoverability(t *testing.T) {
	http.HandleFunc("/admin.teams.settings.setDiscoverability", okHandler)

	once.Do(startServer)
	api := New("testing-token", OptionAPIURL("http://"+serverAddr+"/"))

	err := api.AdminTeamsSettingsSetDiscoverability(context.Background(), "T12345", TeamDiscoverabilityInviteOnly)
	require.NoError(t, err)
}

func TestAdminTeamsSettingsSetIcon(t *testing.T) {
	http.HandleFunc("/admin.teams.settings.setIcon", okHandler)

	once.Do(startServer)
	api := New("testing-token", OptionAPIURL("http://"+serverAddr+"/"))

	err := api.AdminTeamsSettingsSetIcon(context.Background(), "T12345", "https://example.com/icon.png")
	require.NoError(t, err)
}

func TestAdminTeamsSettingsSetName(t *testing.T) {
	http.HandleFunc("/admin.teams.settings.setName", okHandler)

	once.Do(startServer)
	api := New("testing-token", OptionAPIURL("http://"+serverAddr+"/"))

	err := api.AdminTeamsSettingsSetName(context.Background(), "T12345", "New Name")
	require.NoError(t, err)
}
