package summary

import (
	"errors"
	"fmt"
	"os"
	"strings"
)

type Prompt struct {
	User   string
	System string
}

func LoadPrompt(value string) (string, error) {
	if value == "" {
		return "", errors.New("prompt is required")
	}
	if strings.HasPrefix(value, "@") {
		path := strings.TrimPrefix(value, "@")
		data, err := os.ReadFile(path)
		if err != nil {
			return "", err
		}
		return strings.TrimRight(string(data), "\r\n"), nil
	}
	return value, nil
}

func BuildSystemPrompt() string {
	return strings.TrimSpace(`
You are a precise and concise assistant for log/text processing.
- Follow the user's instruction exactly.
- Use only the provided input; do not invent.
- If critical data is missing, say so briefly.
- Keep the output concise and directly usable.
- Output plain text unless the user explicitly requests a structured format.
`)
}

func BuildUserPrompt(userPrompt string, input string) string {
	return fmt.Sprintf("%s\n\nINPUT:\n%s", userPrompt, input)
}
