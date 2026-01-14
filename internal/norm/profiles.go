package norm

import (
	"embed"
	"fmt"
)

//go:embed profiles/*.yaml
var profilesFS embed.FS

var profileFiles = map[string]string{
	"postgres": "profiles/postgres.yaml",
	"kernel":   "profiles/kernel.yaml",
}

func loadProfile(profile string) (*RuleFile, error) {
	path, ok := profileFiles[profile]
	if !ok {
		return nil, fmt.Errorf("unknown profile: %s", profile)
	}
	data, err := profilesFS.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var rf RuleFile
	if err := yamlUnmarshal(data, &rf); err != nil {
		return nil, err
	}
	return &rf, nil
}
