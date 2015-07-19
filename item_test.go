package slack

import "testing"

func TestNewMessageItem(t *testing.T) {
	c := "C1"
	m := &Message{}
	mi := NewMessageItem(c, m)
	if mi.Type != TYPE_MESSAGE {
		t.Errorf("want Type %s, got %s", mi.Type, TYPE_MESSAGE)
	}
	if mi.Channel != c {
		t.Errorf("got Channel %s, want %s", mi.Channel, c)
	}
	if mi.Message != m {
		t.Errorf("got Message %v, want %v", mi.Message, m)
	}
}

func TestNewFileItem(t *testing.T) {
	f := &File{}
	fi := NewFileItem(f)
	if fi.Type != TYPE_FILE {
		t.Errorf("got Type %s, want %s", fi.Type, TYPE_FILE)
	}
	if fi.File != f {
		t.Errorf("got File %v, want %v", fi.File, f)
	}
}

func TestNewFileCommentItem(t *testing.T) {
	f := &File{}
	c := &Comment{}
	fci := NewFileCommentItem(f, c)
	if fci.Type != TYPE_FILE_COMMENT {
		t.Errorf("got Type %s, want %s", fci.Type, TYPE_FILE_COMMENT)
	}
	if fci.File != f {
		t.Errorf("got File %v, want %v", fci.File, f)
	}
	if fci.Comment != c {
		t.Errorf("got Comment %v, want %v", fci.Comment, c)
	}
}

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
