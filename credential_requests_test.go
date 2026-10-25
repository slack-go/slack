package slack

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"testing"
)

func TestCredentialMethodsKeepSecretsOutOfURL(t *testing.T) {
	for _, tt := range []struct {
		name string
		path string
		form url.Values
		call func(*Client) error
	}{
		{
			name: "rotate explicit tokens",
			path: "/tooling.tokens.rotate",
			form: url.Values{"token": {"explicit-config"}, "refresh_token": {"refresh&token=+?"}},
			call: func(api *Client) error {
				_, err := api.RotateTokensContext(context.Background(), "explicit-config", "refresh&token=+?")
				return err
			},
		},
		{
			name: "rotate configured tokens",
			path: "/tooling.tokens.rotate",
			form: url.Values{"token": {"configured-config"}, "refresh_token": {"configured-refresh"}},
			call: func(api *Client) error {
				_, err := api.RotateTokens("", "")
				return err
			},
		},
		{
			name: "uninstall app",
			path: "/apps.uninstall",
			form: url.Values{"token": {"bot-token"}, "client_id": {"client-id"}, "client_secret": {"client&secret=+?"}},
			call: func(api *Client) error {
				return api.UninstallApp("client-id", "client&secret=+?")
			},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodPost {
					t.Errorf("method = %s, want POST", r.Method)
				}
				if r.URL.Path != tt.path || r.URL.RawQuery != "" {
					t.Errorf("request URL = %s, want %s without query parameters", r.URL, tt.path)
				}
				if err := r.ParseForm(); err != nil {
					t.Error(err)
				}
				if !reflect.DeepEqual(r.PostForm, tt.form) {
					t.Errorf("POST form = %v, want %v", r.PostForm, tt.form)
				}
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(getTestTokenResponse())
			}))
			defer server.Close()

			api := New("bot-token", OptionAPIURL(server.URL+"/"),
				OptionConfigToken("configured-config"), OptionConfigRefreshToken("configured-refresh"))
			if err := tt.call(api); err != nil {
				t.Fatal(err)
			}
		})
	}
}
