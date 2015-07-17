package slack

import "testing"

func TestNewRefToMessage(t *testing.T) {
	ref := NewRefToMessage("chan", "ts")
	if got, want := ref.ChannelId, "chan"; got != want {
		t.Errorf("ChannelId got %s, want %s", got, want)
	}
	if got, want := ref.Timestamp, "ts"; got != want {
		t.Errorf("Timestamp got %s, want %s", got, want)
	}
	if got, want := ref.FileId, ""; got != want {
		t.Errorf("FileId got %s, want %s", got, want)
	}
	if got, want := ref.CommentId, ""; got != want {
		t.Errorf("CommentId got %s, want %s", got, want)
	}
}

func TestNewRefToFile(t *testing.T) {
	ref := NewRefToFile("file")
	if got, want := ref.ChannelId, ""; got != want {
		t.Errorf("ChannelId got %s, want %s", got, want)
	}
	if got, want := ref.Timestamp, ""; got != want {
		t.Errorf("Timestamp got %s, want %s", got, want)
	}
	if got, want := ref.FileId, "file"; got != want {
		t.Errorf("FileId got %s, want %s", got, want)
	}
	if got, want := ref.CommentId, ""; got != want {
		t.Errorf("CommentId got %s, want %s", got, want)
	}
}

func TestNewRefToComment(t *testing.T) {
	ref := NewRefToComment("file_comment")
	if got, want := ref.ChannelId, ""; got != want {
		t.Errorf("ChannelId got %s, want %s", got, want)
	}
	if got, want := ref.Timestamp, ""; got != want {
		t.Errorf("Timestamp got %s, want %s", got, want)
	}
	if got, want := ref.FileId, ""; got != want {
		t.Errorf("FileId got %s, want %s", got, want)
	}
	if got, want := ref.CommentId, "file_comment"; got != want {
		t.Errorf("CommentId got %s, want %s", got, want)
	}
}
