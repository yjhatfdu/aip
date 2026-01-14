package i18n

import "testing"

func TestDetectLangDefaultsToEN(t *testing.T) {
	t.Setenv("AIP_LANG", "")
	t.Setenv("LC_ALL", "")
	t.Setenv("LC_MESSAGES", "")
	t.Setenv("LANG", "")

	if got := Detect(); got != LangEN {
		t.Fatalf("Detect() = %q, want %q", got, LangEN)
	}
}

func TestDetectLangChinese(t *testing.T) {
	t.Setenv("AIP_LANG", "zh_CN.UTF-8")
	t.Setenv("LC_ALL", "")
	t.Setenv("LC_MESSAGES", "")
	t.Setenv("LANG", "")

	if got := Detect(); got != LangZH {
		t.Fatalf("Detect() = %q, want %q", got, LangZH)
	}
}
