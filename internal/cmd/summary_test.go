package cmd

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"aip/internal/llm"
)

func TestSummaryCommandText(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/event-stream")
		flusher, ok := w.(http.Flusher)
		if !ok {
			t.Fatal("no flusher")
		}
		_, _ = w.Write([]byte("data: {\"choices\":[{\"delta\":{\"content\":\"summary \"}}]}\n\n"))
		flusher.Flush()
		_, _ = w.Write([]byte("data: {\"choices\":[{\"delta\":{\"content\":\"output\"}}]}\n\n"))
		flusher.Flush()
		_, _ = w.Write([]byte("data: [DONE]\n\n"))
	}))
	t.Cleanup(server.Close)

	root := newRoot()
	root.SetArgs([]string{
		"summary",
		"summarize",
		"--base-url", server.URL,
		"--api-key", "key",
		"--model", "model",
	})
	root.SetIn(strings.NewReader("hello\n"))
	out := &bytes.Buffer{}
	root.SetOut(out)
	root.SetErr(&bytes.Buffer{})

	if err := root.Execute(); err != nil {
		t.Fatalf("summary error: %v", err)
	}
	if strings.TrimSpace(out.String()) != "summary output" {
		t.Fatalf("unexpected output: %q", out.String())
	}
}

func TestSummaryCommandJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := llm.ChatResponse{
			Model: "model",
			Usage: llm.Usage{PromptTokens: 1, CompletionTokens: 2, TotalTokens: 3},
			Choices: []struct {
				Message llm.ChatMessage `json:"message"`
			}{
				{Message: llm.ChatMessage{Role: "assistant", Content: "ok"}},
			},
		}
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			t.Fatalf("encode: %v", err)
		}
	}))
	t.Cleanup(server.Close)

	root := newRoot()
	root.SetArgs([]string{
		"summary",
		"summarize",
		"--format", "json",
		"--stream=false",
		"--base-url", server.URL,
		"--api-key", "key",
		"--model", "model",
	})
	root.SetIn(strings.NewReader("hello\n"))
	out := &bytes.Buffer{}
	root.SetOut(out)
	root.SetErr(&bytes.Buffer{})

	if err := root.Execute(); err != nil {
		t.Fatalf("summary error: %v", err)
	}
	if !strings.Contains(out.String(), `"output":"ok"`) {
		t.Fatalf("unexpected json output: %q", out.String())
	}
}
