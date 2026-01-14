package i18n

import "os"

type Lang string

const (
	LangEN Lang = "en"
	LangZH Lang = "zh"
)

func Detect() Lang {
	if lang := pickEnv("AIP_LANG", "LC_ALL", "LC_MESSAGES", "LANG"); lang != "" {
		if hasZhPrefix(lang) {
			return LangZH
		}
	}
	return LangEN
}

func pickEnv(keys ...string) string {
	for _, key := range keys {
		if val := os.Getenv(key); val != "" {
			return val
		}
	}
	return ""
}

func hasZhPrefix(v string) bool {
	if len(v) < 2 {
		return false
	}
	if v[0] == 'z' && v[1] == 'h' {
		return true
	}
	if v[0] == 'Z' && v[1] == 'H' {
		return true
	}
	return false
}
