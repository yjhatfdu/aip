package cmd

import (
	"bytes"
	"strings"
	"testing"
)

func TestNormJSONDoesNotEscapeHTML(t *testing.T) {
	root := newRoot()
	root.SetArgs([]string{"norm", "--emit", "jsonl"})
	root.SetIn(strings.NewReader("2024-01-02T03:04:05Z test\n"))

	out := &bytes.Buffer{}
	root.SetOut(out)
	root.SetErr(&bytes.Buffer{})

	if err := root.Execute(); err != nil {
		t.Fatalf("norm error: %v", err)
	}

	got := out.String()
	if strings.Contains(got, "\\u003c") || strings.Contains(got, "\\u003e") {
		t.Fatalf("unexpected escaped angle brackets: %q", got)
	}
	if !strings.Contains(got, "<ts>") {
		t.Fatalf("expected <ts> in output: %q", got)
	}
}
