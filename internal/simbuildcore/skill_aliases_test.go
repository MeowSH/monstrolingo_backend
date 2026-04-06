package simbuildcore

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
	"testing"
)

func withAliasSnapshotForTest(t *testing.T, entries []aliasSnapshotEntry) {
	t.Helper()

	dir := t.TempDir()
	path := filepath.Join(dir, "aliases.json")
	payload := aliasSnapshot{Entries: entries}
	raw, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("marshal alias snapshot: %v", err)
	}
	if err := os.WriteFile(path, raw, 0o600); err != nil {
		t.Fatalf("write alias snapshot: %v", err)
	}

	t.Setenv(aliasSnapshotPathEnv, path)
	resetAliasIndexForTest()
	t.Cleanup(resetAliasIndexForTest)
}

func resetAliasIndexForTest() {
	skillAliasIndexOnce = sync.Once{}
	skillAliasIndexInst = skillAliasIndex{}
}

func TestReadAliasSnapshotEntries_FromFile(t *testing.T) {
	withAliasSnapshotForTest(t, []aliasSnapshotEntry{
		{
			LanguageCode:    "en",
			Category:        "group",
			Alias:           "Guts (Tenacity)",
			AliasNormalized: normalizeSkillName("Guts (Tenacity)"),
			AliasLookupKeys: buildAliasLookupKeys("Guts (Tenacity)"),
			CanonicalName:   "Lord's Soul",
			ActivationLevel: 3,
		},
	})

	entries := readAliasSnapshotEntries()
	if len(entries) != 1 {
		t.Fatalf("expected 1 alias entry, got %d", len(entries))
	}
	if entries[0].CanonicalName != "Lord's Soul" {
		t.Fatalf("unexpected canonical name: %s", entries[0].CanonicalName)
	}
}

func TestReadAliasSnapshotEntries_MissingFile(t *testing.T) {
	t.Setenv(aliasSnapshotPathEnv, filepath.Join(t.TempDir(), "does_not_exist.json"))
	entries := readAliasSnapshotEntries()
	if entries != nil {
		t.Fatalf("expected nil entries when file is missing, got %d entries", len(entries))
	}
}
