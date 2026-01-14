package cmd

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"time"

	"aip/internal/config"
	"aip/internal/i18n"
	"aip/internal/llm"
	"aip/internal/summary"
	"github.com/spf13/cobra"
)

type summaryResult struct {
	Output string    `json:"output"`
	Usage  llm.Usage `json:"usage,omitempty"`
	Model  string    `json:"model,omitempty"`
}

func newSummaryCommand(lang i18n.Lang) *cobra.Command {
	var (
		systemValue string
		format      string
		maxChars    int
		includeHead int
		includeTail int
		baseURL     string
		apiKey      string
		model       string
		stream      bool
	)

	cmd := &cobra.Command{
		Use:   "summary <prompt> [file]",
		Short: i18n.T(lang, "cmd.summary.short"),
		Args:  cobra.RangeArgs(1, 2),
		RunE: func(cmd *cobra.Command, args []string) error {
			userPrompt, err := summary.LoadPrompt(args[0])
			if err != nil {
				return err
			}
			systemPrompt := summary.BuildSystemPrompt()
			if systemValue != "" {
				systemPrompt, err = summary.LoadPrompt(systemValue)
				if err != nil {
					return err
				}
			}

			var (
				reader  io.Reader = cmd.InOrStdin()
				srcFile           = ""
			)
			if len(args) == 2 {
				file, err := os.Open(args[1])
				if err != nil {
					return err
				}
				defer file.Close()
				reader = file
				srcFile = args[1]
			}

			input, err := summary.ReadInput(reader, summary.InputOptions{
				MaxChars:    maxChars,
				IncludeHead: includeHead,
				IncludeTail: includeTail,
			})
			if err != nil {
				if srcFile != "" {
					return fmt.Errorf("%s: %w", srcFile, err)
				}
				return err
			}

			cfgPath, err := config.DefaultPath()
			if err != nil {
				return err
			}
			cfg, err := config.LoadMerged(cfgPath)
			if err != nil && !errors.Is(err, os.ErrNotExist) {
				return err
			}
			cfg = config.Merge(cfg, config.Config{
				BaseURL: baseURL,
				APIKey:  apiKey,
				Model:   model,
			})
			if cfg.BaseURL == "" || cfg.APIKey == "" || cfg.Model == "" {
				return errors.New("missing base_url/api_key/model (set env, config, or flags)")
			}

			client := llm.Client{
				BaseURL: cfg.BaseURL,
				APIKey:  cfg.APIKey,
				Model:   cfg.Model,
			}
			ctx, cancel := context.WithTimeout(cmd.Context(), 2*time.Minute)
			defer cancel()

			req := llm.ChatRequest{
				Model: cfg.Model,
				Messages: []llm.ChatMessage{
					{Role: "system", Content: systemPrompt},
					{Role: "user", Content: summary.BuildUserPrompt(userPrompt, input)},
				},
			}
			out := cmd.OutOrStdout()
			if format == "" {
				format = "text"
			}
			if format == "json" {
				stream = false
			}
			switch format {
			case "text":
				if stream {
					last := byte(0)
					wrote := false
					err := client.StreamComplete(ctx, req, func(delta string) error {
						if delta == "" {
							return nil
						}
						wrote = true
						last = delta[len(delta)-1]
						_, err := io.WriteString(out, delta)
						return err
					})
					if err != nil {
						return err
					}
					if wrote && last != '\n' {
						_, err := fmt.Fprintln(out)
						return err
					}
					return nil
				}
				resp, err := client.Complete(ctx, req)
				if err != nil {
					return err
				}
				if len(resp.Choices) == 0 {
					return errors.New("empty response")
				}
				_, err = fmt.Fprintln(out, resp.Choices[0].Message.Content)
				return err
			case "json":
				resp, err := client.Complete(ctx, req)
				if err != nil {
					return err
				}
				if len(resp.Choices) == 0 {
					return errors.New("empty response")
				}
				payload := summaryResult{
					Output: resp.Choices[0].Message.Content,
					Usage:  resp.Usage,
					Model:  resp.Model,
				}
				enc := json.NewEncoder(out)
				enc.SetEscapeHTML(false)
				return enc.Encode(payload)
			default:
				return fmt.Errorf("unknown format: %s", format)
			}
		},
	}

	cmd.Flags().StringVar(&systemValue, "system", "", "system prompt or @file")
	cmd.Flags().StringVar(&format, "format", "text", "format: text|json")
	cmd.Flags().IntVar(&maxChars, "max-chars", 40000, "max input chars")
	cmd.Flags().IntVar(&includeHead, "include-head", 0, "include head chars")
	cmd.Flags().IntVar(&includeTail, "include-tail", 0, "include tail chars")
	cmd.Flags().StringVar(&baseURL, "base-url", "", "LLM base URL")
	cmd.Flags().StringVar(&apiKey, "api-key", "", "LLM API key")
	cmd.Flags().StringVar(&model, "model", "", "LLM model")
	cmd.Flags().BoolVar(&stream, "stream", true, "stream output")
	return cmd
}
