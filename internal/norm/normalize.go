package norm

import (
	"fmt"
	"strings"
	"time"
)

type Normalizer struct {
	rules    []compiledRule
	preserve map[string]struct{}
	bucket   time.Duration
}

func New(profile string, ruleFilePath string, bucket string) (*Normalizer, error) {
	var extra *RuleFile
	if ruleFilePath != "" {
		rf, err := LoadRuleFile(ruleFilePath)
		if err != nil {
			return nil, err
		}
		extra = &rf
	}
	rules, preserve, err := buildRules(profile, extra)
	if err != nil {
		return nil, err
	}
	compiled, err := compileRules(rules)
	if err != nil {
		return nil, err
	}
	var bucketDur time.Duration
	if bucket != "" {
		if d, err := time.ParseDuration(bucket); err == nil {
			bucketDur = d
		} else {
			return nil, fmt.Errorf("invalid bucket duration: %s", bucket)
		}
	}
	return &Normalizer{
		rules:    compiled,
		preserve: preserve,
		bucket:   bucketDur,
	}, nil
}

func (n *Normalizer) Normalize(line string, src Source) Record {
	raw := strings.TrimRight(line, "\r\n")
	sig := raw
	vars := map[string][]string{}
	ts := ""

	for _, rule := range n.rules {
		if !rule.re.MatchString(sig) {
			continue
		}
		sig = rule.re.ReplaceAllStringFunc(sig, func(m string) string {
			if _, ok := n.preserve[m]; ok {
				return m
			}
			if rule.isTS && ts == "" {
				ts = m
			}
			vars[rule.name] = append(vars[rule.name], m)
			return rule.replace
		})
	}

	record := Record{
		Raw:  raw,
		Sig:  sig,
		TS:   ts,
		Vars: vars,
		Src:  src,
	}

	if n.bucket > 0 && ts != "" {
		if parsed, err := parseTime(ts); err == nil {
			record.Bucket = parsed.Truncate(n.bucket).Format(time.RFC3339)
		}
	}
	if len(record.Vars) == 0 {
		record.Vars = nil
	}
	return record
}

func parseTime(value string) (time.Time, error) {
	layouts := []string{
		time.RFC3339Nano,
		time.RFC3339,
		"2006-01-02 15:04:05",
		"2006-01-02 15:04:05.000",
	}
	for _, layout := range layouts {
		if parsed, err := time.Parse(layout, value); err == nil {
			return parsed, nil
		}
	}
	return time.Time{}, fmt.Errorf("unsupported time: %s", value)
}
