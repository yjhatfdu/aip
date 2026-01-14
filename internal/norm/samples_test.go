package norm

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestNormalizeSamples(t *testing.T) {
	n, err := New("generic", "", "")
	if err != nil {
		t.Fatalf("new normalizer: %v", err)
	}

	t.Run("postgresql", func(t *testing.T) {
		lines := loadSample(t, "postgresql.log")
		if len(lines) == 0 {
			t.Fatal("empty postgresql sample")
		}
		var withTS int
		for _, line := range lines {
			rec := n.Normalize(line, Source{})
			if strings.Contains(rec.Sig, "<ts>") {
				withTS++
			}
		}
		if withTS == 0 {
			t.Fatal("expected at least one timestamp normalization in postgresql sample")
		}
	})

	t.Run("linux", func(t *testing.T) {
		lines := loadSample(t, "linux.log")
		if len(lines) == 0 {
			t.Fatal("empty linux sample")
		}
		var withIP int
		for _, line := range lines {
			rec := n.Normalize(line, Source{})
			if strings.Contains(rec.Sig, "<ip>") {
				withIP++
			}
		}
		if withIP == 0 {
			t.Fatal("expected at least one ip normalization in linux sample")
		}
	})
}

func loadSample(t *testing.T, name string) []string {
	t.Helper()
	path := filepath.Join("testdata", name)
	f, err := os.Open(path)
	if err != nil {
		t.Fatalf("open sample %s: %v", name, err)
	}
	defer f.Close()

	var lines []string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		t.Fatalf("scan sample %s: %v", name, err)
	}
	return lines
}
