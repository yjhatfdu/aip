package llm

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestComplete(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/chat/completions" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		if r.Header.Get("Authorization") == "" {
			t.Fatal("missing auth header")
		}
		var req ChatRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("decode: %v", err)
		}
		if req.Model == "" {
			t.Fatal("missing model")
		}
		resp := ChatResponse{
			Model: req.Model,
			Usage: Usage{PromptTokens: 1, CompletionTokens: 2, TotalTokens: 3},
			Choices: []struct {
				Message ChatMessage `json:"message"`
			}{
				{Message: ChatMessage{Role: "assistant", Content: "ok"}},
			},
		}
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			t.Fatalf("encode: %v", err)
		}
	}))
	t.Cleanup(server.Close)

	client := Client{
		BaseURL: server.URL,
		APIKey:  "test",
		Model:   "test-model",
	}
	resp, err := client.Complete(context.Background(), ChatRequest{
		Model: client.Model,
		Messages: []ChatMessage{
			{Role: "user", Content: "hi"},
		},
	})
	if err != nil {
		t.Fatalf("Complete error: %v", err)
	}
	if resp.Choices[0].Message.Content != "ok" {
		t.Fatalf("unexpected response: %+v", resp)
	}
}

func TestStreamComplete(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/event-stream")
		flusher, ok := w.(http.Flusher)
		if !ok {
			t.Fatal("no flusher")
		}
		_, _ = w.Write([]byte("data: {\"choices\":[{\"delta\":{\"content\":\"he\"}}]}\n\n"))
		flusher.Flush()
		_, _ = w.Write([]byte("data: {\"choices\":[{\"delta\":{\"content\":\"llo\"}}]}\n\n"))
		flusher.Flush()
		_, _ = w.Write([]byte("data: [DONE]\n\n"))
	}))
	t.Cleanup(server.Close)

	var out strings.Builder
	client := Client{
		BaseURL: server.URL,
		APIKey:  "test",
		Model:   "test-model",
	}
	err := client.StreamComplete(context.Background(), ChatRequest{
		Model: client.Model,
		Messages: []ChatMessage{
			{Role: "user", Content: "hi"},
		},
	}, func(delta string) error {
		out.WriteString(delta)
		return nil
	})
	if err != nil {
		t.Fatalf("StreamComplete error: %v", err)
	}
	if out.String() != "hello" {
		t.Fatalf("unexpected output: %q", out.String())
	}
}
