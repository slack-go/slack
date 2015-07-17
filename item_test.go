package slack

import "testing"

func TestNewMessageItem(t *testing.T) {
	m := &Message{}
	mi := NewMessageItem(m)
	if mi.Type != TYPE_MESSAGE {
		t.Errorf("want Type %s, got %s", TYPE_MESSAGE, mi.Type)
	}
	if m != mi.Message {
		t.Errorf("want Message %v, got %v", m, mi.Message)
	}
}

func TestNewFileItem(t *testing.T) {
	f := &File{}
	fi := NewFileItem(f)
	if fi.Type != TYPE_FILE {
		t.Errorf("want Type %s, got %s", TYPE_FILE, fi.Type)
	}
	if f != fi.File {
		t.Errorf("want File %v, got %v", f, fi.File)
	}
}

func TestNewCommentItem(t *testing.T) {
	c := &Comment{}
	ci := NewCommentItem(c)
	if ci.Type != TYPE_COMMENT {
		t.Errorf("want Type %s, got %s", TYPE_COMMENT, ci.Type)
	}
	if c != ci.Comment {
		t.Errorf("want Comment %v, got %v", c, ci.Comment)
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
