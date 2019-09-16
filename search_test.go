package slack

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestSearchMessagesSuccess(t *testing.T) {
	searchMessagesSuccessResponse := `{
    "messages": {
        "matches": [
            {
                "channel": {
                    "id": "C12345678",
                    "is_ext_shared": false,
                    "is_mpim": false,
                    "is_org_shared": false,
                    "is_pending_ext_shared": false,
                    "is_private": true,
                    "is_shared": false,
                    "name": "general",
                    "pending_shared": []
                },
                "iid": "cb64bdaa-c1e8-4631-8a91-0f78080113e9",
                "permalink": "https://hitchhikers.slack.com/archives/C12345678/p1508284197000015",
                "team": "T12345678",
                "text": "The meaning of life the universe and everything is 42.",
                "ts": "1508284197.000015",
                "type": "message",
                "user": "U2U85N1RV",
                "username": "roach"
            },
            {
                "channel": {
                    "id": "G024BE91L",
                    "is_ext_shared": false,
                    "is_mpim": true,
                    "is_org_shared": false,
                    "is_pending_ext_shared": false,
                    "is_private": false,
                    "is_shared": false,
                    "name": "mpdm-user1--user2--user3-1",
                    "pending_shared": []
                },
                "iid": "9a00d3c9-bd2d-45b0-988b-6cff99ae2a90",
                "permalink": "https://hitchhikers.slack.com/archives/C12345678/p1508795665000236",
                "team": "T12345678",
                "text": "The meaning of life the universe and everything is 101010",
                "ts": "1508795665.000236",
                "type": "message",
                "user": "",
                "username": "robot overlord"
            }
        ],
        "pagination": {
            "first": 1,
            "last": 2,
            "page": 1,
            "page_count": 1,
            "per_page": 20,
            "total_count": 2
        },
        "paging": {
            "count": 20,
            "page": 1,
            "pages": 1,
            "total": 2
        },
        "total": 2
    },
    "ok": true,
    "query": "The meaning of life the universe and everything"
}`

	expectedResponse := &SearchMessages{
		Matches: []SearchMessage{
			{
				Type: "message",
				IID:  "cb64bdaa-c1e8-4631-8a91-0f78080113e9",
				Channel: CtxChannel{
					ID:            "C12345678",
					Name:          "general",
					IsPrivate:     true,
					PendingShared: []interface{}{},
				},
				User:      "U2U85N1RV",
				Username:  "roach",
				Timestamp: "1508284197.000015",
				Text:      "The meaning of life the universe and everything is 42.",
				Permalink: "https://hitchhikers.slack.com/archives/C12345678/p1508284197000015",
			},
			{
				Type: "message",
				IID:  "9a00d3c9-bd2d-45b0-988b-6cff99ae2a90",
				Channel: CtxChannel{
					ID:            "G024BE91L",
					Name:          "mpdm-user1--user2--user3-1",
					IsMPIM:        true,
					PendingShared: []interface{}{},
				},
				Username:  "robot overlord",
				Timestamp: "1508795665.000236",
				Text:      "The meaning of life the universe and everything is 101010",
				Permalink: "https://hitchhikers.slack.com/archives/C12345678/p1508795665000236",
			},
		},
		Paging: Paging{
			Count: 20,
			Total: 2,
			Page:  1,
			Pages: 1,
		},
		Pagination: Pagination{
			TotalCount: 2,
			Page:       1,
			PerPage:    20,
			PageCount:  1,
			First:      1,
			Last:       2,
		},
		Total: 2,
	}

	http.HandleFunc("/search.messages", func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte(searchMessagesSuccessResponse))
		assert.NoError(t, err, "Writing the search.messages response should not result in an error")
	})
	once.Do(startServer)
	api := New(validToken, OptionAPIURL("http://"+serverAddr+"/"))
	resp, err := api.SearchMessages("pickleface", NewSearchParameters())
	assert.NoError(t, err, "The search.messages API should not error")
	assert.Equal(t, expectedResponse, resp, "Unexpected search.messages response")
}
