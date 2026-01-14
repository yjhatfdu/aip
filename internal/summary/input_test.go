package summary

import (
	"os"
	"strings"
	"testing"
)

func TestReadInputMaxChars(t *testing.T) {
	input := strings.Repeat("a", 10)
	got, err := ReadInput(strings.NewReader(input), InputOptions{MaxChars: 5})
	if err != nil {
		t.Fatalf("ReadInput error: %v", err)
	}
	if got != "aaaaa" {
		t.Fatalf("got %q", got)
	}
}

func TestReadInputHeadTail(t *testing.T) {
	input := "0123456789"
	got, err := ReadInput(strings.NewReader(input), InputOptions{IncludeHead: 3, IncludeTail: 3, MaxChars: 100})
	if err != nil {
		t.Fatalf("ReadInput error: %v", err)
	}
	if !strings.HasPrefix(got, "012") || !strings.HasSuffix(got, "789") {
		t.Fatalf("got %q", got)
	}
	if !strings.Contains(got, "...") {
		t.Fatalf("expected ellipsis: %q", got)
	}
}

func TestLoadPromptFromFile(t *testing.T) {
	dir := t.TempDir()
	path := dir + "/p.txt"
	if err := os.WriteFile(path, []byte("hello\n"), 0o644); err != nil {
		t.Fatalf("write: %v", err)
	}
	got, err := LoadPrompt("@" + path)
	if err != nil {
		t.Fatalf("LoadPrompt: %v", err)
	}
	if got != "hello" {
		t.Fatalf("got %q", got)
	}
}
