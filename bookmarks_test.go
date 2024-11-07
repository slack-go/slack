package slack

import (
	"net/http"
	"testing"
)

var listBookmarksResp = `{
  "ok": true,
  "bookmarks": [
    {
      "id": "Bk12345",
      "channel_id": "C12345",
      "title": "Homepage",
      "link": "https://app.incident.io/incidents/12",
      "emoji": ":globe_with_meridians:",
      "icon_url": null,
      "type": "link",
      "entity_id": null,
      "date_created": 1644767331,
      "date_updated": 0,
      "rank": "U",
      "last_updated_by_user_id": "U12345",
      "last_updated_by_team_id": "T12345",
      "shortcut_id": null,
      "app_id": null
    }
  ]
}`

func listBookmarks(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	rw.Write([]byte(listBookmarksResp))
}

func TestBookmarkList(t *testing.T) {
	http.HandleFunc("/bookmarks.list", listBookmarks)
	once.Do(startServer)

	api := New("testing-token", OptionAPIURL("http://"+serverAddr+"/"))
	bookmarks, err := api.ListBookmarks("C12345")
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
		return
	}

	if len(bookmarks) != 1 {
		t.Fatalf("expected 1 bookmark, got %d", len(bookmarks))
		return
	}

	bookmark := bookmarks[0]

	if bookmark.ID != "Bk12345" {
		t.Fatalf("expected bookmark ID Bk12345, got %s", bookmark.ID)
	}
}

var singleBookmarkResp = `{
  "ok": true,
  "bookmark": {
    "id": "Bk12345",
    "channel_id": "C12345",
    "title": "Homepage",
    "link": "https://app.incident.io/incidents/12",
    "emoji": ":globe_with_meridians:",
    "icon_url": null,
    "type": "link",
    "entity_id": null,
    "date_created": 1644767331,
    "date_updated": 0,
    "rank": "U",
    "last_updated_by_user_id": "U12345",
    "last_updated_by_team_id": "T12345",
    "shortcut_id": null,
    "app_id": null
  }
}`

func getBookmark(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	rw.Write([]byte(singleBookmarkResp))
}

func TestAddBookmark(t *testing.T) {
	http.HandleFunc("/bookmarks.add", getBookmark)
	once.Do(startServer)

	api := New("testing-token", OptionAPIURL("http://"+serverAddr+"/"))
	bookmark, err := api.AddBookmark("C12345", AddBookmarkParameters{
		Title: "Homepage",
		Type:  "link",
	})

	if err != nil {
		t.Errorf("Unexpected error: %s", err)
		return
	}

	if bookmark.ID != "Bk12345" {
		t.Errorf("bookmark ID should be Bk12345, got %s", bookmark.ID)
	}
}

func TestEditBookmark(t *testing.T) {
	http.HandleFunc("/bookmarks.edit", getBookmark)
	once.Do(startServer)

	emoji := ":siren:"
	title := "hello2"

	api := New("testing-token", OptionAPIURL("http://"+serverAddr+"/"))
	bookmark, err := api.EditBookmark("C12345", "Bk12345", EditBookmarkParameters{
		Emoji: &emoji,
		Title: &title,
	})

	if err != nil {
		t.Errorf("Unexpected error: %s", err)
		return
	}

	if bookmark.ID != "Bk12345" {
		t.Errorf("bookmark ID should be Bk12345, got %s", bookmark.ID)
	}
}

func okResponse(rw http.ResponseWriter, r *http.Request) {
	rw.Write([]byte(`{"ok": true}`))
}

func TestRemoveBookmark(t *testing.T) {
	http.HandleFunc("/bookmarks.remove", getBookmark)
	once.Do(startServer)

	api := New("testing-token", OptionAPIURL("http://"+serverAddr+"/"))
	err := api.RemoveBookmark("C12345", "Bk12345")

	if err != nil {
		t.Errorf("Unexpected error: %s", err)
		return
	}
}
