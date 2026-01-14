package norm

type Source struct {
	File string `json:"file,omitempty"`
	Line int    `json:"line,omitempty"`
	Host string `json:"host,omitempty"`
}

type Record struct {
	Raw    string              `json:"raw"`
	Sig    string              `json:"sig"`
	TS     string              `json:"ts,omitempty"`
	Bucket string              `json:"bucket,omitempty"`
	Vars   map[string][]string `json:"vars,omitempty"`
	Src    Source              `json:"src,omitempty"`
}

type Rule struct {
	Name    string `yaml:"name"`
	Type    string `yaml:"type"`
	Class   string `yaml:"class,omitempty"`
	Pattern string `yaml:"pattern,omitempty"`
	Replace string `yaml:"replace,omitempty"`
}

type RuleFile struct {
	Version     int      `yaml:"version"`
	Description string   `yaml:"description,omitempty"`
	Rules       []Rule   `yaml:"rules"`
	Preserve    []string `yaml:"preserve"`
}
