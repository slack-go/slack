package slack

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestUnMuteChat(t *testing.T) {
	for _, tt := range []struct {
		name       string
		muted      string
		channel    string
		wantMuted  string
		wantUpdate bool
	}{
		{"first", "C1,C2,C3", "C1", "C2,C3", true},
		{"middle", "C1,C2,C3", "C2", "C1,C3", true},
		{"last", "C1,C2,C3", "C3", "C1,C2", true},
		{"only", "C1", "C1", "", true},
		{"not muted", "C1,C2", "C3", "C1,C2", false},
		{"none muted", "", "C1", "", false},
	} {
		t.Run(tt.name, func(t *testing.T) {
			updates := make(chan url.Values, 1)
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				muted := tt.muted
				switch r.URL.Path {
				case "/users.prefs.get":
				case "/users.prefs.set":
					if err := r.ParseForm(); err != nil {
						t.Error(err)
						http.Error(w, err.Error(), http.StatusBadRequest)
						return
					}
					updates <- r.PostForm
					muted = r.PostForm.Get("muted_channels")
				default:
					t.Errorf("unexpected request to %s", r.URL.Path)
					http.NotFound(w, r)
					return
				}
				if err := json.NewEncoder(w).Encode(UserPrefsCarrier{
					SlackResponse: SlackResponse{Ok: true},
					UserPrefs:     &UserPrefs{MutedChannels: muted},
				}); err != nil {
					t.Error(err)
				}
			}))
			defer server.Close()

			api := New("testing-token", OptionAPIURL(server.URL+"/"))
			response, err := api.UnMuteChat(tt.channel)
			if err != nil {
				t.Fatal(err)
			}
			select {
			case values := <-updates:
				if !tt.wantUpdate {
					t.Fatal("unexpected preferences update for an unmuted channel")
				}
				if !values.Has("muted_channels") || values.Get("muted_channels") != tt.wantMuted {
					t.Errorf("muted_channels = %v, want %q", values["muted_channels"], tt.wantMuted)
				}
				if response == nil || response.UserPrefs == nil || response.UserPrefs.MutedChannels != tt.wantMuted {
					t.Errorf("unexpected preferences response: %+v", response)
				}
			default:
				if tt.wantUpdate {
					t.Fatal("no preferences update sent for a muted channel")
				}
				if response != nil {
					t.Errorf("expected nil response for an unmuted channel, got %+v", response)
				}
			}
		})
	}
}

func TestJSONTime_UnmarshalJSON(t *testing.T) {
	type args struct {
		buf []byte
	}
	tests := []struct {
		name    string
		args    args
		wantTr  JSONTime
		wantErr bool
	}{
		{
			"acceptable int64 timestamp",
			args{[]byte(`1643435556`)},
			JSONTime(1643435556),
			false,
		},
		{
			"acceptable string timestamp",
			args{[]byte(`"1643435556"`)},
			JSONTime(1643435556),
			false,
		},
		{
			"null",
			args{[]byte(`null`)},
			JSONTime(0),
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var tr JSONTime
			if err := tr.UnmarshalJSON(tt.args.buf); (err != nil) != tt.wantErr {
				t.Errorf("JSONTime.UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
