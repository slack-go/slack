package slack

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestListEventAuthorizationsDeliversWarnings(t *testing.T) {
	var warnings int
	api := New("test-token",
		OptionAppLevelToken("test-app-token"),
		OptionHTTPClient(httpClientFunc(func(r *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(strings.NewReader(`{"ok":true,"warning":"test_warning","authorizations":[{"team_id":"T1","user_id":"U1"}]}`)),
			}, nil
		})),
		OptionOnWarning(func(path string, request any, warning *Warning) {
			warnings++
			require.Equal(t, "apps.event.authorizations.list", path)
		}),
	)
	authorizations, err := api.ListEventAuthorizationsContext(t.Context(), "event-context")
	require.NoError(t, err)
	require.Len(t, authorizations, 1)
	require.Equal(t, "U1", authorizations[0].UserID)
	require.Equal(t, 1, warnings)
}

func TestListEventAuthorizations(t *testing.T) {
	http.HandleFunc("/apps.event.authorizations.list", testListEventAuthorizationsHandler)
	once.Do(startServer)

	api := New("", OptionAppLevelToken("test-token"), OptionAPIURL("http://"+serverAddr+"/"))

	authorizations, err := api.ListEventAuthorizations("1-message-T012345678-DR12345678")

	switch {
	case err != nil:
		t.Errorf("Failed, but should have succeeded")
	case len(authorizations) != 1:
		t.Errorf("Didn't get 1 authorization")
	case authorizations[0].UserID != "U123456789":
		t.Errorf("User ID is wrong")
	}
}

func testListEventAuthorizationsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	response, _ := json.Marshal(listEventAuthorizationsResponse{
		SlackResponse: SlackResponse{Ok: true},
		Authorizations: []EventAuthorization{
			{
				UserID: "U123456789",
				TeamID: "T012345678",
			},
		},
	})
	w.Write(response)
}

func TestUninstallApp(t *testing.T) {
	http.HandleFunc("/apps.uninstall", testUninstallAppHandler)
	once.Do(startServer)

	api := New("test-token", OptionAPIURL("http://"+serverAddr+"/"))

	err := api.UninstallApp("", "")

	if err != nil {
		t.Errorf("Failed, but should have succeeded")
	}
}

func testUninstallAppHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	response, _ := json.Marshal(SlackResponse{Ok: true})
	w.Write(response)
}
