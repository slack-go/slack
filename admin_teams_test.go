package slack

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

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
