package simbuildcore

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestPublicReadmeAndParser_DoNotLeakSourceDomain(t *testing.T) {
	domainToken := strings.Join([]string{"mhwilds", ".", "wiki", "-db", ".", "com"}, "")
	simPathToken := strings.Join([]string{"wiki", "-db", ".", "com", "/sim/"}, "")
	bannedTokens := []string{
		domainToken,
		simPathToken,
	}
	files := []string{
		filepath.Clean("../../README.md"),
		"url_parser.go",
	}

	for _, path := range files {
		raw, err := os.ReadFile(path)
		if err != nil {
			t.Fatalf("read %s: %v", path, err)
		}
		content := strings.ToLower(string(raw))
		for _, token := range bannedTokens {
			if strings.Contains(content, token) {
				t.Fatalf("public file %s contains banned token %q", path, token)
			}
		}
	}
}
