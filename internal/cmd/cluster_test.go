package cmd

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/yjhatfdu/aip/internal/cluster"
)

func TestClusterCommandJSONL(t *testing.T) {
	input := strings.Join([]string{
		`{"sig":"alpha","raw":"a","ts":"2024-01-01T00:00:01Z"}`,
		`{"sig":"alpha","raw":"a2","ts":"2024-01-01T00:00:02Z"}`,
		`{"sig":"beta","raw":"b","ts":"2024-01-01T00:00:03Z"}`,
	}, "\n") + "\n"

	root := newRoot()
	root.SetArgs([]string{"cluster", "--threshold", "64", "--min-cluster", "1", "--samples", "1", "--format", "jsonl"})
	root.SetIn(strings.NewReader(input))
	out := &bytes.Buffer{}
	root.SetOut(out)
	root.SetErr(&bytes.Buffer{})

	if err := root.Execute(); err != nil {
		t.Fatalf("cluster error: %v", err)
	}

	var got cluster.Cluster
	line := strings.TrimSpace(out.String())
	if err := json.Unmarshal([]byte(line), &got); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if got.Count != 3 {
		t.Fatalf("count mismatch: %d", got.Count)
	}
	if got.Repr != "alpha" {
		t.Fatalf("repr mismatch: %q", got.Repr)
	}
}

func TestClusterCommandSampleFormat(t *testing.T) {
	input := strings.Join([]string{
		`{"sig":"alpha","raw":"sample-a","ts":"2024-01-01T00:00:01Z"}`,
		`{"sig":"beta","raw":"sample-b","ts":"2024-01-01T00:00:02Z"}`,
	}, "\n") + "\n"

	root := newRoot()
	root.SetArgs([]string{"cluster", "--threshold", "64", "--min-cluster", "1", "--samples", "1", "--format", "sample"})
	root.SetIn(strings.NewReader(input))
	out := &bytes.Buffer{}
	root.SetOut(out)
	root.SetErr(&bytes.Buffer{})

	if err := root.Execute(); err != nil {
		t.Fatalf("cluster error: %v", err)
	}
	if !strings.Contains(out.String(), "sample-a") && !strings.Contains(out.String(), "sample-b") {
		t.Fatalf("expected sample output, got: %q", out.String())
	}
}

func TestClusterCommandAutoDetectText(t *testing.T) {
	input := "error one\nerror two\n"

	root := newRoot()
	root.SetArgs([]string{"cluster", "--threshold", "64", "--min-cluster", "1"})
	root.SetIn(strings.NewReader(input))
	out := &bytes.Buffer{}
	root.SetOut(out)
	root.SetErr(&bytes.Buffer{})

	if err := root.Execute(); err != nil {
		t.Fatalf("cluster error: %v", err)
	}
	if !strings.Contains(out.String(), `"count":2`) {
		t.Fatalf("expected clustered count, got: %q", out.String())
	}
}
