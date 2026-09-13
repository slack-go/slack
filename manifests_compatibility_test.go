package slack

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"sync/atomic"
	"testing"
)

func newManifestCompatibilityClient(t *testing.T, handler http.HandlerFunc) *Client {
	t.Helper()

	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)

	return New(
		"unused-token",
		OptionAPIURL(server.URL+"/"),
		OptionConfigToken("configured-token"),
	)
}

func TestCreateManifestWithOptions(t *testing.T) {
	manifest := &Manifest{Display: Display{Name: "example-app"}}
	wantManifest, err := json.Marshal(manifest)
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name         string
		existingCall bool
		options      []CreateManifestOption
		wantTeam     string
	}{
		{name: "existing method omits team", existingCall: true},
		{name: "empty", options: []CreateManifestOption{CreateManifestOptionTeamID("")}},
		{name: "present", options: []CreateManifestOption{CreateManifestOptionTeamID("T123")}, wantTeam: "T123"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			api := newManifestCompatibilityClient(t, func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path != "/apps.manifest.create" {
					t.Errorf("path = %q, want %q", r.URL.Path, "/apps.manifest.create")
				}
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
				if tt.wantTeam != "" && (!hasTeam || len(gotTeam) != 1 || gotTeam[0] != tt.wantTeam) {
					t.Errorf("team_id = %q, want %q", gotTeam, tt.wantTeam)
				}

				w.Header().Set("Content-Type", "application/json")
				_, _ = io.WriteString(w, `{
					"ok": true,
					"app_id": "A123",
					"credentials": {
						"client_id": "client-id",
						"client_secret": "client-secret",
						"verification_token": "verification-token",
						"signing_secret": "signing-secret"
					},
					"oauth_authorize_url": "https://slack.com/oauth/v2/authorize?client_id=client-id"
				}`)
			})

			var response *ManifestResponse
			var err error
			if tt.existingCall {
				response, err = api.CreateManifest(manifest, "per-call-token")
			} else {
				response, err = api.CreateManifestWithOptions(manifest, "per-call-token", tt.options...)
			}
			if err != nil {
				t.Fatalf("create manifest error = %v", err)
			}
			if response.AppID != "A123" {
				t.Errorf("AppID = %q, want %q", response.AppID, "A123")
			}
			wantCredentials := &ManifestCredentials{
				ClientID:          "client-id",
				ClientSecret:      "client-secret",
				VerificationToken: "verification-token",
				SigningSecret:     "signing-secret",
			}
			if !reflect.DeepEqual(response.Credentials, wantCredentials) {
				t.Errorf("Credentials = %#v, want %#v", response.Credentials, wantCredentials)
			}
			if response.OAuthAuthorizeURL != "https://slack.com/oauth/v2/authorize?client_id=client-id" {
				t.Errorf("OAuthAuthorizeURL = %q", response.OAuthAuthorizeURL)
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
			api := newManifestCompatibilityClient(t, func(w http.ResponseWriter, r *http.Request) {
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
	api := newManifestCompatibilityClient(t, func(w http.ResponseWriter, r *http.Request) {
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
	api := newManifestCompatibilityClient(t, func(w http.ResponseWriter, r *http.Request) {
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
	api := newManifestCompatibilityClient(t, func(w http.ResponseWriter, r *http.Request) {
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
			api := newManifestCompatibilityClient(t, func(w http.ResponseWriter, r *http.Request) {
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
		api := newManifestCompatibilityClient(t, func(w http.ResponseWriter, r *http.Request) {
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

func TestManifestTypedCompatibility(t *testing.T) {
	api := newManifestCompatibilityClient(t, func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			t.Errorf("ParseForm() error = %v", err)
		}
		if got := r.PostForm.Get("manifest"); got != "null" {
			t.Errorf("manifest = %q, want %q", got, "null")
		}
		if got := r.PostForm.Get("token"); got != "per-call-token" {
			t.Errorf("token = %q, want %q", got, "per-call-token")
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = io.WriteString(w, `{"ok":true}`)
	})

	var create func(*Manifest, string) (*ManifestResponse, error) = api.CreateManifest
	var createContext func(context.Context, *Manifest, string) (*ManifestResponse, error) = api.CreateManifestContext
	var update func(*Manifest, string, string) (*UpdateManifestResponse, error) = api.UpdateManifest
	var updateContext func(context.Context, *Manifest, string, string) (*UpdateManifestResponse, error) = api.UpdateManifestContext
	var export func(string, string) (*Manifest, error) = api.ExportManifest
	var exportContext func(context.Context, string, string) (*Manifest, error) = api.ExportManifestContext
	var validate func(*Manifest, string, string) (*ManifestResponse, error) = api.ValidateManifest
	var validateContext func(context.Context, *Manifest, string, string) (*ManifestResponse, error) = api.ValidateManifestContext
	_, _, _, _, _, _, _, _ = createContext, updateContext, export, exportContext, validateContext, create, update, validate

	if _, err := create(nil, "per-call-token"); err != nil {
		t.Fatalf("CreateManifest() error = %v", err)
	}
	if _, err := update(nil, "per-call-token", "A123"); err != nil {
		t.Fatalf("UpdateManifest() error = %v", err)
	}
	if _, err := validate(nil, "per-call-token", ""); err != nil {
		t.Fatalf("ValidateManifest() error = %v", err)
	}
}
