package config

import (
	"bytes"
	"path/filepath"
	"strings"
	"testing"
)

func TestRenderAndParseTOML(t *testing.T) {
	cfg := Config{
		BaseURL: "https://example.com",
		APIKey:  "key-123",
		Model:   "gpt-test",
	}
	out := renderTOML(cfg)
	parsed := parseTOML(strings.NewReader(out))

	if parsed.BaseURL != cfg.BaseURL {
		t.Fatalf("BaseURL mismatch: got %q want %q", parsed.BaseURL, cfg.BaseURL)
	}
	if parsed.APIKey != cfg.APIKey {
		t.Fatalf("APIKey mismatch: got %q want %q", parsed.APIKey, cfg.APIKey)
	}
	if parsed.Model != cfg.Model {
		t.Fatalf("Model mismatch: got %q want %q", parsed.Model, cfg.Model)
	}
}

func TestWizardUsesDefaultsOnEmptyInput(t *testing.T) {
	current := Config{
		BaseURL: "https://base",
		APIKey:  "secret",
		Model:   "model-a",
	}
	in := bytes.NewBufferString("\n\n\n")
	out := &bytes.Buffer{}

	cfg, err := Wizard(in, out, current)
	if err != nil {
		t.Fatalf("wizard error: %v", err)
	}
	if cfg.BaseURL != current.BaseURL || cfg.APIKey != current.APIKey || cfg.Model != current.Model {
		t.Fatalf("defaults not preserved: got %+v want %+v", cfg, current)
	}
}

func TestWizardOverridesValues(t *testing.T) {
	current := Config{
		BaseURL: "https://base",
		APIKey:  "secret",
		Model:   "model-a",
	}
	in := bytes.NewBufferString("https://new\nnewkey\nnewmodel\n")
	out := &bytes.Buffer{}

	cfg, err := Wizard(in, out, current)
	if err != nil {
		t.Fatalf("wizard error: %v", err)
	}
	if cfg.BaseURL != "https://new" {
		t.Fatalf("BaseURL not overridden: %q", cfg.BaseURL)
	}
	if cfg.APIKey != "newkey" {
		t.Fatalf("APIKey not overridden: %q", cfg.APIKey)
	}
	if cfg.Model != "newmodel" {
		t.Fatalf("Model not overridden: %q", cfg.Model)
	}
}

func TestLoadMergedUsesEnvOverrides(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.toml")
	cfg := Config{
		BaseURL: "https://file",
		APIKey:  "file-key",
		Model:   "file-model",
	}
	if err := Save(path, cfg); err != nil {
		t.Fatalf("save error: %v", err)
	}

	t.Setenv("AIP_BASE_URL", "https://env")
	t.Setenv("AIP_API_KEY", "env-key")
	t.Setenv("AIP_MODEL", "env-model")

	merged, err := LoadMerged(path)
	if err != nil {
		t.Fatalf("LoadMerged error: %v", err)
	}
	if merged.BaseURL != "https://env" {
		t.Fatalf("BaseURL override failed: %q", merged.BaseURL)
	}
	if merged.APIKey != "env-key" {
		t.Fatalf("APIKey override failed: %q", merged.APIKey)
	}
	if merged.Model != "env-model" {
		t.Fatalf("Model override failed: %q", merged.Model)
	}
}
