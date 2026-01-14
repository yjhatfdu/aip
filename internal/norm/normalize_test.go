package norm

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestNormalizeGeneric(t *testing.T) {
	n, err := New("generic", "", "")
	if err != nil {
		t.Fatalf("new normalizer: %v", err)
	}
	line := "2024-01-02 03:04:05 user 123 from 10.0.0.1 /var/log/app id=550e8400-e29b-41d4-a716-446655440000"
	rec := n.Normalize(line, Source{})

	wantSig := "<ts> user <number> from <ip> <path> id=<uuid>"
	if rec.Sig != wantSig {
		t.Fatalf("sig mismatch: got %q want %q", rec.Sig, wantSig)
	}
	if rec.TS != "2024-01-02 03:04:05" {
		t.Fatalf("ts mismatch: got %q", rec.TS)
	}
	if got := rec.Vars["uuid"]; len(got) != 1 || got[0] == "" {
		t.Fatalf("uuid vars missing: %#v", rec.Vars["uuid"])
	}
}

func TestBucketTruncate(t *testing.T) {
	n, err := New("generic", "", "1h")
	if err != nil {
		t.Fatalf("new normalizer: %v", err)
	}
	line := "2024-01-02T03:04:05Z job=sync"
	rec := n.Normalize(line, Source{})
	if rec.Bucket != "2024-01-02T03:00:00Z" {
		t.Fatalf("bucket mismatch: got %q", rec.Bucket)
	}
}

func TestRulesFileOverrides(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "rules.yaml")
	data := strings.Join([]string{
		"version: 1",
		"rules:",
		"  - name: pid",
		"    type: regex",
		"    pattern: 'pid=\\d+'",
		"    replace: 'pid=<pid>'",
	}, "\n")
	if err := os.WriteFile(path, []byte(data), 0o644); err != nil {
		t.Fatalf("write rules: %v", err)
	}

	n, err := New("generic", path, "")
	if err != nil {
		t.Fatalf("new normalizer: %v", err)
	}
	rec := n.Normalize("pid=123 status=ok", Source{})
	if !strings.Contains(rec.Sig, "pid=<pid>") {
		t.Fatalf("rule not applied: %q", rec.Sig)
	}
	if got := rec.Vars["pid"]; len(got) != 1 || got[0] != "pid=123" {
		t.Fatalf("vars mismatch: %#v", rec.Vars["pid"])
	}
}

func TestProfilePostgres(t *testing.T) {
	n, err := New("postgres", "", "")
	if err != nil {
		t.Fatalf("new normalizer: %v", err)
	}
	line := "2023-03-07 09:06:08 CET [130096] 6537ac6c.2397d5 FATAL 57P03: duration: 778.166 ms  statement: select 1"
	rec := n.Normalize(line, Source{})
	if !strings.Contains(rec.Sig, "<ts>") {
		t.Fatalf("expected ts placeholder: %q", rec.Sig)
	}
	if !strings.Contains(rec.Sig, "[<pid>]") {
		t.Fatalf("expected pid placeholder: %q", rec.Sig)
	}
	if !strings.Contains(rec.Sig, "duration: <ms> ms") {
		t.Fatalf("expected duration placeholder: %q", rec.Sig)
	}
	if !strings.Contains(rec.Sig, "<sess>") {
		t.Fatalf("expected session placeholder: %q", rec.Sig)
	}
	if !strings.Contains(rec.Sig, "<sqlstate>") {
		t.Fatalf("expected sqlstate placeholder: %q", rec.Sig)
	}
	if !strings.Contains(rec.Sig, "FATAL") {
		t.Fatalf("expected FATAL to remain: %q", rec.Sig)
	}
	if !strings.Contains(rec.Sig, "statement: <stmt>") {
		t.Fatalf("expected statement placeholder: %q", rec.Sig)
	}
}

func TestProfilePostgresDoesNotMaskPorts(t *testing.T) {
	n, err := New("postgres", "", "")
	if err != nil {
		t.Fatalf("new normalizer: %v", err)
	}
	line := "connection received: host=127.0.0.1 port=49134"
	rec := n.Normalize(line, Source{})
	if strings.Contains(rec.Sig, "<sqlstate>") {
		t.Fatalf("unexpected sqlstate match: %q", rec.Sig)
	}
	if !strings.Contains(rec.Sig, "port=<number>") {
		t.Fatalf("expected port to be normalized as number: %q", rec.Sig)
	}
}

func TestProfileKernel(t *testing.T) {
	n, err := New("kernel", "", "")
	if err != nil {
		t.Fatalf("new normalizer: %v", err)
	}
	line := "Jun 14 15:16:01 combo sshd(pam_unix)[19939]: authentication failure"
	rec := n.Normalize(line, Source{})
	if !strings.Contains(rec.Sig, "<ts>") {
		t.Fatalf("expected ts placeholder: %q", rec.Sig)
	}
	if !strings.Contains(rec.Sig, "[<pid>]") {
		t.Fatalf("expected pid placeholder: %q", rec.Sig)
	}
}
