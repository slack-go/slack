package slack

import (
	"net/http"
	"testing"
)

func getAuditLogs(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	response := []byte(`{"entries": [
			    {
			      "id": "0123a45b-6c7d-8900-e12f-3456789gh0i1",
			      "date_create": 1521214343,
			      "action": "user_login",
			      "actor": {
			        "type": "user",
			        "user": {
			          "id": "W123AB456",
			          "name": "Charlie Parker",
			          "email": "bird@slack.com"
			        }
			      },
			      "entity": {
			        "type": "user",
			        "user": {
			          "id": "W123AB456",
			          "name": "Charlie Parker",
			          "email": "bird@slack.com"
			        }
			      },
			      "context": {
			        "location": {
			          "type": "enterprise",
			          "id": "E1701NCCA",
			          "name": "Birdland",
			          "domain": "birdland"
			        },
			        "ua": "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/64.0.3282.186 Safari/537.36",
			        "ip_address": "1.23.45.678"
			      }
			    }
			  ]
  }`)
	rw.Write(response)
}

func TestGetAuditLogs(t *testing.T) {
	http.HandleFunc("/audit/v1/logs", getAuditLogs)

	once.Do(startServer)
	api := New("testing-token", OptionAPIURL("http://"+serverAddr+"/"))

	events, nextCursor, err := api.GetAuditLogs(AuditLogParameters{})
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
		return
	}

	if len(events) != 1 {
		t.Fatal("Should have been 1 event")
	}

	// test the first login
	event1 := events[0]

	if event1.Action != "user_login" {
		t.Fatal(ErrIncorrectResponse)
	}
	if event1.Entity.User.Email != "bird@slack.com" {
		t.Fatal(ErrIncorrectResponse)
	}
	if event1.Context.Location.Domain != "birdland" {
		t.Fatal(ErrIncorrectResponse)
	}
	if event1.DateCreate != 1521214343 {
		t.Fatal(ErrIncorrectResponse)
	}
	if event1.Context.IPAddress != "1.23.45.678" {
		t.Fatal(ErrIncorrectResponse)
	}

	if nextCursor != "" {
		t.Fatal(ErrIncorrectResponse)
	}
}
