package slack

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"sync/atomic"
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
		"is_mcp_enabled": true,
		"org_deploy_enabled": true,
		"token_rotation_enabled": true
	}`
	// Slack accepts token_management_enabled=false only with org_deploy_enabled=true.
	payload := `{
		"display_information": {"name": "test"},
		"features": {"unfurl_domains": ["example.com"]},
		"oauth_config": {"pkce_enabled": true, "token_management_enabled": false},
		"settings": ` + settings + `
	}`

	var manifest Manifest
	require.NoError(t, json.Unmarshal([]byte(payload), &manifest))

	assert.Equal(t, []string{"example.com"}, manifest.Features.UnfurlDomains)
	require.NotNil(t, manifest.OAuthConfig.TokenManagementEnabled)
	assert.False(t, *manifest.OAuthConfig.TokenManagementEnabled)
	assert.True(t, manifest.OAuthConfig.PKCEEnabled)
	assert.True(t, manifest.Settings.IsMCPEnabled)
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
	assert.Equal(t, true, sent["oauth_config"]["pkce_enabled"])
	sentSettings, err := json.Marshal(sent["settings"])
	require.NoError(t, err)
	assert.JSONEq(t, settings, string(sentSettings))
}

// TestManifestKeepsFeatureKeys checks that rich_previews, search, siws_links and
// outgoing_domains from an export are sent back on update.
func TestManifestKeepsFeatureKeys(t *testing.T) {
	// Shape from https://docs.slack.dev/reference/app-manifest
	const richPreviews = `{"is_active": true, "entity_types": ["slack#/entities/task"]}`
	const search = `{"search_function_callback_id": "search", "search_filters_function_callback_id": "search_filters"}`
	const siwsLinks = `{"initiate_uri": "https://example.com/siws"}`
	payload := `{
		"display_information": {"name": "test"},
		"features": {"rich_previews": ` + richPreviews + `, "search": ` + search + `},
		"settings": {"siws_links": ` + siwsLinks + `},
		"outgoing_domains": ["api.example.com"]
	}`

	var manifest Manifest
	require.NoError(t, json.Unmarshal([]byte(payload), &manifest))

	require.NotNil(t, manifest.Features.RichPreviews)
	assert.True(t, manifest.Features.RichPreviews.IsActive)
	assert.Equal(t, []string{"slack#/entities/task"}, manifest.Features.RichPreviews.EntityTypes)
	require.NotNil(t, manifest.Features.Search)
	assert.Equal(t, "search", manifest.Features.Search.SearchFunctionCallbackID)
	assert.Equal(t, "search_filters", manifest.Features.Search.SearchFiltersFunctionCallbackID)
	require.NotNil(t, manifest.Settings.SIWSLinks)
	assert.Equal(t, "https://example.com/siws", manifest.Settings.SIWSLinks.InitiateURI)
	assert.Equal(t, []string{"api.example.com"}, manifest.OutgoingDomains)

	out, err := json.Marshal(manifest)
	require.NoError(t, err)
	var sent, features, settings map[string]json.RawMessage
	require.NoError(t, json.Unmarshal(out, &sent))
	require.NoError(t, json.Unmarshal(sent["features"], &features))
	require.NoError(t, json.Unmarshal(sent["settings"], &settings))

	assert.JSONEq(t, richPreviews, string(features["rich_previews"]))
	assert.JSONEq(t, search, string(features["search"]))
	assert.JSONEq(t, siwsLinks, string(settings["siws_links"]))
	assert.JSONEq(t, `["api.example.com"]`, string(sent["outgoing_domains"]))
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

func newManifestTestClient(t *testing.T, handler http.HandlerFunc) *Client {
	t.Helper()

	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)

	return New(
		"unused-token",
		OptionAPIURL(server.URL+"/"),
		OptionConfigToken("configured-token"),
	)
}

func TestCreateManifestTeamID(t *testing.T) {
	manifest := &Manifest{Display: Display{Name: "example-app"}}
	wantManifest, err := json.Marshal(manifest)
	require.NoError(t, err)

	tests := []struct {
		name     string
		options  []CreateManifestOption
		wantTeam string
	}{
		{name: "no option"},
		{name: "empty", options: []CreateManifestOption{CreateManifestOptionTeamID("")}},
		{name: "set", options: []CreateManifestOption{CreateManifestOptionTeamID("T123")}, wantTeam: "T123"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			api := newManifestTestClient(t, func(w http.ResponseWriter, r *http.Request) {
				if err := r.ParseForm(); err != nil {
					t.Errorf("ParseForm() error = %v", err)
				}
				if got := r.PostForm.Get("token"); got != "per-call-token" {
					t.Errorf("token = %q, want %q", got, "per-call-token")
				}
				if got := r.PostForm.Get("manifest"); got != string(wantManifest) {
					t.Errorf("manifest = %q, want %q", got, wantManifest)
				}
				gotTeam, hasTeam := r.PostForm["team_id"]
				if tt.wantTeam == "" && hasTeam {
					t.Errorf("team_id = %q, want omitted", gotTeam)
				}
				if tt.wantTeam != "" && r.PostForm.Get("team_id") != tt.wantTeam {
					t.Errorf("team_id = %q, want %q", gotTeam, tt.wantTeam)
				}
				w.Header().Set("Content-Type", "application/json")
				_, _ = io.WriteString(w, `{"ok":true}`)
			})

			if _, err := api.CreateManifest(manifest, "per-call-token", tt.options...); err != nil {
				t.Fatalf("CreateManifest() error = %v", err)
			}
		})
	}
}

func TestManifestRawRequestsPreserveJSON(t *testing.T) {
	rawManifest := json.RawMessage(" \n{\n  \"display_information\": {\"name\": \"example-app\", \"future_field\": {\"enabled\": true}},\n  \"future_root\": [1, {\"name\": \"value\"}]\n}\t ")

	tests := []struct {
		name       string
		path       string
		wantToken  string
		wantAppID  string
		wantTeamID string
		call       func(*Client) error
	}{
		{
			name:       "create",
			path:       "/apps.manifest.create",
			wantToken:  "per-call-token",
			wantTeamID: "T123",
			call: func(api *Client) error {
				_, err := api.CreateManifestRaw(rawManifest, "per-call-token", CreateManifestOptionTeamID("T123"))
				return err
			},
		},
		{
			name:      "update",
			path:      "/apps.manifest.update",
			wantToken: "configured-token",
			wantAppID: "A123",
			call: func(api *Client) error {
				_, err := api.UpdateManifestRaw(rawManifest, "", "A123")
				return err
			},
		},
		{
			name:      "validate",
			path:      "/apps.manifest.validate",
			wantToken: "configured-token",
			wantAppID: "A123",
			call: func(api *Client) error {
				_, err := api.ValidateManifestRaw(rawManifest, "", "A123")
				return err
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			api := newManifestTestClient(t, func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path != tt.path {
					t.Errorf("path = %q, want %q", r.URL.Path, tt.path)
				}
				if err := r.ParseForm(); err != nil {
					t.Errorf("ParseForm() error = %v", err)
				}
				if got := r.PostForm.Get("manifest"); got != string(rawManifest) {
					t.Errorf("manifest bytes changed:\ngot:  %q\nwant: %q", got, rawManifest)
				}
				if got := r.PostForm.Get("token"); got != tt.wantToken {
					t.Errorf("token = %q, want %q", got, tt.wantToken)
				}
				if got := r.PostForm.Get("app_id"); got != tt.wantAppID {
					t.Errorf("app_id = %q, want %q", got, tt.wantAppID)
				}
				if got := r.PostForm.Get("team_id"); got != tt.wantTeamID {
					t.Errorf("team_id = %q, want %q", got, tt.wantTeamID)
				}
				w.Header().Set("Content-Type", "application/json")
				_, _ = io.WriteString(w, `{"ok":true}`)
			})

			if err := tt.call(api); err != nil {
				t.Fatalf("raw manifest call error = %v", err)
			}
		})
	}
}

func TestExportManifestRawPreservesJSON(t *testing.T) {
	rawManifest := json.RawMessage(`{"_metadata":{"major_version":1},"future_root": { "nested": [1,{"enabled":true}] }}`)
	api := newManifestTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/apps.manifest.export" {
			t.Errorf("path = %q, want %q", r.URL.Path, "/apps.manifest.export")
		}
		if err := r.ParseForm(); err != nil {
			t.Errorf("ParseForm() error = %v", err)
		}
		if got := r.PostForm.Get("token"); got != "configured-token" {
			t.Errorf("token = %q, want %q", got, "configured-token")
		}
		if got := r.PostForm.Get("app_id"); got != "A123" {
			t.Errorf("app_id = %q, want %q", got, "A123")
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = io.WriteString(w, `{"ok":true,"manifest":`+string(rawManifest)+`}`)
	})

	got, err := api.ExportManifestRaw("", "A123")
	if err != nil {
		t.Fatalf("ExportManifestRaw() error = %v", err)
	}
	if string(got) != string(rawManifest) {
		t.Fatalf("manifest bytes changed:\ngot:  %q\nwant: %q", got, rawManifest)
	}
}

func TestExportManifestRawPreservesNull(t *testing.T) {
	api := newManifestTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = io.WriteString(w, `{"ok":true,"manifest":null}`)
	})

	got, err := api.ExportManifestRaw("token", "A123")
	if err != nil {
		t.Fatalf("ExportManifestRaw() error = %v", err)
	}
	if string(got) != "null" {
		t.Fatalf("manifest = %q, want %q", got, "null")
	}
}

func TestManifestRawRejectsNonObjectJSON(t *testing.T) {
	var requests atomic.Int32
	api := newManifestTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		requests.Add(1)
		_, _ = io.WriteString(w, `{"ok":true}`)
	})

	methods := []struct {
		name string
		call func(json.RawMessage) error
	}{
		{
			name: "create",
			call: func(manifest json.RawMessage) error {
				_, err := api.CreateManifestRaw(manifest, "token")
				return err
			},
		},
		{
			name: "update",
			call: func(manifest json.RawMessage) error {
				_, err := api.UpdateManifestRaw(manifest, "token", "A123")
				return err
			},
		},
		{
			name: "validate",
			call: func(manifest json.RawMessage) error {
				_, err := api.ValidateManifestRaw(manifest, "token", "")
				return err
			},
		},
	}
	manifests := []struct {
		name    string
		value   json.RawMessage
		wantErr string
	}{
		{name: "empty", wantErr: "valid JSON"},
		{name: "invalid", value: json.RawMessage(`{"display_information":`), wantErr: "valid JSON"},
		{name: "null", value: json.RawMessage(`null`), wantErr: "JSON object"},
		{name: "array", value: json.RawMessage(`[]`), wantErr: "JSON object"},
		{name: "string", value: json.RawMessage(`"manifest"`), wantErr: "JSON object"},
	}

	for _, method := range methods {
		for _, manifest := range manifests {
			t.Run(method.name+"/"+manifest.name, func(t *testing.T) {
				err := method.call(manifest.value)
				if err == nil || !strings.Contains(err.Error(), manifest.wantErr) {
					t.Fatalf("error = %v, want containing %q", err, manifest.wantErr)
				}
			})
		}
	}

	if got := requests.Load(); got != 0 {
		t.Fatalf("requests = %d, want 0", got)
	}
}

func TestValidateManifestRawAppIDAndErrors(t *testing.T) {
	rawManifest := json.RawMessage(`{"display_information":{"name":"example-app"}}`)

	for _, appID := range []string{"", "A123"} {
		name := "omitted"
		if appID != "" {
			name = "present"
		}
		t.Run(name, func(t *testing.T) {
			api := newManifestTestClient(t, func(w http.ResponseWriter, r *http.Request) {
				if err := r.ParseForm(); err != nil {
					t.Errorf("ParseForm() error = %v", err)
				}
				got, present := r.PostForm["app_id"]
				if appID == "" && present {
					t.Errorf("app_id = %q, want omitted", got)
				}
				if appID != "" && (!present || len(got) != 1 || got[0] != appID) {
					t.Errorf("app_id = %q, want %q", got, appID)
				}
				_, _ = io.WriteString(w, `{"ok":true,"errors":[]}`)
			})

			if _, err := api.ValidateManifestRaw(rawManifest, "token", appID); err != nil {
				t.Fatalf("ValidateManifestRaw() error = %v", err)
			}
		})
	}

	t.Run("validation error", func(t *testing.T) {
		api := newManifestTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_, _ = io.WriteString(w, `{
				"ok": false,
				"error": "invalid_manifest",
				"errors": [{
					"message": "Interactivity requires a Request URL",
					"pointer": "/settings/interactivity"
				}]
			}`)
		})

		response, err := api.ValidateManifestRaw(rawManifest, "token", "A123")
		if err == nil {
			t.Fatal("ValidateManifestRaw() error = nil, want Slack error")
		}
		var slackErr SlackErrorResponse
		if !errors.As(err, &slackErr) || slackErr.Err != "invalid_manifest" {
			t.Fatalf("error = %#v, want SlackErrorResponse for invalid_manifest", err)
		}
		wantErrors := []ManifestValidationError{{
			Message: "Interactivity requires a Request URL",
			Pointer: "/settings/interactivity",
		}}
		if !reflect.DeepEqual(response.Errors, wantErrors) {
			t.Fatalf("Errors = %#v, want %#v", response.Errors, wantErrors)
		}
	})
}
