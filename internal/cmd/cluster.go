package cmd

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"aip/internal/cluster"
	"aip/internal/i18n"
	"github.com/spf13/cobra"
)

func newClusterCommand(lang i18n.Lang) *cobra.Command {
	var (
		algo       string
		field      string
		threshold  int
		bands      int
		bandBits   int
		minCluster int
		samples    int
		timeField  string
		format     string
	)

	cmd := &cobra.Command{
		Use:   "cluster [file]",
		Short: i18n.T(lang, "cmd.cluster.short"),
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if algo == "" {
				algo = "simhash"
			}
			if algo != "simhash" {
				return fmt.Errorf("unsupported algo: %s", algo)
			}
			if field == "" {
				field = "sig"
			}
			if timeField == "" {
				timeField = "ts"
			}
			if format == "" {
				format = "jsonl"
			}
			if bands == 0 {
				bands = 8
			}
			if bandBits == 0 {
				bandBits = 8
			}
			if minCluster == 0 {
				minCluster = 2
			}
			if samples == 0 {
				samples = 2
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

			infos, err := readClusterInput(reader, field, timeField)
			if err != nil {
				if srcFile != "" {
					return fmt.Errorf("%s: %w", srcFile, err)
				}
				return err
			}
			clusters, err := cluster.ClusterSigs(infos, cluster.Params{
				Threshold:  threshold,
				Bands:      bands,
				BandBits:   bandBits,
				MinCluster: minCluster,
				Samples:    samples,
			})
			if err != nil {
				return err
			}

			out := cmd.OutOrStdout()
			switch format {
			case "jsonl":
				enc := json.NewEncoder(out)
				enc.SetEscapeHTML(false)
				for _, c := range clusters {
					if err := enc.Encode(c); err != nil {
						return err
					}
				}
			case "json":
				payload := map[string]any{"clusters": clusters}
				enc := json.NewEncoder(out)
				enc.SetEscapeHTML(false)
				return enc.Encode(payload)
			case "text":
				for _, c := range clusters {
					if _, err := fmt.Fprintf(out, "%d\t%s\n", c.Count, c.Repr); err != nil {
						return err
					}
				}
			case "sample":
				for _, c := range clusters {
					sample := ""
					if len(c.Samples) > 0 {
						sample = c.Samples[0].Raw
					}
					if sample == "" {
						sample = c.Repr
					}
					if _, err := fmt.Fprintln(out, sample); err != nil {
						return err
					}
				}
			default:
				return fmt.Errorf("unknown format: %s", format)
			}
			return nil
		},
	}

	cmd.Flags().StringVar(&algo, "algo", "simhash", "cluster algorithm: simhash|minhash|embed")
	cmd.Flags().StringVar(&field, "field", "sig", "input field to cluster")
	cmd.Flags().IntVar(&threshold, "threshold", 4, "simhash hamming distance threshold")
	cmd.Flags().IntVar(&bands, "bands", 8, "LSH bands")
	cmd.Flags().IntVar(&bandBits, "band-bits", 8, "LSH band bits")
	cmd.Flags().IntVar(&minCluster, "min-cluster", 2, "minimum cluster size")
	cmd.Flags().IntVar(&samples, "samples", 2, "samples per cluster")
	cmd.Flags().StringVar(&timeField, "time-field", "ts", "time field name")
	cmd.Flags().StringVar(&format, "format", "jsonl", "format: jsonl|json|text|sample")
	return cmd
}

type inputRecord struct {
	Sig string
	Raw string
	TS  string
}

type sigAgg struct {
	cluster.SigInfo
}

func readClusterInput(r io.Reader, field, timeField string) ([]cluster.SigInfo, error) {
	scanner := bufio.NewScanner(r)
	scanner.Buffer(make([]byte, 0, 64*1024), 4*1024*1024)

	const (
		inputUnknown = iota
		inputJSONL
		inputText
	)
	mode := inputUnknown

	sigs := map[string]*sigAgg{}
	line := 0
	for scanner.Scan() {
		line++
		rawLine := strings.TrimSpace(scanner.Text())
		if rawLine == "" {
			continue
		}
		if mode == inputUnknown {
			if obj, ok := parseJSONObject(rawLine); ok {
				mode = inputJSONL
				if err := addJSONRecord(sigs, obj, field, timeField, line); err != nil {
					return nil, err
				}
				continue
			}
			mode = inputText
		}
		if mode == inputText {
			entry, ok := sigs[rawLine]
			if !ok {
				entry = &sigAgg{SigInfo: cluster.SigInfo{Sig: rawLine, Sample: rawLine}}
				sigs[rawLine] = entry
			}
			entry.Count++
			continue
		}
		obj, ok := parseJSONObject(rawLine)
		if !ok {
			return nil, fmt.Errorf("line %d: invalid json", line)
		}
		if err := addJSONRecord(sigs, obj, field, timeField, line); err != nil {
			return nil, err
		}
	}
	if err := scanner.Err(); err != nil && !errors.Is(err, io.EOF) {
		return nil, err
	}
	infos := make([]cluster.SigInfo, 0, len(sigs))
	for _, entry := range sigs {
		infos = append(infos, entry.SigInfo)
	}
	return infos, nil
}

func parseJSONObject(line string) (map[string]any, bool) {
	var obj map[string]any
	if err := json.Unmarshal([]byte(line), &obj); err != nil {
		return nil, false
	}
	return obj, true
}

func addJSONRecord(sigs map[string]*sigAgg, obj map[string]any, field, timeField string, line int) error {
	sigVal, ok := obj[field]
	if !ok {
		return fmt.Errorf("line %d: missing field %q", line, field)
	}
	sig := fmt.Sprint(sigVal)
	raw := ""
	if val, ok := obj["raw"]; ok {
		raw = fmt.Sprint(val)
	}
	ts := ""
	if timeField != "" {
		if val, ok := obj[timeField]; ok {
			ts = fmt.Sprint(val)
		}
	}
	entry, ok := sigs[sig]
	if !ok {
		entry = &sigAgg{SigInfo: cluster.SigInfo{Sig: sig, Sample: raw, SampleTS: ts}}
		sigs[sig] = entry
	}
	entry.Count++
	if entry.FirstTS == "" || (ts != "" && ts < entry.FirstTS) {
		entry.FirstTS = ts
	}
	if entry.LastTS == "" || (ts != "" && ts > entry.LastTS) {
		entry.LastTS = ts
	}
	return nil
}
