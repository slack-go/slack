package slack

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

func getTestBookmark(channelID, bookmarkID string) Bookmark {
	return Bookmark{
		ID:        bookmarkID,
		ChannelID: channelID,
		Title:     "bookmark",
		Type:      "link",
		Link:      "https://example.com",
		IconURL:   "https://example.com/icon.png",
	}
}

func addBookmarkLinkHandler(rw http.ResponseWriter, r *http.Request) {
	channelID := r.FormValue("channel_id")
	title := r.FormValue("title")
	bookmarkType := r.FormValue("type")
	link := r.FormValue("link")

	rw.Header().Set("Content-Type", "application/json")

	if bookmarkType == "link" && link != "" && channelID != "" && title != "" {
		bookmark := getTestBookmark(channelID, "Bk123RBZG8GZ")
		bookmark.Title = title
		bookmark.Type = bookmarkType
		bookmark.Link = link

		resp, _ := json.Marshal(&addBookmarkResponse{
			SlackResponse: SlackResponse{Ok: true},
			Bookmark:      bookmark})
		rw.Write(resp)
	} else {
		rw.Write([]byte(`{ "ok": false, "error": "errored" }`))
	}
}

func TestAddBookmarkLink(t *testing.T) {
	http.HandleFunc("/bookmarks.add", addBookmarkLinkHandler)
	once.Do(startServer)
	api := New("testing-token", OptionAPIURL("http://"+serverAddr+"/"))
	params := AddBookmarkParameters{
		Title: "test",
		Type:  "link",
		Link:  "https://example.com",
	}
	_, err := api.AddBookmark("CXXXXXXXX", params)
	if err != nil {
		t.Fatalf("Unexpected error: %s", err)
	}
}

func listBookmarksHandler(rw http.ResponseWriter, r *http.Request) {
	channelID := r.FormValue("channel_id")

	rw.Header().Set("Content-Type", "application/json")

	if channelID != "" {
		bookmarks := []Bookmark{
			getTestBookmark(channelID, "Bk001"),
			getTestBookmark(channelID, "Bk002"),
			getTestBookmark(channelID, "Bk003"),
			getTestBookmark(channelID, "Bk004"),
		}

		resp, _ := json.Marshal(&listBookmarksResponse{
			SlackResponse: SlackResponse{Ok: true},
			Bookmarks:     bookmarks})
		rw.Write(resp)
	} else {
		rw.Write([]byte(`{ "ok": false, "error": "errored" }`))
	}
}

func TestListBookmarks(t *testing.T) {
	http.HandleFunc("/bookmarks.list", listBookmarksHandler)
	once.Do(startServer)
	api := New("testing-token", OptionAPIURL("http://"+serverAddr+"/"))
	channel := "CXXXXXXXX"
	bookmarks, err := api.ListBookmarks(channel)
	if err != nil {
		t.Fatalf("Unexpected error: %s", err)
	}

	if !reflect.DeepEqual([]Bookmark{
		getTestBookmark(channel, "Bk001"),
		getTestBookmark(channel, "Bk002"),
		getTestBookmark(channel, "Bk003"),
		getTestBookmark(channel, "Bk004"),
	}, bookmarks) {
		t.Fatal(ErrIncorrectResponse)
	}
}

func removeBookmarkHandler(bookmark *Bookmark) func(rw http.ResponseWriter, r *http.Request) {
	return func(rw http.ResponseWriter, r *http.Request) {
		channelID := r.FormValue("channel_id")
		bookmarkID := r.FormValue("bookmark_id")

		rw.Header().Set("Content-Type", "application/json")

		if channelID == bookmark.ChannelID && bookmarkID == bookmark.ID {
			rw.Write([]byte(`{ "ok": true }`))
		} else {
			rw.Write([]byte(`{ "ok": false, "error": "errored" }`))
		}
	}
}

func TestRemoveBookmark(t *testing.T) {
	channel := "CXXXXXXXX"
	bookmark := getTestBookmark(channel, "BkXXXXX")
	http.HandleFunc("/bookmarks.remove", removeBookmarkHandler(&bookmark))
	once.Do(startServer)
	api := New("testing-token", OptionAPIURL("http://"+serverAddr+"/"))

	err := api.RemoveBookmark(channel, bookmark.ID)
	if err != nil {
		t.Fatalf("Unexpected error: %s", err)
	}
}

func editBookmarkHandler(bookmarks []Bookmark) func(rw http.ResponseWriter, r *http.Request) {
	return func(rw http.ResponseWriter, r *http.Request) {
		channelID := r.FormValue("channel_id")
		bookmarkID := r.FormValue("bookmark_id")

		rw.Header().Set("Content-Type", "application/json")
		if err := r.ParseForm(); err != nil {
			httpTestErrReply(rw, true, fmt.Sprintf("err parsing form: %s", err.Error()))
			return
		}

		for _, bookmark := range bookmarks {
			if bookmark.ID == bookmarkID && bookmark.ChannelID == channelID {
				if v := r.Form.Get("link"); v != "" {
					bookmark.Link = v
				}
				// Emoji and title require special handling since empty string sets to null
				if _, ok := r.Form["emoji"]; ok {
					bookmark.Emoji = r.Form.Get("emoji")
				}
				if _, ok := r.Form["title"]; ok {
					bookmark.Title = r.Form.Get("title")
				}
				resp, _ := json.Marshal(&editBookmarkResponse{
					SlackResponse: SlackResponse{Ok: true},
					Bookmark:      bookmark})
				rw.Write(resp)
				return
			}
		}
		// Fail if the bookmark doesn't exist
		rw.Write([]byte(`{ "ok": false, "error": "not_found" }`))
	}
}

func TestEditBookmark(t *testing.T) {
	channel := "CXXXXXXXX"
	bookmarks := []Bookmark{
		getTestBookmark(channel, "Bk001"),
		getTestBookmark(channel, "Bk002"),
		getTestBookmark(channel, "Bk003"),
		getTestBookmark(channel, "Bk004"),
	}
	http.HandleFunc("/bookmarks.edit", editBookmarkHandler(bookmarks))
	once.Do(startServer)
	api := New("testing-token", OptionAPIURL("http://"+serverAddr+"/"))

	smileEmoji := ":smile:"
	empty := ""
	title := "hello, world!"
	changes := []struct {
		ID     string
		Params EditBookmarkParameters
	}{
		{ // add emoji
			ID:     "Bk001",
			Params: EditBookmarkParameters{Emoji: &smileEmoji},
		},
		{ // delete emoji
			ID:     "Bk001",
			Params: EditBookmarkParameters{Emoji: &empty},
		},
		{ // add title
			ID:     "Bk002",
			Params: EditBookmarkParameters{Title: &title},
		},
		{ // delete title
			ID:     "Bk002",
			Params: EditBookmarkParameters{Title: &empty},
		},
		{ // Change multiple fields at once
			ID: "Bk003",
			Params: EditBookmarkParameters{
				Title: &title,
				Emoji: &empty,
				Link:  "https://example.com/changed",
			},
		},
		{ // noop
			ID: "Bk004",
		},
	}

	for _, change := range changes {
		bookmark, err := api.EditBookmark(channel, change.ID, change.Params)
		if err != nil {
			t.Fatalf("Unexpected error: %s", err)
		}
		if change.ID != bookmark.ID {
			t.Fatalf("expected to modify bookmark with ID = %s, got %s", change.ID, bookmark.ID)
		}
		if change.Params.Emoji != nil && bookmark.Emoji != *change.Params.Emoji {
			t.Fatalf("expected bookmark.Emoji = %s, got %s", *change.Params.Emoji, bookmark.Emoji)
		}
		if change.Params.Title != nil && bookmark.Title != *change.Params.Title {
			t.Fatalf("expected bookmark.Title = %s, got %s", *change.Params.Title, bookmark.Emoji)
		}
		if change.Params.Link != "" && change.Params.Link != bookmark.Link {
			t.Fatalf("expected bookmark.Link = %s, got %s", change.Params.Link, bookmark.Link)
		}
	}

	// Cover the final case of trying to edit a bookmark which doesn't exist
	bookmark, err := api.EditBookmark(channel, "BkMissing", EditBookmarkParameters{})
	if err == nil {
		t.Fatalf("Expected not found error, but got bookmark %s", bookmark.ID)
	}
}
