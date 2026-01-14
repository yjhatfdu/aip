package config

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

type Config struct {
	BaseURL string
	APIKey  string
	Model   string
}

func LoadMerged(path string) (Config, error) {
	var cfg Config
	if loaded, err := Load(path); err == nil {
		cfg = loaded
	} else if !os.IsNotExist(err) {
		return Config{}, err
	}
	applyEnv(&cfg)
	return cfg, nil
}

func Merge(base Config, overrides Config) Config {
	out := base
	if overrides.BaseURL != "" {
		out.BaseURL = overrides.BaseURL
	}
	if overrides.APIKey != "" {
		out.APIKey = overrides.APIKey
	}
	if overrides.Model != "" {
		out.Model = overrides.Model
	}
	return out
}

func applyEnv(cfg *Config) {
	if val := pickEnv("AIP_BASE_URL", "OPENAI_BASE_URL"); val != "" {
		cfg.BaseURL = val
	}
	if val := pickEnv("AIP_API_KEY", "OPENAI_API_KEY"); val != "" {
		cfg.APIKey = val
	}
	if val := pickEnv("AIP_MODEL"); val != "" {
		cfg.Model = val
	}
}

func pickEnv(keys ...string) string {
	for _, key := range keys {
		if val := os.Getenv(key); val != "" {
			return val
		}
	}
	return ""
}

func GetByKey(cfg Config, key string) (string, bool) {
	switch key {
	case "base_url":
		return cfg.BaseURL, cfg.BaseURL != ""
	case "api_key":
		return cfg.APIKey, cfg.APIKey != ""
	case "model":
		return cfg.Model, cfg.Model != ""
	default:
		return "", false
	}
}

func SetByKey(cfg *Config, key, value string) bool {
	switch key {
	case "base_url":
		cfg.BaseURL = value
	case "api_key":
		cfg.APIKey = value
	case "model":
		cfg.Model = value
	default:
		return false
	}
	return true
}

func DefaultPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".aip", "config.toml"), nil
}

func EnsureDir(path string) error {
	dir := filepath.Dir(path)
	if dir == "." || dir == "" {
		return errors.New("invalid config directory")
	}
	return os.MkdirAll(dir, 0o755)
}

func Load(path string) (Config, error) {
	f, err := os.Open(path)
	if err != nil {
		return Config{}, err
	}
	defer f.Close()
	return parseTOML(f), nil
}

func Save(path string, cfg Config) error {
	if err := EnsureDir(path); err != nil {
		return err
	}
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.WriteString(renderTOML(cfg))
	return err
}

func Wizard(r io.Reader, w io.Writer, current Config) (Config, error) {
	reader := bufio.NewReader(r)
	cfg := current

	baseURL, err := prompt(reader, w, "AIP_BASE_URL", current.BaseURL)
	if err != nil {
		return Config{}, err
	}
	apiKey, err := prompt(reader, w, "AIP_API_KEY", current.APIKey)
	if err != nil {
		return Config{}, err
	}
	model, err := prompt(reader, w, "AIP_MODEL", current.Model)
	if err != nil {
		return Config{}, err
	}

	cfg.BaseURL = baseURL
	cfg.APIKey = apiKey
	cfg.Model = model
	return cfg, nil
}

func prompt(r *bufio.Reader, w io.Writer, name, current string) (string, error) {
	if current != "" {
		fmt.Fprintf(w, "%s [%s]: ", name, current)
	} else {
		fmt.Fprintf(w, "%s: ", name)
	}
	line, err := r.ReadString('\n')
	if err != nil && !errors.Is(err, io.EOF) {
		return "", err
	}
	line = strings.TrimSpace(line)
	if line == "" {
		return current, nil
	}
	return line, nil
}

func renderTOML(cfg Config) string {
	var b strings.Builder
	writeKV(&b, "base_url", cfg.BaseURL)
	writeKV(&b, "api_key", cfg.APIKey)
	writeKV(&b, "model", cfg.Model)
	return b.String()
}

func writeKV(b *strings.Builder, key, val string) {
	if val == "" {
		return
	}
	b.WriteString(key)
	b.WriteString(" = \"")
	b.WriteString(escape(val))
	b.WriteString("\"\n")
}

func escape(v string) string {
	v = strings.ReplaceAll(v, "\\", "\\\\")
	v = strings.ReplaceAll(v, "\"", "\\\"")
	return v
}

func parseTOML(r io.Reader) Config {
	var cfg Config
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.TrimSpace(parts[0])
		val := strings.TrimSpace(parts[1])
		val = strings.Trim(val, "\"")
		switch key {
		case "base_url":
			cfg.BaseURL = val
		case "api_key":
			cfg.APIKey = val
		case "model":
			cfg.Model = val
		}
	}
	return cfg
}
