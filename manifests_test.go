package slack

import (
	"encoding/json"
	"net/http"
	"reflect"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateManifest(t *testing.T) {
	http.HandleFunc("/apps.manifest.create", handleCreateManifest)
	once.Do(startServer)

	api := New("testing-token", OptionAPIURL("http://"+serverAddr+"/"))

	manif := getTestManifest()
	resp, err := api.CreateManifest(&manif, "token")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
		return
	}

	if !reflect.DeepEqual(resp, getTestCreateManifestResponse()) {
		t.Fatal(ErrIncorrectResponse)
	}
}

func handleCreateManifest(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Content-Type", "application/json")

	// Shape from https://docs.slack.dev/reference/methods/apps.manifest.create
	rw.Write([]byte(`{
		"ok": true,
		"app_id": "A012ABCD0A0",
		"credentials": {
			"client_id": "1234567890.1234567890123",
			"client_secret": "abcdefghijklmnopqrstuvwxyz012345",
			"verification_token": "abcdefghijklmnopqrstuvwx",
			"signing_secret": "0123456789abcdef0123456789abcdef"
		},
		"oauth_authorize_url": "https://slack.com/oauth/v2/authorize?client_id=1234567890.1234567890123&scope=commands,workflow.steps:execute"
	}`))
}

func TestDeleteManifest(t *testing.T) {
	http.HandleFunc("/apps.manifest.delete", handleDeleteManifest)
	expectedResponse := SlackResponse{Ok: true}

	once.Do(startServer)
	api := New("testing-token", OptionAPIURL("http://"+serverAddr+"/"))

	resp, err := api.DeleteManifest("token", "app id")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
		return
	}

	if !reflect.DeepEqual(expectedResponse, *resp) {
		t.Fatal(ErrIncorrectResponse)
	}
}

func handleDeleteManifest(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Content-Type", "application/json")

	response, _ := json.Marshal(SlackResponse{Ok: true})
	rw.Write(response)
}

func TestExportManifest(t *testing.T) {
	http.HandleFunc("/apps.manifest.export", handleExportManifest)
	expectedResponse := getTestManifest()

	api := New("testing-token", OptionAPIURL("http://"+serverAddr+"/"))

	resp, err := api.ExportManifest("token", "app id")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
		return
	}

	if !reflect.DeepEqual(expectedResponse, *resp) {
		t.Fatal(ErrIncorrectResponse)
	}
}

func handleExportManifest(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Content-Type", "application/json")

	response, _ := json.Marshal(ExportManifestResponse{Manifest: getTestManifest()})
	rw.Write(response)
}

func TestUpdateManifest(t *testing.T) {
	http.HandleFunc("/apps.manifest.update", handleUpdateManifest)
	expectedResponse := UpdateManifestResponse{AppId: "app id"}

	api := New("testing-token", OptionAPIURL("http://"+serverAddr+"/"))

	manif := getTestManifest()
	resp, err := api.UpdateManifest(&manif, "token", "app id")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
		return
	}

	if !reflect.DeepEqual(expectedResponse, *resp) {
		t.Fatal(ErrIncorrectResponse)
	}
}

func handleUpdateManifest(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Content-Type", "application/json")

	response, _ := json.Marshal(UpdateManifestResponse{AppId: "app id"})
	rw.Write(response)
}

func TestValidateManifest(t *testing.T) {
	http.HandleFunc("/apps.manifest.validate", handleValidateManifest)
	expectedResponse := ManifestResponse{SlackResponse: SlackResponse{Ok: true}}

	api := New("testing-token", OptionAPIURL("http://"+serverAddr+"/"))

	manif := getTestManifest()
	resp, err := api.ValidateManifest(&manif, "token", "app id")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
		return
	}

	if !reflect.DeepEqual(expectedResponse, *resp) {
		t.Fatal(ErrIncorrectResponse)
	}
}

func handleValidateManifest(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Content-Type", "application/json")

	response, _ := json.Marshal(ManifestResponse{SlackResponse: SlackResponse{Ok: true}})
	rw.Write(response)
}

func getTestManifest() Manifest {
	return Manifest{
		Display: Display{
			Name:        "test",
			Description: "this is a test",
		},
		Features: Features{
			AgentView: &AgentView{
				AgentDescription: "this is a test agent",
				SuggestedPrompts: []ManifestSuggestedPrompt{
					{Title: "Test prompt", Message: "This is a test prompt"},
				},
				Actions: []ManifestAgentAction{
					{Name: "test_action", Description: "This is a test action"},
				},
			},
			AssistantView: &AssistantView{
				AssistantDescription: "this is a test assistant",
				SuggestedPrompts: []ManifestSuggestedPrompt{
					{Title: "Test prompt", Message: "This is a test prompt"},
				},
			},
		},
	}
}

func TestOAuthScopesOptionalFields(t *testing.T) {
	scopes := OAuthScopes{
		Bot:          []string{"chat:write", "commands"},
		User:         []string{"users:read"},
		BotOptional:  []string{"files:read", "reactions:read"},
		UserOptional: []string{"channels:read"},
	}

	data, err := json.Marshal(scopes)
	if err != nil {
		t.Fatalf("Marshal error: %s", err)
	}

	var roundtrip OAuthScopes
	if err := json.Unmarshal(data, &roundtrip); err != nil {
		t.Fatalf("Unmarshal error: %s", err)
	}

	if !reflect.DeepEqual(scopes, roundtrip) {
		t.Errorf("Round-trip mismatch: got %+v, want %+v", roundtrip, scopes)
	}

	// Verify omitempty: empty optional fields should not appear
	minimal := OAuthScopes{Bot: []string{"chat:write"}}
	data, err = json.Marshal(minimal)
	if err != nil {
		t.Fatalf("Marshal error: %s", err)
	}
	s := string(data)
	if strings.Contains(s, "bot_optional") {
		t.Errorf("Expected bot_optional to be omitted from JSON: %s", s)
	}
	if strings.Contains(s, "user_optional") {
		t.Errorf("Expected user_optional to be omitted from JSON: %s", s)
	}
}

func TestFeaturesAgentView(t *testing.T) {
	features := Features{
		AgentView: &AgentView{
			AgentDescription: "this is a test agent",
			SuggestedPrompts: []ManifestSuggestedPrompt{
				{Title: "Test prompt", Message: "This is a test prompt"},
			},
			Actions: []ManifestAgentAction{
				{Name: "test_action", Description: "This is a test action"},
			},
		},
	}

	data, err := json.Marshal(features)
	if err != nil {
		t.Fatalf("Marshal error: %s", err)
	}

	s := string(data)
	if !strings.Contains(s, `"agent_view"`) {
		t.Errorf("Expected agent_view to be present in JSON: %s", s)
	}
	if !strings.Contains(s, `"agent_description"`) {
		t.Errorf("Expected agent_description to be present in JSON: %s", s)
	}
	if !strings.Contains(s, `"suggested_prompts"`) {
		t.Errorf("Expected suggested_prompts to be present in JSON: %s", s)
	}
	if !strings.Contains(s, `"actions"`) {
		t.Errorf("Expected actions to be present in JSON: %s", s)
	}

	var roundtrip Features
	if err := json.Unmarshal(data, &roundtrip); err != nil {
		t.Fatalf("Unmarshal error: %s", err)
	}

	if !reflect.DeepEqual(features, roundtrip) {
		t.Errorf("Round-trip mismatch: got %+v, want %+v", roundtrip, features)
	}
}

func TestFeaturesAssistantView(t *testing.T) {
	features := Features{
		AssistantView: &AssistantView{
			AssistantDescription: "this is a test assistant",
			SuggestedPrompts: []ManifestSuggestedPrompt{
				{Title: "Test prompt", Message: "This is a test prompt"},
			},
		},
	}

	data, err := json.Marshal(features)
	if err != nil {
		t.Fatalf("Marshal error: %s", err)
	}

	s := string(data)
	if !strings.Contains(s, `"assistant_view"`) {
		t.Errorf("Expected assistant_view to be present in JSON: %s", s)
	}
	if !strings.Contains(s, `"assistant_description"`) {
		t.Errorf("Expected assistant_description to be present in JSON: %s", s)
	}
	if !strings.Contains(s, `"suggested_prompts"`) {
		t.Errorf("Expected suggested_prompts to be present in JSON: %s", s)
	}

	var roundtrip Features
	if err := json.Unmarshal(data, &roundtrip); err != nil {
		t.Fatalf("Unmarshal error: %s", err)
	}

	if !reflect.DeepEqual(features, roundtrip) {
		t.Errorf("Round-trip mismatch: got %+v, want %+v", roundtrip, features)
	}
}

func TestFeaturesAIViewsOmittedWhenNil(t *testing.T) {
	features := Features{
		BotUser: BotUser{
			DisplayName: "bot",
		},
	}

	data, err := json.Marshal(features)
	if err != nil {
		t.Fatalf("Marshal error: %s", err)
	}

	s := string(data)
	if strings.Contains(s, "agent_view") {
		t.Errorf("Expected agent_view to be omitted from JSON: %s", s)
	}
	if strings.Contains(s, "assistant_view") {
		t.Errorf("Expected assistant_view to be omitted from JSON: %s", s)
	}

	// Verify the description fields are always emitted, since Slack requires them
	withViews := Features{
		AgentView:     &AgentView{},
		AssistantView: &AssistantView{},
	}

	data, err = json.Marshal(withViews)
	if err != nil {
		t.Fatalf("Marshal error: %s", err)
	}

	s = string(data)
	if !strings.Contains(s, `"agent_description":""`) {
		t.Errorf("Expected agent_description to be present in JSON: %s", s)
	}
	if !strings.Contains(s, `"assistant_description":""`) {
		t.Errorf("Expected assistant_description to be present in JSON: %s", s)
	}
}

// TestManifestKeepsSettingsKeys checks that keys from an export are sent back on
// update instead of being dropped.
func TestManifestKeepsSettingsKeys(t *testing.T) {
	// Shape from https://docs.slack.dev/reference/app-manifest
	const settings = `{
		"event_subscriptions": {
			"metadata_subscriptions": [
				{"app_id": "A012ABCD0A0", "event_type": "task_created"}
			]
		},
		"function_runtime": "remote",
		"incoming_webhooks": {"incoming_webhooks_enabled": false},
		"org_deploy_enabled": true,
		"token_rotation_enabled": true
	}`
	// Slack accepts token_management_enabled=false only with org_deploy_enabled=true.
	payload := `{
		"display_information": {"name": "test"},
		"features": {"unfurl_domains": ["example.com"]},
		"oauth_config": {"token_management_enabled": false},
		"settings": ` + settings + `
	}`

	var manifest Manifest
	require.NoError(t, json.Unmarshal([]byte(payload), &manifest))

	assert.Equal(t, []string{"example.com"}, manifest.Features.UnfurlDomains)
	require.NotNil(t, manifest.OAuthConfig.TokenManagementEnabled)
	assert.False(t, *manifest.OAuthConfig.TokenManagementEnabled)
	assert.True(t, manifest.Settings.TokenRotationEnabled)
	assert.Equal(t, ManifestFunctionRuntimeRemote, manifest.Settings.FunctionRuntime)
	require.NotNil(t, manifest.Settings.IncomingWebhooks)
	assert.False(t, manifest.Settings.IncomingWebhooks.IncomingWebhooksEnabled)
	require.NotNil(t, manifest.Settings.EventSubscriptions)
	assert.Equal(t,
		[]ManifestMetadataSubscription{{AppID: "A012ABCD0A0", EventType: "task_created"}},
		manifest.Settings.EventSubscriptions.MetadataSubscriptions,
	)

	out, err := json.Marshal(manifest)
	require.NoError(t, err)
	var sent map[string]map[string]any
	require.NoError(t, json.Unmarshal(out, &sent))

	assert.Equal(t, []any{"example.com"}, sent["features"]["unfurl_domains"])
	assert.Equal(t, false, sent["oauth_config"]["token_management_enabled"])
	sentSettings, err := json.Marshal(sent["settings"])
	require.NoError(t, err)
	assert.JSONEq(t, settings, string(sentSettings))
}

func getTestCreateManifestResponse() *ManifestResponse {
	return &ManifestResponse{
		AppId: "A012ABCD0A0",
		Credentials: &ManifestCredentials{
			ClientId:          "1234567890.1234567890123",
			ClientSecret:      "abcdefghijklmnopqrstuvwxyz012345",
			VerificationToken: "abcdefghijklmnopqrstuvwx",
			SigningSecret:     "0123456789abcdef0123456789abcdef",
		},
		OAuthAuthorizeUrl: "https://slack.com/oauth/v2/authorize?client_id=1234567890.1234567890123&scope=commands,workflow.steps:execute",
		SlackResponse: SlackResponse{
			Ok: true,
		},
	}
}
