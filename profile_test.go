package slack

import (
	"encoding/json"
	"net/http"
	"testing"
)

func setUserProfile(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	response, _ := json.Marshal(struct {
		Ok bool `json:"ok"`
	}{
		Ok: true,
	})
	rw.Write(response)
}

func TestSetUserProfile(t *testing.T) {
	http.HandleFunc("/users.profile.set", setUserProfile)
	profile := getTestUserProfile()

	once.Do(startServer)
	api := New("testing-token", OptionAPIURL("http://"+serverAddr+"/"))

	err := api.SetUserProfile("UXXXXXXXX", &profile)
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
		return
	}
}
