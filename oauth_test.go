package slack

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func getOpenIDConnectUserInfo(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	rw.Write([]byte(`{
		"ok": true,
		"sub": "U0R7JM",
		"https://slack.com/user_id": "U0R7JM",
		"https://slack.com/team_id": "T0R7GR",
		"email": "krane@slack-corp.com",
		"email_verified": true,
		"date_email_verified": 1622128723,
		"name": "krane",
		"picture": "https://secure.gravatar.com/....png",
		"given_name": "Bront",
		"family_name": "Kansen",
		"locale": "en-US",
		"https://slack.com/team_name": "Slack Corp",
		"https://slack.com/team_domain": "slackcorp",
		"https://slack.com/user_image_24": "...",
		"https://slack.com/user_image_32": "...",
		"https://slack.com/user_image_48": "...",
		"https://slack.com/user_image_72": "...",
		"https://slack.com/user_image_192": "...",
		"https://slack.com/user_image_512": "...",
		"https://slack.com/team_image_default": true
	}`))
}

func TestGetOpenIDConnectUserInfo(t *testing.T) {
	http.HandleFunc("/openid.connect.userInfo", getOpenIDConnectUserInfo)

	once.Do(startServer)
	api := New("testing-token", OptionAPIURL("http://"+serverAddr+"/"))

	resp, err := api.GetOpenIDConnectUserInfo()
	require.NoError(t, err)

	assert.Equal(t, "U0R7JM", resp.Sub)
	assert.Equal(t, "U0R7JM", resp.UserID)
	assert.Equal(t, "T0R7GR", resp.TeamID)
	assert.Equal(t, "krane@slack-corp.com", resp.Email)
	assert.True(t, resp.EmailVerified)
	assert.Equal(t, int64(1622128723), resp.DateEmailVerified)
	assert.Equal(t, "krane", resp.Name)
	assert.Equal(t, "Bront", resp.GivenName)
	assert.Equal(t, "Kansen", resp.FamilyName)
	assert.Equal(t, "en-US", resp.Locale)
	assert.Equal(t, "Slack Corp", resp.TeamName)
	assert.Equal(t, "slackcorp", resp.TeamDomain)
	assert.True(t, resp.TeamImageDefault)
}
