package norm

import (
	"errors"
	"fmt"
	"os"
	"regexp"
)

var tokenClasses = map[string]string{
	"uuid":   `\b[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}\b`,
	"hex":    `\b0x[0-9a-fA-F]+\b|\b[0-9a-fA-F]{16,}\b`,
	"number": `\b\d+\b`,
	"ip":     `\b(?:\d{1,3}\.){3}\d{1,3}\b`,
	"mac":    `\b(?:[0-9a-fA-F]{2}:){5}[0-9a-fA-F]{2}\b`,
	"path":   `(?:[A-Za-z]:\\(?:[^\\\s]+\\)*)|(?:/[A-Za-z0-9._\-/]+)`,
}

var defaultRules = []Rule{
	{
		Name:    "ts",
		Type:    "regex",
		Pattern: `\b\d{4}-\d{2}-\d{2}[T ]\d{2}:\d{2}:\d{2}(?:\.\d+)?(?:Z|[+-]\d{2}:\d{2})?\b`,
		Replace: "<ts>",
	},
	{Name: "uuid", Type: "token", Class: "uuid", Replace: "<uuid>"},
	{Name: "ip", Type: "token", Class: "ip", Replace: "<ip>"},
	{Name: "mac", Type: "token", Class: "mac", Replace: "<mac>"},
	{Name: "path", Type: "token", Class: "path", Replace: "<path>"},
	{Name: "hex", Type: "token", Class: "hex", Replace: "<hex>"},
	{Name: "number", Type: "token", Class: "number", Replace: "<number>"},
}

type compiledRule struct {
	name    string
	kind    string
	re      *regexp.Regexp
	replace string
	isTS    bool
}

func LoadRuleFile(path string) (RuleFile, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return RuleFile{}, err
	}
	var rf RuleFile
	if err := yamlUnmarshal(data, &rf); err != nil {
		return RuleFile{}, err
	}
	return rf, nil
}

func buildRules(profile string, extra *RuleFile) ([]Rule, map[string]struct{}, error) {
	var base []Rule
	preserve := map[string]struct{}{}
	switch profile {
	case "", "generic":
		base = append([]Rule(nil), defaultRules...)
	default:
		profileRules, err := loadProfile(profile)
		if err != nil {
			return nil, nil, err
		}
		base = append(base, profileRules.Rules...)
		for _, item := range profileRules.Preserve {
			if item != "" {
				preserve[item] = struct{}{}
			}
		}
		base = append(base, defaultRules...)
	}
	if extra != nil {
		for _, item := range extra.Preserve {
			if item != "" {
				preserve[item] = struct{}{}
			}
		}
		base = append(extra.Rules, base...)
	}
	return base, preserve, nil
}

func compileRules(rules []Rule) ([]compiledRule, error) {
	out := make([]compiledRule, 0, len(rules))
	for _, rule := range rules {
		if rule.Name == "" || rule.Type == "" {
			return nil, errors.New("rule requires name and type")
		}
		switch rule.Type {
		case "token":
			pattern, ok := tokenClasses[rule.Class]
			if !ok {
				return nil, fmt.Errorf("unknown token class: %s", rule.Class)
			}
			re, err := regexp.Compile(pattern)
			if err != nil {
				return nil, err
			}
			replace := rule.Replace
			if replace == "" {
				replace = "<" + rule.Class + ">"
			}
			out = append(out, compiledRule{
				name:    rule.Name,
				kind:    rule.Type,
				re:      re,
				replace: replace,
				isTS:    rule.Name == "ts",
			})
		case "regex":
			if rule.Pattern == "" {
				return nil, fmt.Errorf("regex rule %q missing pattern", rule.Name)
			}
			re, err := regexp.Compile(rule.Pattern)
			if err != nil {
				return nil, err
			}
			replace := rule.Replace
			if replace == "" {
				replace = "<" + rule.Name + ">"
			}
			out = append(out, compiledRule{
				name:    rule.Name,
				kind:    rule.Type,
				re:      re,
				replace: replace,
				isTS:    rule.Name == "ts",
			})
		default:
			return nil, fmt.Errorf("unknown rule type: %s", rule.Type)
		}
	}
	return out, nil
}
