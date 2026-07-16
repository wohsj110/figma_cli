package figmaurl

import "testing"

func TestParseRawKey(t *testing.T) {
	got, err := Parse("abc123")
	if err != nil {
		t.Fatal(err)
	}
	if got.FileKey != "abc123" || got.NodeID != "" {
		t.Fatalf("unexpected target: %#v", got)
	}
}

func TestParseFigmaURL(t *testing.T) {
	got, err := Parse("https://www.figma.com/design/AbCdEf/My-File?node-id=12-34&m=dev")
	if err != nil {
		t.Fatal(err)
	}
	if got.FileKey != "AbCdEf" {
		t.Fatalf("FileKey = %q", got.FileKey)
	}
	if got.NodeID != "12:34" {
		t.Fatalf("NodeID = %q", got.NodeID)
	}
}

func TestParseInvalidURL(t *testing.T) {
	if _, err := Parse("https://www.figma.com/community/file/123"); err == nil {
		t.Fatal("expected error")
	}
}
