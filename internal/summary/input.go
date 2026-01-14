package summary

import (
	"errors"
	"io"
	"strings"
)

type InputOptions struct {
	MaxChars    int
	IncludeHead int
	IncludeTail int
}

func ReadInput(r io.Reader, opts InputOptions) (string, error) {
	if opts.MaxChars <= 0 {
		opts.MaxChars = 40000
	}
	data, err := io.ReadAll(r)
	if err != nil {
		return "", err
	}
	text := string(data)
	if opts.IncludeHead > 0 || opts.IncludeTail > 0 {
		head := text
		tail := ""
		if opts.IncludeHead > 0 && len(text) > opts.IncludeHead {
			head = text[:opts.IncludeHead]
		}
		if opts.IncludeTail > 0 && len(text) > opts.IncludeTail {
			tail = text[len(text)-opts.IncludeTail:]
		}
		if tail != "" {
			text = head + "\n...\n" + tail
		} else {
			text = head
		}
	}
	if len(text) > opts.MaxChars {
		text = text[:opts.MaxChars]
	}
	if strings.TrimSpace(text) == "" {
		return "", errors.New("empty input")
	}
	return text, nil
}
