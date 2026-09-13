package slack

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func newUsersLocaleClient(t *testing.T, handler http.HandlerFunc) *Client {
	t.Helper()

	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)

	return New("test-token", OptionAPIURL(server.URL+"/"))
}

func TestGetUserInfoIncludeLocale(t *testing.T) {
	tests := []struct {
		name       string
		wantLocale string
		wantUser   string
		wantUsers  string
		call       func(*Client) error
	}{
		{
			name:       "existing single-user method defaults true",
			wantLocale: "true",
			wantUser:   "U123",
			call: func(api *Client) error {
				_, err := api.GetUserInfo("U123")
				return err
			},
		},
		{
			name:       "new entry point defaults true",
			wantLocale: "true",
			wantUser:   "U123",
			call: func(api *Client) error {
				_, err := api.GetUserInfoWithOptions("U123")
				return err
			},
		},
		{
			name:       "new entry point accepts false",
			wantLocale: "false",
			wantUser:   "U123",
			call: func(api *Client) error {
				_, err := api.GetUserInfoWithOptions("U123", GetUserInfoOptionIncludeLocale(false))
				return err
			},
		},
		{
			name:       "existing multi-user method defaults true",
			wantLocale: "true",
			wantUsers:  "U123,U456",
			call: func(api *Client) error {
				_, err := api.GetUsersInfo("U123", "U456")
				return err
			},
		},
		{
			name:       "new multi-user entry point accepts false",
			wantLocale: "false",
			wantUsers:  "U123,U456",
			call: func(api *Client) error {
				_, err := api.GetUsersInfoWithOptions([]string{"U123", "U456"}, GetUserInfoOptionIncludeLocale(false))
				return err
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			api := newUsersLocaleClient(t, func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path != "/users.info" {
					t.Errorf("path = %q, want %q", r.URL.Path, "/users.info")
				}
				if err := r.ParseForm(); err != nil {
					t.Errorf("ParseForm() error = %v", err)
				}
				if got := r.PostForm.Get("include_locale"); got != tt.wantLocale {
					t.Errorf("include_locale = %q, want %q", got, tt.wantLocale)
				}
				if got := r.PostForm.Get("user"); got != tt.wantUser {
					t.Errorf("user = %q, want %q", got, tt.wantUser)
				}
				if got := r.PostForm.Get("users"); got != tt.wantUsers {
					t.Errorf("users = %q, want %q", got, tt.wantUsers)
				}
				w.Header().Set("Content-Type", "application/json")
				_, _ = io.WriteString(w, `{"ok":true,"user":{"id":"U123"},"users":[{"id":"U123"}]}`)
			})

			if err := tt.call(api); err != nil {
				t.Fatalf("users.info call error = %v", err)
			}
		})
	}
}

func TestGetUsersIncludeLocale(t *testing.T) {
	tests := []struct {
		name       string
		options    []GetUsersOption
		wantLocale string
		wantPages  int
	}{
		{name: "existing default", wantLocale: "true", wantPages: 1},
		{name: "explicit false", options: []GetUsersOption{GetUsersOptionIncludeLocale(false)}, wantLocale: "false", wantPages: 2},
		{name: "explicit true", options: []GetUsersOption{GetUsersOptionIncludeLocale(true)}, wantLocale: "true", wantPages: 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			requests := 0
			api := newUsersLocaleClient(t, func(w http.ResponseWriter, r *http.Request) {
				requests++
				if r.URL.Path != "/users.list" {
					t.Errorf("path = %q, want %q", r.URL.Path, "/users.list")
				}
				if err := r.ParseForm(); err != nil {
					t.Errorf("ParseForm() error = %v", err)
				}
				if got := r.PostForm.Get("include_locale"); got != tt.wantLocale {
					t.Errorf("include_locale = %q, want %q", got, tt.wantLocale)
				}
				w.Header().Set("Content-Type", "application/json")
				if requests < tt.wantPages {
					_, _ = io.WriteString(w, `{"ok":true,"members":[],"response_metadata":{"next_cursor":"next"}}`)
					return
				}
				_, _ = io.WriteString(w, `{"ok":true,"members":[],"response_metadata":{"next_cursor":""}}`)
			})

			if _, err := api.GetUsers(tt.options...); err != nil {
				t.Fatalf("GetUsers() error = %v", err)
			}
			if requests != tt.wantPages {
				t.Fatalf("requests = %d, want %d", requests, tt.wantPages)
			}
		})
	}
}

func TestUsersLocaleExistingSignatures(t *testing.T) {
	api := &Client{}

	var getUserInfo func(string) (*User, error) = api.GetUserInfo
	var getUserInfoContext func(context.Context, string) (*User, error) = api.GetUserInfoContext
	var getUsersInfo func(...string) (*[]User, error) = api.GetUsersInfo
	var getUsersInfoContext func(context.Context, ...string) (*[]User, error) = api.GetUsersInfoContext
	var getUsers func(...GetUsersOption) ([]User, error) = api.GetUsers
	var getUsersContext func(context.Context, ...GetUsersOption) ([]User, error) = api.GetUsersContext
	var getUsersPaginated func(...GetUsersOption) UserPagination = api.GetUsersPaginated

	_, _, _, _, _, _, _ = getUserInfo, getUserInfoContext, getUsersInfo, getUsersInfoContext, getUsers, getUsersContext, getUsersPaginated
}
