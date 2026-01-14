package cmd

import (
	"bytes"
	"path/filepath"
	"strings"
	"testing"
)

func TestConfigSetAndGet(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("USERPROFILE", home)
	t.Setenv("AIP_LANG", "en")

	root := newRoot()
	root.SetArgs([]string{"config", "set", "base_url", "https://example.com"})
	root.SetOut(&bytes.Buffer{})
	root.SetErr(&bytes.Buffer{})
	if err := root.Execute(); err != nil {
		t.Fatalf("config set error: %v", err)
	}

	root = newRoot()
	root.SetArgs([]string{"config", "get", "base_url"})
	out := &bytes.Buffer{}
	root.SetOut(out)
	root.SetErr(&bytes.Buffer{})
	if err := root.Execute(); err != nil {
		t.Fatalf("config get error: %v", err)
	}
	got := strings.TrimSpace(out.String())
	if got != "https://example.com" {
		t.Fatalf("config get = %q", got)
	}
}

func TestConfigPath(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("USERPROFILE", home)
	t.Setenv("AIP_LANG", "en")

	root := newRoot()
	root.SetArgs([]string{"config", "path"})
	out := &bytes.Buffer{}
	root.SetOut(out)
	root.SetErr(&bytes.Buffer{})
	if err := root.Execute(); err != nil {
		t.Fatalf("config path error: %v", err)
	}
	got := strings.TrimSpace(out.String())
	want := filepath.Join(home, ".aip", "config.toml")
	if got != want {
		t.Fatalf("config path = %q want %q", got, want)
	}
}

func TestConfigGetUnknownKey(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("USERPROFILE", home)
	t.Setenv("AIP_LANG", "en")

	root := newRoot()
	root.SetArgs([]string{"config", "get", "unknown"})
	root.SetOut(&bytes.Buffer{})
	root.SetErr(&bytes.Buffer{})
	err := root.Execute()
	if err == nil {
		t.Fatal("expected error for unknown key")
	}
}
