package slack

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"testing"
)

func TestMigrationExchangeRequest(t *testing.T) {
	tests := []struct {
		name   string
		users  []string
		teamID string
		toOld  bool
		want   url.Values
	}{
		{
			name:  "single_user",
			users: []string{"U1"},
			want:  url.Values{"users": {"U1"}},
		},
		{
			name:   "multiple_users",
			users:  []string{"U1", "U2"},
			teamID: "T1",
			toOld:  true,
			want: url.Values{
				"users":   {"U1,U2"},
				"team_id": {"T1"},
				"to_old":  {"true"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodGet || r.URL.Path != "/migration.exchange" {
					t.Errorf("request = %s %s, want GET /migration.exchange", r.Method, r.URL.Path)
				}
				if got := r.URL.Query(); !reflect.DeepEqual(got, tt.want) {
					t.Errorf("query = %#v, want %#v", got, tt.want)
				}
				w.Header().Set("Content-Type", "application/json")
				_, _ = w.Write([]byte(`{"ok":true}`))
			}))
			defer server.Close()

			api := New("testing-token", OptionAPIURL(server.URL+"/"))
			if _, _, err := api.MigrationExchange(t.Context(), tt.teamID, tt.toOld, tt.users); err != nil {
				t.Fatalf("MigrationExchange: %v", err)
			}
		})
	}
}
