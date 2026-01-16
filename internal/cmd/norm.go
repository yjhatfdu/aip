package cmd

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/yjhatfdu/aip/internal/i18n"
	"github.com/yjhatfdu/aip/internal/norm"
)

func newNormCommand(lang i18n.Lang) *cobra.Command {
	var (
		profile   string
		rulesPath string
		emit      string
		bucket    string
	)

	cmd := &cobra.Command{
		Use:   "norm [file]",
		Short: i18n.T(lang, "cmd.norm.short"),
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			n, err := norm.New(profile, rulesPath, bucket)
			if err != nil {
				return err
			}

			var (
				reader  io.Reader = cmd.InOrStdin()
				srcFile           = ""
			)
			if len(args) == 1 {
				file, err := os.Open(args[0])
				if err != nil {
					return err
				}
				defer file.Close()
				reader = file
				srcFile = args[0]
			}

			if emit == "" {
				emit = "jsonl"
			}

			scanner := bufio.NewScanner(reader)
			scanner.Buffer(make([]byte, 0, 64*1024), 4*1024*1024)
			line := 0
			out := cmd.OutOrStdout()
			enc := json.NewEncoder(out)
			enc.SetEscapeHTML(false)

			for scanner.Scan() {
				line++
				record := n.Normalize(scanner.Text(), norm.Source{
					File: srcFile,
					Line: line,
				})
				switch emit {
				case "sig":
					if _, err := fmt.Fprintln(out, record.Sig); err != nil {
						return err
					}
				case "jsonl":
					if err := enc.Encode(record); err != nil {
						return err
					}
				case "tsv":
					sig := sanitizeTSV(record.Sig)
					ts := sanitizeTSV(record.TS)
					raw := sanitizeTSV(record.Raw)
					if _, err := fmt.Fprintf(out, "%s\t%s\t%s\n", sig, ts, raw); err != nil {
						return err
					}
				default:
					return fmt.Errorf("unknown emit: %s", emit)
				}
			}
			if err := scanner.Err(); err != nil && !errors.Is(err, io.EOF) {
				return err
			}
			return nil
		},
	}

	cmd.Flags().StringVar(&profile, "profile", "generic", "norm profile generic|postgres|kernel")
	cmd.Flags().StringVar(&rulesPath, "rules", "", "rules file path (YAML)")
	cmd.Flags().StringVar(&emit, "emit", "jsonl", "emit: sig|jsonl|tsv")
	cmd.Flags().StringVar(&bucket, "bucket", "", "bucket duration (e.g. 1m, 1h)")
	return cmd
}

func sanitizeTSV(value string) string {
	return strings.ReplaceAll(value, "\t", " ")
}
