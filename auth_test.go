package slack

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func getTeamList(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	response := []byte(`{
    "ok": true,
    "teams": [
        {
            "name": "Shinichi's workspace",
            "id": "T12345678"
        },
        {
            "name": "Migi's workspace",
            "id": "T12345679"
        }
    ],
    "response_metadata": {
        "next_cursor": "dXNlcl9pZDo5MTQyOTI5Mzkz"
    }
}`)
	rw.Write(response)
}

func TestListTeams(t *testing.T) {
	http.HandleFunc("/auth.teams.list", getTeamList)

	once.Do(startServer)
	api := New("testing-token", OptionAPIURL("http://"+serverAddr+"/"))

	teams, cursor, err := api.ListTeams(ListTeamsParameters{})
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
		return
	}

	require.Len(t, teams, 2)
	assert.Equal(t, "T12345678", teams[0].ID)
	assert.Equal(t, "Shinichi's workspace", teams[0].Name)

	assert.Equal(t, "T12345679", teams[1].ID)
	assert.Equal(t, "Migi's workspace", teams[1].Name)

	assert.Equal(t, "dXNlcl9pZDo5MTQyOTI5Mzkz", cursor)
}

func TestListTeamsLimit(t *testing.T) {
	methods := []struct {
		name string
		call func(*Client, ListTeamsParameters) ([]Team, string, error)
	}{
		{"ListTeams", (*Client).ListTeams},
		{"ListTeamsContext", func(api *Client, params ListTeamsParameters) ([]Team, string, error) {
			return api.ListTeamsContext(context.Background(), params)
		}},
	}
	for _, method := range methods {
		t.Run(method.name, func(t *testing.T) {
			for _, tc := range []struct {
				name  string
				limit int
				want  []string
			}{
				{"zero omits limit", 0, nil},
				{"minimum", 1, []string{"1"}},
				{"page size", 50, []string{"50"}},
				{"maximum", 1000, []string{"1000"}},
			} {
				t.Run(tc.name, func(t *testing.T) {
					requests := make(chan url.Values, 1)
					server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
						assert.Equal(t, http.MethodPost, r.Method)
						assert.Equal(t, "/auth.teams.list", r.URL.Path)
						if err := r.ParseForm(); err != nil {
							http.Error(w, err.Error(), http.StatusBadRequest)
							return
						}
						requests <- r.PostForm
						w.Header().Set("Content-Type", "application/json")
						w.Write([]byte(`{"ok":true,"teams":[],"response_metadata":{"next_cursor":"next"}}`))
					}))
					t.Cleanup(server.Close)
					api := New("testing-token", OptionAPIURL(server.URL+"/"))
					includeIcon := false
					_, cursor, err := method.call(api, ListTeamsParameters{
						Limit: tc.limit, Cursor: "current", IncludeIcon: &includeIcon,
					})
					require.NoError(t, err)
					assert.Equal(t, "next", cursor)
					want := url.Values{
						"token": {"testing-token"}, "cursor": {"current"}, "include_icon": {"false"},
					}
					if tc.want != nil {
						want["limit"] = tc.want
					}
					require.Len(t, requests, 1)
					assert.Equal(t, want, <-requests)
				})
			}
		})
	}
}
