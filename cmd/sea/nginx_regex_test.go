package main

import (
	"regexp"
	"strings"
	"testing"
)

// extractChatGPTRules compiles regexes from the generated nginx config.
func extractChatGPTRules(conf string) ([]*regexp.Regexp, []string, error) {
	var patterns []string
	lines := strings.Split(conf, "\n")
	inBlock := false
	re := regexp.MustCompile(`~\*\(\?i\)\^([^\\]+)\\b\s+chatgpt;`)
	for _, l := range lines {
		if strings.Contains(l, "# ChatGPT queries") {
			inBlock = true
			continue
		}
		if strings.Contains(l, "# explicit wiki keywords") {
			break
		}
		if inBlock {
			l = strings.TrimSpace(l)
			if m := re.FindStringSubmatch(l); len(m) == 2 {
				p := strings.TrimPrefix(strings.TrimSuffix(m[1], ")"), "(")
				patterns = append(patterns, p)
			}
		}
	}
	if len(patterns) == 0 {
		return nil, nil, nil
	}
	full := "(?i)^(" + strings.Join(patterns, "|") + ")\\b"
	rx, err := regexp.Compile(full)
	if err != nil {
		return nil, nil, err
	}
	var words []string
	for _, p := range patterns {
		words = append(words, strings.Split(p, "|")...)
	}
	return []*regexp.Regexp{rx}, words, nil
}

func matchAny(rs []*regexp.Regexp, s string) bool {
	for _, r := range rs {
		if r.MatchString(s) {
			return true
		}
	}
	return false
}

func TestChatGPTRouteAnchoring(t *testing.T) {
	cfg := Config{}
	out, err := generateNginx(cfg)
	if err != nil {
		t.Fatalf("failed to generate nginx: %v", err)
	}
	regs, words, err := extractChatGPTRules(out)
	if err != nil {
		t.Fatalf("failed to extract regex: %v", err)
	}
	if len(regs) == 0 {
		t.Fatal("no ChatGPT regex found in config")
	}
	for _, w := range words {
		start := w + " test"
		if !matchAny(regs, start) {
			t.Errorf("word %q not matched at start", w)
		}
		notStart := "hello " + w + " there"
		if matchAny(regs, notStart) {
			t.Errorf("word %q matched when not at start", w)
		}
	}
	if matchAny(regs, "completely unrelated") {
		t.Error("unrelated query incorrectly matched")
	}
}
