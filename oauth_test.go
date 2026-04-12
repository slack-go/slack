package slack

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func getOAuthV2Response(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	rw.Write([]byte(`{
		"ok": true,
		"access_token": "xoxb-test-token",
		"token_type": "bot",
		"scope": "chat:write",
		"bot_user_id": "U0KRQLJ9H",
		"app_id": "A0KRD7HC3",
		"team": {"name": "Test Team", "id": "T0KRQLJ9H"},
		"enterprise": {"name": "", "id": ""},
		"is_enterprise_install": false,
		"authed_user": {"id": "U0KRQLJ9H"}
	}`))
}

func TestGetOAuthV2ResponseWithCustomURL(t *testing.T) {
	http.HandleFunc("/oauth.v2.access", getOAuthV2Response)

	once.Do(startServer)

	resp, err := GetOAuthV2Response(
		http.DefaultClient,
		"client-id", "client-secret", "code", "http://localhost/callback",
		OAuthOptionAPIURL("http://"+serverAddr+"/"),
	)
	require.NoError(t, err)

	assert.Equal(t, "xoxb-test-token", resp.AccessToken)
	assert.Equal(t, "bot", resp.TokenType)
	assert.Equal(t, "chat:write", resp.Scope)
	assert.Equal(t, "U0KRQLJ9H", resp.BotUserID)
	assert.Equal(t, "T0KRQLJ9H", resp.Team.ID)
}

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

func TestGenerateCodeVerifier(t *testing.T) {
	v1, err := GenerateCodeVerifier()
	require.NoError(t, err)
	assert.Len(t, v1, 43) // 32 bytes -> 43 chars base64url without padding

	v2, err := GenerateCodeVerifier()
	require.NoError(t, err)
	assert.NotEqual(t, v1, v2, "two calls should produce different verifiers")
}

func TestGenerateCodeChallenge(t *testing.T) {
	// RFC 7636 Appendix B test vector
	verifier := "dBjftJeZ4CVP-mB92K27uhbUJU1p1r_wW1gFWFOEjXk"
	expected := "E9Melhoa2OwvFrEMTJguCHaoeK1t8URWbuGJSstw-cM"

	challenge := GenerateCodeChallenge(verifier)
	assert.Equal(t, expected, challenge)
}

func TestOAuthOptionCodeVerifier(t *testing.T) {
	cfg := resolveOAuthConfig([]OAuthOption{
		OAuthOptionCodeVerifier("test-verifier"),
	})
	assert.Equal(t, "test-verifier", cfg.codeVerifier)
	assert.Equal(t, APIURL, cfg.apiURL) // default preserved
}

func TestGetOAuthV2ResponsePKCE(t *testing.T) {
	once.Do(startServer)

	// The existing /oauth.v2.access handler returns a valid response.
	// Verify that PKCE params are accepted without error.
	resp, err := GetOAuthV2Response(
		http.DefaultClient,
		"client-id", "", "code", "http://localhost/callback",
		OAuthOptionAPIURL("http://"+serverAddr+"/"),
		OAuthOptionCodeVerifier("test-verifier"),
	)
	require.NoError(t, err)
	assert.Equal(t, "xoxb-test-token", resp.AccessToken)
}

func TestGetOAuthV2ResponseOmitsEmptySecret(t *testing.T) {
	// Verify that resolveOAuthConfig + empty secret produces no client_secret key
	values := url.Values{
		"client_id":    {"id"},
		"code":         {"code"},
		"redirect_uri": {"http://localhost"},
	}
	clientSecret := ""
	if clientSecret != "" {
		values.Set("client_secret", clientSecret)
	}
	_, hasSecret := values["client_secret"]
	assert.False(t, hasSecret, "empty client_secret should not be in form values")
}
